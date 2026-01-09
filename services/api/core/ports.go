package core

import "context"

type Pinger interface {
	Ping(context.Context) error
}

type Loginer interface {
	Login(name, password string) (string, error)
}
