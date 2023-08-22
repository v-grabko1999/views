package views

import (
	"html/template"
	"io"
	"log"
	"sync"

	"golang.org/x/sync/singleflight"
)

type Config struct {
	//Dir specifies the path to the directory with templates
	Dir string

	//Extensions specifies what file extensions the templates have
	Extensions []string

	//Compress specifies whether to compress templates
	//true - compresses templates during parsing
	//false - does not compress templates
	Compress bool

	//Dev enable/disable development mode
	//true  - parses templates every time
	//false - parses templates once at the very beginning
	Dev bool
}

func New(cfg Config) *Views {
	return &Views{cfg.Dir, cfg.Extensions, cfg.Dev, cfg.Compress, now(), &sync.RWMutex{}, funcMap(&cfg)}
}

type Views struct {
	root         string
	extensions   []string
	dev          bool
	compressHtml bool
	ex           *tpl
	m            *sync.RWMutex
	cacheFuncMap template.FuncMap
}

func (v *Views) Load() error {
	v.m.Lock()
	defer v.m.Unlock()

	n := now()
	n.Funcs(v.cacheFuncMap)

	if err := n.ParseDir(v.root, v.extensions, v.compressHtml, v.dev); err != nil {
		return err
	}

	v.ex = n
	return nil
}

var exGroup singleflight.Group

func (v *Views) Execute(wr io.Writer, name string, data interface{}) error {
	if v.dev {
		_, err, _ := exGroup.Do("load", func() (interface{}, error) {
			log.Println("views mode dev, reload views.")
			return nil, v.Load()
		})

		if err != nil {
			return err
		}
	}
	return v.ex.ExecuteTemplate(wr, name, data)
}

func (v *Views) Func(fn template.FuncMap) {
	v.m.Lock()
	defer v.m.Unlock()

	for k, val := range fn {
		_, ok := v.cacheFuncMap[k]
		if ok {
			panic("Views: Функция `" + k + "` уже обьявлена")
		}
		v.cacheFuncMap[k] = val
	}
}
