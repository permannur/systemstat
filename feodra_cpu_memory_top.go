package parser

import (
	"fmt"
	"strconv"
	"github.com/permannur/command"
)

// FedoraTop top -b -n 1 -E k -w 512 -1 -p 0
// Top version procps-ng 3.3.17
func FedoraTop() (cpu *Cpu, memory *Memory, err error) {
	fedoraTop := &fedoraTopParser{
		cpu:    &Cpu{},
		memory: &Memory{},
	}
	err = command.NewCommand(`(Cpu(\d+).+?(\d+)\.(\d+).+?(\d+)\.(\d+).+?(\d+)\.(\d+).+?(\d+)\.(\d+).+?(\d+)\.(\d+).+?(\d+)\.(\d+).+?(\d+)\.(\d+).+?(\d+)\.(\d+))|(KiB Mem.+?(\d+) total.+?(\d+) free.+?(\d+) used.+?(\d+) buff\/cache\nKiB Swap.+?(\d+) total.+?(\d+) free.+?(\d+) used\..+?(\d+) avail)`,
		fedoraTop,
		"top",
		"-b", "-n", "1", "-E", "k", "-w", "512", "-1", "-p", "0")
	if err != nil {
		err = fmt.Errorf("FedoraTop: %s", err)
		return
	}
	cpu = fedoraTop.cpu
	memory = fedoraTop.memory
	return
}

type fedoraTopParser struct {
	cpu    *Cpu
	memory *Memory
}

func (c *fedoraTopParser) Parse(parsedList [][]string) (err error) {
	if c.cpu == nil {
		c.cpu = &Cpu{}
	}
	if c.memory == nil {
		c.memory = &Memory{Unit: "KiB"}
	}
	respCpu, respMemory := c.cpu, c.memory
	var totalSum, usedSum int64
	var cpuDetailInx int
	for _, row := range parsedList {
		if len(row) != 28 {
			err = fmt.Errorf("fedoraTopParser.Parse: parsed list length mismatch, want: %d, have: %d", 28, len(row))
			return
		}
		if len(row[1]) > 0 {
			var total, available, used, idle1, idle2 int64
			idle1, err = strconv.ParseInt(row[9], 10, 64)
			if err != nil {
				err = fmt.Errorf("fedoraTopParser.Parse: must be a number at %d index: %s", 9, row[9])
				return
			}
			idle2, err = strconv.ParseInt(row[10], 10, 64)
			if err != nil {
				err = fmt.Errorf("fedoraTopParser.Parse: must be a number at %d index: %s", 10, row[10])
				return
			}
			total = 1000
			available = idle1*10 + idle2
			used = total - available
			if len(respCpu.DetailList) == cpuDetailInx {
				respCpu.DetailList = append(respCpu.DetailList, &CpuDetail{})
			}
			respCpuDetail := respCpu.DetailList[cpuDetailInx]
			respCpuDetail.Name = "Cpu" + row[2]
			respCpuDetail.Percent = calcPercent(total, used)
			totalSum += total
			usedSum += used
			cpuDetailInx++
		}
		if len(row[19]) > 0 {
			var total, available, used int64
			total, err = strconv.ParseInt(row[20], 10, 64)
			if err != nil {
				err = fmt.Errorf("fedoraTopParser.Parse: must be a number at %d index: %s", 20, row[20])
				return
			}
			available, err = strconv.ParseInt(row[27], 10, 64)
			if err != nil {
				err = fmt.Errorf("fedoraTopParser.Parse: must be a number at %d index: %s", 27, row[27])
				return
			}
			used = total - available
			respMemory.Total = total
			respMemory.Used = used
			respMemory.Percent = calcPercent(total, used)
		}
	}
	respCpu.Percent = calcPercent(totalSum, usedSum)
	return
}
