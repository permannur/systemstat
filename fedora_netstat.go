package systemstat

import (
	"fmt"
	"strconv"
	"strings"
)

// FedoraNetstat netstat -e -i
func FedoraNetstat() (receive, send *Network, err error) {
	fedoraNetstat := &fedoraNetstatParser{
		receive: &Network{},
		send:    &Network{},
	}
	err = newCommand(`(\w+):(.|\n)+?bytes\s+(\d+)(.|\n)+?bytes\s+(\d+)`, fedoraNetstat, "netstat", "-e", "-i")
	if err != nil {
		err = fmt.Errorf("FedoraNetstat: %s", err)
		return
	}
	receive = fedoraNetstat.receive
	send = fedoraNetstat.send
	return
}

type fedoraNetstatParser struct {
	receive, send *Network
}

func (c *fedoraNetstatParser) Parse(parsedList [][]string) (err error) {
	if c.receive == nil {
		c.receive = &Network{Unit: "bytes"}
	}
	if c.send == nil {
		c.send = &Network{Unit: "bytes"}
	}
	respReceive, respSend := c.receive, c.send
	var receivedSum, sentSum int64
	for inx, row := range parsedList {
		if len(row) != 6 {
			err = fmt.Errorf("fedoraNetstatParser.Parse: parsed list length mismatch, want: %d, have: %d", 6, len(row))
			return
		}
		var received, sent int64
		received, err = strconv.ParseInt(row[3], 10, 64)
		if err != nil {
			err = fmt.Errorf("fedoraNetstatParser.Parse: must be a number at %d index: %s", 3, row[3])
			return
		}
		sent, err = strconv.ParseInt(row[5], 10, 64)
		if err != nil {
			err = fmt.Errorf("fedoraNetstatParser.Parse: must be a number at %d index: %s", 5, row[5])
			return
		}
		if len(respReceive.DetailList) == inx {
			respReceive.DetailList = append(respReceive.DetailList, &NetworkDetail{})
		}
		if len(respSend.DetailList) == inx {
			respSend.DetailList = append(respSend.DetailList, &NetworkDetail{})
		}
		respReceiveDetail, respSendDetail := respReceive.DetailList[inx], respSend.DetailList[inx]
		respReceiveDetail.Name = row[1]
		respReceiveDetail.PerSecond = received - respReceiveDetail.Total
		respReceiveDetail.Total = received
		respSendDetail.Name = row[1]
		respSendDetail.PerSecond = sent - respSendDetail.Total
		respSendDetail.Total = sent
		if !strings.HasPrefix(row[1], "lo") {
			receivedSum += received
			sentSum += sent
		}
	}
	respReceive.PerSecond = receivedSum - respReceive.Total
	respReceive.Total = receivedSum
	respSend.PerSecond = sentSum - respSend.Total
	respSend.Total = sentSum
	return
}
