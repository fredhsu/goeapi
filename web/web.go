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
	//b, err = json.Marshal([]string{"hi", "there"})
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func myHandler(rw http.ResponseWriter, req *http.Request) {
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

func main() {
	// Create base path "/v1"
	http.HandleFunc("/show/interfaces", myHandler)
	http.ListenAndServe(":8081", nil)
	// Take in a post with the IP address of the switch
}
