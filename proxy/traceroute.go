package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
)

// Wrapper of traceroute, IPv4
func tracerouteIPv4Wrapper(httpW http.ResponseWriter, httpR *http.Request) {
	tracerouteRealHandler(false, httpW, httpR)
}

// Wrapper of traceroute, IPv6
func tracerouteIPv6Wrapper(httpW http.ResponseWriter, httpR *http.Request) {
	tracerouteRealHandler(true, httpW, httpR)
}

func tracerouteTryExecute(cmd []string, args [][]string) ([]byte, error) {
	var output []byte
	var err error
	for i := range cmd {
		instance := exec.Command(cmd[i], args[i]...)
		output, err = instance.Output()
		if err == nil {
			return output, err
		}
	}
	return output, err
}

// Real handler of traceroute requests
func tracerouteRealHandler(useIPv6 bool, httpW http.ResponseWriter, httpR *http.Request) {
	query := string(httpR.URL.Query().Get("q"))
	if query == "" {
		invalidHandler(httpW, httpR)
	} else {
		var result []byte
		var err error
		if runtime.GOOS == "freebsd" || runtime.GOOS == "netbsd" {
			if useIPv6 {
				result, err = tracerouteTryExecute(
					[]string{
						"traceroute6",
						"traceroute",
					},
					[][]string{
						{"-a", "-q1", "-w1", query},
						{"-a", "-q1", "-w1", query},
					},
				)
			} else {
				result, err = tracerouteTryExecute(
					[]string{
						"traceroute",
						"traceroute6",
					},
					[][]string{
						{"-a", "-q1", "-w1", query},
						{"-a", "-q1", "-w1", query},
					},
				)
			}
		} else if runtime.GOOS == "openbsd" {
			if useIPv6 {
				result, err = tracerouteTryExecute(
					[]string{
						"traceroute6",
						"traceroute",
					},
					[][]string{
						{"-A", "-q1", "-w1", query},
						{"-A", "-q1", "-w1", query},
					},
				)
			} else {
				result, err = tracerouteTryExecute(
					[]string{
						"traceroute",
						"traceroute6",
					},
					[][]string{
						{"-A", "-q1", "-w1", query},
						{"-A", "-q1", "-w1", query},
					},
				)
			}
		} else if runtime.GOOS == "linux" {
			if useIPv6 {
				result, err = tracerouteTryExecute(
					[]string{
						"traceroute",
						"traceroute",
						"traceroute",
						"traceroute",
					},
					[][]string{
						{"-6", "-A", "-q1", "-N32", "-w1", query},
						{"-4", "-A", "-q1", "-N32", "-w1", query},
						{"-6", "-q1", "-w1", query},
						{"-4", "-q1", "-w1", query},
					},
				)
			} else {
				result, err = tracerouteTryExecute(
					[]string{
						"traceroute",
						"traceroute",
						"traceroute",
						"traceroute",
					},
					[][]string{
						{"-4", "-A", "-q1", "-N32", "-w1", query},
						{"-6", "-A", "-q1", "-N32", "-w1", query},
						{"-4", "-q1", "-w1", query},
						{"-6", "-q1", "-w1", query},
					},
				)
			}
		} else {
			httpW.WriteHeader(http.StatusInternalServerError)
			httpW.Write([]byte("traceroute not supported on this node.\n"))
			return
		}
		if err != nil {
			httpW.WriteHeader(http.StatusInternalServerError)
			httpW.Write([]byte(fmt.Sprintln("traceroute returned error:", err.Error(), ", please check if IPv4/IPv6 is selected correctly.")))
			return
		}
		httpW.Write(result)
	}
}
