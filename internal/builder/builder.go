package builder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/parser"
)

type Builder struct {
	stdlibLocation     string        // path to standart library module
	thirdPartyLocation string        // path to third-party modules
	parser             parser.Parser // parser is needed to parse manifest files
}

func (b Builder) Build(ctx context.Context, workdir string) (compiler.Build, error) {
	// build user module
	mods := map[string]compiler.RawModule{}
	entryMod, err := b.buildModule(ctx, workdir, mods)
	if err != nil {
		return compiler.Build{}, fmt.Errorf("build entry mod: %w", err)
	}
	mods["entry"] = entryMod

	// build stdlib module
	stdMod, err := b.buildModule(ctx, b.stdlibLocation, mods)
	if err != nil {
		return compiler.Build{}, fmt.Errorf("build entry mod: %w", err)
	}
	mods["std"] = stdMod

	return compiler.Build{
		EntryModule: "entry",
		Modules:     mods,
	}, nil
}

func (b Builder) buildModule(
	ctx context.Context,
	workdir string,
	mods map[string]compiler.RawModule,
) (compiler.RawModule, error) {
	rawManifest, err := readManifestYaml(workdir)
	if err != nil {
		return compiler.RawModule{}, fmt.Errorf("read manifest yaml: %w", err)
	}

	manifest, err := b.parser.ParseManifest(rawManifest)
	if err != nil {
		return compiler.RawModule{}, fmt.Errorf("parse manifest: %w", err)
	}

	// process module deps (download if needed and build)
	for name, dep := range manifest.Deps {
		depPath := fmt.Sprintf("%s/%s_%s", b.thirdPartyLocation, dep.Addr, dep.Version)
		if _, err := os.Stat(depPath); err != nil { // check if directory with this dependency exists
			if os.IsNotExist(err) { // if not, clone the repo
				if _, err := git.PlainClone(depPath, false, &git.CloneOptions{
					URL:           fmt.Sprintf("https://%s", dep.Addr),
					ReferenceName: plumbing.NewTagReferenceName(dep.Version),
				}); err != nil {
					return compiler.RawModule{}, err
				}
			} else { // if it's an unknown error then return
				return compiler.RawModule{}, fmt.Errorf("os stat: %w", err)
			}
		}

		rawMod, err := b.buildModule(ctx, depPath, mods)
		if err != nil {
			return compiler.RawModule{}, fmt.Errorf("build module: %v: %w", name, err)
		}

		mods[name] = rawMod
	}

	pkgs := map[string]compiler.RawPackage{}

	if err := walk(workdir, pkgs); err != nil {
		return compiler.RawModule{}, fmt.Errorf("walk: %w", err)
	}

	return compiler.RawModule{
		Manifest: manifest,
		Packages: pkgs,
	}, nil
}

func readManifestYaml(workdir string) ([]byte, error) {
	rawManifest, err := os.ReadFile(workdir + "/neva.yml")
	if err == nil {
		return rawManifest, nil
	}

	rawManifest, err = os.ReadFile(workdir + "/neva.yaml")
	if err != nil {
		return nil, fmt.Errorf("os read file: %w", err)
	}

	return rawManifest, nil
}

// walk recursively traverses the directory assuming that is a neva module (set of packages with source code files).
func walk(rootPath string, prog map[string]compiler.RawPackage) error {
	if err := filepath.Walk(rootPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("filepath walk: %s: %w", filePath, err)
		}

		ext := filepath.Ext(info.Name())
		if ext != ".neva" {
			return nil
		}

		bb, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		pkgName := getPkgName(rootPath, filePath)
		if _, ok := prog[pkgName]; !ok {
			prog[pkgName] = compiler.RawPackage{}
		}

		fileName := strings.TrimSuffix(info.Name(), ext)
		prog[pkgName][fileName] = bb

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func getPkgName(rootPath, filePath string) string {
	dirPath := filepath.Dir(filePath)
	if dirPath == rootPath { // current directory is root directory
		ss := strings.Split(dirPath, "/")
		return ss[len(ss)-1] // return last element
	}
	return strings.TrimPrefix(dirPath, rootPath+"/")
}

func MustNew(stdPkgPath, thirdPartyLocation string, parser parser.Parser) Builder {
	return Builder{
		stdlibLocation:     stdPkgPath,
		thirdPartyLocation: thirdPartyLocation,
		parser:             parser,
	}
}
