package main

import (
	"log"
	"os"

	app "flux-panel/golang-backend/internal/app"
	"flux-panel/golang-backend/internal/app/util"
	dbpkg "flux-panel/golang-backend/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// load .env if present
	util.LoadEnv()
	if err := dbpkg.Init(); err != nil {
		log.Fatalf("db init error: %v", err)
	}

	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	app.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "6365"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
