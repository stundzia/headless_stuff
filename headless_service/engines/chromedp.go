package engines

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/stundzia/headless_stuff/headless_service/errors"
	"github.com/stundzia/headless_stuff/headless_service/headers"
	"github.com/stundzia/headless_stuff/headless_service/metrics"
	"github.com/stundzia/headless_stuff/headless_service/models"
	"github.com/stundzia/headless_stuff/headless_service/utils"
	"time"
)

func getHeaders(uaType string) network.Headers {
	if uaType == "" {
		uaType = "desktop"
	}
	headersMap := headers.GetHeadersByUaType(uaType)
	headersJson, _ := json.Marshal(headersMap)
	return headersJson
}


func RenderPageChromeDp(targetUrl string, waitTime time.Duration, proxyUrlString string, uaType string, uaString string, headless bool) models.HeadlessResponse {
	userAgent := uaString
	headersObj := getHeaders(uaType)
	if len(uaString) < 3 {
		userAgent = headers.GetUaByUaType(uaType)
		fmt.Println("UA: ", userAgent)
	}
	var proxyStr string
	if len(proxyUrlString) > 5 {
		proxy := utils.ParseProxyUrlComponents(proxyUrlString)
		proxyStr = fmt.Sprintf("%s:%d", proxy.Host, proxy.Port)
		//headersObj = setProxyAuthHeader(headersObj, proxy.BasicAuthString())
		fmt.Println(headersObj)
	}
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", headless),
		chromedp.ProxyServer(proxyStr),
		chromedp.UserAgent(userAgent),
	}

	allocContext, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocContext)

	defer cancel()

	var res string
	var resUrl string
	var title string

	fmt.Println(headersObj)
	start := time.Now()
	err := chromedp.Run(
		ctx,
		network.Enable(),
		network.SetExtraHTTPHeaders(headersObj),
		chromedp.Navigate(targetUrl),
		chromedp.Location(&resUrl),
		chromedp.Sleep(waitTime),
		chromedp.WaitNotPresent(nil),
		chromedp.OuterHTML("html", &res),
		chromedp.Title(&title),
		)
	errors.HandleGenericError(err)
	fmt.Println("Done in: ", time.Now().Sub(start))

	go func() {
		metrics.ChromeDPGetResOkTotal.Inc()
	}()

	fmt.Println("Title: ", title)

	resp := models.HeadlessResponse{
		TargetUrl:     targetUrl,
		ResUrl:        resUrl,
		Body:          res,
		Title:         title,
		StatusCode:    200,
		ContentString: res,
	}

	return resp
}