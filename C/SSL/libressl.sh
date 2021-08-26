#!/bin/bash

# edit these version numbers to suit your needs, or define them before running the script

echo "BUILD_TARGETS environment variable can be set as a string split by ':' as you would a PATH variable. Ditto LINK_TARGETS"
# example: 
#   export BUILD_TARGETS="simulator_x86_64:catalyst_x86_64:macos_x86_64:ios-arm64e"

IFS=':' read -r -a libressl_build_targets <<< "$BUILD_TARGETS"
IFS=':' read -r -a libressl_link_targets <<< "$LINK_TARGETS"

if [ -z "$IOS" ]
then
  IOS=`xcrun -sdk iphoneos --show-sdk-version`
fi

if [ -z "$MIN_IOS_VERSION" ]
then
  MIN_IOS_VERSION=13.0
fi

if [ -z "$LIBRESSL" ]
then
  LIBRESSL=3.3.3
fi

if [ -z "$MACOSX" ]
then
  MACOSX=`xcrun --sdk macosx --show-sdk-version|cut -d '.' -f 1-2`
fi

declare -a all_targets=("ios-arm64" "ios-arm64e" "simulator_x86_64" "simulator_x86_64h" "simulator_arm64e" "simulator_arm64" "catalyst_x86_64" "catalyst_arm64" "macos_x86_64" "macos_x86_64h" "macos_arm64")
declare -a old_targets=("simulator_x86_64" "catalyst_x86_64" "macos_x86_64" "ios-arm64")
declare -a appleSiliconTargets=("simulator_arm64" "simulator_x86_64" "catalyst_x86_64" "catalyst_arm64" "macos_arm64" "macos_x86_64" "ios-arm64")

if [ -z "$libressl_build_targets" ]
then
  #declare -a libressl_build_targets=("simulator_x86_64" "simulator_arm64" "catalyst_x86_64" "catalyst_arm64" "macos_x86_64" "macos_arm64" "ios-arm64")
  declare -a libressl_build_targets=("simulator_x86_64" "catalyst_x86_64" "catalyst_arm64" "macos_x86_64" "macos_arm64" "ios-arm64")
fi

if [ -z "$libressl_link_targets" ]
then
  #declare -a libressl_link_targets=("simulator_x86_64" "simulator_arm64" "catalyst_x86_64" "catalyst_arm64" "macos_x86_64" "macos_arm64" "ios-arm64")
  declare -a libressl_link_targets=("simulator_x86_64" "catalyst_x86_64" "catalyst_arm64" "macos_x86_64" "macos_arm64" "ios-arm64")
fi

set -e

XCODE=`/usr/bin/xcode-select -p`

# download LibreSSL
if [ ! -e "libressl-$LIBRESSL.tar.gz" ]
then
    curl -OL "https://ftp.openbsd.org/pub/OpenBSD/LibreSSL/libressl-${LIBRESSL}.tar.gz"
    tar -zxf "libressl-${LIBRESSL}.tar.gz"
fi

# create a staging directory (we need this for include files later on)
PREFIX=$(pwd)/build/libressl-build    # this is where we build libressl
OUTPUT=$(pwd)/Fat/libressl            # after we build, we put libressls outputs here
XCFRAMEWORKS=$(pwd)/output/           # this is where we produce the resulting XCFrameworks: libcrypto.xcframework and libssl.xcframework

mkdir -p $PREFIX
mkdir -p $OUTPUT
mkdir -p $XCFRAMEWORKS

for target in "${libressl_build_targets[@]}"
do
  mkdir -p $PREFIX/$target;
  mkdir -p $OUTPUT/$target/lib;
  mkdir -p $OUTPUT/$target/include;
done

cd libressl-${LIBRESSL}

# this cleans everything out of the build directory so we can have a clean build
if [ -e "./Makefile" ]
then
  # since we clean before we build, do we still need this??
    make distclean
fi

# some bash'isms
elementIn () { # source https://stackoverflow.com/questions/3685970/check-if-a-bash-array-contains-a-value
  local e match="$1"
  shift
  for e; do [[ "$e" == "$match" ]] && return 0; done
  return 1
}

