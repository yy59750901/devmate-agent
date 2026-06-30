#!/usr/bin/env python3
"""CLI demo for DevMate requirement analysis."""

from __future__ import annotations

import argparse
import json
import sys
import textwrap
import urllib.error
import urllib.request
from typing import Any

DEFAULT_BASE_URL = "http://localhost:8080"
DEFAULT_REQUIREMENT = "用户希望增加订单退款功能，支持部分退款、原路退回、退款失败重试，并记录操作审计日志。"
DEFAULT_CONTEXT = "当前系统已有订单和支付模块，退款依赖第三方支付通道，支付和退款结果通过异步回调通知。"


def call_requirement_api(base_url: str, requirement: str, context: str, timeout: float) -> tuple[int, dict[str, Any]]:
    payload = {
        "requirement": requirement,
        "context": context,
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


def print_section(title: str, value: Any) -> None:
    print(f"\n{title}:")
    if isinstance(value, list):
        if not value:
            print("  - 无")
            return
        for item in value:
            print(f"  - {item}")
        return
    text = str(value or "无")
    wrapped = textwrap.fill(text, width=88, subsequent_indent="  ")
    print(f"  {wrapped}")


def print_success(response: dict[str, Any]) -> None:
    output = response.get("output") or {}
    result = output.get("result") or {}
    llm = output.get("llm") or {}
    usage = llm.get("usage") or {}

    print(f"Task: {response.get('id')} status={response.get('status')}")
    print(f"Prompt version: {output.get('prompt_version', 'unknown')}")
    print_section("Summary", result.get("summary"))
    print_section("APIs", result.get("apis", []))
    print_section("Tables / data objects", result.get("tables", []))
    print_section("Risks", result.get("risks", []))
    print_section("Test cases", result.get("test_cases", []))
    print_section("Questions", result.get("questions", []))
    print("\nLLM:")
    print(f"  model: {llm.get('model', 'unknown')}")
    print(f"  finish_reason: {llm.get('finish_reason', 'unknown')}")
    print(f"  prompt_tokens: {usage.get('prompt_tokens', 0)}")
    print(f"  completion_tokens: {usage.get('completion_tokens', 0)}")
    print(f"  total_tokens: {usage.get('total_tokens', 0)}")
    print(f"  latency_ms: {llm.get('latency_ms', 0)}")


def print_failure(status_code: int, response: dict[str, Any]) -> None:
    print(f"Request failed: http_status={status_code}")
    error = response.get("error")
    if error:
        print(json.dumps(error, ensure_ascii=False, indent=2))
    else:
        print(json.dumps(response, ensure_ascii=False, indent=2))


def main() -> int:
    parser = argparse.ArgumentParser(description="Run a requirement analysis demo request.")
    parser.add_argument("--base-url", default=DEFAULT_BASE_URL, help="DevMate backend base URL")
    parser.add_argument("--requirement", default=DEFAULT_REQUIREMENT, help="Requirement text")
    parser.add_argument("--context", default=DEFAULT_CONTEXT, help="Optional business context")
    parser.add_argument("--timeout", type=float, default=600.0, help="HTTP timeout seconds")
    args = parser.parse_args()

    try:
        status_code, response = call_requirement_api(args.base_url, args.requirement, args.context, args.timeout)
    except Exception as exc:  # noqa: BLE001 - CLI should show connection errors directly.
        print(f"Request error: {exc}", file=sys.stderr)
        return 1

    if status_code == 200 and response.get("status") == "succeeded":
        print_success(response)
        return 0
    print_failure(status_code, response)
    return 1


if __name__ == "__main__":
    sys.exit(main())
