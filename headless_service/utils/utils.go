package utils

import (
	"regexp"
	"fmt"
	"strconv"
	"encoding/base64"
)

type Proxy struct {
	Url string
	UserPassString string
	Host string
	Port int
}

func (proxy *Proxy) BasicAuthString() string {
	userPass := base64.StdEncoding.EncodeToString([]byte(proxy.UserPassString))
	return fmt.Sprintf("Basic %s", userPass)
}

func ParseProxyUrlComponents(proxyUrl string) *Proxy {
	fmt.Println(proxyUrl)
	reg, err := regexp.Compile(`([^\/@,;]+:[^\/@,;]+)@([^,;:]+):([^,;:]+)`)
	if err != nil {
		panic(err)
	}
	res := reg.FindAllStringSubmatch(proxyUrl, -1)
	userpass := ""
	host := ""
	port := 0
	if len(res) > 0 && len(res[0]) == 4 {
		userpass = res[0][1]
		host = res[0][2]
		portString := res[0][3]
		portI64, err := strconv.ParseInt(portString, 10, 32)
		if err == nil {
			port = int(portI64)
		}
	}
	return &Proxy{
		Url:proxyUrl,
		UserPassString:userpass,
		Host:host,
		Port:port,
	}
}