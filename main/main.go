package main

import (
	"fmt"
	"github.com/fredhsu/go-eapi"
	"github.com/mitchellh/mapstructure"
)

func main() {
	cmds := []string{"show version", "show interfaces"}
	url := "https://admin:admin@dbrl3-leaf1/command-api/"
	jr := eapi.Call(url, cmds, "json")
	var sv eapi.ShowVersion
	err := mapstructure.Decode(jr.Result[0], &sv)
	if err != nil {
		panic(err)
	}
	fmt.Println("\nVersion: ", sv.Version)
	//configCmds := []string{"enable", "configure", "interface ethernet 1", "descr go"}
	configCmds := []string{"enable", "configure", "aaa root secret arista"}
	jr = eapi.Call(url, configCmds, "json")
	fmt.Println("result: ", jr.Result)
	fmt.Println("error: ", jr.Error)
	cmds = []string{"show interfaces"}
	jr = eapi.Call(url, cmds, "json")
	var si eapi.ShowInterfaces
	err = mapstructure.Decode(jr.Result[0], &si)
	if err != nil {
		panic(err)
	}
	// fmt.Println("result: ", si)
	fmt.Println("result: ", si.Interfaces["Ethernet10"].Description)
	fmt.Println("result: ", si.Interfaces["Ethernet10"].InterfaceStatistics)
	fmt.Println("result: ", si.Interfaces["Ethernet10"].Mtu)
	fmt.Println("result: ", si.Interfaces["Ethernet10"].LineProtocolStatus)
	fmt.Printf("result: %+v \n", si.Interfaces["Ethernet10"].InterfaceAddress)
	fmt.Printf("result: %+v \n", si.Interfaces["Ethernet10"].InterfaceCounters.OutErrorsDetail)
	// for k, v := range si.Interfaces["Ethernet10"].InterfaceCounters.OutErrorsDetail {
	// 	fmt.Println(k, ":", v)
	// }

}
