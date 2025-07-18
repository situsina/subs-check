package main

import (
	"github.com/bestruirui/bestsub/internal/api/common"
	"github.com/bestruirui/bestsub/internal/api/server"
	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/core/task"
	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/utils/info"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/bestruirui/bestsub/internal/utils/shutdown"
)

func main() {

	info.Banner()

	cfg := config.Base()

	if err := log.Initialize(cfg.Log.Level, cfg.Log.Path, cfg.Log.Output); err != nil {
		panic(err)
	}
	if err := database.Initialize(cfg.Database.Type, cfg.Database.Path); err != nil {
		panic(err)
	}

	if err := server.Initialize(); err != nil {
		panic(err)
	}
	if err := task.Initialize(); err != nil {
		panic(err)
	}

	task.Start()
	server.Start()

	shutdown.Register(server.Close)        // 关闭顺序
	shutdown.Register(database.Close)      //   ↓↓
	shutdown.Register(common.CloseSession) //   ↓↓
	shutdown.Register(task.Shutdown)       //   ↓↓
	shutdown.Register(log.Close)           //   ↓↓

	shutdown.Listen()
}
