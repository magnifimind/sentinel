import subprocess

from langchain_core.tools import tool


@tool
def check_docker_containers(filter_status: str = "") -> str:
    """List Docker containers and their status. Optionally filter by status (running, exited, etc)."""
    try:
        cmd = ["docker", "ps", "-a", "--format", "{{.Names}}\t{{.Status}}\t{{.Image}}"]
        if filter_status:
            cmd.extend(["--filter", f"status={filter_status}"])
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=10)
        if result.returncode != 0:
            return f"Docker error: {result.stderr.strip()}"
        if not result.stdout.strip():
            return "No containers found."
        lines = result.stdout.strip().split("\n")
        output = f"Found {len(lines)} container(s):\n"
        for line in lines:
            parts = line.split("\t")
            if len(parts) == 3:
                name, status, image = parts
                output += f"  {name}: {status} ({image})\n"
            else:
                output += f"  {line}\n"
        return output
    except FileNotFoundError:
        return "Docker is not installed or not in PATH."
    except subprocess.TimeoutExpired:
        return "Docker command timed out."
