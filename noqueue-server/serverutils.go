package main

import (
	"net"
	"sync"
)

type Topic string

type Client struct {
	conn net.Conn 
	uniqueName string 
}

type SafeClients struct {
	mut sync.Mutex
	clients map[Topic][]*Client
}

func (client *Client) String() string {
	return client.uniqueName
}

func (client *Client) sendMessage(message string) error {
	data := []byte(message)
	_, err := client.conn.Write(data)
	return err
}

func createClient(conn net.Conn) *Client {
	return &Client{conn: conn, uniqueName: conn.RemoteAddr().String()}
}
