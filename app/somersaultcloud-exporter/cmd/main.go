package main

import "SomersaultCloud/app/somersaultcloud-exporter/api/route"

func main() {
	_, err := InitializeApp()
	if err != nil {
		return
	}
	route.Setup()
}
