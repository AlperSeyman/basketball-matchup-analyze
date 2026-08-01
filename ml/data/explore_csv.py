import csv

for filename in ["eoinamoore/PlayerStatistics.csv", "eoinamoore/PlayerStatisticsExtended.csv", "eoinamoore/TeamStatistics.csv"]:
    with open(filename, newline="", encoding="utf-8") as f:
        reader = csv.reader(f)
        header = next(reader)
    print(f"\n--- {filename} ---")
    for col in header:
        print(" ", col)