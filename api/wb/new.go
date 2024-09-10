package wb

import "github.com/ProSellers/go-honeyapi/internal/controllers"

type Api struct {
	Auth  string
	Proxy *controllers.Proxy
}

func New() *Api {
	return &Api{}
}

func (a *Api) SetProxy(p *controllers.Proxy) {
	a.Proxy = p
}
