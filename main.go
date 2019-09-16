package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
)

const defaultListenPort int = 8080
const outputFormat string = "%-20s\t%-10s\t%-10s\t%-10s\n"
const notAvailable string = "N/A"
const refreshDelayMs time.Duration = 5000

var containers []types.Container

func main() {
	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err = cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	var listenPort int
	if envListenPort := os.Getenv("LISTEN_PORT"); envListenPort != "" {
		listenPort, _ = strconv.Atoi(envListenPort)
	} else {
		listenPort = defaultListenPort
	}

	// Listen in the background
	go func(socket string) {
		fmt.Printf("Listening on %s\n", socket)
		http.HandleFunc("/", httpHandler)
		http.ListenAndServe(socket, nil)
	}(":" + strconv.Itoa(listenPort))

	for {
		time.Sleep(refreshDelayMs * time.Millisecond)
		containers, err = cli.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	renderHTML(w)
}

func renderHTML(w http.ResponseWriter) {
	fmt.Fprintf(w, outputFormat, "NAME", "STATE", "NETWORK", "PORT")
	for _, c := range containers {
		//fmt.Printf("%+v\n", c)
		var name string
		var ports []string

		if len(c.Names) > 0 {
			name = c.Names[0][1:]
		} else {
			name = notAvailable
		}

		/*
			if len(c.Ports) > 0 {
				port = strconv.Itoa(int(c.Ports[0].PublicPort))
			} else {
				port = notAvailable
			}
		*/

		for _, p := range c.Ports {
			if p.PublicPort == 0 {
				continue
			} else {
				ports = append(ports, p.Type+"/"+strconv.Itoa(int(p.PublicPort)))
			}
		}

		mode := c.HostConfig.NetworkMode
		if mode == "host" {
			mode = "<host>"
		}

		fmt.Fprintf(w, outputFormat, name, c.State, mode, strings.Join(ports, ", "))
	}
}
