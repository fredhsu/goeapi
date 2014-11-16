package main

import (
	"fmt"
	"github.com/fredhsu/go-eapi"
	"github.com/mitchellh/mapstructure"
)

// TODO: Create something to track a topology or bgp neighbor list
// TODO: Create database of switches to be able to find info about a neighbor
// TODO: Find a way to 'walk' down switches, maybe recursively
// TODO: path tracer type functionality - may need to go down multiple paths using goroutines

func main() {
	// TODO: Extend this to use go routines to hit multiple switches at once
	cmds := []string{"show version", "show interfaces"}
	url := "https://admin:admin@dbrl3-leaf1/command-api/"
	jr := eapi.Call(url, cmds, "json")
	var sv eapi.ShowVersion
	err := mapstructure.Decode(jr.Result[0], &sv)
	if err != nil {
		panic(err)
	}
	fmt.Println("\nVersion: ", sv.Version)
	var si eapi.ShowInterfaces
	err = mapstructure.Decode(jr.Result[1], &si)
	if err != nil {
		panic(err)
	}
	fmt.Println("result: ", si.Interfaces["Ethernet10"].Description)
	fmt.Println("result: ", si.Interfaces["Ethernet10"].InterfaceStatistics)
	fmt.Println("result: ", si.Interfaces["Ethernet10"].Mtu)
	fmt.Println("result: ", si.Interfaces["Ethernet10"].LineProtocolStatus)
	fmt.Printf("result: %+v \n", si.Interfaces["Ethernet10"].InterfaceAddress)
	fmt.Printf("result: %+v \n", si.Interfaces["Ethernet10"].InterfaceCounters.OutErrorsDetail)
	//configCmds := []string{"enable", "configure", "interface ethernet 1", "descr go"}
	//configCmds := []string{"enable", "configure", "aaa root secret arista"}
	//jr = eapi.Call(url, configCmds, "json")
	//fmt.Println("result: ", jr.Result)

	cmds = []string{"show ip route", "show ip bgp neighbors"}
	jr = eapi.Call(url, cmds, "text")
	//fmt.Println(jr.Result[0]["output"])
	//fmt.Println(jr.Result[1]["output"])
	out := fmt.Sprintf("%v", jr.Result[0]["output"])
	routes := eapi.ParseShowIpRoute(out)
	for _, route := range routes {
		if len(route.NextHops) > 1 {
			fmt.Printf(" %+v\n", route)
		}
	}
    cvxcmds := []string{"show network physical-network"}
    jr = eapi.Call(url, cvxcmds, "json")
    fmt.Println(jr.Result[0])

}
