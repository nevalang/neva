package src

type Program struct {
	Packages map[PkgRef]Package
	RootPkg  PkgRef
}
