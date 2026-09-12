from google import genai
from schemas import MatchUpData
from config import load

def analyze_matchup(data: MatchUpData) -> str:
    config = load()
    prompt = f"""
You are a basketball analyst. Compare these two NBA teams using ONLY real stat below.
Write a short, clear analysis (3-4 paragraphs) for a fan reading it before the game.

Team A: {data.team_a}
Team A stats: {data.team_a_stats}

Team B: {data.team_b}
Team B stats: {data.team_b_stats}

If a stat value looks like a zero or missing, mention that the season has not started yet
instead of guessing a reason
"""
    client = genai.Client()

    response = client.interactions.create(
        model = config.models.gemini,
        input=prompt
    )

    return response.output_text

