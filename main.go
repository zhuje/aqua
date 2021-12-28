package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	//"io/ioutil"
	"log"
	"net/http"
	"strconv"
)


// X Get all hosts list
// X Get all containers list
// X Get Host by ID
// X Get containers by ID
// X Get all containers for specific host
// X Create a new container in the database via API with json request format (example below)


const hostID string = "id"
const hostUUID string = "uuid"
const hostName string = "name"
const hostIPAddress string = "ipAddress"

type Host struct {
	id int
	uuid string
	name string
	ipAddress string
}

//id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
//host_id INTEGER NOT NULL,
//name TEXT NOT NULL,
//image_name INTEGER NOT NULL, // NOTE ROEE said change to TEXT 12/23 email
//CONSTRAINT containers_FK FOREIGN KEY (host_id) REFERENCES hosts(id)

const containerID string = "id"
const containerHostID string = "host_id"
const containerName string = "name"
const containerImageName string = "image_name"

type Container struct {
	id int
	host_id int
	name string
	image_name string
}

type Cheese struct {
	HostID int `json:"host_id"`
	Name  string `json:"name"`
	ImageName string `json:"image_name"`
}

var db *sql.DB

// middleware to fetch from create database object
// middleware for JSON header

func server(router *gin.Engine){

	router.Use(AddHeader())


	router.GET("/hosts", getAllHosts)           // get all hosts
	router.GET("/hosts/:id", getHostByID)   // get a host by host id

	router.GET("/containers", getAllContainers) // get all containers
	router.GET("/containers/:id", getContainerByID) // get a container by container id
	router.GET("/containers/hosts/:id", getContainersByHostID) // get all container with host id
	router.POST("/cheese", postContainer) // get all container with host id

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

// HANDLERS
func getAllHosts(gctx *gin.Context) {

	row, err := db.Query("SELECT * FROM hosts")
	if err != nil {
		log.Fatal(err)
	}
	host := Host{}
	listOfObjects := make([]interface{},0)

	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		resp := make(map[string]string)
		err := row.Scan(&host.id, &host.uuid, &host.name, &host.ipAddress)
		if err != nil {
			log.Fatal("something went wrong while scanning database: ", err)
		}
		// create map of container attributes
		resp[hostID] = strconv.Itoa(host.id)
		resp[hostUUID] = host.uuid
		resp[hostName] = host.name
		resp[hostIPAddress] = host.ipAddress
		// append this to the listOfObjects
		listOfObjects = append(listOfObjects, resp)
	}

	gctx.JSON(http.StatusOK, gin.H{"Hosts": listOfObjects})
	return
}

func getAllContainers(gctx *gin.Context) {
	row, err := db.Query("SELECT * FROM  containers")
	if err != nil {
		log.Fatal(err)
	}
	container := Container{}
	listOfContainers := make([]interface{},0)

	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		resp := make(map[string]interface{})
		err := row.Scan(&container.id, &container.host_id, &container.name, &container.image_name)
		if err != nil {
			log.Fatal("something went wrong while scanning database: ", err)
		}
		// create map of container attributes
		resp[containerID] = container.id
		resp[containerHostID] = container.host_id
		resp[containerName] = container.name
		resp[containerImageName] = container.image_name
		// append this to the listOfContainers
		listOfContainers = append(listOfContainers, resp)
	}
	gctx.JSON(http.StatusOK, gin.H{"Containers": listOfContainers})
	return
}

func getHostByID(gctx *gin.Context){

	key := gctx.Params.ByName("id")

	row, err := db.Query("SELECT * FROM hosts WHERE id = ? ", key)
	if err != nil {
		log.Fatal(err)
	}

	host := Host{}
	listOfObjects := make([]interface{},0)

	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		resp := make(map[string]string)
		err := row.Scan(&host.id, &host.uuid, &host.name, &host.ipAddress)
		if err != nil {
			log.Fatal("something went wrong while scanning database: ", err)
		}
		// create map of container attributes
		resp[hostID] = strconv.Itoa(host.id)
		resp[hostUUID] = host.uuid
		resp[hostName] = host.name
		resp[hostIPAddress] = host.ipAddress
		// append this to the listOfObjects
		listOfObjects = append(listOfObjects, resp)
	}

	if len(listOfObjects) == 0 {
		gctx.JSON(http.StatusOK, gin.H{"Hosts": "No host found with ID : " + key })
		return
	}
	gctx.JSON(http.StatusOK, gin.H{"Hosts": listOfObjects})
	return
}

func getContainerByID(gctx *gin.Context){
	key := gctx.Params.ByName("id")

	row, err := db.Query("SELECT * FROM containers WHERE id = ?", key)
	if err != nil {
		log.Fatal(err)
	}

	container := Container{}
	listOfContainers := make([]interface{},0)

	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		resp := make(map[string]interface{})
		err := row.Scan(&container.id, &container.host_id, &container.name, &container.image_name)
		if err != nil {
			log.Fatal("something went wrong while scanning database: ", err)
		}
		// create map of container attributes
		resp[containerID] = container.id
		resp[containerHostID] = container.host_id
		resp[containerName] = container.name
		resp[containerImageName] = container.image_name
		// append this to the listOfContainers
		listOfContainers = append(listOfContainers, resp)
	}

	if len(listOfContainers) == 0 {
		gctx.JSON(http.StatusOK, gin.H{"Containers": "No host with that ID was found"})
		return
	}

	gctx.JSON(http.StatusOK, gin.H{"Logs": listOfContainers})
	return
}

func getContainersByHostID(gctx *gin.Context){
	key := gctx.Params.ByName("id")

	row, err := db.Query("SELECT * FROM containers WHERE host_id = ? ", key)
	if err != nil {
		log.Fatal(err)
	}

	container := Container{}
	listOfContainers := make([]interface{},0)

	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		resp := make(map[string]interface{})
		err := row.Scan(&container.id, &container.host_id, &container.name, &container.image_name)
		if err != nil {
			log.Fatal("something went wrong while scanning database: ", err)
		}
		// create map of container attributes
		resp[containerID] = container.id
		resp[containerHostID] = container.host_id
		resp[containerName] = container.name
		resp[containerImageName] = container.image_name
		// append this to the listOfContainers
		listOfContainers = append(listOfContainers, resp)
	}

	if len(listOfContainers) == 0 {
		gctx.JSON(http.StatusOK, gin.H{"Containers": "No host was found with host id :" + key})
		return
	}

	gctx.JSON(http.StatusOK, gin.H{"Logs": listOfContainers})
	return
}

func postContainer(gctx *gin.Context){
	var cheese Cheese
	if err:= gctx.BindJSON(&cheese);err!=nil{
		gctx.AbortWithError(http.StatusBadRequest,err)
		return
	}
	fmt.Println(cheese)

	insertStudentSQL := `INSERT INTO containers(host_id, name, image_name) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		gctx.JSON(http.StatusInternalServerError, gin.H{"Error": "An error occurred while preparing to insert into database."})
		return
	}
	_, err = statement.Exec(cheese.HostID, cheese.Name, cheese.ImageName)
	if err != nil {
		gctx.JSON(http.StatusInternalServerError, gin.H{"Error": "An error occurred while trying to insert into database."})
		return
	}

	gctx.JSON(http.StatusCreated, gin.H{"Container Created": cheese})
	return
}



func main() {
	db, _ = sql.Open("sqlite3", "./aqua.db")
	defer db.Close() // Defer Closing the database

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	//router := gin.Default() // default allows you too see console print outs
	server(router)
	router.Run()

}
