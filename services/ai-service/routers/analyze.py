from fastapi import APIRouter, Depends
from schemas import AnalyzerRequest, PlayerAnalyzerRequest
from tasks import run_team_analysis, run_player_analysis
from dependencies import get_current_id

router = APIRouter()

@router.post("/api/analyze/submit/team")
async def submit_team_analysis(request: AnalyzerRequest, user_id: int = Depends(get_current_id)):
    result = run_team_analysis.delay(request.team_a_id, request.team_a_name, request.team_b_id, request.team_b_name)
    return {"task_id": result.id}


@router.post("/api/analyze/submit/player")
async def submit_player_analysis(request: PlayerAnalyzerRequest, user_id: int = Depends(get_current_id)):
    result = run_player_analysis.delay(request.player_a_id, request.player_a_name, request.player_b_id, request.player_b_name)
    return {"task_id": result.id}