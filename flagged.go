// `flagged` utilizes struct-tags to register flags. It can default the variable explicitly
// based on the `value:` struct-tag as well as from an environment variable based on the
// `env:` struct-tag. If there isn't at least a `usage:` struct-tag, the element is ignored.
// The flag can be named explicitly using the `flag:` struct-tag.
package flagged

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

var (
	// Show some detailed output about the registered flags.
	FlaggedDebugging = false
	// A global prefix.
	Prefix = ""
	// The separator to use.
	Separator = "."
)

var selfz, selfd, selfu = os.Args[0], "", ""

func init() {
	selfd, selfu = "", filepath.Base(selfz)
	if ln, err := filepath.EvalSymlinks(selfz); err != nil {
		log.Fatal(err)
	} else {
		if dir, err := filepath.Abs(filepath.Dir(ln)); err != nil {
			log.Fatal(err)
		} else {
			selfd = dir
		}
	}
}

// Returns the executable's basename and absolute path.
func Program() (string, string) {
	return selfu, selfd
}

type options uint8

const (
	// Pass flagged.Continue to flagged.Parse functions if you want to register different structs.
	// Not passing Continue will cause the flagged.Parse functions to call flag.Parse, meaning that
	// any subsequent calls to flagged.Parse will not have any impact.
	Continue options = iota
)

//
type environment map[string]string

//
func (env environment) get(key string, def string) string {
	if item, ok := env[key]; ok {
		return item
	} else {
		return def
	}
}

//
func getenvironment(data []string, getkeyval func(item string) (key, val string)) environment {
	environment := make(map[string]string)
	for _, item := range data {
		key, val := getkeyval(item)
		environment[key] = val
	}
	return environment
}

//
var (
	re_deCamel = regexp.MustCompile(`([A-Z])`)
	Usage      = usage()
	internal   func()
)

//
func usage() func() {
	var once sync.Once
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %[1]s [flags]\n", selfu)
		fmt.Fprintf(os.Stderr, "flags:\n")
		flag.PrintDefaults()
	}
	return func() {
		once.Do(func() {
			if internal == nil {
				internal = flag.Usage
			}
			flag.Usage = internal
		})
		internal()
	}
}

// Pass a pointer to the flagged-tagged struct and Parse will register all the flags.
// Given a tagged struct:
//
// 	var setting struct {
// 		aString string  `usage:"A String."`
// 	}
//
// Register the flags and view the usage:
//
// 	  flagged.Parse(&setting)
func Parse(value_ interface{}, options_ ...options) {
	ParseWithPrefix(value_, Prefix, options_...)
}

//
func unpack(s []string, vars ...positional) error {
	i, ls, lv := 0, len(s), len(vars)
	for ; i < lv; i++ {
		param := vars[i]
		value := ""
		if i < ls {
			value = s[i]
		} else if param.optional {
			value = param.value
		} else {
			return fmt.Errorf("positional parameter %d, %s, is not optional, has no default, and was not provided", i+1, param.name)
			continue
		}

		if param.parameter != nil {
			v := param

			switch p := (v.parameter).(type) {
			case *string:
				*p = value
			case *float64:
				if f, err := strconv.ParseFloat(value, 64); err != nil {
					return fmt.Errorf("positional parameter %d is not of type %s: %+v", i+1, v.type_, value)
				} else {
					*p = f
				}
			case *int64, *uint64, *int, *uint:
				var default_ interface{}
				switch p.(type) {
				case *int64, *int:
					if i, err := strconv.ParseInt(value, 10, 64); err == nil {
						default_ = i
					}
				case *uint64, *uint:
					if i, err := strconv.ParseUint(value, 10, 64); err == nil {
						default_ = i
					}
				}

				if default_ == nil {
					return fmt.Errorf("positional parameter %d is not of type %s: %+v", i+1, v.type_, value)
				}

				switch p := p.(type) {
				case *int64:
					*p = default_.(int64)
				case *int:
					*p = int(default_.(int64))
				case *uint64:
					*p = default_.(uint64)
				case *uint:
					*p = uint(default_.(uint64))
				}
			}
		}
	}
	if i < ls {
		fmt.Fprintf(os.Stderr, "Additional positional arguments ignored: %+v\n", s[i:])
	}

	return nil
}

