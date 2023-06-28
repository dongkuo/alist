ANDROID_HOME=/d/software/Android/Sdk
ANDROID_NDK_BIN=$(ANDROID_HOME)/ndk/25.2.9519653/toolchains/llvm/prebuilt/windows-x86_64/bin

BUILD_MODULE=./main.go
TARGET_DIR=../../jniLibs
TARGET_NAME=libalist.so

android-armeabi-v7a:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=arm \
	GOARM=7 \
	CC=$(ANDROID_NDK_BIN)/armv7a-linux-androideabi21-clang \
	go build -buildmode=c-shared -o $(TARGET_DIR)/armeabi-v7a/${TARGET_NAME} ${BUILD_MODULE}

android-arm64-v8a:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=arm64 \
	CC=$(ANDROID_NDK_BIN)/aarch64-linux-android21-clang \
	go build -buildmode=c-shared -o $(TARGET_DIR)/arm64-v8a/${TARGET_NAME} ${BUILD_MODULE}

android-x86:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=386 \
	CC=$(ANDROID_NDK_BIN)/i686-linux-android21-clang \
	go build -buildmode=c-shared -o $(TARGET_DIR)/x86/${TARGET_NAME} ${BUILD_MODULE}

android-x86_64:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=amd64 \
	CC=$(ANDROID_NDK_BIN)/x86_64-linux-android21-clang \
	go build -buildmode=c-shared -o $(TARGET_DIR)/x86_64/${TARGET_NAME} ${BUILD_MODULE}

android: android-armeabi-v7a android-arm64-v8a android-x86 android-x86_64
quick: android-arm64-v8a android-x86