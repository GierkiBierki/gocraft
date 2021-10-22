package main

import (
	"fmt"

	"github.com/g3n/engine/loader/obj"
)

type Entity struct {
	pos Vec3
	obiekty []obj.Object
}

func (e *Entity) Draw()  {
	for _, o := range e.obiekty {
		// vxs := make([]int, 3)
		fmt.Print(o)
		// for _, f := range o.Faces {
		// 	vxs = append(vxs, f.Vertices...)
		// }

		// var VAO uint32
		// gl.GenVertexArrays(1, &VAO)
	
		// var VBO uint32
		// gl.GenBuffers(1, &VBO)
		
		// // Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
		// gl.BindVertexArray(VAO)
	
		// // copy vertices data into VBO (it needs to be bound first)
		// gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
		// gl.BufferData(gl.ARRAY_BUFFER, len(vxs)*4, gl.Ptr(vxs), gl.STATIC_DRAW)

		// gl.BindVertexArray(0)
	}

}

func ZaladujSwinke(initialPos Vec3) Entity {
	dec, err := obj.Decode("../assets/pig/Pig.obj", "../assets/pig/Pig.mtl")
	if err != nil {
		panic(err)
	}
	return Entity {
		pos: initialPos,
		obiekty: dec.Objects,
	};
	// mat := material.NewStandard(math32.NewColor("DarkBlue"))
	// for _, o := range dec.Objects {
	// 	mesh, err := dec.NewMesh(&o)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	scene.Add(mesh)
	// }
}
