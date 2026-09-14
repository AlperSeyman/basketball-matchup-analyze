from pydantic import BaseModel

class AnalyzerRequest(BaseModel):
    team_a_id: str
    team_a_name: str
    team_b_id: str
    team_b_name: str

class MatchUpData(BaseModel):
    team_a: str
    team_b: str
    team_a_stats: dict
    team_b_stats: dict
    injuries: list[str]
    recent_head_to_head: list[dict]
    source_url: list[str]

class PlayerMatchUpData(BaseModel):
    player_a: str
    player_b: str
    player_a_stats: dict
    player_b_stats: dict

class PlayerAnalyzerRequest(BaseModel):
    player_a_id: str
    player_a_name: str
    player_b_id: str
    player_b_name: str