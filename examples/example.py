import requests

BASE_URL = "http://localhost:8080"
API_KEY = "your-api-key"
APP_NAME = "mytool"
HWID = "test-hwid-123"

def create_key():
    res = requests.post(f"{BASE_URL}/create-key", json={
        "app": APP_NAME, "duration": "30d"
    }, headers={"X-Api-Key": API_KEY})
    print("Create Key:", res.json())
    return res.json().get("key")

def login(key, hwid):
    res = requests.post(f"{BASE_URL}/login", json={
        "app": APP_NAME, "key": key, "hwid": hwid
    })
    print("Login:", res.json())

def delete_key(key):
    res = requests.post(f"{BASE_URL}/delete-key", json={
        "app": APP_NAME, "key": key
    }, headers={"X-Api-Key": API_KEY})
    print("Delete Key:", res.json())

if __name__ == "__main__":
    key = create_key()
    login(key, HWID)
    login(key, HWID)
    login(key, "another-hwid")
    delete_key(key)
