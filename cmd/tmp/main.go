package main

import (
	"bytes"
	"io/fs"
	"os"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/backend/golang"
)

var efs = golang.Efs
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

	bb, err := (golang.Backend{}).GenerateTarget(nil, compiler.LowLvlProgram{
		Ports: map[compiler.LLPortAddr]uint8{
			{Name: "start"}:                            0,
			{Name: "exit"}:                             0,
			{Path: "printer.in", Name: "v", Idx: 0}:    0,
			{Path: "printer.out", Name: "v", Idx: 0}:   0,
			{Path: "trigger.in", Name: "sigs", Idx: 0}: 0,
			{Path: "trigger.in", Name: "v", Idx: 0}:    0,
			{Path: "trigger.out", Name: "v", Idx: 0}:   0,
			{Path: "giver.out", Name: "code", Idx: 0}:  0,
		},
		// Routines: compiler.LLRoutines{
		// 	Giver: map[compiler.LLPortAddr]compiler.LLMsg{
		// 		{
		// 			Path: "giver.out",
		// 			Name: "code",
		// 			Idx:  0,
		// 		}: {
		// 			Type: compiler.LLIntMsg,
		// 			Int:  0,
		// 		},
		// 	},
		// 	Func: []compiler.LLFunc{
		// 		{
		// 			Ref: compiler.LLFuncRef{
		// 				Pkg:  "io",
		// 				Name: "Print",
		// 			},
		// 			IO: compiler.LLFuncIO{
		// 				In: []compiler.LLPortAddr{
		// 					{Path: "printer.in", Name: "v", Idx: 0},
		// 				},
		// 				Out: []compiler.LLPortAddr{
		// 					{Path: "printer.out", Name: "v", Idx: 0},
		// 				},
		// 			},
		// 		},
		// 		{
		// 			Ref: compiler.LLFuncRef{
		// 				Pkg:  "flow",
		// 				Name: "Trigger",
		// 			},
		// 			IO: compiler.LLFuncIO{
		// 				In: []compiler.LLPortAddr{
		// 					{Path: "trigger.in", Name: "sigs", Idx: 0},
		// 					{Path: "trigger.in", Name: "v", Idx: 0},
		// 				},
		// 				Out: []compiler.LLPortAddr{
		// 					{Path: "trigger.out", Name: "v", Idx: 0},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		Net: []compiler.LLConnection{
			{
				SenderSide: compiler.LLConnectionSide{
					PortAddr: compiler.LLPortAddr{
						Path: "",
						Name: "start",
						Idx:  0,
					},
					Selectors: []compiler.LLSelector{},
				},
				ReceiverSides: []compiler.LLConnectionSide{
					{
						PortAddr: compiler.LLPortAddr{
							Path: "printer.in",
							Name: "v",
							Idx:  0,
						},
						Selectors: []compiler.LLSelector{},
					},
				},
			},
			{
				SenderSide: compiler.LLConnectionSide{
					PortAddr: compiler.LLPortAddr{
						Path: "printer.out",
						Name: "v",
						Idx:  0,
					},
					Selectors: []compiler.LLSelector{},
				},
				ReceiverSides: []compiler.LLConnectionSide{
					{
						PortAddr: compiler.LLPortAddr{
							Path: "trigger.in",
							Name: "sigs",
							Idx:  0,
						},
						Selectors: []compiler.LLSelector{},
					},
				},
			},
			{
				SenderSide: compiler.LLConnectionSide{
					PortAddr: compiler.LLPortAddr{
						Path: "giver.out",
						Name: "code",
						Idx:  0,
					},
					Selectors: []compiler.LLSelector{},
				},
				ReceiverSides: []compiler.LLConnectionSide{
					{
						PortAddr: compiler.LLPortAddr{
							Path: "trigger.in",
							Name: "v",
							Idx:  0,
						},
						Selectors: []compiler.LLSelector{},
					},
				},
			},
			{
				SenderSide: compiler.LLConnectionSide{
					PortAddr: compiler.LLPortAddr{
						Path: "trigger.out",
						Name: "v",
						Idx:  0,
					},
					Selectors: []compiler.LLSelector{},
				},
				ReceiverSides: []compiler.LLConnectionSide{
					{
						PortAddr: compiler.LLPortAddr{
							Path: "",
							Name: "exit",
							Idx:  0,
						},
						Selectors: []compiler.LLSelector{},
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
		fullPath := basePath + "/internal/compiler/backend/golang/" + path
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

	_, err = f.WriteString("module github.com/nevalang/neva")
	if err != nil {
		panic(err)
	}
}
