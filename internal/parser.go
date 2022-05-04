package internal

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"reflect"
	"strings"
	"unicode"
)

type Parser struct {
	EnabledTag bool
}

func (p Parser) Parse(src string) (File, error) {
	file, err := parser.ParseFile(token.NewFileSet(), src, nil, parser.ParseComments)
	if err != nil {
		return File{}, err
	}
	return p.parseFile(file), nil
}

func (p Parser) parseFile(file *ast.File) File {
	ims := p.parseImports(file)

	var ss []Struct
	for s := range p.parseStructs(file) {
		ss = append(ss, s)
	}
	return File{
		PkgName: file.Name.String(),
		Structs: ss,
		Imports: ims,
	}
}

func (Parser) parseImports(file *ast.File) map[string]string {
	imports := make(map[string]string, len(file.Imports))
	for _, im := range file.Imports {
		name := ""
		if im.Name != nil {
			name = im.Name.Name
		}
		imports[im.Path.Value[1:len(im.Path.Value)-1]] = name
	}
	return imports
}

func (p Parser) parseStructs(file *ast.File) <-chan Struct {
	ch := make(chan Struct)
	go func() {
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}

				fs := make([]Field, 0, len(st.Fields.List))
				for _, f := range st.Fields.List {
					fs = append(fs, p.parseField(f, ts.Name.Name)...)
				}

				ch <- Struct{
					Name:   ts.Name.Name,
					Fields: fs,
				}
			}
		}
		close(ch)
	}()
	return ch
}

func (p Parser) parseField(f *ast.Field, structName string) []Field {
	fs := make([]Field, 0, len(f.Names))
	for _, n := range f.Names {
		comments := []string{}
		if f.Doc != nil {
			comments = make([]string, 0, len(f.Doc.List))
			for _, c := range f.Doc.List {
				if c != nil {
					comments = append(comments, c.Text)
				}
			}
		}

		if p.EnabledTag {
			if f.Tag != nil {
				tag := f.Tag.Value[1 : len(f.Tag.Value)-1] // remove back quotes(``)
				init, get, set := p.parseTag(tag)

				if !isPrivateField(n.Name) {
					log.Printf("[goro]: %s in %s is public. "+
						"Public fields do not require getter and setter.", n.Name, structName)
					get = false
					set = false
				}

				fs = append(fs, Field{
					Name:     n.Name,
					Type:     types.ExprString(f.Type),
					Comments: comments,
					Init:     init,
					Get:      get,
					Set:      set,
				})
			}
		} else if isPrivateField(n.Name) {
			fs = append(fs, Field{
				Name:     n.Name,
				Type:     types.ExprString(f.Type),
				Comments: comments,
				Init:     true,
				Get:      true,
			})
		}
	}
	return fs
}

func (Parser) parseTag(tag string) (init bool, get bool, set bool) {
	t := reflect.StructTag(tag)
	v := t.Get("goro")
	if v == "" {
		return
	}
	for _, s := range strings.Split(v, ",") {
		switch s {
		case "init":
			init = true
		case "get":
			get = true
		case "set":
			set = true
		}
	}
	return
}

func isPrivateField(name string) bool {
	if len(name) == 0 {
		return false
	}
	return name[0] == '_' || unicode.IsLower(rune(name[0]))
}
