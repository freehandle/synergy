package main

import (
	"errors"
	"fmt"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/breeze/network/trusted"
)

type ConnectionPool map[crypto.Token]*CachedConnection

func (p ConnectionPool) Broadcast(data []byte) {
	for _, conn := range p {
		conn.Send(data)
	}
}

func (p ConnectionPool) Connect(c *CachedConnection) {
	p[c.conn.Token] = c
}

func (p ConnectionPool) Drop(token crypto.Token) {
	if conn, ok := p[token]; ok {
		if conn != nil {
			conn.Live = false
			conn.Close()
		}
		delete(p, token)
	}
}

type CachedConnection struct {
	Live    bool
	conn    *trusted.SignedConnection
	ready   bool
	receive chan []byte
	queue   chan struct{}
}

func (c *CachedConnection) Send(data []byte) {
	if len(data) == 0 {
		fmt.Println("empty data")
		return
	}
	if c.Live {
		c.receive <- data
	}
}

func (c *CachedConnection) SendDirect(data []byte) error {
	if (!c.Live) || c.ready {
		return errors.New("connection is dead")
	}
	if err := c.conn.Send(data); err != nil {
		fmt.Println("error sending data:", err)
		c.conn.Shutdown()
		c.Live = false
		c.Close()
		return err
	}
	return nil
}

func (c *CachedConnection) Ready() {
	c.ready = true
	if c.Live {
		c.queue <- struct{}{}
	}
}

func (c *CachedConnection) Close() {
	c.receive <- nil
}

func NewCachedConnection(conn *trusted.SignedConnection) *CachedConnection {

	cached := &CachedConnection{
		Live:    true,
		conn:    conn,
		ready:   false,
		receive: make(chan []byte),
		queue:   make(chan struct{}, 2),
	}

	msgCache := make([][]byte, 0)

	// send loop
	go func() {
		defer func() {
			conn.Shutdown()
			cached.Live = false
			close(cached.receive)
			close(cached.queue)
			fmt.Println("shut down connection")
		}()
		for {
			select {
			case <-cached.queue:
				if N := len(msgCache); N > 0 {
					data := msgCache[0]
					msgCache = msgCache[1:]
					if err := conn.Send(data); err != nil {
						fmt.Println("error sending data:", err)
						return
					}
					if N > 1 {
						// this will never block because there will be one
						// buffer slot on the channel
						cached.queue <- struct{}{}
					}
				}
			case data := <-cached.receive:
				if data == nil {
					fmt.Println("shutting down connection")
					return
				}
				msgCache = append(msgCache, data)
				if cached.ready && len(cached.queue) < cap(cached.queue) {
					cached.queue <- struct{}{}
				}
			}
		}
	}()

	return cached
}
