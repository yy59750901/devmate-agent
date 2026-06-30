from fastapi import FastAPI

from app.schemas.requirement import RequirementAnalysisRequest, RequirementAnalysisResult
from app.workflows.requirement_analysis import analyze_requirement

app = FastAPI(title="DevMate Agent Service", version="0.1.0")


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok", "service": "devmate-agent"}


@app.post("/agent/requirement-analysis", response_model=RequirementAnalysisResult)
def requirement_analysis(request: RequirementAnalysisRequest) -> RequirementAnalysisResult:
    return analyze_requirement(request)
