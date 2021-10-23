package main

import (
	"fmt"
	"strings"

	"github.com/faiface/glhf"
	"github.com/faiface/mainthread"
	"github.com/gierkibierki/gocraft/meshview"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Entity struct {
	pos Vec3
	mesh *Mesh
}

func (e *Entity) Draw()  {
	fmt.Printf("Rysuje swinie...")
	e.mesh.Draw()
}

func ZaladujSwinke(initialPos Vec3) *Entity {
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
			// glhf.Attr{Name: "tex", Type: glhf.Vec2},
			// glhf.Attr{Name: "normal", Type: glhf.Vec3},
		}, glhf.AttrFormat{
			// glhf.Attr{Name: "frag_colour", Type: glhf.Vec4},
		}, vertexShaderSource, fragmentShaderSource)
		
	
		fmt.Println("swinoshader.size", shader.VertexFormat().Size())
	
		if err != nil {
			fmt.Println("Swinoshader sie wyjebal", err)
			panic(err)
		}

		swinka.mesh = NewMesh(shader, meshData.Buffer)
	})


	return swinka
}

const (
	vertexShaderSource = `
		#version 330 core

		in vec3 pos;

		void main() {
			gl_Position = vec4(pos, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 330 core
		out vec4 frag_colour;
		void main() {
			frag_colour = vec4(1, 1, 1, 1.0);
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

func drawTriangles(points []float32) {    
	vao := makeVao(points)
    gl.BindVertexArray(vao)
    gl.DrawArrays(gl.TRIANGLES, 0, int32(len(points) / 3))
}

func createShaders() uint32 {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}


