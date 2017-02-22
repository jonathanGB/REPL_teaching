package route

import (
	"github.com/jonathanGB/REPL-Teaching/app/controllers"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"net/http"
)

func UserRoutes(router *gin.Engine, s *mgo.Session) {
	uc := controllers.NewUserController(s)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/users/signup/")
	})

	u := router.Group("/users")
	{
		u.GET("/signup", func(c *gin.Context) {
			c.HTML(http.StatusOK, "signup", gin.H{
				"title": "signup",
			})
		})

		u.POST("/signup", uc.CreateUser)

		u.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login", gin.H{
				"title": "login",
			})
		})

		u.POST("/login", uc.LoginUser)

		u.GET("/logout", func(c *gin.Context) {
			c.HTML(http.StatusOK, "logout", gin.H{
				"title": "logout",
				"name":  "user1234",
			})
		})
	}
}