package lumino_node_client

import (
	"fmt"
	"log"

	"example.com/lumino-node-client/client"
	"example.com/lumino-node-client/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	luminoClient, err := client.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create Lumino client: %v", err)
	}

	if err := luminoClient.Run(); err != nil {
		log.Fatalf("Lumino client error: %v", err)
	}

	fmt.Println("Lumino Node Client is running...")
	// Add code here to keep the application running and handle shutdown gracefully
}
