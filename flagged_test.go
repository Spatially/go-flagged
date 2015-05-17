package flagged

import (
	"flag"
	"testing"
)

//
var setting struct {
	unnamed struct {
		aString string  `usage:"A String."`
		aBool   bool    `usage:"A Bool."`
		aFloat  float64 `usage:"A Float."`
		ints    struct {
			anInt    int    `usage:"An Int."`
			anInt64  int64  `usage:"An Int64."`
			anUint   uint   `usage:"An Uint."`
			anUint64 uint64 `usage:"An Uint64."`
		}
	}
	named struct {
		aString string  `flag:"string" env:"STRING" value:"default string" usage:"A String."`
		aBool   bool    `flag:"bool" env:"BOOL" value:"true" usage:"A Bool."`
		aFloat  float64 `flag:"float" env:"FLOAT" value:"123.456" usage:"A Float."`
		ints    struct {
			anInt    int    `flag:"int" env:"INT" value:"123" usage:"An Int."`
			anInt64  int64  `flag:"int64" env:"INT64" value:"456" usage:"An Int64."`
			anUint   uint   `flag:"uint" env:"UINT" value:"789" usage:"An Uint."`
			anUint64 uint64 `flag:"uint64" env:"UINT64" value:"987" usage:"An Uint64."`
		}
		errs struct {
			Bool     bool    `flag:"bool" env:"BOOL" value:"ok" usage:"A Bool."`
			Float    float64 `flag:"float" env:"FLOAT" value:"xyz" usage:"A Float."`
			anInt    int     `flag:"int" env:"INT" value:"abc" usage:"An Int."`
			anInt64  int64   `flag:"int64" env:"INT64" value:"0xff" usage:"An Int64."`
			anUint   uint    `flag:"uint" env:"UINT" value:"-1" usage:"An Uint."`
			anUint64 uint64  `flag:"uint64" env:"UINT64" value:"-twelve" usage:"An Uint64."`
		}
	}
	ignored struct {
		aString string
		aBool   bool
		aFloat  float64
		ints    struct {
			anInt    int
			anInt64  int64
			anUint   uint
			anUint64 uint64
		}
	}
}

//
func TestFlagged(t *testing.T) {
	FlaggedDebugging = true
	Prefix = "test"
	Parse(&setting)
	t.Logf("%+v", setting)
	flag.Usage()
}
