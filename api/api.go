package api

import (
	"encoding/json"
	"fmt"
	"github.com/irth/im_client/proto"
	"net"
)

type Plugin struct {
	Name string
}

type ConnData struct {
	Addr   string
	Secret string
}

type Conn struct {
	conn      net.Conn
	Connected bool
	encoder   *json.Encoder
	decoder   *json.Decoder
}

func NewConn(p Plugin, d ConnData) Conn {
	c := Conn{}
	c.conn, _ = net.Dial("tcp", d.Addr)
	c.encoder = json.NewEncoder(c.conn)
	c.decoder = json.NewDecoder(c.conn)
	c.SendMessage(proto.InitMessage{
		p.Name,
		d.Secret,
	})
	var res proto.AuthMessage
	c.decoder.Decode(&res)
	fmt.Println(res.Success)
	return c
}

func (c Conn) SendMessage(m interface{}) {
	c.encoder.Encode(m)
}
