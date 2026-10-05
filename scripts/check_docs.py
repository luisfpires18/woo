"""Validate local Markdown links and the task register. No external network required."""
from pathlib import Path
import re
import sys
from urllib.parse import unquote

ROOT = Path(__file__).resolve().parents[1]
errors = []
markdown = list(ROOT.glob("*.md")) + list((ROOT / "docs").rglob("*.md"))
for path in markdown:
    content = re.sub(r"```[\s\S]*?```", "", path.read_text(encoding="utf-8"))
    for target in re.findall(r"\]\(([^)]+)\)", content):
        if re.match(r"^[a-zA-Z][a-zA-Z0-9+.-]*:", target) or target.startswith("#"):
            continue
        target = unquote(target.split("#", 1)[0])
        if target and not (path.parent / target).exists():
            errors.append(f"{path.relative_to(ROOT)}: missing link {target}")

roadmap = (ROOT / "tasks_roadmap.md").read_text(encoding="utf-8")
rows = re.findall(r"^\| (\d{3}) \| ([A-Z0-9-]+) \| ([^|]+) \| ([A-Z]+) \| \[[^]]+\]\(([^)]+)\) \|$", roadmap, re.M)
if not rows:
    errors.append("No task rows found")
keys = set()
referenced = set()
active = 0
for expected, (number, key, title, status, target) in enumerate(rows, 1):
    if int(number) != expected:
        errors.append(f"Task order must be contiguous: expected {expected:03}, got {number}")
    if key in keys:
        errors.append(f"Duplicate task key {key}")
    keys.add(key)
    if status not in {"TODO", "WIP", "DONE"}:
        errors.append(f"{key}: invalid status {status}")
    active += status == "WIP"
    path = ROOT / target
    referenced.add(path.resolve())
    if not path.is_file():
        errors.append(f"{key}: missing detail {target}")
        continue
    text = path.read_text(encoding="utf-8")
    if not path.name.startswith(number + "-") or not text.startswith(f"# {number}: "):
        errors.append(f"{key}: detail number mismatch")
    if f"Task key: {key}\n" not in text or f"Status: {status}\n" not in text:
        errors.append(f"{key}: detail identity/status mismatch")
    match = re.search(r"^Dependencies: (.+)$", text, re.M)
    if not match:
        errors.append(f"{key}: dependencies missing")
if active > 1:
    errors.append("Only one WIP task is permitted")
for path in (ROOT / "docs/tasks").glob("[0-9][0-9][0-9]-*.md"):
    if path.resolve() not in referenced:
        errors.append(f"Unregistered detail: {path.relative_to(ROOT)}")
for number, key, title, status, target in rows:
    path = ROOT / target
    if path.is_file():
        dependency = re.search(r"^Dependencies: (.+)$", path.read_text(encoding="utf-8"), re.M)
        if dependency and dependency[1] != "None":
            for dep in dependency[1].split(", "):
                if dep not in keys:
                    errors.append(f"{key}: unknown dependency {dep}")
if errors:
    print("\n".join(errors))
    sys.exit(1)
print(f"Checked {len(markdown)} Markdown files and {len(rows)} tasks: OK")
