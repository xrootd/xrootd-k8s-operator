package utils

import (
	"reflect"

	"github.com/BurntSushi/ty"
	"github.com/BurntSushi/ty/fun"
)

func call2(f reflect.Value, args ...reflect.Value) (reflect.Value, reflect.Value) {
	ret := f.Call(args)
	return ret[0], ret[1]
}

// Map has a parametric type:
//
//	func Map(f func(A) B, xs []A) []B
//
// Map returns the list corresponding to the return value of applying
// `f` to each element in `xs`.
func Map(f, xs interface{}) interface{} {
	return fun.Map(f, xs)
}

// MapWithError has a parametric type:
//
//	func Map(f func(A) (B, error), xs []A) ([]B, error)
//
// Map returns the list corresponding to the return value of applying
// `f` to each element in `xs`.
func MapWithError(f, xs interface{}) (result interface{}, err error) {
	chk := ty.Check(
		new(func(func(ty.A) (ty.B, error), []ty.A) []ty.B),
		f, xs)
	vf, vxs, tys := chk.Args[0], chk.Args[1], chk.Returns[0]

	xsLen := vxs.Len()
	vys := reflect.MakeSlice(tys, xsLen, xsLen)
	for i := 0; i < xsLen; i++ {
		vy, ey := call2(vf, vxs.Index(i))
		err := ey.Interface()
		if err != nil {
			break
		}
		vys.Index(i).Set(vy)
	}
	return vys.Interface(), err
}
