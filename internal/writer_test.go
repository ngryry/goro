package internal

import (
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	if _, err := os.Stat("../generated"); os.IsNotExist(err) {
		if err := os.Mkdir("../generated", 0777); err != nil {
			t.Fatalf("failed to mkdir ../generated. err: %s", err.Error())
		}
	}

	testcases := []struct {
		name string
		file File
		dst  string
		want string
		err  error
	}{
		{
			name: "",
			file: File{
				PkgName: "test",
				Structs: []Struct{
					{
						Name: "StaticImage",
						Fields: []Field{
							{
								Name:     "file",
								Type:     "os.File",
								Comments: []string{"// image file", "// MIME is png or jpeg."},
								Init:     true,
								Get:      true,
								Set:      false,
							},
							{
								Name:     "Category",
								Type:     "string",
								Comments: []string{},
								Init:     true,
								Get:      false,
								Set:      false,
							},
						},
					},
					{
						Name: "staticHTML",
						Fields: []Field{
							{
								Name:     "file",
								Type:     "os.File",
								Comments: []string{"// html file"},
								Init:     true,
								Get:      false,
								Set:      false,
							},
							{
								Name:     "name",
								Type:     "string",
								Comments: []string{},
								Init:     false,
								Get:      true,
								Set:      true,
							},
							{
								Name:     "Category",
								Type:     "string",
								Comments: []string{},
								Init:     true,
								Get:      false,
								Set:      false,
							},
						},
					},
				},
				Imports: map[string]string{"os": ""},
			},
			dst:  "../generated/write_test.txt",
			want: "../testdata/write_test.txt",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			w := Writer{}
			err := w.Write(tc.file, tc.dst)

			if (err != nil) != (tc.err != nil) {
				t.Errorf("Write returns unexpected error. want: %v, got: %v", tc.err, err)
			} else if err != nil && err.Error() != tc.err.Error() {
				t.Errorf("Write returns unexpected error. want: %s, got: %s", tc.err.Error(), err.Error())
			}
			want, err := os.ReadFile(tc.want)
			if err != nil {
				t.Fatal("failed to open want file.")
			}
			got, err := os.ReadFile(tc.dst)
			if err != nil {
				t.Fatal("failed to open dst file.")
			}
			if string(want) != string(got) {
				t.Errorf("Write generates unexpected Go file.")
			}
		})
	}
}
