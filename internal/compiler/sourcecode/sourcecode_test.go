package sourcecode

import (
	"testing"

	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

func TestPackage_GetInteropableComponents(t *testing.T) {
	tests := []struct {
		name        string
		pkg         Package
		wantExports []string
	}{
		{
			name: "single_valid_export",
			pkg: Package{
				"file.neva": File{
					Entities: map[string]Entity{
						"FormatUser": {
							IsPublic: true,
							Kind:     ComponentEntity,
							Component: []Component{
								{
									Interface: Interface{
										IO: IO{
											In: map[string]Port{
												"input": {
													TypeExpr: ts.Expr{
														Inst: &ts.InstExpr{
															Ref: core.EntityRef{Name: "string"},
														},
													},
												},
											},
											Out: map[string]Port{
												"output": {
													TypeExpr: ts.Expr{
														Inst: &ts.InstExpr{
															Ref: core.EntityRef{Name: "string"},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantExports: []string{"FormatUser"},
		},
		{
			name: "multiple_valid_exports",
			pkg: Package{
				"file.neva": File{
					Entities: map[string]Entity{
						"FormatUser": {
							IsPublic: true,
							Kind:     ComponentEntity,
							Component: []Component{
								{
									Interface: Interface{
										IO: IO{
											In: map[string]Port{
												"input": {},
											},
											Out: map[string]Port{
												"output": {},
											},
										},
									},
								},
							},
						},
						"ProcessData": {
							IsPublic: true,
							Kind:     ComponentEntity,
							Component: []Component{
								{
									Interface: Interface{
										IO: IO{
											In: map[string]Port{
												"data": {},
											},
											Out: map[string]Port{
												"result": {},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantExports: []string{"FormatUser", "ProcessData"},
		},
		{
			name: "ignore_private_components",
			pkg: Package{
				"file.neva": File{
					Entities: map[string]Entity{
						"FormatUser": {
							IsPublic: true,
							Kind:     ComponentEntity,
							Component: []Component{
								{
									Interface: Interface{
										IO: IO{
											In: map[string]Port{
												"input": {},
											},
											Out: map[string]Port{
												"output": {},
											},
										},
									},
								},
							},
						},
						"helper": {
							IsPublic: false,
							Kind:     ComponentEntity,
							Component: []Component{
								{
									Interface: Interface{
										IO: IO{
											In: map[string]Port{
												"input": {},
											},
											Out: map[string]Port{
												"output": {},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantExports: []string{"FormatUser"},
		},
		{
			name: "ignore_overloaded_components",
			pkg: Package{
				"file.neva": File{
					Entities: map[string]Entity{
						"ValidComponent": {
							IsPublic: true,
							Kind:     ComponentEntity,
							Component: []Component{
								{
									Interface: Interface{
										IO: IO{
											In: map[string]Port{
												"input": {},
											},
											Out: map[string]Port{
												"output": {},
											},
										},
									},
								},
							},
						},
						"OverloadedComponent": {
							IsPublic: true,
							Kind:     ComponentEntity,
							Component: []Component{
								{
									Interface: Interface{
										IO: IO{
											In: map[string]Port{
												"input": {},
											},
											Out: map[string]Port{
												"output": {},
											},
										},
									},
								},
								{
									Interface: Interface{
										IO: IO{
											In: map[string]Port{
												"input": {},
											},
											Out: map[string]Port{
												"output": {},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantExports: []string{"ValidComponent"},
		},
		{
			name: "ignore_non_component_entities",
			pkg: Package{
				"file.neva": File{
					Entities: map[string]Entity{
						"ValidComponent": {
							IsPublic: true,
							Kind:     ComponentEntity,
							Component: []Component{
								{
									Interface: Interface{
										IO: IO{
											In: map[string]Port{
												"input": {},
											},
											Out: map[string]Port{
												"output": {},
											},
										},
									},
								},
							},
						},
						"SomeType": {
							IsPublic: true,
							Kind:     TypeEntity,
						},
						"SomeConst": {
							IsPublic: true,
							Kind:     ConstEntity,
						},
						"SomeInterface": {
							IsPublic: true,
							Kind:     InterfaceEntity,
						},
					},
				},
			},
			wantExports: []string{"ValidComponent"},
		},
		{
			name: "no_exports_in_package",
			pkg: Package{
				"file.neva": File{
					Entities: map[string]Entity{
						"helper": {
							IsPublic: false,
							Kind:     ComponentEntity,
							Component: []Component{
								{
									Interface: Interface{
										IO: IO{
											In: map[string]Port{
												"input": {},
											},
											Out: map[string]Port{
												"output": {},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantExports: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exports := tt.pkg.GetInteropableComponents()

			if len(exports) != len(tt.wantExports) {
				t.Errorf("GetInteropableComponents() got %d exports, want %d", len(exports), len(tt.wantExports))
				return
			}

			// create a map for easier comparison (order doesn't matter)
			gotNames := make(map[string]bool)
			for _, exp := range exports {
				gotNames[exp.Name] = true
			}

			for _, wantName := range tt.wantExports {
				if !gotNames[wantName] {
					t.Errorf("GetInteropableComponents() missing export %q", wantName)
				}
			}
		})
	}
}
