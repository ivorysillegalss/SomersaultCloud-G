package main

func main() {
	app, err := InitializeApp()
	app.Setup()
	if err != nil {
		return
	}
}
