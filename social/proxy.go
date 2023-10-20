package social

import (
	"fmt"
	"log"
	"sync"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/breeze/network/trusted"
	"github.com/lienkolabs/breeze/util"
	"github.com/lienkolabs/synergy/social/state"
)

type Proxy struct {
	mu      sync.Mutex
	state   *state.State
	conn    *trusted.SignedConnection
	viewers []chan uint64
	epoch   uint64
}

func (p *Proxy) Stop() {
	p.conn.Shutdown()
}

func (p *Proxy) State() *state.State {
	return p.state
}

func (p *Proxy) Epoch() uint64 {
	return p.epoch
}

func (p *Proxy) Action(data []byte) {
	undressed := Undress(data)
	if err := p.conn.Send(undressed); err != nil {
		log.Printf("error sending action: %v", err)
	}
}

func (p *Proxy) Register() chan uint64 {
	viewer := make(chan uint64)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.viewers = append(p.viewers, viewer)
	return viewer
}

func SelfProxyState(host string, hostToken crypto.Token, credential crypto.PrivateKey, genesis *state.State) *Proxy {
	conn, err := trusted.Dial(host, credential, hostToken)
	if err != nil {
		log.Fatalf("could not connect to host: %v", err)
	}
	proxy := &Proxy{
		mu:      sync.Mutex{},
		state:   genesis,
		conn:    conn,
		viewers: make([]chan uint64, 0),
		epoch:   0,
	}

	go func() {
		for {
			data, err := conn.Read()
			if err != nil {
				log.Printf("error reading from host: %v", err)
				continue
			}
			if data[0] == 0 {
				if len(data) == 9 {
					proxy.epoch, _ = util.ParseUint64(data, 1)
					proxy.mu.Lock()
					for _, v := range proxy.viewers {
						v <- proxy.epoch
					}
					proxy.mu.Unlock()
				} else {
					log.Print("invalid epoch message")
				}
			} else if data[0] == 1 {
				if len(data) > 1 {
					action := data[1:]
					if err := proxy.state.Action(action); err != nil {
						log.Printf("invalid action: %v", err)
					} else {
						fmt.Println("action received")
					}
				}
			} else if data[0] == 2 {
				blocks := ParseMultiBlocks(data)
				if len(blocks) == 0 {
					log.Printf("invalid multiblocv: %v", err)
				} else {
					log.Printf("multiple blocks: %v", len(blocks))
				}

				for _, block := range blocks {
					proxy.mu.Lock()
					proxy.epoch = block.epoch
					for _, v := range proxy.viewers {
						v <- proxy.epoch
					}
					proxy.mu.Unlock()
					for _, action := range block.actions {
						if err := proxy.state.Action(action); err != nil {
							log.Printf("invalid action: %v", err)
						}
					}
				}
			} else {
				log.Printf("invalid message type: %v", data[0])
			}
		}
	}()
	return proxy
}

type blockdata struct {
	epoch   uint64
	actions [][]byte
}

func ParseMultiBlocks(data []byte) []*blockdata {
	if len(data) < 9 {
		return nil
	}
	blocks := make([]*blockdata, 0)
	position := 1
	for {
		block := blockdata{}
		block.epoch, position = util.ParseUint64(data, position)
		block.actions, position = util.ParseActionsArray(data, position)
		blocks = append(blocks, &block)
		if position >= len(data) {
			return blocks
		}
	}
}
