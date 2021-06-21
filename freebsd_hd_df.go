package parser

import (
	"fmt"
	"strconv"
	"github.com/permannur/command"
)

// FreeBsdDf df -kT
func FreeBsdDf(diskList ...string) (hd *Hd, err error) {
	diskListInRegex := ""
	for inx, val := range diskList {
		if inx != 0 {
			diskListInRegex += "|"
		}
		diskListInRegex += val
	}
	freeBsdDf := freeBsdDfParser{
		hd: &Hd{},
	}
	err = command.NewCommand(`(\S+)\s+(\w+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)(?:\.(\d+))?%\s+(`+diskListInRegex+`)\n`,
		freeBsdDf,
		"df",
		"-kT")
	if err != nil {
		err = fmt.Errorf("FreeBsdDf: %s", err)
		return
	}
	hd = freeBsdDf.hd
	return
}

type freeBsdDfParser struct {
	hd *Hd
}

func (c freeBsdDfParser) Parse(parsedList [][]string) (err error) {
	if c.hd == nil {
		c.hd = &Hd{Unit: "KiB"}
	}
	resp := c.hd
	var totalSum, usedSum int64
	for inx, row := range parsedList {
		if len(row) != 9 {
			err = fmt.Errorf("freeBsdDfParser.Parse: parsed list length mismatch, want: %d, have: %d", 9, len(row))
			return
		}
		var total, used int64
		total, err = strconv.ParseInt(row[3], 10, 64)
		if err != nil {
			err = fmt.Errorf("freeBsdDfParser.Parse: must be a number at %d index: %s", 3, row[3])
			return
		}
		used, err = strconv.ParseInt(row[4], 10, 64)
		if err != nil {
			err = fmt.Errorf("freeBsdDfParser.Parse: must be a number at %d index: %s", 4, row[4])
			return
		}
		if len(resp.DetailList) == inx {
			resp.DetailList = append(resp.DetailList, &HdDetail{})
		}
		respDetail := resp.DetailList[inx]
		respDetail.Name = row[8]
		respDetail.Total = total
		respDetail.Used = used
		respDetail.Percent = calcPercent(total, used)
		totalSum += total
		usedSum += used
	}
	resp.Total = totalSum
	resp.Used = usedSum
	resp.Percent = calcPercent(totalSum, usedSum)
	return
}
