package api

import (
    "encoding/json"
    "math/rand"
    "net/http"
	"strconv"
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
    end := now.Add(parseDuration(req.Duration))

    data[key] = map[string]string{
        "duration": req.Duration,
        "start":    now.Format(time.RFC3339),
        "end":      end.Format(time.RFC3339),
        "hwid":     "",
    }

    SaveData(req.App, data)
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

    // 🔒 Expiration check
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
        json.NewEncoder(w).Encode(common.MsgSignedUp)
    } else if entry["hwid"] == req.HWID {
        json.NewEncoder(w).Encode(common.MsgLoginSuccess)
    } else {
        json.NewEncoder(w).Encode(common.MsgInvalidLicense)
    }
}


func parseDuration(d string) time.Duration {
    // Very basic parser for formats like "30d", "15h"
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
