# `usage:"Declarative flags."`

`flagged` utilizes struct-tags to register flags. If there isn't at least a `usage:` struct-tag, the element is ignored.

# Status

[![Build Status](https://travis-ci.org/Urban4M/go-flagged.png?branch=master)](https://travis-ci.org/Urban4M/go-flagged)


# Installation

```
go get -v github.com/Urban4M/go-flagged
```


# Documentation

See [GoDoc](http://godoc.org/github.com/Urban4M/go-flagged) or [Go Walker](http://gowalker.org/github.com/Urban4M/go-flagged) for automatically generated documentation.


## usage:

Given a tagged struct:

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
		ignored string
	}

Register the flags and view the usage:

	  flagged.Parse(&setting)
	  flag.Usage()

Produces:

	Usage of ./ex:
	  -a.bool=false: A Bool.
	  -a.float=0: A Float.
	  -a.string="": A String.
	  -ints.an.int=0: An Int.
	  -ints.an.int64=0: An Int64.
	  -ints.an.uint=0: An Uint.
	  -ints.an.uint64=0: An Uint64.

Additional take-aways from the above are:

1. All flags are converted to lowercase.
1. Camel-case prepends a `.` before lowercasing.
1. Nested stucts prepend the struct name to the flag name.

_Continuing ... (each following section below builds upon the prior section.)_

## value:

It can default the variable explicitly based on the `value:` struct-tag.

	var setting struct {
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

Which produces:

	Usage of ./ex:
	  -a.bool=true: A Bool.
	  -a.float=123.456: A Float.
	  -a.string="default string": A String.
	  -ints.an.int=-1: An Int.
	  -ints.an.int64=-2: An Int64.
	  -ints.an.uint=1: An Uint.
	  -ints.an.uint64=2: An Uint64.

## env:

It can also default the variable from an environment variable based on the `env:` struct-tag. The environment variable value will override a `value:`.

	var setting struct {
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

Which produces:

	Usage of ./ex:
	  -a.bool=true: A Bool. (ABOOL)
	  -a.float=123.456: A Float. (AFLOAT)
	  -a.string="default string": A String. (ASTRING)
	  -ints.an.int=-1: An Int. (ANINT)
	  -ints.an.int64=-2: An Int64. (ANINT64)
	  -ints.an.uint=1: An Uint. (ANUINT)
	  -ints.an.uint64=2: An Uint64. (ANUINT64)

## flag:

The flag can be named explicitly using the `flag:` struct-tag.

**Be careful** with this. Duplicate registrations causes `flag` to panic. `flagged` recovers from the panic and outputs the error but the flag will obviously not be registered.

	var setting struct {
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

Which produces:

	Usage of ./ex:
	  -flag.bool=true: A Bool. (ABOOL)
	  -flag.float=123.456: A Float. (AFLOAT)
	  -flag.string="default string": A String. (ASTRING)
	  -ints.flag.int=-1: An Int. (ANINT)
	  -ints.flag.int64=-2: An Int64. (ANINT64)
	  -ints.flag.uint=1: An Uint. (ANUINT)
	  -ints.flag.uint64=2: An Uint64. (ANUINT64)