// Pass a pointer to the flagged-tagged struct and a prefix string and ParseWithPrefix will
// register all the flags prefixed with `prefix`.
func ParseWithPrefix(value_ interface{}, prefix string, options_ ...options) {
	// Just so that this package does not keep the environment itself.
	var env = getenvironment(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = splits[1]
		return
	})
	positionals := parser(value_, prefix, parental{}, env)

	parse := true
	for _, o := range options_ {
		switch o {
		case Continue:
			parse = false
		}
	}
	if parse {

		if positionals == nil || len(positionals) == 0 {
			flag.Parse()
		} else {

			sp := ""
			oa := []string{}
			for _, p := range positionals {
				if !p.optional {
					sp = fmt.Sprintf("%s %s", sp, p.name)
				} else {
					oa = append([]string{p.name}, oa...)
				}
			}
			so := ""
			for _, p := range oa {
				so = fmt.Sprintf(" [%s%s]", p, so)
			}
			internal = func() {
				fmt.Fprintf(os.Stderr, "usage: %[1]s [flags]%s%s\n", selfu, sp, so)
				fmt.Fprintf(os.Stderr, "positional:\n")
				for _, p := range positionals {
					fmt.Fprintf(os.Stderr, "  %s", p.name)
					if p.optional && p.value != "" {
						switch p.parameter.(type) {
						case *string:
							fmt.Fprintf(os.Stderr, `="%s"`, p.value)
						default:
							fmt.Fprintf(os.Stderr, `=%s`, p.value)
						}
					}
					fmt.Fprintf(os.Stderr, " (%s)", p.type_)
					fmt.Fprintf(os.Stderr, " - %s\n", p.usage)
				}
				fmt.Fprintf(os.Stderr, "flags:\n")
				flag.PrintDefaults()
			}

			flag.Parse()

			if err := unpack(flag.Args(), positionals...); err != nil {
				fmt.Fprintln(os.Stderr, err)
				if Usage != nil {
					Usage()
					os.Exit(1)
				}
			}
		}
	}
}

// TODO optional can probably be enhanced to support optional based on type.
type positional struct {
	name, value, usage string
	parameter          interface{}
	optional           bool // Once a positional is optional, all subsequent are optional.
	type_              reflect.Kind
}

//
type parental struct {
	name, usage string
	positional  bool
}

