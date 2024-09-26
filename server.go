package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/synergy/api"
	"github.com/freehandle/synergy/config"
	"github.com/freehandle/synergy/social/index"
	"github.com/freehandle/synergy/social/state"
)

var pks []crypto.PrivateKey = []crypto.PrivateKey{
	{118, 35, 197, 163, 215, 20, 35, 190, 110, 151, 246, 231, 86, 177, 156, 89, 122, 69, 28, 233, 185, 150, 126, 169, 237, 173, 83, 120, 145, 238, 242, 137,
		171, 216, 111, 131, 116, 217, 38, 148, 28, 178, 174, 63, 166, 4, 50, 6, 20, 133, 15, 153, 41, 252, 164, 165, 2, 127, 163, 204, 24, 24, 188, 240},
	{152, 224, 227, 154, 131, 1, 186, 147, 73, 37, 4, 253, 11, 148, 195, 67, 86, 85, 28, 162, 78, 239, 168, 42, 204, 222, 144, 41, 186, 246, 250, 57, 125, 202,
		107, 133, 63, 39, 136, 246, 120, 222, 29, 73, 106, 213, 95, 132, 50, 130, 162, 42, 95, 159, 10, 246, 213, 217, 160, 125, 181, 194, 37, 174},
	{125, 86, 238, 128, 237, 4, 143, 47, 214, 72, 71, 47, 72, 45, 214, 45, 178, 98, 105, 154, 171, 151, 73, 183, 234, 120, 128, 38, 174, 253, 105, 162, 189,
		253, 40, 134, 214, 5, 229, 224, 171, 175, 152, 114, 72, 167, 9, 215, 75, 171, 3, 255, 30, 255, 110, 127, 9, 3, 129, 24, 230, 246, 109, 184},
}

var gatewayPK = crypto.PrivateKey{121, 98, 124, 72, 181, 150, 37, 34, 195, 97, 127, 65, 198, 38, 114, 116, 94, 244, 191, 249, 171, 114, 54, 232, 84, 87, 151, 146, 40, 249, 220, 89, 52, 170, 195, 171,
	223, 79, 238, 175, 43, 29, 241, 31, 238, 42, 141, 254, 202, 212, 102, 132, 0, 53, 249, 84, 179, 102, 229, 5, 205, 10, 145, 246}

func server4(pass string) {
	indexer := index.NewIndex()
	genesis := state.GenesisState(indexer)
	indexer.SetState(genesis)
	attorneySecret := pks[0]
	//_, attorneySecret := crypto.RandomAsymetricKey()

	gateway := make(chan []byte)

	vault := &config.SecretsVault{
		Secrets: make(map[crypto.Token]crypto.PrivateKey),
	}
	vault.Secrets[attorneySecret.PublicKey()] = attorneySecret

	cookieStore := api.OpenCokieStore("cookies.dat", genesis)
	passwordManager := api.NewFilePasswordManager("passwords.dat")

	config := api.ServerConfig{
		Vault:         vault,
		Attorney:      attorneySecret.PublicKey(),
		Ephemeral:     attorneySecret.PublicKey(),
		Passwords:     passwordManager,
		CookieStore:   cookieStore,
		Indexer:       indexer,
		Gateway:       gateway,
		State:         genesis,
		GenesisTime:   genesis.GenesisTime,
		EmailPassword: pass,
		Port:          3000,
	}
	attorney, finalize := api.NewGeneralAttorneyServer(config)
	if attorney == nil {
		err := <-finalize
		log.Printf("error creating attorney: %v", err)
		return
	}
	// network.NewProxy("localhost:4100", gatewayPK.PublicKey(), attorneySecret, gateway, attorney)
}

