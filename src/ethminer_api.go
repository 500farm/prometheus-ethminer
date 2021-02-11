package main

import (
	"strconv"
	"strings"
)

type EthminerAPIResponse struct {
	ID      int    `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  Result `json:"result"`
	Error   Error  `json:"error"`
}

type Result struct {
	Connection Connection `json:"connection"`
	Devices    []Device   `json:"devices"`
	Host       Host       `json:"host"`
	Mining     Mining     `json:"mining"`
	Monitors   Monitors   `json:"monitors"`
}

type Connection struct {
	Connected bool   `json:"connected"`
	Switches  int    `json:"switches"`
	URI       string `json:"uri"`
}

type Device struct {
	Index    int            `json:"_index"`
	Mode     string         `json:"_mode"`
	Hardware DeviceHardware `json:"hardware"`
	Mining   DeviceMining   `json:"mining"`
}

type DeviceHardware struct {
	Name    string     `json:"name"`
	PCIID   string     `json:"pci"`
	Sensors [3]float64 `json:"sensors"`
	Type    string     `json:"type"`
}

type DeviceMining struct {
	Hashrate    string    `json:"hashrate"`
	PauseReason string    `json:"pause_reason"`
	Paused      bool      `json:"paused"`
	Segment     [2]string `json:"segment"`
	Shares      [4]int    `json:"shares"`
}

type Host struct {
	Name    string `json:"name"`
	Runtime int    `json:"runtime"`
	Version string `json:"version"`
}

type Mining struct {
	Difficulty   float64 `json:"difficulty"`
	Epoch        int     `json:"epoch"`
	EpochChanges int     `json:"epoch_changes"`
	Hashrate     string  `json:"hashrate"`
	Shares       [4]int  `json:"shares"`
}

type Monitors struct {
	Temperatures [2]int `json:"temperatures"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func parseHashrate(value string) uint64 {
	n, err := strconv.ParseUint(strings.ReplaceAll(value, "0x", ""), 16, 64)
	if err != nil {
		return 0
	}
	return n
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
