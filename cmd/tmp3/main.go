// generator/main.go
package main

import (
	"bytes"
	"io/fs"
	"os"

	"github.com/emil14/neva/internal"
	"github.com/emil14/neva/internal/compiler/backend/golang"
	"github.com/emil14/neva/internal/compiler/ir"
)

var efs = internal.RuntimeFiles
var basePath = "/home/evaleev/projects/tmp"

func main() {
	if err := os.RemoveAll(basePath); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		panic(err)
	}

	putGoMod()

	putRuntime()

	bb, err := (golang.Backend{}).GenerateTarget(nil, ir.Program{
		Ports: map[ir.PortAddr]uint8{
			{Port: "start"}:                            0,
			{Port: "exit"}:                             0,
			{Path: "printer_in", Port: "v", Idx: 0}:    0,
			{Path: "printer_out", Port: "v", Idx: 0}:   0,
			{Path: "trigger_in", Port: "sigs", Idx: 0}: 0,
			{Path: "trigger_in", Port: "v", Idx: 0}:    0,
			{Path: "trigger_out", Port: "v", Idx: 0}:   0,
			{Path: "giver_out", Port: "code", Idx: 0}:  0,
		},
		Routines: ir.Routines{
			Giver: map[ir.PortAddr]ir.Msg{
				{
					Path: "giver_out",
					Port: "code",
					Idx:  0,
				}: {
					Type: ir.IntMsg,
					Int:  0,
				},
			},
			Component: []ir.ComponentRef{
				{
					Pkg:  "io",
					Name: "Print",
					PortAddrs: ir.ComponentPortAddrs{
						In: []ir.PortAddr{
							{Path: "printer_in", Port: "v", Idx: 0},
						},
						Out: []ir.PortAddr{
							{Path: "printer_out", Port: "v", Idx: 0},
						},
					},
				},
				{
					Pkg:  "flow",
					Name: "Trigger",
					PortAddrs: ir.ComponentPortAddrs{
						In: []ir.PortAddr{
							{Path: "trigger_in", Port: "sigs", Idx: 0},
							{Path: "trigger_in", Port: "v", Idx: 0},
						},
						Out: []ir.PortAddr{
							{Path: "trigger_out", Port: "v", Idx: 0},
						},
					},
				},
			},
		},
		Connections: []ir.Connection{
			{
				SenderSide: ir.ConnectionSide{
					PortAddr: ir.PortAddr{
						Path: "",
						Port: "start",
						Idx:  0,
					},
					Selectors: []ir.Selector{},
				},
				ReceiverSides: []ir.ConnectionSide{
					{
						PortAddr: ir.PortAddr{
							Path: "printer_in",
							Port: "v",
							Idx:  0,
						},
						Selectors: []ir.Selector{},
					},
				},
			},
			{
				SenderSide: ir.ConnectionSide{
					PortAddr: ir.PortAddr{
						Path: "printer_out",
						Port: "v",
						Idx:  0,
					},
					Selectors: []ir.Selector{},
				},
				ReceiverSides: []ir.ConnectionSide{
					{
						PortAddr: ir.PortAddr{
							Path: "trigger_in",
							Port: "sigs",
							Idx:  0,
						},
						Selectors: []ir.Selector{},
					},
				},
			},
			{
				SenderSide: ir.ConnectionSide{
					PortAddr: ir.PortAddr{
						Path: "giver_out",
						Port: "code",
						Idx:  0,
					},
					Selectors: []ir.Selector{},
				},
				ReceiverSides: []ir.ConnectionSide{
					{
						PortAddr: ir.PortAddr{
							Path: "trigger_in",
							Port: "v",
							Idx:  0,
						},
						Selectors: []ir.Selector{},
					},
				},
			},
			{
				SenderSide: ir.ConnectionSide{
					PortAddr: ir.PortAddr{
						Path: "trigger_out",
						Port: "v",
						Idx:  0,
					},
					Selectors: []ir.Selector{},
				},
				ReceiverSides: []ir.ConnectionSide{
					{
						PortAddr: ir.PortAddr{
							Path: "",
							Port: "exit",
							Idx:  0,
						},
						Selectors: []ir.Selector{},
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// write main.go
	var buf bytes.Buffer
	if _, err := buf.Write(bb); err != nil {
		panic(err)
	}
	if err := os.WriteFile(basePath+"/"+"main.go", buf.Bytes(), os.ModePerm); err != nil {
		panic(err)
	}
}

func putRuntime() {
	// prepare directory structure and collect files to create
	files := map[string][]byte{}
	if err := fs.WalkDir(efs, "runtime", func(path string, d fs.DirEntry, err error) error {
		fullPath := basePath + "/internal/" + path
		if d.IsDir() {
			if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
				return err
			}
			return nil
		}

		bb, err := efs.ReadFile(path)
		if err != nil {
			return err
		}

		files[fullPath] = bb
		return nil
	}); err != nil {
		panic(err)
	}
	// create files
	for path, bb := range files {
		if err := os.WriteFile(path, bb, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func putGoMod() {
	f, err := os.Create(basePath + "/go.mod")
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = f.WriteString("module github.com/emil14/neva")
	if err != nil {
		panic(err)
	}
}
