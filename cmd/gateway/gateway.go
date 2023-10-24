package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/socket"
	"github.com/freehandle/synergy/network"
	"github.com/freehandle/synergy/social/state"
)

func GetState(chain *blockchain) (*state.State, error) {
	genesis := state.GenesisState(nil)
	if genesis == nil {
		return nil, errors.New("could not create genesis state")
	}
	for _, block := range chain.blocks {
		for _, action := range block.data {
			if !IsAxeNonVoid(action) {
				if err := genesis.Action(action); err != nil {
					log.Printf("blockchain has invalid action: %v", err)
				}
			}
		}
	}
	return genesis, nil
}

func IsAxeNonVoid(action []byte) bool {
	if len(action) < 15 {
		return false
	}
	if action[0] != 0 || action[1] != 0 || action[10] != 1 || action[11] != 0 || action[12] != 0 || action[13] != 0 || action[14] == 0 {
		return false
	}
	return true
}

func NewActionsGateway(port int, credentials crypto.PrivateKey, chain *blockchain) (chan socket.Message, error) {
	validate := socket.AcceptAllConnections
	listeners, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, err
	}

	var genesis *state.State
	if genesis, err = GetState(chain); err != nil {
		return nil, err
	}

	pool := make(ConnectionPool)
	incorporate := make(chan *CachedConnection)

	shutDown := make(chan crypto.Token) // receive connection shutdown
	messages := make(chan socket.Message)
	click := time.NewTicker(time.Second)

	go func() {
		for {
			if conn, err := listeners.Accept(); err == nil {
				socketConn, err := socket.PromoteConnection(conn, credentials, validate)
				if err != nil {
					conn.Close()
				} else {
					cached := NewCachedConnection(socketConn)
					incorporate <- cached
					socketConn.Listen(messages, shutDown)
				}
			} else {
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-click.C:
				// next block and broadcast
				chain.NewBlock(pool)
			case cached := <-incorporate:
				pool.Connect(cached)
				// start sync node
				go chain.Sync(cached, len(chain.blocks)-1, len(chain.current.data))
			case token := <-shutDown:
				pool.Drop(token)
			case msg := <-messages:
				axe := IsAxeNonVoid(msg.Data)
				if axe {
					chain.NewAction(msg.Data, pool)
				} else {
					undressed := network.BreezeToSynergy(msg.Data)
					if err := genesis.Action(undressed); err == nil {
						chain.NewAction(undressed, pool)
					} else {
						log.Println(err)
					}
				}
			}
		}
	}()

	return messages, nil
}
