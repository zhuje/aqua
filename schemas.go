package main

const hostID string = "id"
const hostUUID string = "uuid"
const hostName string = "name"
const hostIPAddress string = "ipAddress"

const containerID string = "id"
const containerHostID string = "host_id"
const containerName string = "name"
const containerImageName string = "image_name"

type Host struct {
	ID        int    `json:"id"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	IpAddress string `json:"ipAddress"`
}

type Container struct {
	ID        int    `json:"id"`
	HostID    int    `json:"host_id"`
	Name      string `json:"name"`
	ImageName string `json:"image_name"`
}
