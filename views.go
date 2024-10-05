package views

import (
	"errors"
	"html/template"
	"io"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
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

	Log              LogFunc
	VersionFilePatch string
	VersionSize      int
}

type LogFunc func(string)

func New(cfg Config) *Views {
	if cfg.Log == nil {
		cfg.Log = func(str string) {
			log.Println(str)
		}
	}

	if len(cfg.VersionFilePatch) == 0 {
		cfg.VersionSize = 8
		cfg.VersionFilePatch = os.TempDir() + "./golang.views.tml.version.txt"
	}

	if cfg.VersionSize < 1 {
		cfg.VersionSize = 4
	}

	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.Add("text/html", &html.Minifier{})
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.Add("form57/javascript", &js.Minifier{
		KeepVarNames: true,
	})

	return &Views{cfg.Dir, cfg.Extensions, cfg.Dev, cfg.Compress, now(cfg.Log), &sync.RWMutex{}, funcMap(&cfg), cfg.Log, m}
}

type Views struct {
	root         string
	extensions   []string
	dev          bool
	compressHtml bool
	ex           *tpl
	m            *sync.RWMutex
	cacheFuncMap template.FuncMap
	log          LogFunc
	min          *minify.M
}

func (v *Views) Load() error {
	v.m.Lock()
	defer v.m.Unlock()

	n := now(v.log)
	n.Funcs(v.cacheFuncMap)

	if err := n.ParseDir(v.root, v.extensions, v.dev); err != nil {
		return err
	}

	v.ex = n
	return nil
}

var exGroup singleflight.Group

func (v *Views) Execute(wr io.Writer, name string, data interface{}) error {
	if v.dev {
		_, err, _ := exGroup.Do("load", func() (interface{}, error) {
			v.log("views mode dev, reload views.")
			return nil, v.Load()
		})

		if err != nil {
			return err
		}
	}

	if v.compressHtml {
		cl := v.min.Writer("text/html", wr)
		defer cl.Close()
		return v.ex.ExecuteTemplate(cl, name, data)
	} else {
		return v.ex.ExecuteTemplate(wr, name, data)
	}
}

func (v *Views) Func(fn template.FuncMap) error {
	v.m.Lock()
	defer v.m.Unlock()

	for k, val := range fn {
		_, ok := v.cacheFuncMap[k]
		if ok {
			return errors.New("Views: the " + k + " function is already registered")
		}
		v.cacheFuncMap[k] = val
	}
	return nil
}
