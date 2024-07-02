# views

```Go

   v := views.New(views.Config{
		Dir:        "./examples",
		Extensions: []string{".tmpl"},
		Compress:   true,
		Dev:        false,

		VersionFilePatch: "./examples/version.txt",
		VersionSize:      12,

		Log: func(str string) {
			log.Println(str)
		},
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
coverage:  79.7% of statements
ok      views   0.279s
```