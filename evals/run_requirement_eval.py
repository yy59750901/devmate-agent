#!/usr/bin/env python3
"""Run a minimal rule-based evaluation for requirement analysis."""

from __future__ import annotations

import argparse
import json
import sys
import urllib.error
import urllib.request
from pathlib import Path
from typing import Any

DEFAULT_BASE_URL = "http://localhost:8080"
DEFAULT_SAMPLES = Path(__file__).with_name("requirement_samples.jsonl")
REQUIRED_ARRAY_FIELDS = ["apis", "tables", "risks", "test_cases", "questions"]


def load_samples(path: Path) -> list[dict[str, Any]]:
    samples: list[dict[str, Any]] = []
    with path.open("r", encoding="utf-8") as file:
        for line_no, line in enumerate(file, start=1):
            line = line.strip()
            if not line:
                continue
            try:
                sample = json.loads(line)
            except json.JSONDecodeError as exc:
                raise ValueError(f"invalid JSON at {path}:{line_no}: {exc}") from exc
            samples.append(sample)
    return samples


def call_requirement_api(base_url: str, sample: dict[str, Any], timeout: float) -> tuple[int, dict[str, Any]]:
    payload = {
        "requirement": sample["requirement"],
        "context": sample.get("context", ""),
    }
    data = json.dumps(payload, ensure_ascii=False).encode("utf-8")
    request = urllib.request.Request(
        f"{base_url.rstrip('/')}/api/analyze/requirement",
        data=data,
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    try:
        with urllib.request.urlopen(request, timeout=timeout) as response:
            body = response.read().decode("utf-8")
            return response.status, json.loads(body)
    except urllib.error.HTTPError as exc:
        body = exc.read().decode("utf-8")
        try:
            parsed = json.loads(body)
        except json.JSONDecodeError:
            parsed = {"raw_body": body}
        return exc.code, parsed


def collect_text(value: Any) -> str:
    if isinstance(value, str):
        return value
    if isinstance(value, list):
        return "\n".join(collect_text(item) for item in value)
    if isinstance(value, dict):
        return "\n".join(collect_text(item) for item in value.values())
    return ""


def evaluate_response(sample: dict[str, Any], status_code: int, response: dict[str, Any]) -> dict[str, Any]:
    failures: list[str] = []
    if status_code != 200:
        failures.append(f"http_status={status_code}")
    if response.get("status") != "succeeded":
        failures.append(f"task_status={response.get('status')}")

    output = response.get("output") or {}
    result = output.get("result") or {}
    llm = output.get("llm") or {}
    usage = llm.get("usage") or {}

    summary = result.get("summary", "")
    if not isinstance(summary, str) or not summary.strip():
        failures.append("summary_empty")

    for field in REQUIRED_ARRAY_FIELDS:
        if not isinstance(result.get(field), list):
            failures.append(f"{field}_not_array")

    if not isinstance(usage.get("total_tokens"), int):
        failures.append("total_tokens_missing")

    output_text = collect_text(result)
    expected_keywords = sample.get("expected_keywords", [])
    keyword_hits = [keyword for keyword in expected_keywords if keyword in output_text]
    if expected_keywords and not keyword_hits:
        failures.append("keyword_hits_empty")

    return {
        "id": sample.get("id", ""),
        "passed": not failures,
        "failures": failures,
        "keyword_hits": len(keyword_hits),
        "keyword_total": len(expected_keywords),
        "total_tokens": usage.get("total_tokens", 0),
        "latency_ms": llm.get("latency_ms", 0),
        "error": response.get("error"),
    }


def print_result(result: dict[str, Any]) -> None:
    status = "PASS" if result["passed"] else "FAIL"
    print(
        f"[{status}] {result['id']} "
        f"keyword_hits={result['keyword_hits']}/{result['keyword_total']} "
        f"total_tokens={result['total_tokens']} latency_ms={result['latency_ms']}"
    )
    if result["failures"]:
        print(f"  failures: {', '.join(result['failures'])}")
    if result.get("error"):
        print(f"  error: {json.dumps(result['error'], ensure_ascii=False)}")


def main() -> int:
    parser = argparse.ArgumentParser(description="Run requirement analysis minimal evaluation.")
    parser.add_argument("--base-url", default=DEFAULT_BASE_URL, help="DevMate backend base URL")
    parser.add_argument("--samples", type=Path, default=DEFAULT_SAMPLES, help="JSONL sample file")
    parser.add_argument("--timeout", type=float, default=90.0, help="HTTP timeout seconds per sample")
    args = parser.parse_args()

    samples = load_samples(args.samples)
    print(f"Loaded {len(samples)} samples from {args.samples}")

    results: list[dict[str, Any]] = []
    for sample in samples:
        try:
            status_code, response = call_requirement_api(args.base_url, sample, args.timeout)
            result = evaluate_response(sample, status_code, response)
        except Exception as exc:  # noqa: BLE001 - CLI should report all sample failures.
            result = {
                "id": sample.get("id", ""),
                "passed": False,
                "failures": ["request_exception"],
                "keyword_hits": 0,
                "keyword_total": len(sample.get("expected_keywords", [])),
                "total_tokens": 0,
                "latency_ms": 0,
                "error": str(exc),
            }
        results.append(result)
        print_result(result)

    total = len(results)
    passed = sum(1 for result in results if result["passed"])
    failed = total - passed
    pass_rate = (passed / total * 100) if total else 0.0
    total_tokens = sum(int(result.get("total_tokens") or 0) for result in results)

    print("\nSummary:")
    print(f"  total: {total}")
    print(f"  passed: {passed}")
    print(f"  failed: {failed}")
    print(f"  pass_rate: {pass_rate:.1f}%")
    print(f"  total_tokens: {total_tokens}")
    return 0 if failed == 0 else 1


if __name__ == "__main__":
    sys.exit(main())
