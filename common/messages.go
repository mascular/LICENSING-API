package common

type APIResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

var (
    MsgUnauthorized      = APIResponse{false, "Unauthorized"}
    MsgKeyCreated        = APIResponse{true, "License key created"}
    MsgKeyDeleted        = APIResponse{true, "License key deleted"}
    MsgSignedUp          = APIResponse{true, "Successfully signed up"}
    MsgLoginSuccess      = APIResponse{true, "Logged in"}
    MsgInvalidLicense    = APIResponse{false, "Invalid license or hardware ID"}
    MsgKeyNotFound       = APIResponse{false, "License key not found"}
    MsgBadRequest        = APIResponse{false, "Invalid request"}
    MsgExpiredLicense = APIResponse{false, "License key expired"}
)
