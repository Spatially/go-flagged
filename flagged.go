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
	"reflect"
	"regexp"
	"strconv"
	"strings"
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
var re_deCamel = regexp.MustCompile(`([A-Z])`)

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
	parser(value_, prefix, "", env)
	parse := true
	for _, o := range options_ {
		switch o {
		case Continue:
			parse = false
		}
	}
	if parse {
		flag.Parse()
	}
}

//
func parser(value_ interface{}, prefix string, parent string, env environment) {
	switch t := value_.(type) {
	case reflect.Value:
		switch t.Kind() {
		case reflect.Ptr:
			parser(t.Elem(), prefix, parent, env)
		case reflect.Struct:
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
					parent = fmt.Sprintf("%s.%s", parent, f_.Name)
					parser(field, strings.ToLower(name), parent, env)
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

						envrionment_ := tag.Get("env")
						if envrionment_ != "" {
							value = env.get(envrionment_, value)
							usage = fmt.Sprintf("%s (%s)", usage, envrionment_)
						}

						destination := fmt.Sprintf("%s.%s", parent, f_.Name)

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
							}

							var default_ interface{} = value

							if err := func() (err error) {
								defer func() {
									if r := recover(); r != nil {
										err = fmt.Errorf("%+v", r)
									}
								}()

								switch kind {
								case reflect.Bool:
									p := (*bool)(unsafe.Pointer(field.Addr().Pointer()))
									default_ = value == "true"
									flag.BoolVar(p, name, default_.(bool), description)

								case reflect.String:
									p := (*string)(unsafe.Pointer(field.Addr().Pointer()))
									flag.StringVar(p, name, value, description)

								case reflect.Float64:
									p := (*float64)(unsafe.Pointer(field.Addr().Pointer()))
									if f, err := strconv.ParseFloat(value, 64); err != nil {
										default_ = field.Float()
									} else {
										default_ = f
									}
									flag.Float64Var(p, name, default_.(float64), description)

								case reflect.Int64, reflect.Uint64, reflect.Int, reflect.Uint:

									switch kind {
									case reflect.Int64, reflect.Int:
										if i, err := strconv.ParseInt(value, 10, 64); err != nil {
											default_ = field.Int()
										} else {
											default_ = i
										}
									case reflect.Uint64, reflect.Uint:
										if i, err := strconv.ParseUint(value, 10, 64); err != nil {
											default_ = field.Uint()
										} else {
											default_ = i
										}
									}

									switch kind {
									case reflect.Int64:
										p := (*int64)(unsafe.Pointer(field.Addr().Pointer()))
										flag.Int64Var(p, name, default_.(int64), description)
									case reflect.Int:
										p := (*int)(unsafe.Pointer(field.Addr().Pointer()))
										flag.IntVar(p, name, int(default_.(int64)), description)
									case reflect.Uint64:
										p := (*uint64)(unsafe.Pointer(field.Addr().Pointer()))
										flag.Uint64Var(p, name, default_.(uint64), description)
									case reflect.Uint:
										p := (*uint)(unsafe.Pointer(field.Addr().Pointer()))
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
		parser(reflect.ValueOf(value_), prefix, parent, env)
	}
}
