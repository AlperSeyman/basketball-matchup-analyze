import pandas as pd
from xgboost import  XGBRegressor

df = pd.read_parquet("data/features.parquet")



feature_columns = [
    "points_last5_avg", "assists_last5_avg", "rebounds_last5_avg",
    "fouls_last5_avg", "minutes_last5_avg",
    "home", "rest_days", "points_allowed_last5_avg",
]

points_model = XGBRegressor()
points_model.load_model("models/points_model.json")

rebounds_model = XGBRegressor()
rebounds_model.load_model("models/rebounds_model.json")

assists_model = XGBRegressor()
assists_model.load_model("models/assists_model.json")

fouls_model = XGBRegressor()
fouls_model.load_model("models/fouls_model.json")

sample = df[df["gameDate"] >= "2025-10-01"].iloc[[0]]


print("Player:", sample["firstName"].values[0], sample["lastName"].values[0])
print("Game date:", sample["gameDate"].values[0])

X_sample = sample[feature_columns]

print("Predicted points:", points_model.predict(X_sample)[0], "| Actual points:", sample["points"].values[0])
print("Predicted rebounds:", rebounds_model.predict(X_sample)[0], "| Actual rebounds:", sample["reboundsTotal"].values[0])
print("Predicted assists:", assists_model.predict(X_sample)[0], "| Actual assists:", sample["assists"].values[0])
print("Predicted fouls:", fouls_model.predict(X_sample)[0], "| Actual fouls:", sample["foulsPersonal"].values[0])