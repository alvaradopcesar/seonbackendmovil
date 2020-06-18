package utils

import "github.com/gin-gonic/gin"

// Ping Para ver si esta activo el servicio
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// Version para saber la version de la App
func Version(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "1.0",
	})
}
