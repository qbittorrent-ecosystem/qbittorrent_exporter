package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/ultram4rine/qbittorrent_exporter/client"
	"github.com/ultram4rine/qbittorrent_exporter/collector"
)

func getEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

const html = `
	<!DOCTYPE html>
	<title>qBittorrent Exporter</title>
	<h1>qBittorrent Exporter</h1>
	<p><a href=/metrics>Metrics</a></p>
`

func main() {
	var (
		qbittorrentAddr = getEnv("QBITTORRENT_ADDR", "http://localhost:8080")
		qBittorrentUser = getEnv("QBITTORRENT_USER", "")
		qBittorrentPass = getEnv("QBITTORRENT_PASS", "")
		exporterPort    = getEnv("EXPORTER_PORT", ":9177")
		metricsPrefix   = getEnv("METRICS_PREFIX", "qbittorrent")
	)

	c, err := client.NewQBittorrentClient(qbittorrentAddr, qBittorrentUser, qBittorrentPass)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, html); err != nil {
			log.Warnf("Error while sending a response for the '/' path: %v", err)
		}
	})

	collector := collector.NewQBittorrentCollector(c, metricsPrefix, make(map[string]string))
	prometheus.MustRegister(collector)
	log.Fatal(http.ListenAndServe(exporterPort, nil))
}
