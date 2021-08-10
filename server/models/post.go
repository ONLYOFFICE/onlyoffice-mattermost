package models

type PostPermission struct {
	FileId      string
	Username    string
	Permissions Permissions
}
