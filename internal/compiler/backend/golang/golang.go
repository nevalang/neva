package golang

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"unicode"

	"github.com/nevalang/neva/internal/compiler"
)

//go:embed tmpl/main.go.tmpl
var efs embed.FS

type Backend struct{}

var errExecTmpl = errors.New("execute template")
var ErrWrongGoVersion = errors.New("wrong Go version")

const wantGoVersion = "1.21"

func (b Backend) GenerateTarget(ctx context.Context, prog compiler.LowLvlProgram, basePath string) error {
	tmpl, err := template.New("main.go.tmpl").Funcs(template.FuncMap{
		"getMsg":           getMsg,
		"getPorts":         getPortsFunc(prog.Ports),
		"getPortChVarName": getPortChVarName,
		"getConnComment":   getConnComment,
	}).ParseFS(efs, "tmpl/main.go.tmpl")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, prog); err != nil {
		return errors.Join(errExecTmpl, err)
	}

	if err := os.RemoveAll(basePath); err != nil {
		return err
	}

	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return err
	}

	// write main.go
	if _, err := buf.Write(nil); err != nil {
		return err
	}
	if err := os.WriteFile(basePath+"/"+"main.go", buf.Bytes(), os.ModePerm); err != nil {
		return err
	}

	// go.mod
	if err := putGoMod(basePath); err != nil {
		return err
	}

	// runtime
	if err := putRuntime(basePath); err != nil {
		return err
	}

	// check go compiler version
	versionStr, err := goVersionOut()
	if err != nil {
		return err
	}
	if !strings.Contains(versionStr, "1.21") {
		return fmt.Errorf("%w: want %s, got %s", ErrWrongGoVersion, wantGoVersion, versionStr)
	}

	// compile go code
	cmd := exec.Command("go", "build", basePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

var errUnknownMsgType = errors.New("unknown msg type")

func getMsg(msg compiler.LLMsg) (string, error) {
	switch msg.Type {
	case compiler.LLIntMsg:
		return fmt.Sprintf("runtime.NewIntMsg(%d)", msg.Int), nil
	}
	return "", fmt.Errorf("%w: %v", errUnknownMsgType, msg.Type)
}

func getConnComment(conn compiler.LLConnection) string {
	s := fmtPortAddr(conn.SenderSide.PortAddr) + " -> "

	for _, rcvr := range conn.ReceiverSides {
		s += fmtPortAddr(rcvr.PortAddr)
	}

	return "// " + s
}

func fmtPortAddr(addr compiler.LLPortAddr) string {
	return fmt.Sprintf("%s.%s[%d]", addr.Path, addr.Name, addr.Idx)
}

func getPortChVarName(addr compiler.LLPortAddr) string {
	path := handleSpecialChars(addr.Path)
	port := addr.Name
	if path != "" {
		port = uppercaseFirstLetter(addr.Name)
	}
	return fmt.Sprintf("%s%s%dPort", path, port, addr.Idx)
}

func getPortsFunc(ports map[compiler.LLPortAddr]uint8) func(path, port string) string {
	return func(path, port string) string {
		var s string
		for addr := range ports {
			if addr.Path == path && addr.Name == port {
				s = s + getPortChVarName(addr) + ","
			}
		}
		return s
	}
}

func handleSpecialChars(portPath string) string {
	var (
		buffer          bytes.Buffer
		shouldUppercase bool
	)

	for i := 0; i < len(portPath); i++ {
		if portPath[i] == '.' || portPath[i] == '/' {
			shouldUppercase = true
			continue
		}
		s := string(portPath[i])
		if shouldUppercase {
			s = strings.ToUpper(s)
			shouldUppercase = false
		}
		buffer.WriteString(s)
	}

	return buffer.String()
}

func uppercaseFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	bb := []byte(s)
	bb[0] = byte(unicode.ToUpper(rune(bb[0])))
	return string(bb)
}

func putGoMod(basePath string) error {
	f, err := os.Create(basePath + "/go.mod")
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString("module github.com/nevalang/neva")
	if err != nil {
		return err
	}

	return nil
}

func putRuntime(basePath string) error {
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
		return err
	}
	// create files
	for path, bb := range files {
		if err := os.WriteFile(path, bb, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func goVersionOut() (string, error) {
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
