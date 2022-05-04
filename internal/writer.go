package internal

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"strconv"
	"strings"
)

type Writer struct{}

func (w Writer) Write(f File, dst string) error {
	af := &ast.File{}
	af.Name = ast.NewIdent(f.PkgName)

	if len(f.Imports) > 0 {
		af.Decls = append(af.Decls, w.buildImports(f.Imports))
	}
	for _, s := range f.Structs {
		af.Decls = append(af.Decls, w.buildConstructor(s))
		af.Decls = append(af.Decls, w.buildFieldFunctions(s)...)
	}

	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := format.Node(file, token.NewFileSet(), af); err != nil {
		return err
	}
	return nil
}

func (Writer) buildImports(ims map[string]string) *ast.GenDecl {
	res := ast.GenDecl{
		Tok: token.IMPORT,
	}
	for path, name := range ims {
		spec := ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote(path),
			},
		}
		if name != "" {
			spec.Name = &ast.Ident{Name: name}
		}
		res.Specs = append(res.Specs, &spec)
	}
	return &res
}

func constructorName(name string) string {
	if name == "" {
		return ""
	}
	return fmt.Sprintf("New%s%s", strings.ToUpper(name[:1]), name[1:])
}

func (Writer) buildConstructor(s Struct) *ast.FuncDecl {
	params := make([]*ast.Field, 0, len(s.Fields))
	resKV := make([]ast.Expr, 0, len(s.Fields))
	for _, f := range s.Fields {
		if f.Init {
			p := fmt.Sprintf("m%s%s", strings.ToUpper(f.Name[:1]), f.Name[1:])
			params = append(params, &ast.Field{
				Names: []*ast.Ident{
					{Name: p},
				},
				Type: &ast.BasicLit{
					Kind:  token.STRING,
					Value: f.Type,
				},
			})
			resKV = append(resKV, &ast.KeyValueExpr{
				Key:   ast.NewIdent(f.Name),
				Value: ast.NewIdent(p),
			})
		}
	}

	results := []*ast.Field{
		{
			Type: &ast.BasicLit{
				Kind:  token.STRING,
				Value: s.Name,
			},
		},
	}

	fn := constructorName(s.Name)
	return &ast.FuncDecl{
		Name: ast.NewIdent(fn),
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: fmt.Sprintf("\n// %s is constructor for %s\n", fn, s.Name),
				},
			},
		},
		Type: &ast.FuncType{
			Params:  &ast.FieldList{List: params},
			Results: &ast.FieldList{List: results},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.Ident{
								Name: s.Name,
							},
							Elts: resKV,
						},
					},
				},
			},
		},
	}
}

func (w Writer) buildFieldFunctions(s Struct) []ast.Decl {
	recv := strings.ToLower(s.Name[:1])
	fs := make([]ast.Decl, 0, len(s.Fields)*2)
	for _, f := range s.Fields {
		if f.Get {
			fs = append(fs, w.buildGetter(recv, s.Name, f))
		}
		if f.Set {
			fs = append(fs, w.buildSetter(recv, s.Name, f))
		}
	}
	return fs
}

func getterName(name string) string {
	if name == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(name[:1]), name[1:])
}

func (w Writer) buildGetter(recv string, recvType string, f Field) ast.Decl {
	results := []*ast.Field{
		{
			Type: &ast.BasicLit{
				Kind:  token.STRING,
				Value: f.Type,
			},
		},
	}
	fn := getterName(f.Name)
	return &ast.FuncDecl{
		Name: ast.NewIdent(fn),
		Doc: &ast.CommentGroup{
			List: w.buildComments(f.Comments, fn, f.Name),
		},
		Type: &ast.FuncType{
			Results: &ast.FieldList{List: results},
		},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(recv)},
					Type:  ast.NewIdent(recvType),
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent(recv),
							Sel: ast.NewIdent(f.Name),
						},
					},
				},
			},
		},
	}
}

func setterName(name string) string {
	if name == "" {
		return ""
	}
	return fmt.Sprintf("Set%s%s", strings.ToUpper(name[:1]), name[1:])
}

func (w Writer) buildSetter(recv string, recvType string, f Field) ast.Decl {
	results := []*ast.Field{
		{
			Type: &ast.BasicLit{
				Kind:  token.STRING,
				Value: f.Type,
			},
		},
	}
	fn := setterName(f.Name)
	return &ast.FuncDecl{
		Name: ast.NewIdent(fn),
		Doc: &ast.CommentGroup{
			List: w.buildComments(f.Comments, fn, f.Name),
		},
		Type: &ast.FuncType{
			Results: &ast.FieldList{List: results},
		},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(recv)},
					Type:  ast.NewIdent(recvType),
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent(recv),
							Sel: ast.NewIdent(f.Name),
						},
					},
				},
			},
		},
	}
}

func (Writer) buildComments(comments []string, funcName string, fieldName string) []*ast.Comment {
	var res []*ast.Comment
	if len(comments) > 0 {
		res = make([]*ast.Comment, len(comments))
		for i, c := range comments {
			if i == 0 {
				// godoc for function start with function name
				prefixs := []string{
					fmt.Sprintf("// %s ", fieldName),
					fmt.Sprintf("// %s ", fieldName),
					"// ",
					"//",
				}
				for _, p := range prefixs {
					if strings.HasPrefix(c, p) {
						c = fmt.Sprintf("// %s %s", funcName, c[len(p):])
						break
					}
				}
				res[i] = &ast.Comment{Text: fmt.Sprintf("\n%s", c)}
			} else {
				res[i] = &ast.Comment{Text: c}
			}
		}
	} else {
		res = []*ast.Comment{
			{Text: fmt.Sprintf("\n// %s", funcName)},
		}
	}
	return res
}
