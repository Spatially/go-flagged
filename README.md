# `usage:"Declarative flags."`

`flagged` utilizes struct-tags to register flags and positional parameters. If there isn't at least a `usage:` struct-tag, the element is ignored.

# Status

[![Build Status](https://travis-ci.org/Spatially/go-flagged.png?branch=master)](https://travis-ci.org/Spatially/go-flagged)


# Installation

```{bash}
go get -v github.com/Spatially/go-flagged
```


# Documentation

See [GoDoc](http://godoc.org/github.com/Spatially/go-flagged) or [Go Walker](http://gowalker.org/github.com/Spatially/go-flagged) for automatically generated documentation.


## usage:

Given a tagged struct:

```{go}
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
```

Register the flags and view the usage:

```{go}
flagged.Parse(&setting)
flag.Usage()
```

Produces:

	usage: ex [flags]
	flags:
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

```{go}
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
```

Which produces:

	usage: ex [flags]
	flags:
	  -a.bool=true: A Bool.
	  -a.float=123.456: A Float.
	  -a.string="default string": A String.
	  -ints.an.int=-1: An Int.
	  -ints.an.int64=-2: An Int64.
	  -ints.an.uint=1: An Uint.
	  -ints.an.uint64=2: An Uint64.

## env:

It can also default the variable from an environment variable based on the `env:` struct-tag. The environment variable value will override a `value:`.

```{go}
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
```

Which produces:

	usage: ex [flags]
	flags:
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

```{go}
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
```

Which produces:

	usage: ex [flags]
	flags:
	  -flag.bool=true: A Bool. (ABOOL)
	  -flag.float=123.456: A Float. (AFLOAT)
	  -flag.string="default string": A String. (ASTRING)
	  -ints.flag.int=-1: An Int. (ANINT)
	  -ints.flag.int64=-2: An Int64. (ANINT64)
	  -ints.flag.uint=1: An Uint. (ANUINT)
	  -ints.flag.uint64=2: An Uint64. (ANUINT64)

### flag aliases:

The flag can be a comma-delimited list of flag names. Using `_` will include the standard struct-derived named as a flag.

```{go}
var setting struct {
	aliased struct {
		x int `flag:"i,_" env:"ALIASED" value:"123" usage:"An Int."`
		y int `flag:"yy" env:"ALIASED" value:"345" usage:"An Int."`
		z int `flag:"_,a" env:"ALIASED" value:"567" usage:"An Int."`
	}
}
```

Which produces:

	usage: ex [flags]
	flags:
		-a=567: alias for -aliased.z: An Int. (ALIASED)
		-aliased.x=123: alias for -i: An Int. (ALIASED)
		-aliased.z=567: An Int. (ALIASED)
		-i=123: An Int. (ALIASED)
		-yy=345: An Int. (ALIASED)

## positional:

```{go}
var setting struct {
	_positional struct {
		first  string  `usage:"Just the first"`
		second int     `usage:"Just the second"`
		third  float64 `usage:"Just the third" value:"3.0"`
		fourth uint64  `usage:"Just the fourth" value:"4"`
		fifth  int64   `usage:"Just the fifth"`
		sixth  string  `usage:"Just the sixth" value:"sixth"`
	}
}
```

Which produces:

	usage: ex [flags] first second [third [fourth [fifth [sixth]]]]
	positional:
	  first (string) - Just the first
	  second (int) - Just the second
	  third=3.0 (float64) - Just the third
	  fourth=4 (uint64) - Just the fourth
	  fifth=0 (int64) - Just the fifth
	  sixth="sixth" (string) - Just the sixth

Note: All positional arguments including and after the first with a `value:` tag are optional arguments.
