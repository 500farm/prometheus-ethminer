package main

import (
	"strconv"
	"strings"
)

type EthminerAPIResponse struct {
	id      int    `json:"id"`
	jsonRPC string `json:"jsonrpc"`
	result  Result `json:"result"`
	error   Error  `json:"error"`
}

type Result struct {
	connection Connection `json:"connection"`
	devices    []Device   `json:"devices"`
	host       Host       `json:"host"`
	mining     Mining     `json:"mining"`
	monitors   Monitors   `json:"monitors"`
}

type Connection struct {
	connected int    `json:"connected"`
	switches  int    `json:"switches"`
	uri       string `json:"uri"`
}

type Device struct {
	_index   int            `json:"_index"`
	_mode    string         `json:"_mode"`
	hardware DeviceHardware `json:"hardware"`
	mining   DeviceMining   `json:"mining"`
}

type DeviceHardware struct {
	name    string `json:"name"`
	pci     string `json:"pci"`
	sensors []int  `json:"sensors"`
	_type   string `json:"type"`
}

type DeviceMining struct {
	hashrate    string   `json:"hashrate"`
	pauseReason string   `json:"pause_reason"`
	paused      int      `json:"paused"`
	segment     []string `json:"segment"`
	shares      []int    `json:"shares"`
}

type Host struct {
	name    string `json:"name"`
	runtime int    `json:"runtime"`
	version string `json:"version"`
}

type Mining struct {
	difficulty   int    `json:"difficulty"`
	epoch        int    `json:"epoch"`
	epochChanges int    `json:"epoch_changes"`
	hashrate     string `json:"hashrate"`
	shares       []int  `json:"shares"`
}

type Monitors struct {
	temperatures []int `json:"temperatures"`
}

type Error struct {
	code    int    `json:"code"`
	message string `json:"message"`
}

func parseHashrate(value string) uint64 {
	n, err := strconv.ParseUint(strings.ReplaceAll(value, "0x", ""), 16, 64)
	if err != nil {
		return 0
	}
	return n
}
