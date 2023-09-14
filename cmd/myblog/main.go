package main

import (
	"fmt"
	"log"
	//"net/http"

	"github.com/gin-gonic/gin"
	"subham.com/antino-assignment/internal/api"
	"subham.com/antino-assignment/pkg/blog"
)

func main() {
	// Initialize the Gin router
	router := gin.Default()

	// Initialize routes from the api package

	dsn := "root:iAMGOD@997@tcp(localhost:3306)/myblogdb"
	db := blog.InitializeDB(dsn)
	defer db.Close()

	api.InitRoutes(router)
	// Start the server
	port := 8080
	log.Printf("Server is running on port %d", port)
	router.Run(fmt.Sprintf(":%d", port))
}
