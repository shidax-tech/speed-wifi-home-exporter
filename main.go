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

	ErrorCount           prometheus.Counter
	TotalUploadBytes     prometheus.Counter
	TotalDownloadBytes   prometheus.Counter
	MonthlyUploadBytes   prometheus.Gauge
	MonthlyDownloadBytes prometheus.Gauge
}

func NewSpeedWiFiHomeCollector(namespace string, address string) SpeedWiFiHomeCollector {
	return SpeedWiFiHomeCollector{
		NewMonthClient(address),
		prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "error_count",
		}),
		prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "total_upload_bytes",
		}),
		prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "total_download_bytes",
		}),
		prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "monthly_upload_bytes",
		}),
		prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "monthly_download_bytes",
		}),
	}
}

func (c SpeedWiFiHomeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.ErrorCount.Desc()
	ch <- c.TotalUploadBytes.Desc()
	ch <- c.TotalDownloadBytes.Desc()
	ch <- c.MonthlyUploadBytes.Desc()
	ch <- c.MonthlyDownloadBytes.Desc()
}

func (c SpeedWiFiHomeCollector) Collect(ch chan<- prometheus.Metric) {
	stat, err := c.monthClient.Collect()
	if err != nil {
		log.Printf("Failed to fetch: %s\n", err.Error())
		c.ErrorCount.Inc()
	}

	c.ErrorCount.Collect(ch)

	ch <- prometheus.MustNewConstMetric(c.TotalUploadBytes.Desc(), prometheus.CounterValue, float64(stat.TotalUploaded))
	ch <- prometheus.MustNewConstMetric(c.TotalDownloadBytes.Desc(), prometheus.CounterValue, float64(stat.TotalDownloaded))

	c.MonthlyUploadBytes.Set(float64(stat.MonthlyUploaded))
	c.MonthlyUploadBytes.Collect(ch)
	c.MonthlyDownloadBytes.Set(float64(stat.MonthlyDownloaded))
	c.MonthlyDownloadBytes.Collect(ch)
}

func main() {
	listen := flag.String("listen", "127.0.0.1:9999", "The address to listen")
	address := flag.String("address", "192.168.100.1", "The address of router")
	flag.Parse()

	c := NewSpeedWiFiHomeCollector("speed_wifi_home", *address)
	prometheus.MustRegister(c)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<h1>Speed Wi-Fi Home Exporter</h1><a href="/metrics">metrics</a>`)
	})

	log.Printf("listen on %s", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
