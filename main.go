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
	"github.com/freehandle/breeze/middleware/social"
	"github.com/freehandle/breeze/util"
	"github.com/freehandle/handles/attorney"
	"github.com/freehandle/safe"
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

func launchSafeServer(ctx context.Context, cancel context.CancelFunc, pk crypto.PrivateKey, password string, gateway safe.Sender, receive chan []byte) {
	cfg := safe.SafeConfig{
		Credentials: pk,
		HtmlPath:    "../safe/",
		Path:        ".",
		Port:        7000,
		//ServerName:  "/safe",
	}
	errSignal := safe.NewLocalServer(ctx, cfg, password, gateway, receive)
	err := <-errSignal
	log.Printf("error creating safe server: %v", err)
	cancel()
}

func launchSynergyServer(pk crypto.PrivateKey, gateway chan []byte, receive chan []byte, pass string) {
	indexer := index.NewIndex()
	genesis := state.GenesisState(indexer)
	indexer.SetState(genesis)
	attorneySecret := pk
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
		//ServerName:    "/synergy",
	}
	attorney, finalize := api.NewGeneralAttorneyServer(config)
	if attorney == nil {
		err := <-finalize
		log.Printf("error creating attorney: %v", err)
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
	for _, env := range envs {
		if strings.HasPrefix(env, "FREEHANDLE_SECRET=") {
			emailPassword, _ = strings.CutPrefix(env, "FREEHANDLE_SECRET=")
		}
	}

	safeListener := make(chan []byte)
	synergyListener := make(chan []byte)
	sender := make(chan []byte)

	ctxBack := context.Background()
	ctx, cancel := context.WithCancel(ctxBack)

	go launchLocalChain(ctx, []chan []byte{synergyListener, safeListener}, sender)

	_, pk := crypto.RandomAsymetricKey()

	go launchSafeServer(ctx, cancel, pk, emailPassword, ByArraySender(sender), safeListener)
	go launchSynergyServer(pk, sender, synergyListener, emailPassword)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c
	fmt.Println("Got signal:", s)
	cancel()
	time.Sleep(5 * time.Second)
}
