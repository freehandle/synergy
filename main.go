package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	envs := os.Environ()
	var emailPassword string
	for _, env := range envs {
		if strings.HasPrefix(env, "FREEHANDLE_SECRET=") {
			emailPassword, _ = strings.CutPrefix(env, "FREEHANDLE_SECRET=")
			fmt.Println(emailPassword)
		}
	}

	server4(emailPassword)
	for true {

	}
}
