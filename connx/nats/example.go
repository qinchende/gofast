package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

const serverUrl = "nats://x.x.x.x:4222"

// 发布订阅模型：订阅1
func NatsSub11() {
	// Connect to a server
	nc, _ := nats.Connect(serverUrl)
	defer nc.Close()

	// Simple Async Subscriber
	if _, err := nc.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Sub11: %s\n", m.Data)
	}); err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(120 * time.Second)
	fmt.Printf("Sub11 exit, bye bye...\n")
}

// 发布订阅模型：订阅2
func NatsSub12() {
	// Connect to a server
	nc, _ := nats.Connect(serverUrl)
	defer nc.Close()

	// Simple Async Subscriber
	if _, err := nc.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Sub12: %s\n", m.Data)
	}); err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(120 * time.Second)
	fmt.Printf("Sub12 exit, bye bye...\n")
}

// 发布订阅模型：发布
func NatsPub1() {
	// Connect to a server
	nc, _ := nats.Connect(serverUrl)
	defer nc.Close()

	for i := 0; i < 1; i++ {
		// Simple Publisher
		if err := nc.Publish("foo", []byte(fmt.Sprintf("hello %d", i))); err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("Pub1 finished, bye bye...\n")
}

func RunExample() {
	//go NatsSub11()
	//go NatsSub12()
	go NatsPub1()
	time.Sleep(125 * time.Second)
	fmt.Printf("NatsDemo finished, bye bye...\n")
}
