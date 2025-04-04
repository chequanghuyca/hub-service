package common

import "log"

const (
	DbType     = 1
	DbTypeUser = 2
)

const (
	CurrentUser = "user"
)

func AppRecover() {
	if err := recover(); err != nil {
		log.Panicln("Recovered error:", err)
	}
}
