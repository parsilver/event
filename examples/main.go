package examples

import (
	"fmt"
	"os"
)

func Run() {
	// Get the example name from arguments or use "simple" as default
	exampleName := "simple"
	if len(os.Args) > 1 {
		exampleName = os.Args[1]
	}

	fmt.Printf("Running example: %s\n\n", exampleName)

	switch exampleName {
	case "simple":
		fmt.Println("=== Simple Example ===")
		SimpleExample()
	case "middleware":
		fmt.Println("=== Middleware Example ===")
		MiddlewareExample()
	default:
		fmt.Printf("Unknown example: %s\n", exampleName)
		fmt.Println("Available examples: simple, middleware")
	}
}