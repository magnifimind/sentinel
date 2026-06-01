import shutil

from langchain_core.tools import tool


@tool
def check_disk_usage(path: str = "/") -> str:
    """Check disk usage for a given filesystem path. Returns total, used, free space and percent used."""
    try:
        usage = shutil.disk_usage(path)
        total_gb = usage.total / (1024**3)
        used_gb = usage.used / (1024**3)
        free_gb = usage.free / (1024**3)
        pct = (usage.used / usage.total) * 100
        return (
            f"Disk usage for {path}:\n"
            f"  Total: {total_gb:.1f} GB\n"
            f"  Used:  {used_gb:.1f} GB ({pct:.1f}%)\n"
            f"  Free:  {free_gb:.1f} GB\n"
            f"  Status: {'CRITICAL' if pct > 90 else 'WARNING' if pct > 80 else 'OK'}"
        )
    except OSError as e:
        return f"Error checking disk at {path}: {e}"
