from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    # Sentinel API
    api_base_url: str = "http://sentinel-api:8080"
    agent_id: str = "infra-monitor"

    # Ollama / LLM
    ollama_base_url: str = "http://sentinel-model-server:11434"
    model_name: str = "llama3.1:8b-instruct-q4_0"

    # Kafka
    kafka_brokers: str = "sentinel-kafka:9092"

    # Monitoring schedule
    check_interval_seconds: int = 300  # 5 minutes

    # Target host to monitor (the r740 itself, from inside the cluster)
    monitor_host: str = "localhost"

    model_config = {"env_prefix": "SENTINEL_"}

    @property
    def kafka_broker_list(self) -> list[str]:
        return [b.strip() for b in self.kafka_brokers.split(",")]


settings = Settings()
