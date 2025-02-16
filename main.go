package main

import (
	"fmt"

	"github.com/bevane/gator/internal/config"
)

func main() {
	configFile, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config file: %v", err)
		return
	}
	err = configFile.SetUser("bevane")
	if err != nil {
		fmt.Printf("Error setting user in config: %v", err)
		return
	}
	updatedConfigFile, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config file: %v", err)
		return
	}
	fmt.Println(updatedConfigFile)
}
