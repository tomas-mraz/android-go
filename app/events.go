package app

// #cgo LDFLAGS: -landroid
//
// #include <android/input.h>
// #include <android/native_activity.h>
// #include <android/native_window.h>
import "C"

import (
	"time"
	"unsafe"

	"github.com/tomas-mraz/android-go/android"
)

type LifecycleEvent struct {
	Activity *android.NativeActivity
	Kind     LifecycleEventKind
}

type LifecycleEventKind string

const (
	OnCreate  LifecycleEventKind = "onCreate"
	OnDestroy LifecycleEventKind = "onDestroy"
	OnStart   LifecycleEventKind = "onStart"
	OnStop    LifecycleEventKind = "onStop"
	OnPause   LifecycleEventKind = "onPause"
	OnResume  LifecycleEventKind = "onResume"
)

//export onCreate
func onCreate(activity *C.ANativeActivity) {
	defaultApp.initWG.Wait()

	defaultApp.mux.Lock()
	defaultApp.activity = android.NewNativeActivityRef(unsafe.Pointer(activity))
	defaultApp.activity.Deref()
	defaultApp.mux.Unlock()

	event := LifecycleEvent{
		Activity: defaultApp.activity,
		Kind:     OnCreate,
	}
	defaultApp.lifecycleEvents <- event
}

//export onDestroy
func onDestroy(activity *C.ANativeActivity) {
	defaultApp.initWG.Wait()

	event := LifecycleEvent{
		Activity: android.NewNativeActivityRef(unsafe.Pointer(activity)),
		Kind:     OnDestroy,
	}
	defaultApp.lifecycleEvents <- event
}

//export onStart
func onStart(activity *C.ANativeActivity) {
	defaultApp.initWG.Wait()

	event := LifecycleEvent{
		Activity: android.NewNativeActivityRef(unsafe.Pointer(activity)),
		Kind:     OnStart,
	}
	defaultApp.lifecycleEvents <- event
}

//export onStop
func onStop(activity *C.ANativeActivity) {
	defaultApp.initWG.Wait()

	event := LifecycleEvent{
		Activity: android.NewNativeActivityRef(unsafe.Pointer(activity)),
		Kind:     OnStop,
	}
	defaultApp.lifecycleEvents <- event
}

//export onPause
func onPause(activity *C.ANativeActivity) {
	defaultApp.initWG.Wait()

	event := LifecycleEvent{
		Activity: android.NewNativeActivityRef(unsafe.Pointer(activity)),
		Kind:     OnPause,
	}
	defaultApp.lifecycleEvents <- event
}

//export onResume
func onResume(activity *C.ANativeActivity) {
	defaultApp.initWG.Wait()

	event := LifecycleEvent{
		Activity: android.NewNativeActivityRef(unsafe.Pointer(activity)),
		Kind:     OnResume,
	}
	defaultApp.lifecycleEvents <- event
}

type SaveStateFunc func(activity *android.NativeActivity, size uintptr) unsafe.Pointer

//export onSaveInstanceState
func onSaveInstanceState(activity *C.ANativeActivity, outSize *C.size_t) unsafe.Pointer {
	defaultApp.initWG.Wait()

	// https://developer.android.com/training/basics/activity-lifecycle/recreating.html
	fn := defaultApp.getSaveInstanceStateFunc()
	if fn == nil {
		return nil
	}
	activityRef := android.NewNativeActivityRef(unsafe.Pointer(activity))
	result := fn(activityRef, uintptr(unsafe.Pointer(outSize)))
	return result
}

type WindowFocusEvent struct {
	Activity *android.NativeActivity
	HasFocus bool
}

//export onWindowFocusChanged
func onWindowFocusChanged(activity *C.ANativeActivity, hasFocus C.int) {
	defaultApp.initWG.Wait()

	out := defaultApp.getWindowFocusEventsOut()
	if out == nil {
		return
	}
	event := WindowFocusEvent{
		Activity: android.NewNativeActivityRef(unsafe.Pointer(activity)),
		HasFocus: hasFocus > 0,
	}
	select {
	case out <- event:
		// dispatched
	case <-time.After(defaultApp.maxDispatchTime):
		// timed out
	}
}

type NativeWindowEvent struct {
	Activity *android.NativeActivity
	Window   *android.NativeWindow
	Kind     NativeWindowEventKind
}

type NativeWindowEventKind string

const (
	NativeWindowCreated      NativeWindowEventKind = "nativeWindowCreated"
	NativeWindowRedrawNeeded NativeWindowEventKind = "nativeWindowRedrawNeeded"
	NativeWindowDestroyed    NativeWindowEventKind = "nativeWindowDestroyed"
)

//export onNativeWindowCreated
func onNativeWindowCreated(activity *C.ANativeActivity, window *C.ANativeWindow) {
	defaultApp.initWG.Wait()

	out := defaultApp.getNativeWindowEventsOut()
	if out == nil {
		return
	}
	event := NativeWindowEvent{
		Activity: android.NewNativeActivityRef(unsafe.Pointer(activity)),
		Window:   (*android.NativeWindow)(window),
		Kind:     NativeWindowCreated,
	}
	select {
	case out <- event:
		// dispatched
	case <-time.After(defaultApp.maxDispatchTime):
		// timed out
	}
}

