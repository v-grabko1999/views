package views_test

import (
	"html/template"
	"os"
	"strings"
	"views"
)

func ExampleNew() {
	xt := views.New(views.Config{
		Dir:        "./examples",
		Extensions: []string{".tmpl"},
	})

	xt.Func(template.FuncMap{
		"tolower": strings.ToLower,
	})
	_ = xt.Load()

	_ = xt.Execute(os.Stdout, "child.tmpl", nil)
	/* Output: Hello from child.tmpl
	Hello from partials/question.tmpl */
}
