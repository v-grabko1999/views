package views

import (
	"html/template"
	"strings"

	"github.com/thanhpk/randstr"
)

func funcMap(cfg *Config) template.FuncMap {
	cFuncMap := make(map[string]interface{})

	if cfg.Dev {
		cFuncMap["version"] = func() string {
			return genBuildVersion()
		}
	} else {
		versionApp := genBuildVersion()
		cFuncMap["version"] = func() string {
			return strings.Clone(versionApp)
		}
	}

	return cFuncMap
}

func genBuildVersion() string {
	return randstr.Hex(8)
}
