# libsqcloud XCFramework for iOS, macOS & Catalyst

A script to compile libsqcloud (C SDK for SQLite Cloud) and its dependencies (libtls) to an XCFramework supporting the latest Apple OS versions

Instructions:
Build:
```
bash build-apple.sh
```

The resulting XCFrameworks `libsqcloud.xcframework` will be saved in the "output" directory.
The script also compiles the libtls library using the `libressl.sh` script, if needed, and includes it into the libsqcloud.a files.
