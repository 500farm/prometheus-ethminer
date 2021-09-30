# ethminer_exporter
Prometheus exporter reporting Ethminer stats.  

### Setup

```sh
$ make 
$ sudo make install
```

After that, the exporter will be started automatically by systemd on startup.

### Usage

By default, the `ethminer_exporter` listens on port 8555.

API endpoints are discovered automatically by looking for open port 3333 on localhost and running Docker containers.

For the explanation of Ethminer API output, see https://github.com/ethereum-mining/ethminer/blob/master/docs/API_DOCUMENTATION.md#miner_getstatdetail.

CLI args (specify in `/etc/systemd/system/ethminer_exporter.service`):

```
--listen
    Address to listen on (default 8555).

--discover-api-port
    Port on which to look for Ethminer API on localhost and running docker containers (default 3333).

--net-timeout
    Connection and read timeout for Ethminer API (default 1s)
```
