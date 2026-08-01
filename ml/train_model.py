import os

import pandas as pd
from sklearn.metrics import mean_absolute_error
from xgboost import XGBRegressor

df = pd.read_parquet("data/features.parquet")

train_df = df[df["gameDate"] < "2025-10-01"]
test_df = df[df["gameDate"] >= "2025-10-01"]

print("Train rows:", train_df.shape[0])
print("Test rows:", test_df.shape[0])
print("Train date range:", train_df["gameDate"].min(), "to", train_df["gameDate"].max())
print("Test date range:", test_df["gameDate"].min(), "to", test_df["gameDate"].max())

feature_columns = [
    "points_last5_avg", "assists_last5_avg", "rebounds_last5_avg",
    "fouls_last5_avg", "minutes_last5_avg",
    "home", "rest_days", "points_allowed_last5_avg",
]

target_columns = ["points", "reboundsTotal", "assists", "foulsPersonal"]

X_train = train_df[feature_columns]
y_train = train_df[target_columns]

X_test = test_df[feature_columns]
y_test = test_df[target_columns]

print("X_train shape:", X_train.shape)
print("y_train shape:", y_train.shape)
print("X_test shape:", X_test.shape)
print("y_test shape:", y_test.shape)

print("Training points model...")
points_model = XGBRegressor()
points_model.fit(X_train, y_train["points"])

print("Training rebounds model...")
rebounds_model = XGBRegressor()
rebounds_model.fit(X_train, y_train["reboundsTotal"])

print("Training assists model...")
assists_model = XGBRegressor()
assists_model.fit(X_train, y_train["assists"])

print("Training fouls model...")
fouls_model = XGBRegressor()
fouls_model.fit(X_train, y_train["foulsPersonal"])

print("All models trained!")

points_pred = points_model.predict(X_test)
rebounds_pred = rebounds_model.predict(X_test)
assists_pred = assists_model.predict(X_test)
fouls_pred = fouls_model.predict(X_test)

print("Points MAE:", mean_absolute_error(y_test["points"], points_pred))
print("Rebounds MAE:", mean_absolute_error(y_test["reboundsTotal"], rebounds_pred))
print("Assists MAE:", mean_absolute_error(y_test["assists"], assists_pred))
print("Fouls MAE:", mean_absolute_error(y_test["foulsPersonal"], fouls_pred))

baseline_points_mae = mean_absolute_error(y_test["points"], X_test["points_last5_avg"])
baseline_rebounds_mae = mean_absolute_error(y_test["reboundsTotal"], X_test["rebounds_last5_avg"])
baseline_assists_mae = mean_absolute_error(y_test["assists"], X_test["assists_last5_avg"])
baseline_fouls_mae = mean_absolute_error(y_test["foulsPersonal"], X_test["fouls_last5_avg"])

print("Baseline Points MAE:", baseline_points_mae)
print("Baseline Rebounds MAE:", baseline_rebounds_mae)
print("Baseline Assists MAE:", baseline_assists_mae)
print("Baseline Fouls MAE:", baseline_fouls_mae)

os.makedirs("models", exist_ok=True)

points_model.save_model("models/points_model.json")
rebounds_model.save_model("models/rebounds_model.json")
assists_model.save_model("models/assists_model.json")
fouls_model.save_model("models/fouls_model.json")

print("Models saved.")
