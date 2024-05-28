import json
import argparse
import requests

parser = argparse.ArgumentParser("port")
parser.add_argument("port", type=int)
args = parser.parse_args()

print("Start pics test...")
host = f"http://localhost:{args.port}"
query = "apple,doctor"
expected = "https://imgs.xkcd.com/comics/an_apple_a_day.png"

print("Fetching pics...")
resp = requests.get(f"{host}/api/pics?search={query}")
if resp.status_code != 200:
    print("Failed to fetch pics")
    print(resp)
    exit(1)
pics = json.loads(resp.text)
if expected not in pics:
    print("Expected comic not found")
    print(resp)
    exit(1)
print("pics test PASS")
