package main

import "github.com/irth/im_client/api"

func main() {
	p := api.Plugin{"testplugin"}
	d := api.ConnData{
		"localhost:9123",
		"sekret",
	}
	c := api.NewConn(p, d)
	c.Subscribe("foo")
}
