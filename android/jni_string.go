package android

/*
#include <jni.h>
#include "jni_call.h"
*/
import "C"

import (
	"unsafe"
)

// JNIEnvGetStringUTFChars copies the JNI UTF chars into Go-managed memory and
// releases the JNI buffer before returning.
func JNIEnvGetStringUTFChars(env *JNIEnv, str Jstring) string {
	cEnv := (*C.JNIEnv)(unsafe.Pointer(env))
	cStr := (C.jstring)(unsafe.Pointer(str))
	ret := C.JNIEnv_GetStringUTFChars(cEnv, cStr, nil)
	if ret == nil {
		return ""
	}
	defer C.JNIEnv_ReleaseStringUTFChars(cEnv, cStr, ret)

	return C.GoString(ret)
}
