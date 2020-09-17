package metrics

import (
	"github.com/stundzia/headless_stuff/headless_service/errors"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus"
)


var (
	SurfGetResTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "surf_get_res_total",
		Help: "The total number of processed Surf jobs.",
	},
	[]string{"res", "code"},
	)
)

var (
	SurfLatencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "surf_job_latency",
		Help:    "Time taken to complete Surf job.",
		Buckets: []float64{2, 4, 6, 8, 10}, //defining small buckets as this app should not take more than 1 sec to respond
	}, []string{"code"}) // this will be partitioned by the HTTP code.
)

var (
	ChromeDPGetResOkTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chromedp_get_res_ok_total",
		Help: "The total number of processed ChromeDP jobs.",
	})
)

func RunMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2112", nil)
	errors.HandleGenericError(err)
}