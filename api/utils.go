package api

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "os"
    "waguri-auth/common"
    // "waguri-auth"
)
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        apiKey := r.Header.Get("X-Api-Key")
        if apiKey != common.AppConfig.APIKey {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    }
}

func LoadData(app string) (map[string]map[string]string, error) {
    path := "data/" + app + ".json"
    data := make(map[string]map[string]string)

    if _, err := os.Stat(path); os.IsNotExist(err) {
        return data, nil
    }

    content, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    json.Unmarshal(content, &data)
    return data, nil
}

func SaveData(app string, data map[string]map[string]string) error {
    content, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }
    return ioutil.WriteFile("data/"+app+".json", content, 0644)
}
