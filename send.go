package main

import (
	"time"

	"github.com/Moarqi/gorf/pulsetrain"
)

func main() {
	pulsetrain.Init()
	pulsetrain.Send(49010498, 15, true, true)
	time.Sleep(time.Second * 2)
	pulsetrain.Send(49010498, 15, false, true)
}
