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

	"SplendifeList-Server-Go/models"
	"SplendifeList-Server-Go/services/item_list_service"
	"SplendifeList-Server-Go/services/item_service"
)

type DbConfigsPrivate struct {
	Password string `toml:"password"`
}

type DbConfigsPublic struct {
	User           string `toml:"user"`
	Database       string `toml:"database"`
	UnixSocketPath string `toml:"unix_socket_path"`
}

func dbConnect(dbPublicConfigs DbConfigsPublic,
	dbPrivateConfigs DbConfigsPrivate) (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"%s:%s@unix(%s)/%s",
		dbPublicConfigs.User,
		dbPrivateConfigs.Password,
		dbPublicConfigs.UnixSocketPath,
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
	publicConfigsFile, err := os.Open("config.toml")
	if err != nil {
		return publicConfigs, privateConfigs, err
	}
	defer publicConfigsFile.Close()

	rawPublicConfigs, err := toml.LoadReader(publicConfigsFile)
	if err != nil {
		return publicConfigs, privateConfigs, err
	}

	err = rawPublicConfigs.Unmarshal(&publicConfigs)
	if err != nil {
		return publicConfigs, privateConfigs, err
	}

	// Get private configs
	privateConfigsFile, err := os.Open("config_private.toml")
	if err != nil {
		return publicConfigs, privateConfigs, err
	}
	defer privateConfigsFile.Close()

	rawPrivateConfigs, err := toml.LoadReader(privateConfigsFile)
	if err != nil {
		return publicConfigs, privateConfigs, err
	}

	err = rawPrivateConfigs.Unmarshal(&privateConfigs)
	if err != nil {
		return publicConfigs, privateConfigs, err
	}

	return publicConfigs, privateConfigs, nil
}

func runProgram() error {
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
		return err
	}

	dbConnection, err := dbConnect(dbPublicConfigs, dbPrivateConfigs)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConnection.Close()

	baseRoute := "api"

	// Test route
	app.GET(baseRoute, func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "Connected to server!"})
	})

	// Lists
	app.GET(baseRoute+"/lists", func(context *gin.Context) {
		var itemLists []models.ItemList
		itemLists, err := item_list_service.GetAllLists(dbConnection)
		if err != nil {
			context.JSON(http.StatusInternalServerError,
				gin.H{"error": err.Error()})

			return
		}

		context.JSON(http.StatusOK, itemLists)
	})
	app.POST(baseRoute+"/lists", func(context *gin.Context) {
		var newItemList models.ItemList
		err = context.BindJSON(&newItemList)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		newItemListId, err := item_list_service.CreateList(dbConnection, newItemList)
		if err != nil {
			context.JSON(http.StatusInternalServerError,
				gin.H{"error": err.Error()})

			return
		}

		context.JSON(
			http.StatusOK,
			gin.H{"message": fmt.Sprintf("New list ID: %d", newItemListId)},
		)
	})

	// Items
	app.GET(baseRoute+"/lists/:listId/items",
		func(context *gin.Context) {
			userId := context.Query("userId")
			var items []models.Item
			items, err := item_service.GetItemsByItemList(dbConnection,
				context.Param("listId"), userId)
			if err != nil {
				context.JSON(http.StatusInternalServerError,
					gin.H{"error": err.Error()})

				return
			}

			context.JSON(http.StatusOK, items)
		})

	app.Run(":8000") // Defaults to localhost:8080 when no port given

	return nil
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			msgStart := "Application error: "
			fmt.Println(msgStart, r)
			log.Println(msgStart, r)
		}
	}()

	if err := runProgram(); err != nil {
		log.Fatal(err)
	}
}
