import subprocess

from langchain_core.tools import tool


@tool
def scan_system_logs(log_file: str = "/var/log/syslog", lines: int = 50, pattern: str = "") -> str:
    """Scan system log files for recent entries. Optionally grep for a pattern. Returns the last N lines."""
    try:
        if pattern:
            cmd = ["grep", "-i", pattern, log_file]
            result = subprocess.run(cmd, capture_output=True, text=True, timeout=10)
            if result.returncode == 1:
                return f"No matches for '{pattern}' in {log_file}"
            if result.returncode != 0:
                return f"Error reading {log_file}: {result.stderr.strip()}"
            matched = result.stdout.strip().split("\n")
            tail = matched[-lines:]
            return f"Found {len(matched)} matches for '{pattern}' in {log_file} (showing last {len(tail)}):\n" + "\n".join(tail)
        else:
            cmd = ["tail", "-n", str(lines), log_file]
            result = subprocess.run(cmd, capture_output=True, text=True, timeout=10)
            if result.returncode != 0:
                return f"Error reading {log_file}: {result.stderr.strip()}"
            return f"Last {lines} lines of {log_file}:\n{result.stdout.strip()}"
    except subprocess.TimeoutExpired:
        return f"Timed out reading {log_file}"
