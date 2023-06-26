ANDROID_HOME=/Users/derker/Library/Android/sdk
ANDROID_NDK_BIN=$(ANDROID_HOME)/ndk/25.2.9519653/toolchains/llvm/prebuilt/darwin-x86_64/bin

BUILD_MODULE=./main.go
TARGET_DIR=/Users/derker/Code/android/alist-android/app/src/main/jniLibs
TARGET_NAME=libalist.so

android-armv7a:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=arm \
	GOARM=7 \
	CC=$(ANDROID_NDK_BIN)/armv7a-linux-androideabi21-clang \
	go build -buildmode=c-shared -o $(TARGET_DIR)/armeabi-v7a/${TARGET_NAME} ${BUILD_MODULE}

android-arm64:
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

android: android-armv7a android-arm64 android-x86 android-x86_64