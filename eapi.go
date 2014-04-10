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
    MemTotal float64 `json:"memTotal"`
    BootupTimestap float64 `json:"bootupTimestamp"`
    MemFree float64 `json:"memFree"`
    Version string `json:"version"`
    Architecture string `json:"architecture"`
    InternalBuildId string `json:"internalBuildId"`
    HardwareRevision string `json:"hardwareRevision"`
}

func eapiCall(url string, cmds []string) JsonRpcResponse {
    p := Parameters{1, cmds, "json"}
    req := Request{"2.0", "runCmds", p, "1"}
    buf, err := json.Marshal(req)
    if err != nil {
        panic(err)
    }
    resp, err := http.Post(url, "application/json", bytes.NewReader(buf))
    defer resp.Body.Close()

    if err != nil {
        panic(err)
    }
    return decodeEapiResponse(resp)
}

func decodeEapiResponse(resp *http.Response) JsonRpcResponse {
    dec := json.NewDecoder(resp.Body)
    var v JsonRpcResponse
    if err := dec.Decode(&v); err != nil {
        log.Println(err)
    }
    return v
}

func showVersion (jr JsonRpcResponse) ShowVersion {
    tmp := jr.Result[0].(map[string]interface{})
    var sv ShowVersion
    sv = ShowVersion {
        ModelName: tmp["modelName"].(string),
        InternalVersion: tmp["internalVersion"].(string),
        SystemMacAddress: tmp["systemMacAddress"].(string),
        SerialNumber: tmp["serialNumber"].(string),
        MemTotal: tmp["memTotal"].(float64),
        BootupTimestap: tmp["bootupTimestamp"].(float64),
        MemFree: tmp["memFree"].(float64),
        Version: tmp["version"].(string),
        Architecture: tmp["architecture"].(string),
        InternalBuildId: tmp["internalBuildId"].(string),
        HardwareRevision: tmp["hardwareRevision"].(string),
    }
    return sv
}

func main() {
    cmds := []string{"show version"}
    url := "http://admin:admin@192.168.56.101/command-api/"
    jr := eapiCall(url, cmds)
    fmt.Println("result: ", jr.Result)
    sv := showVersion(jr)
    fmt.Println("\nVersion: ", sv.Version)
}

