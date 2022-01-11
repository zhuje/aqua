package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

var db *sql.DB

func server(router *gin.Engine) {
	router.Use(AddHeader())
	router.GET("/hosts", getAllHosts)                          // get all hosts
	router.GET("/hosts/:id", getHostByID)                      // get a host by host id
	router.GET("/containers", getAllContainers)                // get all containers
	router.GET("/containers/:id", getContainerByID)            // get a container by container id
	router.GET("/containers/hosts/:id", getContainersByHostID) // get all container with host id
	router.POST("/containers", postContainer)                  // post a new container
}

func getAllHosts(gctx *gin.Context) {
	row, err := db.Query("SELECT * FROM hosts")
	defer row.Close()
	if err != nil {
		queryError(gctx, err)
		return
	}
	emitHostRecords(gctx, row)
}

func getHostByID(gctx *gin.Context) {
	key := gctx.Params.ByName("id") // parse url for parameter
	row, err := db.Query("SELECT * FROM hosts WHERE id = ? ", key)
	if err != nil {
		queryError(gctx, err)
		return
	}
	emitHostRecords(gctx, row)
}

func getAllContainers(gctx *gin.Context) {
	row, err := db.Query("SELECT * FROM containers")
	defer row.Close()
	if err != nil {
		queryError(gctx, err)
		return
	}
	emitContainerRecords(gctx, row)
}

func getContainerByID(gctx *gin.Context) {
	key := gctx.Params.ByName("id") // parse url for parameter
	row, err := db.Query("SELECT * FROM containers WHERE id = ?", key)
	defer row.Close()
	if err != nil {
		queryError(gctx, err)
		return
	}
	emitContainerRecords(gctx, row)
}

func getContainersByHostID(gctx *gin.Context) {
	key := gctx.Params.ByName("id") // parse url for parameter
	row, err := db.Query("SELECT * FROM containers WHERE host_id = ? ", key)
	if err != nil {
		queryError(gctx, err)
	}
	emitContainerRecords(gctx, row)
}

func postContainer(gctx *gin.Context) {
	container := Container{}

	// bind response body to container object
	if err := gctx.BindJSON(&container); err != nil {
		gctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// create query and post to database
	insertContainerSQL := `INSERT INTO containers(host_id, name, image_name) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertContainerSQL)
	if err != nil {
		gctx.JSON(http.StatusInternalServerError, gin.H{"Error":
			"An error occurred while preparing to insert into database." + err.Error()})
		return
	}
	_, err = statement.Exec(container.HostID, container.Name, container.ImageName)
	if err != nil {
		gctx.JSON(http.StatusInternalServerError, gin.H{"Error":
			"An error occurred while trying to insert into database. " + err.Error() })
		return
	}

	gctx.Writer.WriteHeader(http.StatusCreated)
	return
}

func AddHeader() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		gctx.Header("Access-Control-Allow-Origin", "*")
		gctx.Header("Access-Control-Allow-Methods", " POST, GET")
		gctx.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers")
		gctx.Next()
	}
}

func emitHostRecords(gctx *gin.Context, row *sql.Rows) {
	host := Host{}
	listOfHosts := make([]Host, 0)
	var err error

	// iterate through rows, place into struct, and collect into a slice
	for row.Next() {
		err = row.Scan(&host.ID, &host.UUID, &host.Name, &host.IpAddress)
		if err != nil {
			queryError(gctx, err)
			return
		}
		listOfHosts = append(listOfHosts, host)
	}

	// emit response
	if len(listOfHosts) == 0 {
		gctx.JSON(http.StatusOK, gin.H{"Hosts": errorEmptyQuery})
	} else {
		gctx.JSON(http.StatusOK, gin.H{"Hosts": listOfHosts})
	}
}

func emitContainerRecords(gctx *gin.Context, row *sql.Rows) {
	container := Container{}
	listOfContainers := make([]Container, 0)
	var err error

	// iterate through rows, place into struct, and collect into a slice
	for row.Next() {
		err = row.Scan(&container.ID, &container.HostID, &container.Name, &container.ImageName)
		if err != nil {
			queryError(gctx, err)
			return
		}
		listOfContainers = append(listOfContainers, container)
	}

	// emit response
	if len(listOfContainers) == 0 {
		gctx.JSON(http.StatusOK, gin.H{"Containers": errorEmptyQuery})
	} else {
		gctx.JSON(http.StatusOK, gin.H{"Containers": listOfContainers})
	}

}

func queryError(gctx *gin.Context, err error){
	gctx.JSON(http.StatusInternalServerError, gin.H{"Error": errorDuringQuery + err.Error()})
}
