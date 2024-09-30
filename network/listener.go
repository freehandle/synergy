package network

import (
	"log"

	"github.com/freehandle/breeze/util"
	"github.com/freehandle/handles/attorney"
	"github.com/freehandle/synergy/api"
)

type Gateway interface {
	SendAction(action []byte) error
}

type SynergyNode struct {
	Gateway Gateway
	Axe     *HandlesDB
	General *api.AttorneyGeneral
}

// 0 = new block
// 1 = action?

type Signal struct {
	Signal byte
	Data   []byte
}

func ByteArrayToSignal(receive chan []byte) chan *Signal {
	signals := make(chan *Signal)
	go func() {
		for {
			action, ok := <-receive
			if !ok {
				return
			}
			if len(action) == 8 {
				signals <- &Signal{Signal: 0, Data: action}
				continue
			}
			signals <- &Signal{Signal: 1, Data: action}
		}
	}()
	return signals
}

// canal é um primitivo de sincronia
// canal <- value manda para o canal
// value = <-canal recebe do canal
func NewSynergyNode(axe *HandlesDB, attorneyGeneral *api.AttorneyGeneral, signals chan *Signal) {
	for {
		signal := <-signals
		if signal.Signal == 0 {
			epoch, _ := util.ParseUint64(signal.Data, 0)
			attorneyGeneral.SetEpoch(epoch)
		} else if signal.Signal == 1 {
			if attorney.IsAxeNonVoid(signal.Data) {
				if attorney.Kind(signal.Data) == attorney.GrantPowerOfAttorneyType {
					grant := attorney.ParseGrantPowerOfAttorney(signal.Data)
					if grant != nil {
						if attorneyGeneral.Token.Equal(grant.Attorney) {
							if user, ok := axe.TokenToHandle[grant.Author]; ok {
								attorneyGeneral.IncorporateGrantPower(user.Handle, grant)
							}
						}
					}
				} else if attorney.Kind(signal.Data) == attorney.RevokePowerOfAttorneyType {
					revoke := attorney.ParseGrantPowerOfAttorney(signal.Data)
					if revoke != nil {
						if attorneyGeneral.Token.Equal(revoke.Attorney) {
							if user, ok := axe.TokenToHandle[revoke.Author]; ok {
								attorneyGeneral.IncorporateRevokePower(user.Handle)
							}
						}
					}
				} else if attorney.Kind(signal.Data) == attorney.JoinNetworkType {
					join := attorney.ParseJoinNetwork(signal.Data)
					if join != nil {
						axe.IncorporateJoin(signal.Data)
					}
				} else if attorney.Kind(signal.Data) == attorney.UpdateInfoType {
					axe.IncorporateUpdate(signal.Data)
				}
			}
			synergyAction := axe.Incorporate(signal.Data)
			if synergyAction != nil {
				action := BreezeToSynergy(signal.Data)
				attorneyGeneral.Incorporate(action)
			}
		} else {
			log.Printf("invalid signal: %v", signal.Signal)
		}
	}
}
