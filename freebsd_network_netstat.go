package systemstat

import (
	"fmt"
	"strconv"
	"strings"
)

// FreeBsdNetstat netstat -ib
func FreeBsdNetstat() (in, out *Network, err error) {
	freeBsdNetstat := &freeBsdNetstatParser{
		in:  &Network{},
		out: &Network{},
	}
	err = newCommand(`(\S+)\s+(?:(\d+)|(?:-))\s+(\S+)\s+(\S+)\s+(?:(\d+)|(?:-))\s+(?:(\d+)|(?:-))\s+(?:(\d+)|(?:-))\s+(?:(\d+)|(?:-))\s+(?:(\d+)|(?:-))\s+(?:(\d+)|(?:-))\s+(?:(\d+)|(?:-))\s+(?:(\d+)|(?:-))\n`,
		freeBsdNetstat,
		"netstat",
		"-ib")
	if err != nil {
		err = fmt.Errorf("FreeBsdNetstat: %s", err)
		return
	}
	in = freeBsdNetstat.in
	out = freeBsdNetstat.out
	return
}

type freeBsdNetstatParser struct {
	in  *Network
	out *Network
}

func (c *freeBsdNetstatParser) Parse(parsedList [][]string) (err error) {
	if c.in == nil {
		c.in = &Network{Unit: "bytes"}
	}
	if c.out == nil {
		c.out = &Network{Unit: "bytes"}
	}
	respIn, respOut := c.in, c.out
	var inSum, outSum int64
	for inx, row := range parsedList {
		if len(row) != 13 {
			err = fmt.Errorf("freeBsdNetstatParser.Parse: parsed list length mismatch, want: %d, have: %d", 13, len(row))
			return
		}
		var inTotal, outTotal int64
		if row[8] != "" {
			inTotal, err = strconv.ParseInt(row[8], 10, 64)
			if err != nil {
				err = fmt.Errorf("freeBsdNetstatParser.Parse: must be a number at %d index: %s", 8, row[8])
				return
			}
		}
		if row[11] != "" {
			outTotal, err = strconv.ParseInt(row[11], 10, 64)
			if err != nil {
				err = fmt.Errorf("freeBsdNetstatParser.Parse: must be a number at %d index: %s", 11, row[11])
				return
			}
		}
		if len(respIn.DetailList) == inx {
			respIn.DetailList = append(respIn.DetailList, &NetworkDetail{})
		}
		if len(respOut.DetailList) == inx {
			respOut.DetailList = append(respOut.DetailList, &NetworkDetail{})
		}
		respInDetail, respOutDetail := respIn.DetailList[inx], respOut.DetailList[inx]
		respInDetail.Name = row[1]
		respInDetail.PerSecond = inTotal - respInDetail.Total
		respInDetail.Total = inTotal
		respOutDetail.Name = row[1]
		respOutDetail.PerSecond = outTotal - respOutDetail.Total
		respOutDetail.Total = outTotal
		if !strings.HasPrefix(row[1], "lo") {
			inSum += inTotal
			outSum += outTotal
		}
	}
	respIn.PerSecond = inSum - respIn.Total
	respIn.Total = inSum
	respOut.PerSecond = outSum - respOut.Total
	respOut.Total = outSum
	return
}
