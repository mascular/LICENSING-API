package common

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var (
	MsgUnauthorized    = APIResponse{false, "Unauthorized"}
	MsgKeyCreated      = APIResponse{true, "License key created"}
	MsgKeyDeleted      = APIResponse{true, "License key deleted"}
	MsgSignedUp        = APIResponse{true, "Successfully signed up"}
	MsgLoginSuccess    = APIResponse{true, "Logged in"}
	MsgInvalidLicense  = APIResponse{false, "Invalid license or hardware ID"}
	MsgKeyNotFound     = APIResponse{false, "License key not found"}
	MsgBadRequest      = APIResponse{false, "Invalid request"}
	MsgExpiredLicense  = APIResponse{false, "License key expired"}
)

// Reusable response writers

func Success(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(APIResponse{true, message})
}

func BadRequest(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(APIResponse{false, message})
}

func InternalError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(APIResponse{false, message})
}

func MethodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(APIResponse{false, "Method not allowed"})
}
