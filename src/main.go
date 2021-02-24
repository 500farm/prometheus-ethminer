package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	listenAddress = kingpin.Flag(
		"listen",
		"Address to listen on.",
	).Default("0.0.0.0:8555").String()
	discoverAPIPort = kingpin.Flag(
		"discover-api-port",
		"Port on which to look for Ethminer API on localhost and running docker containers.",
	).Default("3333").Int()
	netTimeout = kingpin.Flag(
		"net-timeout",
		"Connection and read timeout for Ethminer API.",
	).Default("1s").Duration()
	updateInterval = kingpin.Flag(
		"update-interval",
		"How often to query third-party API for updates (does not affect Ethminer metrics which are realtime).",
	).Default("1m").Duration()
)

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	registry := prometheus.NewRegistry()

	target := r.URL.Query().Get("target")
	var targets []string
	if target != "" {
		targets = []string{target}
	} else {
		targets = discoverTargets(*discoverAPIPort)
	}

	e, _ := newEthminerCollector(targets, *netTimeout)
	registry.MustRegister(e)

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func main() {
	kingpin.Version(version.Print("ethminer_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Infoln("Starting ethminer exporter")

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metricsHandler(w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head>
		<title>Ethminer Exporter</title>
		<style>
		label{
		display:inline-block;
		width:75px;
		}
		form label {
		margin: 10px;
		}
		form input {
		margin: 10px;
		}
		</style>
		</head>
		<body>
		<h1>Ethminer Exporter</h1>
		<form action="/metrics">
		<label>Target:</label> <input type="text" name="target" placeholder="X.X.X.X" value="1.2.3.4"><br>
		<input type="submit" value="Submit">
		</form>
		</body>
		</html>`))
	})

	log.Infoln("Listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
