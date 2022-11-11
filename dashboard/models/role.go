package models

import "github.com/google/uuid"

type Permission struct {
	ID   int    `json:"permId"`
	Name string `json:"permname"`
}

type Role struct {
	ID          int           `json:"roleId"`
	Name        string        `json:"roleName"`
	Permissions []*Permission `json:"permissions"`
	UserId      uuid.UUID     `json:"userId"`
}
