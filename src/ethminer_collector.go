package main

import (
	"encoding/json"
	"net"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	namespace = "ethminer_"
)

type EthminerCollector struct {
	target string
	started_timestamp, connected, last_share_timestamp, hashrate, found_shares_total, rejected_shares_total,
	failed_shares_total, fan_speed_percent, power_draw_watts, temperature_degrees, paused *prometheus.Desc
}

func newEthminerCollector(target string) (*EthminerCollector, error) {
	deviceLabels := []string{"device", "name", "type", "mode"}

	return &EthminerCollector{
		target: target,

		// Global
		started_timestamp: prometheus.NewDesc(
			namespace+"started_timestamp",
			"Ethminer start time (unix timestamp)",
			[]string{"version"},
			nil,
		),
		connected: prometheus.NewDesc(
			namespace+"connected",
			"Is Ethminer connected to the pool",
			[]string{"uri"},
			nil,
		),

		// Per-device
		last_share_timestamp: prometheus.NewDesc(
			namespace+"last_share_timestamp",
			"Per-device: Last found share time (unix timestamp)",
			deviceLabels,
			nil,
		),
		hashrate: prometheus.NewDesc(
			namespace+"hashrate",
			"Per-device: Hashrate (H/s)",
			deviceLabels,
			nil,
		),
		found_shares_total: prometheus.NewDesc(
			namespace+"found_shares_total",
			"Per-device: Number of found shares",
			deviceLabels,
			nil,
		),
		rejected_shares_total: prometheus.NewDesc(
			namespace+"rejected_shares_total",
			"Per-device: Number of shares rejected by the pool",
			deviceLabels,
			nil,
		),
		failed_shares_total: prometheus.NewDesc(
			namespace+"failed_shares_total",
			"Per-device: Number of failed shares (always 0 if --no-eval is set)",
			deviceLabels,
			nil,
		),
		fan_speed_percent: prometheus.NewDesc(
			namespace+"fan_speed_percent",
			"Per-device: Fan speed (0-100%)",
			deviceLabels,
			nil,
		),
		power_draw_watts: prometheus.NewDesc(
			namespace+"power_draw_watts",
			"Per-device: Power draw (W)",
			deviceLabels,
			nil,
		),
		temperature_degrees: prometheus.NewDesc(
			namespace+"temperature_degrees",
			"Per-device: Temperature (degrees celsius)",
			deviceLabels,
			nil,
		),
		paused: prometheus.NewDesc(
			namespace+"paused",
			"Per-device: Is device paused",
			append(deviceLabels, "reason"),
			nil,
		),
	}, nil
}

func (e *EthminerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.started_timestamp
	ch <- e.connected

	ch <- e.last_share_timestamp
	ch <- e.hashrate
	ch <- e.found_shares_total
	ch <- e.rejected_shares_total
	ch <- e.failed_shares_total
	ch <- e.fan_speed_percent
	ch <- e.power_draw_watts
	ch <- e.temperature_degrees
	ch <- e.paused
}

func (e *EthminerCollector) Collect(ch chan<- prometheus.Metric) {
	conn, err := net.Dial("tcp", e.target)
	if err != nil {
		log.Errorln(err)
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("connection_error", "Error connecting to target", nil, nil), err)
		return
	}
	defer conn.Close()

	message := "{\"id\":0, \"jsonrpc\": \"2.0\", \"method\":\"miner_getstatdetail\"}\n"
	conn.Write([]byte(message))
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}

	ethstats := new(EthminerAPIResponse)
	if err := json.Unmarshal(buf[:n], ethstats); err != nil {
		log.Errorln(err)
	}
	if ethstats.error.code != 0 {
		log.Errorln(ethstats.error.message)
	}

	result := ethstats.result

	ch <- prometheus.MustNewConstMetric(
		e.started_timestamp,
		prometheus.GaugeValue,
		float64(time.Now().Unix()-int64(result.host.runtime)),
		result.host.version,
	)
	ch <- prometheus.MustNewConstMetric(
		e.connected,
		prometheus.GaugeValue,
		float64(result.connection.connected),
		result.connection.uri,
	)

	for _, device := range result.devices {
		labelValues := []string{
			device.hardware.pci,
			device.hardware.name,
			device.hardware._type,
			device._mode,
		}
		ch <- prometheus.MustNewConstMetric(
			e.last_share_timestamp,
			prometheus.GaugeValue,
			float64(time.Now().Unix()-int64(device.mining.shares[3])),
			labelValues...,
		)
		ch <- prometheus.MustNewConstMetric(
			e.hashrate,
			prometheus.GaugeValue,
			float64(parseHashrate(device.mining.hashrate)),
			labelValues...,
		)
		ch <- prometheus.MustNewConstMetric(
			e.found_shares_total,
			prometheus.CounterValue,
			float64(device.mining.shares[0]),
			labelValues...,
		)
		ch <- prometheus.MustNewConstMetric(
			e.rejected_shares_total,
			prometheus.CounterValue,
			float64(device.mining.shares[1]),
			labelValues...,
		)
		ch <- prometheus.MustNewConstMetric(
			e.failed_shares_total,
			prometheus.CounterValue,
			float64(device.mining.shares[2]),
			labelValues...,
		)
		ch <- prometheus.MustNewConstMetric(
			e.fan_speed_percent,
			prometheus.GaugeValue,
			float64(device.hardware.sensors[1]),
			labelValues...,
		)
		ch <- prometheus.MustNewConstMetric(
			e.power_draw_watts,
			prometheus.GaugeValue,
			float64(device.hardware.sensors[2]),
			labelValues...,
		)
		ch <- prometheus.MustNewConstMetric(
			e.temperature_degrees,
			prometheus.GaugeValue,
			float64(device.hardware.sensors[0]),
			labelValues...,
		)
		ch <- prometheus.MustNewConstMetric(
			e.paused,
			prometheus.GaugeValue,
			float64(device.mining.paused),
			append(labelValues, device.mining.pauseReason)...,
		)
	}
}
