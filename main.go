package main

import (
	"context"
	"os"
	"strings"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/middleware/social"
	"github.com/freehandle/breeze/util"
	"github.com/freehandle/handles/attorney"
)

const (
	handlesProtocolCode = 1
	breezeProtocolCode  = 0
	notarypath          = ""
	blocksPath          = ""
	blocksName          = "chain"
)

func launchGenesis(ctx context.Context, credentials crypto.PrivateKey, listeners []chan []byte) (*social.LocalBlockChain[*attorney.Mutations, *attorney.MutatingState], error) {
	genesis := attorney.NewGenesisState(notarypath)
	IO, err := util.OpenMultiFileStore(blocksPath, blocksName)
	if err != nil {
		return nil, err
	}
	cfg := &social.LocalChainConfig[*attorney.Mutations, *attorney.MutatingState]{
		Credentials:  credentials,
		ProtocolCode: handlesProtocolCode,
		Interval:     1,
		Listeners:    listeners,
		Genesis:      genesis,
	}
	local, err := social.OpenChain[*attorney.Mutations, *attorney.MutatingState](IO, cfg)
	return local, err
}

func main() {

	/*token, _ := crypto.RandomAsymetricKey()
	signin := actions.Signin{
		Epoch:   10,
		Author:  token,
		Reasons: "I am the best",
	}
	data := network.SynergyToBreeze(signin.Serialize())
	network.BreezeToSynergy(data)
	*/

	envs := os.Environ()
	var emailPassword string
	for _, env := range envs {
		if strings.HasPrefix(env, "FREEHANDLE_SECRET=") {
			emailPassword, _ = strings.CutPrefix(env, "FREEHANDLE_SECRET=")
		}
	}

	server4(emailPassword)
	for true {

	}
}
