package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html"
	"log"
	"net/http"
	"strconv"
)


// X Get all hosts list
//- Get all containers list
//- Get Host by ID
//- Get containers by ID
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

// middleware to fetch from create database object
// middleware for JSON header

func server(){
	http.HandleFunc("/hosts", getAllHosts)           // get all hosts
	http.HandleFunc("/containers", getAllContainers) // get all containers
	http.HandleFunc("/hosts/:hostID", getHostByID)           // get a host by host id
	http.HandleFunc("/containers/test", getContainerByID) // get a container by container id
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// HANDLERS

func getAllHosts(w http.ResponseWriter, r *http.Request) {
	//_, err := fmt.Fprintf(w, "Hello Cheese Pizza, %q", html.EscapeString(r.URL.Path))
	//if err != nil {
	//	panic(err)
	//}

	db, _ := sql.Open("sqlite3", "./aqua.db")
	defer db.Close() // Defer Closing the database

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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	containers := map[string]interface{}{"hosts": listOfObjects}
	jsonResp, err := json.Marshal(containers)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}



func getAllContainers(w http.ResponseWriter, r *http.Request) {

	db, _ := sql.Open("sqlite3", "./aqua.db")
	defer db.Close() // Defer Closing the database

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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	containers := map[string]interface{}{"containers":listOfContainers}
	jsonResp, err := json.Marshal(containers)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

func getHostByID(w http.ResponseWriter, r *http.Request ){

	db, _ := sql.Open("sqlite3", "./aqua.db")
	defer db.Close() // Defer Closing the database

	row, err := db.Query("SELECT * FROM host WHERE id = hostID")
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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	containers := map[string]interface{}{"containers":listOfContainers}
	jsonResp, err := json.Marshal(containers)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

func getContainerByID(w http.ResponseWriter, r *http.Request){
	_, err := fmt.Fprintf(w, "IMPLEMENT ME, %q", html.EscapeString(r.URL.Path))
	if err != nil {
		panic(err)
	}
}

func getContainersByHostID(w http.ResponseWriter, r *http.Request){
	_, err := fmt.Fprintf(w, "IMPLEMENT ME, %q", html.EscapeString(r.URL.Path))
	if err != nil {
		panic(err)
	}
}



func displayHostID(db *sql.DB) {
	row, err := db.Query("SELECT * FROM hosts")
	if err != nil {
		log.Fatal(err)
	}

	//resp := make(map[string]string)
	host := Host{}
	listOfContainers := make([]interface{},0)

	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		resp := make(map[string]string)

		err := row.Scan(&host.id, &host.uuid, &host.name, &host.ipAddress)
		if err != nil {
			log.Fatal("something went wrong while scanning database: ", err)
		}
		resp[hostID] = strconv.Itoa(host.id)
		resp[hostUUID] = host.uuid
		resp[hostName] = host.name
		resp[hostIPAddress] = host.ipAddress

		log.Println(resp)
		listOfContainers = append(listOfContainers, resp)
		log.Println(listOfContainers)
		// log.Println("Pizza: ", host.id, host.uuid, host.name, host.ipAddress)
	}
}


func displayContainers(){
	db, _ := sql.Open("sqlite3", "./aqua.db")
	defer db.Close() // Defer Closing the database

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
	log.Println(listOfContainers)
}

func insertStudent(db *sql.DB, ) {
	log.Println("Inserting student record ...")

	host_id := 2
	name := "nginx"
	image_name := "apple"

	insertStudentSQL := `INSERT INTO containers(host_id, name, image_name ) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(host_id, name, image_name)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func testPathBase(){

	db, _ := sql.Open("sqlite3", "./aqua.db")
	defer db.Close() // Defer Closing the database

	row, err := db.Query("SELECT * FROM  containers WHERE id = 4")
	if err != nil {
		log.Fatal(err)
	}
	container := Container{}

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
		log.Println(resp)

	}



}

func main() {

	// JZ to delete -- for testing only
	//db, _ := sql.Open("sqlite3", "./aqua.db")
	//defer db.Close() // Defer Closing the database
	////displayHostID(db)
	//insertStudent(db)

	//server()
	testPathBase()

}
