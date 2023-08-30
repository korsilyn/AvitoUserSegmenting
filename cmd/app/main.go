package main

import (
	"avito-user-segmenting/core/app"
)

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
