package main

import (
	"fmt"
	"log"
	"resource-service/cmd/api/bootstrap"
)

func main() {
	err := bootstrap.Run()
	if err != nil {
		log.Fatal(fmt.Printf("Se ha producido un error al arrancar el servidor: %e", err))
	}
}
