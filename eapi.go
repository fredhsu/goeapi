package goeapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Parameters struct {
	Version int      `json:"version"`
	Cmds    []string `json:"cmds"`
	Format  string   `json:"format"`
}

type Request struct {
	Jsonrpc string     `json:"jsonrpc"`
	Method  string     `json:"method"`
	Params  Parameters `json:"params"`
	Id      string     `json:"id"`
}

type RawJsonRpcResponse struct {
	Jsonrpc string                 `json:"jsonrpc"`
	Result  []json.RawMessage      `json:"result"`
	Error   map[string]interface{} `json:"error"`
	Id      string                 `json:"id"`
}

type JsonRpcResponse struct {
	Jsonrpc string                   `json:"jsonrpc"`
	Result  []map[string]interface{} `json:"result"`
	Error   map[string]interface{}   `json:"error"`
	Id      string                   `json:"id"`
}

type ShowInterfacesResponse struct {
	Jsonrpc string
	Result  []ShowInterfaces
	Error   map[string]interface{}
	Id      string
}

type ShowVersion struct {
	ModelName        string  `json:"modelName"`
	InternalVersion  string  `json:"internalVersion"`
	SystemMacAddress string  `json:"systemMacAddress"`
	SerialNumber     string  `json:"serialNumber"`
	MemTotal         float64 `json:"memTotal"`
	BootupTimestap   float64 `json:"bootupTimestamp"`
	MemFree          float64 `json:"memFree"`
	Version          string  `json:"version"`
	Architecture     string  `json:"architecture"`
	InternalBuildId  string  `json:"internalBuildId"`
	HardwareRevision string  `json:"hardwareRevision"`
}

type ShowLldpNeighbors struct {
	TablesDeletes        int            `json:"tablesDeletes"`
	TablesAgeOuts        int            `json:"tablesAgeOuts"`
	TablesDrops          int            `json:"tablesDrops"`
	TablesInserts        int            `json:"tablesInserts"`
	TablesLastChangeTime int            `json:"tablesLastChangeTime"`
	LldpNeighbors        []LldpNeighbor `json:"lldpNeighbors"`
}

type LldpNeighbor struct {
	NeighborDevice string
	NeighborPort   string
	Port           string
	Ttl            int
}

type ShowDirectFlowFlows struct {
	Flows []Flow
}

type Flow struct {
	Priority   int
	MatchBytes int
	Name       string
}

//TODO: May need to replace Interface with a generic interface
type ShowInterfaces struct {
	Interfaces map[string]Interface `json:"interfaces"`
}

type ShowVlanResponse struct {
	Jsonrpc string
	Result  []ShowVlan
	Error   map[string]interface{}
	Id      string
}

type ShowVlan struct {
	SourceDetail string              `json:"sourceDetail"`
	Vlans        map[string]VlanInfo `json:"vlans"`
}

type VlanInfo struct {
	Status     string                    `json:"status"`
	Name       string                    `json:"name"`
	Interfaces map[string]InterfaceNames `json:"interfaces"`
	Dynamic    bool                      `json:"dynamic"`
}

type InterfaceNames struct {
	NameVlan map[string]VlanProperty
}

type VlanProperty struct {
	Property map[string]bool `json:"privatePromoted"`
}

// todo: account for different interface types, this is specific to Eth
type Interface struct {
	Bandwidth                 int
	BurnedInAddress           string
	Description               string //`json:"description"`
	ForwardingModel           string
	Hardware                  string
	InterfaceAddress          []InterfaceAddress
	InterfaceCounters         EthInterfaceCounters
	InterfaceMembership       string
	InterfaceStatistics       InterfaceStatistics
	InterfaceStatus           string
	L2Mtu                     int
	LastStatusChangeTimestamp float64
	LineProtocolStatus        string
	Mtu                       int
	Name                      string
	PhysicalAddress           string
}

type InterfaceAddress struct {
	BroadcastAddress       string
	PrimaryIp              IpAddress
	SecondaryIps           interface{}
	SecondaryIpOrderedList []IpAddress
	VirtualIp              IpAddress
}

type IpAddress struct {
	Address string
	MaskLen int
}

type EthInterfaceCounters struct {
	CounterRefreshTime float64
	InBroadcastPkts    int
	InDiscards         int
	InMulticastPkts    int
	InOctets           int
	InUcastPkts        int
	InputErrorsDetail  PhysicalInputErrors
	LastClear          float64
	LinkStatusChanges  int
	OutBroadcastPkts   int
	OutDiscards        int
	OutMulticastPkts   int
	OutOctets          int
	OutUcastPkts       int
	OutErrorsDetail    PhysicalOutputErrors
	TotalInErrors      int
	TotalOutErrors     int
}

type PhysicalInputErrors struct {
	AlignmentErrots int
	FcsErrors       int
	GiantFrames     int
	RuntFrames      int
	RxPause         int
	SymbolErrors    int
}

type PhysicalOutputErrors struct {
	Collisions            int
	DeferredTransmissions int
	LateCollisions        int
	TxPause               int
}

type InterfaceStatistics struct {
	InBitsRate     float64
	OutBitsRate    float64
	InPktsRate     float64
	UpdateInterval float64
	OutPktsRate    float64
}

type OpenstackNetworks struct {
	Regions   map[string]OpenstackRegion
	Neighbors interface{}
	Hosts     interface{}
}

type OpenstackRegion struct {
	Tenants          interface{}
	ServiceEndPoints interface{}
	RegionName       string
}

func call(url string, cmds []string, format string) *http.Response {
	p := Parameters{1, cmds, format}
	req := Request{"2.0", "runCmds", p, "1"}
	buf, err := json.Marshal(req)
	resp := new(http.Response)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	if strings.HasPrefix(url, "https") {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}
	resp, err = client.Post(url, "application/json", bytes.NewReader(buf))

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return resp
}

func RawCall(url string, cmds []string, format string) []byte {
	resp := call(url, cmds, format)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error calling")
		fmt.Println(err)
	}
	return body
}

func Call(url string, cmds []string, format string) JsonRpcResponse {
	p := Parameters{1, cmds, format}
	req := Request{"2.0", "runCmds", p, "1"}
	buf, err := json.Marshal(req)
	resp := new(http.Response)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	if strings.HasPrefix(url, "https") {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}
	resp, err = client.Post(url, "application/json", bytes.NewReader(buf))
	defer resp.Body.Close()

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return decodeEapiResponse(resp)
}

func CallShowInterfaces(url, intf string) ShowInterfacesResponse {
	cmds := []string{"enable", "show interfaces " + intf}
	resp := call(url, cmds, "json")
	dec := json.NewDecoder(resp.Body)
	var v ShowInterfacesResponse
	if err := dec.Decode(&v); err != nil {
		log.Println(err)
	}
	return v
}

func CallShowVlan(url, intf string) ShowVlanResponse {
	cmds := []string{"enable", "show vlan " + intf}
	resp := call(url, cmds, "json")
	dec := json.NewDecoder(resp.Body)
	var v ShowVlanResponse
	if err := dec.Decode(&v); err != nil {
		log.Println(err)
	}
	return v
}

func decodeEapiResponse(resp *http.Response) JsonRpcResponse {
	dec := json.NewDecoder(resp.Body)
	var v JsonRpcResponse
	if err := dec.Decode(&v); err != nil {
		log.Println(err)
	}
	return v
}
