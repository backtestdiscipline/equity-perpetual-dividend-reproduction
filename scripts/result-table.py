#!/usr/bin/env python3
"""Render the sealed result table; no external inputs or calculations."""
import json
from pathlib import Path

p = Path(__file__).resolve().parents[1] / '.research-platform/runs/20261111T125325Z-0a66a12c1a04/artifacts/result.json'
r = json.loads(p.read_text(encoding='utf-8'))
print('| Contract | Ex-date | Dividend ($) | Special rate | API mark ($) | Absolute residual |')
print('|---|---|---:|---:|---:|---:|')
for e in r['events']:
    print(f"| {e['symbol']}USDT | {e['ex_date']} | {e['dividend']} | {e['rate']} | {e['mark']} | {e['paper_residual']} |")
