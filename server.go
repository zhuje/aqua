package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
)

var db *sql.DB

func server(router *gin.Engine) {
	router.Use(AddHeader())
	router.GET("/hosts", getAllHosts)                          // get all hosts
	router.GET("/hosts/:id", getHostByID)                      // get a host by host id
	router.GET("/containers", getAllContainers)                // get all containers
	router.GET("/containers/:id", getContainerByID)            // get a container by container id
	router.GET("/containers/hosts/:id", getContainersByHostID) // get all container with host id
	router.POST("/containers", postContainer)                  // get all container with host id
}

func getAllHosts(gctx *gin.Context) {
	row, err := db.Query("SELECT * FROM hosts")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	listOfHosts := fetchHostRecords(row)
	gctx.JSON(http.StatusOK, gin.H{"Hosts": listOfHosts})
	return
}

func getAllContainers(gctx *gin.Context) {
	row, err := db.Query("SELECT * FROM  containers")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	listOfContainers := fetchContainerRecords(row)
	gctx.JSON(http.StatusOK, gin.H{"Containers": listOfContainers})
	return
}

func getHostByID(gctx *gin.Context) {
	// parse url to get the parameter, id
	key := gctx.Params.ByName("id")

	query, args := "SELECT * FROM hosts WHERE id = ?", key
	row, err := db.Query(query, args)
	if err != nil {
		log.Fatal(err)
	}

	//row, err := db.Query("SELECT * FROM hosts WHERE id = ? ", key)
	//if err != nil {
	//	log.Fatal(err)
	//}

	listOfHosts := fetchHostRecords(row)
	if len(listOfHosts) == 0 {
		gctx.JSON(http.StatusOK, gin.H{"Hosts": "No host found with ID : " + key})
		return
	}
	gctx.JSON(http.StatusOK, gin.H{"Hosts": listOfHosts})
	return
}

func getContainerByID(gctx *gin.Context) {
	// parse url to get the parameter, id
	key := gctx.Params.ByName("id")

	row, err := db.Query("SELECT * FROM containers WHERE id = ?", key)
	if err != nil {
		log.Fatal(err)
	}

	listOfContainers := fetchContainerRecords(row)
	if len(listOfContainers) == 0 {
		gctx.JSON(http.StatusOK, gin.H{"Containers": "No host with that ID was found"})
		return
	}

	gctx.JSON(http.StatusOK, gin.H{"Logs": listOfContainers})
	return
}

func getContainersByHostID(gctx *gin.Context) {
	// parse url to get the parameter, id
	key := gctx.Params.ByName("id")

	row, err := db.Query("SELECT * FROM containers WHERE host_id = ? ", key)
	if err != nil {
		log.Fatal(err)
	}

	listOfContainers := fetchContainerRecords(row)
	if len(listOfContainers) == 0 {
		gctx.JSON(http.StatusOK, gin.H{"Containers": "No host was found with host id :" + key})
		return
	}

	gctx.JSON(http.StatusOK, gin.H{"Logs": listOfContainers})
	return
}

func postContainer(gctx *gin.Context) {
	container := Container{}

	// bind response body to container object
	if err := gctx.BindJSON(&container); err != nil {
		gctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

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

func AddHeader() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		gctx.Header("Access-Control-Allow-Origin", "*")
		gctx.Header("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
		gctx.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		if gctx.Request.Method == "OPTIONS" {
			gctx.AbortWithStatus(204)
		}
		gctx.Next()
	}
}

func sqlQuery(query string, args...interface{}) *sql.Rows{

}

func fetchHostRecords(row *sql.Rows) []interface{} {
	host := Host{}
	listOfObjects := make([]interface{}, 0)

	// Iterate and fetch the records from result cursor
	for row.Next() {
		resp := make(map[string]string)
		err := row.Scan(&host.ID, &host.UUID, &host.Name, &host.IpAddress)
		if err != nil {
			log.Fatal("something went wrong while scanning database: ", err)
		}
		// create map of container attributes
		resp[hostID] = strconv.Itoa(host.ID)
		resp[hostUUID] = host.UUID
		resp[hostName] = host.Name
		resp[hostIPAddress] = host.IpAddress

		listOfObjects = append(listOfObjects, resp)
	}
	return listOfObjects
}

func fetchContainerRecords(row *sql.Rows) []interface{} {
	container := Container{}
	listOfContainers := make([]interface{}, 0)

	// Iterate and fetch the records from result cursor
	for row.Next() {
		resp := make(map[string]interface{})
		err := row.Scan(&container.ID, &container.HostID, &container.Name, &container.ImageName)
		if err != nil {
			log.Fatal("something went wrong while scanning database: ", err)
		}
		// create map of container attributes
		resp[containerID] = container.ID
		resp[containerHostID] = container.HostID
		resp[containerName] = container.Name
		resp[containerImageName] = container.ImageName

		listOfContainers = append(listOfContainers, resp)
	}
	return listOfContainers
}
