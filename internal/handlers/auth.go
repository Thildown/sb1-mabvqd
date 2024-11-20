package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginPage(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("user_id") != nil {
		c.Redirect(http.StatusSeeOther, "/deploy")
		return
	}
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Connexion",
	})
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Dans un environnement de production, utilisez une vraie authentification
	if username == "admin" && password == "admin" {
		session := sessions.Default(c)
		session.Set("user_id", username)
		session.Save()
		c.Redirect(http.StatusSeeOther, "/deploy")
		return
	}

	c.HTML(http.StatusUnauthorized, "login.html", gin.H{
		"error": "Identifiants invalides",
	})
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusSeeOther, "/")
}