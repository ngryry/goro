# goro
Goro (go-readonly) is a Go code generator to automate the creation of constructor, getter and setter.

## Instalation
Install `goro` tool.
```bash
go install github.com/ngryry/goro@v0.1.0
```
**Note** Please make sure `$GOPATH/bin` is in  `$PATH`. 

## Usage
`goro` generates constructors, field getters and setters for all structs in input file, and writes generated code to the output file.

You specify the input file with `-s` flag and the output file with `-d` flag.

`goro` has two generation mode: basic and tag.

### Basic mode
Basic mode generates getters for **private** fields and struct constructors.
Generated constructor initializes only the **private** fields of struct.
Basic mode is enabled by runing without `-t` flag.
```bash
goro -s source.go -d target.go
```
```go
// source.go
package main

import "os"

type A struct {
    // file is a file
    file os.File
    Name string
}
```
```go
// target.go
package main

import "os"

// NewA is constructor for A
func NewA(mFile os.File) A {
    return A{file: mFile}
}

// File is a file
func (a A) File() os.File {
    return a.file
}
```

### Tag mode
You can specify which fields to generate getters and setters for, and which fields to initialize in the struct constructor, by setting `goro` tag.

`goro` supports following tags:
- `init`: initialize the field in constructor
- `get`: generate getter
- `set`: generate setter

Tag mode is enabled by runing with `-t` flag.
```bash
goro -s source.go -d target.go -t
```
```go
// source.go
package main

import "os"

type A struct {
    // file is a file
    file  os.File `goro:"init,get,set"`
    Name  string  `goro:"init,get,set"`
    owner string
}
```
```go
// target.go
package main

import "os"

// NewA is constructor for A
func NewA(mFile os.File, mName string) A {
    return A{file: mFile, Name: mName}
}

// File is a file
func (a A) File() os.File {
    return a.file
}

// SetFile is setter for file
func (a *A) SetFile(file os.File) {
    return a.file = file
}
```
**Note** Name field has a tag `get` and `set`, but goro don't generate getter and setter. This is because public fields do not require getters and setters to reference and edit them.
