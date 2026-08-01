import pandas as pd

df = pd.read_csv("data/eoinamoore/PlayerStatistics.csv")
df["gameDate"] = pd.to_datetime(df["gameDate"])

exclude_types = ["All-Star Game", "Preseason", "Pre Season"]
df = df[df["gameType"].notna() & ~df["gameType"].isin(exclude_types)]

team_df = pd.read_csv("data/eoinamoore/TeamStatistics.csv")
team_df["gameDate"] = pd.to_datetime(team_df["gameDate"])
team_df = team_df[team_df["gameType"].notna() & ~team_df["gameType"].isin(exclude_types)]

team_df = team_df.sort_values(["teamId", "gameDate"])
team_df["points_allowed_last5_avg"] = team_df.groupby("teamId")["opponentScore"].transform(lambda x: x.shift(1).rolling(5).mean())

df = df.sort_values(["personId", "gameDate"])
df["points_last5_avg"] = df.groupby("personId")["points"].transform(lambda x: x.shift(1).rolling(5).mean())
df["assists_last5_avg"] = df.groupby("personId")["assists"].transform(lambda x: x.shift(1).rolling(5).mean())
df["rebounds_last5_avg"] = df.groupby("personId")["reboundsTotal"].transform(lambda x: x.shift(1).rolling(5).mean())
df["fouls_last5_avg"] = df.groupby("personId")["foulsPersonal"].transform(lambda x: x.shift(1).rolling(5).mean())
df["minutes_last5_avg"] = df.groupby("personId")["numMinutes"].transform(lambda x: x.shift(1).rolling(5).mean())
df["rest_days"] = df.groupby("personId")["gameDate"].diff().dt.days

df = df.merge(
    team_df[["gameId", "teamId", "points_allowed_last5_avg"]],
    left_on=["gameId", "opponentteamId"],
    right_on=["gameId", "teamId"],
    how="left",
)

feature_columns = [
    "points_last5_avg", "assists_last5_avg", "rebounds_last5_avg",
    "fouls_last5_avg", "minutes_last5_avg",
    "home", "rest_days", "points_allowed_last5_avg",
]
df = df.dropna(subset=feature_columns)

print("Date range in final data:", df["gameDate"].min(), "to", df["gameDate"].max())
print("Rows remaining after dropping:", df.shape)
print("Missing values in feature columns:\n", df[feature_columns].isna().sum())

columns_to_save = [
    "personId", "firstName", "lastName", "gameId", "gameDate", "opponentteamName",
    "points_last5_avg", "assists_last5_avg", "rebounds_last5_avg", "fouls_last5_avg", "minutes_last5_avg",
    "home", "rest_days", "points_allowed_last5_avg",
    "points", "reboundsTotal", "assists", "foulsPersonal",
]

df[columns_to_save].to_parquet("data/features.parquet")
