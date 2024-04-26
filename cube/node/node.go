package node

import (
	"cube/stats"
	"cube/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Node struct {
	Name            string
	Ip              string
	Api             string
	Cores           int64
	Memory          int64
	MemoryAllocated int64
	Disk            int64
	DiskAllocated   int64
	Stats           stats.Stats
	Role            string
	TaskCount       int
}

func New(name string, api string, role string) *Node {
	return &Node{
		Name: name,
		Api:  api,
		Role: role,
	}
}

func (n *Node) GetStats() (*stats.Stats, error) {
	url := fmt.Sprintf("%s/stats", n.Api)
	resp, err := utils.HTTPWithRetry(http.Get, url)
	if err != nil {
		msg := fmt.Sprintf("Unable to connect to %v. Permanent failure.\n", n.Api)
		log.Println(msg)
		return nil, errors.New(msg)
	}

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Error retrieving stats from %v: %v", n.Api, err)
		log.Println(msg)
		return nil, errors.New(msg)
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var stats stats.Stats
	err = json.Unmarshal(body, &stats)
	if err != nil {
		msg := fmt.Sprintf("error decoding message while getting stats for node %s", n.Name)
		log.Println(msg)
		return nil, errors.New(msg)
	}

	if stats.MemStats == nil || stats.DiskStats == nil {
		return nil, fmt.Errorf("error getting stats from node %s", n.Name)
	}

	n.Memory = int64(stats.MemTotalKb())
	n.Disk = int64(stats.DiskTotal())
	n.Stats = stats

	return &n.Stats, nil
}

// See discussion from this StackOverflow thread:
// https://stackoverflow.com/questions/23367857/accurate-calculation-of-cpu-usage-given-in-percentage-in-linux
func (n *Node) CalculateCpuUsage() (*float64, error) {
	stat1, err := n.GetStats()
	if err != nil {
		return nil, err
	}
	time.Sleep(3 * time.Second)

	stat2, err := n.GetStats()
	if err != nil {
		return nil, err
	}

	stat1Idle := stat1.CpuStats.Idle + stat1.CpuStats.IOWait
	stat2Idle := stat2.CpuStats.Idle + stat2.CpuStats.IOWait

	stat1NonIdle := stat1.CpuStats.User + stat1.CpuStats.Nice + stat1.CpuStats.System + stat1.CpuStats.IRQ + stat1.CpuStats.SoftIRQ + stat1.CpuStats.Steal
	stat2NonIdle := stat2.CpuStats.User + stat2.CpuStats.Nice + stat2.CpuStats.System + stat2.CpuStats.IRQ + stat2.CpuStats.SoftIRQ + stat2.CpuStats.Steal

	stat1Total := stat1Idle + stat1NonIdle
	stat2Total := stat2Idle + stat2NonIdle

	total := stat2Total - stat1Total
	idle := stat2Idle - stat1Idle

	var cpuPercentUsage float64
	if total == 0 && idle == 0 {
		cpuPercentUsage = 0.00
	} else {
		cpuPercentUsage = (float64(total) - float64(idle)) / float64(total)
	}
	return &cpuPercentUsage, nil
}
