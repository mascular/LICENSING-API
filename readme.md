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
  "Port": "8080"
}
```

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