# views

```Go

    v := views.New(views.Config{
		Dir:        "./examples",
		Extensions: []string{".tmpl"},
		Compress:   true,
		Dev:        false,
	})

	v.Func(template.FuncMap{
		"tolower": strings.ToLower,
	})

	if err := v.Load(); err != nil {
		log.Fatalln("views load", err)
	}

	if err := v.Execute(os.Stdout, "child.tmpl", nil); err != nil {
		log.Fatalln("views Execute", err)
	}


```

```code
PASS
coverage: 83.8% of statements
ok      views   0.279s
```