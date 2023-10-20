package api

import (
	"log"
	"os"
	"sync"

	"github.com/freehandle/breeze/crypto"
)

type filePasswordManager struct {
	mu        sync.Mutex
	file      os.File
	hashes    []crypto.Hash
	passwords map[crypto.Token]int
}

func (f *filePasswordManager) Check(user crypto.Token, password crypto.Hash) bool {
	if n, ok := f.passwords[user]; ok {
		if n >= len(f.hashes) {
			log.Printf("unexpected error in file password manager")
			return false
		}
		return password.Equal(f.hashes[n])
	}
	return false
}

func (f *filePasswordManager) Has(user crypto.Token) bool {
	_, ok := f.passwords[user]
	return ok
}

func (f *filePasswordManager) Set(user crypto.Token, password crypto.Hash, email string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	data := append(user[:], password[:]...)
	n, ok := f.passwords[user]
	if ok {
		if n > len(f.hashes) {
			log.Printf("unexpected error in file password manager")
			return false
		}
		if n, err := f.file.WriteAt(data, int64(n)*2*crypto.Size); n != 64 || err != nil {
			log.Printf("unexpected error in file password manager: %v", err)
			return false
		}
		f.hashes[n] = password
	}
	f.file.Seek(0, 2)
	if n, err := f.file.Write(data); n != 64 || err != nil {
		log.Printf("unexpected error in file password manager: %v", err)
		return false
	}
	f.hashes = append(f.hashes, password)
	f.passwords[user] = len(f.hashes) - 1
	return true
}

func NewFilePasswordManager(filename string) PasswordManager {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("could not open password manager file: %v", err)
	}
	stat, err := file.Stat()
	if err != nil {
		log.Fatalf("could not stat password manager file: %v", err)
	}
	size := stat.Size()
	if size%64 != 0 {
		log.Fatal("corrupted password manager file: incompatible size")
	}
	manager := filePasswordManager{
		file:      *file,
		hashes:    make([]crypto.Hash, size/64),
		passwords: make(map[crypto.Token]int),
	}
	entry := make([]byte, 64)
	for n := 0; n < int(size)/64; n++ {
		if n, err := file.Read(entry); n != 64 {
			log.Fatalf("corrupted password manager file: %v", err)
		}
		var token crypto.Token
		copy(token[:], entry[:32])
		copy(manager.hashes[n][:], entry[32:])
		manager.passwords[token] = n
	}
	return &manager
}

type PasswordManager interface {
	Check(user crypto.Token, password crypto.Hash) bool
	Set(user crypto.Token, password crypto.Hash, email string) bool
	Has(user crypto.Token) bool
}
