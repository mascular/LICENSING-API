import requests
import json

BASE_URL = "http://localhost:90"
API_KEY = "your-secret-api-key"  # replace with your real key
APP_NAME = "mytool"
HWID = "test-hwid-lifetime"

HEADERS = {
    "X-Api-Key": API_KEY,
    "Content-Type": "application/json"
}

def create_lifetime_key():
    payload = {
        "app": APP_NAME,
        "duration": "0"  # Lifetime
    }
    r = requests.post(f"{BASE_URL}/create-key", json=payload, headers=HEADERS)
    print("[Create Lifetime Key]", r.status_code, r.json())
    return r.json().get("key")

def login_with_key(key):
    payload = {
        "app": APP_NAME,
        "key": key,
        "hwid": HWID
    }
    r = requests.post(f"{BASE_URL}/login", json=payload)
    print("[Login]", r.status_code, r.json())

def key_info(key):
    payload = {
        "app": APP_NAME,
        "key": key
    }
    r = requests.post(f"{BASE_URL}/key-info", json=payload, headers=HEADERS)
    print("[Key Info]", r.status_code, json.dumps(r.json(), indent=2))

def main():
    key = create_lifetime_key()
    if key:
        login_with_key(key)
        key_info(key)

if __name__ == "__main__":
    main()
