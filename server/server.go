package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	// Deal with an error event.
	if err != nil {
		panic(err)
	}

}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
	for {
		conn, _ := ln.Accept()
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
	reader := bufio.NewReader(client)
	for {
		msg, _ := reader.ReadString('\n')
		msgs <- Message{clientid, msg}
	}
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above.
	ln, err := net.Listen("tcp", *portPtr)
	handleError(err)
	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	id := 0
	//Start accepting connections
	go acceptConns(ln, conns)
	for {
		select {
		case conn := <-conns:
			//TODO Deal with a new connection
			// - assign a client ID
			// - add the client to the clients channel
			// - start to asynchronously handle messages from this client
			clients[id] = conn
			fmt.Println("client connected", id)
			go handleClient(conn, id, msgs)
			id++
		case msg := <-msgs:
			//TODO Deal with a new message
			fmt.Println(msg.sender, ":", msg.message)
			// // Send the message to all clients that aren't the sender
			for i := 0; i < id; i++ {
				if msg.sender != i {
					fmtedMsg := fmt.Sprintf("%d%s%s", msg.sender, ":", msg.message)
					fmt.Fprint(clients[i], fmtedMsg)
				}
			}
		}
	}
}
