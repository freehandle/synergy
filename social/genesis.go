package social

import (
	"sync"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/synergy/social/index"
	"github.com/freehandle/synergy/social/state"
)

// Cria os usuários teste numa instância do estado
func TestGenesisState(users map[crypto.Token]string, indexer *index.Index) *state.State {
	genesis := state.GenesisState(indexer)
	for user, handle := range users {
		genesis.Members[crypto.HashToken(user)] = handle
		genesis.MembersIndex[handle] = user
	}
	return genesis
}

type Gatewayer interface {
	Stop()
	Action([]byte)
	Register() chan uint64
	State() *state.State
}

type Gateway struct {
	mu       *sync.Mutex        // bloqueia ações simultâneas - só deixa passar uma pessoa alterando por vez
	incoming chan []byte        //canal para receber as mensagens
	newBlock []chan uint64      // canal para informar que se forma novo bloco
	stop     chan chan struct{} // fechar o gateway
	state    *state.State       // estado da blockchain
}

func (g *Gateway) State() *state.State {
	return g.state
}

func (g *Gateway) Stop() {
	resp := make(chan struct{})
	g.stop <- resp
	<-resp
}

func (g *Gateway) Action(data []byte) {
	g.incoming <- data
}

func (g *Gateway) Register() chan uint64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	blockEvent := make(chan uint64)
	g.newBlock = append(g.newBlock, blockEvent)
	return blockEvent
}

// seria mais ou menos um nó da rede
// tem o routing, tem a validacao, distribuicao (3 pecas da matriz de funcionalidades independentes)
func SelfGateway(engine *state.State) *Gateway {
	gateway := &Gateway{
		mu:       &sync.Mutex{},
		incoming: make(chan []byte), // criando os canais
		newBlock: make([]chan uint64, 0),
		stop:     make(chan chan struct{}),
		state:    engine,
	}

	ticker := time.NewTicker(time.Second)
	// manda um sinal a cada segundo pra simular o bloco

	// usa go para executar isso em paralelo - abre um processo só pra essa funcao
	go func() {
		// fica apenas ouvindo os canais
		for {
			select {

			// novo bloco
			case <-ticker.C:
				gateway.mu.Lock()
				// nesse tempo só eu posso alterar
				engine.Epoch += 1
				// atualiza a epoch e avisa todos que foi atualizado
				for _, emit := range gateway.newBlock {
					emit <- engine.Epoch
				}
				// destrava
				gateway.mu.Unlock()

			// chegou uma acao pra processar
			case action := <-gateway.incoming:
				// tira tudo que é do breeze e axé e fica só com o que diz respeito ao synergy
				undressed := Undress(action)
				// engine (que é o estado) incorpora acao, se nao conseguir devolve o erro
				if err := engine.Action(undressed); err != nil {
				} else {
				}
			// fechar o processo
			case resp := <-gateway.stop:
				gateway.mu.Lock()
				defer gateway.mu.Unlock()
				for _, event := range gateway.newBlock {
					close(event)
				}
				resp <- struct{}{}
				// unico return aqui que sai do for
				return
			}
		}
	}()
	return gateway
}

// removendo tudo que nao é do synergy
func Undress(data []byte) []byte {
	// ignore first byte (breeze version)
	head := data[1 : 8+crypto.TokenSize+1]
	// ignore protocol id and tail
	// tail = wallet signature + fee + wallet + attorney signature + attorney
	tailSize := 2*crypto.SignatureSize + 2*crypto.TokenSize + 8
	return append(head, data[8+crypto.TokenSize+1+4:len(data)-tailSize]...)
}
