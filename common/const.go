package common

import "log"

const (
	DbTypeUser = 1
)

const (
	CurrentUser = "user"
)

type Requester interface {
	GetUserId() int
	GetEmail() string
	GetRole() string
}

// Define roles for documentation and consistent usage
const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleClient     = "client"
)

func AppRecover() {
	if err := recover(); err != nil {
		log.Panicln("Recovered error:", err)
	}
}
