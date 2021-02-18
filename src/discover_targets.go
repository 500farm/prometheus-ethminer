package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/prometheus/common/log"
)

func getContainerIds(port int) []string {
	command := "docker ps --filter 'expose=" + fmt.Sprint(port) + "' | egrep -o '^[0-9a-f]{12}'"
	cmd := exec.Command("/bin/sh", "-c", command)
	out, err := cmd.Output()
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) && ee.ExitCode() == 127 {
			// exit code 127 means docker not found, do not log it
		} else {
			log.Errorln(err)
		}
		return []string{}
	}
	return strings.Split(strings.Trim(string(out), "\n"), "\n")
}

func getContainerIPAddress(id string) string {
	command := "docker inspect " + id + " --format '{{.NetworkSettings.IPAddress}}'"
	cmd := exec.Command("/bin/sh", "-c", command)
	out, err := cmd.Output()
	if err != nil {
		log.Errorln(err)
		return ""
	}
	return strings.Trim(string(out), "\n")
}

func discoverTargets(port int) []string {
	result := []string{"127.0.0.1:" + fmt.Sprint(port)}
	for _, id := range getContainerIds(port) {
		ip := getContainerIPAddress(id)
		if ip != "" {
			result = append(result, ip+":"+fmt.Sprint(port))
		}
	}
	log.Infoln("Discovered targets: ", result)
	return result
}
