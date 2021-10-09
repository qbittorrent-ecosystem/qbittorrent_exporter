package collector

import (
	"fmt"
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/ultram4rine/qbittorrent_exporter/client"
)

type QBittorrentCollector struct {
	namespace string

	qBittorrentClient *client.QBittorrentClient
	upMetric          prometheus.Gauge
	connectedMetric   prometheus.Gauge
	firewalledMetric  prometheus.Gauge
	metrics           map[string]*prometheus.Desc

	mutex sync.Mutex
}

// NewNginxCollector creates an NginxCollector.
func NewQBittorrentCollector(qBittorrentClient *client.QBittorrentClient, namespace string, constLabels map[string]string) *QBittorrentCollector {
	return &QBittorrentCollector{
		namespace:         namespace,
		qBittorrentClient: qBittorrentClient,
		upMetric: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   namespace,
			Name:        "up",
			Help:        "Whether if server is alive or not",
			ConstLabels: constLabels,
		}),
		connectedMetric: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   namespace,
			Name:        "connected",
			Help:        "Whether if server is connected or not",
			ConstLabels: constLabels,
		}),
		firewalledMetric: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   namespace,
			Name:        "firewalled",
			Help:        "Whether if server is under a firewall or not",
			ConstLabels: constLabels,
		}),
		metrics: map[string]*prometheus.Desc{
			"dht_nodes":          newMetric(namespace, "dht_nodes", "DHT nodes connected to", constLabels),
			"dl_info_data_total": newMetric(namespace, "dl_info_data_total", "Data downloaded this session (bytes)", constLabels),
			"up_info_data_total": newMetric(namespace, "up_info_data_total", "Data uploaded this session (bytes)", constLabels),
		},
	}
}

// Describe sends the super-set of all possible descriptors of qBittorrent metrics
// to the provided channel.
func (c *QBittorrentCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.upMetric.Desc()
	for _, m := range c.metrics {
		ch <- m
	}
}

// Collect fetches metrics from qBittorrent and sends them to the provided channel.
func (c *QBittorrentCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock() // To protect metrics from concurrent collects
	defer c.mutex.Unlock()

	status, err := c.qBittorrentClient.GetStatus()
	if err != nil {
		c.upMetric.Set(0)
		ch <- c.upMetric
		log.Printf("Error getting status: %v", err)
		return
	}

	categories, err := c.qBittorrentClient.GetCategories()
	if err != nil {
		c.upMetric.Set(0)
		ch <- c.upMetric
		log.Printf("Error getting torrents: %v", err)
		return
	}

	torrents, err := c.qBittorrentClient.GetTorrents()
	if err != nil {
		c.upMetric.Set(0)
		ch <- c.upMetric
		log.Printf("Error getting torrents: %v", err)
		return
	}

	c.upMetric.Set(1)
	ch <- c.upMetric

	var connected float64 = 0
	if status.Connection == "connected" {
		connected = 1
	}
	c.connectedMetric.Set(connected)
	ch <- c.connectedMetric

	var firewalled float64 = 0
	if status.Connection == "firewalled" {
		connected = 1
	}
	c.firewalledMetric.Set(firewalled)
	ch <- c.firewalledMetric

	ch <- prometheus.MustNewConstMetric(c.metrics["dht_nodes"],
		prometheus.GaugeValue, float64(status.DHTNodes))
	ch <- prometheus.MustNewConstMetric(c.metrics["dl_info_data_total"],
		prometheus.CounterValue, float64(status.Downloaded))
	ch <- prometheus.MustNewConstMetric(c.metrics["up_info_data_total"],
		prometheus.CounterValue, float64(status.Uploaded))

	categories = append(categories, "Uncategorized")
	for _, category := range categories {
		var category_torrents []client.Torrent
		for _, t := range torrents {
			if t.Category == category || (t.Category == "" && category == "Uncategorized") {
				category_torrents = append(category_torrents, t)
			}
		}

		for _, status := range TORRENT_STATUSES {
			var (
				status_prop     = fmt.Sprintf("is_%s", status)
				status_torrents []client.Torrent
			)

			for _, t := range category_torrents {
				check := stateFuncs[status_prop]
				if check(t.State) {
					status_torrents = append(status_torrents, t)
				}
			}
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					c.namespace+"_torrents_count",
					"Number of torrents in status under category",
					nil,
					prometheus.Labels{
						"status":   status,
						"category": category,
					},
				),
				prometheus.GaugeValue, float64(len(status_torrents)),
			)
		}
	}
}
