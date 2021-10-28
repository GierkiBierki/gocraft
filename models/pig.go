package models

import (
	"fmt"

	"github.com/faiface/glhf"
	"github.com/gierkibierki/gocraft/mesh"
	"github.com/gierkibierki/gocraft/meshview"
	"github.com/go-gl/mathgl/mgl32"
)

type Entity struct {
	pos mgl32.Vec3
	mesh *mesh.Mesh
	shader *glhf.Shader
}

func (e *Entity) Draw()  {
    e.shader.Begin()
	// s := g.camera.State();
	// fmt.Printf("Rysuje swinie %.f %.f %.f \n", s.X, s.Y, s.Z)
	proj := mgl32.Translate3D(0, 0, 10).Mul4(mgl32.Scale3D(10, 10, 10))
	e.shader.SetUniformAttr(0, proj)
	e.mesh.Draw()
	e.shader.End()
}

func (e *Entity) Test(){}

func LoadModel(initialPos mgl32.Vec3) *Entity {
	path := "./assets/pig/Pig.obj"
	var err error;

	meshData, err := meshview.LoadMesh(path)
	if err != nil {
		fmt.Println("Error loding mesh", err)
		panic(err)
	}

	var shader *glhf.Shader;


	model := &Entity {
		pos: initialPos,
	}

	fmt.Println("Tu jestem")
	shader, err = glhf.NewShader(glhf.AttrFormat{
		glhf.Attr{Name: "pos", Type: glhf.Vec3},
		// glhf.Attr{Name: "trans", Type: glhf.Vec2},
		// glhf.Attr{Name: "normal", Type: glhf.Vec3},
	}, glhf.AttrFormat{
		glhf.Attr{ Name: "trans", Type: glhf.Mat4 },
		// glhf.Attr{Name: "frag_colour", Type: glhf.Vec4},
	}, vertexShaderSource, fragmentShaderSource)
	

	if err != nil {
		fmt.Println("shader err", err)
		panic(err)
	}

	model.mesh = mesh.NewMesh(shader, meshData.Buffer)
	model.shader = shader;


	return model
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
			frag_colour = vec4(1, 1, 1, 1);
		}
	` + "\x00"
)
