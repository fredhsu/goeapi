package main

import (
	"fmt"
	"log"

	"github.com/fredhsu/goeapi"
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
	jr := goeapi.Call(url, cmds, "json")
	var sv goeapi.ShowVersion
	err := mapstructure.Decode(jr.Result[0], &sv)
	if err != nil {
		panic(err)
	}
	fmt.Println("\nVersion: ", sv.Version)
	var si goeapi.ShowInterfaces
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
	//jr = goeapi.Call(url, configCmds, "json")
	//fmt.Println("result: ", jr.Result)

	cmds = []string{"show ip route", "show ip bgp neighbors"}
	jr = goeapi.Call(url, cmds, "text")
	//fmt.Println(jr.Result[0]["output"])
	//fmt.Println(jr.Result[1]["output"])
	out := fmt.Sprintf("%v", jr.Result[0]["output"])
	routes := goeapi.ParseShowIpRoute(out)
	for _, route := range routes {
		if len(route.NextHops) > 1 {
			fmt.Printf(" %+v\n", route)
		}
	}

	cvx := goeapi.Call(url, []string{"show network physical-topology hosts"}, "json")
	fmt.Println(cvx.Result[0]["hosts"])
	showInt := goeapi.CallShowInterfaces(url, "")
	fmt.Println(showInt)
	e := goeapi.EosNode{"https",
		"fredlhsu",
		"arista",
		"172.28.170.27"}
	l, err := e.ShowLldpNeighbors()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(l)
}
