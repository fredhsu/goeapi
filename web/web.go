package main

import (
	"fmt"
	//"io/ioutil"
	"encoding/json"
	"github.com/fredhsu/go-eapi"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

func toJson(i eapi.ShowInterfaces) (s string) {
	b, err := json.Marshal([]eapi.ShowInterfaces{i})
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func showInterfacesHandler(rw http.ResponseWriter, req *http.Request) {
	cmds := []string{"show interfaces"}
	url := "https://admin:admin@dbrl3-leaf1/command-api/"
	jr := eapi.Call(url, cmds, "json")
	var si eapi.ShowInterfaces
	err := mapstructure.Decode(jr.Result[0], &si)
	if err != nil {
		panic(err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(rw, toJson(si))
	return
}

func showIpRouteHandler(rw http.ResponseWriter, req *http.Request) {
	cmds := []string{"show ip route"}
	url := "https://admin:admin@dbrl3-leaf1/command-api/"
	jr := eapi.Call(url, cmds, "text")
	out := fmt.Sprintf("%v", jr.Result[0]["output"])
	sir := eapi.ParseShowIpRoute(out)

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	jsonOut, err := json.Marshal(sir)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(rw, string(jsonOut))
	//TODO: JSON Marshal of IP.IPNet gives weird data
	// Actually it seems to IP.IPMask
	// Can I do a toString of IP.IPMask to fix it?
	return
}

func main() {
	// Create base path "/v1"
	http.HandleFunc("/show/interfaces", showInterfacesHandler)
	http.HandleFunc("/show/ip/route", showIpRouteHandler)
	// Create something that uses goroutines to go fetch multiple switch info
	http.ListenAndServe(":8081", nil)
	// Take in a post with the IP address of the switch
}
