//cmd/collyclicker/main.go

package main

import (
	"log"

	"github.com/lukusbeaur/collyclicker-core/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("error: %v", err)
	}
}
