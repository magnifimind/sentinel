"""LangGraph agent for infrastructure monitoring.

Flow: observe → decide (LLM) → check policy → act → report
"""

import logging
import time
from typing import Annotated, TypedDict

from langchain_core.messages import AIMessage, BaseMessage, HumanMessage, SystemMessage
from langchain_ollama import ChatOllama
from langgraph.graph import END, StateGraph
from langgraph.graph.message import add_messages
from langgraph.prebuilt import ToolNode

from sentinel_agent.clients.events import EventProducer
from sentinel_agent.clients.policy import PolicyClient
from sentinel_agent.config import settings
from sentinel_agent.tools import ALL_TOOLS

logger = logging.getLogger(__name__)

SYSTEM_PROMPT = """You are Sentinel, an infrastructure monitoring agent for the r740 server.

Your job is to run routine checks on the system and report findings. You have three tools:
- check_disk_usage: Check disk space on filesystem paths
- check_docker_containers: List Docker containers and their status
- scan_system_logs: Read recent log entries, optionally filtering by pattern

On each monitoring cycle:
1. Check disk usage on /
2. Check Docker container status
3. Scan syslog for recent errors

After gathering observations, summarize the system state. If you find issues that need
action (e.g., cleaning up disk space, restarting a container), propose the specific
command you would run. Do NOT execute destructive commands without stating them first.

Be concise. Report facts, not opinions."""


class AgentState(TypedDict):
    messages: Annotated[list[BaseMessage], add_messages]
    policy_result: str
    proposed_action: str
    cycle_complete: bool


def build_graph(
    policy_client: PolicyClient,
    event_producer: EventProducer,
) -> StateGraph:
    """Build the LangGraph agent graph."""

    llm = ChatOllama(
        base_url=settings.ollama_base_url,
        model=settings.model_name,
        temperature=0,
    ).bind_tools(ALL_TOOLS)

    tool_node = ToolNode(ALL_TOOLS)

    def observe(state: AgentState) -> dict:
        """Kick off the monitoring cycle by asking the LLM to run checks."""
        return {
            "messages": [
                SystemMessage(content=SYSTEM_PROMPT),
                HumanMessage(content="Run your monitoring checks now."),
            ],
        }

    def reason(state: AgentState) -> dict:
        """LLM reasons about observations and may call tools or propose actions."""
        start = time.time()
        response = llm.invoke(state["messages"])
        latency_ms = int((time.time() - start) * 1000)

        # Record decision event
        event_producer.send_decision({
            "decision_type": "monitoring_cycle",
            "action": "observe_and_reason",
            "policy_result": "n/a",
            "success": True,
            "llm_provider": "ollama",
            "model": settings.model_name,
            "latency_ms": latency_ms,
            "context": {"message_count": len(state["messages"])},
        })

        return {"messages": [response]}

    def should_use_tools(state: AgentState) -> str:
        """Route: if the LLM called tools, go to tool_node. Otherwise check for actions."""
        last = state["messages"][-1]
        if isinstance(last, AIMessage) and last.tool_calls:
            return "tools"
        return "check_actions"

    def check_actions(state: AgentState) -> dict:
        """Parse the LLM's final response for proposed actions and check policy."""
        last = state["messages"][-1]
        content = last.content if isinstance(last.content, str) else str(last.content)

        # Look for command-like proposals in the response
        action_keywords = ["rm ", "docker rm", "systemctl", "kubectl", "shutdown", "reboot"]
        proposed = ""
        for line in content.split("\n"):
            stripped = line.strip().lstrip("$ ").lstrip("`").rstrip("`")
            if any(kw in stripped.lower() for kw in action_keywords):
                proposed = stripped
                break

        if not proposed:
            # No action proposed — just an observation report
            event_producer.send_outcome({
                "action": "observation_only",
                "success": True,
                "summary": content[:500],
            })
            event_producer.flush()
            return {"cycle_complete": True, "proposed_action": "", "policy_result": "n/a"}

        # Check policy before acting
        result = policy_client.evaluate(proposed, {"source": "monitoring_cycle"})
        decision = result.get("decision", "block")

        event_producer.send_decision({
            "decision_type": "action_proposal",
            "action": proposed,
            "policy_result": decision,
            "success": decision == "approve",
            "llm_provider": "ollama",
            "model": settings.model_name,
            "latency_ms": 0,
            "context": {"reason": result.get("reason", "")},
        })

        if decision == "block":
            logger.warning("action BLOCKED by policy: %s — %s", proposed, result.get("reason"))
            event_producer.send_outcome({
                "action": proposed,
                "success": False,
                "policy_result": "blocked",
                "reason": result.get("reason", ""),
            })
        elif decision == "escalate":
            logger.info("action ESCALATED: %s — requires human approval", proposed)
            event_producer.send_outcome({
                "action": proposed,
                "success": False,
                "policy_result": "escalated",
                "reason": result.get("reason", ""),
            })
        else:
            logger.info("action APPROVED: %s", proposed)
            event_producer.send_action({"action": proposed, "policy_result": "approved"})
            event_producer.send_outcome({
                "action": proposed,
                "success": True,
                "policy_result": "approved",
            })

        event_producer.flush()
        return {
            "cycle_complete": True,
            "proposed_action": proposed,
            "policy_result": decision,
        }

    # Build the graph
    graph = StateGraph(AgentState)

    graph.add_node("observe", observe)
    graph.add_node("reason", reason)
    graph.add_node("tools", tool_node)
    graph.add_node("check_actions", check_actions)

    graph.set_entry_point("observe")
    graph.add_edge("observe", "reason")
    graph.add_conditional_edges("reason", should_use_tools, {"tools": "tools", "check_actions": "check_actions"})
    graph.add_edge("tools", "reason")  # after tool results, go back to LLM
    graph.add_edge("check_actions", END)

    return graph.compile()
