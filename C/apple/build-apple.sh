#!/bin/bash

# edit these version numbers to suit your needs, or define them before running the script

echo "BUILD_TARGETS environment variable can be set as a string split by ':' as you would a PATH variable. Ditto LINK_TARGETS"
# example: 
#   export BUILD_TARGETS="simulator_x86_64:catalyst_x86_64:macos_x86_64:ios-arm64e"

IFS=':' read -r -a sqcloud_build_targets <<< "$BUILD_TARGETS"
IFS=':' read -r -a sqcloud_link_targets <<< "$LINK_TARGETS"

if [ -z "$IOS" ]
then
  IOS=`xcrun -sdk iphoneos --show-sdk-version`
fi

if [ -z "$MIN_IOS_VERSION" ]
then
  MIN_IOS_VERSION=13.0
fi

if [ -z "$SQCLOUD_VERSION" ]
then
  SQCLOUD_VERSION=1.0.0
fi

if [ -z "$MACOSX" ]
then
  MACOSX=`xcrun --sdk macosx --show-sdk-version|cut -d '.' -f 1-2`
fi

declare -a all_targets=("ios-arm64" "ios-arm64e" "simulator_x86_64" "simulator_x86_64h" "simulator_arm64e" "simulator_arm64" "catalyst_x86_64" "catalyst_arm64" "macos_x86_64" "macos_x86_64h" "macos_arm64")
declare -a old_targets=("simulator_x86_64" "catalyst_x86_64" "macos_x86_64" "ios-arm64")
declare -a appleSiliconTargets=("simulator_arm64" "simulator_x86_64" "catalyst_x86_64" "catalyst_arm64" "macos_arm64" "macos_x86_64" "ios-arm64")

if [ -z "$sqcloud_build_targets" ]
then
  declare -a sqcloud_build_targets=("simulator_x86_64" "simulator_arm64" "catalyst_x86_64" "catalyst_arm64" "macos_x86_64" "macos_arm64" "ios-arm64")
  #declare -a sqcloud_build_targets=("simulator_x86_64" "catalyst_x86_64" "catalyst_arm64" "macos_x86_64" "macos_arm64" "ios-arm64")
fi

if [ -z "$sqcloud_link_targets" ]
then
  declare -a sqcloud_link_targets=("simulator_x86_64" "simulator_arm64" "catalyst_x86_64" "catalyst_arm64" "macos_x86_64" "macos_arm64" "ios-arm64")
  #declare -a sqcloud_link_targets=("simulator_x86_64" "catalyst_x86_64" "catalyst_arm64" "macos_x86_64" "macos_arm64" "ios-arm64")
fi

set -e

XCODE=`/usr/bin/xcode-select -p`

# create a staging directory (we need this for include files later on)
#PREFIX=$(pwd)/build/sqcloud-build    # this is where we build sqcloud
OUTPUT=$(pwd)/fat/sqcloud            # after we build, we put sqcloud outputs here
XCFRAMEWORKS=$(pwd)/output/      # this is where we produce the resulting XCFrameworks: libcrypto.xcframework and libssl.xcframework
LIBTLS_FRAMEWORK=$(pwd)/output/libtls.xcframework

#mkdir -p $PREFIX
mkdir -p $OUTPUT
mkdir -p $XCFRAMEWORKS

for target in "${sqcloud_build_targets[@]}"
do
  #mkdir -p $PREFIX/$target;
  mkdir -p $OUTPUT/$target/lib;
  mkdir -p $OUTPUT/$target/bin;
  mkdir -p $OUTPUT/$target/include;
done

# build libtls, if needed
if [ ! -e $LIBTLS_FRAMEWORK ]
then
  ./libressl.sh
fi

# some bash'isms
elementIn () { # source https://stackoverflow.com/questions/3685970/check-if-a-bash-array-contains-a-value
  local e match="$1"
  shift
  for e; do [[ "$e" == "$match" ]] && return 0; done
  return 1
}

makeSQCloud() {
  # only build the files we need (libsqcloud, include files)
  CC=$CC CFLAGS=$CFLAGS LIBTLS=$LIBTLS TARGET=$target make clean libsqcloud.a cli
}

