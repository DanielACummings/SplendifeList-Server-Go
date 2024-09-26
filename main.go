package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pelletier/go-toml"
)

type DbConfigsPrivate struct {
	Password string `toml:"password"`
}

type DbConfigsPublic struct {
	User     string `toml:"user"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Database string `toml:"database"`
}

func dbConnect(dbPublicConfigs DbConfigsPublic, dbPrivateConfigs DbConfigsPrivate) (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		dbPublicConfigs.User,
		dbPrivateConfigs.Password,
		dbPublicConfigs.Host,
		dbPublicConfigs.Port,
		dbPublicConfigs.Database,
	)
	dbConnection, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	return dbConnection, nil
}

func dbGetConfigs() (DbConfigsPublic, DbConfigsPrivate, error) {
	var publicConfigs DbConfigsPublic
	var privateConfigs DbConfigsPrivate

	// Get public configs
	publicConfigFile, err := os.Open("config.toml")
	if err != nil {
		return publicConfigs, privateConfigs, err
	}
	defer publicConfigFile.Close()

	rawPublicConfigs, err := toml.LoadReader(publicConfigFile)
	if err != nil {
		return publicConfigs, privateConfigs, err
	}

	err = rawPublicConfigs.Unmarshal(&publicConfigs)
	if err != nil {
		return publicConfigs, privateConfigs, err
	}

	// Get private configs
	privateConfigFile, err := os.Open("private_config.toml")
	if err != nil {
		return publicConfigs, privateConfigs, err
	}
	defer privateConfigFile.Close()

	rawPrivateConfigs, err := toml.LoadReader(privateConfigFile)
	if err != nil {
		return publicConfigs, privateConfigs, err
	}

	err = rawPrivateConfigs.Unmarshal(&publicConfigs)
	if err != nil {
		return publicConfigs, privateConfigs, err
	}

	return publicConfigs, privateConfigs, nil
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			msgStart := "Application error: "
			fmt.Println(msgStart, r)
			log.Println(msgStart, r)
		}
	}()

	app := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
		MaxAge:        12 * time.Hour,
	}
	app.Use(cors.New(corsConfig))

	dbPublicConfigs, dbPrivateConfigs, err := dbGetConfigs()
	if err != nil {
		panic(err)
	}

	dbConnection, err := dbConnect(dbPublicConfigs, dbPrivateConfigs)

	baseRoute := "api"

	// Test route
	app.GET(baseRoute, func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "Connected to server!"})
	})

	// Lists

	app.Run(":8000") // Defaults to localhost:8080 when no port given
}
