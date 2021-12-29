package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
)

var db *sql.DB

const errorDBQuery string = "Error occurred while querying database."
const errorNoMatch string = "No matching entry found with key : "
const errorEmptyTable string = "There are currently no entries for in this table."

func server(router *gin.Engine) {
	//router.Use(AddHeader())
	router.GET("/hosts", getAllHosts)                          // get all hosts
	router.GET("/hosts/:id", getHostByID)                      // get a host by host id
	router.GET("/containers", getAllContainers)                // get all containers
	router.GET("/containers/:id", getContainerByID)            // get a container by container id
	router.GET("/containers/hosts/:id", getContainersByHostID) // get all container with host id
	router.POST("/containers", postContainer)                  // get all container with host id
}

func getAllHosts(gctx *gin.Context) {
	row, err := db.Query("SELECT * FROM hosts")
	defer row.Close()
	if err != nil {
		gctx.JSON(http.StatusOK, gin.H{"Error": errorDBQuery + err.Error()})
		return
	}
	fetchHostRecords(gctx, row)
}

func getHostByID(gctx *gin.Context) {
	key := gctx.Params.ByName("id") // parse url for parameter
	row, err := db.Query("SELECT * FROM hosts WHERE id = ? ", key)
	if err != nil {
		gctx.JSON(http.StatusInternalServerError, gin.H{"Hosts": errorDBQuery + key})
		return
	}
	fetchHostRecords(gctx, row)
}

func getAllContainers(gctx *gin.Context) {
	row, err := db.Query("SELECT * FROM containers")
	defer row.Close()
	if err != nil {
		gctx.JSON(http.StatusInternalServerError, gin.H{"Containers": errorDBQuery})
		return
	}
	fetchContainerRecords(gctx, row)
}

func getContainerByID(gctx *gin.Context) {
	key := gctx.Params.ByName("id") // parse url for parameter
	row, err := db.Query("SELECT * FROM containers WHERE id = ?", key)
	defer row.Close()
	if err != nil {
		gctx.JSON(http.StatusInternalServerError, gin.H{"Containers": errorDBQuery})
		return
	}
	fetchContainerRecords(gctx, row)
}

func getContainersByHostID(gctx *gin.Context) {
	key := gctx.Params.ByName("id") // parse url for parameter
	row, err := db.Query("SELECT * FROM containers WHERE host_id = ? ", key)
	if err != nil {
		gctx.JSON(http.StatusInternalServerError, gin.H{"Containers": errorDBQuery})
	}
	fetchContainerRecords(gctx, row)
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
		gctx.JSON(http.StatusInternalServerError, gin.H{"Error": "An error occurred while preparing to insert into database."})
		return
	}
	_, err = statement.Exec(container.HostID, container.Name, container.ImageName)
	if err != nil {
		gctx.JSON(http.StatusInternalServerError, gin.H{"Error": "An error occurred while trying to insert into database."})
		return
	}

	gctx.Writer.WriteHeader(http.StatusCreated)
	return
}

//// AddHeader is  helper function to create header
//func AddHeader() gin.HandlerFunc {
//	return func(gctx *gin.Context) {
//		gctx.Header("Access-Control-Allow-Origin", "*")
//		gctx.Header("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
//		gctx.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
//		if gctx.Request.Method == "OPTIONS" {
//			gctx.AbortWithStatus(204)
//		}
//		gctx.Next()
//	}
//}

// fetchHostRecords	is a helper function to collect the result of a database query
func fetchHostRecords(gctx *gin.Context, row *sql.Rows) {
	host := Host{}
	listOfHosts := make([]interface{}, 0)
	var err error

	// Iterate and fetch the records from result cursor
	for row.Next() {
		resp := make(map[string]string)
		err = row.Scan(&host.ID, &host.UUID, &host.Name, &host.IpAddress)
		if err != nil {
			gctx.JSON(http.StatusInternalServerError, gin.H{"Error": errorDBQuery})
			return
		}
		// create map of attributes
		resp[hostID] = strconv.Itoa(host.ID)
		resp[hostUUID] = host.UUID
		resp[hostName] = host.Name
		resp[hostIPAddress] = host.IpAddress

		listOfHosts = append(listOfHosts, resp)
	}

	if len(listOfHosts) == 0 {
		gctx.JSON(http.StatusOK, gin.H{"Containers": errorEmptyTable})
	} else {
		gctx.JSON(http.StatusOK, gin.H{"Containers": listOfHosts})
	}
}

// fetchContainerRecords is a helper function to collect the result of a database query
func fetchContainerRecords(gctx *gin.Context, row *sql.Rows) {
	container := Container{}
	listOfContainers := make([]interface{}, 0)
	var err error

	// Iterate and fetch the records from result cursor
	for row.Next() {
		resp := make(map[string]interface{})
		err = row.Scan(&container.ID, &container.HostID, &container.Name, &container.ImageName)
		if err != nil {
			gctx.JSON(http.StatusInternalServerError, gin.H{"Error": errorDBQuery})
			return
		}
		// create map of attributes
		resp[containerID] = container.ID
		resp[containerHostID] = container.HostID
		resp[containerName] = container.Name
		resp[containerImageName] = container.ImageName

		listOfContainers = append(listOfContainers, resp)
	}

	if len(listOfContainers) == 0 {
		gctx.JSON(http.StatusOK, gin.H{"Containers": errorEmptyTable})
	} else {
		gctx.JSON(http.StatusOK, gin.H{"Containers": listOfContainers})
	}

}
