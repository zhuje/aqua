package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
)


// X Get all hosts list
// X Get all containers list
// X Get Host by ID
// X Get containers by ID
//- Get all containers for specific host


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

var db *sql.DB

// middleware to fetch from create database object
// middleware for JSON header

func server(router *gin.Engine){

	router.Use(AddHeader())


	router.GET("/hosts", getAllHosts)           // get all hosts
	router.GET("/containers", getAllContainers) // get all containers
	router.GET("/hosts/:id", getHostByID)   // get a host by host id
	router.GET("/containers/:id", getContainerByID) // get a container by container id
	router.GET("/containers/:host", getContainersByHostID) // get all container with host id

	// log.Fatal(http.ListenAndServe(":8081", nil))
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

	gctx.JSON(http.StatusOK, gin.H{"Containers": listOfObjects})
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
	gctx.JSON(http.StatusOK, gin.H{"Logs": listOfObjects})
	return
}

func getContainerByID(gctx *gin.Context){
	key := gctx.Params.ByName("id")

	row, err := db.Query("SELECT * FROM containers WHERE id = ? ", key)
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
	key := gctx.Params.ByName("hostID")

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
		gctx.JSON(http.StatusOK, gin.H{"Containers": "No host with that ID was found"})
		return
	}

	gctx.JSON(http.StatusOK, gin.H{"Logs": listOfContainers})
	return
}



//func displayHostID(db *sql.DB) {
//	row, err := db.Query("SELECT * FROM hosts")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	//resp := make(map[string]string)
//	host := Host{}
//	listOfContainers := make([]interface{},0)
//
//	defer row.Close()
//	for row.Next() { // Iterate and fetch the records from result cursor
//		resp := make(map[string]string)
//
//		err := row.Scan(&host.id, &host.uuid, &host.name, &host.ipAddress)
//		if err != nil {
//			log.Fatal("something went wrong while scanning database: ", err)
//		}
//		resp[hostID] = strconv.Itoa(host.id)
//		resp[hostUUID] = host.uuid
//		resp[hostName] = host.name
//		resp[hostIPAddress] = host.ipAddress
//
//		log.Println(resp)
//		listOfContainers = append(listOfContainers, resp)
//		log.Println(listOfContainers)
//		// log.Println("Pizza: ", host.id, host.uuid, host.name, host.ipAddress)
//	}
//}
//
//
//func displayContainers(){
//	db, _ := sql.Open("sqlite3", "./aqua.db")
//	defer db.Close() // Defer Closing the database
//
//	row, err := db.Query("SELECT * FROM  containers")
//	if err != nil {
//		log.Fatal(err)
//	}
//	container := Container{}
//	listOfContainers := make([]interface{},0)
//
//	defer row.Close()
//	for row.Next() { // Iterate and fetch the records from result cursor
//		resp := make(map[string]interface{})
//		err := row.Scan(&container.id, &container.host_id, &container.name, &container.image_name)
//		if err != nil {
//			log.Fatal("something went wrong while scanning database: ", err)
//		}
//		// create map of container attributes
//		resp[containerID] = container.id
//		resp[containerHostID] = container.host_id
//		resp[containerName] = container.name
//		resp[containerImageName] = container.image_name
//
//		// append this to the listOfContainers
//		listOfContainers = append(listOfContainers, resp)
//	}
//	log.Println(listOfContainers)
//}
//
//func insertStudent(db *sql.DB, ) {
//	log.Println("Inserting student record ...")
//
//	host_id := 2
//	name := "nginx"
//	image_name := "apple"
//
//	insertStudentSQL := `INSERT INTO containers(host_id, name, image_name ) VALUES (?, ?, ?)`
//	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
//	// This is good to avoid SQL injections
//	if err != nil {
//		log.Fatalln(err.Error())
//	}
//	_, err = statement.Exec(host_id, name, image_name)
//	if err != nil {
//		log.Fatalln(err.Error())
//	}
//}
//
//func testPathBase(){
//
//	db, _ := sql.Open("sqlite3", "./aqua.db")
//	defer db.Close() // Defer Closing the database
//
//	row, err := db.Query("SELECT * FROM  containers WHERE id = 4")
//	if err != nil {
//		log.Fatal(err)
//	}
//	container := Container{}
//
//	defer row.Close()
//	for row.Next() { // Iterate and fetch the records from result cursor
//		resp := make(map[string]interface{})
//		err := row.Scan(&container.id, &container.host_id, &container.name, &container.image_name)
//		if err != nil {
//			log.Fatal("something went wrong while scanning database: ", err)
//		}
//		// create map of container attributes
//		resp[containerID] = container.id
//		resp[containerHostID] = container.host_id
//		resp[containerName] = container.name
//		resp[containerImageName] = container.image_name
//
//		// append this to the listOfContainers
//		log.Println(resp)
//
//	}
//}
//
//func testServer(router *gin.Engine){
//	router.GET("/containers/:id", func(c *gin.Context) {
//		id := c.Params.ByName("id")
//
//		c.JSON(200, gin.H{
//			"message": id,
//		})
//	})
//	router.Run()
//
//}
//
//func testPath2(gctx *gin.Context){
//	param := gctx.Params.ByName("id")
//	log.Print(param)
//
//}
//
//func initDatabase(){
//
//}

func main() {

	db, _ = sql.Open("sqlite3", "./aqua.db")
	defer db.Close() // Defer Closing the database

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	server(router)
	router.Run()


	// JZ to delete -- for testing only
	//db, _ := sql.Open("sqlite3", "./aqua.db")
	//defer db.Close() // Defer Closing the database
	////displayHostID(db)
	//insertStudent(db)

	//server()

}
