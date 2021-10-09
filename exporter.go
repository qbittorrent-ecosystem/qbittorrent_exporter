package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/ultram4rine/qbittorrent_exporter/client"
	"github.com/ultram4rine/qbittorrent_exporter/collector"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	address  = kingpin.Flag("address", "Address of qBittorrent.").Short('a').Required().String()
	username = kingpin.Flag("username", "Username for qBittorrent.").Short('u').Required().String()
	password = kingpin.Flag("password", "Password for qBittorrent.").Short('p').Required().String()
)

func main() {
	kingpin.Parse()

	c, err := client.NewQBittorrentClient(*address, *username, *password)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/metrics", promhttp.Handler())
	collector := collector.NewQBittorrentCollector(c, "qbittorrent", make(map[string]string))
	prometheus.MustRegister(collector)
	log.Fatal(http.ListenAndServe(":9177", nil))
}
