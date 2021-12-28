package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

const dataSource string = "./aqua.db"

func main() {
	db, _ = sql.Open("sqlite3", dataSource)
	defer db.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	server(router)
	router.Run()
}
