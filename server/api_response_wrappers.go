package main

import "models"

type UserinfoWrapper struct {
	Id          string             `json:"id"`
	Username    string             `json:"username"`
	Permissions models.Permissions `json:"permissions"`
}
