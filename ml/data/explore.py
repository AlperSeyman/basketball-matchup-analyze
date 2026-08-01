import sqlite3

con = sqlite3.connect("wyattowalsh/nba.sqlite")
cur = con.cursor()

cur.execute("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name;")
tables = [row[0] for row in cur.fetchall()]


for table in tables:
    print(f"\n--- {table} ---")
    cur.execute(f"PRAGMA table_info({table});")
    for col in cur.fetchall():
        print(" ",col[1])

