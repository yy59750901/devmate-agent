from pydantic import BaseModel, Field


class RequirementAnalysisRequest(BaseModel):
    requirement: str = Field(..., min_length=1, description="Product or engineering requirement text")


class RequirementAnalysisResult(BaseModel):
    summary: str
    apis: list[str]
    tables: list[str]
    risks: list[str]
    test_cases: list[str]
    questions: list[str]
