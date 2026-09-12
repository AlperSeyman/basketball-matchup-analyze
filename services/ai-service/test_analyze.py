from tasks import gather_matchup_data
from gemini_client import analyze_matchup

data = gather_matchup_data("2", "Boston Celtics", "13", "Los Angeles")
result = analyze_matchup(data)
print(result)