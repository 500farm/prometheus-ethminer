module prometheus-ethminer

go 1.15

require github.com/prometheus/client_golang v1.9.0

require github.com/prometheus/common v0.17.0

require (
	github.com/containerd/containerd v1.5.6 // indirect
	github.com/docker/docker v20.10.8+incompatible
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)
