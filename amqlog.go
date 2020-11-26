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

func checktopic(s *stomp.Subscription, f *os.File, name string) {
	for {
		msg := <-s.C
		currentTime := time.Now()
		_, _ = f.WriteString(currentTime.Format("2006.01.02 15:04:05") +
			name + string(msg.Body) + "\n")
		fmt.Printf(currentTime.Format("2006.01.02 15:04:05") +
			"Wrote bytes " + name + "\n")
	}
}

func main() {
	// Connect to ActiveMQ on joe (the computer)
	// Use HeartBeat to tell the connection to not timeout.
	conn, err := stomp.Dial("tcp", "joe.lowell.edu:61613",
		stomp.ConnOpt.HeartBeat(6*time.Hour, 6*time.Hour))
	defer conn.Disconnect()

	if err != nil {
		fmt.Println(err)
	}

	// Create and open a log file on disk.
	f, err := os.Create("/home/pi/go/amqoutput")
	check(err)
	defer f.Close()

	// Subscribe to various topics.
	wrs, err := conn.Subscribe("/topic/wrs.loisTelemetry", stomp.AckAuto)
	defer wrs.Unsubscribe()
	aos, err := conn.Subscribe("/topic/AOS.AOSPubDataSV.AOSDataPacket",
		stomp.AckAuto)
	defer wrs.Unsubscribe()

	// The listeners for the various subscriptions are handled as goroutines.
	fmt.Println("starting wrs")
	go checktopic(wrs, f, "wrs")
	fmt.Println("starting aos")
	go checktopic(aos, f, "aos")

	// Tell the main program to wait forever.
	select {}
}
