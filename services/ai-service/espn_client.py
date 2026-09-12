import httpx

SITE_BASE = "https://site.api.espn.com/apis/site/v2/sports/basketball/nba"
CORE_BASE = "https://site.web.api.espn.com/apis/common/v3/sports/basketball/nba"

def get_player_stats(player_id: str) -> dict:
    response = httpx.get(f"{CORE_BASE}/athletes/{player_id}/stats")
    return response.json()


def get_team_stats(team_id: str) -> dict:
    response = httpx.get(f"{SITE_BASE}/teams/{team_id}/statistics")
    return response.json()


def get_standings() -> dict:
    response = httpx.get("https://site.api.espn.com/apis/v2/sports/basketball/nba/standings")
    return response.json()


def get_team_record(team_id: str, standings_data: dict) -> dict:
    for conference in standings_data["children"]:
        for entry in conference["standings"]["entries"]:
            if entry["team"]["id"] == team_id:
                record = {}
                for stat in entry["stats"]:
                    record[stat["name"]] = stat["displayValue"]
                return record
    return {}


def get_player_season_averages(stats_data: dict) -> dict:
    for category in stats_data["categories"]:
        if category["name"] == "averages":
            latest_season = category["statistics"][-1]
            return dict(zip(category["names"], latest_season["stats"]))
    return {}