makeLibreSSL() {
  # only build the files we need (libcrypto, libssl, include files)
  make -C crypto clean all install
  make -C ssl clean all install
  make -C include install
}

moveLibreSSLOutputInPlace() {
  local target=$1
  local output=$2
  cp crypto/.libs/libcrypto.a $OUTPUT/$target/lib
  cp ssl/.libs/libssl.a $OUTPUT/$target/lib
  rsync -am --include='*.h' -f 'hide,! */' include/* $OUTPUT/$target/include
}

needsRebuilding() {
  local target=$1
  test crypto/.libs/libcrypto.a -nt Makefile
  timestampCompare=$?
  if [ $timestampCompare -eq 1 ]; then
    return 0
  else
    arch=`/usr/bin/lipo -archs crypto/.libs/libcrypto.a`
    if [ "$arch" == "$target" ]; then
      return 1
    else
      return 0
    fi
  fi
}

##############################################
##  iOS Simulator x86_64h libssl Compilation
##############################################

target=simulator_x86_64h
if needsRebuilding "$target" && elementIn "$target" "${libressl_build_targets[@]}"; then


  printf "\n\n--> iOS Simulator x86_64h libssl Compilation"

  DEVROOT=$XCODE/Platforms/iPhoneSimulator.platform/Developer
  SDKROOT=$DEVROOT/SDKs/iPhoneSimulator${IOS}.sdk

  ./configure --host=x86_64-apple-darwin --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang" \
    CPPFLAGS="-I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch x86_64h -miphoneos-version-min=${MIN_IOS_VERSION} -pipe -no-cpp-precomp -isysroot $SDKROOT" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD=$DEVROOT/usr/bin/ld

  makeLibreSSL
  printf "\n\n--> XX iOS Simulator x86_64h libssl Compilation"
  moveLibreSSLOutputInPlace $target $OUTPUT

fi;

#############################################
##  iOS Simulator x86_64 libssl Compilation
#############################################

target=simulator_x86_64
if needsRebuilding "$target" && elementIn "$target" "${libressl_build_targets[@]}"; then

  printf "\n\n--> iOS Simulator x86_64 libssl Compilation"

  DEVROOT=$XCODE/Platforms/iPhoneSimulator.platform/Developer
  SDKROOT=$DEVROOT/SDKs/iPhoneSimulator${IOS}.sdk

  echo "prefix: $PREFIX/$target"
  echo "SDKROOT: $SDKROOT"
  echo "CPPFLAGS: $CPPFLAGS"
  echo "IOS: $IOS"

  ./configure --host=x86_64-apple-darwin --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang" \
    CPPFLAGS="-I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch x86_64 -miphoneos-version-min=${MIN_IOS_VERSION} -pipe -no-cpp-precomp -isysroot $SDKROOT" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD=$DEVROOT/usr/bin/ld

  makeLibreSSL
  echo -printf "\n\n--> XX iOS Simulator x86_64 libssl Compilation"
  moveLibreSSLOutputInPlace $target $OUTPUT

fi;

#############################################
##  iOS Simulator arm64e libssl Compilation
#############################################

target=simulator_arm64e
if needsRebuilding "$target" && elementIn "$target" "${libressl_build_targets[@]}"; then

  printf "\n\n--> iOS Simulator arm64e libssl Compilation"

  DEVROOT=$XCODE/Platforms/iPhoneSimulator.platform/Developer
  SDKROOT=$DEVROOT/SDKs/iPhoneSimulator${IOS}.sdk

  #./configure --build=aarch64-apple-darwin --host=aarch64-apple-darwin22 --prefix="$PREFIX/$target" \
  ./configure --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -target arm64-apple-ios${IOS}-simulator" \
    CPPFLAGS="-I$SDKROOT/usr/include/ -target arm64-apple-ios${IOS}-simulator" \
    CFLAGS="$CPPFLAGS -arch arm64e -miphoneos-version-min=${MIN_IOS_VERSION} -pipe -no-cpp-precomp -isysroot $SDKROOT" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD=$DEVROOT/usr/bin/ld

  makeLibreSSL
  printf "\n\n--> XX iOS Simulator arm64e libssl Compilation"
  moveLibreSSLOutputInPlace $target $OUTPUT

fi;

#############################################
##  iOS Simulator arm64 libssl Compilation
#############################################

target=simulator_arm64
if needsRebuilding "$target" && elementIn "$target" "${libressl_build_targets[@]}"; then

  printf "\n\n--> iOS Simulator arm64 libssl Compilation"

  DEVROOT=$XCODE/Platforms/iPhoneSimulator.platform/Developer
  SDKROOT=$DEVROOT/SDKs/iPhoneSimulator${IOS}.sdk

  ./configure --host=aarch64-apple-darwin --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -target arm64-apple-ios${IOS}-simulator" \
    CPPFLAGS="-I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch arm64 -miphoneos-version-min=${MIN_IOS_VERSION} -pipe -no-cpp-precomp -isysroot $SDKROOT" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD=$DEVROOT/usr/bin/ld

  makeLibreSSL
  printf "\n\n--> XX iOS Simulator arm64 libssl Compilation"
  moveLibreSSLOutputInPlace $target $OUTPUT

fi;


##################################
##  iOS arm64 libssl Compilation
##################################

target=ios-arm64
if needsRebuilding "$target" && elementIn "$target" "${libressl_build_targets[@]}"; then

  printf "\n\n--> iOS arm64 libssl Compilation"

  DEVROOT=$XCODE/Platforms/iPhoneOS.platform/Developer
  SDKROOT=$DEVROOT/SDKs/iPhoneOS${IOS}.sdk

  ./configure --host=aarch64-apple-darwin --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -isysroot $SDKROOT" \
    CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch arm64 -miphoneos-version-min=${MIN_IOS_VERSION} -pipe -no-cpp-precomp" \
    CPP="/usr/bin/cpp -D__arm__=1 $CPPFLAGS" \
    LD=$DEVROOT/usr/bin/ld

  makeLibreSSL
  printf "\n\n--> XX iOS arm64 libssl Compilation"
  moveLibreSSLOutputInPlace $target $OUTPUT

fi;

###################################
##  iOS arm64e libssl Compilation
###################################

target=ios-arm64e
if elementIn "$target" "${libressl_build_targets[@]}"; then

  printf "\n\n--> iOS arm64e libssl Compilation"

  DEVROOT=$XCODE/Platforms/iPhoneOS.platform/Developer
  SDKROOT=$DEVROOT/SDKs/iPhoneOS${IOS}.sdk

  ./configure --host=aarch64-apple-darwin19 --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -isysroot $SDKROOT" \
    CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch arm64e -miphoneos-version-min=${MIN_IOS_VERSION} -pipe -no-cpp-precomp" \
    CPP="/usr/bin/cpp -D__arm__=1 $CPPFLAGS" \
    LD=$DEVROOT/usr/bin/ld

  makeLibreSSL
  printf "\n\n--> XX iOS arm64e libssl Compilation"
  moveLibreSSLOutputInPlace $target $OUTPUT

fi;

##############################################
##  macOS Catalyst x86_64 libssl Compilation
##############################################

target=catalyst_x86_64
if needsRebuilding "$target" && elementIn "$target" "${libressl_build_targets[@]}"; then

  printf "\n\n--> macOS Catalyst x86_64 libssl Compilation"

  DEVROOT=$XCODE/Platforms/MacOSX.platform/Developer
  SDKROOT=$DEVROOT/SDKs/MacOSX${MACOSX}.sdk

  ./configure --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -target x86_64-apple-ios${IOS}-macabi -isysroot $SDKROOT" \
    CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch x86_64 -pipe -no-cpp-precomp" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD="/usr/bin/ld"

  makeLibreSSL
  printf "\n\n--> XX macOS Catalyst x86_64 libssl Compilation"
  moveLibreSSLOutputInPlace $target $OUTPUT

fi;

#############################################
##  macOS Catalyst arm64 libssl Compilation
#############################################

target=catalyst_arm64
if needsRebuilding "$target" && elementIn "$target" "${libressl_build_targets[@]}"; then

  printf "\n\n--> macOS Catalyst arm64 libssl Compilation"

  DEVROOT=$XCODE/Platforms/MacOSX.platform/Developer
  SDKROOT=$DEVROOT/SDKs/MacOSX${MACOSX}.sdk

  ./configure --build=aarch64-apple-darwin --host=aarch64-apple-darwin22 --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -target arm64-apple-ios${IOS}-macabi -isysroot $SDKROOT" \
    CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch arm64 -pipe -no-cpp-precomp" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD="/usr/bin/ld"

  makeLibreSSL
  printf "\n\n--> XX macOS Catalyst arm64 libssl Compilation"
  moveLibreSSLOutputInPlace $target $OUTPUT

fi;


#####################################
##  macOS x86_64 libssl Compilation
#####################################

target=macos_x86_64
if needsRebuilding "$target" && elementIn "$target" "${libressl_build_targets[@]}"; then

  printf "\n\n--> macOS x86_64 libssl Compilation"

  DEVROOT=$XCODE/Platforms/MacOSX.platform/Developer
  SDKROOT=$DEVROOT/SDKs/MacOSX${MACOSX}.sdk

  ./configure --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -target x86_64-apple-darwin -isysroot $SDKROOT" \
    CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch x86_64 -pipe -no-cpp-precomp" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD="/usr/bin/clang -target x86_64-apple-darwin"

  makeLibreSSL
  printf "\n\n--> XX macOS x86_64 libssl Compilation"
  moveLibreSSLOutputInPlace $target $OUTPUT

fi;


######################################
##  macOS x86_64h libssl Compilation
######################################

target=macos_x86_64h
if needsRebuilding "$target" && elementIn "$target" "${libressl_build_targets[@]}"; then

  printf "\n\n--> macOS x86_64h libssl Compilation"

  DEVROOT=$XCODE/Platforms/MacOSX.platform/Developer
  SDKROOT=$DEVROOT/SDKs/MacOSX${MACOSX}.sdk

  ./configure --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -isysroot $SDKROOT" \
    CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch x86_64h -pipe -no-cpp-precomp" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD="/usr/bin/ld"

  makeLibreSSL
  printf "\n\n--> XX macOS x86_64h libssl Compilation"
  moveLibreSSLOutputInPlace $target $OUTPUT

fi;

#####################################
##  macOS arm64 libssl Compilation
#####################################

target=macos_arm64
if needsRebuilding "$target" && elementIn "$target" "${libressl_build_targets[@]}"; then

  printf "\n\n--> macOS arm64 libssl Compilation"

  DEVROOT=$XCODE/Platforms/MacOSX.platform/Developer
  SDKROOT=$DEVROOT/SDKs/MacOSX${MACOSX}.sdk

  ./configure --build=aarch64-apple-darwin --host=aarch64-apple-darwin22 --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -isysroot $SDKROOT" \
    CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch arm64 -pipe -no-cpp-precomp" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD="/usr/bin/ld"

  makeLibreSSL
  printf "\n\n--> XX macOS arm64 libssl Compilation"
  moveLibreSSLOutputInPlace $target $OUTPUT

fi;

# TODO: This one isn't working - use "host" to cross compile
#####################################
##  macOS arm64e libssl Compilation
#####################################

target=macos_arm64e
if needsRebuilding "$target" && elementIn "$target" "${libressl_build_targets[@]}"; then

  printf "\n\n--> macOS arm64e libssl Compilation"

  DEVROOT=$XCODE/Platforms/MacOSX.platform/Developer
  SDKROOT=$DEVROOT/SDKs/MacOSX${MACOSX}.sdk

  ./configure --prefix="$PREFIX/$target" \
    CC="/usr/bin/clang -target arm64-apple-darwin -isysroot $SDKROOT" \
    CPPFLAGS="-fembed-bitcode -I$SDKROOT/usr/include/" \
    CFLAGS="$CPPFLAGS -arch arm64e -pipe -no-cpp-precomp" \
    CPP="/usr/bin/cpp $CPPFLAGS" \
    LD="/usr/bin/ld -target arm64-apple-darwin"

  makeLibreSSL
  printf "\n\n--> XX macOS arm64e libssl Compilation"
  moveLibreSSLOutputInPlace $target $OUTPUT

fi;


####################################
## lipo & XCFrameworks for LibreSSL
####################################

## lipo & XCFramework

macos=()
catalyst=()
simulator=()
ios=()


for target in "${libressl_link_targets[@]}"
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

XCFRAMEWORK_LIBSSL_CMD="xcodebuild -create-xcframework"
XCFRAMEWORK_LIBCRYPTO_CMD="xcodebuild -create-xcframework"

framework_targets=()

if [ ${#ios[@]} -gt 0 ]; then
  lipo_libssl="lipo -create "
  lipo_libcrypto="lipo -create "

  framework_targets+=("ios")
  mkdir -p $OUTPUT/ios/lib
  mkdir -p $OUTPUT/ios/include

  for target in "${ios[@]}"
  do
    lipo_libssl="$lipo_libssl $OUTPUT/$target/lib/libssl.a"
    lipo_libcrypto="$lipo_libcrypto $OUTPUT/$target/lib/libcrypto.a"
    rsync -a $OUTPUT/$target/include/* $OUTPUT/ios/include
  done

  lipo_libssl="$lipo_libssl -output $OUTPUT/ios/lib/libssl.a"
  echo $lipo_libssl
  eval $lipo_libssl

  lipo_libcrypto="$lipo_libcrypto -output $OUTPUT/ios/lib/libcrypto.a"
  echo $lipo_libcrypto
  eval $lipo_libcrypto

  XCFRAMEWORK_LIBSSL_CMD="$XCFRAMEWORK_LIBSSL_CMD -library $OUTPUT/ios/lib/libssl.a"
  XCFRAMEWORK_LIBSSL_CMD="$XCFRAMEWORK_LIBSSL_CMD -headers $OUTPUT/ios/include"

  XCFRAMEWORK_LIBCRYPTO_CMD="$XCFRAMEWORK_LIBCRYPTO_CMD -library $OUTPUT/ios/lib/libcrypto.a"
  XCFRAMEWORK_LIBCRYPTO_CMD="$XCFRAMEWORK_LIBCRYPTO_CMD -headers $OUTPUT/ios/include"
fi

if [ ${#catalyst[@]} -gt 0 ]; then
  lipo_libssl="lipo -create "
  lipo_libcrypto="lipo -create "

  framework_targets+=("catalyst")
  mkdir -p $OUTPUT/catalyst/lib
  mkdir -p $OUTPUT/catalyst/include

  for target in "${catalyst[@]}"
  do
    lipo_libssl="$lipo_libssl $OUTPUT/$target/lib/libssl.a"
    lipo_libcrypto="$lipo_libcrypto $OUTPUT/$target/lib/libcrypto.a"
    rsync -a $OUTPUT/$target/include/* $OUTPUT/catalyst/include
  done

  lipo_libssl="$lipo_libssl -output $OUTPUT/catalyst/lib/libssl.a"
  echo $lipo_libssl
  eval $lipo_libssl

  lipo_libcrypto="$lipo_libcrypto -output $OUTPUT/catalyst/lib/libcrypto.a"
  echo $lipo_libcrypto
  eval $lipo_libcrypto

  XCFRAMEWORK_LIBSSL_CMD="$XCFRAMEWORK_LIBSSL_CMD -library $OUTPUT/catalyst/lib/libssl.a"
  XCFRAMEWORK_LIBSSL_CMD="$XCFRAMEWORK_LIBSSL_CMD -headers $OUTPUT/catalyst/include"

  XCFRAMEWORK_LIBCRYPTO_CMD="$XCFRAMEWORK_LIBCRYPTO_CMD -library $OUTPUT/catalyst/lib/libcrypto.a"
  XCFRAMEWORK_LIBCRYPTO_CMD="$XCFRAMEWORK_LIBCRYPTO_CMD -headers $OUTPUT/catalyst/include"
fi

if [ ${#macos[@]} -gt 0 ]; then
  lipo_libssl="lipo -create "
  lipo_libcrypto="lipo -create "

  framework_targets+=("macos")
  mkdir -p $OUTPUT/macos/lib
  mkdir -p $OUTPUT/macos/include

  for target in "${macos[@]}"
  do
    lipo_libssl="$lipo_libssl $OUTPUT/$target/lib/libssl.a"
    lipo_libcrypto="$lipo_libcrypto $OUTPUT/$target/lib/libcrypto.a"
    rsync -a $OUTPUT/$target/include/* $OUTPUT/macos/include
  done

  lipo_libssl="$lipo_libssl -output $OUTPUT/macos/lib/libssl.a"
  echo $lipo_libssl
  eval $lipo_libssl

  lipo_libcrypto="$lipo_libcrypto -output $OUTPUT/macos/lib/libcrypto.a"
  echo $lipo_libcrypto
  eval $lipo_libcrypto

  XCFRAMEWORK_LIBSSL_CMD="$XCFRAMEWORK_LIBSSL_CMD -library $OUTPUT/macos/lib/libssl.a"
  XCFRAMEWORK_LIBSSL_CMD="$XCFRAMEWORK_LIBSSL_CMD -headers $OUTPUT/macos/include"

  XCFRAMEWORK_LIBCRYPTO_CMD="$XCFRAMEWORK_LIBCRYPTO_CMD -library $OUTPUT/macos/lib/libcrypto.a"
  XCFRAMEWORK_LIBCRYPTO_CMD="$XCFRAMEWORK_LIBCRYPTO_CMD -headers $OUTPUT/macos/include"
fi

if [ ${#simulator[@]} -gt 0 ]; then
  lipo_libssl="lipo -create "
  lipo_libcrypto="lipo -create "

  framework_targets+=("simulator")
  mkdir -p $OUTPUT/simulator/lib
  mkdir -p $OUTPUT/simulator/include

  for target in "${simulator[@]}"
  do
    lipo_libssl="$lipo_libssl $OUTPUT/$target/lib/libssl.a"
    lipo_libcrypto="$lipo_libcrypto $OUTPUT/$target/lib/libcrypto.a"
    rsync -a $OUTPUT/$target/include/* $OUTPUT/simulator/include
  done

  lipo_libssl="$lipo_libssl -output $OUTPUT/simulator/lib/libssl.a"
  echo $lipo_libssl
  eval $lipo_libssl

  lipo_libcrypto="$lipo_libcrypto -output $OUTPUT/simulator/lib/libcrypto.a"
  echo $lipo_libcrypto
  eval $lipo_libcrypto

  XCFRAMEWORK_LIBSSL_CMD="$XCFRAMEWORK_LIBSSL_CMD -library $OUTPUT/simulator/lib/libssl.a"
  XCFRAMEWORK_LIBSSL_CMD="$XCFRAMEWORK_LIBSSL_CMD -headers $OUTPUT/simulator/include"

  XCFRAMEWORK_LIBCRYPTO_CMD="$XCFRAMEWORK_LIBCRYPTO_CMD -library $OUTPUT/simulator/lib/libcrypto.a"
  XCFRAMEWORK_LIBCRYPTO_CMD="$XCFRAMEWORK_LIBCRYPTO_CMD -headers $OUTPUT/simulator/include"
fi

XCFRAMEWORK_LIBSSL_CMD="$XCFRAMEWORK_LIBSSL_CMD -output $XCFRAMEWORKS/libssl.xcframework"
XCFRAMEWORK_LIBCRYPTO_CMD="$XCFRAMEWORK_LIBCRYPTO_CMD -output $XCFRAMEWORKS/libcrypto.xcframework"


echo $XCFRAMEWORK_LIBSSL_CMD
eval $XCFRAMEWORK_LIBSSL_CMD
echo $XCFRAMEWORK_LIBCRYPTO_CMD
eval $XCFRAMEWORK_LIBCRYPTO_CMD


cd ..
