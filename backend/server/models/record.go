package models

type Record struct {
	UUID string `json:"uuid" binding:"required"`
}
