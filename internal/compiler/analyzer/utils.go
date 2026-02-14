package analyzer

import (
	"fmt"

	src "github.com/nevalang/neva/pkg/ast"
)

type netNodesUsage map[string]netNodeUsage

func (n netNodesUsage) trackInportUsage(addr src.PortAddr) error {
	if _, ok := n[addr.Node]; !ok {
		n[addr.Node] = netNodeUsage{
			In:  portsUsage{},
			Out: portsUsage{},
		}
	}
	return n[addr.Node].trackInportUsage(addr)
}

func (n netNodesUsage) trackOutportUsage(addr src.PortAddr) error {
	if _, ok := n[addr.Node]; !ok {
		n[addr.Node] = netNodeUsage{
			In:  portsUsage{},
			Out: portsUsage{},
		}
	}
	return n[addr.Node].trackOutportUsage(addr)
}

type netNodeUsage struct {
	In, Out portsUsage
}

func (n netNodeUsage) trackOutportUsage(addr src.PortAddr) error {
	return n.Out.trackSlotUsage(addr)
}

func (n netNodeUsage) trackInportUsage(addr src.PortAddr) error {
	return n.In.trackSlotUsage(addr)
}

// portsUsage maps port name to slots used, slots map is nil for single ports
type portsUsage map[string]map[uint8]struct{}

func (p portsUsage) trackSlotUsage(addr src.PortAddr) error {
	isArrayBypass := src.IsArrayBypassIdx(addr.Idx)
	if _, ok := p[addr.Port]; !ok {
		if addr.Idx != nil && !isArrayBypass {
			p[addr.Port] = map[uint8]struct{}{
				*addr.Idx: {},
			}
		} else {
			p[addr.Port] = nil
		}
		return nil
	}

	if isArrayBypass {
		// Array-bypass consumes the whole port, so any additional usage is a duplicate.
		return fmt.Errorf("port '%v' is used twice", addr)
	}

	if addr.Idx == nil {
		return fmt.Errorf("port '%v' is used twice", addr)
	}

	if p[addr.Port] == nil {
		return fmt.Errorf("port '%v' mixes array-bypass with indexed usage", addr)
	}

	if _, ok := p[addr.Port][*addr.Idx]; ok {
		return fmt.Errorf("port '%v' is used twice", addr)
	}

	p[addr.Port][*addr.Idx] = struct{}{}

	return nil
}
