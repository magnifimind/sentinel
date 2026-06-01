import json
import logging
from datetime import datetime, timezone

from kafka import KafkaProducer

from sentinel_agent.config import settings

logger = logging.getLogger(__name__)


class EventProducer:
    """Produces agent events to Kafka topics."""

    def __init__(self) -> None:
        self.producer = KafkaProducer(
            bootstrap_servers=settings.kafka_broker_list,
            value_serializer=lambda v: json.dumps(v, default=str).encode("utf-8"),
            key_serializer=lambda k: k.encode("utf-8") if k else None,
            acks="all",
            retries=3,
        )

    def send_decision(self, decision: dict) -> None:
        """Send a decision event to agent.decisions topic."""
        event = self._wrap_event("decision", decision)
        self.producer.send("agent.decisions", key=settings.agent_id, value=event)
        logger.debug("sent decision event: %s", decision.get("action", "unknown"))

    def send_action(self, action: dict) -> None:
        """Send an action event to agent.actions topic."""
        event = self._wrap_event("action", action)
        self.producer.send("agent.actions", key=settings.agent_id, value=event)
        logger.debug("sent action event: %s", action.get("action", "unknown"))

    def send_outcome(self, outcome: dict) -> None:
        """Send an outcome event to agent.outcomes topic."""
        event = self._wrap_event("outcome", outcome)
        self.producer.send("agent.outcomes", key=settings.agent_id, value=event)
        logger.debug("sent outcome event")

    def flush(self) -> None:
        self.producer.flush(timeout=5)

    def close(self) -> None:
        self.producer.close(timeout=5)

    def _wrap_event(self, event_type: str, payload: dict) -> dict:
        return {
            "agent_id": settings.agent_id,
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "type": event_type,
            "payload": payload,
        }
