package main

import (
	"seonbackendmovil/rest"
	"seonbackendmovil/utils"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

import (
	"fmt"
	"log"
)

func main() {

	today := time.Now()
	filename := fmt.Sprintf("%s_%d%d%d.log", "ProductService2", today.Year(), today.Month(), today.Day())
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error creando archivo de Log")
	}

	log.SetOutput(f)

	r := gin.New()

	r.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true))
	logrus.SetOutput(f)
	logger := logrus.New()
	logger.Level = logrus.ErrorLevel
	logger.Out = os.Stderr
	r.Use(ginrus.Ginrus(logger, time.RFC3339, false))

	r.Use(utils.Cors())

	authMiddleware := &utils.FirebaseGinMiddleware{
		TokenLookup: "x-access-token",
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
	}

	r.GET("/ping", utils.Ping)
	r.GET("/ver", utils.Version)

	r.Use(authMiddleware.MiddlewareFunc())
	r.GET("/category/parents", rest.GetCategoryrelationsChildsAndParents)
	r.GET("/getcategorybyparents/:id", rest.GetCategoryRelationsByID)
	r.GET("/search/keyword", rest.GetSearchKeyword)
	r.GET("/getprodbycategory", rest.GetProdbyCategory)
	r.GET("/getrelproducts", rest.GetRelProducts)
	r.GET("/product/:id", rest.GetProductByID)

	r.Run(":8082")

}
