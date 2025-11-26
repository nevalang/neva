package analyzer

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
	"github.com/nevalang/neva/pkg"
)

// semverCheck ensures that module is compatible with existing compiler
// by checking it's version against semver. It uses minor as major.
func (a Analyzer) semverCheck(mod src.Module, modRef core.ModuleRef) *compiler.Error {
	moduleVersion, semverErr := semver.NewVersion(mod.Manifest.LanguageVersion)
	if semverErr != nil {
		return &compiler.Error{
			Message: fmt.Sprintf(
				"Module %v has invalid language Version: %v",
				modRef, semverErr,
			),
		}
	}

	compilerVersion, semverErr := semver.NewVersion(pkg.Version)
	if semverErr != nil {
		return &compiler.Error{
			Message: fmt.Sprintf(
				"Compiler version is invalid: %v", semverErr,
			),
		}
	}

	// major versions should be strictly equal
	// if got major more than ours, then compatibility in that program is broken
	// and vice versa if got major less than ours
	if moduleVersion.Major() != compilerVersion.Major() {
		return &compiler.Error{
			Message: fmt.Sprintf(
				"Incompatible compiler versions: module %v wants %v while current is %v",
				modRef, mod.Manifest.LanguageVersion, pkg.Version,
			),
		}
	}

	// if majors are equal, then minor should be less or equal
	// so we make sure module don't want any features we don't have
	if moduleVersion.Minor() > compilerVersion.Minor() {
		return &compiler.Error{
			Message: fmt.Sprintf(
				"Incompatible compiler versions: module %v wants %v while current is %v",
				modRef, mod.Manifest.LanguageVersion, pkg.Version,
			),
		}
	}

	// at this point we sure we have same majors and got.Minor >= want.Minor

	// if module's minor is less than ours then we don't care about patches
	// because with newer minor we sure have all the patches of the previous one
	if moduleVersion.Minor() < compilerVersion.Minor() {
		// note that we don't fix previous minors and instead force to update
		return nil
	}

	// if we here then we sure than minors are equal (as well as majors)
	// this is the only case where we actually care about patches

	// it's ok if we have some patches that module doesn't rely on
	// but it's not ok if module a wants some patch we don't really have
	if moduleVersion.Patch() > compilerVersion.Patch() {
		return &compiler.Error{
			Message: fmt.Sprintf(
				"Incompatible compiler versions: module %v wants %v while current is %v",
				modRef, mod.Manifest.LanguageVersion, pkg.Version,
			),
		}
	}

	// versions are strictly equal
	return nil
}
