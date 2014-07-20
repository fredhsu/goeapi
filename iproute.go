package eapi

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type NextHop struct {
	Address   net.IP
	Interface string
}

type Route struct {
	Codes     []string
	Network   *net.IPNet
	AdminDist int
	Metric    int
	NextHops  []NextHop
}

func (r *Route) appendNextHop(nh NextHop) {
	r.NextHops = append(r.NextHops, nh)
}

func ParseCodes(fields []string) ([]string, int) {
	var netIndex int
	codes := []string{}
	for i, f := range fields {
		if len(f) <= 2 {
			codes = append(codes, f)
		} else {
			netIndex = i
			break
		}
	}
	return codes, netIndex
}

func ParseMetric(s string) (int, int) {
	s = strings.Trim(s, "[]")
	values := strings.Split(s, "/")
	ad, err := strconv.Atoi(values[0])
	if err != nil {
		fmt.Println(err)
	}
	metric, err := strconv.Atoi(values[1])
	if err != nil {
		fmt.Println(err)
	}
	return ad, metric
}

func NewIpRoute(s string) Route {
	var metric, ad int
	var nh NextHop
	fields := strings.Fields(s)
	codes, index := ParseCodes(fields)
	_, ipnet, err := net.ParseCIDR(fields[index])
	if err != nil {
		fmt.Println(err)
	}
	index += 1
	if fields[index] == "is" {
		// Directly connected
		ad = 0
		metric = 0
		nh = NextHop{Interface: fields[len(fields)-1]}
	} else {
		ad, metric = ParseMetric(fields[index])
		index += 1
		nh = ParseNextHop(fields[index:])
	}
	return Route{Codes: codes, Network: ipnet, AdminDist: ad, Metric: metric, NextHops: []NextHop{nh}}
}

func ParseNextHop(fields []string) NextHop {
	if fields[0] != "via" {
		fmt.Println("Error, not the right format")
	}
	nhIp := net.ParseIP(strings.TrimSuffix(fields[1], ","))
	nhInt := fields[2]
	return NextHop{Address: nhIp, Interface: nhInt}
}

func ParseShowIpRoute(s string) []Route {
	routes := []Route{}
	reader := strings.NewReader(s)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		// First scan up to Gateway of Last Resort
		if strings.HasPrefix(scanner.Text(), "Gateway") {
			break
		}

	}
	var r Route
	for scanner.Scan() {
		// Next skip empty lines
		s := strings.TrimSpace(scanner.Text())
		if s != "" {
			if strings.HasPrefix(s, "via") {
				r.appendNextHop(ParseNextHop(strings.Fields(s)))
			} else {
				if len(r.Codes) > 0 {
					routes = append(routes, r)
				}
				r = NewIpRoute(s)
			}
		}
	}
	return routes
}
