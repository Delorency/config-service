package configHandler

import (
	cs "main/internal/service/configService"
	"net/http"
)

type ConfigHandler struct {
	service cs.ConfigServiceI
}

type ConfigHandlerI interface {
	CreateConfig(http.ResponseWriter, *http.Request)
	UpdateConfig(http.ResponseWriter, *http.Request)
	GetConfig(http.ResponseWriter, *http.Request)
}

func NewConfigHandler(cfgService cs.ConfigServiceI) ConfigHandlerI {
	return &ConfigHandler{
		service: cfgService,
	}
}
