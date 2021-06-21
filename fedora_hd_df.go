package parser

import (
	"fmt"
	"strconv"
)

// FedoraDf df -B KiB
func FedoraDf(diskList ...string) (hd *Hd, err error) {
	diskListInRegex := ""
	for inx, val := range diskList {
		if inx != 0 {
			diskListInRegex += "|"
		}
		diskListInRegex += val
	}
	fedoraDf := &fedoraDfParser{
		hd: &Hd{},
	}
	err = newCommand(`(\d+)KiB\s+(\d+)KiB\s+(\d+)KiB\s+(\d+)%\s+(`+diskListInRegex+`)\n`, fedoraDf, "df", "-B", "KiB")
	if err != nil {
		err = fmt.Errorf("FedoraDf: %s", err)
		return
	}
	hd = fedoraDf.hd
	return
}

type fedoraDfParser struct {
	hd *Hd
}

func (c *fedoraDfParser) Parse(parsedList [][]string) (err error) {
	if c.hd == nil {
		c.hd = &Hd{Unit: "KiB"}
	}
	resp := c.hd
	var totalSum, usedSum int64
	for inx, row := range parsedList {
		if len(row) != 6 {
			err = fmt.Errorf("fedoraDfParser.Parse: parsed list length mismatch, want: %d, have: %d", 6, len(row))
			return
		}
		var used, total int64
		total, err = strconv.ParseInt(row[1], 10, 64)
		if err != nil {
			err = fmt.Errorf("fedoraDfParser.Parse: must be a number at %d index: %s", 1, row[1])
			return
		}
		used, err = strconv.ParseInt(row[2], 10, 64)
		if err != nil {
			err = fmt.Errorf("fedoraDfParser.Parse: must be a number at %d index: %s", 2, row[2])
			return
		}
		if len(resp.DetailList) == inx {
			resp.DetailList = append(resp.DetailList, &HdDetail{})
		}
		respDetail := resp.DetailList[inx]
		respDetail.Name = row[5]
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
