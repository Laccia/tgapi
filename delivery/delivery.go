package delivery

import (
	"myapp/app/configs"
	"myapp/logic"
)

type Delivery struct {
	Config *configs.Server
	Logic  *logic.Logic
}

func New(conf *configs.Server) *Delivery {
	return &Delivery{
		Config: conf,
		Logic:  logic.New(conf),
	}
}
