package main

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type EthereumCollector struct {
	block_time_seconds              prometheus.Gauge
	block_reward_eth                prometheus.Gauge
	last_block_number               prometheus.Counter
	difficulty_hashes               prometheus.Gauge
	network_hashrate_hashes_per_sec prometheus.Gauge
	eth_price_dollars               prometheus.Gauge
}

func newEthereumCollector() (*EthereumCollector, error) {
	namespace := "ethereum_"

	return &EthereumCollector{
		block_time_seconds: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "block_time_seconds",
			Help:      "Time it took to find the last block, in seconds",
		}),
		block_reward_eth: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "block_reward_eth",
			Help:      "Reward for the last found block, in ETH",
		}),
		last_block_number: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "last_block_number",
			Help:      "Number of the last found block",
		}),
		difficulty_hashes: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "difficulty_hashes",
			Help:      "Last block difficulty in hashes",
		}),
		network_hashrate_hashes_per_sec: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "network_hashrate_hashes_per_sec",
			Help:      "Current network hasrate, in H/s",
		}),
		eth_price_dollars: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "eth_price_dollars",
			Help:      "Current ETH price, in USD",
		}),
	}, nil
}

func (e *EthereumCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.block_time_seconds.Desc()
	ch <- e.block_reward_eth.Desc()
	ch <- e.last_block_number.Desc()
	ch <- e.difficulty_hashes.Desc()
	ch <- e.network_hashrate_hashes_per_sec.Desc()
	ch <- e.eth_price_dollars.Desc()
}

func (e *EthereumCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- e.block_time_seconds
	ch <- e.block_reward_eth
	ch <- e.last_block_number
	ch <- e.difficulty_hashes
	ch <- e.network_hashrate_hashes_per_sec
	ch <- e.eth_price_dollars
}

func (e *EthereumCollector) Update(info *EthereumInfo) {
	t, _ := strconv.ParseFloat(info.BlockTime, 64)
	e.block_time_seconds.Set(t)
	e.block_reward_eth.Set(info.BlockReward)
	e.last_block_number.Set(float64(info.LastBlockNumber))
	e.difficulty_hashes.Set(info.Difficulty)
	e.network_hashrate_hashes_per_sec.Set(float64(info.NetworkHashRate))
	e.eth_price_dollars.Set(info.ETHUSDPrice)
}
