package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/breeze/network/trusted"
	"github.com/lienkolabs/synergy/social/state"
)

func GetState(chain *blockchain) (*state.State, error) {
	genesis := state.GenesisState(nil)
	if genesis == nil {
		return nil, errors.New("could not create genesis state")
	}
	for _, block := range chain.blocks {
		for _, action := range block.data {
			if err := genesis.Action(action); err != nil {
				return nil, fmt.Errorf("blockchain has invalid action: %v", err)
			}

		}
	}
	return genesis, nil
}

func NewActionsGateway(port int, credentials crypto.PrivateKey, chain *blockchain) (chan trusted.Message, error) {
	validate := trusted.AcceptAllConnections
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
	messages := make(chan trusted.Message)
	click := time.NewTicker(time.Second)

	go func() {
		for {
			if conn, err := listeners.Accept(); err == nil {
				trustedConn, err := trusted.PromoteConnection(conn, credentials, validate)
				if err != nil {
					conn.Close()
				} else {
					cached := NewCachedConnection(trustedConn)
					incorporate <- cached
					trustedConn.Listen(messages, shutDown)
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
				if err := genesis.Action(msg.Data); err == nil {
					// incorporate to chain and broadcast
					chain.NewAction(msg.Data, pool)
				} else {
					log.Println(err)
				}
			}
		}
	}()

	return messages, nil
}
