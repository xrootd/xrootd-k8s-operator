package utils

import "github.com/BurntSushi/ty/fun"

func Map(f, xs interface{}) interface{} {
	return fun.Map(f, xs)
}