//
func parser(value_ interface{}, prefix string, parent parental, env environment) (ps []positional) {
	// defer func() { debug.PrintStack() }()

	switch t := value_.(type) {
	case reflect.Value:
		switch t.Kind() {
		case reflect.Ptr:
			ps = parser(t.Elem(), prefix, parent, env)
		case reflect.Struct:
			optional := false
			for f, fs := 0, t.NumField(); f < fs; f++ {
				field := t.Field(f)
				switch kind := field.Kind(); kind {
				case reflect.Ptr:
				case reflect.Struct:
					t_ := t.Type()
					f_ := t_.Field(f)
					name := strings.ToLower(prefix)
					if name == "" {
						name = f_.Name
					} else {
						name = fmt.Sprintf("%s%s%s", name, Separator, f_.Name)
					}
					if parent.positional {
						log.Printf("WARNING structs within _positional is unsupported: %s", name)
						return
					}
					p := parental{
						name:       fmt.Sprintf("%s.%s", parent.name, f_.Name),
						usage:      fmt.Sprintf("%s%s", parent.usage, f_.Tag.Get("usage")),
						positional: f_.Name == "_positional",
					}
					if ps_ := parser(field, strings.ToLower(name), p, env); p.positional && len(ps_) > 0 {
						ps = ps_
					}
				default:
					if field.CanAddr() {
						t_ := t.Type()
						f_ := t_.Field(f)
						tag := f_.Tag
						if tag == "" {
							continue
						}
						usage := tag.Get("usage")
						if usage == "" {
							continue
						}
						name := tag.Get("flag")
						value := tag.Get("value")

						var p interface{}
						switch kind {
						case reflect.Bool:
							p = (*bool)(unsafe.Pointer(field.Addr().Pointer()))
						case reflect.String:
							p = (*string)(unsafe.Pointer(field.Addr().Pointer()))
						case reflect.Float64:
							p = (*float64)(unsafe.Pointer(field.Addr().Pointer()))
						case reflect.Int64:
							p = (*int64)(unsafe.Pointer(field.Addr().Pointer()))
						case reflect.Int:
							p = (*int)(unsafe.Pointer(field.Addr().Pointer()))
						case reflect.Uint64:
							p = (*uint64)(unsafe.Pointer(field.Addr().Pointer()))
						case reflect.Uint:
							p = (*uint)(unsafe.Pointer(field.Addr().Pointer()))
						}

						if parent.positional {
							optional = optional || tag.Get("value") != ""
							if optional && value == "" {
								switch p.(type) {
								case *string:
								case *float64:
									value = "0.0"
								default:
									value = "0"
								}
							}
							ps = append(ps, positional{
								name:      f_.Name,
								value:     value,
								usage:     usage,
								type_:     kind,
								parameter: p,
								optional:  optional,
							})
							continue
						}

						envrionment_ := tag.Get("env")
						if envrionment_ != "" {
							value = env.get(envrionment_, value)
							usage = fmt.Sprintf("%s (%s)", usage, envrionment_)
						}

						destination := fmt.Sprintf("%s.%s", parent.name, f_.Name)

						aliased := ""
						names := strings.Split(name, ",")
						{
							previous := ""
							for i, name := range names {
								name = strings.TrimSpace(name)
								if name == "_" || (i == 0 && name == "") {
									name = f_.Name
									name = fmt.Sprintf("%s%s", strings.ToLower(name[:1]), strings.ToLower(re_deCamel.ReplaceAllString(name[1:], `.$1`)))
									if prefix != "" {
										name = fmt.Sprintf("%s%s%s", prefix, Separator, name)
									}
								} else if i > 0 && aliased == "" {
									aliased = previous
								}
								if aliased == "" {
									aliased = name
								}
								previous = name
								names[i] = name
							}
						}

						for _, name := range names {
							if name == "" {
								continue
							}
							description := usage
							if aliased != "" && name != aliased {
								description = fmt.Sprintf("alias for -%s: %s", aliased, usage)
							} else {
								description = fmt.Sprintf("%s%s", parent.usage, description)
							}

							var default_ interface{} = value

							if err := func() (err error) {
								defer func() {
									if r := recover(); r != nil {
										err = fmt.Errorf("%+v", r)
									}
								}()

								switch p := p.(type) {
								case *bool:
									default_ = value == "true"
									flag.BoolVar(p, name, default_.(bool), description)

								case *string:
									flag.StringVar(p, name, value, description)

								case *float64:
									if f, err := strconv.ParseFloat(value, 64); err != nil {
										default_ = field.Float()
									} else {
										default_ = f
									}
									flag.Float64Var(p, name, default_.(float64), description)

								case *int64, *uint64, *int, *uint:

									switch p.(type) {
									case *int64, *int:
										if i, err := strconv.ParseInt(value, 10, 64); err != nil {
											default_ = field.Int()
										} else {
											default_ = i
										}
									case *uint64, *uint:
										if i, err := strconv.ParseUint(value, 10, 64); err != nil {
											default_ = field.Uint()
										} else {
											default_ = i
										}
									}

									switch p := p.(type) {
									case *int64:
										flag.Int64Var(p, name, default_.(int64), description)
									case *int:
										flag.IntVar(p, name, int(default_.(int64)), description)
									case *uint64:
										flag.Uint64Var(p, name, default_.(uint64), description)
									case *uint:
										flag.UintVar(p, name, uint(default_.(uint64)), description)
									}

								default:
									return fmt.Errorf("%T flags are not currently supported", field.Interface())
								}
								if FlaggedDebugging {
									fmt.Printf("%8s -%-40s env:%-20s value:%-20v usage:%s (%s)\n", field.Kind(), name, envrionment_, default_, description, destination)
								}
								return nil
							}(); err != nil {
								log.Println(err)
								continue
							}

						}
					}
				}
			}
		}
	default:
		ps = parser(reflect.ValueOf(value_), prefix, parent, env)
	}
	return
}
