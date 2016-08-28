package main

import (
	"encoding/json"
	"fmt"
	"github.com/irth/im_client/api"
	"github.com/irth/im_client/proto"
	"net"
)

var plugins []api.Plugin

func handlePlugin(conn net.Conn) {
	d := json.NewDecoder(conn)
	e := json.NewEncoder(conn)
	var m proto.InitMessage

	d.Decode(&m)
	if m.Secret == "sekret" {
		e.Encode(proto.AuthMessage{
			true,
		})
	} else {
		e.Encode(proto.AuthMessage{
			false,
		})
	}
	fmt.Println(m.Name)
}

func main() {
	ln, _ := net.Listen("tcp", ":9123")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("error ln.Accept()")
		}
		go handlePlugin(conn)
	}
}
