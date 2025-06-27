package main

import (
	"log"
	"runtime"
	"time"

	"github.com/tomas-mraz/android-go/android"
	"github.com/tomas-mraz/android-go/app"
	"github.com/tomas-mraz/android-go/egl"
	gl "github.com/tomas-mraz/android-go/gles"
)

func init() {
	app.SetLogTag("EGLActivity")
}

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	nativeWindowEvents := make(chan app.NativeWindowEvent, 1)
	windowFocusEvents := make(chan app.WindowFocusEvent, 1)
	inputQueueEvents := make(chan app.InputQueueEvent, 1)
	inputQueueChan := make(chan *android.InputQueue, 1)
	var displayHandle *egl.DisplayHandle
	var windowFocused bool

	type vec3 struct {
		X, Y, Z float32
	}
	sensorEvents := make(chan vec3, 100)
	sensorMan := NewSensorMan(20*time.Millisecond, func(event *android.SensorEvent) {
		x, y, z := event.Acceleration()
		select {
		case sensorEvents <- vec3{x, y, z}:
		default:
		}
	})

	stateX := float32(0.5)
	stateY := float32(0.5)
	stateZ := float32(0.5)
	stateRiseY := true

	app.Main(func(a app.NativeActivity) {
		a.HandleNativeWindowEvents(nativeWindowEvents)
		a.HandleWindowFocusEvents(windowFocusEvents)
		a.HandleInputQueueEvents(inputQueueEvents)
		go app.HandleInputQueues(inputQueueChan, func() {
			a.InputQueueHandled()
		}, app.LogInputEvents)
		a.InitDone()
		for {
			select {
			case vec := <-sensorEvents:
				// log.Printf("accelerometer x=%0.3f y=%0.3f z=%0.3f", vec.X, vec.Y, vec.Z)
				stateX = 0.5 + vec.X/10.0
				stateZ = 0.5 + vec.Y/10.0
				if stateRiseY {
					stateY += 0.01
					if stateY >= 1 {
						stateRiseY = false
					}
				} else {
					stateY -= 0.01
					if stateY <= 0 {
						stateRiseY = true
					}
				}
				draw(displayHandle, stateX, stateY, stateZ)
			case <-a.LifecycleEvents():
			case event := <-windowFocusEvents:
				if event.HasFocus && !windowFocused {
					windowFocused = true
					sensorMan.Start()
				}
				if !event.HasFocus && windowFocused {
					windowFocused = false
					sensorMan.Stop()
				}
				draw(displayHandle, stateX, stateY, stateZ)
			case event := <-inputQueueEvents:
				switch event.Kind {
				case app.QueueCreated:
					inputQueueChan <- event.Queue
				case app.QueueDestroyed:
					inputQueueChan <- nil
				}
			case event := <-nativeWindowEvents:
				switch event.Kind {
				case app.NativeWindowRedrawNeeded:
					draw(displayHandle, stateX, stateY, stateZ)
					a.NativeWindowRedrawDone()
				case app.NativeWindowCreated:
					expectedSurface := map[int32]int32{
						egl.SurfaceType: egl.WindowBit,
						egl.RedSize:     8,
						egl.GreenSize:   8,
						egl.BlueSize:    8,
					}
					if handle, err := egl.NewDisplayHandle(event.Window, expectedSurface); err != nil {
						log.Fatalln("EGL error:", err)
					} else {
						displayHandle = handle
						log.Printf("EGL display res: %dx%d", handle.Width, handle.Height)
					}
					initGL()
				case app.NativeWindowDestroyed:
					displayHandle.Destroy()
				}
			}
		}
	})
}

func initGL() {
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.FASTEST)
	gl.Enable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.ShadeModel(gl.SMOOTH)
}

func draw(handle *egl.DisplayHandle, x, y, z float32) {
	gl.ClearColor(x, y, z, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	handle.SwapBuffers()
}
