from celery import Celery
from config import load
from schemas import MatchUpData, PlayerMatchUpData
from espn_client import get_standings, get_team_record, get_player_stats, get_player_season_averages


config = load()

celery_app = Celery("ai-service", broker=config.redis_url, backend=config.redis_url)


def gather_matchup_data(team_a_id: str, team_a_name: str, team_b_id: str, team_b_name: str) -> MatchUpData:
    standings_data = get_standings()
    team_a_stats = get_team_record(team_a_id, standings_data)
    team_b_stats = get_team_record(team_b_id, standings_data)

    return MatchUpData(
        team_a=team_a_name,
        team_b=team_b_name,
        team_a_stats=team_a_stats,
        team_b_stats=team_b_stats,
        injuries=[],
        recent_head_to_head=[],
        source_url=[],
    )

def gather_player_matchup_data(player_a_id: str, player_a_name: str, player_b_id: str, player_b_name: str) -> PlayerMatchUpData:
    player_a_data = get_player_stats(player_a_id)
    player_a_stats = get_player_season_averages(player_a_data)

    player_b_data = get_player_stats(player_b_id)
    player_b_stats = get_player_season_averages(player_b_data)

    return PlayerMatchUpData(
        player_a=player_a_name,
        player_b=player_b_name,
        player_a_stats=player_a_stats,
        player_b_stats=player_b_stats,
    )