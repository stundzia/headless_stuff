package main

import (
	"github.com/stundzia/headless_stuff/headless_service/engines"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var result uint64

func load(site <-chan string, wg *sync.WaitGroup) {
	//ctx, cancel := chromedp.NewContext(context.Background())
	//defer cancel()

	for s := range site {
		//var res string
		res := engines.RenderPageChromeDp(s, 0, "", "desktop", "", true)


		if len(res.ContentString) > 3000 {
			atomic.AddUint64(&result, 1)
		}
		wg.Done()
	}
}

func performanceTest() {
	workStream := make(chan string)
	wg := &sync.WaitGroup{}

	workers := 10
	work := 1000

	wg.Add(work)

	for i := 0; i < workers; i++ {
		go load(workStream, wg)
	}

	start := time.Now()
	for i := 0; i < work; i++ {
		workStream <- "https://www.google.com/search?q=adidas"
	}
	close(workStream)
	wg.Wait()

	log.Println("FINISHED ", time.Since(start), " result ",  atomic.LoadUint64(&result))
}

func main() {
	//performanceTest()
	engines.TryCDP()
	//http.StartServer("8099")
}
