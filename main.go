package main

import (
	"fmt"
	"syscall/js"
)

var wait = make(chan struct{}, 0)
var glo js.Value
var scene js.Value
var camera js.Value
var renderer js.Value
var mesh js.Value
var animator js.Func

func main() {
	fmt.Println("Go/WASM main()")
	glo := js.Global()
	doc := glo.Get("document")
	app := doc.Call("getElementById", "app")
	app.Set("style", "float:left; width:100%; height:100%;")
	cw := app.Get("clientWidth")
	ch := app.Get("clientHeight")

	// camera
	three := glo.Get("THREE")
	cameraCon := three.Get("PerspectiveCamera")
	camera = cameraCon.New(70, cw.Float()/ch.Float(), 0.1, 10)
	cpos := camera.Get("position")
	cpos.Set("x", 0)
	cpos.Set("y", 1.6) //mimic VR pos
	cpos.Set("z", 0)

	// scene
	sceneCon := three.Get("Scene")
	scene = sceneCon.New()

	// bg
	colCon := three.Get("Color")
	col := colCon.New(0x5050a0)
	scene.Set("background", col)

	// cube
	boxCon := three.Get("BoxGeometry")
	box := boxCon.New(0.2, 0.2, 0.2)
	matCon := three.Get("MeshNormalMaterial")
	mat := matCon.New()
	meshCon := three.Get("Mesh")
	mesh = meshCon.New(box, mat)
	meshLoc := mesh.Get("position")
	meshLoc.Set("y", 1.6)
	meshLoc.Set("z", -1)
	scene.Call("add", mesh)

	// room
	lineCon := three.Get("LineSegments")
	boxLineCon := three.Get("BoxLineGeometry")
	matlineCon := three.Get("LineBasicMaterial")
	room := lineCon.New(boxLineCon.New(6, 6, 6, 10, 10, 10), matlineCon.New("{ color: 0x808080 }"))
	roomPos := room.Get("position")
	roomPos.Set("y", 3)
	scene.Call("add", room)

	// renderer
	rendCon := three.Get("WebGLRenderer")
	renderer = rendCon.New("{antialias: true}")
	renderer.Call("setSize", cw.Int(), ch.Int())
	rdom := renderer.Get("domElement")
	app.Call("appendChild", rdom)

	// WebVR
	nav := glo.Get("navigator")
	displayFunc := nav.Get("getVRDisplays") //first see if function exists!

	// animator
	if displayFunc == js.Undefined() {
		fmt.Println("no VR")
		animator = js.FuncOf(animate)
		glo.Set("animate", animator)
		glo.Call("animate")
	} else {
		displays := nav.Call("getVRDisplays")
		fmt.Println("got displays", displays)
		vrAtt := renderer.Get("vr")
		vrAtt.Set("enabled", true)
		animator = js.FuncOf(animateVR)
		renderer.Call("setAnimationLoop", animator)
		body := doc.Call("getElementById", "body")
		webvr := glo.Get("WEBVR")
		button := webvr.Call("createButton", renderer)
		body.Call("appendChild", button)
	}
	defer animator.Release()

	fmt.Println("wait")
	<-wait
	fmt.Println("exit")
}

func rotate() {
	rot := mesh.Get("rotation")
	rx := rot.Get("x")
	ry := rot.Get("y")
	rot.Set("x", rx.Float()+0.01)
	rot.Set("y", ry.Float()+0.02)
}

func animateVR(this js.Value, args []js.Value) interface{} {
	renderer.Call("render", scene, camera)
	rotate()
	return nil
}

func animate(this js.Value, args []js.Value) interface{} {
	js.Global().Call("requestAnimationFrame", animator)
	rotate()
	renderer.Call("render", scene, camera)
	return nil
}
