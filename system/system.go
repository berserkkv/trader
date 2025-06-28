package system

import (
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type SystemUsage struct {
	CpuUsage  float64 `json:"cpu_usage"`
	MemUsage  float64 `json:"mem_usage"`
	DiskUsage float64 `json:"disk_usage"`
}

func GetSystemUsage(c *gin.Context) {

	cpuUsage, _ := GetCPUUsage()
	_, _, memUsage, _ := GetMemUsage()
	_, _, diskUsage, _ := GetDiskUsage("/")

	s := SystemUsage{
		CpuUsage:  cpuUsage,
		MemUsage:  memUsage,
		DiskUsage: diskUsage,
	}

	c.JSON(http.StatusOK, s)

}

func GetCPUUsage() (float64, error) {
	readCPU := func() (float64, float64, error) {
		data, err := ioutil.ReadFile("/proc/stat")
		if err != nil {
			return 0, 0, err
		}
		fields := strings.Fields(string(data))
		if len(fields) < 3 {
			return 0, 0, errors.New("invalid /proc/stat output")
		}
		user, _ := strconv.ParseFloat(fields[1], 64)
		nice, _ := strconv.ParseFloat(fields[2], 64)
		system, _ := strconv.ParseFloat(fields[3], 64)
		idle, _ := strconv.ParseFloat(fields[4], 64)
		iowait, _ := strconv.ParseFloat(fields[5], 64)
		irq, _ := strconv.ParseFloat(fields[6], 64)
		softirq, _ := strconv.ParseFloat(fields[7], 64)

		total := user + nice + system + idle + iowait + irq + softirq
		return total, idle, nil
	}

	total1, idle1, _ := readCPU()
	// wait 100ms
	select {
	case <-time.After(100 * 1e6):
	}
	total2, idle2, _ := readCPU()
	totalDelta := total2 - total1
	idleDelta := idle2 - idle1
	usage := (1.0 - idleDelta/totalDelta) * 100
	return usage, nil
}

func GetMemUsage() (total, used, usagePercent float64, err error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	mem := make(map[string]float64)
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		val, _ := strconv.ParseFloat(fields[1], 64)
		mem[fields[0][:len(fields[0])-1]] = val
	}
	total = mem["MemTotal"] / 1e6
	free := mem["MemFree"] + mem["Buffers"] + mem["Cached"]
	used = (mem["MemTotal"] - free) / 1e6
	usagePercent = (used / total) * 100
	return
}

func GetDiskUsage(path string) (total, used, usagePercent float64, err error) {
	var statfs syscall.Statfs_t
	err = syscall.Statfs(path, &statfs)
	if err != nil {
		return
	}
	total = float64(statfs.Blocks*uint64(statfs.Bsize)) / 1e9
	free := float64(statfs.Bfree*uint64(statfs.Bsize)) / 1e9
	used = total - free
	usagePercent = (used / total) * 100
	return
}
