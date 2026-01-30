package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/middleware/simple"
	"github.com/freehandle/breeze/middleware/social"
	"github.com/freehandle/breeze/util"
	"github.com/freehandle/handles/attorney"
	"github.com/freehandle/synergy/api"
	"github.com/freehandle/synergy/config"
	"github.com/freehandle/synergy/network"
	"github.com/freehandle/synergy/social/index"
	"github.com/freehandle/synergy/social/state"
)

const (
	handlesProtocolCode = 1
	breezeProtocolCode  = 0
	notarypath          = ""
	blocksPath          = ""
	blocksName          = "chain"
)

type ByArraySender chan []byte

func (b ByArraySender) Send(data []byte) error {
	b <- data
	return nil
}

func launchLocalChain(ctx context.Context, listeners []chan []byte, receiver chan []byte) error {
	genesis := attorney.NewGenesisState(notarypath)
	IO, err := util.OpenMultiFileStore(".", "blocos")
	if err != nil {
		return err
	}
	defer func() {
		IO.Close()
		log.Println("blockchain IO closed")
	}()

	chain := &social.LocalBlockChain[*attorney.Mutations, *attorney.MutatingState]{
		Interval:  time.Second,
		Listeners: listeners,
		Receiver:  receiver,
		IO:        IO,
	}
	if err = chain.LoadState(genesis, IO, listeners); err != nil {
		return err
	}
	return <-chain.Start(ctx)
}

func launchSynergyServer(gateway chan []byte, receive chan []byte, synergyPass, emailPass string, vault *config.SecretsVault) {
	indexer := index.NewIndex()
	genesis := state.GenesisState(indexer)
	indexer.SetState(genesis)

	attorneySecret := vault.PK
	cookieStore := api.OpenCokieStore("cookies.dat", genesis)
	passwordManager := api.NewFilePasswordManager("passwords.dat")
	config := api.ServerConfig{
		Vault:       vault,
		Attorney:    attorneySecret.PublicKey(),
		Ephemeral:   attorneySecret.PublicKey(),
		Passwords:   passwordManager,
		CookieStore: cookieStore,
		Indexer:     indexer,
		Gateway:     gateway,
		State:       genesis,
		GenesisTime: genesis.GenesisTime,
		Hostname:    "localhost:3000",
		Mail:        &api.SMTPGmail{From: "freemyhandle@gmail.com", Password: emailPass},
		Port:        3000,
		Safe:        8090,
		//ServerName:    "/synergy",
	}
	attorney, finalize := api.NewGeneralAttorneyServer(config)
	if attorney == nil {
		err := <-finalize
		log.Fatalf("error creating attorney: %v\n", err)
		return
	}
	handles := &network.HandlesDB{
		TokenToHandle: make(map[crypto.Token]network.UserInfo),
		HandleToToken: make(map[string]crypto.Token),
		Attorneys:     make(map[crypto.Token]struct{}),
		SynergyApp:    attorneySecret.PublicKey(),
	}
	genesis.Axe = handles
	signal := network.ByteArrayToSignal(receive)
	network.NewSynergyNode(handles, attorney, signal)
}

func main() {

	envs := os.Environ()
	var emailPassword string
	var synergyPassword string
	for _, env := range envs {
		if strings.HasPrefix(env, "FREEHANDLE_SECRET=") {
			emailPassword, _ = strings.CutPrefix(env, "FREEHANDLE_SECRET=")
		} else if strings.HasPrefix(env, "SYNERGY_SECRET=") {
			synergyPassword, _ = strings.CutPrefix(env, "SYNERGY_SECRET=")
		}
	}

	vault, err := config.OpenVaultFromPassword([]byte(synergyPassword), "synergyvault.dat")
	if err != nil {
		log.Fatalf("error opening vault: %v\n", err)
		return
	}
	vault.Secrets[vault.PK.PublicKey()] = vault.PK

	//safeListener := make(chan []byte)
	//synergyListener := make(chan []byte)

	ctxBack := context.Background()
	ctx, cancel := context.WithCancel(ctxBack)

	// HARD CODED ENDERECO DA CHAIN
	synergyListener := simple.DissociateActions(ctx, simple.NewBlockReader(ctx, "/home/lienko/setembro/handles/cmd/proxy-handles", "blocos", time.Second))
	//safeListener := simple.DissociateActions(ctx, simple.NewBlockReader(ctx, "", "blocos", time.Second))

	//breezeToken, _ := crypto.RandomAsymetricKey()
	breezeToken := crypto.TokenFromString("91ad274d06c4be307a332a0e59449ad25ae2c65e4ad5a8f0af87067ac2fc3a54")
	fmt.Println("Using breeze token:", breezeToken.String())
	sender, err := simple.Gateway(ctx, 7000, breezeToken, vault.PK)
	if err != nil {
		log.Fatalf("error creating gateway: %v", err)
	}

	//sender := make(chan []byte)

	//go launchLocalChain(ctx, []chan []byte{synergyListener, safeListener}, sender)

	/* Initialize Safe Server */
	// cfg := safe.SafeConfig{
	// 	Credentials: vault.PK,
	// 	HtmlPath:    "../safe/",
	// 	Path:        ".",
	// 	Port:        7000,
	// 	//ServerName:  "/safe",
	// }

	//errSignal, safe := safe.NewLocalServer(ctx, cfg, synergyPassword, ByArraySender(sender), safeListener)

	// go func() {
	// 	err := <-errSignal
	// 	log.Printf("error creating safe server: %v", err)
	// 	cancel()
	// }()

	/* Initilize Synergy Server */

	go launchSynergyServer(sender, synergyListener, synergyPassword, emailPassword, vault)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	//go NewSafeRestAPI(8000, safe)

	s := <-c
	fmt.Println("Got signal:", s)
	cancel()
	time.Sleep(5 * time.Second)
}

/*
  curl -X POST http://localhost:8000/ -H "Content-Type: application/json" -d '{"handle": "ruben"}'
*/
