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
	GetConfigList(http.ResponseWriter, *http.Request)
	GetActualConfig(http.ResponseWriter, *http.Request)
	GetVersion(http.ResponseWriter, *http.Request)
	Rollback(w http.ResponseWriter, r *http.Request)
}

func NewConfigHandler(cfgService cs.ConfigServiceI) ConfigHandlerI {
	return &ConfigHandler{
		service: cfgService,
	}
}
