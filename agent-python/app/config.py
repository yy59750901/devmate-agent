from __future__ import annotations

from dataclasses import dataclass, field
from pathlib import Path
from typing import Any

import yaml


AGENT_ROOT = Path(__file__).resolve().parents[1]
PROJECT_ROOT = AGENT_ROOT.parent
DEFAULT_CONFIG_PATH = AGENT_ROOT / "config" / "config.local.yaml"


@dataclass(frozen=True)
class ModelConfig:
    provider: str
    base_url: str
    model: str
    api_key: str


@dataclass(frozen=True)
class RAGConfig:
    index_version: str = "rag-docs-v1"
    storage_dir: str = "storage/rag_index"
    document_paths: list[str] = field(default_factory=lambda: ["../docs", "../README.md"])


@dataclass(frozen=True)
class AppConfig:
    llm: ModelConfig
    embedding: ModelConfig
    rag: RAGConfig


def load_config(config_path: Path = DEFAULT_CONFIG_PATH) -> AppConfig:
    if not config_path.exists():
        raise FileNotFoundError(
            f"config file not found: {config_path}. Copy agent-python/config/config.example.yaml to config.local.yaml first."
        )

    raw = yaml.safe_load(config_path.read_text(encoding="utf-8")) or {}
    return AppConfig(
        llm=_load_model_config(raw, "llm"),
        embedding=_load_model_config(raw, "embedding"),
        rag=_load_rag_config(raw.get("rag") or {}),
    )


def resolve_agent_path(path: str) -> Path:
    candidate = Path(path)
    if candidate.is_absolute():
        return candidate
    return (AGENT_ROOT / candidate).resolve()


def project_relative_path(path: Path) -> str:
    try:
        return str(path.resolve().relative_to(PROJECT_ROOT))
    except ValueError:
        return str(path.resolve())


def _load_model_config(raw: dict[str, Any], section: str) -> ModelConfig:
    value = raw.get(section) or {}
    config = ModelConfig(
        provider=str(value.get("provider", "openai-compatible")),
        base_url=str(value.get("base_url", "")).rstrip("/"),
        model=str(value.get("model", "")),
        api_key=str(value.get("api_key", "")),
    )
    missing = [name for name in ("base_url", "model", "api_key") if not getattr(config, name)]
    if missing:
        raise ValueError(f"{section}.{', '.join(missing)} is required")
    return config


def _load_rag_config(raw: dict[str, Any]) -> RAGConfig:
    document_paths = raw.get("document_paths") or ["../docs", "../README.md"]
    return RAGConfig(
        index_version=str(raw.get("index_version", "rag-docs-v1")),
        storage_dir=str(raw.get("storage_dir", "storage/rag_index")),
        document_paths=[str(item) for item in document_paths],
    )
