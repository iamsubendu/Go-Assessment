package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"net"
	"strings"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	x = make(map[string]*list.List)
)

var (
	openConnection = make(map[net.Conn]bool)
	newConnection  = make(chan net.Conn)
	deadConnection = make(chan net.Conn)
)

func main() {
	ln, err := net.Listen("tcp", ":9000")
	logFatal(err)

	defer ln.Close()

	go func() {

		for {
			conn, err := ln.Accept()
			logFatal(err)

			openConnection[conn] = true
			newConnection <- conn

		}
	}()

	for {
		select {
		case conn := <-newConnection:
			go Store(conn)
		case conn := <-deadConnection:
			for item := range openConnection {
				if item == conn {
					break
				}
			}
		}
	}

}

func Store(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		usernameAndMessage, err := reader.ReadString('\n')
		logFatal(err)

		//fmt.Print(usernameAndMessage)

		parts := strings.Split(usernameAndMessage, ":-")
		username := parts[0]
		message := parts[1]

		fmt.Println(message)
		switch msg := strings.Trim(message, "\r\n"); msg {
		case "#getFirstMsg":
			fmt.Println("first")
			first(conn, username)
		case "#getLastMsg":
			fmt.Println("last")
			last(conn, username)
		case "#dequeue":
			fmt.Println("dequeue")
			Dequeue(conn, username)
		default:

			if _, ok := x[username]; ok {

				x[username].PushBack(message)
				fmt.Println("Messages Fom ", username, ":-")
				for e := x[username].Front(); e != nil; e = e.Next() {
					fmt.Print(e.Value)
				}
				fmt.Println("---------------------")

			} else {
				x[username] = list.New()
				x[username].PushBack(message)
				fmt.Println("Messages Fom ", username, ":-")
				for e := x[username].Front(); e != nil; e = e.Next() {
					fmt.Print(e.Value)
				}
				fmt.Println("---------------------")

			}

		}

	}

	deadConnection <- conn

}

func first(connection net.Conn, username string) {
	first := x[username].Front().Value
	first_msg := fmt.Sprintf("Fisrt Messaage:- %v", first)
	// first_msg := []byte(fmt.Sprintf("%v", first.(interface{})))
	connection.Write([]byte(first_msg))
	fmt.Println("fisrt Functin Called", first)
}

func last(connection net.Conn, username string) {
	first := x[username].Back().Value
	first_msg := fmt.Sprintf("Last Messaage:- %v", first)
	// first_msg := []byte(fmt.Sprintf("%v", first.(interface{})))
	connection.Write([]byte(first_msg))
	fmt.Println("Last Functin Called", first)
}

func Dequeue(connection net.Conn, username string) {
	if x[username].Len() > 0 {
		ele := x[username].Front()
		x[username].Remove(ele)
		msg := "Removed First Message , Because this is FIFO Queue ."
		msg1 := fmt.Sprintf("%v ....\n", msg)
		connection.Write([]byte(msg1))

	} else {
		fmt.Errorf("Pop Error: Queue is empty")
		msg := "Pop Error: Queue is empty"
		connection.Write([]byte(msg))
	}