moveSQCloudOutputInPlace() {
  local target=$1
  local output=$2
  echo "cp libsqcloud.a $OUTPUT/$target/lib"
  cp libsqcloud.a $OUTPUT/$target/lib
  cp sqlitecloud-cli $OUTPUT/$target/bin
  rsync -am --include='sqcloud.h' --exclude="" -f 'hide,! */' ../*.h $OUTPUT/$target/include
}

needsRebuilding() {
  local target=$1
  test libsqcloud.a -nt Makefile
  timestampCompare=$?
  if [ $timestampCompare -eq 1 ]; then
    return 0
  else
    arch=`/usr/bin/lipo -archs libsqcloud.a`
    if [ "$arch" == "$target" ]; then
      return 1
    else
      return 0
    fi
  fi
}

##############################################
##  iOS Simulator x86_64h Compilation
##############################################

target=simulator_x86_64h
if needsRebuilding "$target" && elementIn "$target" "${sqcloud_build_targets[@]}"; then


  printf "\n\n--> iOS Simulator x86_64h Compilation"

  DEVROOT=$XCODE/Platforms/iPhoneSimulator.platform/Developer
  SDKROOT=$DEVROOT/SDKs/iPhoneSimulator${IOS}.sdk
  CC="/usr/bin/clang"
  CPPFLAGS="-I$SDKROOT/usr/include/"
  CFLAGS="$CPPFLAGS -arch x86_64h -miphoneos-version-min=${MIN_IOS_VERSION} -isysroot $SDKROOT"
  
#  ./configure --host=x86_64-apple-darwin --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang" \
    CPPFLAGS="-I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch x86_64h -miphoneos-version-min=${MIN_IOS_VERSION} -pipe -no-cpp-precomp -isysroot $SDKROOT" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD=$DEVROOT/usr/bin/ld

  makeSQCloud
  printf "\n\n--> XX iOS Simulator x86_64h libsqcloud Compilation\n"
  moveSQCloudOutputInPlace $target $OUTPUT

fi;

#############################################
##  iOS Simulator x86_64 libsqcloud Compilation
#############################################

target=simulator_x86_64
if needsRebuilding "$target" && elementIn "$target" "${sqcloud_build_targets[@]}"; then

  printf "\n\n--> iOS Simulator x86_64 libsqcloud Compilation\n"

  DEVROOT=$XCODE/Platforms/iPhoneSimulator.platform/Developer
  SDKROOT=$DEVROOT/SDKs/iPhoneSimulator${IOS}.sdk
  CC="/usr/bin/clang"
  CPPFLAGS="-I$SDKROOT/usr/include/"
  CFLAGS="$CPPFLAGS -arch x86_64 -miphoneos-version-min=${MIN_IOS_VERSION} -pipe -no-cpp-precomp -isysroot $SDKROOT"

  #echo "prefix: $PREFIX/$target"
  echo "SDKROOT: $SDKROOT"
  echo "CPPFLAGS: $CPPFLAGS"
  echo "IOS: $IOS"

  #./configure --host=x86_64-apple-darwin --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang" \
    CPPFLAGS="-I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch x86_64 -miphoneos-version-min=${MIN_IOS_VERSION} -pipe -no-cpp-precomp -isysroot $SDKROOT" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD=$DEVROOT/usr/bin/ld

  makeSQCloud
  echo -printf "\n\n--> XX iOS Simulator x86_64 libsqcloud Compilation\n"
  moveSQCloudOutputInPlace $target $OUTPUT

fi;

#############################################
##  iOS Simulator arm64e libsqcloud Compilation
#############################################

target=simulator_arm64e
if needsRebuilding "$target" && elementIn "$target" "${sqcloud_build_targets[@]}"; then

  printf "\n\n--> iOS Simulator arm64e libsqcloud Compilation\n"

  DEVROOT=$XCODE/Platforms/iPhoneSimulator.platform/Developer
  SDKROOT=$DEVROOT/SDKs/iPhoneSimulator${IOS}.sdk
  CC="/usr/bin/clang -target arm64-apple-ios${IOS}-simulator"
  CPPFLAGS="-I$SDKROOT/usr/include/ -target arm64-apple-ios${IOS}-simulator"
  CFLAGS="$CPPFLAGS -arch arm64e -miphoneos-version-min=${MIN_IOS_VERSION} -pipe -no-cpp-precomp -isysroot $SDKROOT"
  
  #./configure --build=aarch64-apple-darwin --host=aarch64-apple-darwin22 --prefix="$PREFIX/$target" \
  ./configure --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -target arm64-apple-ios${IOS}-simulator" \
    CPPFLAGS="-I$SDKROOT/usr/include/ -target arm64-apple-ios${IOS}-simulator" \
    CFLAGS="$CPPFLAGS -arch arm64e -miphoneos-version-min=${MIN_IOS_VERSION} -pipe -no-cpp-precomp -isysroot $SDKROOT" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD=$DEVROOT/usr/bin/ld

  makeSQCloud
  printf "\n\n--> XX iOS Simulator arm64e libsqcloud Compilation\n"
  moveSQCloudOutputInPlace $target $OUTPUT

fi;

#############################################
##  iOS Simulator arm64 libsqcloud Compilation
#############################################

target=simulator_arm64
if needsRebuilding "$target" && elementIn "$target" "${sqcloud_build_targets[@]}"; then

  printf "\n\n--> iOS Simulator arm64 libsqcloud Compilation\n"

  DEVROOT=$XCODE/Platforms/iPhoneSimulator.platform/Developer
  SDKROOT=$DEVROOT/SDKs/iPhoneSimulator${IOS}.sdk
  CC="/usr/bin/clang -target arm64-apple-ios${IOS}-simulator"
  CPPFLAGS="-I$SDKROOT/usr/include/"
  CFLAGS="$CPPFLAGS -arch arm64 -miphoneos-version-min=${MIN_IOS_VERSION} -pipe -no-cpp-precomp -isysroot $SDKROOT"
  
  #./configure --host=aarch64-apple-darwin --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -target arm64-apple-ios${IOS}-simulator" \
    CPPFLAGS="-I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch arm64 -miphoneos-version-min=${MIN_IOS_VERSION} -pipe -no-cpp-precomp -isysroot $SDKROOT" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD=$DEVROOT/usr/bin/ld

  makeSQCloud
  printf "\n\n--> XX iOS Simulator arm64 libsqcloud Compilation\n"
  moveSQCloudOutputInPlace $target $OUTPUT

fi;


##################################
##  iOS arm64 libsqcloud Compilation
##################################

target=ios-arm64
if needsRebuilding "$target" && elementIn "$target" "${sqcloud_build_targets[@]}"; then

  printf "\n\n--> iOS arm64 libsqcloud Compilation\n"

  DEVROOT=$XCODE/Platforms/iPhoneOS.platform/Developer
  SDKROOT=$DEVROOT/SDKs/iPhoneOS${IOS}.sdk
  CC="/usr/bin/clang -isysroot $SDKROOT" \
  CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
  CFLAGS="$CPPFLAGS -arch arm64 -miphoneos-version-min=${MIN_IOS_VERSION} -pipe -no-cpp-precomp" \
  
  #./configure --host=aarch64-apple-darwin --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -isysroot $SDKROOT" \
    CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch arm64 -miphoneos-version-min=${MIN_IOS_VERSION} -pipe -no-cpp-precomp" \
    CPP="/usr/bin/cpp -D__arm__=1 $CPPFLAGS" \
    LD=$DEVROOT/usr/bin/ld

  makeSQCloud
  printf "\n\n--> XX iOS arm64 libsqcloud Compilation\n"
  moveSQCloudOutputInPlace $target $OUTPUT

fi;

###################################
##  iOS arm64e libsqcloud Compilation
###################################

target=ios-arm64e
if elementIn "$target" "${sqcloud_build_targets[@]}"; then

  printf "\n\n--> iOS arm64e libsqcloud Compilation\n"

  DEVROOT=$XCODE/Platforms/iPhoneOS.platform/Developer
  SDKROOT=$DEVROOT/SDKs/iPhoneOS${IOS}.sdk
  CC="/usr/bin/clang -isysroot $SDKROOT" \
  CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
  CFLAGS="$CPPFLAGS -arch arm64e -miphoneos-version-min=${MIN_IOS_VERSION} -pipe -no-cpp-precomp" \
  
  #./configure --host=aarch64-apple-darwin19 --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -isysroot $SDKROOT" \
    CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch arm64e -miphoneos-version-min=${MIN_IOS_VERSION} -pipe -no-cpp-precomp" \
    CPP="/usr/bin/cpp -D__arm__=1 $CPPFLAGS" \
    LD=$DEVROOT/usr/bin/ld

  makeSQCloud
  printf "\n\n--> XX iOS arm64e libsqcloud Compilation\n"
  moveSQCloudOutputInPlace $target $OUTPUT

fi;

##############################################
##  macOS Catalyst x86_64 libsqcloud Compilation
##############################################

target=catalyst_x86_64
if needsRebuilding "$target" && elementIn "$target" "${sqcloud_build_targets[@]}"; then

  printf "\n\n--> macOS Catalyst x86_64 libsqcloud Compilation\n"

  DEVROOT=$XCODE/Platforms/MacOSX.platform/Developer
  SDKROOT=$DEVROOT/SDKs/MacOSX${MACOSX}.sdk
  CC="/usr/bin/clang -target x86_64-apple-ios${IOS}-macabi -isysroot $SDKROOT" \
  CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
  CFLAGS="$CPPFLAGS -arch x86_64 -pipe -no-cpp-precomp" \
  
  #./configure --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -target x86_64-apple-ios${IOS}-macabi -isysroot $SDKROOT" \
    CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch x86_64 -pipe -no-cpp-precomp" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD="/usr/bin/ld"

  makeSQCloud
  printf "\n\n--> XX macOS Catalyst x86_64 libsqcloud Compilation\n"
  moveSQCloudOutputInPlace $target $OUTPUT

fi;

#############################################
##  macOS Catalyst arm64 libsqcloud Compilation
#############################################

target=catalyst_arm64
if needsRebuilding "$target" && elementIn "$target" "${sqcloud_build_targets[@]}"; then

  printf "\n\n--> macOS Catalyst arm64 libsqcloud Compilation\n"

  DEVROOT=$XCODE/Platforms/MacOSX.platform/Developer
  SDKROOT=$DEVROOT/SDKs/MacOSX${MACOSX}.sdk
  CC="/usr/bin/clang -target arm64-apple-ios${IOS}-macabi -isysroot $SDKROOT" \
  CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
  CFLAGS="$CPPFLAGS -arch arm64 -pipe -no-cpp-precomp" \
  
  #./configure --build=aarch64-apple-darwin --host=aarch64-apple-darwin22 --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -target arm64-apple-ios${IOS}-macabi -isysroot $SDKROOT" \
    CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch arm64 -pipe -no-cpp-precomp" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD="/usr/bin/ld"

  makeSQCloud
  printf "\n\n--> XX macOS Catalyst arm64 libsqcloud Compilation\n"
  moveSQCloudOutputInPlace $target $OUTPUT

fi;


#####################################
##  macOS x86_64 libsqcloud Compilation
#####################################

target=macos_x86_64
if needsRebuilding "$target" && elementIn "$target" "${sqcloud_build_targets[@]}"; then

  printf "\n\n--> macOS x86_64 libsqcloud Compilation\n"

  DEVROOT=$XCODE/Platforms/MacOSX.platform/Developer
  SDKROOT=$DEVROOT/SDKs/MacOSX${MACOSX}.sdk
  CC="/usr/bin/clang -target x86_64-apple-darwin -isysroot $SDKROOT" \
  CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
  CFLAGS="$CPPFLAGS -arch x86_64 -pipe -no-cpp-precomp" \
  
  #./configure --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -target x86_64-apple-darwin -isysroot $SDKROOT" \
    CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch x86_64 -pipe -no-cpp-precomp" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD="/usr/bin/clang -target x86_64-apple-darwin"

  makeSQCloud
  printf "\n\n--> XX macOS x86_64 libsqcloud Compilation\n"
  moveSQCloudOutputInPlace $target $OUTPUT

fi;


######################################
##  macOS x86_64h libsqcloud Compilation
######################################

target=macos_x86_64h
if needsRebuilding "$target" && elementIn "$target" "${sqcloud_build_targets[@]}"; then

  printf "\n\n--> macOS x86_64h libsqcloud Compilation\n"

  DEVROOT=$XCODE/Platforms/MacOSX.platform/Developer
  SDKROOT=$DEVROOT/SDKs/MacOSX${MACOSX}.sdk
  CC="/usr/bin/clang -isysroot $SDKROOT" \
  CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
  CFLAGS="$CPPFLAGS -arch x86_64h -pipe -no-cpp-precomp" \
  
  #./configure --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -isysroot $SDKROOT" \
    CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch x86_64h -pipe -no-cpp-precomp" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD="/usr/bin/ld"

  makeSQCloud
  printf "\n\n--> XX macOS x86_64h libsqcloud Compilation\n"
  moveSQCloudOutputInPlace $target $OUTPUT

fi;

#####################################
##  macOS arm64 libsqcloud Compilation
#####################################

target=macos_arm64
if needsRebuilding "$target" && elementIn "$target" "${sqcloud_build_targets[@]}"; then

  printf "\n\n--> macOS arm64 libsqcloud Compilation\n"

  DEVROOT=$XCODE/Platforms/MacOSX.platform/Developer
  SDKROOT=$DEVROOT/SDKs/MacOSX${MACOSX}.sdk
  CC="/usr/bin/clang -isysroot $SDKROOT" \
  CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
  CFLAGS="$CPPFLAGS -arch arm64 -pipe -no-cpp-precomp" \
  
  #./configure --host=arm-apple-darwin --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -isysroot $SDKROOT" \
    CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch arm64 -pipe -no-cpp-precomp" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD="/usr/bin/ld"

  makeSQCloud
  printf "\n\n--> XX macOS arm64 libsqcloud Compilation\n"
  moveSQCloudOutputInPlace $target $OUTPUT

fi;

# TODO: This one isn't working - use "host" to cross compile
#####################################
##  macOS arm64e libsqcloud Compilation
#####################################

target=macos_arm64e
if needsRebuilding "$target" && elementIn "$target" "${sqcloud_build_targets[@]}"; then

  printf "\n\n--> macOS arm64e libsqcloud Compilation\n"

  DEVROOT=$XCODE/Platforms/MacOSX.platform/Developer
  SDKROOT=$DEVROOT/SDKs/MacOSX${MACOSX}.sdk
  CC="/usr/bin/clang -target arm64-apple-darwin -isysroot $SDKROOT" \
  CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
  CFLAGS="$CPPFLAGS -arch arm64e -pipe -no-cpp-precomp" \
  
  #./configure --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -target arm64-apple-darwin -isysroot $SDKROOT" \
    CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch arm64e -pipe -no-cpp-precomp" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD="/usr/bin/ld -target arm64-apple-darwin"

  makeSQCloud
  printf "\n\n--> XX macOS arm64e libsqcloud Compilation\n"
  moveSQCloudOutputInPlace $target $OUTPUT

fi;


####################################
## lipo & XCFrameworks for SQCloud
####################################

## lipo & XCFramework

macos=()
catalyst=()
simulator=()
ios=()


for target in "${sqcloud_link_targets[@]}"
do
  if [[ $target == "ios-"* ]]; then
    ios+=($target)
  fi
  if [[ $target == "simulator_"* ]]; then
    simulator+=($target)
  fi
  if [[ $target == "catalyst_"* ]]; then
    catalyst+=($target)
  fi
  if [[ $target == "macos_"* ]]; then
    macos+=($target)
  fi
done

XCFRAMEWORK_LIBSQCLOUD_CMD="xcodebuild -create-xcframework"

framework_targets=()

if [ ${#ios[@]} -gt 0 ]; then
  lipo_libsqcloud="lipo -create "

  framework_targets+=("ios")
  mkdir -p $OUTPUT/ios/lib
  mkdir -p $OUTPUT/ios/include

  for target in "${ios[@]}"
  do
      lipo_libsqcloud="$lipo_libsqcloud $OUTPUT/$target/lib/libsqcloud.a"
    rsync -a $OUTPUT/$target/include/* $OUTPUT/ios/include
  done

  lipo_libsqcloud="$lipo_libsqcloud -output $OUTPUT/ios/lib/libsqcloud.a"
  echo $lipo_libsqcloud
  eval $lipo_libsqcloud


  XCFRAMEWORK_LIBSQCLOUD_CMD="$XCFRAMEWORK_LIBSQCLOUD_CMD -library $OUTPUT/ios/lib/libsqcloud.a"
  XCFRAMEWORK_LIBSQCLOUD_CMD="$XCFRAMEWORK_LIBSQCLOUD_CMD -headers $OUTPUT/ios/include"
fi

if [ ${#catalyst[@]} -gt 0 ]; then
  lipo_libsqcloud="lipo -create "

  framework_targets+=("catalyst")
  mkdir -p $OUTPUT/catalyst/lib
  mkdir -p $OUTPUT/catalyst/include

  for target in "${catalyst[@]}"
  do
    lipo_libsqcloud="$lipo_libsqcloud $OUTPUT/$target/lib/libsqcloud.a"
    rsync -a $OUTPUT/$target/include/* $OUTPUT/catalyst/include
  done

  lipo_libsqcloud="$lipo_libsqcloud -output $OUTPUT/catalyst/lib/libsqcloud.a"
  echo $lipo_libsqcloud
  eval $lipo_libsqcloud

  XCFRAMEWORK_LIBSQCLOUD_CMD="$XCFRAMEWORK_LIBSQCLOUD_CMD -library $OUTPUT/catalyst/lib/libsqcloud.a"
  XCFRAMEWORK_LIBSQCLOUD_CMD="$XCFRAMEWORK_LIBSQCLOUD_CMD -headers $OUTPUT/catalyst/include"
fi

if [ ${#macos[@]} -gt 0 ]; then
  lipo_libsqcloud="lipo -create "

  framework_targets+=("macos")
  mkdir -p $OUTPUT/macos/lib
  mkdir -p $OUTPUT/macos/include

  for target in "${macos[@]}"
  do
    lipo_libsqcloud="$lipo_libsqcloud $OUTPUT/$target/lib/libsqcloud.a"
    rsync -a $OUTPUT/$target/include/* $OUTPUT/macos/include
  done

  lipo_libsqcloud="$lipo_libsqcloud -output $OUTPUT/macos/lib/libsqcloud.a"
  echo $lipo_libsqcloud
  eval $lipo_libsqcloud

  XCFRAMEWORK_LIBSQCLOUD_CMD="$XCFRAMEWORK_LIBSQCLOUD_CMD -library $OUTPUT/macos/lib/libsqcloud.a"
  XCFRAMEWORK_LIBSQCLOUD_CMD="$XCFRAMEWORK_LIBSQCLOUD_CMD -headers $OUTPUT/macos/include"
fi

if [ ${#simulator[@]} -gt 0 ]; then
  lipo_libsqcloud="lipo -create "

  framework_targets+=("simulator")
  mkdir -p $OUTPUT/simulator/lib
  mkdir -p $OUTPUT/simulator/include

  for target in "${simulator[@]}"
  do
    lipo_libsqcloud="$lipo_libsqcloud $OUTPUT/$target/lib/libsqcloud.a"
    rsync -a $OUTPUT/$target/include/* $OUTPUT/simulator/include
  done

  lipo_libsqcloud="$lipo_libsqcloud -output $OUTPUT/simulator/lib/libsqcloud.a"
  echo $lipo_libsqcloud
  eval $lipo_libsqcloud

  XCFRAMEWORK_LIBSQCLOUD_CMD="$XCFRAMEWORK_LIBSQCLOUD_CMD -library $OUTPUT/simulator/lib/libsqcloud.a"
  XCFRAMEWORK_LIBSQCLOUD_CMD="$XCFRAMEWORK_LIBSQCLOUD_CMD -headers $OUTPUT/simulator/include"
fi

XCFRAMEWORK_LIBSQCLOUD_CMD="$XCFRAMEWORK_LIBSQCLOUD_CMD -output $XCFRAMEWORKS/libsqcloud.xcframework"

echo $XCFRAMEWORK_LIBSQCLOUD_CMD
eval $XCFRAMEWORK_LIBSQCLOUD_CMD

cd ..
