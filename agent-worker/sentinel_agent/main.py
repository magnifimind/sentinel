"""Sentinel agent-worker entrypoint.

Runs the LangGraph monitoring agent on a configurable schedule.
"""

import logging
import signal
import sys
import time

from sentinel_agent.clients.events import EventProducer
from sentinel_agent.clients.policy import PolicyClient
from sentinel_agent.config import settings
from sentinel_agent.graph import build_graph

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(name)s %(message)s",
    datefmt="%Y-%m-%dT%H:%M:%S",
)
logger = logging.getLogger("sentinel_agent")


def main() -> None:
    logger.info(
        "sentinel agent-worker starting (agent_id=%s, interval=%ds, model=%s)",
        settings.agent_id,
        settings.check_interval_seconds,
        settings.model_name,
    )

    policy_client = PolicyClient()
    event_producer = EventProducer()
    agent = build_graph(policy_client, event_producer)

    shutdown = False

    def handle_signal(sig: int, _frame: object) -> None:
        nonlocal shutdown
        logger.info("received signal %d, shutting down", sig)
        shutdown = True

    signal.signal(signal.SIGINT, handle_signal)
    signal.signal(signal.SIGTERM, handle_signal)

    logger.info("agent-worker ready, entering monitoring loop")

    while not shutdown:
        try:
            logger.info("starting monitoring cycle")
            result = agent.invoke({
                "messages": [],
                "policy_result": "",
                "proposed_action": "",
                "cycle_complete": False,
            })
            logger.info(
                "monitoring cycle complete (policy_result=%s, action=%s)",
                result.get("policy_result", "n/a"),
                result.get("proposed_action", "none"),
            )
        except Exception:
            logger.exception("monitoring cycle failed")

        # Sleep in small increments so we can respond to signals
        for _ in range(settings.check_interval_seconds):
            if shutdown:
                break
            time.sleep(1)

    policy_client.close()
    event_producer.close()
    logger.info("agent-worker stopped")


if __name__ == "__main__":
    main()