//export onNativeWindowRedrawNeeded
func onNativeWindowRedrawNeeded(activity *C.ANativeActivity, window *C.ANativeWindow) {
	defaultApp.initWG.Wait()

	out := defaultApp.getNativeWindowEventsOut()
	if out == nil {
		return
	}
	event := NativeWindowEvent{
		Activity: android.NewNativeActivityRef(unsafe.Pointer(activity)),
		Window:   (*android.NativeWindow)(window),
		Kind:     NativeWindowRedrawNeeded,
	}
	select {
	case <-defaultApp.nativeWindowRedrawDone:
	default:
		// skip check
	}
	select {
	case out <- event:
		// dispatched
	case <-time.After(defaultApp.maxDispatchTime):
		return // timed out
	}
	// The drawing window for this native activity needs to be
	// redrawn. To avoid transient artifacts during screen changes
	// (such resizing after rotation), applications should not return
	// from this function until they have finished drawing their
	// window in its current state.
	//
	// Refer to
	// https://developer.android.com/ndk/reference/struct_a_native_activity_callbacks.html
	<-defaultApp.nativeWindowRedrawDone
}

//export onNativeWindowDestroyed
func onNativeWindowDestroyed(activity *C.ANativeActivity, window *C.ANativeWindow) {
	defaultApp.initWG.Wait()

	out := defaultApp.getNativeWindowEventsOut()
	if out == nil {
		return
	}
	event := NativeWindowEvent{
		Activity: android.NewNativeActivityRef(unsafe.Pointer(activity)),
		Window:   (*android.NativeWindow)(window),
		Kind:     NativeWindowDestroyed,
	}
	select {
	case out <- event:
		// dispatched
	case <-time.After(defaultApp.maxDispatchTime):
		// timed out
	}
}

type InputQueueEvent struct {
	Activity *android.NativeActivity
	Queue    *android.InputQueue
	Kind     InputQueueEventKind
}

type InputQueueEventKind string

const (
	QueueCreated   InputQueueEventKind = "queueCreated"
	QueueDestroyed InputQueueEventKind = "queueDestroyed"
)

//export onInputQueueCreated
func onInputQueueCreated(activity *C.ANativeActivity, queue *C.AInputQueue) {
	defaultApp.initWG.Wait()

	out := defaultApp.getInputQueueEventsOut()
	if out == nil {
		return
	}
	event := InputQueueEvent{
		Activity: android.NewNativeActivityRef(unsafe.Pointer(activity)),
		Queue:    (*android.InputQueue)(queue),
		Kind:     QueueCreated,
	}
	select {
	case out <- event:
		// dispatched
	case <-time.After(defaultApp.maxDispatchTime):
		return // timed out
	}

	<-defaultApp.inputQueueHandled
}

//export onInputQueueDestroyed
func onInputQueueDestroyed(activity *C.ANativeActivity, queue *C.AInputQueue) {
	defaultApp.initWG.Wait()

	out := defaultApp.getInputQueueEventsOut()
	if out == nil {
		return
	}
	event := InputQueueEvent{
		Activity: android.NewNativeActivityRef(unsafe.Pointer(activity)),
		Queue:    (*android.InputQueue)(queue),
		Kind:     QueueDestroyed,
	}
	select {
	case out <- event:
		// dispatched
	case <-time.After(defaultApp.maxDispatchTime):
		return // timed out
	}

	<-defaultApp.inputQueueHandled
}

type ContentRectEvent struct {
	Activity *android.NativeActivity
	Rect     *android.Rect
}

//export onContentRectChanged
func onContentRectChanged(activity *C.ANativeActivity, rect *C.ARect) {
	defaultApp.initWG.Wait()

	out := defaultApp.getContentRectEventsOut()
	if out == nil {
		return
	}
	event := ContentRectEvent{
		Activity: android.NewNativeActivityRef(unsafe.Pointer(activity)),
		Rect:     android.NewRectRef(unsafe.Pointer(rect)),
	}
	select {
	case out <- event:
		// dispatched
	case <-time.After(defaultApp.maxDispatchTime):
		// timed out
	}
}

type ActivityEvent struct {
	Activity *android.NativeActivity
	Kind     ActivityEventKind
}

type ActivityEventKind string

const (
	OnConfigurationChanged ActivityEventKind = "onConfigurationChanged"
	OnLowMemory            ActivityEventKind = "onLowMemory"
)

//export onConfigurationChanged
func onConfigurationChanged(activity *C.ANativeActivity) {
	defaultApp.initWG.Wait()

	out := defaultApp.getActivityEventsOut()
	if out == nil {
		return
	}
	event := ActivityEvent{
		Activity: android.NewNativeActivityRef(unsafe.Pointer(activity)),
		Kind:     OnConfigurationChanged,
	}
	select {
	case out <- event:
		// dispatched
	case <-time.After(defaultApp.maxDispatchTime):
		// timed out
	}
}

//export onLowMemory
func onLowMemory(activity *C.ANativeActivity) {
	defaultApp.initWG.Wait()

	out := defaultApp.getActivityEventsOut()
	if out == nil {
		return
	}
	event := ActivityEvent{
		Activity: android.NewNativeActivityRef(unsafe.Pointer(activity)),
		Kind:     OnLowMemory,
	}
	select {
	case out <- event:
		// dispatched
	case <-time.After(defaultApp.maxDispatchTime):
		// timed out
	}
}
