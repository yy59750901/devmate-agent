from app.schemas.requirement import RequirementAnalysisRequest, RequirementAnalysisResult


def analyze_requirement(request: RequirementAnalysisRequest) -> RequirementAnalysisResult:
    """First-stage mock workflow.

    This deliberately returns a deterministic structure before connecting a real LLM.
    The next step is to replace this implementation with an LLM-backed workflow,
    then evolve it into a LangGraph workflow.
    """

    requirement = request.requirement.strip()
    return RequirementAnalysisResult(
        summary=f"待分析需求：{requirement}",
        apis=[
            "POST /api/example：提交业务请求",
            "GET /api/example/{id}：查询处理结果",
        ],
        tables=[
            "example_tasks：保存任务状态、输入、输出和错误信息",
        ],
        risks=[
            "当前为规则化 mock 输出，尚未接入真实 LLM。",
            "后续需要增加结构化输出校验、重试和异常兜底。",
        ],
        test_cases=[
            "输入正常需求时应返回 summary、apis、tables、risks、test_cases、questions。",
            "输入空字符串时接口应返回参数校验错误。",
        ],
        questions=[
            "该需求的核心业务对象是什么？",
            "是否需要权限控制、审计日志或异步处理？",
        ],
    )
