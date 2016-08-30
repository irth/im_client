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
	Name          string
	Subscriptions []proto.Event
	Mutex         sync.Mutex
}

var pluginsLock sync.Mutex
var plugins *list.List

func initPlugin(conn net.Conn) {
	defer conn.Close()
	d := json.NewDecoder(conn)
	e := json.NewEncoder(conn)

	// first message that's supposed to be sent by the plugin is so called init message
	// it contains the key used to authorize the plugin
	var m proto.Message
	d.Decode(&m)

	if m.Type != "InitMessage" {
		// first message is supposed to be an InitMessage
		return
	}

	im := m.InitMessage

	// TODO: make secret configurable
	if im.Secret != "sekret" {
		// sorry, wrong key, plugin not authorized
		e.Encode(proto.InitMessageResponse{false})
		return
	}

	e.Encode(proto.InitMessageResponse{true}) // secret key is correct, plugin is authorized
	plugin := Plugin{
		Name: im.Name,
	}

	// every plugin has its own goroutine so we need to be careful when accessing global variables
	pluginsLock.Lock()

	//add plugin to plugin list
	el := plugins.PushFront(&plugin)
	defer func() {
		// remove the plugin from the list when its goroutine exits
		pluginsLock.Lock()
		plugins.Remove(el)
		pluginsLock.Unlock()
	}()

	// let go of the plugin list after modyfing it
	pluginsLock.Unlock()

	handlePlugin(&plugin, d, e)
}

func handlePlugin(plugin *Plugin, d *json.Decoder, e *json.Encoder) {
	var err error = nil
	var m proto.Message
	for { // listen until we can't decode new messages anymore
		err = d.Decode(&m)
		if err != nil {
			break
		}
		switch m.Type {
		case "SubscribeMessage":
			sm := m.SubscribeMessage
			plugin.Mutex.Lock() // acquire a lock on the plugin so that nobody else can modify it while we do.
			plugin.Subscriptions = append(plugin.Subscriptions, sm.Events...)
			fmt.Println(plugin.Subscriptions)
			plugin.Mutex.Unlock() // let go of the lock
		}

	}
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
		go initPlugin(conn)
	}
}
