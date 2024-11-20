package main

import (
	"azure-vm-deployer/internal/handlers"
	"log"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Configuration des sessions
	store := cookie.NewStore([]byte("change_this_secret_in_production"))
	r.Use(sessions.Sessions("vm-deployer", store))

	// Configuration des templates et fichiers statiques
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", filepath.Join(".", "static"))

	// Routes publiques
	r.GET("/", handlers.LoginPage)
	r.POST("/login", handlers.Login)
	r.GET("/logout", handlers.Logout)

	// Routes protégées
	authorized := r.Group("/")
	authorized.Use(handlers.RequireAuth())
	{
		authorized.GET("/deploy", handlers.DeployPage)
		authorized.POST("/deploy", handlers.DeployVM)
		authorized.GET("/status/:id", handlers.VMStatus)
		authorized.GET("/vms", handlers.ListVMs)
	}

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}