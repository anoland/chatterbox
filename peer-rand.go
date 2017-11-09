package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/secondbit/wendy"
)

type applicationCallback struct{}

func (app *applicationCallback) OnError(err error) {
	log.Printf("\nERROR %v %v\n\n", app, err)
}
func (app *applicationCallback) OnDeliver(msg wendy.Message) {
	log.Printf("Delivered message: %v\n", msg.String())
}
func (app *applicationCallback) OnForward(msg *wendy.Message, nextId wendy.NodeID) bool {
	//log.Printf("Forward %v %v\n", msg.String(), nextId)
	return true
}
func (app *applicationCallback) OnNewLeaves(leafset []*wendy.Node) {
	log.Printf("New %v\n", leafset)
}
func (app *applicationCallback) OnNodeJoin(node wendy.Node) {
	log.Printf("Join %v\n", node)
}
func (app *applicationCallback) OnNodeExit(node wendy.Node) {
	log.Printf("Exit %v\n", node)
}
func (app *applicationCallback) OnHeartbeat(node wendy.Node) {
	log.Printf("Beat %v\n", node)
}

func main() {
	// Most of this is taken directly from the README (with typos corrected)

	// The port for _this_ node. I do this instead of using a random value so as to know the _cluster_ port.
	this_port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
	//hostname, err := os.Hostname()
	if err != nil {
		panic(err.Error())
	}
	//id_string := fmt.Sprintf("%v %v", this_port, hostname)
	id, err := wendy.NodeIDFromBytes([]byte(uuid.NewV1().String()))
	if err != nil {
		panic(err.Error())
	}
	node := wendy.NewNode(id, "127.0.0.1", "127.0.0.1", "none", this_port)
	//log.Printf("%#v", node)

	credentials := wendy.Passphrase("I <3 Gophers.")
	cluster := wendy.NewCluster(node, credentials)
	cluster.RegisterCallback(&applicationCallback{})
	go func() {
		defer cluster.Stop()
		log.Printf("Listening %v\n", this_port)
		err := cluster.Listen()
		if err != nil {
			panic(err.Error())
		}
	}()

	// If there are two parameters, join the cluster.
	if len(os.Args) > 2 {
		cluster_port, err := strconv.Atoi(os.Args[2])
		if err == nil {
			//log.Printf("cluster port %v\n", cluster_port)
			err := cluster.Join("127.0.0.1", cluster_port)
			if err != nil {
				log.Println(err.Error())
			}
		} else {
			log.Println(err)
			os.Exit(2)
		}

		// Otherwise, send a random message to the cluster every 5 seconds.
	} else {
		ticker := time.NewTicker(5 * time.Second)
		quit := make(chan struct{})
		go func() {
			bn := 0
			for {
				select {
				case <-ticker.C:
					h := sha256.New()
					io.WriteString(h, fmt.Sprintf("%d", bn%10))
					bn += 1
					id, err := wendy.NodeIDFromBytes(h.Sum(nil))
					if err != nil {
						panic(err.Error())
					}

					purpose := byte(16)
					m := fmt.Sprintf("[%v.%v]", rand.Int63(), int32(time.Now().Unix()))
					//log.Println(m)
					msg := cluster.NewMessage(purpose, id, []byte(m))
					err = cluster.Send(msg)
					if err != nil {
						log.Println(err.Error())
					}
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()
	}

	select {}
}
