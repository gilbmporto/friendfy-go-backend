package main

import (
	"fmt"
	"friendfy-api/src/config"
	"friendfy-api/src/router"
	"log"
	"net/http"
)

func main() {
	config.Setup()
	r := router.Generate()

	fmt.Printf("Running on port %d...\n", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
