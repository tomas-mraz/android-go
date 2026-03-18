CPP := $(shell which cpp 2>/dev/null)
#ANDROID_HOME must be set
#ANDROID_NDK_HOME must be set
ANDROID_API ?= 33
ANDROID_TOOLCHAIN_DIR ?= $(ANDROID_NDK_HOME)/toolchains/llvm/prebuilt/linux-x86_64/bin
ANDROID_CC ?= $(ANDROID_TOOLCHAIN_DIR)/aarch64-linux-android$(ANDROID_API)-clang
ANDROID_CXX ?= $(ANDROID_TOOLCHAIN_DIR)/aarch64-linux-android$(ANDROID_API)-clang++
ANDROID_BUILD_ENV = CC="$(ANDROID_CC)" CXX="$(ANDROID_CXX)" GOOS=android GOARCH=arm64 CGO_ENABLED=1

all: gen-android gen-egl gen-gles gen-gles2 gen-gles3 gen-gles31

gen-android:
	test -n "$(CPP)"
	test -x "$(CPP)"
	CPP="$(CPP)" c-for-go -out=. -ccdefs=true android.yml

gen-egl:
	c-for-go -out=. egl.yml

gen-gles:
	c-for-go -out=. gles.yml

gen-gles2:
	c-for-go -out=. gles2.yml

gen-gles3:
	c-for-go -out=. gles3.yml

gen-gles31:
	c-for-go -out=. gles31.yml

clean: clean-egl clean-gles clean-gles2 clean-gles3 clean-gles31
	rm -f android/cgo_helpers.go android/cgo_helpers.h android/cgo_helpers.c
	rm -f android/doc.go android/types.go android/const.go
	rm -f android/android.go

clean-egl:
	rm -f egl/cgo_helpers.go egl/cgo_helpers.h egl/cgo_helpers.c
	rm -f egl/doc.go egl/types.go egl/const.go
	rm -f egl/egl.go

clean-gles:
	rm -f gles/cgo_helpers.go gles/cgo_helpers.h gles/cgo_helpers.c
	rm -f gles/doc.go gles/types.go gles/const.go
	rm -f gles/gles.go

clean-gles2:
	rm -f gles2/cgo_helpers.go gles2/cgo_helpers.h gles2/cgo_helpers.c
	rm -f gles2/doc.go gles2/types.go gles2/const.go
	rm -f gles2/gles2.go

clean-gles3:
	rm -f gles3/cgo_helpers.go gles3/cgo_helpers.h gles3/cgo_helpers.c
	rm -f gles3/doc.go gles3/types.go gles3/const.go
	rm -f gles3/gles3.go

clean-gles31:
	rm -f gles31/cgo_helpers.go gles31/cgo_helpers.h gles31/cgo_helpers.c
	rm -f gles31/doc.go gles31/types.go gles31/const.go
	rm -f gles31/gles31.go

check-android-build-env:
	test -n "$(ANDROID_NDK_HOME)"
	test -x "$(ANDROID_CC)"
	test -x "$(ANDROID_CXX)"

test: check-android-build-env
	cd android && $(ANDROID_BUILD_ENV) go build

test-egl: check-android-build-env
	cd egl && $(ANDROID_BUILD_ENV) go build

test-gles:
	cd gles && go build

test-gles2:
	cd gles2 && go build

test-gles3:
	cd gles3 && go build

test-gles31:
	cd gles31 && go build
