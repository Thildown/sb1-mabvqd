package main

import (
	"html/template"
	"log"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func loadTemplates() *template.Template {
	templ := template.Must(template.ParseGlob("templates/*.html"))
	return templ
}

func main() {
	r := gin.Default()

	// Configuration des sessions
	store := cookie.NewStore([]byte("change_this_secret_in_production"))
	r.Use(sessions.Sessions("vm-deployer", store))

	// Configuration des templates
	r.SetHTMLTemplate(loadTemplates())
	
	// Configuration des fichiers statiques
	r.Static("/static", filepath.Join("static"))

	// Routes publiques
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{
			"title": "Connexion",
		})
	})
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		if username == "admin" && password == "admin" {
			session := sessions.Default(c)
			session.Set("user_id", username)
			session.Save()
			c.Redirect(302, "/deploy")
			return
		}

		c.HTML(401, "login.html", gin.H{
			"title": "Connexion",
			"error": "Identifiants invalides",
		})
	})
	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(302, "/")
	})

	// Routes protégées
	authorized := r.Group("/")
	authorized.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.Redirect(302, "/")
			c.Abort()
			return
		}
		c.Next()
	})
	{
		authorized.GET("/deploy", func(c *gin.Context) {
			c.HTML(200, "deploy.html", gin.H{
				"title": "Déployer une VM",
				"user_id": sessions.Default(c).Get("user_id"),
			})
		})
	}

	log.Fatal(r.Run(":8080"))
}