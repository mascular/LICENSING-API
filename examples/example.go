package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

var baseURL = "http://localhost:8080"
var apiKey = "your-api-key"
var app = "mytool"

func makeRequest(method, url string, body interface{}, withKey bool) map[string]interface{} {
    jsonData, _ := json.Marshal(body)
    req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    if withKey {
        req.Header.Set("X-Api-Key", apiKey)
    }

    client := &http.Client{}
    res, _ := client.Do(req)
    defer res.Body.Close()

    responseBody, _ := ioutil.ReadAll(res.Body)
    var result map[string]interface{}
    json.Unmarshal(responseBody, &result)
    return result
}

func main() {
    // Create key
    createRes := makeRequest("POST", baseURL+"/create-key", map[string]string{
        "app": app, "duration": "30d",
    }, true)
    fmt.Println("Create Key:", createRes)

    key := createRes["key"].(string)

    // First login
    loginRes := makeRequest("POST", baseURL+"/login", map[string]string{
        "app": app, "key": key, "hwid": "abc-123",
    }, false)
    fmt.Println("Login:", loginRes)

    // Delete key
    deleteRes := makeRequest("POST", baseURL+"/delete-key", map[string]string{
        "app": app, "key": key,
    }, true)
    fmt.Println("Delete:", deleteRes)
}
