package app

import (
	"github.com/muchlist/mini_pos/configs"
	"github.com/muchlist/mini_pos/utils/logger"
)

func RunApp() {
	// Init config, logger dan db
	configs.InitConfig()
	logger.InitLogger()
}
