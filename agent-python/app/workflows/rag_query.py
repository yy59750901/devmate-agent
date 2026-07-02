from __future__ import annotations

import json
import logging
from pathlib import Path
from threading import Lock
from typing import Any

import httpx

from llama_index.core import Settings, SimpleDirectoryReader, StorageContext, VectorStoreIndex, load_index_from_storage
from llama_index.core.schema import NodeWithScore
from llama_index.embeddings.openai_like import OpenAILikeEmbedding
from llama_index.llms.openai_like import OpenAILike

from app.config import AppConfig, load_config, project_relative_path, resolve_agent_path
from app.schemas.rag import RAGMetadata, RAGQueryRequest, RAGQueryResponse, RAGSource

_logger = logging.getLogger("rag_query")
_logger.setLevel(logging.INFO)
_logger.propagate = True
if not _logger.handlers:
    _logger.addHandler(logging.StreamHandler())
_index_lock = Lock()
_index_cache: VectorStoreIndex | None = None
_index_config: AppConfig | None = None
_debug_patched = False


def query_rag(request: RAGQueryRequest) -> RAGQueryResponse:
    config = load_config()
    index = get_index(config)
    query_engine = index.as_query_engine(similarity_top_k=request.top_k)
    response = query_engine.query(request.question.strip())
    sources = [_source_from_node(node) for node in getattr(response, "source_nodes", [])]
    return RAGQueryResponse(
        answer=str(response),
        sources=sources,
        metadata=RAGMetadata(index_version=config.rag.index_version, top_k=request.top_k),
    )


def get_index(config: AppConfig) -> VectorStoreIndex:
    global _index_cache, _index_config
    with _index_lock:
        if _index_cache is not None and _index_config == config:
            return _index_cache

        configure_llama_index(config)
        storage_dir = resolve_agent_path(config.rag.storage_dir)
        if _has_persisted_index(storage_dir):
            storage_context = StorageContext.from_defaults(persist_dir=str(storage_dir))
            index = load_index_from_storage(storage_context)
        else:
            documents = load_documents(config)
            index = VectorStoreIndex.from_documents(documents)
            storage_dir.mkdir(parents=True, exist_ok=True)
            index.storage_context.persist(persist_dir=str(storage_dir))

        _index_cache = index
        _index_config = config
        return index


def configure_llama_index(config: AppConfig) -> None:
    Settings.llm = OpenAILike(
        model=config.llm.model,
        api_key=config.llm.api_key,
        api_base=config.llm.base_url,
        is_chat_model=True,
    )

    Settings.embed_model = OpenAILikeEmbedding(
        model_name=config.embedding.model,
        api_key=config.embedding.api_key,
        api_base=config.embedding.base_url,
    )

    _install_debug_logging()


def _install_debug_logging() -> None:
    """Intercept httpx to log full LLM / Embedding request and response bodies."""
    global _debug_patched
    if _debug_patched:
        return
    _debug_patched = True

    import sys as _sys

    def _log_line(msg: str) -> None:
        _sys.stderr.write(f"[RAG DEBUG] {msg}\n")
        _sys.stderr.flush()

    _log_line("Debug logging enabled for LLM / Embedding calls")

    _orig_send = httpx.Client.send
    _talk = {"embeddings", "chat/completions"}

    def _send(client: httpx.Client, request: httpx.Request, *args: Any, **kwargs: Any) -> httpx.Response:  # noqa: ARG001
        if request.method == "POST" and any(endpoint in str(request.url) for endpoint in _talk):
            body_bytes = request.read()
            try:
                body_json = json.loads(body_bytes)
            except json.JSONDecodeError:
                body_json = {"_raw": body_bytes.decode(errors="replace")}
            safe_body = _mask_api_key(body_json)
            _log_line(f"LLM/EMBED REQUEST | url={request.url} | body={json.dumps(safe_body, ensure_ascii=False)}")

            response = _orig_send(client, request, *args, **kwargs)  # type: ignore[arg-type]

            try:
                resp_json = response.json()
            except json.JSONDecodeError:
                resp_json = {"_raw": response.text}
            _log_line(f"LLM/EMBED RESPONSE | status={response.status_code} | body={json.dumps(resp_json, ensure_ascii=False)}")
            return response

        return _orig_send(client, request, *args, **kwargs)  # type: ignore[arg-type]

    httpx.Client.send = _send  # type: ignore[method-assign]


def _mask_api_key(body: dict[str, Any]) -> dict[str, Any]:
    """Remove api_key fields before logging."""
    masked = json.loads(json.dumps(body, ensure_ascii=False))
    for key in ("api_key", "Authorization", "authorization"):
        if key in masked:
            masked[key] = "***"
    if "headers" in masked and isinstance(masked["headers"], dict):
        for sensitive in ("authorization", "api-key", "x-api-key"):
            if sensitive in masked["headers"]:
                masked["headers"][sensitive] = "***"
    return masked

def load_documents(config: AppConfig) -> list[Any]:
    input_files = collect_document_files(config)
    if not input_files:
        raise ValueError("no documents found for RAG index")
    return SimpleDirectoryReader(input_files=[str(path) for path in input_files]).load_data()


def collect_document_files(config: AppConfig) -> list[Path]:
    files: list[Path] = []
    for configured_path in config.rag.document_paths:
        path = resolve_agent_path(configured_path)
        if path.is_file() and path.suffix.lower() == ".md":
            files.append(path)
        elif path.is_file() and path.name == "README.md":
            files.append(path)
        elif path.is_dir():
            files.extend(sorted(path.glob("*.md")))
    return sorted(set(files))


def _has_persisted_index(storage_dir: Path) -> bool:
    return storage_dir.exists() and any(storage_dir.iterdir())


def _source_from_node(node_with_score: NodeWithScore) -> RAGSource:
    node = node_with_score.node
    metadata = node.metadata or {}
    file_path = Path(str(metadata.get("file_path") or metadata.get("filename") or "unknown"))
    content = node.get_content(metadata_mode="none").strip()
    return RAGSource(
        path=project_relative_path(file_path) if file_path != Path("unknown") else "unknown",
        title=_title_from_content(content, file_path),
        snippet=_snippet(content),
        score=node_with_score.score,
    )


def _title_from_content(content: str, file_path: Path) -> str:
    for line in content.splitlines():
        stripped = line.strip()
        if stripped.startswith("#"):
            return stripped.lstrip("#").strip()
    if file_path.name:
        return file_path.name
    return "unknown"


def _snippet(content: str, max_length: int = 300) -> str:
    compact = " ".join(content.split())
    if len(compact) <= max_length:
        return compact
    return compact[:max_length].rstrip() + "..."
