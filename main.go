package main

import (
	"fmt"
	//"fmt"
	"reflect"
	"strings"
)

func main() {

	/*factory := func(name string) func() {
		return func() {
			fmt.Println(name)
		}
	}
	f1 := factory("f1")
	f2 := factory("f2")

	pf1 := reflect.ValueOf(f1)
	pf2 := reflect.ValueOf(f2)

	fmt.Println(pf1.Pointer(), pf2.Pointer())
	fmt.Println(pf1.Pointer() == pf2.Pointer())

	f1()
	f2()*/
	factory := func(name string) func() {
		return func() {
			fmt.Println(name)
		}
	}
	fmt.Println(signature(factory))
}

func signature(f interface{}) string {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return "<not a function>"
	}

	buf := strings.Builder{}
	buf.WriteString("func (")
	for i := 0; i < t.NumIn(); i++ {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(t.In(i).String())
	}
	buf.WriteString(")")
	if numOut := t.NumOut(); numOut > 0 {
		if numOut > 1 {
			buf.WriteString(" (")
		} else {
			buf.WriteString(" ")
		}
		for i := 0; i < t.NumOut(); i++ {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(t.Out(i).String())
		}
		if numOut > 1 {
			buf.WriteString(")")
		}
	}

	return buf.String()
}
