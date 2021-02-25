package main

import (
	"encoding/json"
	"errors"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	namespace = "ethminer_"
)

type EthminerCollector struct {
	targets    []string
	netTimeout time.Duration
	started_timestamp, connected, last_share_timestamp, hashrate, found_shares_total, rejected_shares_total,
	failed_shares_total, fan_speed_percent, power_draw_watts, temperature_degrees, paused *prometheus.Desc
}

func newEthminerCollector(targets []string, netTimeout time.Duration) (*EthminerCollector, error) {
	deviceLabels := []string{"api_endpoint", "device", "name", "type", "mode"}

	return &EthminerCollector{
		targets:    targets,
		netTimeout: netTimeout,

		// Global
		started_timestamp: prometheus.NewDesc(
			namespace+"started_timestamp",
			"Ethminer start time (unix timestamp)",
			[]string{"api_endpoint", "version"},
			nil,
		),
		connected: prometheus.NewDesc(
			namespace+"connected",
			"Is Ethminer connected to the pool",
			[]string{"api_endpoint", "uri"},
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
	errOut := func(code string, message string, err error) {
		log.Errorln(err)
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc(code, message, nil, nil), err)
	}

	for _, target := range e.targets {
		conn, err := net.DialTimeout("tcp", target, e.netTimeout)
		if err != nil {
			// intentionally ignore connection failures
			continue
		}
		conn.SetDeadline(time.Now().Add(e.netTimeout))
		defer conn.Close()

		message := "{\"id\":0, \"jsonrpc\": \"2.0\", \"method\":\"miner_getstatdetail\"}\n"
		conn.Write([]byte(message))
		buf := make([]byte, 65536)
		n, err := conn.Read(buf)
		if err != nil {
			errOut("read_error", "Error reading from "+target, err)
			continue
		}

		ethstats := new(EthminerAPIResponse)
		if err := json.Unmarshal(buf[:n], ethstats); err != nil {
			errOut("parse_error", "Invalid JSON from "+target, err)
			continue
		}
		if ethstats.Error.Code != 0 {
			message := ethstats.Error.Message
			errOut("api_error", "Error from Ethminer API on "+target+": "+message, errors.New(message))
			continue
		}

		result := ethstats.Result

		ch <- prometheus.MustNewConstMetric(
			e.started_timestamp,
			prometheus.GaugeValue,
			float64(time.Now().Unix()-int64(result.Host.Runtime)),
			target,
			result.Host.Version,
		)
		ch <- prometheus.MustNewConstMetric(
			e.connected,
			prometheus.GaugeValue,
			float64(boolToInt(result.Connection.Connected)),
			target,
			result.Connection.URI,
		)

		for _, device := range result.Devices {
			labelValues := []string{
				target,
				strings.ToUpper(device.Hardware.PCIID),
				regexp.MustCompile(`\s+[\d\.]+\s*GB$`).ReplaceAllString(device.Hardware.Name, ""),
				device.Hardware.Type,
				device.Mode,
			}
			ch <- prometheus.MustNewConstMetric(
				e.last_share_timestamp,
				prometheus.GaugeValue,
				float64(time.Now().Unix()-int64(device.Mining.Shares[3])),
				labelValues...,
			)
			ch <- prometheus.MustNewConstMetric(
				e.hashrate,
				prometheus.GaugeValue,
				float64(parseHashrate(device.Mining.Hashrate)),
				labelValues...,
			)
			ch <- prometheus.MustNewConstMetric(
				e.found_shares_total,
				prometheus.CounterValue,
				float64(device.Mining.Shares[0]),
				labelValues...,
			)
			ch <- prometheus.MustNewConstMetric(
				e.rejected_shares_total,
				prometheus.CounterValue,
				float64(device.Mining.Shares[1]),
				labelValues...,
			)
			ch <- prometheus.MustNewConstMetric(
				e.failed_shares_total,
				prometheus.CounterValue,
				float64(device.Mining.Shares[2]),
				labelValues...,
			)
			ch <- prometheus.MustNewConstMetric(
				e.fan_speed_percent,
				prometheus.GaugeValue,
				float64(device.Hardware.Sensors[1]),
				labelValues...,
			)
			ch <- prometheus.MustNewConstMetric(
				e.power_draw_watts,
				prometheus.GaugeValue,
				float64(device.Hardware.Sensors[2]),
				labelValues...,
			)
			ch <- prometheus.MustNewConstMetric(
				e.temperature_degrees,
				prometheus.GaugeValue,
				float64(device.Hardware.Sensors[0]),
				labelValues...,
			)
			ch <- prometheus.MustNewConstMetric(
				e.paused,
				prometheus.GaugeValue,
				float64(boolToInt(device.Mining.Paused)),
				append(labelValues, device.Mining.PauseReason)...,
			)
		}
	}
}
