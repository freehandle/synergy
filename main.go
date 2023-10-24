package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	/*token, _ := crypto.RandomAsymetricKey()
	signin := actions.Signin{
		Epoch:   10,
		Author:  token,
		Reasons: "I am the best",
	}
	data := network.SynergyToBreeze(signin.Serialize())
	network.BreezeToSynergy(data)
	*/

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
