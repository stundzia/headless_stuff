package headers

import (
	"fmt"
	"strings"
	"math/rand"
)

type HeadlessHeaders struct {
	UserAgent string `json:"user_agent"`
	AcceptLanguage string `json:"accept_language"`
	AcceptEncoding string `json:"accept_encoding"`

}

const DefaultUA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"

func GetHeadersByUaType(uaType string) map[string]string {
	typeSubtype := strings.Split(uaType, "_")
	var headers map[string]string
	if len(typeSubtype) == 1 {
		i := 0
		headersMap := HeadersMap[typeSubtype[0]]
		headerKeys := make([]string, len(headersMap))
		for k := range headersMap {
			headerKeys[i] = k
			i++
		}
		fmt.Println("head key len: ", len(headerKeys))
		headers = headersMap[headerKeys[rand.Intn(len(headerKeys))]]
	} else {
		headers = HeadersMap[typeSubtype[0]][typeSubtype[1]]
	}
	return headers
}

func GetUaByUaType(uaType string) string {
	if uaType == "" {
		uaType = "desktop"
	}
	typeSubtype := strings.Split(uaType, "_")
	var uaList []string
	if len(typeSubtype) == 1 {
		i := 0
		UasMap := UAsMap[typeSubtype[0]]
		UasKeys := make([]string, len(UasMap))
		for k := range UasMap {
			UasKeys[i] = k
			i++
		}
		uaList = UasMap[UasKeys[rand.Intn(len(UasKeys))]]
	} else {
		uaList = UAsMap[typeSubtype[0]][typeSubtype[1]]
	}
	if uaList == nil {
		return DefaultUA
	}
	return uaList[rand.Intn(len(uaList))]
}