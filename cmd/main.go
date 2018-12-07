package main

import (
	"log"
	"tech-db-server/app"
)

func main() {
	err := app.ListenAndServe(5000)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

}
