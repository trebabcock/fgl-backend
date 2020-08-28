package main

import (
	app "fgl-backend/app"
	config "fgl-backend/config"
	db "fgl-backend/db"
)

func main() {
	config := config.GetConfig()
	dbConfig := db.GetConfig()

	app := &app.App{}
	app.Initialize(dbConfig)
	app.Run(":" + config.Port)
}
