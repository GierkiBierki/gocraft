package main

import (
	"os"

	"github.com/sheenobu/go-obj/obj"
)

func LoadObj(filepath string) *obj.Object {
	reader, err := os.Open("tree.obj")
	if err != nil {
		panic(err)
	}

	plik, err := obj.NewReader(reader).Read()
	if err != nil {
		panic(err)
	}
	return plik

}
