package model

import "vk-test-task/pkg/web"

type ValidationErrors struct {
	Errors []web.ValidationError `json:"errors"`
}

type BadRequestInvalidBodyResponse struct {
	Status  string `json:"status" example:"ERROR"`
	MsgCode string `json:"msg_code" example:"invalid_body"`
}

type BadRequestInvalidIDResponse struct {
	Status  string `json:"status" example:"ERROR"`
	MsgCode string `json:"msg_code" example:"invalid_id"`
}

type BadRequestInvalidQueryResponse struct {
	Status  string `json:"status" example:"ERROR"`
	MsgCode string `json:"msg_code" example:"invalid_request_query"`
}

type BadRequestInvalidHeaderResponse struct {
	Status  string `json:"status" example:"ERROR"`
	MsgCode string `json:"msg_code" example:"invalid_header"`
}

type MovieNotFoundResponse struct {
	Status  string `json:"status" example:"ERROR"`
	MsgCode string `json:"msg_code" example:"movie_not_found"`
}

type StarNotFoundResponse struct {
	Status  string `json:"status" example:"ERROR"`
	MsgCode string `json:"msg_code" example:"star_not_found"`
}

type ValidationResponse struct {
	Status  string           `json:"status" example:"ERROR"`
	MsgCode string           `json:"msg_code" example:"go_validation"`
	Data    ValidationErrors `json:"data" example:"{\"errors\": [{\"tag\": \"<tag>\", \"field\": \"<field>\", \"param\": \"<param>\"}]}"`
}

type UnauthorizedResponse struct {
	Status  string `json:"status" example:"ERROR"`
	MsgCode string `json:"msg_code" example:"general_unauthorized"`
}

type InternalResponse struct {
	Status  string `json:"status" example:"ERROR"`
	MsgCode string `json:"msg_code" example:"general_internal"`
}

type ForbiddenResponse struct {
	Status  string `json:"status" example:"ERROR"`
	MsgCode string `json:"msg_code" example:"general_forbidden"`
}

type ConflictUsernameResponse struct {
	Status  string `json:"status" example:"ERROR"`
	MsgCode string `json:"msg_code" example:"username_is_taken"`
}

type WrongCredentialsResponse struct {
	Status  string `json:"status" example:"ERROR"`
	MsgCode string `json:"msg_code" example:"wrong_credentials"`
}
