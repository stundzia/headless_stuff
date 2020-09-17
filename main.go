package main

import (
	"github.com/stundzia/headless_stuff/headless_service/http"
	"github.com/stundzia/headless_stuff/headless_service/metrics"
)

func main() {
	go metrics.RunMetrics()
	http.StartServer("8099")
}
