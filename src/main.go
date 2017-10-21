package  main

import (
	"./probe"
	"./webserver"
	"time"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("=============================================")
	fmt.Println("=============================================")
	fmt.Println("======= Docker Swarm Service Lister =========")
	fmt.Println("=============================================")
	stacks := probe.DiscoverStacks()
	webserverData := &webserver.WebserverData{Stacks: stacks}

	c := make(chan bool)
	go webserver.StartServer("7777", webserverData, c)
	fmt.Println("> Started the web server, now polling swarm")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	for i := 1; ; i++ { // this is still infinite
		t := time.NewTicker(time.Second * 30)
		select {
		case <-stop:
			fmt.Println("> Shutting down polling")
			break
		case <-t.C:
			fmt.Println("> Updating Stacks")
			webserverData.UpdateStacks(probe.DiscoverStacks())
			continue
		}
		break // only reached if the quitCh case happens
	}
	fmt.Println("> Shutting down webserver")
	c <- true
	if b := <-c; b {
		fmt.Println("> Webserver shut down")
	}
	fmt.Println("> Shut down app")
}
