package servinterface

import "github.com/ExtraWhy/internal-libs/config"

type ServiceInterface interface {
	DoRun(conf *config.RequestService) error
}
