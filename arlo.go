package arlo

import (
	"github.com/jeffreydwalter/arlo-golang/internal/request"
)

type Arlo struct {
	user         string
	pass         string
	client       *request.Client
	Account      Account
	Basestations Basestations
	Cameras      Cameras
}

func newArlo(user string, pass string) *Arlo {

	c, _ := request.NewClient(BaseUrl)
	arlo := &Arlo{
		user:   user,
		pass:   pass,
		client: c,
	}

	return arlo
}
