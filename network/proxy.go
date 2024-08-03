package network

/*
import (
	"context"
	"log"

	"github.com/freehandle/breeze/consensus/messages"
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/middleware/social"
	"github.com/freehandle/breeze/socket"
	"github.com/freehandle/breeze/util"
	"github.com/freehandle/synergy/api"
)

func LaunchProxy(axeHost, gatewayHost string, axeToken, gatewayToken crypto.Token, credentials crypto.PrivateKey, gateway chan []byte, attorneyGeneral *api.AttorneyGeneral) {

	connGateway, err := socket.Dial("", gatewayHost, credentials, gatewayToken)
	if err != nil {
		log.Fatalf("could not connect to axe host: %v", err)
	}

	//	connAxe, err := socket.Dial(axeHost, credentials, axeToken)
	//	if err != nil {
	//		log.Fatalf("could not connect to axe host: %v", err)
	//	}

	axe := &HandlesDB{
		TokenToHandle: make(map[crypto.Token]UserInfo),
		HandleToToken: make(map[string]crypto.Token),
		Attorneys:     make(map[crypto.Token]struct{}),
		SynergyApp:    credentials.PublicKey(),
	}
	attorneyGeneral.RegisterAxeDataBase(axe)

	signal := make(chan *Signal)

	// get actions incorporate to HandlesDB and foward to attorney
	go NewSynergyNode(axe, attorneyGeneral, signal)

	// get data from host andforwar to sinergy node

	go SocialProtocolProxy(axeHost, axeToken, credentials, 1, signal)
	//go SelfProxyState(connAxe, signal)

	// gateway
	go func() {
		for {
			action := <-gateway
			//undressed := BreezeToSynergy(action)
			action = append([]byte{messages.MsgActionSubmit}, action...)
			if err := connGateway.Send(action); err != nil {
				log.Printf("error sending action: %v", err)
			}
		}
	}()

}

func NewProxy(host string, token crypto.Token, credentials crypto.PrivateKey, gateway chan []byte, attorneyGeneral *api.AttorneyGeneral) {
	conn, err := socket.Dial("", host, credentials, token)
	if err != nil {
		log.Fatalf("could not connect to host: %v", err)
	}
	axe := &HandlesDB{
		TokenToHandle: make(map[crypto.Token]UserInfo),
		HandleToToken: make(map[string]crypto.Token),
		Attorneys:     make(map[crypto.Token]struct{}),
		SynergyApp:    credentials.PublicKey(),
	}
	signal := make(chan *Signal)

	// get actions incorporate to HandlesDB and foward to attorney
	go NewSynergyNode(axe, attorneyGeneral, signal)

	// get data from host andforwar to sinergy node
	go SelfProxyState(conn, signal)

	// gateway
	go func() {
		for {
			action := <-gateway
			//undressed := BreezeToSynergy(action)
			if err := conn.Send(action); err != nil {
				log.Printf("error sending action: %v", err)
			}
		}
	}()
}

func SocialProtocolProxy(address string, token crypto.Token, credentials crypto.PrivateKey, epoch uint64, signal chan *Signal) {
	ctx := context.Background()
	listener := social.SocialProtocolBlockListener(ctx, address, token, credentials, epoch)
	for {
		block := <-listener
		if block != nil {

			signal <- &Signal{
				Signal: 0,
				Data:   util.Uint64ToBytes(block.Epoch),
			}
			for _, action := range block.Actions {
				signal <- &Signal{
					Signal: 1,
					Data:   action,
				}
			}
		}
	}
}

func SelfProxyState(conn *socket.SignedConnection, signal chan *Signal) {
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
				//log.Printf("invalid multiblocv: %v", err)
			} else {
				//log.Printf("multiple blocks: %v", len(blocks))
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
*/
