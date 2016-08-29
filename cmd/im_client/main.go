package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/irth/im_client/proto"
	"net"
	"sync"
)

type Plugin struct {
	Name string
}

var pluginsLock sync.Mutex
var plugins *list.List

func handlePlugin(conn net.Conn) {
	defer conn.Close()
	d := json.NewDecoder(conn)
	e := json.NewEncoder(conn)

	// first message that's supposed to be sent by the plugin is so called init message
	// it contains the key used to authorize the plugin
	var m proto.InitMessage
	d.Decode(&m)

	// TODO: make secret configurable
	if m.Secret != "sekret" {
		// sorry, wrong key, plugin not authorized
		e.Encode(proto.InitMessageResponse{false})
		return
	}

	e.Encode(proto.InitMessageResponse{true}) // secret key is correct, plugin is authorized

	// every plugin has its own goroutine so we need to be careful when accessing global variables
	pluginsLock.Lock()

	//add plugin to plugin list
	el := plugins.PushFront(m.Name)
	defer func() {
		// remove the plugin from the list when its goroutine exits
		pluginsLock.Lock()
		plugins.Remove(el)
		pluginsLock.Unlock()
	}()

	// let go of the plugin list after modyfing it
	pluginsLock.Unlock()
}

func main() {
	ln, _ := net.Listen("tcp", ":9123")
	pluginsLock = sync.Mutex{}
	plugins = list.New()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("error ln.Accept()")
		}
		go handlePlugin(conn)
	}
}
