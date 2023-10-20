package network

import (
	"log"

	"github.com/lienkolabs/breeze/util"
	"github.com/lienkolabs/synergy/api"
)

type Gateway interface {
	SendAction(action []byte) error
}

type SynergyNode struct {
	Gateway Gateway
	Axe     *AxeDB
	General *api.AttorneyGeneral
}

// 0 = new block
// 1 = action?

type Signal struct {
	Signal byte
	Data   []byte
}

// canal é um primitivo de sincronia
// canal <- value manda para o canal
// value = <-canal recebe do canal
func NewSynergyNode(axe *AxeDB, attorney *api.AttorneyGeneral, signals chan *Signal) {
	for {
		signal := <-signals
		if signal.Signal == 0 {
			epoch, _ := util.ParseUint64(signal.Data, 0)
			if epoch >= 0 {
				attorney.SetEpoch(epoch)
			} else {
				//log.Printf("invalid new block epoch: %v", epoch)
			}
		} else if signal.Signal == 1 {
			synergyAction := axe.Incorporate2(signal.Data)
			if synergyAction != nil {
				attorney.Incorporate(signal.Data)
			}
		} else {
			log.Printf("invalid signal: %v", signal.Signal)
		}
	}
}
