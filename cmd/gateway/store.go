package main

import (
	"log"
	"os"
	"sync"

	"github.com/lienkolabs/breeze/util"
)

type block struct {
	data [][]byte
}

func multiblock(startEpoch uint64, blocks []*block) []byte {
	bytes := []byte{2}
	for n, block := range blocks {
		util.PutUint64(startEpoch+uint64(n), &bytes)
		util.PutActionsArray(block.data, &bytes)
	}
	return bytes
}

type blockchain struct {
	mu      sync.Mutex
	io      *os.File
	blocks  []*block
	current *block
}

// sync a new connection
func (b *blockchain) Sync(conn *CachedConnection, epoch, actionCount int) {
	b.mu.Lock()
	currentBlockCache := make([][]byte, actionCount)
	for n := 0; n < actionCount; n++ {
		currentBlockCache[n] = b.blocks[epoch].data[n]
	}
	b.mu.Unlock()
	for n := 0; n <= (epoch-1)/1000; n++ {
		count := 1000
		if (n+1)*1000 > (epoch - 1) {
			count = (epoch - 1) % 1000
		}
		blocks := make([]*block, count)
		for c := 0; c < count; c++ {
			blocks[c] = &block{
				data: make([][]byte, 0),
			}
			blocks[c].data = append(blocks[c].data, b.blocks[n*1000+c].data...)
		}
		conn.SendDirect(multiblock(uint64(n*1000), blocks))
	}
	conn.SendDirect(newBlockBytes(uint64(epoch)))
	for _, action := range currentBlockCache {
		conn.SendDirect(append([]byte{actionsignal}, action...))
	}
	/*for n := 0; n <= epoch; n++ {
		conn.SendDirect(newBlockBytes(uint64(n)))
		if n == epoch {
			for _, action := range currentBlockCache {
				conn.SendDirect(append([]byte{actionsignal}, action...))
			}
		} else {
			for _, action := range b.blocks[n].data {
				conn.SendDirect(append([]byte{actionsignal}, action...))
			}
		}
	}*/
	conn.Ready()
}

const (
	blocksignal  byte = 0
	actionsignal byte = 1
)

func newBlockBytes(epoch uint64) []byte {
	data := []byte{blocksignal}
	util.PutUint64(epoch, &data)
	return data
}

func (b *blockchain) NewBlock(pool ConnectionPool) {
	newBlock := &block{data: make([][]byte, 0)}
	b.current = newBlock
	b.blocks = append(b.blocks, b.current)
	// pool = nil when reading data from file at initialization
	if pool != nil {
		epoch := uint64(len(b.blocks) - 1)
		data := newBlockBytes(epoch)
		if n, err := b.io.Write(data); n != len(data) || err != nil {
			log.Fatalf("could not write block: %v", err)
		}
		pool.Broadcast(data)
	}
}

func (b *blockchain) NewAction(action []byte, pool ConnectionPool) {
	b.mu.Lock()
	b.current.data = append(b.current.data, action)
	b.mu.Unlock() // pool = nil when reading data from file at initialization
	if pool != nil {
		data := []byte{actionsignal}
		util.PutUint64(uint64(len(action)), &data)
		data = append(data, action...)
		if n, err := b.io.Write(data); n != len(data) || err != nil {
			log.Fatalf("could not write action: %v", err)
		}
		pool.Broadcast(append([]byte{actionsignal}, action...))
	}
}

func (b *blockchain) Close() {
	b.io.Close()
}

func OpenBlockchain() (*blockchain, bool) {
	exists := true
	if stat, err := os.Stat("../../chain.dat"); err != nil || stat.Size() == 0 {
		exists = false
	}
	file, err := os.OpenFile("../../chain.dat", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("could not access chain file: %v\n", err)
	}
	b := &blockchain{
		mu:     sync.Mutex{},
		blocks: make([]*block, 0),
		io:     file,
	}
	newBlock := &block{data: make([][]byte, 0)}
	b.current = newBlock
	b.blocks = append(b.blocks, b.current)
	if !exists {
		data := make([]byte, 9)
		data[0] = blocksignal
		// write the creation of genesis block
		if n, err := file.Write(data); n != len(data) || err != nil {
			log.Fatalf("could not write block: %v", err)
		}
		return b, false
	}
	signal := make([]byte, 9)
	for {
		if n, err := b.io.Read(signal); n != 9 || err != nil {
			break
		}
		number, _ := util.ParseUint64(signal, 1)
		if signal[0] == blocksignal {
			if number == 0 {
				if len(b.blocks) != 1 {
					log.Fatal("blockchain file corrupted: genesis block elsewhere")
				}
			} else if number != uint64(len(b.blocks)) {
				log.Fatal("blockchain file corrupted: block out of order")
			} else {
				b.NewBlock(nil)
			}
		} else if signal[0] == actionsignal {
			data := make([]byte, int(number))
			if n, _ := b.io.Read(data); n != len(data) {
				log.Fatal("blockchain file corrupted: incomplete action")
			}
			b.NewAction(data, nil)
		} else {
			log.Fatal("blockchain file corrupted: invalid data type")
		}
	}
	return b, true
}
