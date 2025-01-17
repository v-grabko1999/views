package views

import (
	"errors"
	"html/template"
	"os"
	"strings"

	"github.com/greatcloak/decimal"
	"github.com/thanhpk/randstr"
)

func funcMap(cfg *Config) template.FuncMap {
	cFuncMap := map[string]interface{}{
		"to_lower":   strings.ToLower,
		"to_upper":   strings.ToUpper,
		"trim_space": strings.TrimSpace,
		"rand_str":   randstr.Hex,

		"multiplication_decimal_string":     MultiplicationDecimalString,
		"multiplication_decimal_string_fix": MultiplicationDecimalStringFix,
		"multiplication_decimal":            MultiplicationDecimal,
		"multiplication_decimal_fix":        MultiplicationDecimalFix,

		"decimal_to_string":   DecimalToString,
		"decimal_from_string": DecimalFromString,

		"is_dev": func() bool {
			return cfg.Dev
		},
	}

	if cfg.Dev {
		delVersionApp(cfg)
		cFuncMap["version"] = func() string {
			return randstr.Hex(cfg.VersionSize)
		}
	} else {
		versionApp := getVersionApp(cfg)
		cFuncMap["version"] = func() string {
			return strings.Clone(versionApp)
		}
	}

	return cFuncMap
}

func DecimalToString(x decimal.Decimal) string {
	return x.String()
}

func DecimalFromString(x string) decimal.Decimal {
	return decimal.RequireFromString(x)
}

func MultiplicationDecimal(n1, n2 decimal.Decimal) decimal.Decimal {
	return n1.Mul(n2)
}

func MultiplicationDecimalString(n1, n2 string) string {
	return decimal.RequireFromString(n1).Mul(decimal.RequireFromString(n2)).String()
}

func MultiplicationDecimalFix(fix int32, n1, n2 decimal.Decimal) string {
	return n1.Mul(n2).StringFixedBank(fix)
}

func MultiplicationDecimalStringFix(fix int32, n1, n2 string) string {
	return decimal.RequireFromString(n1).Mul(decimal.RequireFromString(n2)).StringFixedBank(fix)
}

func getVersionApp(cfg *Config) (versionApp string) {
	file, err := os.ReadFile(cfg.VersionFilePatch)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			versionApp = setVersionApp(cfg)
		} else {
			panic(err)
		}
	} else {
		versionApp = string(file)
		if len(versionApp) != cfg.VersionSize {
			versionApp = setVersionApp(cfg)
		}
	}
	return
}

func delVersionApp(cfg *Config) {
	if err := os.Remove(cfg.VersionFilePatch); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
	}
}

func setVersionApp(cfg *Config) (versionApp string) {
	versionApp = randstr.Hex(cfg.VersionSize)
	if err := os.WriteFile(cfg.VersionFilePatch, []byte(versionApp), 0777); err != nil {
		panic(err)
	}
	return
}
