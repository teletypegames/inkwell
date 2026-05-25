package main

import (
	"log"

	"pncdsl/domain"
	"pncdsl/pncdsl"
)

func main() {
	if err := pncdsl.Run(domain.Build()); err != nil {
		log.Fatal(err)
	}
}
