package main

import (
	"fmt"

	"github.com/gierkibierki/gocraft/meshview"
)

type Entity struct {
	pos Vec3
	mesh *meshview.MeshData
}

func (e *Entity) Draw()  {

}

func ZaladujSwinke(initialPos Vec3) Entity {
	path := "./assets/pig/Pig.obj"
	mesh, err := meshview.LoadMesh(path)
	if err != nil {
		fmt.Println("Error przy ladowaniu swini", err)
		panic(err)
	}
	return Entity {
		pos: initialPos,
		mesh: mesh,
	};
}
