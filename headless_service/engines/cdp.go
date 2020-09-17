package engines

import (
	"context"
	"fmt"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"time"

	"github.com/stundzia/headless_stuff/headless_service/errors"
	"github.com/stundzia/headless_stuff/headless_service/models"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/rpcc"
)

func getClient() (c *cdp.Client, ctx context.Context, cancel context.CancelFunc, conn *rpcc.Conn) {
	// Launch chrome first with `google-chrome --headless --remote-debugging-port=9222`
	timeout := 60 * time.Second
	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Use the DevTools HTTP/JSON API to manage targets (e.g. pages, webworkers).
	devt := devtool.New("http://127.0.0.1:9222")
	pt, err := devt.Get(ctx, devtool.Page)
	if err != nil {
		pt, err = devt.Create(ctx)
		errors.HandleCDPError(err, "devt.Create")
	}

	// Initiate a new RPC connection to the Chrome DevTools Protocol target.
	conn, err = rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	errors.HandleCDPError(err, "conn init")

	c = cdp.NewClient(conn)
	return c, ctx, cancel, conn
}

func RenderPageChromeCDP(targetUrl string, waitTime time.Duration, proxyUrlString string, uaType string, uaString string, headless bool) models.HeadlessResponse {
	c, ctx, cancel, conn := getClient()
	fmt.Println(c)
	fmt.Println(ctx)
	defer cancel()
	defer conn.Close() // Leaving connections open will leak memory.

	// Open a DOMContentEventFired client to buffer this event.
	domContent, err := c.Page.DOMContentEventFired(ctx)
	errors.HandleCDPError(err, "domContent")
	defer domContent.Close()

	// Enable events on the Page domain, it's often preferrable to create
	// event clients before enabling events so that we don't miss any.
	err = c.Page.Enable(ctx)
	errors.HandleGenericError(err)

	// Create the Navigate arguments with the optional Referrer field set.
	navArgs := page.NewNavigateArgs(targetUrl).
		SetReferrer("https://duckduckgo.com")
	nav, err := c.Page.Navigate(ctx, navArgs)
	errors.HandleCDPError(err,"Page.Navigate")

	// Wait until we have a DOMContentEventFired event.
	_, err = domContent.Recv()
	errors.HandleCDPError(err,"domContent.Recv")

	fmt.Printf("Page loaded with frame ID: %s\n", nav.FrameID)

	// Fetch the document root node. We can pass nil here
	// since this method only takes optional arguments.
	doc, err := c.DOM.GetDocument(ctx, nil)
	errors.HandleCDPError(err,"DOM.GetDocument")

	// Get the outer HTML for the page.
	result, err := c.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{
		NodeID: &doc.Root.NodeID,
	})
	errors.HandleCDPError(err,"GetOuterHTML")

	fmt.Printf("HTML: %s\n", result.OuterHTML)


	// Capture a screenshot of the current page.
	//screenshotName := "screenshot.jpg"
	//screenshotArgs := page.NewCaptureScreenshotArgs().
	//	SetFormat("jpeg").
	//	SetQuality(80)
	//screenshot, err := c.Page.CaptureScreenshot(ctx, screenshotArgs)
	//errors.HandleCDPError(err, "CaptureScreenshot")
	//
	//err = ioutil.WriteFile(screenshotName, screenshot.Data, 0644)
	//errors.HandleGenericError(err)




	userAgent := uaString
	headersObj := getHeaders(uaType)
	fmt.Println(userAgent, headersObj)

	resp := models.HeadlessResponse{
		TargetUrl:     targetUrl,
		ResUrl:        "resUrl",
		Body:          result.OuterHTML,
		Title:         "title",
		StatusCode:    200,
		ContentString: "res",
	}

	return resp
}