package handlars

import (
	"mymachine707/config"
	"mymachine707/storage"
)

// Handler ...
type Handler struct {
	Stg storage.Interfaces
	Cfg config.Config
}

// func NewHandler(sfg storage.Interfaces, cfg config.Config) handler {
// 	return handler{
// 		Sfg: sfg,
// 		Cfg: cfg,
// 	}
// }