/*func server3(pass string) {
	indexer := index.NewIndex()
	genesis := state.GenesisState(indexer)
	indexer.SetState(genesis)

	_, attorneySecret := crypto.RandomAsymetricKey()

	proxy := social.SelfProxyState("localhost:4100", gatewayPK.PublicKey(), attorneySecret, genesis) // simulador de blockchain
	for n := 0; n < len(pks); n++ {
		//	api.NewAttorneyServer(attorneySecret, pks[n].PublicKey(), 3000+n, proxy, indexer)
		indexer.AddMemberToIndex(pks[n].PublicKey(), fmt.Sprintf("user_%v", n))
	}

	vault := vault.SecureVault{
		Secrets: make(map[crypto.Token]crypto.PrivateKey),
	}
	vault.Secrets[attorneySecret.PublicKey()] = attorneySecret

	cookieStore := api.OpenCokieStore("cookies.dat", genesis)
	passwordManager := api.NewFilePasswordManager("passwords.dat")

	config := api.ServerConfig{
		Vault:         &vault,
		Attorney:      attorneySecret.PublicKey(),
		Ephemeral:     attorneySecret.PublicKey(),
		Gateway:       proxy,
		CookieStore:   cookieStore,
		Passwords:     passwordManager,
		EmailPassword: pass,
		Indexer:       indexer,
		Port:          3000,
	}
	err := <-api.NewGeneralAttorneyServer(config)
	fmt.Println(err)
}

func server2() {

	indexer := index.NewIndex()
	genesis := state.GenesisState(indexer)
	indexer.SetState(genesis)

	_, attorneySecret := crypto.RandomAsymetricKey()

	proxy := social.SelfProxyState("localhost:4100", gatewayPK.PublicKey(), attorneySecret, genesis) // simulador de blockchain
	//proxy := social.SelfProxyState("lienko.com:4100", gatewayPK.PublicKey(), attorneySecret, genesis) // simulador de blockchain
	for n := 0; n < len(pks); n++ {
		api.NewAttorneyServer(attorneySecret, pks[n].PublicKey(), 3000+n, proxy, indexer)
		indexer.AddMemberToIndex(pks[n].PublicKey(), fmt.Sprintf("user_%v", n))
	}
}

func server() {
	N := 3
	users := make(map[crypto.Token]string)
	userToken := make([]crypto.Token, N)
	indexer := index.NewIndex()
	for n := 0; n < N; n++ {
		userToken[n] = pks[n].PublicKey()
		users[userToken[n]] = fmt.Sprintf("user_%v", n)
		indexer.AddMemberToIndex(userToken[n], users[userToken[n]])
	}
	state := social.TestGenesisState(users, indexer)
	indexer.SetState(state)
	gateway := social.SelfGateway(state) // simulador de blockchain

	_, attorneySecret := crypto.RandomAsymetricKey()
	for n := 0; n < N; n++ {
		api.NewAttorneyServer(attorneySecret, userToken[n], 3000+n, gateway, indexer)
	}
}
*/

func createNewServer() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not retrieve USER home dir: %v\n", err)
	}

	var files []fs.DirEntry
	path := filepath.Join(homeDir, ".synergy")
	if files, err = os.ReadDir(path); err != nil {
		if err := os.Mkdir(path, fs.ModePerm); err != nil {
			log.Fatalf("could not create directort: %v\n", err)
		}
		if files, err = os.ReadDir(path); err != nil {
			log.Fatalf("unexpected error: %v\n", err)
		}
	}
	var instructionGateway, protocolGateway string
	if len(files) > 0 {
		return
	}
	token, _ := crypto.RandomAsymetricKey()
	fmt.Printf("You must grant power of attorney to the application key\n%v\n", token)

	fmt.Println("instruction gateway:")
	fmt.Scanln(&instructionGateway)
	fmt.Printf("Ok, connected to %v gateway\n", instructionGateway)
	fmt.Println("protocol gateway:")
	fmt.Scanln(&protocolGateway)
	fmt.Printf("Ok, connected to %v gateway\n", protocolGateway)

}
