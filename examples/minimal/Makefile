# OS switch for cygwin and mingw [linux|windows]
OPERATING_SYSTEM ?= linux
ANDROID_API ?= 26
ANDROID_TOOLCHAIN_DIR ?= $(ANDROID_NDK_HOME)/toolchains/llvm/prebuilt/$(OPERATING_SYSTEM)-x86_64/bin

build:
	mkdir -p android/app/src/main/jniLibs/armeabi-v7a
	CC="$(ANDROID_TOOLCHAIN_DIR)/armv7a-linux-androideabi$(ANDROID_API)-clang" \
	CXX="$(ANDROID_TOOLCHAIN_DIR)/armv7a-linux-androideabi$(ANDROID_API)-clang++" \
	GOOS=android \
	GOARCH=arm \
	CGO_ENABLED=1 \
	go build -buildmode=c-shared -o android/app/src/main/jniLibs/armeabi-v7a/libgomain.so
#	cd android
#	./gradlew build

build64:
	mkdir -p android/app/src/main/jniLibs/arm64-v8a
	CC="$(ANDROID_TOOLCHAIN_DIR)/aarch64-linux-android$(ANDROID_API)-clang" \
	CXX="$(ANDROID_TOOLCHAIN_DIR)/aarch64-linux-android$(ANDROID_API)-clang++" \
	GOOS=android \
	GOARCH=arm64 \
	CGO_ENABLED=1 \
	go build -buildmode=c-shared -o android/app/src/main/jniLibs/arm64-v8a/libgomain.so
#	cd android
#	./gradlew build
