package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		http.Error(w, fmt.Sprintf("target is not specified."), http.StatusBadRequest)
		return
	}
	e, _ := newEthminerCollector(target)

	registry := prometheus.NewRegistry()
	registry.MustRegister(e)

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func main() {
	listenAddress := kingpin.Flag("listen", "Address to listen on for web interface and telemetry.").Default("0.0.0.0:8555").String()
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
