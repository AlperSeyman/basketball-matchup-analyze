import pandas as pd

df = pd.read_csv("data/eoinamoore/PlayerStatistics.csv")

print("Rows, columns:", df.shape)
print()
print(list(df.columns))
print()
print("Earliest game date:", df["gameDate"].min())
print("Latest game date:", df["gameDate"].max())