package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
)

const outputFormat string = "%-20s\t%-10s\t%-10s\t%-10s\n"
const notAvailable string = "N/A"

func main() {
	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf(outputFormat, "NAME", "STATE", "NETWORK", "PORT")
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
				ports = append(ports, p.Type + "/" + strconv.Itoa(int(p.PublicPort)))
			}
		}

		mode := c.HostConfig.NetworkMode
		if mode == "host" {
			mode = "<host>"
		}

		fmt.Printf(outputFormat, name, c.State, mode, strings.Join(ports, ", "))
	}
}
