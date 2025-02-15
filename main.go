package main

import (
	"fmt"

	"github.com/bevane/gator/internal/config"
)

func main() {
	configFile := config.Read()
	configFile.SetUser("bevane")
	updatedConfigFile := config.Read()
	fmt.Println(updatedConfigFile)
}
