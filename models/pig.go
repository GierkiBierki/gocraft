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
	proj := mgl32.Translate3D(0, 0, 0).Mul4(mgl32.Scale3D(0.01, 0.01, 0.01))
	e.shader.SetUniformAttr(0, proj)
	e.mesh.Draw()
	e.shader.End()
}

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
		glhf.Attr{Name: "pos", Type: glhf.Vec4},
	}, glhf.AttrFormat{
		glhf.Attr{ Name: "trans", Type: glhf.Mat4 },
	}, vertexShaderSource, fragmentShaderSource)
	

	if err != nil {
		fmt.Println("shader err", err)
		panic(err)
	}

	fmt.Println("Shader size", shader.VertexFormat().Size())

	model.mesh = mesh.NewMesh(shader, meshData.Buffer)
	model.shader = shader;


	return model
}

const (
	vertexShaderSource = `
		#version 330 core

		in vec3 pos;

		uniform mat4 trans;

		out vec3 ec_pos;

		void main() {
			gl_Position = trans * vec4(pos, 1.0);
			ec_pos = vec3(gl_Position);
		}
	` + "\x00"

	fragmentShaderSource = `
#version 330 core

in vec3 ec_pos;
out vec4 frag_colour;

vec3 light_direction = normalize(vec3(1, -1.5, 1));
vec3 object_color = vec3(0x5b / 255.0, 0xac / 255.0, 0xe3 / 255.0);

void main() {

	vec3 ec_normal = normalize(cross(dFdx(ec_pos), dFdy(ec_pos)));
	float diffuse = max(0, dot(ec_normal, light_direction)) * 0.9 + 0.15;
	vec3 color = object_color * diffuse;
	frag_colour = vec4(color, 1);
}
`

	// fragmentShaderSource2 = `
	// 	#version 330 core
	// 	out vec4 frag_colour;
	// 	void main() {
	// 		frag_colour = vec4(1, 1, 1, 1);
	// 	}
	// ` + "\x00"
)
