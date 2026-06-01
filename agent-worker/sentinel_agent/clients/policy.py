import logging

import httpx

from sentinel_agent.config import settings

logger = logging.getLogger(__name__)


class PolicyClient:
    """Client for the Sentinel Go API policy evaluation endpoint."""

    def __init__(self) -> None:
        self.base_url = settings.api_base_url.rstrip("/")
        self.client = httpx.Client(base_url=self.base_url, timeout=10.0)

    def evaluate(self, action: str, context: dict | None = None) -> dict:
        """Ask the policy engine whether an action is allowed.

        Returns dict with keys: decision ("approve"|"block"|"escalate"), reason, policy_id.
        """
        payload = {
            "agent_id": settings.agent_id,
            "action": action,
        }
        if context:
            payload["context"] = context

        try:
            resp = self.client.post("/api/v1/policy/evaluate", json=payload)
            resp.raise_for_status()
            result = resp.json()
            logger.info(
                "policy decision: %s for action '%s' (policy: %s)",
                result.get("decision"),
                action,
                result.get("policy_id", "default"),
            )
            return result
        except httpx.HTTPError as e:
            logger.error("policy evaluation failed for '%s': %s", action, e)
            # Fail closed — if we can't reach the policy engine, block the action
            return {"decision": "block", "reason": f"policy engine unreachable: {e}"}

    def close(self) -> None:
        self.client.close()
