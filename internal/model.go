package internal

// File is a Go file.
type File struct {
	PkgName string
	Structs []Struct
	Imports map[string]string // map[import_path]local_name
}

// Struct is a Go struct
type Struct struct {
	Name   string
	Fields []Field
}

// Field is a field of Struct
type Field struct {
	Name     string
	Type     string
	Comments []string
	Init     bool // the field is initialized in constructor
	Get      bool // the field has getter
	Set      bool // the field has sette
}
