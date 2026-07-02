from pydantic import BaseModel, Field


class RAGQueryRequest(BaseModel):
    question: str = Field(..., min_length=1, max_length=1000, description="Question about project documents")
    top_k: int = Field(default=5, ge=1, le=10, description="Number of source chunks to retrieve")


class RAGSource(BaseModel):
    path: str
    title: str
    snippet: str
    score: float | None = None


class RAGMetadata(BaseModel):
    index_version: str
    top_k: int


class RAGQueryResponse(BaseModel):
    answer: str
    sources: list[RAGSource]
    metadata: RAGMetadata
