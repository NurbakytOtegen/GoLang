package models

type Role string

const (
	RoleUser       Role = "USER"
	RoleAdmin      Role = "ADMIN"
	RoleSuperAdmin Role = "SUPER_ADMIN"
)
