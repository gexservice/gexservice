package main

import (
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/transport"
)

func main() {
	forward, _ := transport.NewTransportH("tcp://192.168.1.1:80")
	web.Shared.Handle("/", forward)
	web.ListenAndServe(":9322")
}
