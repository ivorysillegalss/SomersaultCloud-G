package main

import (
	"SomersaultCloud/app/somersaultcloud-ipconfig/bootstrap"
	"log"
)

func main() {
	app, err := InitializeApp()
	if err != nil {
		log.Fatal(err.Error())
	}
	bootstrap.RunIpConfig(app)
	log.Println("start successfully")
}
