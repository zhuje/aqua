package main

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
