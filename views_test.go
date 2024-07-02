package views_test

import (
	"fmt"
	"html/template"
	"os"
	"strings"
	"testing"

	"application/core/views"
)

func TestNew(t *testing.T) {
	v := views.New(views.Config{
		Dir:        "./examples",
		Extensions: []string{".tmpl"},
		Compress:   true,
		Dev:        false,
		Log: func(str string) {
			fmt.Println(str)
		},
	})

	v.Func(template.FuncMap{
		"tolower": strings.ToLower,
	})

	if err := v.Load(); err != nil {
		t.Fatal("views load", err)
	}

	if err := v.Execute(os.Stdout, "child.tmpl", nil); err != nil {
		t.Fatal("views Execute", err)
	}

	t.Log("views OK")
}
