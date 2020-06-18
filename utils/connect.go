package utils

import (
	"database/sql"
	"log"

	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
)

// DbmapMySQLInkafarma variable de conecccion
var DbmapMySQLInkafarma = InitDbConnet()

// InitDbConnet Funcion que carga la conecccion a al BD Global
func InitDbConnet() *gorp.DbMap {
	fileConfig := "ProductService2"
	viper.SetConfigName(fileConfig) // name of config file (without extension) envoca un archivo json y actualiza la informacion
	viper.AddConfigPath(".")        // path to look for the config file in

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Config not found..." + fileConfig)
		log.Println(err)
		panic(err)
	}
	var cadenaCon string

	cadenaCon = viper.GetString("user") + ":" +
		viper.GetString("pass") + "@tcp(" + viper.GetString("server") + ":" +
		viper.GetString("port") + ")/" +
		viper.GetString("schema")

	log.Println("Config found, name = ", cadenaCon)

	db, err := sql.Open("mysql", cadenaCon)
	if err != nil {
		log.Println(err)
	}
	dbmapMySQL := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	return dbmapMySQL
}
