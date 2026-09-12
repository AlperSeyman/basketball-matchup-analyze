from espn_client import get_player_stats, get_player_season_averages


data = get_player_stats("1966")
averages = get_player_season_averages(data)
print(averages)