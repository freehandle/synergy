package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/middleware/social"
	"github.com/freehandle/breeze/socket"
	"github.com/freehandle/breeze/util"
	"github.com/freehandle/handles/attorney"
)

const HandlePort = 8000

func GenesisHandlesServer(ctx context.Context, pk crypto.PrivateKey) error {
	genesis := attorney.NewGenesisState(notarypath)
	sender := make(chan []byte)
	listeners := []chan []byte{sender}
	receiver := make(chan []byte, 2)
	IO, err := util.OpenMultiFileStore(".", "blocos_handles")
	if err != nil {
		return err
	}
	defer func() {
		IO.Close()
		log.Println("handles blockchain IO closed")
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
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return err
	}
	ctxCancel, cancel := context.WithCancel(ctx)
	newConnection := make(chan *socket.SignedConnection)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error listening on port 8000:", err)
				cancel()
				return
			}
			signed, err := socket.PromoteConnection(conn, pk, socket.AcceptAllConnections)
			if err != nil {
				newConnection <- signed
			}
		}
	}()

	go func() {
		abertas := make([]*socket.SignedConnection, 0)
		for {
			select {
			case conn := <-newConnection:
				abertas = append(abertas, conn)
				go func() {
					for {
						data, err := conn.Read()
						if err != nil {
							fmt.Println("Error reading from handle client:", err)
							return
						}
						receiver <- data
					}
				}()
			case data := <-sender:
				for i, conn := range abertas {
					if err := conn.Send(data); err != nil {
						fmt.Println("Error writing to handle client:", err)
						abertas = append(abertas[:i], abertas[i+1:]...)
					}
				}
			case <-ctxCancel.Done():
				listener.Close()
				return
			}
		}
	}()
	return <-chain.Start(ctxCancel)
}

func NewHandleConnector(address string, pk crypto.PrivateKey, token crypto.Token) (chan []byte, error) {
	conn, err := socket.Dial("", fmt.Sprintf("localhost:%d", HandlePort), pk, token)
	if err != nil {
		return nil, err
	}
	sender := make(chan []byte, 2)
	go func() {
		for {
			data, err := conn.Read()
			if err != nil {
				fmt.Println("Error reading from handle server:", err)
				return
			}
			sender <- data
		}
	}()

	go func() {
		for {
			data := <-sender
			if err := conn.Send(data); err != nil {
				fmt.Println("Error writing to handle server:", err)
				return
			}
		}
	}()

	return sender, nil

}
