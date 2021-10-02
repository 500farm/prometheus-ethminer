package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func discoverTargets(port int) []string {
	result := []string{"127.0.0.1:" + fmt.Sprint(port)}
	for _, ip := range getContainerTargets(port) {
		result = append(result, ip+":"+fmt.Sprint(port))
	}
	return result
}

var cli *client.Client

func createClient() error {
	if cli == nil {
		var err error
		cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return err
		}
	}
	return nil
}

func getContainerTargets(port int) []string {
	ips := []string{}
	if err := createClient(); err != nil {
		return ips
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(filters.Arg("expose", strconv.Itoa(port))),
	})
	if err != nil {
		return ips
	}
	for _, container := range containers {
		for _, net := range container.NetworkSettings.Networks {
			ips = append(ips, net.IPAddress)
		}
	}
	return ips
}
