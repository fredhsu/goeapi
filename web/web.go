package main

import (
    "fmt"
    "io/ioutil"
)

type Page struct {
    Title string
    Body []byte
}

func main() {
    http.HandleFunc("/show/interfaces", handler)
    http.ListenAndServe(":8081", nil)
}
