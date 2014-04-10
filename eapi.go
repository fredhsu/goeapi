package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "bytes"
    "log"
)

type Parameters struct {
    Version int `json:"version"`
    Cmds []string `json:"cmds"`
    Format string `json:"format"`
}

type Request struct {
    Jsonrpc string `json:"jsonrpc"`
    Method string `json:"method"`
    Params Parameters `json:"params"`
    Id string `json:"id"`
}

type JsonRpcResponse struct {
    Jsonrpc string `json:"jsonrpc"`
    Result []interface{} `json:"result"`
    Id string `json:"id"`
}

type ShowVersion struct {
    ModelName string `json:"modelName"`
    InternalVersion string `json:"internalVersion"`
    SystemMacAddress string `json:"systemMacAddress"`
    SerialNumber string `json:"serialNumber"`
    MemTotal string `json:"memTotal"`
    BootupTimestap string `json:"bootupTimestamp"`
    MemFree string `json:"memFree"`
    Version string `json:"version"`
    Architecture string `json:"architecture"`
    InternalBuildId string `json:"internalBuildId"`
    HardwareRevision string `json:"hardwareRevision"`
}

func eapiCall(url string, cmds []string) *http.Response { 
    p := Parameters{1, cmds, "json"}
    req := Request{"2.0", "runCmds", p, "1"}
    buf, err := json.Marshal(req)
    if err != nil {
        panic(err)
    }
    resp, err := http.Post(url, "application/json", bytes.NewReader(buf))

    if err != nil {
        panic(err)
    }
    return resp
}

func decodeEapiResponse(resp *http.Response) []interface{} {
    dec := json.NewDecoder(resp.Body)

    // Could I use a different type of variable to decode into
    // a known struct type?
    var v map[string]interface{}
    if err := dec.Decode(&v); err != nil {
        log.Println(err)
        return nil
    }
    return v["result"].([]interface{})
}

func showVersion (resp *http.Response) ShowVersion {
    dec := json.NewDecoder(resp.Body)
    // Could I use a different type of variable to decode into
    // a known struct type?
    var sv ShowVersion 
    if err := dec.Decode(&sv); err != nil {
        log.Println(err)
    }
    fmt.Println(sv)
    return sv
}

func main() {
    cmds := []string{"show version"}
    url := "http://admin:admin@192.168.56.101/command-api/"
    resp := eapiCall(url, cmds)
    /*
    m := decodeEapiResponse(resp)
    fields := m[0].(map[string]interface{})
    fmt.Println(fields["version"])
    */
    m := showVersion(resp)
    fmt.Println(m)
}

