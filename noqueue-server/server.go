package main

import (
	"bufio"
	"fmt"
	"net"
)

const (
	PORT int = 8090
	PUBLISHWORKERS int = 3
	PUBLISHWORKERBUFFER int = 40
)

var pubchan chan pubData
var clientsmap SafeClients

func setupWorkers() {
	pubchan = make(chan pubData, PUBLISHWORKERBUFFER)

	for range PUBLISHWORKERS {
		go handlePublish()
	}
	fmt.Printf("Spawned %d workers for handling PUB commands\n", PUBLISHWORKERS)

}


func startServer() {
	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	setupWorkers()

	clientsmap = SafeClients{
		clients: make(map[Topic][]*Client),
	}

	fmt.Printf("server listening at localhost:%d\n", PORT)
	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept error: %v\n", err)
			continue
		}

		go handleConnection(conn)

	}

}

func handleConnection(conn net.Conn) {
	fmt.Println("handling a new client connection...")
	reader := bufio.NewReader(conn)

	for {

		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Read error: %v\n", err)
			conn.Close()
			break 
		}

		action, data, err := parseMessage(message)
		if err != nil {
			fmt.Printf("Action error: %v\n", err)
			conn.Close()
			break
		}
		shouldExit := handleAction(action, data, conn)
		if shouldExit {
			break
		}
	}
}

func handleAction(action ActionType, data ActionData, conn net.Conn) bool  {
	switch action {
	case PUB:
		data, err := parsePubData(data)
		if err != nil {
			fmt.Printf("PUB data parse error: %v\n", err)
			break // keep alive (for now) or close connection 
		}
		pubchan <- data
	case SUB:
		data, err := parseSubData(data)
		if err != nil {
			fmt.Printf("SUB data parse error: %v\n", err)
			break // keep alive (for now) or close connection 
		}
		handleSubscribe(data, conn)
	case EXIT:
		handleExit(conn)
		return true
	}
	return false
}

func handlePublish() {
	for data := range pubchan {
		count := broadcastMessage(data.topic, data.message)
		fmt.Printf("broadcasted to %d subscribers of %s\n", count, data.topic)
	}
}

func broadcastMessage(topic Topic, message string) int {
	clientsmap.mut.Lock()
	clients := append([]*Client(nil), clientsmap.clients[topic]...)
	defer clientsmap.mut.Unlock()

	count := 0

	for _, client := range clients {
		err := client.sendMessage(message + "\n")
		if err != nil {
			fmt.Printf("Broadcast sending message error: %v\n", err)
			continue
		}
		count++
	}

	return count
}

func handleSubscribe(data subData, conn net.Conn) {
	clientsmap.mut.Lock()
	defer clientsmap.mut.Unlock()
	clients := clientsmap.clients[data.topic];
	for _, client := range clients {
		if client.conn == conn {
			return
		}
	}
	clientsmap.clients[data.topic] = append(clients, createClient(conn))
}

// TODO: make it efficient by having some inverse mapping from client to topics
// so can delete it way faster.
// right now we go through each room and check each client, not efficient
func handleExit(conn net.Conn) {
	clientsmap.mut.Lock()
	defer clientsmap.mut.Unlock()
	defer conn.Close()

	for topic, clients := range clientsmap.clients {
		filtered := clients[:0]
		for _, c := range clients {
			if c.conn != conn {
				filtered = append(filtered, c)
			}
		}
		clientsmap.clients[topic] = filtered
	}
	fmt.Println("A client dropped by EXIT command...")
}
