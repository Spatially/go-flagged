package main

import (
	"flag"
	"log"

	flagged "github.com/Urban4M/go-flagged"
)

//
var setting struct {
	aString string  `usage:"A String."`
	aBool   bool    `usage:"A Bool."`
	aFloat  float64 `usage:"A Float."`
	ints    struct {
		anInt    int    `usage:"An Int."`
		anInt64  int64  `usage:"An Int64."`
		anUint   uint   `usage:"An Uint."`
		anUint64 uint64 `usage:"An Uint64."`
	}
	values struct {
		aString string  `value:"default string" usage:"A String."`
		aBool   bool    `value:"true" usage:"A Bool."`
		aFloat  float64 `value:"123.456" usage:"A Float."`
		ints    struct {
			anInt    int    `value:"-1" usage:"An Int."`
			anInt64  int64  `value:"-2" usage:"An Int64."`
			anUint   uint   `value:"1" usage:"An Uint."`
			anUint64 uint64 `value:"2" usage:"An Uint64."`
		}
	}
	envs struct {
		aString string  `env:"ASTRING" value:"default string" usage:"A String."`
		aBool   bool    `env:"ABOOL" value:"true" usage:"A Bool."`
		aFloat  float64 `env:"AFLOAT" value:"123.456" usage:"A Float."`
		ints    struct {
			anInt    int    `env:"ANINT" value:"-1" usage:"An Int."`
			anInt64  int64  `env:"ANINT64" value:"-2" usage:"An Int64."`
			anUint   uint   `env:"ANUINT" value:"1" usage:"An Uint."`
			anUint64 uint64 `env:"ANUINT64" value:"2" usage:"An Uint64."`
		}
	}
	named struct {
		aString string  `flag:"flag.string" env:"ASTRING" value:"default string" usage:"A String."`
		aBool   bool    `flag:"flag.bool" env:"ABOOL" value:"true" usage:"A Bool."`
		aFloat  float64 `flag:"flag.float" env:"AFLOAT" value:"123.456" usage:"A Float."`
		ints    struct {
			anInt    int    `flag:"flag.int" env:"ANINT" value:"-1" usage:"An Int."`
			anInt64  int64  `flag:"flag.int64" env:"ANINT64" value:"-2" usage:"An Int64."`
			anUint   uint   `flag:"flag.uint" env:"ANUINT" value:"1" usage:"An Uint."`
			anUint64 uint64 `flag:"flag.uint64" env:"ANUINT64" value:"2" usage:"An Uint64."`
		}
	}
}

//
func main() {
	flagged.FlaggedDebugging = true
	// This is just an example of how you can use parts of a struct instead of the entire struct:
	flagged.ParseWithPrefix(&setting.ints, "ok", flagged.Continue)
	flagged.ParseWithPrefix(&setting.named, "weird", flagged.Continue)
	// This is typically all you'll need to call:
	flagged.Parse(&setting)
	flag.Usage()
	log.Printf("%+v", setting)
}
