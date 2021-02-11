# ethminer_exporter
Prometheus exporter reporting Ethminer stats.  

For the list and explanation of fields, see https://github.com/ethereum-mining/ethminer/blob/master/docs/API_DOCUMENTATION.md#miner_getstatdetail.

To enable JSON-RPC API, run Ethminer with `--api-port` and specify the listening port.

### Setup

```sh
$ make build
```

### Usage

```sh
$ bin/ethminer_exporter
```

By default, *ethminer_exporter* starts listening on port 8555. You can set another port with `--listen` argument:
```
$ bin/ethminer_exporter --listen 0.0.0.0:8556
```

Then, visit http://localhost:8555/metrics?target=1.2.3.4:3333 where *1.2.3.4:333* is the IP and port of Ethminer API.

