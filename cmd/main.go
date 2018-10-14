package main

import (
	"tech-db-server/app"
	"log"
)

func main() {
	err := app.ListenAndServe(5000)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

}