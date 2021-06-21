package command

import (
	"fmt"
	"os/exec"
	"regexp"
	"github.com/permannur/updater"
)

type parser interface {
	Parse([][]string) error
}

type command struct {
	cmd    string
	args   []string
	regexC *regexp.Regexp
	parser parser
}

func NewCommand(regexStr string, parser parser, cmdStr string, args ...string) (err error) {
	c := &command{
		cmd:    cmdStr,
		args:   args,
		parser: parser,
	}
	c.regexC, err = regexp.Compile(regexStr)
	if err != nil {
		err = fmt.Errorf("NewCommand: error regexp.Compile, %s", err)
		return
	}
	_, err = exec.LookPath(c.cmd)
	if err != nil {
		err = fmt.Errorf("NewCommand: command not found")
		return
	}
	err = c.Read()
	if err != nil {
		return
	}
	updater.AddReader(c)
	return
}

func (c *command) Read() (err error) {
	var bt []byte
	bt, _ = exec.Command(c.cmd, c.args...).Output()
	parsedList := c.regexC.FindAllStringSubmatch(string(bt), -1)
	if len(parsedList) == 0 {
		err = fmt.Errorf("command.Read: command and regex mismatch")
		return
	}
	err = c.parser.Parse(parsedList)
	if err != nil {
		err = fmt.Errorf("command.Read: parser error, %s", err)
		return
	}
	return
}
