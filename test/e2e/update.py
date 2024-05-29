import argparse
import json
import requests

parser = argparse.ArgumentParser("port")
parser.add_argument("port", type=int)
args = parser.parse_args()
print("Start update test...")
host = f"http://localhost:{args.port}"
login = "bob"
pswd = "bob"
credentials = {"login": login, "password": pswd}
print("Fetching JWT token...")
resp = requests.post(host+"/api/login",data=json.dumps(credentials))
if resp.status_code != 200:
    print("Failed to fetch jwt token")
    print(resp)
    exit(1)
print("Updating...")
jwt = resp.headers.get("Authorization")
resp = requests.post(host+"/api/update",headers={"Authorization":jwt})
if resp.status_code != 200:
    print("Failed to update")
    print(resp)
    exit(1)
print("update test PASS")
