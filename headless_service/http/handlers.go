package http

import (
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/stundzia/headless_stuff/headless_service/engines"
	"github.com/stundzia/headless_stuff/headless_service/errors"
	"github.com/stundzia/headless_stuff/headless_service/models"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const DefaultWaitTime = 3

var static = packr.NewBox("./")

func getWaitTimeFromStr(waitTimeStr string) time.Duration {
	reqWaitTime, err := strconv.ParseInt(waitTimeStr, 10, 32)
	waitTime := DefaultWaitTime
	if err == nil {
		waitTime = int(reqWaitTime)
	}
	res := time.Duration(waitTime) * time.Second
	return res
}

func getWaitTime(waitTime float64) time.Duration {
	return time.Duration(waitTime) * time.Second
}

func renderGet(w http.ResponseWriter, r *http.Request) {

	var requestJson models.RenderRequest
	body, err := ioutil.ReadAll(r.Body)
	errors.HandleGenericError(err)

	err = json.Unmarshal(body, &requestJson)
	errors.HandleJSONError(err)
	waitTime := getWaitTime(requestJson.WaitTime)
	engine := requestJson.Engine

	var headlessRes models.HeadlessResponse
	if engine == "surf" {
		headlessRes = engines.RenderPage(
			requestJson.Url,
			waitTime,
			requestJson.ProxyUrl,
			requestJson.UserAgentType,
			requestJson.UserAgent,
		)
	} else {
		headlessRes = engines.RenderPageChromeDp(
			requestJson.Url,
			waitTime,
			requestJson.ProxyUrl,
			requestJson.UserAgentType,
			requestJson.UserAgent,
			requestJson.Headless,
		)
	}
	msg := map[string]string{}
	msg["url_requested"] = requestJson.Url
	msg["content"] = headlessRes.ContentString
	msg["title"] = headlessRes.Title
	msg["status_code"] = strconv.Itoa(headlessRes.StatusCode)
	res, err := json.Marshal(msg)
	errors.HandleGenericError(err)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	errors.HandleResponseWriteError(err)
}

func index(w http.ResponseWriter, r *http.Request) {
	message, err := static.FindString("static/html/index.html")

	// FIXME: test stuff, remove later.
	fmt.Println(r.Header.Get("user-agent"))


	errors.HandleGenericError(err)
	_, err = w.Write([]byte(message))
	errors.HandleResponseWriteError(err)
}