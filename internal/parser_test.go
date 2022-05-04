package internal

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	testcases := []struct {
		name       string
		src        string
		enabledTag bool
		want       File
		err        error
	}{
		{
			name:       "tag mode is disabled",
			src:        "../testdata/parse_test.txt",
			enabledTag: false,
			want: File{
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
								Name:     "name",
								Type:     "string",
								Comments: []string{},
								Init:     true,
								Get:      true,
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
								Get:      true,
								Set:      false,
							},
							{
								Name:     "name",
								Type:     "string",
								Comments: []string{},
								Init:     true,
								Get:      true,
								Set:      false,
							},
						},
					},
				},
				Imports: map[string]string{"os": ""},
			},
			err: nil,
		},
		{
			name:       "tag mode is enabled",
			src:        "../testdata/parse_test.txt",
			enabledTag: true,
			want: File{
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
								Get:      false,
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
			err: nil,
		},
		{
			name:       "source file is not found",
			src:        "../foo/parse_test.txt",
			enabledTag: true,
			want:       File{},
			err:        errors.New("open ../foo/parse_test.txt: no such file or directory"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			p := Parser{tc.enabledTag}
			got, err := p.Parse(tc.src)

			if (err != nil) != (tc.err != nil) {
				t.Errorf("Parse returns unexpected error. want: %v, got: %v", tc.err, err)
			} else if err != nil && err.Error() != tc.err.Error() {
				t.Errorf("Parse returns unexpected error. want: %s, got: %s", tc.err.Error(), err.Error())
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Parse returns unexpected result. -want +got: %s", diff)
			}
		})
	}
}

func TestParseTag(t *testing.T) {
	testcases := []struct {
		name string
		tag  string
		init bool
		get  bool
		set  bool
	}{
		{
			name: "empty",
			tag:  "",
		},
		{
			name: "invalid",
			tag:  "invalid",
		},
		{
			name: "empty value",
			tag:  `goro:""`,
		},
		{
			name: "init",
			tag:  `goro:"init"`,
			init: true,
		},
		{
			name: "init and get",
			tag:  `goro:"init,get"`,
			init: true,
			get:  true,
		},
		{
			name: "init and get and set",
			tag:  `goro:"init,get,set"`,
			init: true,
			get:  true,
			set:  true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			p := Parser{}
			init, get, set := p.parseTag(tc.tag)

			if init != tc.init || get != tc.get || set != tc.set {
				t.Errorf("parseTag returns unexpected result. want: (init,get,set)=(%t,%t,%t), got: (init,get,set)=(%t,%t,%t)", tc.init, tc.get, tc.set, init, get, set)
			}
		})
	}
}
