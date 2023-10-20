package network

import (
	"log"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/breeze/network/trusted"
	"github.com/lienkolabs/breeze/util"
	"github.com/lienkolabs/synergy/api"
)

type Proxy struct {
	conn *trusted.SignedConnection
}

func (p *Proxy) Stop() {
	p.conn.Shutdown()
}

func (p *Proxy) Action(data []byte) {
	undressed := Undress(data)
	if err := p.conn.Send(undressed); err != nil {
		log.Printf("error sending action: %v", err)
	}
}

func NewProxy(host string, token crypto.Token, credentials crypto.PrivateKey, gateway chan []byte, attorney *api.AttorneyGeneral) {
	conn, err := trusted.Dial(host, credentials, token)
	if err != nil {
		log.Fatalf("could not connect to host: %v", err)
	}
	axe := &AxeDB{
		TokenToHandle: make(map[crypto.Token]UserInfo),
		HandleToToken: make(map[string]crypto.Token),
		Attorneys:     make(map[crypto.Token]struct{}),
		SynergyApp:    credentials.PublicKey(),
	}
	signal := make(chan *Signal)

	// get actions incorporate to axedb and foward to attorney
	go NewSynergyNode(axe, attorney, signal)

	// get data from host andforwar to sinergy node
	go SelfProxyState(conn, signal)

	// gateway
	go func() {
		for {
			action := <-gateway
			undressed := Undress(action)
			if err := conn.Send(undressed); err != nil {
				log.Printf("error sending action: %v", err)
			}
		}
	}()
}

func SelfProxyState(conn *trusted.SignedConnection, signal chan *Signal) {
	for {
		data, err := conn.Read()
		if err != nil {
			log.Printf("error reading from host: %v", err)
			continue
		}
		if data[0] == 0 {
			if len(data) == 9 {
				signal <- &Signal{
					Signal: 0,
					Data:   data[1:],
				}
			} else {
				log.Print("invalid epoch message")
			}
		} else if data[0] == 1 {
			if len(data) > 1 {
				signal <- &Signal{
					Signal: 1,
					Data:   data[1:],
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
				epochBytes := make([]byte, 8)
				util.PutUint64(block.epoch, &epochBytes)
				signal <- &Signal{
					Signal: 0,
					Data:   epochBytes,
				}
				for _, action := range block.actions {
					signal <- &Signal{
						Signal: 1,
						Data:   action,
					}
				}
			}
		} else {
			log.Printf("invalid message type: %v", data[0])
		}
	}
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

func Undress(data []byte) []byte {
	// ignore first byte (breeze version)
	head := data[1 : 8+crypto.TokenSize+1]
	// ignore protocol id and tail
	// tail = wallet signature + fee + wallet + attorney signature + attorney
	tailSize := 2*crypto.SignatureSize + 2*crypto.TokenSize + 8
	return append(head, data[8+crypto.TokenSize+1+4:len(data)-tailSize]...)
}

// como estão no breeze
// 0 (byte) (Versão do Breeze) | 1 (exclsive)
// 0 (byte) (Void do Breeze)   | 2
// Epoch (8 bytes)             | 10
// Author (32 bytes)           | 32
// 0 (Void do Axé)             | 33
// Data .... (info do synergy) | Varia
// Signer (32 bytes)
// Signature (64 bytes)
// Wallet (32 bytes)
// Fee (8 bytes)
// Signature (64 bytes)

// como estão no synergy
// 0  (1 byte)
// Epoch (8 bytes)
// Author (32 bytes)
// ActionKind (1Byte)
// varia com a instrução (X bytes)

const breezeTailSize = 2*crypto.SignatureSize + 8 + 2*crypto.Size

func BreezeVoidToSynergy(action []byte) []byte {
	if len(action) < 33+breezeTailSize {
		return nil
	}
	if action[0] != 0 || action[1] != 0 || action[32] != 0 {
		return nil
	}

	translated = action[2 : 10+32]
	translated = append(translated, action[33:len(action)-breezeTailSize]...)

}
