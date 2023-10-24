package network

import (
	"fmt"
	"log"

	"github.com/freehandle/axe/attorney"
	"github.com/freehandle/breeze/util"
	"github.com/freehandle/synergy/api"
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
func NewSynergyNode(axe *AxeDB, attorneyGeneral *api.AttorneyGeneral, signals chan *Signal) {
	for {
		signal := <-signals
		if signal.Signal == 0 {
			epoch, _ := util.ParseUint64(signal.Data, 0)
			if epoch >= 0 {
				attorneyGeneral.SetEpoch(epoch)
			} else {
				//log.Printf("invalid new block epoch: %v", epoch)
			}
		} else if signal.Signal == 1 {
			fmt.Println("tem...")
			if attorney.IsAxeNonVoid(signal.Data) {
				if attorney.Kind(signal.Data) == attorney.GrantPowerOfAttorneyType {
					grant := attorney.ParseGrantPowerOfAttorney(signal.Data)
					if grant != nil {
						if attorneyGeneral.Token.Equal(grant.Attorney) {
							if user, ok := axe.TokenToHandle[grant.Author]; ok {
								fmt.Printf("OPPPPPPPPPPPPA %+v\n", *grant)
								attorneyGeneral.IncorporateGrantPower(user.Handle, grant)
							}
						}
					}
				}
			}
			synergyAction := axe.Incorporate(signal.Data)
			if synergyAction != nil {
				fmt.Println("mensagem")
				attorneyGeneral.Incorporate(signal.Data)
			}
		} else {
			log.Printf("invalid signal: %v", signal.Signal)
		}
	}
}
