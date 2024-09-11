package main

import "log"

func main() {
	store, err := NewJSONFileStorage("users.json")
	if err != nil {
		log.Fatal(err)
	}
	server := NewAPIServer(":3000", store)
	server.Run()
}