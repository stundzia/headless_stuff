package engines

import (
	"fmt"
	"github.com/headzoo/surf/browser"
	"gopkg.in/headzoo/surf.v1"
	"github.com/stundzia/headless_stuff/headless_service/headers"
	"github.com/stundzia/headless_stuff/headless_service/metrics"
	"github.com/stundzia/headless_stuff/headless_service/models"
	"net/http"
	"net/url"
	"time"
)


func setBowHeaders(bow *browser.Browser, headers map[string]string) {
	for key, value := range headers {
		bow.AddRequestHeader(key, value)
	}
}

func RenderPage(targetUrl string, waitTime time.Duration, proxyUrlString string, uaType string, uaString string) models.HeadlessResponse {
	start := time.Now()

	bow := surf.NewBrowser()

	ua := uaString
	if len(uaString) < 3 {
		ua = headers.GetUaByUaType(uaType)
	}
	bow.SetUserAgent(ua)
	setBowHeaders(bow, headers.DefaultHeaders)

	if len(proxyUrlString) > 5 {
		proxyUrl, err := url.Parse(proxyUrlString)
		if err == nil {
			proxy := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
			bow.SetTransport(proxy)
		} else {
			fmt.Println("Unable to parse proxy.")
		}
	}

	err := bow.Open(targetUrl)
	if err != nil {
		panic(err)
	}

	//errclick := bow.Click("input#oneway-flight")
	//if errclick != nil {
	//	fmt.Println(errclick)
	//}

	time.Sleep(waitTime)

	status_str := ""
	if bow.StatusCode() == 200 {
		status_str = "ok"
	} else {
		status_str = "not_ok"
	}

	go func() {
		metrics.SurfGetResTotal.WithLabelValues(fmt.Sprintf("%s", status_str), fmt.Sprintf("%d", bow.StatusCode())).Inc()
	}()

	defer func() {
		httpDuration := time.Since(start)
		metrics.SurfLatencyHistogram.WithLabelValues(fmt.Sprintf("%d", bow.StatusCode())).Observe(httpDuration.Seconds())
	}()

	resp := models.HeadlessResponse{
		TargetUrl:  targetUrl,
		ResUrl:     bow.Url().String(),
		Body:       bow.Body(),
		Title:      bow.Title(),
		StatusCode: bow.StatusCode(),
	}

	_, err = bow.Download(&resp)

	return resp
}