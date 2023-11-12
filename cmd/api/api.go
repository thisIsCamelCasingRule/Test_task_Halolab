package api

import "Test_task_Halolab/cmd/service"

type API struct {
	Service *service.Service
}

func NewAPI() *API {
	return &API{}
}

func (a *API) InitAPI(service *service.Service) {
	a.Service = service
}
