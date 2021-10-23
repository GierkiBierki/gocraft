package main

import (
	"fmt"

	"github.com/faiface/glhf"
	"github.com/faiface/mainthread"
	"github.com/gierkibierki/gocraft/meshview"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Entity struct {
	pos mgl32.Vec3
	mesh *Mesh
	shader *glhf.Shader
}

func (e *Entity) Draw(g *Game)  {
    e.shader.Begin()
	s := g.camera.State();
	fmt.Printf("Rysuje swinie %.f %.f %.f \n", s.X, s.Y, s.Z)
	// proj := mgl32.Translate3D(s.X +3, s.Y, s.Z +3).Mul4(mgl32.Scale3D(10, 10, 10))
	e.shader.SetUniformAttr(0, mgl32.Scale3D(10, 10, 10))
	e.mesh.Draw()
	e.shader.End()
}

func ZaladujSwinke(initialPos mgl32.Vec3) *Entity {
	path := "./assets/pig/Pig.obj"
	var err error;

	meshData, err := meshview.LoadMesh(path)
	if err != nil {
		fmt.Println("Error przy ladowaniu swini", err)
		panic(err)
	}

	var shader *glhf.Shader;

	fmt.Println("Shader dla swini tworze", len(meshData.Buffer), len(meshData.Buffer) / 4)


	swinka := &Entity {
		pos: initialPos,
	}

	mainthread.Call(func() {
		shader, err = glhf.NewShader(glhf.AttrFormat{
			glhf.Attr{Name: "pos", Type: glhf.Vec3},
			// glhf.Attr{Name: "trans", Type: glhf.Vec2},
			// glhf.Attr{Name: "normal", Type: glhf.Vec3},
		}, glhf.AttrFormat{
			glhf.Attr{ Name: "trans", Type: glhf.Mat4 },
			// glhf.Attr{Name: "frag_colour", Type: glhf.Vec4},
		}, vertexShaderSource, fragmentShaderSource)
		
	
		fmt.Println("swinoshader.size", shader.VertexFormat().Size())
	
		if err != nil {
			fmt.Println("Swinoshader sie wyjebal", err)
			panic(err)
		}

		swinka.mesh = NewMesh(shader, meshData.Buffer)
		swinka.shader = shader;
	})


	return swinka
}

const (
	vertexShaderSource = `
		#version 330 core

		in vec3 pos;

		uniform mat4 trans;

		void main() {
			gl_Position = trans * vec4(pos, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 330 core
		out vec4 frag_colour;
		void main() {
			frag_colour = vec4(1, 0, 0, 1);
		}
	` + "\x00"
)


// TODO: move to rendering.go or sth like that

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
    var vbo uint32
    gl.GenBuffers(1, &vbo)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)
    
    var vao uint32
    gl.GenVertexArrays(1, &vao)
    gl.BindVertexArray(vao)
    gl.EnableVertexAttribArray(0)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
    
    return vao
}
