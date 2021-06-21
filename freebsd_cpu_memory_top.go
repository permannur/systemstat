package systemstat

import (
	"fmt"
	"strconv"
)

// FreeBsdTop Top version FreeBSD
func FreeBsdTop() (cpu *Cpu, memory *Memory, err error) {
	freeBsdTop := &freeBsdTopParser{
		cpu:    &Cpu{},
		memory: &Memory{},
	}
	err = newCommand(`(CPU\s+(\d+)?:\s+(\d+)(?:\.(\d+))?% user,\s+(\d+)(?:\.(\d+))?% nice,\s+(\d+)(?:\.(\d+))?% system,\s+(\d+)(?:\.(\d+))?% interrupt,\s+(\d+)(?:\.(\d+))?% idle\n)(?:(Mem:(?:\s+(\d+)([BKMGTE]) Active[,\n])?(?:\s+(\d+)([BKMGTE]) Inact[,\n])?(?:\s+(\d+)([BKMGTE]) Laundry[,\n])?(?:\s+(\d+)([BKMGTE]) Wired[,\n])?(?:\s+(\d+)([BKMGTE]) Buf[,\n])?(?:\s+(\d+)([BKMGTE]) Free[,\n])?)((.|\n)*(Swap:(?:\s+(\d+)([BKMGTE]) Total[,\n])?(?:\s+(\d+)([BKMGTE]) Free[,\n])?(?:\s+(\d+)(?:\.(\d+))?% Inuse[,\n])?(?:\s+(\d+)([BKMGTE]) In[,\n])?(?:\s+(\d+)([BKMGTE]) Out[,\n])?))?)?`,
		freeBsdTop,
		"top",
		"-b", "-P", "0")
	if err != nil {
		err = fmt.Errorf("FreeBsdTop: %s", err)
		return
	}
	cpu = freeBsdTop.cpu
	memory = freeBsdTop.memory
	return
}

type freeBsdTopParser struct {
	cpu    *Cpu
	memory *Memory
}

func (c *freeBsdTopParser) Parse(parsedList [][]string) (err error) {
	if c.cpu == nil {
		c.cpu = &Cpu{}
	}
	if c.memory == nil {
		c.memory = &Memory{Unit: "byte"}
	}
	respCpu, respMemory := c.cpu, c.memory
	var totalSum, usedSum int64
	var cpuDetailInx int
	for _, row := range parsedList {
		if len(row) != 39 {
			err = fmt.Errorf("freeBsdTopParser.Parse: parsed list length mismatch, want: %d, have: %d", 39, len(row))
			return
		}
		if len(row[1]) > 0 {
			var total, available, used, idle1, idle2 int64
			idle1, err = strconv.ParseInt(row[11], 10, 64)
			if err != nil {
				err = fmt.Errorf("freeBsdTopParser.Parse: must be a number at %d index: %s", 11, row[11])
				return
			}
			if len(row[12]) > 0 {
				idle2, err = strconv.ParseInt(row[12], 10, 64)
				if err != nil {
					err = fmt.Errorf("freeBsdTopParser.Parse: must be a number at %d index: %s", 12, row[12])
					return
				}
			}
			total = 1000
			available = idle1*10 + idle2
			used = total - available
			if len(respCpu.DetailList) == cpuDetailInx {
				respCpu.DetailList = append(respCpu.DetailList, &CpuDetail{})
			}
			respCpuDetail := respCpu.DetailList[cpuDetailInx]
			respCpuDetail.Name = "CPU " + row[2]
			respCpuDetail.Percent = calcPercent(total, used)
			totalSum += total
			usedSum += used
			cpuDetailInx++
		}
		if len(row[13]) > 0 {
			var active, inactive, laundry, wired, buf, free int64
			if len(row[14]) > 0 {
				active, err = strconv.ParseInt(row[14], 10, 64)
				if err != nil {
					err = fmt.Errorf("freeBsdTopParser.Parse: must be a number at %d index: %s", 14, row[14])
					return
				}
				active *= toByte(row[15])
			}
			if len(row[16]) > 0 {
				inactive, err = strconv.ParseInt(row[16], 10, 64)
				if err != nil {
					err = fmt.Errorf("freeBsdTopParser.Parse: must be a number at %d index: %s", 16, row[16])
					return
				}
				inactive *= toByte(row[17])
			}
			if len(row[18]) > 0 {
				laundry, err = strconv.ParseInt(row[18], 10, 64)
				if err != nil {
					err = fmt.Errorf("freeBsdTopParser.Parse: must be a number at %d index: %s", 18, row[18])
					return
				}
				laundry *= toByte(row[19])
			}
			if len(row[20]) > 0 {
				wired, err = strconv.ParseInt(row[20], 10, 64)
				if err != nil {
					err = fmt.Errorf("freeBsdTopParser.Parse: must be a number at %d index: %s", 20, row[20])
					return
				}
				wired *= toByte(row[21])
			}
			if len(row[22]) > 0 {
				buf, err = strconv.ParseInt(row[22], 10, 64)
				if err != nil {
					err = fmt.Errorf("freeBsdTopParser.Parse: must be a number at %d index: %s", 22, row[22])
					return
				}
				buf *= toByte(row[23])
			}
			if len(row[24]) > 0 {
				free, err = strconv.ParseInt(row[24], 10, 64)
				if err != nil {
					err = fmt.Errorf("freeBsdTopParser.Parse: must be a number at %d index: %s", 24, row[24])
					return
				}
				free *= toByte(row[25])
			}
			used := active + inactive + laundry + wired + buf
			respMemory.Total = used + free
			respMemory.Used = used
			respMemory.Percent = calcPercent(respMemory.Total, respMemory.Used)
		}
	}
	respCpu.Percent = calcPercent(totalSum, usedSum)
	return
}

func toByte(unit string) int64 {
	var k int64
	k = 1
	if unit == "B" {
		return k
	}
	k *= 1024
	if unit == "K" {
		return k
	}
	k *= 1024
	if unit == "M" {
		return k
	}
	k *= 1024
	if unit == "G" {
		return k
	}
	k *= 1024
	if unit == "T" {
		return k
	}
	k *= 1024
	if unit == "E" {
		return k
	}
	return 0
}
