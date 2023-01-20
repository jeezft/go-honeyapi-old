package wb

type Api struct {
	Auth string
}

func New() *Api {
	return &Api{}
}
