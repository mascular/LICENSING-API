package api

import (
    "encoding/json"
    "math/rand"
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"
    "waguri-auth/common"
)

type CreateKeyRequest struct {
    App      string `json:"app"`
    Duration string `json:"duration"`
}

type DeleteKeyRequest struct {
    App string `json:"app"`
    Key string `json:"key"`
}

type LoginRequest struct {
    App  string `json:"app"`
    Key  string `json:"key"`
    HWID string `json:"hwid"`
}

func generateKey() string {
    charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    key := make([]byte, 16)
    for i := range key {
        key[i] = charset[rand.Intn(len(charset))]
    }
    return "key-" + string(key)
}

func CreateKeyHandler(w http.ResponseWriter, r *http.Request) {
    var req CreateKeyRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        json.NewEncoder(w).Encode(common.MsgBadRequest)
        return
    }

    data, _ := LoadData(req.App)
    key := generateKey()
now := time.Now().UTC()
var end time.Time
duration := req.Duration

if duration == "0" || duration == "0d" || duration == "lifetime" {
    duration = "lifetime"
    end = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
} else {
    end = now.Add(parseDuration(req.Duration))
}

data[key] = map[string]string{
    "duration": duration,
    "start":    now.Format(time.RFC3339),
    "end":      end.Format(time.RFC3339),
    "hwid":     "",
}

    SaveData(req.App, data)

    common.SendUsageAlert(req.App, key, req.Duration, "Key Created")

    json.NewEncoder(w).Encode(struct {
        common.APIResponse
        Key string `json:"key"`
    }{common.MsgKeyCreated, key})
}

func DeleteKeyHandler(w http.ResponseWriter, r *http.Request) {
    var req DeleteKeyRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        json.NewEncoder(w).Encode(common.MsgBadRequest)
        return
    }

    data, _ := LoadData(req.App)
    if _, exists := data[req.Key]; !exists {
        json.NewEncoder(w).Encode(common.MsgKeyNotFound)
        return
    }

    delete(data, req.Key)
    SaveData(req.App, data)

    common.SendUsageAlert(req.App, req.Key, "", "Key Deleted")

    json.NewEncoder(w).Encode(common.MsgKeyDeleted)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        json.NewEncoder(w).Encode(common.MsgBadRequest)
        return
    }

    data, _ := LoadData(req.App)
    entry, exists := data[req.Key]
    if !exists {
        json.NewEncoder(w).Encode(common.MsgInvalidLicense)
        return
    }

    endTimeStr := entry["end"]
    endTime, err := time.Parse(time.RFC3339, endTimeStr)
    if err != nil || time.Now().UTC().After(endTime) {
        json.NewEncoder(w).Encode(common.MsgExpiredLicense)
        return
    }

    if entry["hwid"] == "" {
        entry["hwid"] = req.HWID
        data[req.Key] = entry
        SaveData(req.App, data)

        common.SendUsageAlert(req.App, req.Key, entry["duration"], "Signed up (first login)")

        json.NewEncoder(w).Encode(common.MsgSignedUp)
    } else if entry["hwid"] == req.HWID {
        common.SendUsageAlert(req.App, req.Key, entry["duration"], "Logged in")

        json.NewEncoder(w).Encode(common.MsgLoginSuccess)
    } else {
        common.SendUsageAlert(req.App, req.Key, entry["duration"], "HWID mismatch")

        json.NewEncoder(w).Encode(common.MsgInvalidLicense)
    }
}

func ResetHWIDHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        common.MethodNotAllowed(w)
        return
    }

    var req struct {
        App string `json:"app"`
        Key string `json:"key"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        common.BadRequest(w, "Invalid JSON")
        return
    }

    if req.App == "" || req.Key == "" {
        common.BadRequest(w, "Missing fields")
        return
    }

    data, err := LoadData(req.App)
    if err != nil {
        common.InternalError(w, "Failed to read license data")
        return
    }

    entry, exists := data[req.Key]
    if !exists {
        common.BadRequest(w, "Key not found")
        return
    }

    entry["hwid"] = ""
    data[req.Key] = entry

    if err := SaveData(req.App, data); err != nil {
        common.InternalError(w, "Failed to save data")
        return
    }

    common.SendUsageAlert(req.App, req.Key, entry["duration"], "HWID Reset")

    common.Success(w, "HWID reset for "+req.Key)
}

func ListAppsHandler(w http.ResponseWriter, r *http.Request) {
    files, err := os.ReadDir("data")
    if err != nil {
        json.NewEncoder(w).Encode(common.APIResponse{false, "Failed to read app list"})
        return
    }

    var apps []string
    for _, f := range files {
        if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
            apps = append(apps, strings.TrimSuffix(f.Name(), ".json"))
        }
    }

    json.NewEncoder(w).Encode(struct {
        Success bool     `json:"success"`
        Apps    []string `json:"apps"`
    }{true, apps})
}

func ListKeysHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        App string `json:"app"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.App == "" {
        json.NewEncoder(w).Encode(common.MsgBadRequest)
        return
    }

    data, err := LoadData(req.App)
    if err != nil {
        json.NewEncoder(w).Encode(common.MsgKeyNotFound)
        return
    }

    keys := make([]map[string]string, 0)
    for key, entry := range data {
        keys = append(keys, map[string]string{
            "key":   key,
            "start": entry["start"],
            "end":   entry["end"],
            "hwid":  entry["hwid"],
        })
    }

    json.NewEncoder(w).Encode(struct {
        Success bool                   `json:"success"`
        Keys    []map[string]string   `json:"keys"`
    }{true, keys})
}

func KeyInfoHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        json.NewEncoder(w).Encode(common.MsgBadRequest)
        return
    }

    var req struct {
        App string `json:"app"`
        Key string `json:"key"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        json.NewEncoder(w).Encode(common.MsgBadRequest)
        return
    }

    if req.App == "" || req.Key == "" {
        json.NewEncoder(w).Encode(common.MsgBadRequest)
        return
    }

    data, err := LoadData(req.App)
    if err != nil {
        json.NewEncoder(w).Encode(common.MsgKeyNotFound)
        return
    }

    entry, exists := data[req.Key]
    if !exists {
        json.NewEncoder(w).Encode(common.MsgKeyNotFound)
        return
    }

    json.NewEncoder(w).Encode(struct {
        common.APIResponse
        Info map[string]string `json:"info"`
    }{
        APIResponse: common.APIResponse{true, "Key details"},
        Info:        entry,
    })
}

func parseDuration(d string) time.Duration {
    if d == "0" {
        // Lifetime: 100 years
        return time.Hour * 24 * 365 * 100
    }
    if len(d) < 2 {
        return time.Hour * 24 * 30
    }

    num := 0
    fmt := d[len(d)-1]
    fmtStr := d[:len(d)-1]

    fmtParsed, err := time.ParseDuration(fmtStr + string(fmt))
    if err != nil {
        switch fmt {
        case 'd':
            num, _ = strconv.Atoi(fmtStr)
            return time.Hour * 24 * time.Duration(num)
        case 'h':
            num, _ = strconv.Atoi(fmtStr)
            return time.Hour * time.Duration(num)
        }
        return time.Hour * 24 * 30
    }
    return fmtParsed
}
