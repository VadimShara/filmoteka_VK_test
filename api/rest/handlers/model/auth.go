package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Token struct {
	AccessToken string `json:"access_token"`
}

type LoginRequest struct {
	Username string `json:"username" example:"testuser"`
	Password string `json:"password" example:"tesTpa$$w0rd"`
}

type SignUpRequest struct {
	Username string `json:"username" example:"testuser"`
	Password string `json:"password" example:"tesTpa$$w0rd"`
	Role     string `json:"role" example:"admin"`
}

type LoginResponse struct {
	Status  string `json:"status" example:"OK"`
	MsgCode string `json:"msg_code" example:"admin_panel_login_success"`
	Data    Token  `json:"data" example:"{\"access_token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTA1NDE4OTUsImlhdCI6MTcxMDU0MDk5NSwianRpIjoiMzE4MTJkNjQtMGZkYy00ZmU2LWE4NGUtYmRlZDBjMjVhMGFhIiwidXNlciI6InRlc3R1c2VyIiwicm9sZSI6InVzZXIiLCJ0b2tlbl90eXBlIjoiYWNjZXNzX3Rva2VuIn0._mTQ8tY52PAVOT21w72HCnUl6epEautV_v6eDh0FTlI\"}"`
}

type SignUpResponse struct {
	Status  string `json:"status" example:"OK"`
	MsgCode string `json:"msg_code" example:"admin_panel_refresh_success"`
	Data    User   `json:"data" example:"{\"id\":1,\"username\":\"testuser\",\"role\":\"user\",\"created_at\":\"2024-01-01T00:00:00Z\",\"updated_at\":\"2024-01-01T00:00:00Z\"}"`
}
