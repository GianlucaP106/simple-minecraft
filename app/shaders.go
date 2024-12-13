package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type ShaderManager struct {
	shaders  map[string]uint32
	rootPath string
}

func newShaderManager(root string) *ShaderManager {
	s := &ShaderManager{}
	s.rootPath = root
	s.shaders = make(map[string]uint32)
	return s
}

func (s *ShaderManager) Add(name string) uint32 {
	vshader := filepath.Join(s.rootPath, name, "vert.glsl")
	fshader := filepath.Join(s.rootPath, name, "frag.glsl")
	vb, err := os.ReadFile(vshader)
	if err != nil {
		panic(err)
	}

	fb, err := os.ReadFile(fshader)
	if err != nil {
		panic(err)
	}

	vsrc := string(vb) + "\x00"
	fsrc := string(fb) + "\x00"
	program := s.createProgram(vsrc, fsrc)
	s.shaders[name] = program
	return program
}

func (s *ShaderManager) Program(name string) uint32 {
	return s.shaders[name]
}

func (s *ShaderManager) createProgram(vertexShaderSource, fragmentShaderSource string) uint32 {
	vertexShader, err := s.compile(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := s.compile(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()

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

		panic(fmt.Errorf("failed to link program: %v", log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program
}

func (s *ShaderManager) compile(source string, shaderType uint32) (uint32, error) {
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
