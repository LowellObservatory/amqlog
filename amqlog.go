package main

import stomp "github.com/go-stomp/stomp"
import (
	"fmt"
	"os"
	"time"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

//Connect to ActiveMQ and listen for messages
func main() {
	conn, err := stomp.Dial("tcp", "joe.lowell.edu:61613",
		stomp.ConnOpt.HeartBeat(6 * time.Hour, 6 * time.Hour))
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create("/home/pi/go/amqoutput")
	check(err)
	defer f.Close()

	subRC1, err := conn.Subscribe("/topic/LOUI.RC1.loisCommandResult", stomp.AckAuto)
	if err != nil {
		fmt.Println(err)
	}
	for {
		msg := <-subRC1.C
        _, err = f.WriteString("RC1: " + string(msg.Body) + "\n\n")
		fmt.Printf("Wrote bytes\n")
		fmt.Println(string(msg.Body))
	}
	subRC2, err := conn.Subscribe("/topic/LOUI.RC2.loisCommandResult", stomp.AckAuto)
	if err != nil {
		fmt.Println(err)
	}
	for {
		msg := <-subRC2.C
        _, err = f.WriteString("RC2: " + string(msg.Body) + "\n\n")
		fmt.Printf("Wrote bytes\n")
		fmt.Println(string(msg.Body))
	}
	subGDR, err := conn.Subscribe("/topic/LOUI.GDR.loisCommandResult", stomp.AckAuto)
	if err != nil {
		fmt.Println(err)
	}
	for {
		msg := <-subGDR.C
        _, err = f.WriteString("GDR: " + string(msg.Body) + "\n\n")
		fmt.Printf("Wrote bytes\n")
		fmt.Println(string(msg.Body))
	}
	subWFS, err := conn.Subscribe("/topic/LOUI.WFS.loisCommandResult", stomp.AckAuto)
	if err != nil {
		fmt.Println(err)
	}
	for {
		msg := <-subWFS.C
        _, err = f.WriteString("WFS: " + string(msg.Body) + "\n\n")
		fmt.Printf("Wrote bytes\n")
		fmt.Println(string(msg.Body))
	}

	err = subRC1.Unsubscribe()
	if err != nil {
		fmt.Println(err)
	}
	err = subRC2.Unsubscribe()
	if err != nil {
		fmt.Println(err)
	}
	err = subGDR.Unsubscribe()
	if err != nil {
		fmt.Println(err)
	}
	err = subWFS.Unsubscribe()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Disconnect()
}
