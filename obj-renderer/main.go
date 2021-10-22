package main

// Based on https://github.com/sheenobu/go-obj/tree/master/cmd/obj-renderer
// https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl
import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sheenobu/go-obj/obj"
)

func initGL(w, h int) *glfw.Window {
	err := glfw.Init()
	if err != nil {
		log.Fatal(err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, gl.TRUE)

	win, err := glfw.CreateWindow(w, h, "gocraft", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	win.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		log.Fatal(err)
	}
	glfw.SwapInterval(1) // enable vsync
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	return win
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

func newProgram(vertexShaderSource io.Reader, fragShaderSource io.Reader) (uint32, error) {

	vbody, err := ioutil.ReadAll(vertexShaderSource)
	if err != nil {
		return 0, err
	}

	fbody, err := ioutil.ReadAll(fragShaderSource)
	if err != nil {
		return 0, err
	}

	vertexShader, err := compileShader(string(vbody), gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(string(fbody), gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	// Create the buffer objects

	gl.BindAttribLocation(program, 0, gl.Str("VertexPosition\x00"))
	gl.BindAttribLocation(program, 1, gl.Str("VertexColor\x00"))

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
    if err := gl.Init(); err != nil {
            panic(err)
    }
    version := gl.GoStr(gl.GetString(gl.VERSION))
    log.Println("OpenGL version", version)

    prog := gl.CreateProgram()
    gl.LinkProgram(prog)
    return prog
}

func main() {
	runtime.LockOSThread()

	// build the main window
	window := initGL(800, 600)
	defer glfw.Terminate()

	filename := "untitled.obj"
	if len(os.Args) == 1 {
		fmt.Printf("No argument, using untitled.obj\n")
	} else {
		filename = os.Args[1]
	}

	// load our OBJ
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	cube, err := obj.NewReader(f).Read()
	if err != nil {
		panic(err)
	}
	var cubeVertices []float32
	var normalVertices []float32

	// convert our object into cube vertices for opengl
	for _, f := range cube.Faces {
		for _, p := range f.Points {

			vx := p.Vertex
			nx := p.Normal

			u := 0.0
			v := 0.0
			if p.Texture != nil {
				u = p.Texture.U
				v = p.Texture.V
			}

			normalVertices = append(normalVertices,
				[]float32{
					float32(nx.Z), float32(nx.Y), float32(nx.X), float32(u), float32(v),
				}...)

			cubeVertices = append(cubeVertices,
				[]float32{
					float32(vx.Z), float32(vx.Y), float32(vx.X), float32(u), float32(v),
				}...)
		}
	}

	// Configure the vertex and fragment shaders
	basicVert, _ := ioutil.ReadFile("basic.vert.glsl")
	basicFrag, _ := ioutil.ReadFile("basic.frag.glsl")

	basicVertStr := strings.TrimSpace(string(basicVert)) + "\x00"
	basicFragStr := strings.TrimSpace(string(basicFrag)) + "\x00"
	program, err := newProgram(bytes.NewBufferString(basicVertStr), bytes.NewBufferString(basicFragStr))
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(800)/600, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projectionMatrix\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{5, 0, 5}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("cameraMatrix\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("modelMatrix\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	modelView := model.Mul4(camera)

	normal := modelView.Inv().Transpose()
	normalUniform := gl.GetUniformLocation(program, gl.Str("normalMatrix\x00"))
	gl.UniformMatrix4fv(normalUniform, 1, false, &normal[0])

	gl.BindFragDataLocation(program, 0, gl.Str("fragColor\x00"))
	// gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("inPosition\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	// Configure the normal data
	var nao uint32
	gl.GenVertexArrays(1, &nao)
	gl.BindVertexArray(nao)

	var nbo uint32
	gl.GenBuffers(1, &nbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, nbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(normalVertices)*4, gl.Ptr(normalVertices), gl.DYNAMIC_DRAW)

	normAttrib := uint32(gl.GetAttribLocation(program, gl.Str("inNormal\x00")))
	gl.EnableVertexAttribArray(normAttrib)
	gl.VertexAttribPointer(normAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	arrayDrawTypeIndex := 0
	arrayDrawTypes := []uint32{
		gl.TRIANGLES,
		gl.TRIANGLES_ADJACENCY,
		gl.TRIANGLE_FAN,
		gl.TRIANGLE_STRIP,
		gl.TRIANGLE_STRIP_ADJACENCY,
		gl.QUADS,
	}

	angle := float32(0.0)
	go func() {
		for {
			select {
			case <-time.After(200 * time.Millisecond):
				angle += 0.1
			}
		}
	}()

	fmt.Printf("cube verticies: %d\n", len(cubeVertices))

	// LOOP:
	for !window.ShouldClose() {
		// Do OpenGL stuff.
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		model = mgl32.Ident4().Mul4(mgl32.HomogRotate3DY(angle))
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		// Render
		gl.UseProgram(program)

		gl.BindVertexArray(vao)
		gl.DrawArrays(arrayDrawTypes[arrayDrawTypeIndex], 0, int32(len(cubeVertices)*4))
		// end OpenGL stuff
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
