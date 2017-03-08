package main

import (
	"github.com/reechou/robot-manager/config"
	"github.com/reechou/robot-manager/controller"
)

func main() {
	controller.NewLogic(config.NewConfig()).Run()
}
