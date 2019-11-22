package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type SpeedWiFiHomeCollector struct {
	monthClient MonthClient

	TotalUploadBytes   prometheus.Counter
	TotalDownloadBytes prometheus.Counter
}

func NewSpeedWiFiHomeCollector(namespace string) SpeedWiFiHomeCollector {
	return SpeedWiFiHomeCollector{
		NewMonthClient(),
		prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "total_upload_bytes",
		}),
		prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "total_download_bytes",
		}),
	}
}

func (c SpeedWiFiHomeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.TotalUploadBytes.Desc()
	ch <- c.TotalDownloadBytes.Desc()
}

func (c SpeedWiFiHomeCollector) Collect(ch chan<- prometheus.Metric) {
	uploaded, downloaded, err := c.monthClient.Collect()
	if err != nil {
		log.Printf("Failed to fetch: %s\n", err.Error())
		return
	}

	ch <- prometheus.MustNewConstMetric(c.TotalUploadBytes.Desc(), prometheus.CounterValue, float64(uploaded))
	ch <- prometheus.MustNewConstMetric(c.TotalDownloadBytes.Desc(), prometheus.CounterValue, float64(downloaded))
}

func main() {
	listen := flag.String("listen", "127.0.0.1:9999", "The address to listen")
	flag.Parse()

	c := NewSpeedWiFiHomeCollector("speed_wifi_home")
	prometheus.MustRegister(c)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<h1>Speed Wi-Fi Home Exporter</h1><a href="/metrics">metrics</a>`)
	})

	log.Printf("listen on %s", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
