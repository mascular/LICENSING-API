# 🔐 License Auth API

A lightweight and secure license authentication system for your tools or apps — built in **Go** with JSON-based storage and HWID binding.

---

## 🚀 Features

- API Key-protected `/create-key` and `/delete-key`
- HWID-based login and validation
- Auto license expiry
- JSON-based data storage (per-app)
- Configurable via `config.json`

---

## ⚙️ Setup

### ✅ 1. Requirements

- Go 1.20+
- `config.json` file in project root

### ✅ 2. Installation

```bash
git clone https://github.com/your-username/license-auth-api.git
cd license-auth-api
go mod tidy
````

### ✅ 3. Configure

Create `config.json`:

```json
{
  "X-Api-Key": "your-secret-api-key",
  "Port": "8080",
  "discord_webhook": "your-discord-webhook-url"
}
```

**Note:** The `discord_webhook` is optional and used for sending usage alerts to Discord when keys are created, deleted, reset, or used.

---

## ▶️ Running the Server

```bash
go run main.go
```

Server will start at: `http://localhost:8080`

---

## 📬 API Endpoints

### 🔐 `POST /create-key`

> Create a license key

* **Headers**:
  `X-Api-Key: your-secret-api-key`

* **Body**:

```json
{
  "app": "mytool",
  "duration": "30d"
}
```

* ✅ **Response**:

```json
{
  "success": true,
  "message": "License key created",
  "key": "key-XYZ123abc"
}
```

---

### 🔓 `POST /login`

> Log in or validate HWID + license

* **Body**:

```json
{
  "app": "mytool",
  "key": "key-XYZ123abc",
  "hwid": "hwid-123"
}
```

* ✅ Responses:

```json
{ "success": true, "message": "Successfully signed up" }
{ "success": true, "message": "Logged in" }
```

* ❌ Invalid:

```json
{ "success": false, "message": "License key expired" }
{ "success": false, "message": "Invalid license or hardware ID" }
```

---

### 🗑️ `POST /delete-key`

> Delete a license key

* **Headers**:
  `X-Api-Key: your-secret-api-key`

* **Body**:

```json
{
  "app": "mytool",
  "key": "key-XYZ123abc"
}
```

* ✅ Response:

```json
{ "success": true, "message": "License key deleted" }
```

---

### 🔄 `POST /reset-hwid`

> Reset HWID binding for a license key

* **Headers**:
  `X-Api-Key: your-secret-api-key`

* **Body**:

```json
{
  "app": "mytool",
  "key": "key-XYZ123abc"
}
```


## ⏰ Duration Format

License duration can be specified using the following formats:

* `30d` - 30 days
* `7d` - 7 days
* `24h` - 24 hours
* `lifetime` or `0d` - Lifetime license (expires year 9999)

---

## 🔐 Authentication

All endpoints **except** `/login` require the `X-Api-Key` header with your configured API key from `config.json`.

---

## 🔔 Discord Notifications

When configured, the API sends Discord webhook notifications for the following events:

* Key Created
* Key Deleted
* HWID Reset
* Signed up (first login)
* Logged in
* HWID mismatch

---
* ✅ Response:

```json
{ "success": true, "message": "HWID reset for key-XYZ123abc" }
```

---

### 📋 `GET /list-apps`

> List all available apps

* **Headers**:
  `X-Api-Key: your-secret-api-key`

* ✅ Response:

```json
{
  "success": true,
  "apps": ["mytool", "otherapp", "thirdapp"]
}
```

---

### 🔑 `POST /list-keys`

> List all license keys for an app

* **Headers**:
  `X-Api-Key: your-secret-api-key`

* **Body**:

```json
{
  "app": "mytool"
}
```

* ✅ Response:

```json
{
  "success": true,
  "keys": [
    {
      "key": "key-XYZ123abc",
      "start": "2025-07-16T12:00:00Z",
      "end": "2025-08-15T12:00:00Z",
      "hwid": "hwid-123"
    }
  ]
}
```

---

### ℹ️ `POST /key-info`

> Get detailed information about a specific license key

* **Headers**:
  `X-Api-Key: your-secret-api-key`

* **Body**:

```json
{
  "app": "mytool",
  "key": "key-XYZ123abc"
}
```

* ✅ Response:

```json
{
  "success": true,
  "message": "Key details",
  "info": {
    "duration": "30d",
    "start": "2025-07-16T12:00:00Z",
    "end": "2025-08-15T12:00:00Z",
    "hwid": "hwid-123"
  }
}
```

---

## 📌 JSON License Format

Each app stores licenses in `data/<app>.json`:

```json
{
  "key-XYZ123abc": {
    "duration": "30d",
    "start": "2025-07-16T12:00:00Z",
    "end": "2025-08-15T12:00:00Z",
    "hwid": "hwid-123"
  }
}
```

---

## 🧪 Client Examples

Client usage in:

* [Python Example](examples/example.py)
* [Node.js Example](examples/example.js)
* [Go Client Example](examples/example.go)

---