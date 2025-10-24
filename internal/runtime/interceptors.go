package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type ProdInterceptor struct{}

func (ProdInterceptor) Prepare() error { return nil }

func (ProdInterceptor) Sent(sender PortSlotAddr, msg Msg) Msg { return msg }

func (ProdInterceptor) Received(receiver PortSlotAddr, msg Msg) Msg { return msg }

type DebugInterceptor struct{ file *os.File }

func (d *DebugInterceptor) Open(filepath string, metadata TraceMetadata) (func() error, error) {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	d.file = file
	if err := d.writeMetadata(metadata); err != nil {
		file.Close()
		return nil, err
	}
	return file.Close, nil
}

type TraceMetadata struct {
	Program  TraceProgramMetadata  `json:"program,omitempty"`
	Compiler TraceCompilerMetadata `json:"compiler,omitempty"`
}

type TraceProgramMetadata struct {
	ModulePath    string `json:"modulePath,omitempty"`
	ModuleVersion string `json:"moduleVersion,omitempty"`
	MainPackage   string `json:"mainPackage,omitempty"`
}

type TraceCompilerMetadata struct {
	Version string `json:"version,omitempty"`
}

func (m TraceMetadata) isZero() bool {
	return m.Program == (TraceProgramMetadata{}) && m.Compiler == (TraceCompilerMetadata{})
}

func (d *DebugInterceptor) writeMetadata(metadata TraceMetadata) error {
	if metadata.isZero() {
		return nil
	}

	payload, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(d.file, "# metadata %s\n", payload); err != nil {
		return err
	}

	return nil
}

func (d *DebugInterceptor) Sent(sender PortSlotAddr, msg Msg) Msg {
	fmt.Fprintf(
		d.file,
		"sent | %v | %v\n",
		d.formatPortSlotAddr(sender), d.formatMsg(msg),
	)
	return msg
}

func (d *DebugInterceptor) Received(receiver PortSlotAddr, msg Msg) Msg {
	fmt.Fprintf(
		d.file,
		"recv | %v | %v\n",
		d.formatPortSlotAddr(receiver),
		d.formatMsg(msg),
	)
	return msg
}

func (d DebugInterceptor) formatMsg(msg Msg) string {
	if strMsg, ok := msg.(StringMsg); ok {
		return fmt.Sprintf("%q", strMsg.Str())
	}
	return fmt.Sprint(msg)
}

func (d DebugInterceptor) formatPortSlotAddr(slotAddr PortSlotAddr) string {
	parts := strings.Split(slotAddr.Path, "/")
	lastPart := parts[len(parts)-1]
	if lastPart == "in" || lastPart == "out" {
		parts = parts[:len(parts)-1]
	}
	slotAddr.Path = strings.Join(parts, "/")

	s := fmt.Sprintf("%v:%v", slotAddr.Path, slotAddr.Port)
	if slotAddr.Index != nil {
		s = fmt.Sprintf("%v[%v]", s, *slotAddr.Index)
	}

	return s
}

func NewDebugInterceptor() *DebugInterceptor {
	return &DebugInterceptor{}
}
