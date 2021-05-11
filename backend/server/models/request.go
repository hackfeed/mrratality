package models

type ParseRequest struct {
	UUID string `json:"uuid" binding:"required"`
}
