package view

import (
	"fmt"
	"strings"

	"github.com/nevalang/neva/pkg/core"
)

func moduleID(modRef core.ModuleRef) string {
	if modRef.Version == "" {
		return "module/" + sanitizeSegment(modRef.Path)
	}
	return "module/" + sanitizeSegment(modRef.Path) + "@" + sanitizeSegment(modRef.Version)
}

func packageID(modRef core.ModuleRef, packageName string) string {
	return moduleID(modRef) + "/package/" + sanitizeSegment(packageName)
}

func fileID(loc core.Location) string {
	return packageID(loc.ModRef, loc.Package) + "/file/" + sanitizeSegment(loc.Filename)
}

func constID(loc core.Location, constName string) string {
	return fileID(loc) + "/const/" + sanitizeSegment(constName)
}

func typeID(loc core.Location, typeName string) string {
	return fileID(loc) + "/type/" + sanitizeSegment(typeName)
}

func interfaceID(loc core.Location, name string) string {
	return fileID(loc) + "/interface/" + sanitizeSegment(name)
}

func componentID(loc core.Location, name string, overloadIndex int) string {
	return fmt.Sprintf("%s/component/%s@%d", fileID(loc), sanitizeSegment(name), overloadIndex)
}

func nodeID(componentViewID string, nodeName string) string {
	return componentViewID + "/node/" + sanitizeSegment(nodeName)
}

func diNodeID(componentViewID string, nodeName string, diName string) string {
	return nodeID(componentViewID, nodeName) + "/di/" + sanitizeSegment(diName)
}

func portID(parentID string, direction string, portName string) string {
	return parentID + "/" + sanitizeSegment(direction) + "/port/" + sanitizeSegment(portName)
}

func importID(fileViewID string, alias string, module string, pkg string) string {
	return fileViewID + "/import/" + sanitizeSegment(alias) + "=" + sanitizeSegment(module) + ":" + sanitizeSegment(pkg)
}

func sanitizeSegment(value string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		" ", "_",
		"\t", "_",
		"\n", "_",
		":", "_",
		"#", "_",
	)
	clean := strings.TrimSpace(replacer.Replace(value))
	if clean == "" {
		return "_"
	}
	return clean
}
