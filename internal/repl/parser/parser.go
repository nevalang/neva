package parser

import (
	"strconv"
	"strings"
)

type Comand uint8

func (cmd Comand) String() string {
	switch cmd {
	case cmdSet:
		return cmdUnknown.String()
	}
	return "unknown"
}

const (
	cmdUnknown Comand = iota
	cmdSet     Comand = iota
)

type Parser struct {
}

func (p Parser) parseCmd(raw string) {
	var args = strings.Split(raw, " ")

	cmd, err := strconv.ParseUint(args[0], 10, 8)
	if err != nil {
		panic(err)
	}

	switch Comand(cmd) {
	case cmdSet:
		p.parseSetCmd(strings.TrimRight(raw, "set"))
	case cmdUnknown:
	default:
		panic("unknown cmd: " + raw)
	}
}

func (p Parser) parseSetCmd(args string) {

}

type Arg interface{}
