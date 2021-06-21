package parser

import (
	"fmt"
	"time"
)

type updater struct {
	interval   time.Duration
	quit       chan interface{}
	readerList []reader
}

type reader interface {
	Read() error
}

var u *updater

func (u *updater) Stop() {
	u.quit <- 1
}

func addReader(rd reader) {
	if u == nil {
		u = &updater{
			interval: time.Second,
			quit:     make(chan interface{}),
		}
		go u.tickFn(time.Tick(u.interval))
	}
	u.readerList = append(u.readerList, rd)
}

func (u *updater) tickFn(tick <-chan time.Time) {
	var err error
	for {
		select {
		case <-tick:
			for _, update := range u.readerList {
				err = update.Read()
				if err != nil {
					err = fmt.Errorf("updater.tickFn: reading error, %s", err)
					return
				}
			}
		case <-u.quit:
			return
		}
	}
}
