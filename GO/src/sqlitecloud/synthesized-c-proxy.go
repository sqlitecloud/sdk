package sqlitecloud

// #cgo CFLAGS: -Wno-multichar -Wno-macro-redefined 
// #cgo LDFLAGS: -L. 
// #include <stdlib.h>
// /*
//    LZ4 - Fast LZ compression algorithm
//    Copyright (C) 2011-present, Yann Collet.
// 
//    BSD 2-Clause License (http://www.opensource.org/licenses/bsd-license.php)
// 
//    Redistribution and use in source and binary forms, with or without
//    modification, are permitted provided that the following conditions are
//    met:
// 
//        * Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
//        * Redistributions in binary form must reproduce the above
//    copyright notice, this list of conditions and the following disclaimer
//    in the documentation and/or other materials provided with the
//    distribution.
// 
//    THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//    "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//    LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//    A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//    OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//    SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
//    LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
//    DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
//    THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//    (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//    OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// 
//    You can contact the author at :
//     - LZ4 homepage : http://www.lz4.org
//     - LZ4 source repository : https://github.com/lz4/lz4
// */
// 
// /*-************************************
// *  Tuning parameters
// **************************************/
// /*
//  * LZ4_HEAPMODE :
//  * Select how default compression functions will allocate memory for their hash table,
//  * in memory stack (0:default, fastest), or in memory heap (1:requires malloc()).
//  */
// #ifndef LZ4_HEAPMODE
// #  define LZ4_HEAPMODE 0
// #endif
// 
// /*
//  * LZ4_ACCELERATION_DEFAULT :
//  * Select "acceleration" for LZ4_compress_fast() when parameter value <= 0
//  */
// #define LZ4_ACCELERATION_DEFAULT 1
// /*
//  * LZ4_ACCELERATION_MAX :
//  * Any "acceleration" value higher than this threshold
//  * get treated as LZ4_ACCELERATION_MAX instead (fix #876)
//  */
// #define LZ4_ACCELERATION_MAX 65537
// 
// 
// /*-************************************
// *  CPU Feature Detection
// **************************************/
// /* LZ4_FORCE_MEMORY_ACCESS
//  * By default, access to unaligned memory is controlled by `memcpy()`, which is safe and portable.
//  * Unfortunately, on some target/compiler combinations, the generated assembly is sub-optimal.
//  * The below switch allow to select different access method for improved performance.
//  * Method 0 (default) : use `memcpy()`. Safe and portable.
//  * Method 1 : `__packed` statement. It depends on compiler extension (ie, not portable).
//  *            This method is safe if your compiler supports it, and *generally* as fast or faster than `memcpy`.
//  * Method 2 : direct access. This method is portable but violate C standard.
//  *            It can generate buggy code on targets which assembly generation depends on alignment.
//  *            But in some circumstances, it's the only known way to get the most performance (ie GCC + ARMv6)
//  * See https://fastcompression.blogspot.fr/2015/08/accessing-unaligned-memory.html for details.
//  * Prefer these methods in priority order (0 > 1 > 2)
//  */
// #ifndef LZ4_FORCE_MEMORY_ACCESS   /* can be defined externally */
// #  if defined(__GNUC__) && \
//   ( defined(__ARM_ARCH_6__) || defined(__ARM_ARCH_6J__) || defined(__ARM_ARCH_6K__) \
//   || defined(__ARM_ARCH_6Z__) || defined(__ARM_ARCH_6ZK__) || defined(__ARM_ARCH_6T2__) )
// #    define LZ4_FORCE_MEMORY_ACCESS 2
// #  elif (defined(__INTEL_COMPILER) && !defined(_WIN32)) || defined(__GNUC__)
// #    define LZ4_FORCE_MEMORY_ACCESS 1
// #  endif
// #endif
// 
// /*
//  * LZ4_FORCE_SW_BITCOUNT
//  * Define this parameter if your target system or compiler does not support hardware bit count
//  */
// #if defined(_MSC_VER) && defined(_WIN32_WCE)   /* Visual Studio for WinCE doesn't support Hardware bit count */
// #  undef  LZ4_FORCE_SW_BITCOUNT  /* avoid double def */
// #  define LZ4_FORCE_SW_BITCOUNT
// #endif
// 
// 
// 
// /*-************************************
// *  Dependency
// **************************************/
// /*
//  * LZ4_SRC_INCLUDED:
//  * Amalgamation flag, whether lz4.c is included
//  */
// #ifndef LZ4_SRC_INCLUDED
// #  define LZ4_SRC_INCLUDED 1
// #endif
// 
// #ifndef LZ4_STATIC_LINKING_ONLY
// #define LZ4_STATIC_LINKING_ONLY
// #endif
// 
// #ifndef LZ4_DISABLE_DEPRECATE_WARNINGS
// #define LZ4_DISABLE_DEPRECATE_WARNINGS /* due to LZ4_decompress_safe_withPrefix64k */
// #endif
// 
// #define LZ4_STATIC_LINKING_ONLY  /* LZ4_DISTANCE_MAX */
// /*
//  *  LZ4 - Fast LZ compression algorithm
//  *  Header File
//  *  Copyright (C) 2011-present, Yann Collet.
// 
//    BSD 2-Clause License (http://www.opensource.org/licenses/bsd-license.php)
// 
//    Redistribution and use in source and binary forms, with or without
//    modification, are permitted provided that the following conditions are
//    met:
// 
//        * Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
//        * Redistributions in binary form must reproduce the above
//    copyright notice, this list of conditions and the following disclaimer
//    in the documentation and/or other materials provided with the
//    distribution.
// 
//    THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//    "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//    LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//    A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//    OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//    SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
//    LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
//    DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
//    THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//    (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//    OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// 
//    You can contact the author at :
//     - LZ4 homepage : http://www.lz4.org
//     - LZ4 source repository : https://github.com/lz4/lz4
// */
// #if defined (__cplusplus)
// extern "C" {
// #endif
// 
// #ifndef LZ4_H_2983827168210
// #define LZ4_H_2983827168210
// 
// /* --- Dependency --- */
// #include <stddef.h>   /* size_t */
// 
// 
// /**
//   Introduction
// 
//   LZ4 is lossless compression algorithm, providing compression speed >500 MB/s per core,
//   scalable with multi-cores CPU. It features an extremely fast decoder, with speed in
//   multiple GB/s per core, typically reaching RAM speed limits on multi-core systems.
// 
//   The LZ4 compression library provides in-memory compression and decompression functions.
//   It gives full buffer control to user.
//   Compression can be done in:
//     - a single step (described as Simple Functions)
//     - a single step, reusing a context (described in Advanced Functions)
//     - unbounded multiple steps (described as Streaming compression)
// 
//   lz4.h generates and decodes LZ4-compressed blocks (doc/lz4_Block_format.md).
//   Decompressing such a compressed block requires additional metadata.
//   Exact metadata depends on exact decompression function.
//   For the typical case of LZ4_decompress_safe(),
//   metadata includes block's compressed size, and maximum bound of decompressed size.
//   Each application is free to encode and pass such metadata in whichever way it wants.
// 
//   lz4.h only handle blocks, it can not generate Frames.
// 
//   Blocks are different from Frames (doc/lz4_Frame_format.md).
//   Frames bundle both blocks and metadata in a specified manner.
//   Embedding metadata is required for compressed data to be self-contained and portable.
//   Frame format is delivered through a companion API, declared in lz4frame.h.
//   The `lz4` CLI can only manage frames.
// */
// 
// /*^***************************************************************
// *  Export parameters
// *****************************************************************/
// /*
// *  LZ4_DLL_EXPORT :
// *  Enable exporting of functions when building a Windows DLL
// *  LZ4LIB_VISIBILITY :
// *  Control library symbols visibility.
// */
// #ifndef LZ4LIB_VISIBILITY
// #  if defined(__GNUC__) && (__GNUC__ >= 4)
// #    define LZ4LIB_VISIBILITY __attribute__ ((visibility ("default")))
// #  else
// #    define LZ4LIB_VISIBILITY
// #  endif
// #endif
// #if defined(LZ4_DLL_EXPORT) && (LZ4_DLL_EXPORT==1)
// #  define LZ4LIB_API __declspec(dllexport) LZ4LIB_VISIBILITY
// #elif defined(LZ4_DLL_IMPORT) && (LZ4_DLL_IMPORT==1)
// #  define LZ4LIB_API __declspec(dllimport) LZ4LIB_VISIBILITY /* It isn't required but allows to generate better code, saving a function pointer load from the IAT and an indirect jump.*/
// #else
// #  define LZ4LIB_API LZ4LIB_VISIBILITY
// #endif
// 
// /*------   Version   ------*/
// #define LZ4_VERSION_MAJOR    1    /* for breaking interface changes  */
// #define LZ4_VERSION_MINOR    9    /* for new (non-breaking) interface capabilities */
// #define LZ4_VERSION_RELEASE  3    /* for tweaks, bug-fixes, or development */
// 
// #define LZ4_VERSION_NUMBER (LZ4_VERSION_MAJOR *100*100 + LZ4_VERSION_MINOR *100 + LZ4_VERSION_RELEASE)
// 
// #define LZ4_LIB_VERSION LZ4_VERSION_MAJOR.LZ4_VERSION_MINOR.LZ4_VERSION_RELEASE
// #define LZ4_QUOTE(str) #str
// #define LZ4_EXPAND_AND_QUOTE(str) LZ4_QUOTE(str)
// #define LZ4_VERSION_STRING LZ4_EXPAND_AND_QUOTE(LZ4_LIB_VERSION)
// 
// LZ4LIB_API int LZ4_versionNumber (void);  /**< library version number; useful to check dll version */
// LZ4LIB_API const char* LZ4_versionString (void);   /**< library version string; useful to check dll version */
// 
// 
// /*-************************************
// *  Tuning parameter
// **************************************/
// /*!
//  * LZ4_MEMORY_USAGE :
//  * Memory usage formula : N->2^N Bytes (examples : 10 -> 1KB; 12 -> 4KB ; 16 -> 64KB; 20 -> 1MB; etc.)
//  * Increasing memory usage improves compression ratio.
//  * Reduced memory usage may improve speed, thanks to better cache locality.
//  * Default value is 14, for 16KB, which nicely fits into Intel x86 L1 cache
//  */
// #ifndef LZ4_MEMORY_USAGE
// # define LZ4_MEMORY_USAGE 14
// #endif
// 
// 
// /*-************************************
// *  Simple Functions
// **************************************/
// /*! LZ4_compress_default() :
//  *  Compresses 'srcSize' bytes from buffer 'src'
//  *  into already allocated 'dst' buffer of size 'dstCapacity'.
//  *  Compression is guaranteed to succeed if 'dstCapacity' >= LZ4_compressBound(srcSize).
//  *  It also runs faster, so it's a recommended setting.
//  *  If the function cannot compress 'src' into a more limited 'dst' budget,
//  *  compression stops *immediately*, and the function result is zero.
//  *  In which case, 'dst' content is undefined (invalid).
//  *      srcSize : max supported value is LZ4_MAX_INPUT_SIZE.
//  *      dstCapacity : size of buffer 'dst' (which must be already allocated)
//  *     @return  : the number of bytes written into buffer 'dst' (necessarily <= dstCapacity)
//  *                or 0 if compression fails
//  * Note : This function is protected against buffer overflow scenarios (never writes outside 'dst' buffer, nor read outside 'source' buffer).
//  */
// LZ4LIB_API int LZ4_compress_default(const char* src, char* dst, int srcSize, int dstCapacity);
// 
// /*! LZ4_decompress_safe() :
//  *  compressedSize : is the exact complete size of the compressed block.
//  *  dstCapacity : is the size of destination buffer (which must be already allocated), presumed an upper bound of decompressed size.
//  * @return : the number of bytes decompressed into destination buffer (necessarily <= dstCapacity)
//  *           If destination buffer is not large enough, decoding will stop and output an error code (negative value).
//  *           If the source stream is detected malformed, the function will stop decoding and return a negative result.
//  * Note 1 : This function is protected against malicious data packets :
//  *          it will never writes outside 'dst' buffer, nor read outside 'source' buffer,
//  *          even if the compressed block is maliciously modified to order the decoder to do these actions.
//  *          In such case, the decoder stops immediately, and considers the compressed block malformed.
//  * Note 2 : compressedSize and dstCapacity must be provided to the function, the compressed block does not contain them.
//  *          The implementation is free to send / store / derive this information in whichever way is most beneficial.
//  *          If there is a need for a different format which bundles together both compressed data and its metadata, consider looking at lz4frame.h instead.
//  */
// LZ4LIB_API int LZ4_decompress_safe (const char* src, char* dst, int compressedSize, int dstCapacity);
// 
// 
// /*-************************************
// *  Advanced Functions
// **************************************/
// #define LZ4_MAX_INPUT_SIZE        0x7E000000   /* 2 113 929 216 bytes */
// #define LZ4_COMPRESSBOUND(isize)  ((unsigned)(isize) > (unsigned)LZ4_MAX_INPUT_SIZE ? 0 : (isize) + ((isize)/255) + 16)
// 
// /*! LZ4_compressBound() :
//     Provides the maximum size that LZ4 compression may output in a "worst case" scenario (input data not compressible)
//     This function is primarily useful for memory allocation purposes (destination buffer size).
//     Macro LZ4_COMPRESSBOUND() is also provided for compilation-time evaluation (stack memory allocation for example).
//     Note that LZ4_compress_default() compresses faster when dstCapacity is >= LZ4_compressBound(srcSize)
//         inputSize  : max supported value is LZ4_MAX_INPUT_SIZE
//         return : maximum output size in a "worst case" scenario
//               or 0, if input size is incorrect (too large or negative)
// */
// LZ4LIB_API int LZ4_compressBound(int inputSize);
// 
// /*! LZ4_compress_fast() :
//     Same as LZ4_compress_default(), but allows selection of "acceleration" factor.
//     The larger the acceleration value, the faster the algorithm, but also the lesser the compression.
//     It's a trade-off. It can be fine tuned, with each successive value providing roughly +~3% to speed.
//     An acceleration value of "1" is the same as regular LZ4_compress_default()
//     Values <= 0 will be replaced by LZ4_ACCELERATION_DEFAULT (currently == 1, see lz4.c).
//     Values > LZ4_ACCELERATION_MAX will be replaced by LZ4_ACCELERATION_MAX (currently == 65537, see lz4.c).
// */
// LZ4LIB_API int LZ4_compress_fast (const char* src, char* dst, int srcSize, int dstCapacity, int acceleration);
// 
// 
// /*! LZ4_compress_fast_extState() :
//  *  Same as LZ4_compress_fast(), using an externally allocated memory space for its state.
//  *  Use LZ4_sizeofState() to know how much memory must be allocated,
//  *  and allocate it on 8-bytes boundaries (using `malloc()` typically).
//  *  Then, provide this buffer as `void* state` to compression function.
//  */
// LZ4LIB_API int LZ4_sizeofState(void);
// LZ4LIB_API int LZ4_compress_fast_extState (void* state, const char* src, char* dst, int srcSize, int dstCapacity, int acceleration);
// 
// 
// /*! LZ4_compress_destSize() :
//  *  Reverse the logic : compresses as much data as possible from 'src' buffer
//  *  into already allocated buffer 'dst', of size >= 'targetDestSize'.
//  *  This function either compresses the entire 'src' content into 'dst' if it's large enough,
//  *  or fill 'dst' buffer completely with as much data as possible from 'src'.
//  *  note: acceleration parameter is fixed to "default".
//  *
//  * *srcSizePtr : will be modified to indicate how many bytes where read from 'src' to fill 'dst'.
//  *               New value is necessarily <= input value.
//  * @return : Nb bytes written into 'dst' (necessarily <= targetDestSize)
//  *           or 0 if compression fails.
//  *
//  * Note : from v1.8.2 to v1.9.1, this function had a bug (fixed un v1.9.2+):
//  *        the produced compressed content could, in specific circumstances,
//  *        require to be decompressed into a destination buffer larger
//  *        by at least 1 byte than the content to decompress.
//  *        If an application uses `LZ4_compress_destSize()`,
//  *        it's highly recommended to update liblz4 to v1.9.2 or better.
//  *        If this can't be done or ensured,
//  *        the receiving decompression function should provide
//  *        a dstCapacity which is > decompressedSize, by at least 1 byte.
//  *        See https://github.com/lz4/lz4/issues/859 for details
//  */
// LZ4LIB_API int LZ4_compress_destSize (const char* src, char* dst, int* srcSizePtr, int targetDstSize);
// 
// 
// /*! LZ4_decompress_safe_partial() :
//  *  Decompress an LZ4 compressed block, of size 'srcSize' at position 'src',
//  *  into destination buffer 'dst' of size 'dstCapacity'.
//  *  Up to 'targetOutputSize' bytes will be decoded.
//  *  The function stops decoding on reaching this objective.
//  *  This can be useful to boost performance
//  *  whenever only the beginning of a block is required.
//  *
//  * @return : the number of bytes decoded in `dst` (necessarily <= targetOutputSize)
//  *           If source stream is detected malformed, function returns a negative result.
//  *
//  *  Note 1 : @return can be < targetOutputSize, if compressed block contains less data.
//  *
//  *  Note 2 : targetOutputSize must be <= dstCapacity
//  *
//  *  Note 3 : this function effectively stops decoding on reaching targetOutputSize,
//  *           so dstCapacity is kind of redundant.
//  *           This is because in older versions of this function,
//  *           decoding operation would still write complete sequences.
//  *           Therefore, there was no guarantee that it would stop writing at exactly targetOutputSize,
//  *           it could write more bytes, though only up to dstCapacity.
//  *           Some "margin" used to be required for this operation to work properly.
//  *           Thankfully, this is no longer necessary.
//  *           The function nonetheless keeps the same signature, in an effort to preserve API compatibility.
//  *
//  *  Note 4 : If srcSize is the exact size of the block,
//  *           then targetOutputSize can be any value,
//  *           including larger than the block's decompressed size.
//  *           The function will, at most, generate block's decompressed size.
//  *
//  *  Note 5 : If srcSize is _larger_ than block's compressed size,
//  *           then targetOutputSize **MUST** be <= block's decompressed size.
//  *           Otherwise, *silent corruption will occur*.
//  */
// LZ4LIB_API int LZ4_decompress_safe_partial (const char* src, char* dst, int srcSize, int targetOutputSize, int dstCapacity);
// 
// 
// /*-*********************************************
// *  Streaming Compression Functions
// ***********************************************/
// typedef union LZ4_stream_u LZ4_stream_t;  /* incomplete type (defined later) */
// 
// LZ4LIB_API LZ4_stream_t* LZ4_createStream(void);
// LZ4LIB_API int           LZ4_freeStream (LZ4_stream_t* streamPtr);
// 
// /*! LZ4_resetStream_fast() : v1.9.0+
//  *  Use this to prepare an LZ4_stream_t for a new chain of dependent blocks
//  *  (e.g., LZ4_compress_fast_continue()).
//  *
//  *  An LZ4_stream_t must be initialized once before usage.
//  *  This is automatically done when created by LZ4_createStream().
//  *  However, should the LZ4_stream_t be simply declared on stack (for example),
//  *  it's necessary to initialize it first, using LZ4_initStream().
//  *
//  *  After init, start any new stream with LZ4_resetStream_fast().
//  *  A same LZ4_stream_t can be re-used multiple times consecutively
//  *  and compress multiple streams,
//  *  provided that it starts each new stream with LZ4_resetStream_fast().
//  *
//  *  LZ4_resetStream_fast() is much faster than LZ4_initStream(),
//  *  but is not compatible with memory regions containing garbage data.
//  *
//  *  Note: it's only useful to call LZ4_resetStream_fast()
//  *        in the context of streaming compression.
//  *        The *extState* functions perform their own resets.
//  *        Invoking LZ4_resetStream_fast() before is redundant, and even counterproductive.
//  */
// LZ4LIB_API void LZ4_resetStream_fast (LZ4_stream_t* streamPtr);
// 
// /*! LZ4_loadDict() :
//  *  Use this function to reference a static dictionary into LZ4_stream_t.
//  *  The dictionary must remain available during compression.
//  *  LZ4_loadDict() triggers a reset, so any previous data will be forgotten.
//  *  The same dictionary will have to be loaded on decompression side for successful decoding.
//  *  Dictionary are useful for better compression of small data (KB range).
//  *  While LZ4 accept any input as dictionary,
//  *  results are generally better when using Zstandard's Dictionary Builder.
//  *  Loading a size of 0 is allowed, and is the same as reset.
//  * @return : loaded dictionary size, in bytes (necessarily <= 64 KB)
//  */
// LZ4LIB_API int LZ4_loadDict (LZ4_stream_t* streamPtr, const char* dictionary, int dictSize);
// 
// /*! LZ4_compress_fast_continue() :
//  *  Compress 'src' content using data from previously compressed blocks, for better compression ratio.
//  * 'dst' buffer must be already allocated.
//  *  If dstCapacity >= LZ4_compressBound(srcSize), compression is guaranteed to succeed, and runs faster.
//  *
//  * @return : size of compressed block
//  *           or 0 if there is an error (typically, cannot fit into 'dst').
//  *
//  *  Note 1 : Each invocation to LZ4_compress_fast_continue() generates a new block.
//  *           Each block has precise boundaries.
//  *           Each block must be decompressed separately, calling LZ4_decompress_*() with relevant metadata.
//  *           It's not possible to append blocks together and expect a single invocation of LZ4_decompress_*() to decompress them together.
//  *
//  *  Note 2 : The previous 64KB of source data is __assumed__ to remain present, unmodified, at same address in memory !
//  *
//  *  Note 3 : When input is structured as a double-buffer, each buffer can have any size, including < 64 KB.
//  *           Make sure that buffers are separated, by at least one byte.
//  *           This construction ensures that each block only depends on previous block.
//  *
//  *  Note 4 : If input buffer is a ring-buffer, it can have any size, including < 64 KB.
//  *
//  *  Note 5 : After an error, the stream status is undefined (invalid), it can only be reset or freed.
//  */
// LZ4LIB_API int LZ4_compress_fast_continue (LZ4_stream_t* streamPtr, const char* src, char* dst, int srcSize, int dstCapacity, int acceleration);
// 
// /*! LZ4_saveDict() :
//  *  If last 64KB data cannot be guaranteed to remain available at its current memory location,
//  *  save it into a safer place (char* safeBuffer).
//  *  This is schematically equivalent to a memcpy() followed by LZ4_loadDict(),
//  *  but is much faster, because LZ4_saveDict() doesn't need to rebuild tables.
//  * @return : saved dictionary size in bytes (necessarily <= maxDictSize), or 0 if error.
//  */
// LZ4LIB_API int LZ4_saveDict (LZ4_stream_t* streamPtr, char* safeBuffer, int maxDictSize);
// 
// 
// /*-**********************************************
// *  Streaming Decompression Functions
// *  Bufferless synchronous API
// ************************************************/
// typedef union LZ4_streamDecode_u LZ4_streamDecode_t;   /* tracking context */
// 
// /*! LZ4_createStreamDecode() and LZ4_freeStreamDecode() :
//  *  creation / destruction of streaming decompression tracking context.
//  *  A tracking context can be re-used multiple times.
//  */
// LZ4LIB_API LZ4_streamDecode_t* LZ4_createStreamDecode(void);
// LZ4LIB_API int                 LZ4_freeStreamDecode (LZ4_streamDecode_t* LZ4_stream);
// 
// /*! LZ4_setStreamDecode() :
//  *  An LZ4_streamDecode_t context can be allocated once and re-used multiple times.
//  *  Use this function to start decompression of a new stream of blocks.
//  *  A dictionary can optionally be set. Use NULL or size 0 for a reset order.
//  *  Dictionary is presumed stable : it must remain accessible and unmodified during next decompression.
//  * @return : 1 if OK, 0 if error
//  */
// LZ4LIB_API int LZ4_setStreamDecode (LZ4_streamDecode_t* LZ4_streamDecode, const char* dictionary, int dictSize);
// 
// /*! LZ4_decoderRingBufferSize() : v1.8.2+
//  *  Note : in a ring buffer scenario (optional),
//  *  blocks are presumed decompressed next to each other
//  *  up to the moment there is not enough remaining space for next block (remainingSize < maxBlockSize),
//  *  at which stage it resumes from beginning of ring buffer.
//  *  When setting such a ring buffer for streaming decompression,
//  *  provides the minimum size of this ring buffer
//  *  to be compatible with any source respecting maxBlockSize condition.
//  * @return : minimum ring buffer size,
//  *           or 0 if there is an error (invalid maxBlockSize).
//  */
// LZ4LIB_API int LZ4_decoderRingBufferSize(int maxBlockSize);
// #define LZ4_DECODER_RING_BUFFER_SIZE(maxBlockSize) (65536 + 14 + (maxBlockSize))  /* for static allocation; maxBlockSize presumed valid */
// 
// /*! LZ4_decompress_*_continue() :
//  *  These decoding functions allow decompression of consecutive blocks in "streaming" mode.
//  *  A block is an unsplittable entity, it must be presented entirely to a decompression function.
//  *  Decompression functions only accepts one block at a time.
//  *  The last 64KB of previously decoded data *must* remain available and unmodified at the memory position where they were decoded.
//  *  If less than 64KB of data has been decoded, all the data must be present.
//  *
//  *  Special : if decompression side sets a ring buffer, it must respect one of the following conditions :
//  *  - Decompression buffer size is _at least_ LZ4_decoderRingBufferSize(maxBlockSize).
//  *    maxBlockSize is the maximum size of any single block. It can have any value > 16 bytes.
//  *    In which case, encoding and decoding buffers do not need to be synchronized.
//  *    Actually, data can be produced by any source compliant with LZ4 format specification, and respecting maxBlockSize.
//  *  - Synchronized mode :
//  *    Decompression buffer size is _exactly_ the same as compression buffer size,
//  *    and follows exactly same update rule (block boundaries at same positions),
//  *    and decoding function is provided with exact decompressed size of each block (exception for last block of the stream),
//  *    _then_ decoding & encoding ring buffer can have any size, including small ones ( < 64 KB).
//  *  - Decompression buffer is larger than encoding buffer, by a minimum of maxBlockSize more bytes.
//  *    In which case, encoding and decoding buffers do not need to be synchronized,
//  *    and encoding ring buffer can have any size, including small ones ( < 64 KB).
//  *
//  *  Whenever these conditions are not possible,
//  *  save the last 64KB of decoded data into a safe buffer where it can't be modified during decompression,
//  *  then indicate where this data is saved using LZ4_setStreamDecode(), before decompressing next block.
// */
// LZ4LIB_API int LZ4_decompress_safe_continue (LZ4_streamDecode_t* LZ4_streamDecode, const char* src, char* dst, int srcSize, int dstCapacity);
// 
// 
// /*! LZ4_decompress_*_usingDict() :
//  *  These decoding functions work the same as
//  *  a combination of LZ4_setStreamDecode() followed by LZ4_decompress_*_continue()
//  *  They are stand-alone, and don't need an LZ4_streamDecode_t structure.
//  *  Dictionary is presumed stable : it must remain accessible and unmodified during decompression.
//  *  Performance tip : Decompression speed can be substantially increased
//  *                    when dst == dictStart + dictSize.
//  */
// LZ4LIB_API int LZ4_decompress_safe_usingDict (const char* src, char* dst, int srcSize, int dstCapcity, const char* dictStart, int dictSize);
// 
// #endif /* LZ4_H_2983827168210 */
// 
// 
// /*^*************************************
//  * !!!!!!   STATIC LINKING ONLY   !!!!!!
//  ***************************************/
// 
// /*-****************************************************************************
//  * Experimental section
//  *
//  * Symbols declared in this section must be considered unstable. Their
//  * signatures or semantics may change, or they may be removed altogether in the
//  * future. They are therefore only safe to depend on when the caller is
//  * statically linked against the library.
//  *
//  * To protect against unsafe usage, not only are the declarations guarded,
//  * the definitions are hidden by default
//  * when building LZ4 as a shared/dynamic library.
//  *
//  * In order to access these declarations,
//  * define LZ4_STATIC_LINKING_ONLY in your application
//  * before including LZ4's headers.
//  *
//  * In order to make their implementations accessible dynamically, you must
//  * define LZ4_PUBLISH_STATIC_FUNCTIONS when building the LZ4 library.
//  ******************************************************************************/
// 
// #ifdef LZ4_STATIC_LINKING_ONLY
// 
// #ifndef LZ4_STATIC_3504398509
// #define LZ4_STATIC_3504398509
// 
// #ifdef LZ4_PUBLISH_STATIC_FUNCTIONS
// #define LZ4LIB_STATIC_API LZ4LIB_API
// #else
// #define LZ4LIB_STATIC_API
// #endif
// 
// 
// /*! LZ4_compress_fast_extState_fastReset() :
//  *  A variant of LZ4_compress_fast_extState().
//  *
//  *  Using this variant avoids an expensive initialization step.
//  *  It is only safe to call if the state buffer is known to be correctly initialized already
//  *  (see above comment on LZ4_resetStream_fast() for a definition of "correctly initialized").
//  *  From a high level, the difference is that
//  *  this function initializes the provided state with a call to something like LZ4_resetStream_fast()
//  *  while LZ4_compress_fast_extState() starts with a call to LZ4_resetStream().
//  */
// LZ4LIB_STATIC_API int LZ4_compress_fast_extState_fastReset (void* state, const char* src, char* dst, int srcSize, int dstCapacity, int acceleration);
// 
// /*! LZ4_attach_dictionary() :
//  *  This is an experimental API that allows
//  *  efficient use of a static dictionary many times.
//  *
//  *  Rather than re-loading the dictionary buffer into a working context before
//  *  each compression, or copying a pre-loaded dictionary's LZ4_stream_t into a
//  *  working LZ4_stream_t, this function introduces a no-copy setup mechanism,
//  *  in which the working stream references the dictionary stream in-place.
//  *
//  *  Several assumptions are made about the state of the dictionary stream.
//  *  Currently, only streams which have been prepared by LZ4_loadDict() should
//  *  be expected to work.
//  *
//  *  Alternatively, the provided dictionaryStream may be NULL,
//  *  in which case any existing dictionary stream is unset.
//  *
//  *  If a dictionary is provided, it replaces any pre-existing stream history.
//  *  The dictionary contents are the only history that can be referenced and
//  *  logically immediately precede the data compressed in the first subsequent
//  *  compression call.
//  *
//  *  The dictionary will only remain attached to the working stream through the
//  *  first compression call, at the end of which it is cleared. The dictionary
//  *  stream (and source buffer) must remain in-place / accessible / unchanged
//  *  through the completion of the first compression call on the stream.
//  */
// LZ4LIB_STATIC_API void LZ4_attach_dictionary(LZ4_stream_t* workingStream, const LZ4_stream_t* dictionaryStream);
// 
// 
// /*! In-place compression and decompression
//  *
//  * It's possible to have input and output sharing the same buffer,
//  * for highly contrained memory environments.
//  * In both cases, it requires input to lay at the end of the buffer,
//  * and decompression to start at beginning of the buffer.
//  * Buffer size must feature some margin, hence be larger than final size.
//  *
//  * |<------------------------buffer--------------------------------->|
//  *                             |<-----------compressed data--------->|
//  * |<-----------decompressed size------------------>|
//  *                                                  |<----margin---->|
//  *
//  * This technique is more useful for decompression,
//  * since decompressed size is typically larger,
//  * and margin is short.
//  *
//  * In-place decompression will work inside any buffer
//  * which size is >= LZ4_DECOMPRESS_INPLACE_BUFFER_SIZE(decompressedSize).
//  * This presumes that decompressedSize > compressedSize.
//  * Otherwise, it means compression actually expanded data,
//  * and it would be more efficient to store such data with a flag indicating it's not compressed.
//  * This can happen when data is not compressible (already compressed, or encrypted).
//  *
//  * For in-place compression, margin is larger, as it must be able to cope with both
//  * history preservation, requiring input data to remain unmodified up to LZ4_DISTANCE_MAX,
//  * and data expansion, which can happen when input is not compressible.
//  * As a consequence, buffer size requirements are much higher,
//  * and memory savings offered by in-place compression are more limited.
//  *
//  * There are ways to limit this cost for compression :
//  * - Reduce history size, by modifying LZ4_DISTANCE_MAX.
//  *   Note that it is a compile-time constant, so all compressions will apply this limit.
//  *   Lower values will reduce compression ratio, except when input_size < LZ4_DISTANCE_MAX,
//  *   so it's a reasonable trick when inputs are known to be small.
//  * - Require the compressor to deliver a "maximum compressed size".
//  *   This is the `dstCapacity` parameter in `LZ4_compress*()`.
//  *   When this size is < LZ4_COMPRESSBOUND(inputSize), then compression can fail,
//  *   in which case, the return code will be 0 (zero).
//  *   The caller must be ready for these cases to happen,
//  *   and typically design a backup scheme to send data uncompressed.
//  * The combination of both techniques can significantly reduce
//  * the amount of margin required for in-place compression.
//  *
//  * In-place compression can work in any buffer
//  * which size is >= (maxCompressedSize)
//  * with maxCompressedSize == LZ4_COMPRESSBOUND(srcSize) for guaranteed compression success.
//  * LZ4_COMPRESS_INPLACE_BUFFER_SIZE() depends on both maxCompressedSize and LZ4_DISTANCE_MAX,
//  * so it's possible to reduce memory requirements by playing with them.
//  */
// 
// #define LZ4_DECOMPRESS_INPLACE_MARGIN(compressedSize)          (((compressedSize) >> 8) + 32)
// #define LZ4_DECOMPRESS_INPLACE_BUFFER_SIZE(decompressedSize)   ((decompressedSize) + LZ4_DECOMPRESS_INPLACE_MARGIN(decompressedSize))  /**< note: presumes that compressedSize < decompressedSize. note2: margin is overestimated a bit, since it could use compressedSize instead */
// 
// #ifndef LZ4_DISTANCE_MAX   /* history window size; can be user-defined at compile time */
// #  define LZ4_DISTANCE_MAX 65535   /* set to maximum value by default */
// #endif
// 
// #define LZ4_COMPRESS_INPLACE_MARGIN                           (LZ4_DISTANCE_MAX + 32)   /* LZ4_DISTANCE_MAX can be safely replaced by srcSize when it's smaller */
// #define LZ4_COMPRESS_INPLACE_BUFFER_SIZE(maxCompressedSize)   ((maxCompressedSize) + LZ4_COMPRESS_INPLACE_MARGIN)  /**< maxCompressedSize is generally LZ4_COMPRESSBOUND(inputSize), but can be set to any lower value, with the risk that compression can fail (return code 0(zero)) */
// 
// #endif   /* LZ4_STATIC_3504398509 */
// #endif   /* LZ4_STATIC_LINKING_ONLY */
// 
// 
// 
// #ifndef LZ4_H_98237428734687
// #define LZ4_H_98237428734687
// 
// /*-************************************************************
//  *  Private Definitions
//  **************************************************************
//  * Do not use these definitions directly.
//  * They are only exposed to allow static allocation of `LZ4_stream_t` and `LZ4_streamDecode_t`.
//  * Accessing members will expose user code to API and/or ABI break in future versions of the library.
//  **************************************************************/
// #define LZ4_HASHLOG   (LZ4_MEMORY_USAGE-2)
// #define LZ4_HASHTABLESIZE (1 << LZ4_MEMORY_USAGE)
// #define LZ4_HASH_SIZE_U32 (1 << LZ4_HASHLOG)       /* required as macro for static allocation */
// 
// #if defined(__cplusplus) || (defined (__STDC_VERSION__) && (__STDC_VERSION__ >= 199901L) /* C99 */)
// # include <stdint.h>
//   typedef  int8_t  LZ4_i8;
//   typedef uint8_t  LZ4_byte;
//   typedef uint16_t LZ4_u16;
//   typedef uint32_t LZ4_u32;
// #else
//   typedef   signed char  LZ4_i8;
//   typedef unsigned char  LZ4_byte;
//   typedef unsigned short LZ4_u16;
//   typedef unsigned int   LZ4_u32;
// #endif
// 
// typedef struct LZ4_stream_t_internal LZ4_stream_t_internal;
// struct LZ4_stream_t_internal {
//     LZ4_u32 hashTable[LZ4_HASH_SIZE_U32];
//     LZ4_u32 currentOffset;
//     LZ4_u32 tableType;
//     const LZ4_byte* dictionary;
//     const LZ4_stream_t_internal* dictCtx;
//     LZ4_u32 dictSize;
// };
// 
// typedef struct {
//     const LZ4_byte* externalDict;
//     size_t extDictSize;
//     const LZ4_byte* prefixEnd;
//     size_t prefixSize;
// } LZ4_streamDecode_t_internal;
// 
// 
// /*! LZ4_stream_t :
//  *  Do not use below internal definitions directly !
//  *  Declare or allocate an LZ4_stream_t instead.
//  *  LZ4_stream_t can also be created using LZ4_createStream(), which is recommended.
//  *  The structure definition can be convenient for static allocation
//  *  (on stack, or as part of larger structure).
//  *  Init this structure with LZ4_initStream() before first use.
//  *  note : only use this definition in association with static linking !
//  *  this definition is not API/ABI safe, and may change in future versions.
//  */
// #define LZ4_STREAMSIZE       16416  /* static size, for inter-version compatibility */
// #define LZ4_STREAMSIZE_VOIDP (LZ4_STREAMSIZE / sizeof(void*))
// union LZ4_stream_u {
//     void* table[LZ4_STREAMSIZE_VOIDP];
//     LZ4_stream_t_internal internal_donotuse;
// }; /* previously typedef'd to LZ4_stream_t */
// 
// 
// /*! LZ4_initStream() : v1.9.0+
//  *  An LZ4_stream_t structure must be initialized at least once.
//  *  This is automatically done when invoking LZ4_createStream(),
//  *  but it's not when the structure is simply declared on stack (for example).
//  *
//  *  Use LZ4_initStream() to properly initialize a newly declared LZ4_stream_t.
//  *  It can also initialize any arbitrary buffer of sufficient size,
//  *  and will @return a pointer of proper type upon initialization.
//  *
//  *  Note : initialization fails if size and alignment conditions are not respected.
//  *         In which case, the function will @return NULL.
//  *  Note2: An LZ4_stream_t structure guarantees correct alignment and size.
//  *  Note3: Before v1.9.0, use LZ4_resetStream() instead
//  */
// LZ4LIB_API LZ4_stream_t* LZ4_initStream (void* buffer, size_t size);
// 
// 
// /*! LZ4_streamDecode_t :
//  *  information structure to track an LZ4 stream during decompression.
//  *  init this structure  using LZ4_setStreamDecode() before first use.
//  *  note : only use in association with static linking !
//  *         this definition is not API/ABI safe,
//  *         and may change in a future version !
//  */
// #define LZ4_STREAMDECODESIZE_U64 (4 + ((sizeof(void*)==16) ? 2 : 0) /*AS-400*/ )
// #define LZ4_STREAMDECODESIZE     (LZ4_STREAMDECODESIZE_U64 * sizeof(unsigned long long))
// union LZ4_streamDecode_u {
//     unsigned long long table[LZ4_STREAMDECODESIZE_U64];
//     LZ4_streamDecode_t_internal internal_donotuse;
// } ;   /* previously typedef'd to LZ4_streamDecode_t */
// 
// 
// 
// /*-************************************
// *  Obsolete Functions
// **************************************/
// 
// /*! Deprecation warnings
//  *
//  *  Deprecated functions make the compiler generate a warning when invoked.
//  *  This is meant to invite users to update their source code.
//  *  Should deprecation warnings be a problem, it is generally possible to disable them,
//  *  typically with -Wno-deprecated-declarations for gcc
//  *  or _CRT_SECURE_NO_WARNINGS in Visual.
//  *
//  *  Another method is to define LZ4_DISABLE_DEPRECATE_WARNINGS
//  *  before including the header file.
//  */
// #ifdef LZ4_DISABLE_DEPRECATE_WARNINGS
// #  define LZ4_DEPRECATED(message)   /* disable deprecation warnings */
// #else
// #  if defined (__cplusplus) && (__cplusplus >= 201402) /* C++14 or greater */
// #    define LZ4_DEPRECATED(message) [[deprecated(message)]]
// #  elif defined(_MSC_VER)
// #    define LZ4_DEPRECATED(message) __declspec(deprecated(message))
// #  elif defined(__clang__) || (defined(__GNUC__) && (__GNUC__ * 10 + __GNUC_MINOR__ >= 45))
// #    define LZ4_DEPRECATED(message) __attribute__((deprecated(message)))
// #  elif defined(__GNUC__) && (__GNUC__ * 10 + __GNUC_MINOR__ >= 31)
// #    define LZ4_DEPRECATED(message) __attribute__((deprecated))
// #  else
// #    pragma message("WARNING: LZ4_DEPRECATED needs custom implementation for this compiler")
// #    define LZ4_DEPRECATED(message)   /* disabled */
// #  endif
// #endif /* LZ4_DISABLE_DEPRECATE_WARNINGS */
// 
// /*! Obsolete compression functions (since v1.7.3) */
// LZ4_DEPRECATED("use LZ4_compress_default() instead")       LZ4LIB_API int LZ4_compress               (const char* src, char* dest, int srcSize);
// LZ4_DEPRECATED("use LZ4_compress_default() instead")       LZ4LIB_API int LZ4_compress_limitedOutput (const char* src, char* dest, int srcSize, int maxOutputSize);
// LZ4_DEPRECATED("use LZ4_compress_fast_extState() instead") LZ4LIB_API int LZ4_compress_withState               (void* state, const char* source, char* dest, int inputSize);
// LZ4_DEPRECATED("use LZ4_compress_fast_extState() instead") LZ4LIB_API int LZ4_compress_limitedOutput_withState (void* state, const char* source, char* dest, int inputSize, int maxOutputSize);
// LZ4_DEPRECATED("use LZ4_compress_fast_continue() instead") LZ4LIB_API int LZ4_compress_continue                (LZ4_stream_t* LZ4_streamPtr, const char* source, char* dest, int inputSize);
// LZ4_DEPRECATED("use LZ4_compress_fast_continue() instead") LZ4LIB_API int LZ4_compress_limitedOutput_continue  (LZ4_stream_t* LZ4_streamPtr, const char* source, char* dest, int inputSize, int maxOutputSize);
// 
// /*! Obsolete decompression functions (since v1.8.0) */
// LZ4_DEPRECATED("use LZ4_decompress_fast() instead") LZ4LIB_API int LZ4_uncompress (const char* source, char* dest, int outputSize);
// LZ4_DEPRECATED("use LZ4_decompress_safe() instead") LZ4LIB_API int LZ4_uncompress_unknownOutputSize (const char* source, char* dest, int isize, int maxOutputSize);
// 
// /* Obsolete streaming functions (since v1.7.0)
//  * degraded functionality; do not use!
//  *
//  * In order to perform streaming compression, these functions depended on data
//  * that is no longer tracked in the state. They have been preserved as well as
//  * possible: using them will still produce a correct output. However, they don't
//  * actually retain any history between compression calls. The compression ratio
//  * achieved will therefore be no better than compressing each chunk
//  * independently.
//  */
// LZ4_DEPRECATED("Use LZ4_createStream() instead") LZ4LIB_API void* LZ4_create (char* inputBuffer);
// LZ4_DEPRECATED("Use LZ4_createStream() instead") LZ4LIB_API int   LZ4_sizeofStreamState(void);
// LZ4_DEPRECATED("Use LZ4_resetStream() instead")  LZ4LIB_API int   LZ4_resetStreamState(void* state, char* inputBuffer);
// LZ4_DEPRECATED("Use LZ4_saveDict() instead")     LZ4LIB_API char* LZ4_slideInputBuffer (void* state);
// 
// /*! Obsolete streaming decoding functions (since v1.7.0) */
// LZ4_DEPRECATED("use LZ4_decompress_safe_usingDict() instead") LZ4LIB_API int LZ4_decompress_safe_withPrefix64k (const char* src, char* dst, int compressedSize, int maxDstSize);
// LZ4_DEPRECATED("use LZ4_decompress_fast_usingDict() instead") LZ4LIB_API int LZ4_decompress_fast_withPrefix64k (const char* src, char* dst, int originalSize);
// 
// /*! Obsolete LZ4_decompress_fast variants (since v1.9.0) :
//  *  These functions used to be faster than LZ4_decompress_safe(),
//  *  but this is no longer the case. They are now slower.
//  *  This is because LZ4_decompress_fast() doesn't know the input size,
//  *  and therefore must progress more cautiously into the input buffer to not read beyond the end of block.
//  *  On top of that `LZ4_decompress_fast()` is not protected vs malformed or malicious inputs, making it a security liability.
//  *  As a consequence, LZ4_decompress_fast() is strongly discouraged, and deprecated.
//  *
//  *  The last remaining LZ4_decompress_fast() specificity is that
//  *  it can decompress a block without knowing its compressed size.
//  *  Such functionality can be achieved in a more secure manner
//  *  by employing LZ4_decompress_safe_partial().
//  *
//  *  Parameters:
//  *  originalSize : is the uncompressed size to regenerate.
//  *                 `dst` must be already allocated, its size must be >= 'originalSize' bytes.
//  * @return : number of bytes read from source buffer (== compressed size).
//  *           The function expects to finish at block's end exactly.
//  *           If the source stream is detected malformed, the function stops decoding and returns a negative result.
//  *  note : LZ4_decompress_fast*() requires originalSize. Thanks to this information, it never writes past the output buffer.
//  *         However, since it doesn't know its 'src' size, it may read an unknown amount of input, past input buffer bounds.
//  *         Also, since match offsets are not validated, match reads from 'src' may underflow too.
//  *         These issues never happen if input (compressed) data is correct.
//  *         But they may happen if input data is invalid (error or intentional tampering).
//  *         As a consequence, use these functions in trusted environments with trusted data **only**.
//  */
// LZ4_DEPRECATED("This function is deprecated and unsafe. Consider using LZ4_decompress_safe() instead")
// LZ4LIB_API int LZ4_decompress_fast (const char* src, char* dst, int originalSize);
// LZ4_DEPRECATED("This function is deprecated and unsafe. Consider using LZ4_decompress_safe_continue() instead")
// LZ4LIB_API int LZ4_decompress_fast_continue (LZ4_streamDecode_t* LZ4_streamDecode, const char* src, char* dst, int originalSize);
// LZ4_DEPRECATED("This function is deprecated and unsafe. Consider using LZ4_decompress_safe_usingDict() instead")
// LZ4LIB_API int LZ4_decompress_fast_usingDict (const char* src, char* dst, int originalSize, const char* dictStart, int dictSize);
// 
// /*! LZ4_resetStream() :
//  *  An LZ4_stream_t structure must be initialized at least once.
//  *  This is done with LZ4_initStream(), or LZ4_resetStream().
//  *  Consider switching to LZ4_initStream(),
//  *  invoking LZ4_resetStream() will trigger deprecation warnings in the future.
//  */
// LZ4LIB_API void LZ4_resetStream (LZ4_stream_t* streamPtr);
// 
// 
// #endif /* LZ4_H_98237428734687 */
// 
// 
// #if defined (__cplusplus)
// }
// #endif
// /* see also "memory routines" below */
// 
// 
// /*-************************************
// *  Compiler Options
// **************************************/
// #if defined(_MSC_VER) && (_MSC_VER >= 1400)  /* Visual Studio 2005+ */
// #  include <intrin.h>               /* only present in VS2005+ */
// #  pragma warning(disable : 4127)   /* disable: C4127: conditional expression is constant */
// #endif  /* _MSC_VER */
// 
// #ifndef LZ4_FORCE_INLINE
// #  ifdef _MSC_VER    /* Visual Studio */
// #    define LZ4_FORCE_INLINE static __forceinline
// #  else
// #    if defined (__cplusplus) || defined (__STDC_VERSION__) && __STDC_VERSION__ >= 199901L   /* C99 */
// #      ifdef __GNUC__
// #        define LZ4_FORCE_INLINE static inline __attribute__((always_inline))
// #      else
// #        define LZ4_FORCE_INLINE static inline
// #      endif
// #    else
// #      define LZ4_FORCE_INLINE static
// #    endif /* __STDC_VERSION__ */
// #  endif  /* _MSC_VER */
// #endif /* LZ4_FORCE_INLINE */
// 
// /* LZ4_FORCE_O2 and LZ4_FORCE_INLINE
//  * gcc on ppc64le generates an unrolled SIMDized loop for LZ4_wildCopy8,
//  * together with a simple 8-byte copy loop as a fall-back path.
//  * However, this optimization hurts the decompression speed by >30%,
//  * because the execution does not go to the optimized loop
//  * for typical compressible data, and all of the preamble checks
//  * before going to the fall-back path become useless overhead.
//  * This optimization happens only with the -O3 flag, and -O2 generates
//  * a simple 8-byte copy loop.
//  * With gcc on ppc64le, all of the LZ4_decompress_* and LZ4_wildCopy8
//  * functions are annotated with __attribute__((optimize("O2"))),
//  * and also LZ4_wildCopy8 is forcibly inlined, so that the O2 attribute
//  * of LZ4_wildCopy8 does not affect the compression speed.
//  */
// #if defined(__PPC64__) && defined(__LITTLE_ENDIAN__) && defined(__GNUC__) && !defined(__clang__)
// #  define LZ4_FORCE_O2  __attribute__((optimize("O2")))
// #  undef LZ4_FORCE_INLINE
// #  define LZ4_FORCE_INLINE  static __inline __attribute__((optimize("O2"),always_inline))
// #else
// #  define LZ4_FORCE_O2
// #endif
// 
// #if (defined(__GNUC__) && (__GNUC__ >= 3)) || (defined(__INTEL_COMPILER) && (__INTEL_COMPILER >= 800)) || defined(__clang__)
// #  define expect(expr,value)    (__builtin_expect ((expr),(value)) )
// #else
// #  define expect(expr,value)    (expr)
// #endif
// 
// #ifndef likely
// #define likely(expr)     expect((expr) != 0, 1)
// #endif
// #ifndef unlikely
// #define unlikely(expr)   expect((expr) != 0, 0)
// #endif
// 
// /* Should the alignment test prove unreliable, for some reason,
//  * it can be disabled by setting LZ4_ALIGN_TEST to 0 */
// #ifndef LZ4_ALIGN_TEST  /* can be externally provided */
// # define LZ4_ALIGN_TEST 1
// #endif
// 
// 
// /*-************************************
// *  Memory routines
// **************************************/
// #ifdef LZ4_USER_MEMORY_FUNCTIONS
// /* memory management functions can be customized by user project.
//  * Below functions must exist somewhere in the Project
//  * and be available at link time */
// void* LZ4_malloc(size_t s);
// void* LZ4_calloc(size_t n, size_t s);
// void  LZ4_free(void* p);
// # define ALLOC(s)          LZ4_malloc(s)
// # define ALLOC_AND_ZERO(s) LZ4_calloc(1,s)
// # define FREEMEM(p)        LZ4_free(p)
// #else
// # include <stdlib.h>   /* malloc, calloc, free */
// # define ALLOC(s)          malloc(s)
// # define ALLOC_AND_ZERO(s) calloc(1,s)
// # define FREEMEM(p)        free(p)
// #endif
// 
// #include <string.h>   /* memset, memcpy */
// #define MEM_INIT(p,v,s)   memset((p),(v),(s))
// 
// 
// /*-************************************
// *  Common Constants
// **************************************/
// #define MINMATCH 4
// 
// #define WILDCOPYLENGTH 8
// #define LASTLITERALS   5   /* see ../doc/lz4_Block_format.md#parsing-restrictions */
// #define MFLIMIT       12   /* see ../doc/lz4_Block_format.md#parsing-restrictions */
// #define MATCH_SAFEGUARD_DISTANCE  ((2*WILDCOPYLENGTH) - MINMATCH)   /* ensure it's possible to write 2 x wildcopyLength without overflowing output buffer */
// #define FASTLOOP_SAFE_DISTANCE 64
// static const int LZ4_minLength = (MFLIMIT+1);
// 
// #define KB *(1 <<10)
// #define MB *(1 <<20)
// #define GB *(1U<<30)
// 
// #define LZ4_DISTANCE_ABSOLUTE_MAX 65535
// #if (LZ4_DISTANCE_MAX > LZ4_DISTANCE_ABSOLUTE_MAX)   /* max supported by LZ4 format */
// #  error "LZ4_DISTANCE_MAX is too big : must be <= 65535"
// #endif
// 
// #define ML_BITS  4
// #define ML_MASK  ((1U<<ML_BITS)-1)
// #define RUN_BITS (8-ML_BITS)
// #define RUN_MASK ((1U<<RUN_BITS)-1)
// 
// 
// /*-************************************
// *  Error detection
// **************************************/
// #if defined(LZ4_DEBUG) && (LZ4_DEBUG>=1)
// #  include <assert.h>
// #else
// #  ifndef assert
// #    define assert(condition) ((void)0)
// #  endif
// #endif
// 
// #define LZ4_STATIC_ASSERT(c)   { enum { LZ4_static_assert = 1/(int)(!!(c)) }; }   /* use after variable declarations */
// 
// #if defined(LZ4_DEBUG) && (LZ4_DEBUG>=2)
// #  include <stdio.h>
//    static int g_debuglog_enable = 1;
// #  define DEBUGLOG(l, ...) {                          \
//         if ((g_debuglog_enable) && (l<=LZ4_DEBUG)) {  \
//             fprintf(stderr, __FILE__ ": ");           \
//             fprintf(stderr, __VA_ARGS__);             \
//             fprintf(stderr, " \n");                   \
//     }   }
// #else
// #  define DEBUGLOG(l, ...) {}    /* disabled */
// #endif
// 
// static int LZ4_isAligned(const void* ptr, size_t alignment)
// {
//     return ((size_t)ptr & (alignment -1)) == 0;
// }
// 
// 
// /*-************************************
// *  Types
// **************************************/
// #include <limits.h>
// #if defined(__cplusplus) || (defined (__STDC_VERSION__) && (__STDC_VERSION__ >= 199901L) /* C99 */)
// # include <stdint.h>
//   typedef  uint8_t BYTE;
//   typedef uint16_t U16;
//   typedef uint32_t U32;
//   typedef  int32_t S32;
//   typedef uint64_t U64;
//   typedef uintptr_t uptrval;
// #else
// # if UINT_MAX != 4294967295UL
// #   error "LZ4 code (when not C++ or C99) assumes that sizeof(int) == 4"
// # endif
//   typedef unsigned char       BYTE;
//   typedef unsigned short      U16;
//   typedef unsigned int        U32;
//   typedef   signed int        S32;
//   typedef unsigned long long  U64;
//   typedef size_t              uptrval;   /* generally true, except OpenVMS-64 */
// #endif
// 
// #if defined(__x86_64__)
//   typedef U64    reg_t;   /* 64-bits in x32 mode */
// #else
//   typedef size_t reg_t;   /* 32-bits in x32 mode */
// #endif
// 
// typedef enum {
//     notLimited = 0,
//     limitedOutput = 1,
//     fillOutput = 2
// } limitedOutput_directive;
// 
// 
// /*-************************************
// *  Reading and writing into memory
// **************************************/
// 
// /**
//  * LZ4 relies on memcpy with a constant size being inlined. In freestanding
//  * environments, the compiler can't assume the implementation of memcpy() is
//  * standard compliant, so it can't apply its specialized memcpy() inlining
//  * logic. When possible, use __builtin_memcpy() to tell the compiler to analyze
//  * memcpy() as if it were standard compliant, so it can inline it in freestanding
//  * environments. This is needed when decompressing the Linux Kernel, for example.
//  */
// #if defined(__GNUC__) && (__GNUC__ >= 4)
// #define LZ4_memcpy(dst, src, size) __builtin_memcpy(dst, src, size)
// #else
// #define LZ4_memcpy(dst, src, size) memcpy(dst, src, size)
// #endif
// 
// static unsigned LZ4_isLittleEndian(void)
// {
//     const union { U32 u; BYTE c[4]; } one = { 1 };   /* don't use static : performance detrimental */
//     return one.c[0];
// }
// 
// 
// #if defined(LZ4_FORCE_MEMORY_ACCESS) && (LZ4_FORCE_MEMORY_ACCESS==2)
// /* lie to the compiler about data alignment; use with caution */
// 
// static U16 LZ4_read16(const void* memPtr) { return *(const U16*) memPtr; }
// static U32 LZ4_read32(const void* memPtr) { return *(const U32*) memPtr; }
// static reg_t LZ4_read_ARCH(const void* memPtr) { return *(const reg_t*) memPtr; }
// 
// static void LZ4_write16(void* memPtr, U16 value) { *(U16*)memPtr = value; }
// static void LZ4_write32(void* memPtr, U32 value) { *(U32*)memPtr = value; }
// 
// #elif defined(LZ4_FORCE_MEMORY_ACCESS) && (LZ4_FORCE_MEMORY_ACCESS==1)
// 
// /* __pack instructions are safer, but compiler specific, hence potentially problematic for some compilers */
// /* currently only defined for gcc and icc */
// typedef union { U16 u16; U32 u32; reg_t uArch; } __attribute__((packed)) unalign;
// 
// static U16 LZ4_read16(const void* ptr) { return ((const unalign*)ptr)->u16; }
// static U32 LZ4_read32(const void* ptr) { return ((const unalign*)ptr)->u32; }
// static reg_t LZ4_read_ARCH(const void* ptr) { return ((const unalign*)ptr)->uArch; }
// 
// static void LZ4_write16(void* memPtr, U16 value) { ((unalign*)memPtr)->u16 = value; }
// static void LZ4_write32(void* memPtr, U32 value) { ((unalign*)memPtr)->u32 = value; }
// 
// #else  /* safe and portable access using memcpy() */
// 
// static U16 LZ4_read16(const void* memPtr)
// {
//     U16 val; LZ4_memcpy(&val, memPtr, sizeof(val)); return val;
// }
// 
// static U32 LZ4_read32(const void* memPtr)
// {
//     U32 val; LZ4_memcpy(&val, memPtr, sizeof(val)); return val;
// }
// 
// static reg_t LZ4_read_ARCH(const void* memPtr)
// {
//     reg_t val; LZ4_memcpy(&val, memPtr, sizeof(val)); return val;
// }
// 
// static void LZ4_write16(void* memPtr, U16 value)
// {
//     LZ4_memcpy(memPtr, &value, sizeof(value));
// }
// 
// static void LZ4_write32(void* memPtr, U32 value)
// {
//     LZ4_memcpy(memPtr, &value, sizeof(value));
// }
// 
// #endif /* LZ4_FORCE_MEMORY_ACCESS */
// 
// 
// static U16 LZ4_readLE16(const void* memPtr)
// {
//     if (LZ4_isLittleEndian()) {
//         return LZ4_read16(memPtr);
//     } else {
//         const BYTE* p = (const BYTE*)memPtr;
//         return (U16)((U16)p[0] + (p[1]<<8));
//     }
// }
// 
// static void LZ4_writeLE16(void* memPtr, U16 value)
// {
//     if (LZ4_isLittleEndian()) {
//         LZ4_write16(memPtr, value);
//     } else {
//         BYTE* p = (BYTE*)memPtr;
//         p[0] = (BYTE) value;
//         p[1] = (BYTE)(value>>8);
//     }
// }
// 
// /* customized variant of memcpy, which can overwrite up to 8 bytes beyond dstEnd */
// LZ4_FORCE_INLINE
// void LZ4_wildCopy8(void* dstPtr, const void* srcPtr, void* dstEnd)
// {
//     BYTE* d = (BYTE*)dstPtr;
//     const BYTE* s = (const BYTE*)srcPtr;
//     BYTE* const e = (BYTE*)dstEnd;
// 
//     do { LZ4_memcpy(d,s,8); d+=8; s+=8; } while (d<e);
// }
// 
// static const unsigned inc32table[8] = {0, 1, 2,  1,  0,  4, 4, 4};
// static const int      dec64table[8] = {0, 0, 0, -1, -4,  1, 2, 3};
// 
// 
// #ifndef LZ4_FAST_DEC_LOOP
// #  if defined __i386__ || defined _M_IX86 || defined __x86_64__ || defined _M_X64
// #    define LZ4_FAST_DEC_LOOP 1
// #  elif defined(__aarch64__) && !defined(__clang__)
//      /* On aarch64, we disable this optimization for clang because on certain
//       * mobile chipsets, performance is reduced with clang. For information
//       * refer to https://github.com/lz4/lz4/pull/707 */
// #    define LZ4_FAST_DEC_LOOP 1
// #  else
// #    define LZ4_FAST_DEC_LOOP 0
// #  endif
// #endif
// 
// #if LZ4_FAST_DEC_LOOP
// 
// LZ4_FORCE_INLINE void
// LZ4_memcpy_using_offset_base(BYTE* dstPtr, const BYTE* srcPtr, BYTE* dstEnd, const size_t offset)
// {
//     assert(srcPtr + offset == dstPtr);
//     if (offset < 8) {
//         LZ4_write32(dstPtr, 0);   /* silence an msan warning when offset==0 */
//         dstPtr[0] = srcPtr[0];
//         dstPtr[1] = srcPtr[1];
//         dstPtr[2] = srcPtr[2];
//         dstPtr[3] = srcPtr[3];
//         srcPtr += inc32table[offset];
//         LZ4_memcpy(dstPtr+4, srcPtr, 4);
//         srcPtr -= dec64table[offset];
//         dstPtr += 8;
//     } else {
//         LZ4_memcpy(dstPtr, srcPtr, 8);
//         dstPtr += 8;
//         srcPtr += 8;
//     }
// 
//     LZ4_wildCopy8(dstPtr, srcPtr, dstEnd);
// }
// 
// /* customized variant of memcpy, which can overwrite up to 32 bytes beyond dstEnd
//  * this version copies two times 16 bytes (instead of one time 32 bytes)
//  * because it must be compatible with offsets >= 16. */
// LZ4_FORCE_INLINE void
// LZ4_wildCopy32(void* dstPtr, const void* srcPtr, void* dstEnd)
// {
//     BYTE* d = (BYTE*)dstPtr;
//     const BYTE* s = (const BYTE*)srcPtr;
//     BYTE* const e = (BYTE*)dstEnd;
// 
//     do { LZ4_memcpy(d,s,16); LZ4_memcpy(d+16,s+16,16); d+=32; s+=32; } while (d<e);
// }
// 
// /* LZ4_memcpy_using_offset()  presumes :
//  * - dstEnd >= dstPtr + MINMATCH
//  * - there is at least 8 bytes available to write after dstEnd */
// LZ4_FORCE_INLINE void
// LZ4_memcpy_using_offset(BYTE* dstPtr, const BYTE* srcPtr, BYTE* dstEnd, const size_t offset)
// {
//     BYTE v[8];
// 
//     assert(dstEnd >= dstPtr + MINMATCH);
// 
//     switch(offset) {
//     case 1:
//         MEM_INIT(v, *srcPtr, 8);
//         break;
//     case 2:
//         LZ4_memcpy(v, srcPtr, 2);
//         LZ4_memcpy(&v[2], srcPtr, 2);
//         LZ4_memcpy(&v[4], v, 4);
//         break;
//     case 4:
//         LZ4_memcpy(v, srcPtr, 4);
//         LZ4_memcpy(&v[4], srcPtr, 4);
//         break;
//     default:
//         LZ4_memcpy_using_offset_base(dstPtr, srcPtr, dstEnd, offset);
//         return;
//     }
// 
//     LZ4_memcpy(dstPtr, v, 8);
//     dstPtr += 8;
//     while (dstPtr < dstEnd) {
//         LZ4_memcpy(dstPtr, v, 8);
//         dstPtr += 8;
//     }
// }
// #endif
// 
// 
// /*-************************************
// *  Common functions
// **************************************/
// static unsigned LZ4_NbCommonBytes (reg_t val)
// {
//     assert(val != 0);
//     if (LZ4_isLittleEndian()) {
//         if (sizeof(val) == 8) {
// #       if defined(_MSC_VER) && (_MSC_VER >= 1800) && defined(_M_AMD64) && !defined(LZ4_FORCE_SW_BITCOUNT)
//             /* x64 CPUS without BMI support interpret `TZCNT` as `REP BSF` */
//             return (unsigned)_tzcnt_u64(val) >> 3;
// #       elif defined(_MSC_VER) && defined(_WIN64) && !defined(LZ4_FORCE_SW_BITCOUNT)
//             unsigned long r = 0;
//             _BitScanForward64(&r, (U64)val);
//             return (unsigned)r >> 3;
// #       elif (defined(__clang__) || (defined(__GNUC__) && ((__GNUC__ > 3) || \
//                             ((__GNUC__ == 3) && (__GNUC_MINOR__ >= 4))))) && \
//                                         !defined(LZ4_FORCE_SW_BITCOUNT)
//             return (unsigned)__builtin_ctzll((U64)val) >> 3;
// #       else
//             const U64 m = 0x0101010101010101ULL;
//             val ^= val - 1;
//             return (unsigned)(((U64)((val & (m - 1)) * m)) >> 56);
// #       endif
//         } else /* 32 bits */ {
// #       if defined(_MSC_VER) && (_MSC_VER >= 1400) && !defined(LZ4_FORCE_SW_BITCOUNT)
//             unsigned long r;
//             _BitScanForward(&r, (U32)val);
//             return (unsigned)r >> 3;
// #       elif (defined(__clang__) || (defined(__GNUC__) && ((__GNUC__ > 3) || \
//                             ((__GNUC__ == 3) && (__GNUC_MINOR__ >= 4))))) && \
//                         !defined(__TINYC__) && !defined(LZ4_FORCE_SW_BITCOUNT)
//             return (unsigned)__builtin_ctz((U32)val) >> 3;
// #       else
//             const U32 m = 0x01010101;
//             return (unsigned)((((val - 1) ^ val) & (m - 1)) * m) >> 24;
// #       endif
//         }
//     } else   /* Big Endian CPU */ {
//         if (sizeof(val)==8) {
// #       if (defined(__clang__) || (defined(__GNUC__) && ((__GNUC__ > 3) || \
//                             ((__GNUC__ == 3) && (__GNUC_MINOR__ >= 4))))) && \
//                         !defined(__TINYC__) && !defined(LZ4_FORCE_SW_BITCOUNT)
//             return (unsigned)__builtin_clzll((U64)val) >> 3;
// #       else
// #if 1
//             /* this method is probably faster,
//              * but adds a 128 bytes lookup table */
//             static const unsigned char ctz7_tab[128] = {
//                 7, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
//                 4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
//                 5, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
//                 4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
//                 6, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
//                 4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
//                 5, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
//                 4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
//             };
//             U64 const mask = 0x0101010101010101ULL;
//             U64 const t = (((val >> 8) - mask) | val) & mask;
//             return ctz7_tab[(t * 0x0080402010080402ULL) >> 57];
// #else
//             /* this method doesn't consume memory space like the previous one,
//              * but it contains several branches,
//              * that may end up slowing execution */
//             static const U32 by32 = sizeof(val)*4;  /* 32 on 64 bits (goal), 16 on 32 bits.
//             Just to avoid some static analyzer complaining about shift by 32 on 32-bits target.
//             Note that this code path is never triggered in 32-bits mode. */
//             unsigned r;
//             if (!(val>>by32)) { r=4; } else { r=0; val>>=by32; }
//             if (!(val>>16)) { r+=2; val>>=8; } else { val>>=24; }
//             r += (!val);
//             return r;
// #endif
// #       endif
//         } else /* 32 bits */ {
// #       if (defined(__clang__) || (defined(__GNUC__) && ((__GNUC__ > 3) || \
//                             ((__GNUC__ == 3) && (__GNUC_MINOR__ >= 4))))) && \
//                                         !defined(LZ4_FORCE_SW_BITCOUNT)
//             return (unsigned)__builtin_clz((U32)val) >> 3;
// #       else
//             val >>= 8;
//             val = ((((val + 0x00FFFF00) | 0x00FFFFFF) + val) |
//               (val + 0x00FF0000)) >> 24;
//             return (unsigned)val ^ 3;
// #       endif
//         }
//     }
// }
// 
// 
// #define STEPSIZE sizeof(reg_t)
// LZ4_FORCE_INLINE
// unsigned LZ4_count(const BYTE* pIn, const BYTE* pMatch, const BYTE* pInLimit)
// {
//     const BYTE* const pStart = pIn;
// 
//     if (likely(pIn < pInLimit-(STEPSIZE-1))) {
//         reg_t const diff = LZ4_read_ARCH(pMatch) ^ LZ4_read_ARCH(pIn);
//         if (!diff) {
//             pIn+=STEPSIZE; pMatch+=STEPSIZE;
//         } else {
//             return LZ4_NbCommonBytes(diff);
//     }   }
// 
//     while (likely(pIn < pInLimit-(STEPSIZE-1))) {
//         reg_t const diff = LZ4_read_ARCH(pMatch) ^ LZ4_read_ARCH(pIn);
//         if (!diff) { pIn+=STEPSIZE; pMatch+=STEPSIZE; continue; }
//         pIn += LZ4_NbCommonBytes(diff);
//         return (unsigned)(pIn - pStart);
//     }
// 
//     if ((STEPSIZE==8) && (pIn<(pInLimit-3)) && (LZ4_read32(pMatch) == LZ4_read32(pIn))) { pIn+=4; pMatch+=4; }
//     if ((pIn<(pInLimit-1)) && (LZ4_read16(pMatch) == LZ4_read16(pIn))) { pIn+=2; pMatch+=2; }
//     if ((pIn<pInLimit) && (*pMatch == *pIn)) pIn++;
//     return (unsigned)(pIn - pStart);
// }
// 
// 
// #ifndef LZ4_COMMONDEFS_ONLY
// /*-************************************
// *  Local Constants
// **************************************/
// static const int LZ4_64Klimit = ((64 KB) + (MFLIMIT-1));
// static const U32 LZ4_skipTrigger = 6;  /* Increase this value ==> compression run slower on incompressible data */
// 
// 
// /*-************************************
// *  Local Structures and types
// **************************************/
// typedef enum { clearedTable = 0, byPtr, byU32, byU16 } tableType_t;
// 
// /**
//  * This enum distinguishes several different modes of accessing previous
//  * content in the stream.
//  *
//  * - noDict        : There is no preceding content.
//  * - withPrefix64k : Table entries up to ctx->dictSize before the current blob
//  *                   blob being compressed are valid and refer to the preceding
//  *                   content (of length ctx->dictSize), which is available
//  *                   contiguously preceding in memory the content currently
//  *                   being compressed.
//  * - usingExtDict  : Like withPrefix64k, but the preceding content is somewhere
//  *                   else in memory, starting at ctx->dictionary with length
//  *                   ctx->dictSize.
//  * - usingDictCtx  : Like usingExtDict, but everything concerning the preceding
//  *                   content is in a separate context, pointed to by
//  *                   ctx->dictCtx. ctx->dictionary, ctx->dictSize, and table
//  *                   entries in the current context that refer to positions
//  *                   preceding the beginning of the current compression are
//  *                   ignored. Instead, ctx->dictCtx->dictionary and ctx->dictCtx
//  *                   ->dictSize describe the location and size of the preceding
//  *                   content, and matches are found by looking in the ctx
//  *                   ->dictCtx->hashTable.
//  */
// typedef enum { noDict = 0, withPrefix64k, usingExtDict, usingDictCtx } dict_directive;
// typedef enum { noDictIssue = 0, dictSmall } dictIssue_directive;
// 
// 
// /*-************************************
// *  Local Utils
// **************************************/
// int LZ4_versionNumber (void) { return LZ4_VERSION_NUMBER; }
// const char* LZ4_versionString(void) { return LZ4_VERSION_STRING; }
// int LZ4_compressBound(int isize)  { return LZ4_COMPRESSBOUND(isize); }
// int LZ4_sizeofState(void) { return LZ4_STREAMSIZE; }
// 
// 
// /*-************************************
// *  Internal Definitions used in Tests
// **************************************/
// #if defined (__cplusplus)
// extern "C" {
// #endif
// 
// int LZ4_compress_forceExtDict (LZ4_stream_t* LZ4_dict, const char* source, char* dest, int srcSize);
// 
// int LZ4_decompress_safe_forceExtDict(const char* source, char* dest,
//                                      int compressedSize, int maxOutputSize,
//                                      const void* dictStart, size_t dictSize);
// 
// #if defined (__cplusplus)
// }
// #endif
// 
// /*-******************************
// *  Compression functions
// ********************************/
// LZ4_FORCE_INLINE U32 LZ4_hash4(U32 sequence, tableType_t const tableType)
// {
//     if (tableType == byU16)
//         return ((sequence * 2654435761U) >> ((MINMATCH*8)-(LZ4_HASHLOG+1)));
//     else
//         return ((sequence * 2654435761U) >> ((MINMATCH*8)-LZ4_HASHLOG));
// }
// 
// LZ4_FORCE_INLINE U32 LZ4_hash5(U64 sequence, tableType_t const tableType)
// {
//     const U32 hashLog = (tableType == byU16) ? LZ4_HASHLOG+1 : LZ4_HASHLOG;
//     if (LZ4_isLittleEndian()) {
//         const U64 prime5bytes = 889523592379ULL;
//         return (U32)(((sequence << 24) * prime5bytes) >> (64 - hashLog));
//     } else {
//         const U64 prime8bytes = 11400714785074694791ULL;
//         return (U32)(((sequence >> 24) * prime8bytes) >> (64 - hashLog));
//     }
// }
// 
// LZ4_FORCE_INLINE U32 LZ4_hashPosition(const void* const p, tableType_t const tableType)
// {
//     if ((sizeof(reg_t)==8) && (tableType != byU16)) return LZ4_hash5(LZ4_read_ARCH(p), tableType);
//     return LZ4_hash4(LZ4_read32(p), tableType);
// }
// 
// LZ4_FORCE_INLINE void LZ4_clearHash(U32 h, void* tableBase, tableType_t const tableType)
// {
//     switch (tableType)
//     {
//     default: /* fallthrough */
//     case clearedTable: { /* illegal! */ assert(0); return; }
//     case byPtr: { const BYTE** hashTable = (const BYTE**)tableBase; hashTable[h] = NULL; return; }
//     case byU32: { U32* hashTable = (U32*) tableBase; hashTable[h] = 0; return; }
//     case byU16: { U16* hashTable = (U16*) tableBase; hashTable[h] = 0; return; }
//     }
// }
// 
// LZ4_FORCE_INLINE void LZ4_putIndexOnHash(U32 idx, U32 h, void* tableBase, tableType_t const tableType)
// {
//     switch (tableType)
//     {
//     default: /* fallthrough */
//     case clearedTable: /* fallthrough */
//     case byPtr: { /* illegal! */ assert(0); return; }
//     case byU32: { U32* hashTable = (U32*) tableBase; hashTable[h] = idx; return; }
//     case byU16: { U16* hashTable = (U16*) tableBase; assert(idx < 65536); hashTable[h] = (U16)idx; return; }
//     }
// }
// 
// LZ4_FORCE_INLINE void LZ4_putPositionOnHash(const BYTE* p, U32 h,
//                                   void* tableBase, tableType_t const tableType,
//                             const BYTE* srcBase)
// {
//     switch (tableType)
//     {
//     case clearedTable: { /* illegal! */ assert(0); return; }
//     case byPtr: { const BYTE** hashTable = (const BYTE**)tableBase; hashTable[h] = p; return; }
//     case byU32: { U32* hashTable = (U32*) tableBase; hashTable[h] = (U32)(p-srcBase); return; }
//     case byU16: { U16* hashTable = (U16*) tableBase; hashTable[h] = (U16)(p-srcBase); return; }
//     }
// }
// 
// LZ4_FORCE_INLINE void LZ4_putPosition(const BYTE* p, void* tableBase, tableType_t tableType, const BYTE* srcBase)
// {
//     U32 const h = LZ4_hashPosition(p, tableType);
//     LZ4_putPositionOnHash(p, h, tableBase, tableType, srcBase);
// }
// 
// /* LZ4_getIndexOnHash() :
//  * Index of match position registered in hash table.
//  * hash position must be calculated by using base+index, or dictBase+index.
//  * Assumption 1 : only valid if tableType == byU32 or byU16.
//  * Assumption 2 : h is presumed valid (within limits of hash table)
//  */
// LZ4_FORCE_INLINE U32 LZ4_getIndexOnHash(U32 h, const void* tableBase, tableType_t tableType)
// {
//     LZ4_STATIC_ASSERT(LZ4_MEMORY_USAGE > 2);
//     if (tableType == byU32) {
//         const U32* const hashTable = (const U32*) tableBase;
//         assert(h < (1U << (LZ4_MEMORY_USAGE-2)));
//         return hashTable[h];
//     }
//     if (tableType == byU16) {
//         const U16* const hashTable = (const U16*) tableBase;
//         assert(h < (1U << (LZ4_MEMORY_USAGE-1)));
//         return hashTable[h];
//     }
//     assert(0); return 0;  /* forbidden case */
// }
// 
// static const BYTE* LZ4_getPositionOnHash(U32 h, const void* tableBase, tableType_t tableType, const BYTE* srcBase)
// {
//     if (tableType == byPtr) { const BYTE* const* hashTable = (const BYTE* const*) tableBase; return hashTable[h]; }
//     if (tableType == byU32) { const U32* const hashTable = (const U32*) tableBase; return hashTable[h] + srcBase; }
//     { const U16* const hashTable = (const U16*) tableBase; return hashTable[h] + srcBase; }   /* default, to ensure a return */
// }
// 
// LZ4_FORCE_INLINE const BYTE*
// LZ4_getPosition(const BYTE* p,
//                 const void* tableBase, tableType_t tableType,
//                 const BYTE* srcBase)
// {
//     U32 const h = LZ4_hashPosition(p, tableType);
//     return LZ4_getPositionOnHash(h, tableBase, tableType, srcBase);
// }
// 
// LZ4_FORCE_INLINE void
// LZ4_prepareTable(LZ4_stream_t_internal* const cctx,
//            const int inputSize,
//            const tableType_t tableType) {
//     /* If the table hasn't been used, it's guaranteed to be zeroed out, and is
//      * therefore safe to use no matter what mode we're in. Otherwise, we figure
//      * out if it's safe to leave as is or whether it needs to be reset.
//      */
//     if ((tableType_t)cctx->tableType != clearedTable) {
//         assert(inputSize >= 0);
//         if ((tableType_t)cctx->tableType != tableType
//           || ((tableType == byU16) && cctx->currentOffset + (unsigned)inputSize >= 0xFFFFU)
//           || ((tableType == byU32) && cctx->currentOffset > 1 GB)
//           || tableType == byPtr
//           || inputSize >= 4 KB)
//         {
//             DEBUGLOG(4, "LZ4_prepareTable: Resetting table in %p", cctx);
//             MEM_INIT(cctx->hashTable, 0, LZ4_HASHTABLESIZE);
//             cctx->currentOffset = 0;
//             cctx->tableType = (U32)clearedTable;
//         } else {
//             DEBUGLOG(4, "LZ4_prepareTable: Re-use hash table (no reset)");
//         }
//     }
// 
//     /* Adding a gap, so all previous entries are > LZ4_DISTANCE_MAX back, is faster
//      * than compressing without a gap. However, compressing with
//      * currentOffset == 0 is faster still, so we preserve that case.
//      */
//     if (cctx->currentOffset != 0 && tableType == byU32) {
//         DEBUGLOG(5, "LZ4_prepareTable: adding 64KB to currentOffset");
//         cctx->currentOffset += 64 KB;
//     }
// 
//     /* Finally, clear history */
//     cctx->dictCtx = NULL;
//     cctx->dictionary = NULL;
//     cctx->dictSize = 0;
// }
// 
// /** LZ4_compress_generic() :
//  *  inlined, to ensure branches are decided at compilation time.
//  *  Presumed already validated at this stage:
//  *  - source != NULL
//  *  - inputSize > 0
//  */
// LZ4_FORCE_INLINE int LZ4_compress_generic_validated(
//                  LZ4_stream_t_internal* const cctx,
//                  const char* const source,
//                  char* const dest,
//                  const int inputSize,
//                  int *inputConsumed, /* only written when outputDirective == fillOutput */
//                  const int maxOutputSize,
//                  const limitedOutput_directive outputDirective,
//                  const tableType_t tableType,
//                  const dict_directive dictDirective,
//                  const dictIssue_directive dictIssue,
//                  const int acceleration)
// {
//     int result;
//     const BYTE* ip = (const BYTE*) source;
// 
//     U32 const startIndex = cctx->currentOffset;
//     const BYTE* base = (const BYTE*) source - startIndex;
//     const BYTE* lowLimit;
// 
//     const LZ4_stream_t_internal* dictCtx = (const LZ4_stream_t_internal*) cctx->dictCtx;
//     const BYTE* const dictionary =
//         dictDirective == usingDictCtx ? dictCtx->dictionary : cctx->dictionary;
//     const U32 dictSize =
//         dictDirective == usingDictCtx ? dictCtx->dictSize : cctx->dictSize;
//     const U32 dictDelta = (dictDirective == usingDictCtx) ? startIndex - dictCtx->currentOffset : 0;   /* make indexes in dictCtx comparable with index in current context */
// 
//     int const maybe_extMem = (dictDirective == usingExtDict) || (dictDirective == usingDictCtx);
//     U32 const prefixIdxLimit = startIndex - dictSize;   /* used when dictDirective == dictSmall */
//     const BYTE* const dictEnd = dictionary ? dictionary + dictSize : dictionary;
//     const BYTE* anchor = (const BYTE*) source;
//     const BYTE* const iend = ip + inputSize;
//     const BYTE* const mflimitPlusOne = iend - MFLIMIT + 1;
//     const BYTE* const matchlimit = iend - LASTLITERALS;
// 
//     /* the dictCtx currentOffset is indexed on the start of the dictionary,
//      * while a dictionary in the current context precedes the currentOffset */
//     const BYTE* dictBase = !dictionary ? NULL : (dictDirective == usingDictCtx) ?
//                             dictionary + dictSize - dictCtx->currentOffset :
//                             dictionary + dictSize - startIndex;
// 
//     BYTE* op = (BYTE*) dest;
//     BYTE* const olimit = op + maxOutputSize;
// 
//     U32 offset = 0;
//     U32 forwardH;
// 
//     DEBUGLOG(5, "LZ4_compress_generic_validated: srcSize=%i, tableType=%u", inputSize, tableType);
//     assert(ip != NULL);
//     /* If init conditions are not met, we don't have to mark stream
//      * as having dirty context, since no action was taken yet */
//     if (outputDirective == fillOutput && maxOutputSize < 1) { return 0; } /* Impossible to store anything */
//     if ((tableType == byU16) && (inputSize>=LZ4_64Klimit)) { return 0; }  /* Size too large (not within 64K limit) */
//     if (tableType==byPtr) assert(dictDirective==noDict);      /* only supported use case with byPtr */
//     assert(acceleration >= 1);
// 
//     lowLimit = (const BYTE*)source - (dictDirective == withPrefix64k ? dictSize : 0);
// 
//     /* Update context state */
//     if (dictDirective == usingDictCtx) {
//         /* Subsequent linked blocks can't use the dictionary. */
//         /* Instead, they use the block we just compressed. */
//         cctx->dictCtx = NULL;
//         cctx->dictSize = (U32)inputSize;
//     } else {
//         cctx->dictSize += (U32)inputSize;
//     }
//     cctx->currentOffset += (U32)inputSize;
//     cctx->tableType = (U32)tableType;
// 
//     if (inputSize<LZ4_minLength) goto _last_literals;        /* Input too small, no compression (all literals) */
// 
//     /* First Byte */
//     LZ4_putPosition(ip, cctx->hashTable, tableType, base);
//     ip++; forwardH = LZ4_hashPosition(ip, tableType);
// 
//     /* Main Loop */
//     for ( ; ; ) {
//         const BYTE* match;
//         BYTE* token;
//         const BYTE* filledIp;
// 
//         /* Find a match */
//         if (tableType == byPtr) {
//             const BYTE* forwardIp = ip;
//             int step = 1;
//             int searchMatchNb = acceleration << LZ4_skipTrigger;
//             do {
//                 U32 const h = forwardH;
//                 ip = forwardIp;
//                 forwardIp += step;
//                 step = (searchMatchNb++ >> LZ4_skipTrigger);
// 
//                 if (unlikely(forwardIp > mflimitPlusOne)) goto _last_literals;
//                 assert(ip < mflimitPlusOne);
// 
//                 match = LZ4_getPositionOnHash(h, cctx->hashTable, tableType, base);
//                 forwardH = LZ4_hashPosition(forwardIp, tableType);
//                 LZ4_putPositionOnHash(ip, h, cctx->hashTable, tableType, base);
// 
//             } while ( (match+LZ4_DISTANCE_MAX < ip)
//                    || (LZ4_read32(match) != LZ4_read32(ip)) );
// 
//         } else {   /* byU32, byU16 */
// 
//             const BYTE* forwardIp = ip;
//             int step = 1;
//             int searchMatchNb = acceleration << LZ4_skipTrigger;
//             do {
//                 U32 const h = forwardH;
//                 U32 const current = (U32)(forwardIp - base);
//                 U32 matchIndex = LZ4_getIndexOnHash(h, cctx->hashTable, tableType);
//                 assert(matchIndex <= current);
//                 assert(forwardIp - base < (ptrdiff_t)(2 GB - 1));
//                 ip = forwardIp;
//                 forwardIp += step;
//                 step = (searchMatchNb++ >> LZ4_skipTrigger);
// 
//                 if (unlikely(forwardIp > mflimitPlusOne)) goto _last_literals;
//                 assert(ip < mflimitPlusOne);
// 
//                 if (dictDirective == usingDictCtx) {
//                     if (matchIndex < startIndex) {
//                         /* there was no match, try the dictionary */
//                         assert(tableType == byU32);
//                         matchIndex = LZ4_getIndexOnHash(h, dictCtx->hashTable, byU32);
//                         match = dictBase + matchIndex;
//                         matchIndex += dictDelta;   /* make dictCtx index comparable with current context */
//                         lowLimit = dictionary;
//                     } else {
//                         match = base + matchIndex;
//                         lowLimit = (const BYTE*)source;
//                     }
//                 } else if (dictDirective==usingExtDict) {
//                     if (matchIndex < startIndex) {
//                         DEBUGLOG(7, "extDict candidate: matchIndex=%5u  <  startIndex=%5u", matchIndex, startIndex);
//                         assert(startIndex - matchIndex >= MINMATCH);
//                         match = dictBase + matchIndex;
//                         lowLimit = dictionary;
//                     } else {
//                         match = base + matchIndex;
//                         lowLimit = (const BYTE*)source;
//                     }
//                 } else {   /* single continuous memory segment */
//                     match = base + matchIndex;
//                 }
//                 forwardH = LZ4_hashPosition(forwardIp, tableType);
//                 LZ4_putIndexOnHash(current, h, cctx->hashTable, tableType);
// 
//                 DEBUGLOG(7, "candidate at pos=%u  (offset=%u \n", matchIndex, current - matchIndex);
//                 if ((dictIssue == dictSmall) && (matchIndex < prefixIdxLimit)) { continue; }    /* match outside of valid area */
//                 assert(matchIndex < current);
//                 if ( ((tableType != byU16) || (LZ4_DISTANCE_MAX < LZ4_DISTANCE_ABSOLUTE_MAX))
//                   && (matchIndex+LZ4_DISTANCE_MAX < current)) {
//                     continue;
//                 } /* too far */
//                 assert((current - matchIndex) <= LZ4_DISTANCE_MAX);  /* match now expected within distance */
// 
//                 if (LZ4_read32(match) == LZ4_read32(ip)) {
//                     if (maybe_extMem) offset = current - matchIndex;
//                     break;   /* match found */
//                 }
// 
//             } while(1);
//         }
// 
//         /* Catch up */
//         filledIp = ip;
//         while (((ip>anchor) & (match > lowLimit)) && (unlikely(ip[-1]==match[-1]))) { ip--; match--; }
// 
//         /* Encode Literals */
//         {   unsigned const litLength = (unsigned)(ip - anchor);
//             token = op++;
//             if ((outputDirective == limitedOutput) &&  /* Check output buffer overflow */
//                 (unlikely(op + litLength + (2 + 1 + LASTLITERALS) + (litLength/255) > olimit)) ) {
//                 return 0;   /* cannot compress within `dst` budget. Stored indexes in hash table are nonetheless fine */
//             }
//             if ((outputDirective == fillOutput) &&
//                 (unlikely(op + (litLength+240)/255 /* litlen */ + litLength /* literals */ + 2 /* offset */ + 1 /* token */ + MFLIMIT - MINMATCH /* min last literals so last match is <= end - MFLIMIT */ > olimit))) {
//                 op--;
//                 goto _last_literals;
//             }
//             if (litLength >= RUN_MASK) {
//                 int len = (int)(litLength - RUN_MASK);
//                 *token = (RUN_MASK<<ML_BITS);
//                 for(; len >= 255 ; len-=255) *op++ = 255;
//                 *op++ = (BYTE)len;
//             }
//             else *token = (BYTE)(litLength<<ML_BITS);
// 
//             /* Copy Literals */
//             LZ4_wildCopy8(op, anchor, op+litLength);
//             op+=litLength;
//             DEBUGLOG(6, "seq.start:%i, literals=%u, match.start:%i",
//                         (int)(anchor-(const BYTE*)source), litLength, (int)(ip-(const BYTE*)source));
//         }
// 
// _next_match:
//         /* at this stage, the following variables must be correctly set :
//          * - ip : at start of LZ operation
//          * - match : at start of previous pattern occurence; can be within current prefix, or within extDict
//          * - offset : if maybe_ext_memSegment==1 (constant)
//          * - lowLimit : must be == dictionary to mean "match is within extDict"; must be == source otherwise
//          * - token and *token : position to write 4-bits for match length; higher 4-bits for literal length supposed already written
//          */
// 
//         if ((outputDirective == fillOutput) &&
//             (op + 2 /* offset */ + 1 /* token */ + MFLIMIT - MINMATCH /* min last literals so last match is <= end - MFLIMIT */ > olimit)) {
//             /* the match was too close to the end, rewind and go to last literals */
//             op = token;
//             goto _last_literals;
//         }
// 
//         /* Encode Offset */
//         if (maybe_extMem) {   /* static test */
//             DEBUGLOG(6, "             with offset=%u  (ext if > %i)", offset, (int)(ip - (const BYTE*)source));
//             assert(offset <= LZ4_DISTANCE_MAX && offset > 0);
//             LZ4_writeLE16(op, (U16)offset); op+=2;
//         } else  {
//             DEBUGLOG(6, "             with offset=%u  (same segment)", (U32)(ip - match));
//             assert(ip-match <= LZ4_DISTANCE_MAX);
//             LZ4_writeLE16(op, (U16)(ip - match)); op+=2;
//         }
// 
//         /* Encode MatchLength */
//         {   unsigned matchCode;
// 
//             if ( (dictDirective==usingExtDict || dictDirective==usingDictCtx)
//               && (lowLimit==dictionary) /* match within extDict */ ) {
//                 const BYTE* limit = ip + (dictEnd-match);
//                 assert(dictEnd > match);
//                 if (limit > matchlimit) limit = matchlimit;
//                 matchCode = LZ4_count(ip+MINMATCH, match+MINMATCH, limit);
//                 ip += (size_t)matchCode + MINMATCH;
//                 if (ip==limit) {
//                     unsigned const more = LZ4_count(limit, (const BYTE*)source, matchlimit);
//                     matchCode += more;
//                     ip += more;
//                 }
//                 DEBUGLOG(6, "             with matchLength=%u starting in extDict", matchCode+MINMATCH);
//             } else {
//                 matchCode = LZ4_count(ip+MINMATCH, match+MINMATCH, matchlimit);
//                 ip += (size_t)matchCode + MINMATCH;
//                 DEBUGLOG(6, "             with matchLength=%u", matchCode+MINMATCH);
//             }
// 
//             if ((outputDirective) &&    /* Check output buffer overflow */
//                 (unlikely(op + (1 + LASTLITERALS) + (matchCode+240)/255 > olimit)) ) {
//                 if (outputDirective == fillOutput) {
//                     /* Match description too long : reduce it */
//                     U32 newMatchCode = 15 /* in token */ - 1 /* to avoid needing a zero byte */ + ((U32)(olimit - op) - 1 - LASTLITERALS) * 255;
//                     ip -= matchCode - newMatchCode;
//                     assert(newMatchCode < matchCode);
//                     matchCode = newMatchCode;
//                     if (unlikely(ip <= filledIp)) {
//                         /* We have already filled up to filledIp so if ip ends up less than filledIp
//                          * we have positions in the hash table beyond the current position. This is
//                          * a problem if we reuse the hash table. So we have to remove these positions
//                          * from the hash table.
//                          */
//                         const BYTE* ptr;
//                         DEBUGLOG(5, "Clearing %u positions", (U32)(filledIp - ip));
//                         for (ptr = ip; ptr <= filledIp; ++ptr) {
//                             U32 const h = LZ4_hashPosition(ptr, tableType);
//                             LZ4_clearHash(h, cctx->hashTable, tableType);
//                         }
//                     }
//                 } else {
//                     assert(outputDirective == limitedOutput);
//                     return 0;   /* cannot compress within `dst` budget. Stored indexes in hash table are nonetheless fine */
//                 }
//             }
//             if (matchCode >= ML_MASK) {
//                 *token += ML_MASK;
//                 matchCode -= ML_MASK;
//                 LZ4_write32(op, 0xFFFFFFFF);
//                 while (matchCode >= 4*255) {
//                     op+=4;
//                     LZ4_write32(op, 0xFFFFFFFF);
//                     matchCode -= 4*255;
//                 }
//                 op += matchCode / 255;
//                 *op++ = (BYTE)(matchCode % 255);
//             } else
//                 *token += (BYTE)(matchCode);
//         }
//         /* Ensure we have enough space for the last literals. */
//         assert(!(outputDirective == fillOutput && op + 1 + LASTLITERALS > olimit));
// 
//         anchor = ip;
// 
//         /* Test end of chunk */
//         if (ip >= mflimitPlusOne) break;
// 
//         /* Fill table */
//         LZ4_putPosition(ip-2, cctx->hashTable, tableType, base);
// 
//         /* Test next position */
//         if (tableType == byPtr) {
// 
//             match = LZ4_getPosition(ip, cctx->hashTable, tableType, base);
//             LZ4_putPosition(ip, cctx->hashTable, tableType, base);
//             if ( (match+LZ4_DISTANCE_MAX >= ip)
//               && (LZ4_read32(match) == LZ4_read32(ip)) )
//             { token=op++; *token=0; goto _next_match; }
// 
//         } else {   /* byU32, byU16 */
// 
//             U32 const h = LZ4_hashPosition(ip, tableType);
//             U32 const current = (U32)(ip-base);
//             U32 matchIndex = LZ4_getIndexOnHash(h, cctx->hashTable, tableType);
//             assert(matchIndex < current);
//             if (dictDirective == usingDictCtx) {
//                 if (matchIndex < startIndex) {
//                     /* there was no match, try the dictionary */
//                     matchIndex = LZ4_getIndexOnHash(h, dictCtx->hashTable, byU32);
//                     match = dictBase + matchIndex;
//                     lowLimit = dictionary;   /* required for match length counter */
//                     matchIndex += dictDelta;
//                 } else {
//                     match = base + matchIndex;
//                     lowLimit = (const BYTE*)source;  /* required for match length counter */
//                 }
//             } else if (dictDirective==usingExtDict) {
//                 if (matchIndex < startIndex) {
//                     match = dictBase + matchIndex;
//                     lowLimit = dictionary;   /* required for match length counter */
//                 } else {
//                     match = base + matchIndex;
//                     lowLimit = (const BYTE*)source;   /* required for match length counter */
//                 }
//             } else {   /* single memory segment */
//                 match = base + matchIndex;
//             }
//             LZ4_putIndexOnHash(current, h, cctx->hashTable, tableType);
//             assert(matchIndex < current);
//             if ( ((dictIssue==dictSmall) ? (matchIndex >= prefixIdxLimit) : 1)
//               && (((tableType==byU16) && (LZ4_DISTANCE_MAX == LZ4_DISTANCE_ABSOLUTE_MAX)) ? 1 : (matchIndex+LZ4_DISTANCE_MAX >= current))
//               && (LZ4_read32(match) == LZ4_read32(ip)) ) {
//                 token=op++;
//                 *token=0;
//                 if (maybe_extMem) offset = current - matchIndex;
//                 DEBUGLOG(6, "seq.start:%i, literals=%u, match.start:%i",
//                             (int)(anchor-(const BYTE*)source), 0, (int)(ip-(const BYTE*)source));
//                 goto _next_match;
//             }
//         }
// 
//         /* Prepare next loop */
//         forwardH = LZ4_hashPosition(++ip, tableType);
// 
//     }
// 
// _last_literals:
//     /* Encode Last Literals */
//     {   size_t lastRun = (size_t)(iend - anchor);
//         if ( (outputDirective) &&  /* Check output buffer overflow */
//             (op + lastRun + 1 + ((lastRun+255-RUN_MASK)/255) > olimit)) {
//             if (outputDirective == fillOutput) {
//                 /* adapt lastRun to fill 'dst' */
//                 assert(olimit >= op);
//                 lastRun  = (size_t)(olimit-op) - 1/*token*/;
//                 lastRun -= (lastRun + 256 - RUN_MASK) / 256;  /*additional length tokens*/
//             } else {
//                 assert(outputDirective == limitedOutput);
//                 return 0;   /* cannot compress within `dst` budget. Stored indexes in hash table are nonetheless fine */
//             }
//         }
//         DEBUGLOG(6, "Final literal run : %i literals", (int)lastRun);
//         if (lastRun >= RUN_MASK) {
//             size_t accumulator = lastRun - RUN_MASK;
//             *op++ = RUN_MASK << ML_BITS;
//             for(; accumulator >= 255 ; accumulator-=255) *op++ = 255;
//             *op++ = (BYTE) accumulator;
//         } else {
//             *op++ = (BYTE)(lastRun<<ML_BITS);
//         }
//         LZ4_memcpy(op, anchor, lastRun);
//         ip = anchor + lastRun;
//         op += lastRun;
//     }
// 
//     if (outputDirective == fillOutput) {
//         *inputConsumed = (int) (((const char*)ip)-source);
//     }
//     result = (int)(((char*)op) - dest);
//     assert(result > 0);
//     DEBUGLOG(5, "LZ4_compress_generic: compressed %i bytes into %i bytes", inputSize, result);
//     return result;
// }
// 
// /** LZ4_compress_generic() :
//  *  inlined, to ensure branches are decided at compilation time;
//  *  takes care of src == (NULL, 0)
//  *  and forward the rest to LZ4_compress_generic_validated */
// LZ4_FORCE_INLINE int LZ4_compress_generic(
//                  LZ4_stream_t_internal* const cctx,
//                  const char* const src,
//                  char* const dst,
//                  const int srcSize,
//                  int *inputConsumed, /* only written when outputDirective == fillOutput */
//                  const int dstCapacity,
//                  const limitedOutput_directive outputDirective,
//                  const tableType_t tableType,
//                  const dict_directive dictDirective,
//                  const dictIssue_directive dictIssue,
//                  const int acceleration)
// {
//     DEBUGLOG(5, "LZ4_compress_generic: srcSize=%i, dstCapacity=%i",
//                 srcSize, dstCapacity);
// 
//     if ((U32)srcSize > (U32)LZ4_MAX_INPUT_SIZE) { return 0; }  /* Unsupported srcSize, too large (or negative) */
//     if (srcSize == 0) {   /* src == NULL supported if srcSize == 0 */
//         if (outputDirective != notLimited && dstCapacity <= 0) return 0;  /* no output, can't write anything */
//         DEBUGLOG(5, "Generating an empty block");
//         assert(outputDirective == notLimited || dstCapacity >= 1);
//         assert(dst != NULL);
//         dst[0] = 0;
//         if (outputDirective == fillOutput) {
//             assert (inputConsumed != NULL);
//             *inputConsumed = 0;
//         }
//         return 1;
//     }
//     assert(src != NULL);
// 
//     return LZ4_compress_generic_validated(cctx, src, dst, srcSize,
//                 inputConsumed, /* only written into if outputDirective == fillOutput */
//                 dstCapacity, outputDirective,
//                 tableType, dictDirective, dictIssue, acceleration);
// }
// 
// 
// int LZ4_compress_fast_extState(void* state, const char* source, char* dest, int inputSize, int maxOutputSize, int acceleration)
// {
//     LZ4_stream_t_internal* const ctx = & LZ4_initStream(state, sizeof(LZ4_stream_t)) -> internal_donotuse;
//     assert(ctx != NULL);
//     if (acceleration < 1) acceleration = LZ4_ACCELERATION_DEFAULT;
//     if (acceleration > LZ4_ACCELERATION_MAX) acceleration = LZ4_ACCELERATION_MAX;
//     if (maxOutputSize >= LZ4_compressBound(inputSize)) {
//         if (inputSize < LZ4_64Klimit) {
//             return LZ4_compress_generic(ctx, source, dest, inputSize, NULL, 0, notLimited, byU16, noDict, noDictIssue, acceleration);
//         } else {
//             const tableType_t tableType = ((sizeof(void*)==4) && ((uptrval)source > LZ4_DISTANCE_MAX)) ? byPtr : byU32;
//             return LZ4_compress_generic(ctx, source, dest, inputSize, NULL, 0, notLimited, tableType, noDict, noDictIssue, acceleration);
//         }
//     } else {
//         if (inputSize < LZ4_64Klimit) {
//             return LZ4_compress_generic(ctx, source, dest, inputSize, NULL, maxOutputSize, limitedOutput, byU16, noDict, noDictIssue, acceleration);
//         } else {
//             const tableType_t tableType = ((sizeof(void*)==4) && ((uptrval)source > LZ4_DISTANCE_MAX)) ? byPtr : byU32;
//             return LZ4_compress_generic(ctx, source, dest, inputSize, NULL, maxOutputSize, limitedOutput, tableType, noDict, noDictIssue, acceleration);
//         }
//     }
// }
// 
// /**
//  * LZ4_compress_fast_extState_fastReset() :
//  * A variant of LZ4_compress_fast_extState().
//  *
//  * Using this variant avoids an expensive initialization step. It is only safe
//  * to call if the state buffer is known to be correctly initialized already
//  * (see comment in lz4.h on LZ4_resetStream_fast() for a definition of
//  * "correctly initialized").
//  */
// int LZ4_compress_fast_extState_fastReset(void* state, const char* src, char* dst, int srcSize, int dstCapacity, int acceleration)
// {
//     LZ4_stream_t_internal* ctx = &((LZ4_stream_t*)state)->internal_donotuse;
//     if (acceleration < 1) acceleration = LZ4_ACCELERATION_DEFAULT;
//     if (acceleration > LZ4_ACCELERATION_MAX) acceleration = LZ4_ACCELERATION_MAX;
// 
//     if (dstCapacity >= LZ4_compressBound(srcSize)) {
//         if (srcSize < LZ4_64Klimit) {
//             const tableType_t tableType = byU16;
//             LZ4_prepareTable(ctx, srcSize, tableType);
//             if (ctx->currentOffset) {
//                 return LZ4_compress_generic(ctx, src, dst, srcSize, NULL, 0, notLimited, tableType, noDict, dictSmall, acceleration);
//             } else {
//                 return LZ4_compress_generic(ctx, src, dst, srcSize, NULL, 0, notLimited, tableType, noDict, noDictIssue, acceleration);
//             }
//         } else {
//             const tableType_t tableType = ((sizeof(void*)==4) && ((uptrval)src > LZ4_DISTANCE_MAX)) ? byPtr : byU32;
//             LZ4_prepareTable(ctx, srcSize, tableType);
//             return LZ4_compress_generic(ctx, src, dst, srcSize, NULL, 0, notLimited, tableType, noDict, noDictIssue, acceleration);
//         }
//     } else {
//         if (srcSize < LZ4_64Klimit) {
//             const tableType_t tableType = byU16;
//             LZ4_prepareTable(ctx, srcSize, tableType);
//             if (ctx->currentOffset) {
//                 return LZ4_compress_generic(ctx, src, dst, srcSize, NULL, dstCapacity, limitedOutput, tableType, noDict, dictSmall, acceleration);
//             } else {
//                 return LZ4_compress_generic(ctx, src, dst, srcSize, NULL, dstCapacity, limitedOutput, tableType, noDict, noDictIssue, acceleration);
//             }
//         } else {
//             const tableType_t tableType = ((sizeof(void*)==4) && ((uptrval)src > LZ4_DISTANCE_MAX)) ? byPtr : byU32;
//             LZ4_prepareTable(ctx, srcSize, tableType);
//             return LZ4_compress_generic(ctx, src, dst, srcSize, NULL, dstCapacity, limitedOutput, tableType, noDict, noDictIssue, acceleration);
//         }
//     }
// }
// 
// 
// int LZ4_compress_fast(const char* source, char* dest, int inputSize, int maxOutputSize, int acceleration)
// {
//     int result;
// #if (LZ4_HEAPMODE)
//     LZ4_stream_t* ctxPtr = ALLOC(sizeof(LZ4_stream_t));   /* malloc-calloc always properly aligned */
//     if (ctxPtr == NULL) return 0;
// #else
//     LZ4_stream_t ctx;
//     LZ4_stream_t* const ctxPtr = &ctx;
// #endif
//     result = LZ4_compress_fast_extState(ctxPtr, source, dest, inputSize, maxOutputSize, acceleration);
// 
// #if (LZ4_HEAPMODE)
//     FREEMEM(ctxPtr);
// #endif
//     return result;
// }
// 
// 
// int LZ4_compress_default(const char* src, char* dst, int srcSize, int maxOutputSize)
// {
//     return LZ4_compress_fast(src, dst, srcSize, maxOutputSize, 1);
// }
// 
// 
// /* Note!: This function leaves the stream in an unclean/broken state!
//  * It is not safe to subsequently use the same state with a _fastReset() or
//  * _continue() call without resetting it. */
// static int LZ4_compress_destSize_extState (LZ4_stream_t* state, const char* src, char* dst, int* srcSizePtr, int targetDstSize)
// {
//     void* const s = LZ4_initStream(state, sizeof (*state));
//     assert(s != NULL); (void)s;
// 
//     if (targetDstSize >= LZ4_compressBound(*srcSizePtr)) {  /* compression success is guaranteed */
//         return LZ4_compress_fast_extState(state, src, dst, *srcSizePtr, targetDstSize, 1);
//     } else {
//         if (*srcSizePtr < LZ4_64Klimit) {
//             return LZ4_compress_generic(&state->internal_donotuse, src, dst, *srcSizePtr, srcSizePtr, targetDstSize, fillOutput, byU16, noDict, noDictIssue, 1);
//         } else {
//             tableType_t const addrMode = ((sizeof(void*)==4) && ((uptrval)src > LZ4_DISTANCE_MAX)) ? byPtr : byU32;
//             return LZ4_compress_generic(&state->internal_donotuse, src, dst, *srcSizePtr, srcSizePtr, targetDstSize, fillOutput, addrMode, noDict, noDictIssue, 1);
//     }   }
// }
// 
// 
// int LZ4_compress_destSize(const char* src, char* dst, int* srcSizePtr, int targetDstSize)
// {
// #if (LZ4_HEAPMODE)
//     LZ4_stream_t* ctx = (LZ4_stream_t*)ALLOC(sizeof(LZ4_stream_t));   /* malloc-calloc always properly aligned */
//     if (ctx == NULL) return 0;
// #else
//     LZ4_stream_t ctxBody;
//     LZ4_stream_t* ctx = &ctxBody;
// #endif
// 
//     int result = LZ4_compress_destSize_extState(ctx, src, dst, srcSizePtr, targetDstSize);
// 
// #if (LZ4_HEAPMODE)
//     FREEMEM(ctx);
// #endif
//     return result;
// }
// 
// 
// 
// /*-******************************
// *  Streaming functions
// ********************************/
// 
// LZ4_stream_t* LZ4_createStream(void)
// {
//     LZ4_stream_t* const lz4s = (LZ4_stream_t*)ALLOC(sizeof(LZ4_stream_t));
//     LZ4_STATIC_ASSERT(LZ4_STREAMSIZE >= sizeof(LZ4_stream_t_internal));    /* A compilation error here means LZ4_STREAMSIZE is not large enough */
//     DEBUGLOG(4, "LZ4_createStream %p", lz4s);
//     if (lz4s == NULL) return NULL;
//     LZ4_initStream(lz4s, sizeof(*lz4s));
//     return lz4s;
// }
// 
// static size_t LZ4_stream_t_alignment(void)
// {
// #if LZ4_ALIGN_TEST
//     typedef struct { char c; LZ4_stream_t t; } t_a;
//     return sizeof(t_a) - sizeof(LZ4_stream_t);
// #else
//     return 1;  /* effectively disabled */
// #endif
// }
// 
// LZ4_stream_t* LZ4_initStream (void* buffer, size_t size)
// {
//     DEBUGLOG(5, "LZ4_initStream");
//     if (buffer == NULL) { return NULL; }
//     if (size < sizeof(LZ4_stream_t)) { return NULL; }
//     if (!LZ4_isAligned(buffer, LZ4_stream_t_alignment())) return NULL;
//     MEM_INIT(buffer, 0, sizeof(LZ4_stream_t_internal));
//     return (LZ4_stream_t*)buffer;
// }
// 
// /* resetStream is now deprecated,
//  * prefer initStream() which is more general */
// void LZ4_resetStream (LZ4_stream_t* LZ4_stream)
// {
//     DEBUGLOG(5, "LZ4_resetStream (ctx:%p)", LZ4_stream);
//     MEM_INIT(LZ4_stream, 0, sizeof(LZ4_stream_t_internal));
// }
// 
// void LZ4_resetStream_fast(LZ4_stream_t* ctx) {
//     LZ4_prepareTable(&(ctx->internal_donotuse), 0, byU32);
// }
// 
// int LZ4_freeStream (LZ4_stream_t* LZ4_stream)
// {
//     if (!LZ4_stream) return 0;   /* support free on NULL */
//     DEBUGLOG(5, "LZ4_freeStream %p", LZ4_stream);
//     FREEMEM(LZ4_stream);
//     return (0);
// }
// 
// 
// #define HASH_UNIT sizeof(reg_t)
// int LZ4_loadDict (LZ4_stream_t* LZ4_dict, const char* dictionary, int dictSize)
// {
//     LZ4_stream_t_internal* dict = &LZ4_dict->internal_donotuse;
//     const tableType_t tableType = byU32;
//     const BYTE* p = (const BYTE*)dictionary;
//     const BYTE* const dictEnd = p + dictSize;
//     const BYTE* base;
// 
//     DEBUGLOG(4, "LZ4_loadDict (%i bytes from %p into %p)", dictSize, dictionary, LZ4_dict);
// 
//     /* It's necessary to reset the context,
//      * and not just continue it with prepareTable()
//      * to avoid any risk of generating overflowing matchIndex
//      * when compressing using this dictionary */
//     LZ4_resetStream(LZ4_dict);
// 
//     /* We always increment the offset by 64 KB, since, if the dict is longer,
//      * we truncate it to the last 64k, and if it's shorter, we still want to
//      * advance by a whole window length so we can provide the guarantee that
//      * there are only valid offsets in the window, which allows an optimization
//      * in LZ4_compress_fast_continue() where it uses noDictIssue even when the
//      * dictionary isn't a full 64k. */
//     dict->currentOffset += 64 KB;
// 
//     if (dictSize < (int)HASH_UNIT) {
//         return 0;
//     }
// 
//     if ((dictEnd - p) > 64 KB) p = dictEnd - 64 KB;
//     base = dictEnd - dict->currentOffset;
//     dict->dictionary = p;
//     dict->dictSize = (U32)(dictEnd - p);
//     dict->tableType = (U32)tableType;
// 
//     while (p <= dictEnd-HASH_UNIT) {
//         LZ4_putPosition(p, dict->hashTable, tableType, base);
//         p+=3;
//     }
// 
//     return (int)dict->dictSize;
// }
// 
// void LZ4_attach_dictionary(LZ4_stream_t* workingStream, const LZ4_stream_t* dictionaryStream) {
//     const LZ4_stream_t_internal* dictCtx = dictionaryStream == NULL ? NULL :
//         &(dictionaryStream->internal_donotuse);
// 
//     DEBUGLOG(4, "LZ4_attach_dictionary (%p, %p, size %u)",
//              workingStream, dictionaryStream,
//              dictCtx != NULL ? dictCtx->dictSize : 0);
// 
//     if (dictCtx != NULL) {
//         /* If the current offset is zero, we will never look in the
//          * external dictionary context, since there is no value a table
//          * entry can take that indicate a miss. In that case, we need
//          * to bump the offset to something non-zero.
//          */
//         if (workingStream->internal_donotuse.currentOffset == 0) {
//             workingStream->internal_donotuse.currentOffset = 64 KB;
//         }
// 
//         /* Don't actually attach an empty dictionary.
//          */
//         if (dictCtx->dictSize == 0) {
//             dictCtx = NULL;
//         }
//     }
//     workingStream->internal_donotuse.dictCtx = dictCtx;
// }
// 
// 
// static void LZ4_renormDictT(LZ4_stream_t_internal* LZ4_dict, int nextSize)
// {
//     assert(nextSize >= 0);
//     if (LZ4_dict->currentOffset + (unsigned)nextSize > 0x80000000) {   /* potential ptrdiff_t overflow (32-bits mode) */
//         /* rescale hash table */
//         U32 const delta = LZ4_dict->currentOffset - 64 KB;
//         const BYTE* dictEnd = LZ4_dict->dictionary + LZ4_dict->dictSize;
//         int i;
//         DEBUGLOG(4, "LZ4_renormDictT");
//         for (i=0; i<LZ4_HASH_SIZE_U32; i++) {
//             if (LZ4_dict->hashTable[i] < delta) LZ4_dict->hashTable[i]=0;
//             else LZ4_dict->hashTable[i] -= delta;
//         }
//         LZ4_dict->currentOffset = 64 KB;
//         if (LZ4_dict->dictSize > 64 KB) LZ4_dict->dictSize = 64 KB;
//         LZ4_dict->dictionary = dictEnd - LZ4_dict->dictSize;
//     }
// }
// 
// 
// int LZ4_compress_fast_continue (LZ4_stream_t* LZ4_stream,
//                                 const char* source, char* dest,
//                                 int inputSize, int maxOutputSize,
//                                 int acceleration)
// {
//     const tableType_t tableType = byU32;
//     LZ4_stream_t_internal* streamPtr = &LZ4_stream->internal_donotuse;
//     const BYTE* dictEnd = streamPtr->dictionary + streamPtr->dictSize;
// 
//     DEBUGLOG(5, "LZ4_compress_fast_continue (inputSize=%i)", inputSize);
// 
//     LZ4_renormDictT(streamPtr, inputSize);   /* avoid index overflow */
//     if (acceleration < 1) acceleration = LZ4_ACCELERATION_DEFAULT;
//     if (acceleration > LZ4_ACCELERATION_MAX) acceleration = LZ4_ACCELERATION_MAX;
// 
//     /* invalidate tiny dictionaries */
//     if ( (streamPtr->dictSize-1 < 4-1)   /* intentional underflow */
//       && (dictEnd != (const BYTE*)source) ) {
//         DEBUGLOG(5, "LZ4_compress_fast_continue: dictSize(%u) at addr:%p is too small", streamPtr->dictSize, streamPtr->dictionary);
//         streamPtr->dictSize = 0;
//         streamPtr->dictionary = (const BYTE*)source;
//         dictEnd = (const BYTE*)source;
//     }
// 
//     /* Check overlapping input/dictionary space */
//     {   const BYTE* sourceEnd = (const BYTE*) source + inputSize;
//         if ((sourceEnd > streamPtr->dictionary) && (sourceEnd < dictEnd)) {
//             streamPtr->dictSize = (U32)(dictEnd - sourceEnd);
//             if (streamPtr->dictSize > 64 KB) streamPtr->dictSize = 64 KB;
//             if (streamPtr->dictSize < 4) streamPtr->dictSize = 0;
//             streamPtr->dictionary = dictEnd - streamPtr->dictSize;
//         }
//     }
// 
//     /* prefix mode : source data follows dictionary */
//     if (dictEnd == (const BYTE*)source) {
//         if ((streamPtr->dictSize < 64 KB) && (streamPtr->dictSize < streamPtr->currentOffset))
//             return LZ4_compress_generic(streamPtr, source, dest, inputSize, NULL, maxOutputSize, limitedOutput, tableType, withPrefix64k, dictSmall, acceleration);
//         else
//             return LZ4_compress_generic(streamPtr, source, dest, inputSize, NULL, maxOutputSize, limitedOutput, tableType, withPrefix64k, noDictIssue, acceleration);
//     }
// 
//     /* external dictionary mode */
//     {   int result;
//         if (streamPtr->dictCtx) {
//             /* We depend here on the fact that dictCtx'es (produced by
//              * LZ4_loadDict) guarantee that their tables contain no references
//              * to offsets between dictCtx->currentOffset - 64 KB and
//              * dictCtx->currentOffset - dictCtx->dictSize. This makes it safe
//              * to use noDictIssue even when the dict isn't a full 64 KB.
//              */
//             if (inputSize > 4 KB) {
//                 /* For compressing large blobs, it is faster to pay the setup
//                  * cost to copy the dictionary's tables into the active context,
//                  * so that the compression loop is only looking into one table.
//                  */
//                 LZ4_memcpy(streamPtr, streamPtr->dictCtx, sizeof(*streamPtr));
//                 result = LZ4_compress_generic(streamPtr, source, dest, inputSize, NULL, maxOutputSize, limitedOutput, tableType, usingExtDict, noDictIssue, acceleration);
//             } else {
//                 result = LZ4_compress_generic(streamPtr, source, dest, inputSize, NULL, maxOutputSize, limitedOutput, tableType, usingDictCtx, noDictIssue, acceleration);
//             }
//         } else {
//             if ((streamPtr->dictSize < 64 KB) && (streamPtr->dictSize < streamPtr->currentOffset)) {
//                 result = LZ4_compress_generic(streamPtr, source, dest, inputSize, NULL, maxOutputSize, limitedOutput, tableType, usingExtDict, dictSmall, acceleration);
//             } else {
//                 result = LZ4_compress_generic(streamPtr, source, dest, inputSize, NULL, maxOutputSize, limitedOutput, tableType, usingExtDict, noDictIssue, acceleration);
//             }
//         }
//         streamPtr->dictionary = (const BYTE*)source;
//         streamPtr->dictSize = (U32)inputSize;
//         return result;
//     }
// }
// 
// 
// /* Hidden debug function, to force-test external dictionary mode */
// int LZ4_compress_forceExtDict (LZ4_stream_t* LZ4_dict, const char* source, char* dest, int srcSize)
// {
//     LZ4_stream_t_internal* streamPtr = &LZ4_dict->internal_donotuse;
//     int result;
// 
//     LZ4_renormDictT(streamPtr, srcSize);
// 
//     if ((streamPtr->dictSize < 64 KB) && (streamPtr->dictSize < streamPtr->currentOffset)) {
//         result = LZ4_compress_generic(streamPtr, source, dest, srcSize, NULL, 0, notLimited, byU32, usingExtDict, dictSmall, 1);
//     } else {
//         result = LZ4_compress_generic(streamPtr, source, dest, srcSize, NULL, 0, notLimited, byU32, usingExtDict, noDictIssue, 1);
//     }
// 
//     streamPtr->dictionary = (const BYTE*)source;
//     streamPtr->dictSize = (U32)srcSize;
// 
//     return result;
// }
// 
// 
// /*! LZ4_saveDict() :
//  *  If previously compressed data block is not guaranteed to remain available at its memory location,
//  *  save it into a safer place (char* safeBuffer).
//  *  Note : you don't need to call LZ4_loadDict() afterwards,
//  *         dictionary is immediately usable, you can therefore call LZ4_compress_fast_continue().
//  *  Return : saved dictionary size in bytes (necessarily <= dictSize), or 0 if error.
//  */
// int LZ4_saveDict (LZ4_stream_t* LZ4_dict, char* safeBuffer, int dictSize)
// {
//     LZ4_stream_t_internal* const dict = &LZ4_dict->internal_donotuse;
//     const BYTE* const previousDictEnd = dict->dictionary + dict->dictSize;
// 
//     if ((U32)dictSize > 64 KB) { dictSize = 64 KB; } /* useless to define a dictionary > 64 KB */
//     if ((U32)dictSize > dict->dictSize) { dictSize = (int)dict->dictSize; }
// 
//     if (safeBuffer == NULL) assert(dictSize == 0);
//     if (dictSize > 0)
//         memmove(safeBuffer, previousDictEnd - dictSize, dictSize);
// 
//     dict->dictionary = (const BYTE*)safeBuffer;
//     dict->dictSize = (U32)dictSize;
// 
//     return dictSize;
// }
// 
// 
// 
// /*-*******************************
//  *  Decompression functions
//  ********************************/
// 
// typedef enum { endOnOutputSize = 0, endOnInputSize = 1 } endCondition_directive;
// typedef enum { decode_full_block = 0, partial_decode = 1 } earlyEnd_directive;
// 
// #undef MIN
// #define MIN(a,b)    ( (a) < (b) ? (a) : (b) )
// 
// /* Read the variable-length literal or match length.
//  *
//  * ip - pointer to use as input.
//  * lencheck - end ip.  Return an error if ip advances >= lencheck.
//  * loop_check - check ip >= lencheck in body of loop.  Returns loop_error if so.
//  * initial_check - check ip >= lencheck before start of loop.  Returns initial_error if so.
//  * error (output) - error code.  Should be set to 0 before call.
//  */
// typedef enum { loop_error = -2, initial_error = -1, ok = 0 } variable_length_error;
// LZ4_FORCE_INLINE unsigned
// read_variable_length(const BYTE**ip, const BYTE* lencheck,
//                      int loop_check, int initial_check,
//                      variable_length_error* error)
// {
//     U32 length = 0;
//     U32 s;
//     if (initial_check && unlikely((*ip) >= lencheck)) {    /* overflow detection */
//         *error = initial_error;
//         return length;
//     }
//     do {
//         s = **ip;
//         (*ip)++;
//         length += s;
//         if (loop_check && unlikely((*ip) >= lencheck)) {    /* overflow detection */
//             *error = loop_error;
//             return length;
//         }
//     } while (s==255);
// 
//     return length;
// }
// 
// /*! LZ4_decompress_generic() :
//  *  This generic decompression function covers all use cases.
//  *  It shall be instantiated several times, using different sets of directives.
//  *  Note that it is important for performance that this function really get inlined,
//  *  in order to remove useless branches during compilation optimization.
//  */
// LZ4_FORCE_INLINE int
// LZ4_decompress_generic(
//                  const char* const src,
//                  char* const dst,
//                  int srcSize,
//                  int outputSize,         /* If endOnInput==endOnInputSize, this value is `dstCapacity` */
// 
//                  endCondition_directive endOnInput,   /* endOnOutputSize, endOnInputSize */
//                  earlyEnd_directive partialDecoding,  /* full, partial */
//                  dict_directive dict,                 /* noDict, withPrefix64k, usingExtDict */
//                  const BYTE* const lowPrefix,  /* always <= dst, == dst when no prefix */
//                  const BYTE* const dictStart,  /* only if dict==usingExtDict */
//                  const size_t dictSize         /* note : = 0 if noDict */
//                  )
// {
//     if (src == NULL) { return -1; }
// 
//     {   const BYTE* ip = (const BYTE*) src;
//         const BYTE* const iend = ip + srcSize;
// 
//         BYTE* op = (BYTE*) dst;
//         BYTE* const oend = op + outputSize;
//         BYTE* cpy;
// 
//         const BYTE* const dictEnd = (dictStart == NULL) ? NULL : dictStart + dictSize;
// 
//         const int safeDecode = (endOnInput==endOnInputSize);
//         const int checkOffset = ((safeDecode) && (dictSize < (int)(64 KB)));
// 
// 
//         /* Set up the "end" pointers for the shortcut. */
//         const BYTE* const shortiend = iend - (endOnInput ? 14 : 8) /*maxLL*/ - 2 /*offset*/;
//         const BYTE* const shortoend = oend - (endOnInput ? 14 : 8) /*maxLL*/ - 18 /*maxML*/;
// 
//         const BYTE* match;
//         size_t offset;
//         unsigned token;
//         size_t length;
// 
// 
//         DEBUGLOG(5, "LZ4_decompress_generic (srcSize:%i, dstSize:%i)", srcSize, outputSize);
// 
//         /* Special cases */
//         assert(lowPrefix <= op);
//         if ((endOnInput) && (unlikely(outputSize==0))) {
//             /* Empty output buffer */
//             if (partialDecoding) return 0;
//             return ((srcSize==1) && (*ip==0)) ? 0 : -1;
//         }
//         if ((!endOnInput) && (unlikely(outputSize==0))) { return (*ip==0 ? 1 : -1); }
//         if ((endOnInput) && unlikely(srcSize==0)) { return -1; }
// 
// 	/* Currently the fast loop shows a regression on qualcomm arm chips. */
// #if LZ4_FAST_DEC_LOOP
//         if ((oend - op) < FASTLOOP_SAFE_DISTANCE) {
//             DEBUGLOG(6, "skip fast decode loop");
//             goto safe_decode;
//         }
// 
//         /* Fast loop : decode sequences as long as output < iend-FASTLOOP_SAFE_DISTANCE */
//         while (1) {
//             /* Main fastloop assertion: We can always wildcopy FASTLOOP_SAFE_DISTANCE */
//             assert(oend - op >= FASTLOOP_SAFE_DISTANCE);
//             if (endOnInput) { assert(ip < iend); }
//             token = *ip++;
//             length = token >> ML_BITS;  /* literal length */
// 
//             assert(!endOnInput || ip <= iend); /* ip < iend before the increment */
// 
//             /* decode literal length */
//             if (length == RUN_MASK) {
//                 variable_length_error error = ok;
//                 length += read_variable_length(&ip, iend-RUN_MASK, (int)endOnInput, (int)endOnInput, &error);
//                 if (error == initial_error) { goto _output_error; }
//                 if ((safeDecode) && unlikely((uptrval)(op)+length<(uptrval)(op))) { goto _output_error; } /* overflow detection */
//                 if ((safeDecode) && unlikely((uptrval)(ip)+length<(uptrval)(ip))) { goto _output_error; } /* overflow detection */
// 
//                 /* copy literals */
//                 cpy = op+length;
//                 LZ4_STATIC_ASSERT(MFLIMIT >= WILDCOPYLENGTH);
//                 if (endOnInput) {  /* LZ4_decompress_safe() */
//                     if ((cpy>oend-32) || (ip+length>iend-32)) { goto safe_literal_copy; }
//                     LZ4_wildCopy32(op, ip, cpy);
//                 } else {   /* LZ4_decompress_fast() */
//                     if (cpy>oend-8) { goto safe_literal_copy; }
//                     LZ4_wildCopy8(op, ip, cpy); /* LZ4_decompress_fast() cannot copy more than 8 bytes at a time :
//                                                  * it doesn't know input length, and only relies on end-of-block properties */
//                 }
//                 ip += length; op = cpy;
//             } else {
//                 cpy = op+length;
//                 if (endOnInput) {  /* LZ4_decompress_safe() */
//                     DEBUGLOG(7, "copy %u bytes in a 16-bytes stripe", (unsigned)length);
//                     /* We don't need to check oend, since we check it once for each loop below */
//                     if (ip > iend-(16 + 1/*max lit + offset + nextToken*/)) { goto safe_literal_copy; }
//                     /* Literals can only be 14, but hope compilers optimize if we copy by a register size */
//                     LZ4_memcpy(op, ip, 16);
//                 } else {  /* LZ4_decompress_fast() */
//                     /* LZ4_decompress_fast() cannot copy more than 8 bytes at a time :
//                      * it doesn't know input length, and relies on end-of-block properties */
//                     LZ4_memcpy(op, ip, 8);
//                     if (length > 8) { LZ4_memcpy(op+8, ip+8, 8); }
//                 }
//                 ip += length; op = cpy;
//             }
// 
//             /* get offset */
//             offset = LZ4_readLE16(ip); ip+=2;
//             match = op - offset;
//             assert(match <= op);
// 
//             /* get matchlength */
//             length = token & ML_MASK;
// 
//             if (length == ML_MASK) {
//                 variable_length_error error = ok;
//                 if ((checkOffset) && (unlikely(match + dictSize < lowPrefix))) { goto _output_error; } /* Error : offset outside buffers */
//                 length += read_variable_length(&ip, iend - LASTLITERALS + 1, (int)endOnInput, 0, &error);
//                 if (error != ok) { goto _output_error; }
//                 if ((safeDecode) && unlikely((uptrval)(op)+length<(uptrval)op)) { goto _output_error; } /* overflow detection */
//                 length += MINMATCH;
//                 if (op + length >= oend - FASTLOOP_SAFE_DISTANCE) {
//                     goto safe_match_copy;
//                 }
//             } else {
//                 length += MINMATCH;
//                 if (op + length >= oend - FASTLOOP_SAFE_DISTANCE) {
//                     goto safe_match_copy;
//                 }
// 
//                 /* Fastpath check: Avoids a branch in LZ4_wildCopy32 if true */
//                 if ((dict == withPrefix64k) || (match >= lowPrefix)) {
//                     if (offset >= 8) {
//                         assert(match >= lowPrefix);
//                         assert(match <= op);
//                         assert(op + 18 <= oend);
// 
//                         LZ4_memcpy(op, match, 8);
//                         LZ4_memcpy(op+8, match+8, 8);
//                         LZ4_memcpy(op+16, match+16, 2);
//                         op += length;
//                         continue;
//             }   }   }
// 
//             if (checkOffset && (unlikely(match + dictSize < lowPrefix))) { goto _output_error; } /* Error : offset outside buffers */
//             /* match starting within external dictionary */
//             if ((dict==usingExtDict) && (match < lowPrefix)) {
//                 if (unlikely(op+length > oend-LASTLITERALS)) {
//                     if (partialDecoding) {
//                         DEBUGLOG(7, "partialDecoding: dictionary match, close to dstEnd");
//                         length = MIN(length, (size_t)(oend-op));
//                     } else {
//                         goto _output_error;  /* end-of-block condition violated */
//                 }   }
// 
//                 if (length <= (size_t)(lowPrefix-match)) {
//                     /* match fits entirely within external dictionary : just copy */
//                     memmove(op, dictEnd - (lowPrefix-match), length);
//                     op += length;
//                 } else {
//                     /* match stretches into both external dictionary and current block */
//                     size_t const copySize = (size_t)(lowPrefix - match);
//                     size_t const restSize = length - copySize;
//                     LZ4_memcpy(op, dictEnd - copySize, copySize);
//                     op += copySize;
//                     if (restSize > (size_t)(op - lowPrefix)) {  /* overlap copy */
//                         BYTE* const endOfMatch = op + restSize;
//                         const BYTE* copyFrom = lowPrefix;
//                         while (op < endOfMatch) { *op++ = *copyFrom++; }
//                     } else {
//                         LZ4_memcpy(op, lowPrefix, restSize);
//                         op += restSize;
//                 }   }
//                 continue;
//             }
// 
//             /* copy match within block */
//             cpy = op + length;
// 
//             assert((op <= oend) && (oend-op >= 32));
//             if (unlikely(offset<16)) {
//                 LZ4_memcpy_using_offset(op, match, cpy, offset);
//             } else {
//                 LZ4_wildCopy32(op, match, cpy);
//             }
// 
//             op = cpy;   /* wildcopy correction */
//         }
//     safe_decode:
// #endif
// 
//         /* Main Loop : decode remaining sequences where output < FASTLOOP_SAFE_DISTANCE */
//         while (1) {
//             token = *ip++;
//             length = token >> ML_BITS;  /* literal length */
// 
//             assert(!endOnInput || ip <= iend); /* ip < iend before the increment */
// 
//             /* A two-stage shortcut for the most common case:
//              * 1) If the literal length is 0..14, and there is enough space,
//              * enter the shortcut and copy 16 bytes on behalf of the literals
//              * (in the fast mode, only 8 bytes can be safely copied this way).
//              * 2) Further if the match length is 4..18, copy 18 bytes in a similar
//              * manner; but we ensure that there's enough space in the output for
//              * those 18 bytes earlier, upon entering the shortcut (in other words,
//              * there is a combined check for both stages).
//              */
//             if ( (endOnInput ? length != RUN_MASK : length <= 8)
//                 /* strictly "less than" on input, to re-enter the loop with at least one byte */
//               && likely((endOnInput ? ip < shortiend : 1) & (op <= shortoend)) ) {
//                 /* Copy the literals */
//                 LZ4_memcpy(op, ip, endOnInput ? 16 : 8);
//                 op += length; ip += length;
// 
//                 /* The second stage: prepare for match copying, decode full info.
//                  * If it doesn't work out, the info won't be wasted. */
//                 length = token & ML_MASK; /* match length */
//                 offset = LZ4_readLE16(ip); ip += 2;
//                 match = op - offset;
//                 assert(match <= op); /* check overflow */
// 
//                 /* Do not deal with overlapping matches. */
//                 if ( (length != ML_MASK)
//                   && (offset >= 8)
//                   && (dict==withPrefix64k || match >= lowPrefix) ) {
//                     /* Copy the match. */
//                     LZ4_memcpy(op + 0, match + 0, 8);
//                     LZ4_memcpy(op + 8, match + 8, 8);
//                     LZ4_memcpy(op +16, match +16, 2);
//                     op += length + MINMATCH;
//                     /* Both stages worked, load the next token. */
//                     continue;
//                 }
// 
//                 /* The second stage didn't work out, but the info is ready.
//                  * Propel it right to the point of match copying. */
//                 goto _copy_match;
//             }
// 
//             /* decode literal length */
//             if (length == RUN_MASK) {
//                 variable_length_error error = ok;
//                 length += read_variable_length(&ip, iend-RUN_MASK, (int)endOnInput, (int)endOnInput, &error);
//                 if (error == initial_error) { goto _output_error; }
//                 if ((safeDecode) && unlikely((uptrval)(op)+length<(uptrval)(op))) { goto _output_error; } /* overflow detection */
//                 if ((safeDecode) && unlikely((uptrval)(ip)+length<(uptrval)(ip))) { goto _output_error; } /* overflow detection */
//             }
// 
//             /* copy literals */
//             cpy = op+length;
// #if LZ4_FAST_DEC_LOOP
//         safe_literal_copy:
// #endif
//             LZ4_STATIC_ASSERT(MFLIMIT >= WILDCOPYLENGTH);
//             if ( ((endOnInput) && ((cpy>oend-MFLIMIT) || (ip+length>iend-(2+1+LASTLITERALS))) )
//               || ((!endOnInput) && (cpy>oend-WILDCOPYLENGTH)) )
//             {
//                 /* We've either hit the input parsing restriction or the output parsing restriction.
//                  * In the normal scenario, decoding a full block, it must be the last sequence,
//                  * otherwise it's an error (invalid input or dimensions).
//                  * In partialDecoding scenario, it's necessary to ensure there is no buffer overflow.
//                  */
//                 if (partialDecoding) {
//                     /* Since we are partial decoding we may be in this block because of the output parsing
//                      * restriction, which is not valid since the output buffer is allowed to be undersized.
//                      */
//                     assert(endOnInput);
//                     DEBUGLOG(7, "partialDecoding: copying literals, close to input or output end")
//                     DEBUGLOG(7, "partialDecoding: literal length = %u", (unsigned)length);
//                     DEBUGLOG(7, "partialDecoding: remaining space in dstBuffer : %i", (int)(oend - op));
//                     DEBUGLOG(7, "partialDecoding: remaining space in srcBuffer : %i", (int)(iend - ip));
//                     /* Finishing in the middle of a literals segment,
//                      * due to lack of input.
//                      */
//                     if (ip+length > iend) {
//                         length = (size_t)(iend-ip);
//                         cpy = op + length;
//                     }
//                     /* Finishing in the middle of a literals segment,
//                      * due to lack of output space.
//                      */
//                     if (cpy > oend) {
//                         cpy = oend;
//                         assert(op<=oend);
//                         length = (size_t)(oend-op);
//                     }
//                 } else {
//                     /* We must be on the last sequence because of the parsing limitations so check
//                      * that we exactly regenerate the original size (must be exact when !endOnInput).
//                      */
//                     if ((!endOnInput) && (cpy != oend)) { goto _output_error; }
//                      /* We must be on the last sequence (or invalid) because of the parsing limitations
//                       * so check that we exactly consume the input and don't overrun the output buffer.
//                       */
//                     if ((endOnInput) && ((ip+length != iend) || (cpy > oend))) {
//                         DEBUGLOG(6, "should have been last run of literals")
//                         DEBUGLOG(6, "ip(%p) + length(%i) = %p != iend (%p)", ip, (int)length, ip+length, iend);
//                         DEBUGLOG(6, "or cpy(%p) > oend(%p)", cpy, oend);
//                         goto _output_error;
//                     }
//                 }
//                 memmove(op, ip, length);  /* supports overlapping memory regions; only matters for in-place decompression scenarios */
//                 ip += length;
//                 op += length;
//                 /* Necessarily EOF when !partialDecoding.
//                  * When partialDecoding, it is EOF if we've either
//                  * filled the output buffer or
//                  * can't proceed with reading an offset for following match.
//                  */
//                 if (!partialDecoding || (cpy == oend) || (ip >= (iend-2))) {
//                     break;
//                 }
//             } else {
//                 LZ4_wildCopy8(op, ip, cpy);   /* may overwrite up to WILDCOPYLENGTH beyond cpy */
//                 ip += length; op = cpy;
//             }
// 
//             /* get offset */
//             offset = LZ4_readLE16(ip); ip+=2;
//             match = op - offset;
// 
//             /* get matchlength */
//             length = token & ML_MASK;
// 
//     _copy_match:
//             if (length == ML_MASK) {
//               variable_length_error error = ok;
//               length += read_variable_length(&ip, iend - LASTLITERALS + 1, (int)endOnInput, 0, &error);
//               if (error != ok) goto _output_error;
//                 if ((safeDecode) && unlikely((uptrval)(op)+length<(uptrval)op)) goto _output_error;   /* overflow detection */
//             }
//             length += MINMATCH;
// 
// #if LZ4_FAST_DEC_LOOP
//         safe_match_copy:
// #endif
//             if ((checkOffset) && (unlikely(match + dictSize < lowPrefix))) goto _output_error;   /* Error : offset outside buffers */
//             /* match starting within external dictionary */
//             if ((dict==usingExtDict) && (match < lowPrefix)) {
//                 if (unlikely(op+length > oend-LASTLITERALS)) {
//                     if (partialDecoding) length = MIN(length, (size_t)(oend-op));
//                     else goto _output_error;   /* doesn't respect parsing restriction */
//                 }
// 
//                 if (length <= (size_t)(lowPrefix-match)) {
//                     /* match fits entirely within external dictionary : just copy */
//                     memmove(op, dictEnd - (lowPrefix-match), length);
//                     op += length;
//                 } else {
//                     /* match stretches into both external dictionary and current block */
//                     size_t const copySize = (size_t)(lowPrefix - match);
//                     size_t const restSize = length - copySize;
//                     LZ4_memcpy(op, dictEnd - copySize, copySize);
//                     op += copySize;
//                     if (restSize > (size_t)(op - lowPrefix)) {  /* overlap copy */
//                         BYTE* const endOfMatch = op + restSize;
//                         const BYTE* copyFrom = lowPrefix;
//                         while (op < endOfMatch) *op++ = *copyFrom++;
//                     } else {
//                         LZ4_memcpy(op, lowPrefix, restSize);
//                         op += restSize;
//                 }   }
//                 continue;
//             }
//             assert(match >= lowPrefix);
// 
//             /* copy match within block */
//             cpy = op + length;
// 
//             /* partialDecoding : may end anywhere within the block */
//             assert(op<=oend);
//             if (partialDecoding && (cpy > oend-MATCH_SAFEGUARD_DISTANCE)) {
//                 size_t const mlen = MIN(length, (size_t)(oend-op));
//                 const BYTE* const matchEnd = match + mlen;
//                 BYTE* const copyEnd = op + mlen;
//                 if (matchEnd > op) {   /* overlap copy */
//                     while (op < copyEnd) { *op++ = *match++; }
//                 } else {
//                     LZ4_memcpy(op, match, mlen);
//                 }
//                 op = copyEnd;
//                 if (op == oend) { break; }
//                 continue;
//             }
// 
//             if (unlikely(offset<8)) {
//                 LZ4_write32(op, 0);   /* silence msan warning when offset==0 */
//                 op[0] = match[0];
//                 op[1] = match[1];
//                 op[2] = match[2];
//                 op[3] = match[3];
//                 match += inc32table[offset];
//                 LZ4_memcpy(op+4, match, 4);
//                 match -= dec64table[offset];
//             } else {
//                 LZ4_memcpy(op, match, 8);
//                 match += 8;
//             }
//             op += 8;
// 
//             if (unlikely(cpy > oend-MATCH_SAFEGUARD_DISTANCE)) {
//                 BYTE* const oCopyLimit = oend - (WILDCOPYLENGTH-1);
//                 if (cpy > oend-LASTLITERALS) { goto _output_error; } /* Error : last LASTLITERALS bytes must be literals (uncompressed) */
//                 if (op < oCopyLimit) {
//                     LZ4_wildCopy8(op, match, oCopyLimit);
//                     match += oCopyLimit - op;
//                     op = oCopyLimit;
//                 }
//                 while (op < cpy) { *op++ = *match++; }
//             } else {
//                 LZ4_memcpy(op, match, 8);
//                 if (length > 16)  { LZ4_wildCopy8(op+8, match+8, cpy); }
//             }
//             op = cpy;   /* wildcopy correction */
//         }
// 
//         /* end of decoding */
//         if (endOnInput) {
//             DEBUGLOG(5, "decoded %i bytes", (int) (((char*)op)-dst));
//            return (int) (((char*)op)-dst);     /* Nb of output bytes decoded */
//        } else {
//            return (int) (((const char*)ip)-src);   /* Nb of input bytes read */
//        }
// 
//         /* Overflow error detected */
//     _output_error:
//         return (int) (-(((const char*)ip)-src))-1;
//     }
// }
// 
// 
// /*===== Instantiate the API decoding functions. =====*/
// 
// LZ4_FORCE_O2
// int LZ4_decompress_safe(const char* source, char* dest, int compressedSize, int maxDecompressedSize)
// {
//     return LZ4_decompress_generic(source, dest, compressedSize, maxDecompressedSize,
//                                   endOnInputSize, decode_full_block, noDict,
//                                   (BYTE*)dest, NULL, 0);
// }
// 
// LZ4_FORCE_O2
// int LZ4_decompress_safe_partial(const char* src, char* dst, int compressedSize, int targetOutputSize, int dstCapacity)
// {
//     dstCapacity = MIN(targetOutputSize, dstCapacity);
//     return LZ4_decompress_generic(src, dst, compressedSize, dstCapacity,
//                                   endOnInputSize, partial_decode,
//                                   noDict, (BYTE*)dst, NULL, 0);
// }
// 
// LZ4_FORCE_O2
// int LZ4_decompress_fast(const char* source, char* dest, int originalSize)
// {
//     return LZ4_decompress_generic(source, dest, 0, originalSize,
//                                   endOnOutputSize, decode_full_block, withPrefix64k,
//                                   (BYTE*)dest - 64 KB, NULL, 0);
// }
// 
// /*===== Instantiate a few more decoding cases, used more than once. =====*/
// 
// LZ4_FORCE_O2 /* Exported, an obsolete API function. */
// int LZ4_decompress_safe_withPrefix64k(const char* source, char* dest, int compressedSize, int maxOutputSize)
// {
//     return LZ4_decompress_generic(source, dest, compressedSize, maxOutputSize,
//                                   endOnInputSize, decode_full_block, withPrefix64k,
//                                   (BYTE*)dest - 64 KB, NULL, 0);
// }
// 
// /* Another obsolete API function, paired with the previous one. */
// int LZ4_decompress_fast_withPrefix64k(const char* source, char* dest, int originalSize)
// {
//     /* LZ4_decompress_fast doesn't validate match offsets,
//      * and thus serves well with any prefixed dictionary. */
//     return LZ4_decompress_fast(source, dest, originalSize);
// }
// 
// LZ4_FORCE_O2
// static int LZ4_decompress_safe_withSmallPrefix(const char* source, char* dest, int compressedSize, int maxOutputSize,
//                                                size_t prefixSize)
// {
//     return LZ4_decompress_generic(source, dest, compressedSize, maxOutputSize,
//                                   endOnInputSize, decode_full_block, noDict,
//                                   (BYTE*)dest-prefixSize, NULL, 0);
// }
// 
// LZ4_FORCE_O2
// int LZ4_decompress_safe_forceExtDict(const char* source, char* dest,
//                                      int compressedSize, int maxOutputSize,
//                                      const void* dictStart, size_t dictSize)
// {
//     return LZ4_decompress_generic(source, dest, compressedSize, maxOutputSize,
//                                   endOnInputSize, decode_full_block, usingExtDict,
//                                   (BYTE*)dest, (const BYTE*)dictStart, dictSize);
// }
// 
// LZ4_FORCE_O2
// static int LZ4_decompress_fast_extDict(const char* source, char* dest, int originalSize,
//                                        const void* dictStart, size_t dictSize)
// {
//     return LZ4_decompress_generic(source, dest, 0, originalSize,
//                                   endOnOutputSize, decode_full_block, usingExtDict,
//                                   (BYTE*)dest, (const BYTE*)dictStart, dictSize);
// }
// 
// /* The "double dictionary" mode, for use with e.g. ring buffers: the first part
//  * of the dictionary is passed as prefix, and the second via dictStart + dictSize.
//  * These routines are used only once, in LZ4_decompress_*_continue().
//  */
// LZ4_FORCE_INLINE
// int LZ4_decompress_safe_doubleDict(const char* source, char* dest, int compressedSize, int maxOutputSize,
//                                    size_t prefixSize, const void* dictStart, size_t dictSize)
// {
//     return LZ4_decompress_generic(source, dest, compressedSize, maxOutputSize,
//                                   endOnInputSize, decode_full_block, usingExtDict,
//                                   (BYTE*)dest-prefixSize, (const BYTE*)dictStart, dictSize);
// }
// 
// LZ4_FORCE_INLINE
// int LZ4_decompress_fast_doubleDict(const char* source, char* dest, int originalSize,
//                                    size_t prefixSize, const void* dictStart, size_t dictSize)
// {
//     return LZ4_decompress_generic(source, dest, 0, originalSize,
//                                   endOnOutputSize, decode_full_block, usingExtDict,
//                                   (BYTE*)dest-prefixSize, (const BYTE*)dictStart, dictSize);
// }
// 
// /*===== streaming decompression functions =====*/
// 
// LZ4_streamDecode_t* LZ4_createStreamDecode(void)
// {
//     LZ4_streamDecode_t* lz4s = (LZ4_streamDecode_t*) ALLOC_AND_ZERO(sizeof(LZ4_streamDecode_t));
//     LZ4_STATIC_ASSERT(LZ4_STREAMDECODESIZE >= sizeof(LZ4_streamDecode_t_internal));    /* A compilation error here means LZ4_STREAMDECODESIZE is not large enough */
//     return lz4s;
// }
// 
// int LZ4_freeStreamDecode (LZ4_streamDecode_t* LZ4_stream)
// {
//     if (LZ4_stream == NULL) { return 0; }  /* support free on NULL */
//     FREEMEM(LZ4_stream);
//     return 0;
// }
// 
// /*! LZ4_setStreamDecode() :
//  *  Use this function to instruct where to find the dictionary.
//  *  This function is not necessary if previous data is still available where it was decoded.
//  *  Loading a size of 0 is allowed (same effect as no dictionary).
//  * @return : 1 if OK, 0 if error
//  */
// int LZ4_setStreamDecode (LZ4_streamDecode_t* LZ4_streamDecode, const char* dictionary, int dictSize)
// {
//     LZ4_streamDecode_t_internal* lz4sd = &LZ4_streamDecode->internal_donotuse;
//     lz4sd->prefixSize = (size_t) dictSize;
//     lz4sd->prefixEnd = (const BYTE*) dictionary + dictSize;
//     lz4sd->externalDict = NULL;
//     lz4sd->extDictSize  = 0;
//     return 1;
// }
// 
// /*! LZ4_decoderRingBufferSize() :
//  *  when setting a ring buffer for streaming decompression (optional scenario),
//  *  provides the minimum size of this ring buffer
//  *  to be compatible with any source respecting maxBlockSize condition.
//  *  Note : in a ring buffer scenario,
//  *  blocks are presumed decompressed next to each other.
//  *  When not enough space remains for next block (remainingSize < maxBlockSize),
//  *  decoding resumes from beginning of ring buffer.
//  * @return : minimum ring buffer size,
//  *           or 0 if there is an error (invalid maxBlockSize).
//  */
// int LZ4_decoderRingBufferSize(int maxBlockSize)
// {
//     if (maxBlockSize < 0) return 0;
//     if (maxBlockSize > LZ4_MAX_INPUT_SIZE) return 0;
//     if (maxBlockSize < 16) maxBlockSize = 16;
//     return LZ4_DECODER_RING_BUFFER_SIZE(maxBlockSize);
// }
// 
// /*
// *_continue() :
//     These decoding functions allow decompression of multiple blocks in "streaming" mode.
//     Previously decoded blocks must still be available at the memory position where they were decoded.
//     If it's not possible, save the relevant part of decoded data into a safe buffer,
//     and indicate where it stands using LZ4_setStreamDecode()
// */
// LZ4_FORCE_O2
// int LZ4_decompress_safe_continue (LZ4_streamDecode_t* LZ4_streamDecode, const char* source, char* dest, int compressedSize, int maxOutputSize)
// {
//     LZ4_streamDecode_t_internal* lz4sd = &LZ4_streamDecode->internal_donotuse;
//     int result;
// 
//     if (lz4sd->prefixSize == 0) {
//         /* The first call, no dictionary yet. */
//         assert(lz4sd->extDictSize == 0);
//         result = LZ4_decompress_safe(source, dest, compressedSize, maxOutputSize);
//         if (result <= 0) return result;
//         lz4sd->prefixSize = (size_t)result;
//         lz4sd->prefixEnd = (BYTE*)dest + result;
//     } else if (lz4sd->prefixEnd == (BYTE*)dest) {
//         /* They're rolling the current segment. */
//         if (lz4sd->prefixSize >= 64 KB - 1)
//             result = LZ4_decompress_safe_withPrefix64k(source, dest, compressedSize, maxOutputSize);
//         else if (lz4sd->extDictSize == 0)
//             result = LZ4_decompress_safe_withSmallPrefix(source, dest, compressedSize, maxOutputSize,
//                                                          lz4sd->prefixSize);
//         else
//             result = LZ4_decompress_safe_doubleDict(source, dest, compressedSize, maxOutputSize,
//                                                     lz4sd->prefixSize, lz4sd->externalDict, lz4sd->extDictSize);
//         if (result <= 0) return result;
//         lz4sd->prefixSize += (size_t)result;
//         lz4sd->prefixEnd  += result;
//     } else {
//         /* The buffer wraps around, or they're switching to another buffer. */
//         lz4sd->extDictSize = lz4sd->prefixSize;
//         lz4sd->externalDict = lz4sd->prefixEnd - lz4sd->extDictSize;
//         result = LZ4_decompress_safe_forceExtDict(source, dest, compressedSize, maxOutputSize,
//                                                   lz4sd->externalDict, lz4sd->extDictSize);
//         if (result <= 0) return result;
//         lz4sd->prefixSize = (size_t)result;
//         lz4sd->prefixEnd  = (BYTE*)dest + result;
//     }
// 
//     return result;
// }
// 
// LZ4_FORCE_O2
// int LZ4_decompress_fast_continue (LZ4_streamDecode_t* LZ4_streamDecode, const char* source, char* dest, int originalSize)
// {
//     LZ4_streamDecode_t_internal* lz4sd = &LZ4_streamDecode->internal_donotuse;
//     int result;
//     assert(originalSize >= 0);
// 
//     if (lz4sd->prefixSize == 0) {
//         assert(lz4sd->extDictSize == 0);
//         result = LZ4_decompress_fast(source, dest, originalSize);
//         if (result <= 0) return result;
//         lz4sd->prefixSize = (size_t)originalSize;
//         lz4sd->prefixEnd = (BYTE*)dest + originalSize;
//     } else if (lz4sd->prefixEnd == (BYTE*)dest) {
//         if (lz4sd->prefixSize >= 64 KB - 1 || lz4sd->extDictSize == 0)
//             result = LZ4_decompress_fast(source, dest, originalSize);
//         else
//             result = LZ4_decompress_fast_doubleDict(source, dest, originalSize,
//                                                     lz4sd->prefixSize, lz4sd->externalDict, lz4sd->extDictSize);
//         if (result <= 0) return result;
//         lz4sd->prefixSize += (size_t)originalSize;
//         lz4sd->prefixEnd  += originalSize;
//     } else {
//         lz4sd->extDictSize = lz4sd->prefixSize;
//         lz4sd->externalDict = lz4sd->prefixEnd - lz4sd->extDictSize;
//         result = LZ4_decompress_fast_extDict(source, dest, originalSize,
//                                              lz4sd->externalDict, lz4sd->extDictSize);
//         if (result <= 0) return result;
//         lz4sd->prefixSize = (size_t)originalSize;
//         lz4sd->prefixEnd  = (BYTE*)dest + originalSize;
//     }
// 
//     return result;
// }
// 
// 
// /*
// Advanced decoding functions :
// *_usingDict() :
//     These decoding functions work the same as "_continue" ones,
//     the dictionary must be explicitly provided within parameters
// */
// 
// int LZ4_decompress_safe_usingDict(const char* source, char* dest, int compressedSize, int maxOutputSize, const char* dictStart, int dictSize)
// {
//     if (dictSize==0)
//         return LZ4_decompress_safe(source, dest, compressedSize, maxOutputSize);
//     if (dictStart+dictSize == dest) {
//         if (dictSize >= 64 KB - 1) {
//             return LZ4_decompress_safe_withPrefix64k(source, dest, compressedSize, maxOutputSize);
//         }
//         assert(dictSize >= 0);
//         return LZ4_decompress_safe_withSmallPrefix(source, dest, compressedSize, maxOutputSize, (size_t)dictSize);
//     }
//     assert(dictSize >= 0);
//     return LZ4_decompress_safe_forceExtDict(source, dest, compressedSize, maxOutputSize, dictStart, (size_t)dictSize);
// }
// 
// int LZ4_decompress_fast_usingDict(const char* source, char* dest, int originalSize, const char* dictStart, int dictSize)
// {
//     if (dictSize==0 || dictStart+dictSize == dest)
//         return LZ4_decompress_fast(source, dest, originalSize);
//     assert(dictSize >= 0);
//     return LZ4_decompress_fast_extDict(source, dest, originalSize, dictStart, (size_t)dictSize);
// }
// 
// 
// /*=*************************************************
// *  Obsolete Functions
// ***************************************************/
// /* obsolete compression functions */
// int LZ4_compress_limitedOutput(const char* source, char* dest, int inputSize, int maxOutputSize)
// {
//     return LZ4_compress_default(source, dest, inputSize, maxOutputSize);
// }
// int LZ4_compress(const char* src, char* dest, int srcSize)
// {
//     return LZ4_compress_default(src, dest, srcSize, LZ4_compressBound(srcSize));
// }
// int LZ4_compress_limitedOutput_withState (void* state, const char* src, char* dst, int srcSize, int dstSize)
// {
//     return LZ4_compress_fast_extState(state, src, dst, srcSize, dstSize, 1);
// }
// int LZ4_compress_withState (void* state, const char* src, char* dst, int srcSize)
// {
//     return LZ4_compress_fast_extState(state, src, dst, srcSize, LZ4_compressBound(srcSize), 1);
// }
// int LZ4_compress_limitedOutput_continue (LZ4_stream_t* LZ4_stream, const char* src, char* dst, int srcSize, int dstCapacity)
// {
//     return LZ4_compress_fast_continue(LZ4_stream, src, dst, srcSize, dstCapacity, 1);
// }
// int LZ4_compress_continue (LZ4_stream_t* LZ4_stream, const char* source, char* dest, int inputSize)
// {
//     return LZ4_compress_fast_continue(LZ4_stream, source, dest, inputSize, LZ4_compressBound(inputSize), 1);
// }
// 
// /*
// These decompression functions are deprecated and should no longer be used.
// They are only provided here for compatibility with older user programs.
// - LZ4_uncompress is totally equivalent to LZ4_decompress_fast
// - LZ4_uncompress_unknownOutputSize is totally equivalent to LZ4_decompress_safe
// */
// int LZ4_uncompress (const char* source, char* dest, int outputSize)
// {
//     return LZ4_decompress_fast(source, dest, outputSize);
// }
// int LZ4_uncompress_unknownOutputSize (const char* source, char* dest, int isize, int maxOutputSize)
// {
//     return LZ4_decompress_safe(source, dest, isize, maxOutputSize);
// }
// 
// /* Obsolete Streaming functions */
// 
// int LZ4_sizeofStreamState(void) { return LZ4_STREAMSIZE; }
// 
// int LZ4_resetStreamState(void* state, char* inputBuffer)
// {
//     (void)inputBuffer;
//     LZ4_resetStream((LZ4_stream_t*)state);
//     return 0;
// }
// 
// void* LZ4_create (char* inputBuffer)
// {
//     (void)inputBuffer;
//     return LZ4_createStream();
// }
// 
// char* LZ4_slideInputBuffer (void* state)
// {
//     /* avoid const char * -> char * conversion warning */
//     return (char *)(uptrval)((LZ4_stream_t*)state)->internal_donotuse.dictionary;
// }
// 
// #endif   /* LZ4_COMMONDEFS_ONLY */
// //
// //  sqcloud.c
// //
// //  Created by Marco Bambini on 08/02/21.
// //
// 
// //
// //  sqcloud.h
// //
// //  Created by Marco Bambini on 08/02/21.
// //
// 
// #ifndef __SQCLOUD_CLI__
// #define __SQCLOUD_CLI__
// 
// #include <stdio.h>
// #include <stdbool.h>
// 
// #ifdef __cplusplus
// extern "C" {
// #endif
// 
// #define SQCLOUD_SDK_VERSION         "0.4.1"
// #define SQCLOUD_SDK_VERSION_NUM     0x000401
// #define SQCLOUD_DEFAULT_PORT        8860
// #define SQCLOUD_DEFAULT_TIMEOUT     12
// 
// // opaque datatypes
// typedef struct SQCloudConnection    SQCloudConnection;
// typedef struct SQCloudResult        SQCloudResult;
// typedef void (*SQCloudPubSubCB)    (SQCloudConnection *connection, SQCloudResult *result, void *data);
// 
// // configuration struct to be passed to the connect function (currently unused)
// typedef struct SQCloudConfigStruct {
//     const char *username;
//     const char *password;
//     const char *database;
//     int timeout;
//     int family;                 // can be: AF_INET, AF_INET6 or AF_UNSPEC
// } SQCloudConfig;
// 
// typedef enum {
//     RESULT_OK,
//     RESULT_ERROR,
//     RESULT_STRING,
//     RESULT_INTEGER,
//     RESULT_FLOAT,
//     RESULT_ROWSET,
//     RESULT_NULL,
//     RESULT_JSON
// } SQCloudResType;
// 
// typedef enum {
//     VALUE_INTEGER = 1,
//     VALUE_FLOAT = 2,
//     VALUE_TEXT = 3,
//     VALUE_BLOB = 4,
//     VALUE_NULL = 5
// } SQCloudValueType;
// 
// SQCloudConnection *SQCloudConnect (const char *hostname, int port, SQCloudConfig *config);
// SQCloudConnection *SQCloudConnectWithString (const char *s);
// SQCloudResult *SQCloudExec (SQCloudConnection *connection, const char *command);
// char *SQCloudUUID (SQCloudConnection *connection);
// void SQCloudDisconnect (SQCloudConnection *connection);
// void SQCloudSetPubSubCallback (SQCloudConnection *connection, SQCloudPubSubCB callback, void *data);
// SQCloudResult *SQCloudSetPubSubOnly (SQCloudConnection *connection);
// 
// bool SQCloudIsError (SQCloudConnection *connection);
// int SQCloudErrorCode (SQCloudConnection *connection);
// const char *SQCloudErrorMsg (SQCloudConnection *connection);
// 
// SQCloudResType SQCloudResultType (SQCloudResult *result);
// uint32_t SQCloudResultLen (SQCloudResult *result);
// char *SQCloudResultBuffer (SQCloudResult *result);
// void SQCloudResultFree (SQCloudResult *result);
// bool SQCloudResultIsOK (SQCloudResult *result);
// 
// SQCloudValueType SQCloudRowsetValueType (SQCloudResult *result, uint32_t row, uint32_t col);
// uint32_t SQCloudRowsetRowsMaxColumnLength (SQCloudResult *result, uint32_t col);
// char *SQCloudRowsetColumnName (SQCloudResult *result, uint32_t col, uint32_t *len);
// uint32_t SQCloudRowsetRows (SQCloudResult *result);
// uint32_t SQCloudRowsetCols (SQCloudResult *result);
// uint32_t SQCloudRowsetMaxLen (SQCloudResult *result);
// char *SQCloudRowsetValue (SQCloudResult *result, uint32_t row, uint32_t col, uint32_t *len);
// int32_t SQCloudRowsetInt32Value (SQCloudResult *result, uint32_t row, uint32_t col);
// int64_t SQCloudRowsetInt64Value (SQCloudResult *result, uint32_t row, uint32_t col);
// float SQCloudRowsetFloatValue (SQCloudResult *result, uint32_t row, uint32_t col);
// double SQCloudRowsetDoubleValue (SQCloudResult *result, uint32_t row, uint32_t col);
// void SQCloudRowsetDump (SQCloudResult *result, uint32_t maxline);
// 
// #ifdef __cplusplus
// }
// #endif
// 
// #endif
// #include <ctype.h>
// #include <stdlib.h>
// #include <string.h>
// #include <stdarg.h>
// #include <assert.h>
// #include <sys/time.h>
// 
// #ifdef _WIN32
// #include <winsock2.h>
// #include <ws2tcpip.h>
// #pragma comment(lib, "Ws2_32.lib")
// #include <Shlwapi.h>
// #include <io.h>
// #include <float.h>
// #include "pthread.h"
// #else
// #include <errno.h>
// #include <netdb.h>
// #include <signal.h>
// #include <unistd.h>
// #include <netinet/in.h>
// #include <sys/socket.h>
// #include <sys/stat.h>
// #include <sys/types.h>
// #include <sys/wait.h>
// #include <arpa/inet.h>
// #include <netinet/tcp.h>
// #include <sys/ioctl.h>
// #include <pthread.h>
// #endif
// 
// // MARK: MACROS -
// #ifdef _WIN32
// #pragma warning (disable: 4005)
// #pragma warning (disable: 4068)
// #define readsocket(a,b,c)                   recv((a), (b), (c), 0L)
// #define writesocket(a,b,c)                  send((a), (b), (c), 0L)
// #else
// #define readsocket                          read
// #define writesocket                         write
// #define closesocket(s)                      close(s)
// #endif
// 
// #define mem_realloc                         realloc
// #define mem_zeroalloc(_s)                   calloc(1,_s)
// #define mem_alloc(_s)                       malloc(_s)
// #define mem_free(_s)                        free(_s)
// #define string_dup(_s)                      strdup(_s)
// #define MIN(a,b)                            (((a)<(b))?(a):(b))
// 
// #define MAX_SOCK_LIST                       6           // maximum number of socket descriptor to try to connect to
//                                                         // this change is required to support IPv4/IPv6 connections
// #define DEFAULT_TIMEOUT                     12          // default connection timeout in seconds
// 
// #define REPLY_OK                            "+2 OK"     // default OK reply
// #define REPLY_OK_LEN                        5           // default OK reply string length
// 
// // https://levelup.gitconnected.com/8-ways-to-measure-execution-time-in-c-c-48634458d0f9
// #define TIME_GET(_t1)                       struct timeval _t1; gettimeofday(&_t1, NULL)
// #define TIME_VAL(_t1, _t2)                  ((double)(_t2.tv_sec - _t1.tv_sec) + (double)((_t2.tv_usec - _t1.tv_usec)*1e-6))
// 
// #define CMD_STRING                          '+'
// #define CMD_ZEROSTRING                      '!'
// #define CMD_ERROR                           '-'
// #define CMD_INT                             ':'
// #define CMD_FLOAT                           ','
// #define CMD_ROWSET                          '*'
// #define CMD_ROWSET_CHUNK                    '/'
// #define CMD_JSON                            '#'
// #define CMD_RAWJSON                         '{'
// #define CMD_NULL                            '_'
// #define CMD_BLOB                            '$'
// #define CMD_COMPRESSED                      '%'
// #define CMD_PUBSUB                          '|'
// #define CMD_COMMAND                         '^'
// #define CMD_RECONNECT                       '@'
// 
// #define CMD_MINLEN                          2
// 
// #define CONNSTRING_KEYVALUE_SEPARATOR       '='
// #define CONNSTRING_TOKEN_SEPARATOR          ';'
// 
// #define DEFAULT_CHUCK_NBUFFERS              20
// #define DEFAULT_CHUNK_MINROWS               2000
// 
// // MARK: - PROTOTYPES -
// 
// static SQCloudResult *internal_socket_read (SQCloudConnection *connection, bool mainfd);
// static bool internal_socket_write (SQCloudConnection *connection, const char *buffer, size_t len, bool mainfd);
// static uint32_t internal_parse_number (char *buffer, uint32_t blen, uint32_t *cstart);
// static SQCloudResult *internal_parse_buffer (SQCloudConnection *connection, char *buffer, uint32_t blen, uint32_t cstart, bool isstatic, bool externalbuffer);
// static bool internal_connect (SQCloudConnection *connection, const char *hostname, int port, SQCloudConfig *config, bool mainfd);
// 
// // MARK: -
// 
// struct SQCloudResult {
//     SQCloudResType  tag;                    // RESULT_OK, RESULT_ERROR, RESULT_STRING, RESULT_INTEGER, RESULT_FLOAT, RESULT_ROWSET, RESULT_NULL
//     
//     bool            ischunk;                // flag used to correctly access the union below
//     union {
//         struct {
//             char        *buffer;            // buffer used by the user (it could be a ptr inside rawbuffer)
//             char        *rawbuffer;         // ptr to the buffer to be freed
//             uint32_t    balloc;             // buffer allocation size
//         };
//         struct {
//             char        **buffers;          // array of buffers used by rowset sent in chunk
//             uint32_t    bcount;             // number of buffers in the array
//             uint32_t    bnum;               // number of pre-allocated buffers
//             uint32_t    brows;              // number of pre-allocated rows
//         };
//     };
//     
//     // common
//     uint32_t        blen;                   // total buffer length (also the sum of buffers)
//     double          time;                   // full execution time (latency + server side time)
//     bool            externalbuffer;         // true if the buffer is managed by the caller code
//                                             // false if the buffer can be freed by the SQCloudResultFree func
//     
//     // used in TYPE_ROWSET only
//     uint32_t        nrows;                  // number of rows
//     uint32_t        ncols;                  // number of columns
//     uint32_t        ndata;                  // number of items stores in data
//     char            **data;                 // data contained in the rowset
//     char            **name;                 // column names
//     uint32_t        *clen;                  // max len for each column (used to display result)
//     uint32_t        maxlen;                 // max len for each row/column
// } _SQCloudResult;
// 
// struct SQCloudConnection {
//     int             fd;
//     char            errmsg[1024];
//     int             errcode;
//     SQCloudResult   *_chunk;
//     
//     // pub/sub
//     char            *uuid;
//     int             pubsubfd;
//     SQCloudPubSubCB callback;
//     void            *data;
//     char            *hostname;
//     int             port;
//     pthread_t       tid;
//     
// } _SQCloudConnection;
// 
// static SQCloudResult SQCloudResultOK = {RESULT_OK, NULL, 0, 0, 0};
// static SQCloudResult SQCloudResultNULL = {RESULT_NULL, NULL, 0, 0, 0};
// 
// // MARK: - UTILS -
// 
// static uint32_t utf8_charbytes (const char *s, uint32_t i) {
//     unsigned char c = (unsigned char)s[i];
//     
//     // determine bytes needed for character, based on RFC 3629
//     if ((c > 0) && (c <= 127)) return 1;
//     if ((c >= 194) && (c <= 223)) return 2;
//     if ((c >= 224) && (c <= 239)) return 3;
//     if ((c >= 240) && (c <= 244)) return 4;
//     
//     // means error
//     return 0;
// }
// 
// static uint32_t utf8_len (const char *s, uint32_t nbytes) {
//     uint32_t pos = 0;
//     uint32_t len = 0;
//     
//     while (pos < nbytes) {
//         ++len;
//         uint32_t n = utf8_charbytes(s, pos);
//         if (n == 0) return 0; // means error
//         pos += n;
//     }
//     
//     return len;
// }
// 
// #if 0
// static char *extract_connection_token (const char *s, char *key, char buffer[256]) {
//     char *target = strstr(s, key);
//     if (!target) return NULL;
//     
//     // find out = separator
//     char *p = target;
//     while (p[0]) {
//         if (p[0] == CONNSTRING_KEYVALUE_SEPARATOR) break;
//         ++p;
//     }
//     
//     // skip =
//     ++p;
//     
//     // skip spaces (if any)
//     while (p[0]) {
//         if (!isspace(p[0])) break;
//         ++p;
//     }
//     
//     // copy value to buffer
//     int len = 0;
//     while (p[0] && len < 255) {
//         if (isspace(p[0])) break;
//         if (p[0] == CONNSTRING_TOKEN_SEPARATOR) break;
//         buffer[len] = p[0];
//         ++len;
//         ++p;
//     }
//     
//     // null-terminate returning value
//     buffer[len] = 0;
//     p = &buffer[0];
//     
//     return p;
// }
// #endif
// 
// // MARK: - PRIVATE -
// 
// static int socket_geterror (int fd) {
//     int err;
//     socklen_t errlen = sizeof(err);
//     
//     int sockerr = getsockopt(fd, SOL_SOCKET, SO_ERROR, (void *)&err, &errlen);
//     if (sockerr < 0) return -1;
//     
//     return ((err == 0 || err == EINTR || err == EAGAIN || err == EINPROGRESS)) ? 0 : err;
// }
// 
// static void *pubsub_thread (void *arg) {
//     SQCloudConnection *connection = (SQCloudConnection *)arg;
//     int fd = connection->pubsubfd;
//     
//     size_t blen = 2048;
//     char *buffer = mem_alloc(blen);
//     if (buffer == NULL) return NULL;
//     
//     char *original = buffer;
//     uint32_t tread = 0;
// 
//     while (1) {
//         fd_set set;
//         FD_ZERO(&set);
//         FD_SET(fd, &set);
//         
//         // wait for read event
//         int rc = select(fd + 1, &set, NULL, NULL, NULL);
//         if (rc <= 0) continue;
//         
//         //  read payload string
//         ssize_t nread = readsocket(fd, buffer, blen);
//         
//         if (nread < 0) {
//             printf("Handle error here %s.", strerror(errno));
//             break;
//             // internal_set_error(connection, 1, "An error occurred while reading data: %s.", strerror(errno));
//             // goto abort_read;
//         }
//         
//         if (nread == 0) {
//             printf("Handle error here %s.", strerror(errno));
//             break;
//             // internal_set_error(connection, 1, "Unexpected EOF found while reading data: %s.", strerror(errno));
//             // goto abort_read;
//         }
//         
//         tread += (uint32_t)nread;
//         blen -= (uint32_t)nread;
//         buffer += nread;
//         
//         uint32_t cstart = 0;
//         uint32_t clen = internal_parse_number (&original[1], tread-1, &cstart);
//         if (clen == 0) continue;
//         
//         // check if read is complete
//         // clen is the lenght parsed in the buffer
//         // cstart is the index of the first space
//         // +1 because we skipped the first character in the internal_parse_number function
//         if (clen + cstart + 1 != tread) {
//             // check buffer allocation and continue reading
//             if (clen + cstart > blen) {
//                 char *clone = mem_alloc(clen + cstart + 1);
//                 if (!clone) {
//                     printf("Handle memory error here %s.", strerror(errno));
//                     break;
//                     // internal_set_error(connection, 1, "Unable to allocate memory: %d.", clen + cstart + 1);
//                     // goto abort_read;
//                 }
//                 memcpy(clone, original, tread);
//                 buffer = original = clone;
//                 blen = (clen + cstart + 1) - tread;
//                 buffer += tread;
//             }
//             
//             continue;
//         }
//         
//         SQCloudResult *result = internal_parse_buffer(connection, original, tread, (clen) ? cstart : 0, false, false);
//         if (result->tag == RESULT_STRING) result->tag = RESULT_JSON;
//         
//         connection->callback(connection, result, connection->data);
//         
//         blen = 2048;
//         buffer = mem_alloc(blen);
//         if (!buffer) break;
//         
//         original = buffer;
//         tread = 0;
//     }
//     
//     return NULL;
// }
// 
// // MARK: -
// 
// static bool internal_init (void) {
//     static bool inited = false;
//     if (inited) return true;
//     
//     #ifdef _WIN32
//     WSADATA wsaData;
//     WSAStartup(MAKEWORD(2,2), &wsaData);
//     #else
//     // IGNORE SIGPIPE and SIGABORT
//     struct sigaction act;
//     act.sa_handler = SIG_IGN;
//     sigemptyset(&act.sa_mask);
//     act.sa_flags = 0;
//     sigaction(SIGPIPE, &act, (struct sigaction *)NULL);
//     sigaction(SIGABRT, &act, (struct sigaction *)NULL);
//     #endif
//     
//     inited = true;
//     return true;
// }
// 
// static bool internal_set_error (SQCloudConnection *connection, int errcode, const char *format, ...) {
//     connection->errcode = errcode;
//     
//     va_list arg;
//     va_start (arg, format);
//     vsnprintf(connection->errmsg, sizeof(connection->errmsg), format, arg);
//     va_end (arg);
//     
//     return false;
// }
// 
// static void internal_parse_uuid (SQCloudConnection *connection, const char *buffer, size_t blen) {
//     // sanity check
//     if (!buffer || blen == 0) return;
//     
//     // expected buffer is PAUTH uuid secret
//     // PUATH -> 5
//     // uuid -> 36
//     // secret -> 36
//     // spaces -> 2
//     if (blen != (5 + 36 + 36 + 2)) return;
//     
//     if (strncmp(buffer, "PAUTH ", 6) != 0) return;
//     
//     // allocate 36 (UUID) + 1 (null-terminated) zero-bytes
//     char *uuid = mem_zeroalloc(37);
//     if (!uuid) return;
//     
//     memcpy(uuid, &buffer[6], 36);
//     connection->uuid = uuid;
// }
// 
// static void internal_clear_error (SQCloudConnection *connection) {
//     connection->errcode = 0;
//     connection->errmsg[0] = 0;
// }
// 
// static bool internal_setup_ssl (SQCloudConnection *connection, SQCloudConfig *config) {
//     return true;
// }
// 
// static SQCloudValueType internal_type (char *buffer) {
//     switch (buffer[0]) {
//         case '+': return VALUE_TEXT;
//         case ':': return VALUE_INTEGER;
//         case ',': return VALUE_FLOAT;
//         case '_': return VALUE_NULL;
//         case '$': return VALUE_BLOB;
//     }
//     return VALUE_NULL;
// }
// 
// static bool internal_has_commandlen (int c) {
//     return ((c == CMD_INT) || (c == CMD_FLOAT) || (c == CMD_NULL)) ? false : true;
// }
// 
// static uint32_t internal_parse_number (char *buffer, uint32_t blen, uint32_t *cstart) {
//     uint32_t value = 0;
//     
//     for (uint32_t i=0; i<blen; ++i) {
//         if (buffer[i] == ' ') {
//             *cstart = i+1;
//             return value;
//         }
//         value = (value * 10) + (buffer[i] - '0');
//     }
//     
//     return 0;
// }
// 
// static char *internal_parse_value (char *buffer, uint32_t *len, uint32_t *cellsize) {
//     // handle special NULL value case
//     if (!buffer || buffer[0] == CMD_NULL) {
//         *len = 0;
//         if (cellsize) *cellsize = 2;
//         return NULL;
//     }
//     
//     // blen originally was hard coded to 24 because the max 64bit value is 20 characters long
//     uint32_t cstart = 0;
//     uint32_t blen = *len;
//     blen = internal_parse_number(&buffer[1], blen, &cstart);
//     
//     // handle decimal/float cases
//     if ((buffer[0] == CMD_INT) || (buffer[0] == CMD_FLOAT)) {
//         *len = cstart - 1;
//         if (cellsize) *cellsize = cstart + 1;
//         return &buffer[1];
//     }
//     
//     *len = (buffer[0] == CMD_ZEROSTRING) ? blen - 1 : blen;
//     if (cellsize) *cellsize = cstart + blen + 1;
//     return &buffer[1+cstart];
// }
// 
// static SQCloudResult *internal_run_command (SQCloudConnection *connection, const char *buffer, size_t blen, bool mainfd) {
//     internal_clear_error(connection);
//     
//     if (!buffer || blen < CMD_MINLEN) return NULL;
//     
//     TIME_GET(tstart);
//     if (!internal_socket_write(connection, buffer, blen, mainfd)) return false;
//     SQCloudResult *result = internal_socket_read(connection, mainfd);
//     TIME_GET(tend);
//     if (result) result->time = TIME_VAL(tstart, tend);
//     return result;
// }
// 
// static SQCloudResult *internal_setup_pubsub (SQCloudConnection *connection, const char *buffer, size_t blen) {
//     // check if pubsub was already setup
//     if (connection->pubsubfd != 0) return &SQCloudResultOK;
//     
//     if (internal_connect(connection, connection->hostname, connection->port, NULL, false)) {
//         SQCloudResult *result = internal_run_command(connection, buffer, blen, false);
//         if (!SQCloudResultIsOK(result)) return result;
//         internal_parse_uuid(connection, buffer, blen);
//         pthread_create(&connection->tid, NULL, pubsub_thread, (void *)connection);
//     } else {
//         return NULL;
//     }
//     
//     return &SQCloudResultOK;
// }
// 
// static SQCloudResult *internal_reconnect (SQCloudConnection *connection, const char *buffer, size_t blen) {
//     // DO RE-CONNECT HERE
//     return NULL;
// }
// 
// static SQCloudResult *internal_rowset_type (SQCloudConnection *connection, char *buffer, uint32_t blen, uint32_t bstart, SQCloudResType type) {
//     SQCloudResult *rowset = (SQCloudResult *)mem_zeroalloc(sizeof(SQCloudResult));
//     if (!rowset) {
//         internal_set_error(connection, 1, "Unable to allocate memory for SQCloudResult: %d.", sizeof(SQCloudResult));
//         return NULL;
//     }
//     
//     rowset->tag = type;
//     rowset->buffer = &buffer[bstart];
//     rowset->rawbuffer = buffer;
//     rowset->blen = blen;
//     rowset->balloc = blen;
//     
//     return rowset;
// }
// 
// static SQCloudResult *internal_parse_rowset (SQCloudConnection *connection, char *buffer, uint32_t blen, uint32_t bstart, uint32_t nrows, uint32_t ncols) {    
//     SQCloudResult *rowset = (SQCloudResult *)mem_zeroalloc(sizeof(SQCloudResult));
//     if (!rowset) {
//         internal_set_error(connection, 1, "Unable to allocate memory for SQCloudResult: %d.", sizeof(SQCloudResult));
//         return NULL;
//     }
//     
//     rowset->tag = RESULT_ROWSET;
//     rowset->buffer = buffer;
//     rowset->rawbuffer = buffer;
//     rowset->blen = blen;
//     rowset->balloc = blen;
//     
//     rowset->nrows = nrows;
//     rowset->ncols = ncols;
//     rowset->data = (char **) mem_alloc(nrows * ncols * sizeof(char *));
//     rowset->name = (char **) mem_alloc(ncols * sizeof(char *));
//     rowset->clen = (uint32_t *) mem_zeroalloc(ncols * sizeof(uint32_t));
//     if (!rowset->data || !rowset->name || !rowset->clen) goto abort_rowset;
//     
//     buffer += bstart;
//     
//     // the first column contains names
//     for (uint32_t i=0; i<ncols; ++i) {
//         uint32_t cstart = 0;
//         uint32_t len = internal_parse_number(&buffer[1], blen, &cstart);
//         rowset->name[i] = buffer;
//         buffer += cstart + len + 1;
//         blen -= cstart + len + 1;
//         if (rowset->clen[i] < len) rowset->clen[i] = len;
//         if (rowset->maxlen < len) rowset->maxlen = len;
//     }
//     
//     // parse values
//     for (uint32_t i=0; i<nrows * ncols; ++i) {
//         uint32_t len = blen, cellsize;
//         char *value = internal_parse_value(buffer, &len, &cellsize);
//         rowset->data[i] = (value) ? buffer : NULL;
//         buffer += cellsize;
//         blen -= cellsize;
//         ++rowset->ndata;
//         if (rowset->clen[i % ncols] < len) rowset->clen[i % ncols] = len;
//         if (rowset->maxlen < len) rowset->maxlen = len;
//     }
//     
//     return rowset;
//     
// abort_rowset:
//     if (rowset->data) mem_free(rowset->data);
//     if (rowset->name) mem_free(rowset->name);
//     if (rowset->clen) mem_free(rowset->clen);
//     if (rowset) mem_free(rowset);
//     
//     internal_set_error(connection, 1, "Unable to allocate internal memory for SQCloudResult.");
//     return NULL;
// }
// 
// static SQCloudResult *internal_parse_rowset_chunck (SQCloudConnection *connection, char *buffer, uint32_t blen, uint32_t bstart, uint32_t idx, uint32_t nrows, uint32_t ncols) {
//     SQCloudResult *rowset = connection->_chunk;
//     bool first_chunk = false;
//     
//     // sanity check
//     if (idx == 1 && connection->_chunk) {
//         // something bad happened here because a first chunk is received while a saved one has not been fully processed
//         // lets try to restart the whole process
//         SQCloudResultFree(connection->_chunk);
//         connection->_chunk = NULL;
//         rowset = NULL;
//     }
//     
//     if (!rowset) {
//         // this should never happen
//         if (idx != 1) return NULL;
//         
//         // allocate a new rowset
//         rowset = (SQCloudResult *)mem_zeroalloc(sizeof(SQCloudResult));
//         if (!rowset) {
//             internal_set_error(connection, 1, "Unable to allocate memory for SQCloudResult: %d.", sizeof(SQCloudResult));
//             return NULL;
//         }
//         first_chunk = true;
//         connection->_chunk = rowset;
//     }
//     
//     if (first_chunk) {
//         rowset->tag = RESULT_ROWSET;
//         rowset->ischunk = true;
//         
//         rowset->buffers = (char **)mem_zeroalloc((sizeof(char *) * DEFAULT_CHUCK_NBUFFERS));
//         if (!rowset->buffers) goto abort_rowset;
//         
//         rowset->bnum = DEFAULT_CHUCK_NBUFFERS;
//         rowset->buffers[0] = buffer;
//         rowset->bcount = 1;
//         
//         rowset->brows = nrows + DEFAULT_CHUNK_MINROWS;
//         rowset->nrows = nrows;
//         rowset->ncols = ncols;
//         rowset->data = (char **) mem_alloc(rowset->brows * ncols * sizeof(char *));
//         rowset->name = (char **) mem_alloc(ncols * sizeof(char *));
//         rowset->clen = (uint32_t *) mem_zeroalloc(ncols * sizeof(uint32_t));
//         if (!rowset->data || !rowset->name || !rowset->clen) goto abort_rowset;
//         
//         buffer += bstart;
//         
//         // first buffer is guarantee to contains column names
//         for (uint32_t i=0; i<ncols; ++i) {
//             uint32_t cstart = 0;
//             uint32_t len = internal_parse_number(&buffer[1], blen, &cstart);
//             rowset->name[i] = buffer;
//             buffer += cstart + len + 1;
//             blen -= cstart + len + 1;
//             if (rowset->clen[i] < len) rowset->clen[i] = len;
//             if (rowset->maxlen < len) rowset->maxlen = len;
//         }
//     }
//     
//     // update total buffer size
//     rowset->blen += blen;
//     
//     // check end-chunk condition
//     if (idx == 0 && nrows == 0 && ncols == 0) {
//         connection->_chunk = NULL;
//         if (!rowset->externalbuffer) mem_free(buffer);
//         return rowset;
//     }
//     
//     // check if a resize is needed in the array of buffers
//     if (rowset->bnum <= rowset->bcount + 1) {
//         uint32_t n = rowset->bnum * 2;
//         char **temp = (char **)mem_realloc(rowset->buffers, (sizeof(char *) * n));
//         if (!temp) goto abort_rowset;
//         rowset->buffers = temp;
//         rowset->bnum = n;
//     }
//     
//     // check if a resize is needed in the ptr data array
//     if (rowset->brows <= rowset->nrows + nrows) {
//         uint32_t n = rowset->brows * 2;
//         char **temp = (char **)mem_realloc(rowset->data, n * ncols * (sizeof(char *)));
//         if (!temp) goto abort_rowset;
//         rowset->data = temp;
//         rowset->brows = n;
//     }
//     
//     // adjust internal fields
//     if (!first_chunk) {
//         rowset->buffers[rowset->bcount++] = buffer;
//         rowset->nrows += nrows;
//         buffer += bstart;
//     }
//     
//     // parse values
//     uint32_t index = rowset->ndata;
//     uint32_t bound = rowset->ndata + (nrows * ncols);
//     for (uint32_t i=index; i<bound; ++i) {
//         uint32_t len = blen, cellsize;
//         char *value = internal_parse_value(buffer, &len, &cellsize);
//         rowset->data[i] = (value) ? buffer : NULL;
//         buffer += cellsize;
//         blen -= cellsize;
//         ++rowset->ndata;
//         if (rowset->clen[i % ncols] < len) rowset->clen[i % ncols] = len;
//         if (rowset->maxlen < len) rowset->maxlen = len;
//     }
//     
//     // this check is for internal usage only
//     if (connection->fd == 0) return rowset;
//     
//     // normal usage
//     // send ACK
//     if (!internal_socket_write(connection, "OK", 2, true)) goto abort_rowset;
//         
//     // read next chunk
//     return internal_socket_read (connection, true);
//     
// abort_rowset:
//     SQCloudResultFree(rowset);
//     connection->_chunk = NULL;
//     return NULL;
// }
// 
// static SQCloudResult *internal_parse_buffer (SQCloudConnection *connection, char *buffer, uint32_t blen, uint32_t cstart, bool isstatic, bool externalbuffer) {
//     if (blen <= 1) return false;
//     
//     // try to check if it is a OK reply: +2 OK
//     if ((blen == REPLY_OK_LEN) && (strncmp(buffer, REPLY_OK, REPLY_OK_LEN) == 0)) {
//         return &SQCloudResultOK;
//     }
//     
//     // if buffer is static (stack based allocation) then it must be duplicated
//     if (buffer[0] != CMD_ERROR && isstatic) {
//         char *clone = mem_alloc(blen);
//         if (!clone) {
//             internal_set_error(connection, 1, "Unable to allocate memory: %d.", blen);
//             return NULL;
//         }
//         memcpy(clone, buffer, blen);
//         buffer = clone;
//         isstatic = false;
//     }
//     
//     // check for compressed reply before the parse step
//     char *zdata = NULL;
//     if (buffer[0] == CMD_COMPRESSED) {
//         // %TLEN CLEN ULEN *0 NROWS NCOLS DATA
//         uint32_t cstart1 = 0;
//         uint32_t cstart2 = 0;
//         uint32_t cstart3 = 0;
//         uint32_t tlen = internal_parse_number(&buffer[1], blen-1, &cstart1);
//         uint32_t clen = internal_parse_number(&buffer[cstart1 + 1], blen-1, &cstart2);
//         uint32_t ulen = internal_parse_number(&buffer[cstart1 + cstart2 + 1], blen-1, &cstart3);
//         
//         // start of compressed buffer
//         zdata = &buffer[tlen - clen + cstart1 + 1];
//         
//         // start of raw uncompressed header
//         char *hstart = &buffer[cstart1 + cstart2 + cstart3 + 1];
//         
//         // try to allocate a buffer big enough to hold uncompressed data + raw header
//         long clonelen = ulen + (zdata - hstart) + 1;
//         char *clone = mem_alloc (clonelen);
//         if (!clone) {
//             internal_set_error(connection, 1, "Unable to allocate memory to uncompress buffer: %d.", clonelen);
//             if (!isstatic && !externalbuffer) mem_free(buffer);
//             return NULL;
//         }
//         
//         // copy raw buffer
//         memcpy(clone, hstart, zdata - hstart);
//         
//         // uncompress buffer and sanity check the result
//         uint32_t rc = LZ4_decompress_safe(zdata, clone + (zdata - hstart), clen, ulen);
//         if (rc <= 0 || rc != ulen) {
//             internal_set_error(connection, 1, "Unable to decompress buffer (err code: %d).", rc);
//             if (!isstatic && !externalbuffer) mem_free(buffer);
//             return NULL;
//         }
//         
//         // decompression is OK so replace buffer
//         if (!isstatic && !externalbuffer) mem_free(buffer);
//         
//         isstatic = false;
//         buffer = clone;
//         blen = ulen;
//         
//         // at this point the buffer used in the SQCloudResult is a newly allocated one (clone)
//         // so externalbuffer flag must be set to false
//         externalbuffer = false;
//     }
//     
//     // parse reply
//     switch (buffer[0]) {
//         case CMD_ZEROSTRING:
//         case CMD_RECONNECT:
//         case CMD_PUBSUB:
//         case CMD_COMMAND:
//         case CMD_STRING:
//         case CMD_JSON: {
//             // +LEN string
//             uint32_t cstart = 0;
//             uint32_t len = internal_parse_number(&buffer[1], blen-1, &cstart);
//             SQCloudResType type = (buffer[0] == CMD_JSON) ? RESULT_JSON : RESULT_STRING;
//             if (buffer[0] == CMD_ZEROSTRING) --len;
//             else if (buffer[0] == CMD_COMMAND) return internal_run_command(connection, &buffer[cstart+1], len, true);
//             else if (buffer[0] == CMD_PUBSUB) return internal_setup_pubsub(connection, &buffer[cstart+1], len);
//             else if (buffer[0] == CMD_RECONNECT) return internal_reconnect(connection, &buffer[cstart+1], len);
//             SQCloudResult *res = internal_rowset_type(connection, buffer, len, cstart+1, type);
//             if (res) res->externalbuffer = externalbuffer;
//             return res;
//         }
//             
//         case CMD_ERROR: {
//             // -LEN ERRCODE ERRMSG
//             uint32_t cstart = 0, cstart2 = 0;
//             uint32_t len = internal_parse_number(&buffer[1], blen-1, &cstart);
//             
//             uint32_t errcode = internal_parse_number(&buffer[cstart + 1], blen-1, &cstart2);
//             connection->errcode = (int)errcode;
//             
//             len -= cstart2;
//             memcpy(connection->errmsg, &buffer[cstart + cstart2 + 1], MIN(len, sizeof(connection->errmsg)));
//             connection->errmsg[len] = 0;
//             
//             // check free buffer
//             if (!isstatic && !externalbuffer) mem_free(buffer);
//             return NULL;
//         }
//         
//         case CMD_ROWSET:
//         case CMD_ROWSET_CHUNK: {
//             // CMD_ROWSET:          *LEN ROWS COLS DATA
//             // CMD_ROWSET_CHUNK:    /LEN IDX ROWS COLS DATA
//             uint32_t cstart1 = 0, cstart2 = 0, cstart3 = 0, cstart4 = 0;
//             
//             internal_parse_number(&buffer[1], blen-1, &cstart1); // parse len (already parsed in blen parameter)
//             uint32_t idx = (buffer[0] == CMD_ROWSET) ? 0 : internal_parse_number(&buffer[cstart1 + 1], blen-1, &cstart2);
//             uint32_t nrows = internal_parse_number(&buffer[cstart1 + cstart2 + 1], blen-1, &cstart3);
//             uint32_t ncols = internal_parse_number(&buffer[cstart1 + cstart2 + + cstart3 + 1], blen-1, &cstart4);
//             
//             uint32_t bstart = cstart1 + cstart2 + cstart3 + cstart4 + 1;
//             SQCloudResult *res = NULL;
//             if (buffer[0] == CMD_ROWSET) res = internal_parse_rowset(connection, buffer, blen, bstart, nrows, ncols);
//             else res = internal_parse_rowset_chunck(connection, buffer, blen, bstart, idx, nrows, ncols);
//             if (res) res->externalbuffer = externalbuffer;
//             
//             // check free buffer
//             if (!res && !isstatic && !externalbuffer) mem_free(buffer);
//             return res;
//         }
//         
//         case CMD_NULL:
//             if (!isstatic && !externalbuffer) mem_free(buffer);
//             return &SQCloudResultNULL;
//             
//         case CMD_INT:
//         case CMD_FLOAT: {
//             // NUMBER case
//             internal_parse_value(buffer, &blen, NULL);
//             SQCloudResult *res = internal_rowset_type(connection, buffer, blen, 1, (buffer[0] == CMD_INT) ? RESULT_INTEGER : RESULT_FLOAT);
//             if (res) res->externalbuffer = externalbuffer;
//             
//             if (!res && !isstatic && !externalbuffer) mem_free(buffer);
//             return res;
//         }
//             
//         case CMD_RAWJSON: {
//             // handle JSON here
//             return &SQCloudResultNULL;
//         }
//     }
//     
//     if (!isstatic && !externalbuffer) mem_free(buffer);
//     return NULL;
// }
// 
// static bool internal_socker_forward_read (SQCloudConnection *connection, bool (*forward_cb) (char *buffer, size_t blen, void *xdata), void *xdata) {
//     char sbuffer[8129];
//     uint32_t blen = sizeof(sbuffer);
//     uint32_t cstart = 0;
//     uint32_t tread = 0;
//     uint32_t clen = 0;
//     
//     char *buffer = sbuffer;
//     char *original = buffer;
//     int fd = connection->fd;
//     
//     while (1) {
//         // perform read operation
//         ssize_t nread = readsocket(fd, buffer, blen);
//         
//         // sanity check read
//         if (nread < 0) {
//             internal_set_error(connection, 1, "An error occurred while reading data: %s.", strerror(errno));
//             goto abort_read;
//         }
//         
//         if (nread == 0) {
//             internal_set_error(connection, 1, "Unexpected EOF found while reading data: %s.", strerror(errno));
//             goto abort_read;
//         }
//         
//         // forward read to callback
//         bool result = forward_cb(buffer, nread, xdata);
//         if (!result) goto abort_read;
//         
//         // update internal counter
//         tread += (uint32_t)nread;
//         
//         // determine command length
//         if (clen == 0) {
//             clen = internal_parse_number (&original[1], tread-1, &cstart);
//             
//             // handle special cases
//             if ((original[0] == CMD_INT) || (original[0] == CMD_FLOAT) || (original[0] == CMD_NULL)) clen = 0;
//             else if (clen == 0) continue;
//         }
//         
//         // check if read is complete
//         if (clen + cstart + 1 == tread) break;
//     }
//     
//     return true;
//     
// abort_read:
//     return false;
// }
// 
// static SQCloudResult *internal_socket_read (SQCloudConnection *connection, bool mainfd) {
//     // most of the time one read will be sufficient
//     char header[1024];
//     char *buffer = (char *)&header;
//     uint32_t blen = sizeof(header);
//     uint32_t tread = 0;
// 
//     int fd = (mainfd) ? connection->fd : connection->pubsubfd;
//     char *original = buffer;
//     while (1) {
//         ssize_t nread = readsocket(fd, buffer, blen);
//         
//         if (nread < 0) {
//             internal_set_error(connection, 1, "An error occurred while reading data: %s.", strerror(errno));
//             goto abort_read;
//         }
//         
//         if (nread == 0) {
//             internal_set_error(connection, 1, "Unexpected EOF found while reading data: %s.", strerror(errno));
//             goto abort_read;
//         }
//         
//         tread += (uint32_t)nread;
//         blen -= (uint32_t)nread;
//         buffer += nread;
//         
//         // parse buffer looking for command length
//         uint32_t cstart = 0;
//         uint32_t clen = 0;
//         
//         if (internal_has_commandlen(original[0])) {
//             clen = internal_parse_number (&original[1], tread-1, &cstart);
//             if (clen == 0) continue;
//         }
//         
//         // check if read is complete
//         // clen is the lenght parsed in the buffer
//         // cstart is the index of the first space
//         // +1 because we skipped the first character in the internal_parse_number function
//         if (clen + cstart + 1 != tread) {
//             // check buffer allocation and continue reading
//             if (clen + cstart > blen) {
//                 char *clone = mem_alloc(clen + cstart + 1);
//                 if (!clone) {
//                     internal_set_error(connection, 1, "Unable to allocate memory: %d.", clen + cstart + 1);
//                     goto abort_read;
//                 }
//                 memcpy(clone, original, tread);
//                 buffer = original = clone;
//                 blen = (clen + cstart + 1) - tread;
//                 buffer += tread;
//             }
//             
//             continue;
//         }
//         
//         // command is complete so parse it
//         return internal_parse_buffer(connection, original, tread, (clen) ? cstart : 0, (original == header), false);
//     }
//     
// abort_read:
//     if (original != (char *)&header) mem_free(original);
//     return NULL;
// }
// 
// static bool internal_socket_write (SQCloudConnection *connection, const char *buffer, size_t len, bool mainfd) {
//     size_t written = 0;
//     
//     int fd = (mainfd) ? connection->fd : connection->pubsubfd;
//     
//     // write header
//     char header[32];
//     int hlen = snprintf(header, sizeof(header), "+%zu ", len);
//     ssize_t n = writesocket(fd, header, hlen);
//     if (n != hlen) return internal_set_error(connection, 1, "An error occurred while writing header data: %s.", strerror(errno));
//     
//     // write buffer
//     while (written < len) {
//         ssize_t nwrote = writesocket(fd, buffer, len);
//         //printf("writesocket connfd:%d nwrote:%d", fd, nwrote);
//         
//         if (nwrote < 0) {
//             return internal_set_error(connection, 1, "An error occurred while writing data: %s.", strerror(errno));
//         } else if (nwrote == 0) {
//             return true;
//         } else {
//             written += nwrote;
//             buffer += nwrote;
//             len -= nwrote;
//         }
//     }
//     
//     return true;
// }
// 
// static void internal_socket_set_timeout (int sockfd, int timeout_secs) {
//     #ifdef _WIN32
//     DWORD timeout = timeout_secs * 1000;
//     setsockopt(sockfd, SOL_SOCKET, SO_RCVTIMEO, (const char*)&timeout, sizeof timeout);
//     setsockopt(sockfd, SOL_SOCKET, SO_SNDTIMEO, (const char*)&timeout, sizeof timeout);
//     #else
//     struct timeval tv;
//     tv.tv_sec = timeout_secs;
//     tv.tv_usec = 0;
//     setsockopt(sockfd, SOL_SOCKET, SO_RCVTIMEO, (const char*)&tv, sizeof tv);
//     setsockopt(sockfd, SOL_SOCKET, SO_SNDTIMEO, (const char*)&tv, sizeof tv);
//     #endif
// }
// 
// static bool internal_connect_apply_config (SQCloudConnection *connection, SQCloudConfig *config) {
//     if (config->timeout) {
//         internal_socket_set_timeout(connection->fd, config->timeout);
//     }
//     
//     if (config->username && config->password) {
//         char buffer[1024];
//         snprintf(buffer, sizeof(buffer), "AUTH USER %s PASS %s", config->username, config->password);
//         SQCloudResult *res = internal_run_command(connection, buffer, strlen(buffer), true);
//         if (res != &SQCloudResultOK) return false;
//     }
//     
//     if (config->database) {
//         char buffer[1024];
//         snprintf(buffer, sizeof(buffer), "USE DATABASE %s", config->database);
//         SQCloudResult *res = internal_run_command(connection, buffer, strlen(buffer), true);
//         if (res != &SQCloudResultOK) return false;
//     }
//     
//     return true;
// }
// 
// static bool internal_connect (SQCloudConnection *connection, const char *hostname, int port, SQCloudConfig *config, bool mainfd) {
//     // ipv4/ipv6 specific variables
//     struct addrinfo hints, *addr_list = NULL, *addr;
//     
//     // ipv6 code from https://www.ibm.com/support/knowledgecenter/ssw_ibm_i_72/rzab6/xip6client.htm
//     memset(&hints, 0, sizeof(hints));
//     hints.ai_family = (config) ? config->family : AF_UNSPEC;
//     hints.ai_socktype = SOCK_STREAM;
//     
//     // get the address information for the server using getaddrinfo()
//     char port_string[256];
//     snprintf(port_string, sizeof(port_string), "%d", port);
//     int rc = getaddrinfo(hostname, port_string, &hints, &addr_list);
//     if (rc != 0 || addr_list == NULL) {
//         return internal_set_error(connection, 1, "Error while resolving getaddrinfo (host %s not found).", hostname);
//     }
//     
//     // begin non-blocking connection loop
//     int sock_index = 0;
//     int sock_current = 0;
//     int sock_list[MAX_SOCK_LIST] = {0};
//     for (addr = addr_list; addr != NULL; addr = addr->ai_next, ++sock_index) {
//         if (sock_index >= MAX_SOCK_LIST) break;
//         if ((addr->ai_family != AF_INET) && (addr->ai_family != AF_INET6)) continue;
//         
//         sock_current = socket(addr->ai_family, addr->ai_socktype, addr->ai_protocol);
//         if (sock_current < 0) continue;
//         
//         // set socket options
//         int len = 1;
//         setsockopt(sock_current, SOL_SOCKET, SO_KEEPALIVE, (const char *) &len, sizeof(len));
//         len = 1;
//         setsockopt(sock_current, IPPROTO_TCP, TCP_NODELAY, (const char *) &len, sizeof(len));
//         #ifdef SO_NOSIGPIPE
//         len = 1;
//         setsockopt(sock_current, SOL_SOCKET, SO_NOSIGPIPE, (const char *) &len, sizeof(len));
//         #endif
//         
//         // by default, an IPv6 socket created on Windows Vista and later only operates over the IPv6 protocol
//         // in order to make an IPv6 socket into a dual-stack socket, the setsockopt function must be called
//         if (addr->ai_family == AF_INET6) {
//             #ifdef _WIN32
//             DWORD ipv6only = 0;
//             #else
//             int   ipv6only = 0;
//             #endif
//             setsockopt(sock_current, IPPROTO_IPV6, IPV6_V6ONLY, (void *)&ipv6only, sizeof(ipv6only));
//         }
//         
//         // turn on non-blocking
//         unsigned long ioctl_blocking = 1;    /* ~0; //TRUE; */
//         ioctl(sock_current, FIONBIO, &ioctl_blocking);
//         
//         // initiate non-blocking connect ignoring return code
//         connect(sock_current, addr->ai_addr, addr->ai_addrlen);
//         
//         // add sock_current to internal list of trying to connect sockets
//         sock_list[sock_index] = sock_current;
//     }
//     
//     // free not more needed memory
//     freeaddrinfo(addr_list);
//     
//     // calculate the connection timeout and reset timers
//     int connect_timeout = (config && config->timeout > 0) ? config->timeout : SQCLOUD_DEFAULT_TIMEOUT;
//     time_t start = time(NULL);
//     time_t now = start;
//     rc = 0;
//     
//     int sockfd = 0;
//     fd_set write_fds;
//     fd_set except_fds;
//     struct timeval tv;
//     
//     while (rc == 0 && ((now - start) < connect_timeout)) {
//         FD_ZERO(&write_fds);
//         FD_ZERO(&except_fds);
//         
//         int nfds = 0;
//         for (int i=0; i<MAX_SOCK_LIST; ++i) {
//             if (sock_list[i]) {
//                 FD_SET(sock_list[i], &write_fds);
//                 FD_SET(sock_list[i], &except_fds);
//                 if (nfds < sock_list[i]) nfds = sock_list[i];
//             }
//         }
//         
//         tv.tv_sec = connect_timeout;
//         tv.tv_usec = 0;
//         rc = select(nfds + 1, NULL, &write_fds, &except_fds, &tv);
//         
//         if (rc == 0) break; // timeout
//         else if (rc == -1) {
//             if (errno == EINTR || errno == EAGAIN || errno == EINPROGRESS) continue;
//             break; // handle error
//         }
//         
//         // check for error first
//         for (int i=0; i<MAX_SOCK_LIST; ++i) {
//             if (sock_list[i] > 0) {
//                 if (FD_ISSET(sock_list[i], &except_fds)) {
//                     closesocket(sock_list[i]);
//                     sock_list[i] = 0;
//                 }
//             }
//         }
//         
//         // check which file descriptor is ready (need to check for socket error also)
//         for (int i=0; i<MAX_SOCK_LIST; ++i) {
//             if (sock_list[i] > 0) {
//                 if (FD_ISSET(sock_list[i], &write_fds)) {
//                     int err = socket_geterror(sock_list[i]);
//                     if (err > 0) {
//                         closesocket(sock_list[i]);
//                         sock_list[i] = 0;
//                     } else {
//                         sockfd = sock_list[i];
//                         break;
//                     }
//                 }
//             }
//         }
//         // check if a valid descriptor has been found
//         if (sockfd != 0) break;
//         
//         // no socket ready yet
//         now = time(NULL);
//         rc = 0;
//     }
//     
//     // close still opened sockets
//     for (int i=0; i<MAX_SOCK_LIST; ++i) {
//         if ((sock_list[i] > 0) && (sock_list[i] != sockfd)) closesocket(sock_list[i]);
//     }
//     
//     // bail if there was an error
//     if (rc < 0) {
//         return internal_set_error(connection, 1, "An error occurred while trying to connect: %s.", strerror(errno));
//     }
//     
//     // bail if there was a timeout
//     if ((time(NULL) - start) >= connect_timeout) {
//         return internal_set_error(connection, 1, "Connection timeout while trying to connect (%d).", connect_timeout);
//     }
//     
//     // turn off non-blocking
//     int ioctl_blocking = 0;    /* ~0; //TRUE; */
//     ioctl(sockfd, FIONBIO, &ioctl_blocking);
//     
//     // SSL on sockfd
//     if (!internal_setup_ssl(connection, config)) return false;
//     
//     // finalize connection
//     if (mainfd) {
//         connection->fd = sockfd;
//         connection->port = port;
//         connection->hostname = strdup(hostname);
//     } else {
//         connection->pubsubfd = sockfd;
//     }
//     return true;
// }
// 
// // MARK: - URL -
// 
// static int url_extract_username_password (const char *s, char b1[512], char b2[512]) {
//     // user:pass@host.com:port/dbname?timeout=10&key2=value2&key3=value3
//     
//     // lookup username (if any)
//     char *username = strchr(s, ':');
//     if (!username) return 0;
//     size_t len = username - s;
//     if (len > 511) return -1;
//     memcpy(b1, s, len);
//     b1[len] = 0;
//     
//     // lookup username (if any)
//     char *password = strchr(s, '@');
//     if (!password) return 0;
//     len = password - username - 1;
//     if (len > 511) return -1;
//     memcpy(b2, username+1, len);
//     b2[len] = 0;
//     
//     return (int)(password - s) + 1;
// }
// 
// static int url_extract_hostname_port (const char *s, char b1[512], char b2[512]) {
//     // host.com:port/dbname?timeout=10&key2=value2&key3=value3
//     
//     // lookup hostname (if any)
//     char *hostname = strchr(s, ':');
//     if (!hostname) hostname = strchr(s, '/');
//     if (!hostname) hostname = strchr(s, '?');
//     if (!hostname) hostname = strchr(s, 0);
//     if (!hostname) return -1;
//     size_t len = hostname - s;
//     if (len > 511) return -1;
//     memcpy(b1, s, len);
//     b1[len] = 0;
//     
//     // lookup port (if any)
//     char *port = strchr(s, ':');
//     if (port) {
//         char *p = port + 1;
//         ++len;
//         
//         int i = 0;
//         while (p[0]) {
//             if ((p[0] == '/') || (p[0] == '?') || (p[0] == 0)) break;
//             if (i+1 > 511) return -1;
//             b2[i++] = p[0];
//             ++len;
//             ++p;
//         }
//         b2[len] = 0;
//     }
//     
//     // adjust returned len
//     if (s[len] != 0) ++len;
//     
//     return (int)len;
// }
// 
// static int url_extract_database (const char *s, char b1[512]) {
//     // dbname?timeout=10&key2=value2&key3=value3
//     
//     // lookup database (if any)
//     char *database = strchr(s, '?');
//     if (database) {
//         size_t len = database - s;
//         if (len > 511) return -1;
//         memcpy(b1, s, len);
//         b1[len] = 0;
//         return (int)(len + 1);
//     }
//     
//     // there is no ? separator character
//     // that means that there should be
//     // no key/value
//     char *guard = strchr(s, '=');
//     if (guard) return 0;
//     
//     // database name is the s string
//     size_t len = strlen(s);
//     if (len > 511) return -1;
//     memcpy(b1, s, len);
//     b1[len] = 0;
//     
//     return (int)len;
// }
// 
// static int url_extract_keyvalue (const char *s, char b1[512], char b2[512]) {
//     // timeout=10&key2=value2&key3=value3
//     
//     // lookup key (if any)
//     char *key = strchr(s, '=');
//     if (!key) return 0;
//     size_t len = key - s;
//     if (len > 511) return -1;
//     memcpy(b1, s, len);
//     b1[len] = 0;
//     
//     // lookup username (if any)
//     char *value = strchr(s, '&');
//     if (!value) value = strchr(s, 0);
//     if (!value) return 0;
//     len = value - key - 1;
//     if (len > 511) return -1;
//     memcpy(b2, key+1, len);
//     b2[len] = 0;
//     
//     return (int)(value - s) + 1;
// }
// 
// // MARK: - RESERVED -
// 
// bool _reserved1 (SQCloudConnection *connection, const char *command, bool (*forward_cb) (char *buffer, size_t blen, void *xdata), void *xdata) {
//     if (!forward_cb) return false;
//     if (!internal_socket_write(connection, command, strlen(command), true)) return false;
//     if (!internal_socker_forward_read(connection, forward_cb, xdata)) return false;
//     return true;
// }
// 
// SQCloudResult *_reserved2 (SQCloudConnection *connection, const char *UUID) {
//     char command[512];
//     snprintf(command, sizeof(command), "SET CLIENT UUID TO %s", UUID);
//     return internal_run_command(connection, command, strlen(command), true);
// }
// 
// SQCloudResult *_reserved3 (char *buffer, uint32_t blen, uint32_t cstart, SQCloudResult *chunk) {
//     SQCloudConnection connection = {0};
//     connection._chunk = chunk;
//     SQCloudResult *res = internal_parse_buffer(&connection, buffer, blen, cstart, false, true);
//     return res;
// }
// 
// uint32_t _reserved4 (char *buffer, uint32_t blen, uint32_t *cstart) {
//     return internal_parse_number(buffer, blen, cstart);
// }
// 
// // MARK: - PUBLIC -
// 
// SQCloudConnection *SQCloudConnect (const char *hostname, int port, SQCloudConfig *config) {
//     internal_init();
//     
//     SQCloudConnection *connection = mem_zeroalloc(sizeof(SQCloudConnection));
//     if (!connection) return NULL;
//     
//     if (internal_connect(connection, hostname, port, config, true)) {
//         if (config) internal_connect_apply_config(connection, config);
//     }
//     
//     return connection;
// }
// 
// SQCloudConnection *SQCloudConnectWithString (const char *s) {
//     // URL STRING FORMAT
//     // sqlitecloud://user:pass@host.com:port/dbname?timeout=10&key2=value2&key3=value3
//     
//     // sanity check
//     const char domain[] = "sqlitecloud://";
//     int n = sizeof(domain) - 1;
//     if (strncmp(s, domain, n) != 0) return NULL;
//     
//     // config struct
//     SQCloudConfig *config = NULL;
//     SQCloudConfig sconfig;
//     
//     // lookup for optional username/password
//     char username[512];
//     char password[512];
//     int rc = url_extract_username_password(&s[n], username, password);
//     if (rc == -1) return NULL;
//     if (rc) {
//         sconfig.username = string_dup(username);
//         sconfig.password = string_dup(password);
//         config = &sconfig;
//     }
//     
//     // lookup for mandatory hostname
//     n += rc;
//     char hostname[512];
//     char port_s[512];
//     rc = url_extract_hostname_port(&s[n], hostname, port_s);
//     if (rc <= 0) return NULL;
//     int port = (int)strtol(port_s, NULL, 0);
//     if (port == 0) port = SQCLOUD_DEFAULT_PORT;
//     
//     // lookup for optional database
//     n += rc;
//     char database[512];
//     rc = url_extract_database(&s[n], database);
//     if (rc == -1) return NULL;
//     if (rc > 0) {
//         sconfig.database = string_dup(database);
//         config = &sconfig;
//     }
//     
//     // lookup for optional key(s)/value(s)
//     n += rc;
//     char key[512];
//     char value[512];
//     while ((rc = url_extract_keyvalue(&s[n], key, value)) > 0) {
//         if (strcasecmp(key, "timeout")) {
//             int timeout = (int)strtol(value, NULL, 0);
//             sconfig.timeout = (timeout) ? timeout : 0;
//             config = &sconfig;
//         }
//         n += rc;
//     }
//     
//     return SQCloudConnect(hostname, port, config);
// }
// 
// SQCloudResult *SQCloudExec (SQCloudConnection *connection, const char *command) {
//     return internal_run_command(connection, command, strlen(command), true);
// }
// 
// void SQCloudDisconnect (SQCloudConnection *connection) {
//     if (!connection) return;
//     
//     // free SSL
//     
//     // try to gracefully close connections
//     if (connection->fd) {
//         closesocket(connection->fd);
//     }
//     
//     if (connection->pubsubfd) {
//         closesocket(connection->pubsubfd);
//     }
//     
//     // free memory
//     if (connection->hostname) {
//         mem_free(connection->hostname);
//     }
//     
//     if (connection->uuid) {
//         mem_free(connection->uuid);
//     }
//     
//     mem_free(connection);
// }
// 
// void SQCloudSetPubSubCallback (SQCloudConnection *connection, SQCloudPubSubCB callback, void *data) {
//     connection->callback = callback;
//     connection->data = data;
// }
// 
// SQCloudResult *SQCloudSetPubSubOnly (SQCloudConnection *connection) {
//     if (!connection->callback) {
//         internal_set_error(connection, 1, "A PubSub callback must be set before executing a PUBSUB ONLY command.");
//         return NULL;
//     }
//     
//     const char *command = "PUBSUB ONLY";
//     return internal_run_command(connection, command, strlen(command), true);
// }
// 
// char *SQCloudUUID (SQCloudConnection *connection) {
//     return connection->uuid;
// }
// 
// // MARK: -
// 
// bool SQCloudIsError (SQCloudConnection *connection) {
//     return (!connection || connection->errcode);
// }
// 
// int SQCloudErrorCode (SQCloudConnection *connection) {
//     return (connection) ? connection->errcode : 666;
// }
// 
// const char *SQCloudErrorMsg (SQCloudConnection *connection) {
//     return (connection) ? connection->errmsg : "Not enoght memory to allocate a SQCloudConnection.";
// }
// 
// // MARK: -
// 
// SQCloudResType SQCloudResultType (SQCloudResult *result) {
//     return (result) ? result->tag : RESULT_ERROR;
// }
// 
// bool SQCloudResultIsOK (SQCloudResult *result) {
//     return (result == &SQCloudResultOK);
// }
// 
// uint32_t SQCloudResultLen (SQCloudResult *result) {
//     return (result) ? result->blen : 0;
// }
// 
// char *SQCloudResultBuffer (SQCloudResult *result) {
//     return (result) ? result->buffer : NULL;
// }
// 
// void SQCloudResultFree (SQCloudResult *result) {
//     if (!result || (result == &SQCloudResultOK) || (result == &SQCloudResultNULL)) return;
//     
//     if (!result->ischunk && !result->externalbuffer) {
//         mem_free(result->rawbuffer);
//     }
//     
//     if (result->tag == RESULT_ROWSET) {
//         mem_free(result->name);
//         mem_free(result->data);
//         mem_free(result->clen);
//         
//         if (result->ischunk && !result->externalbuffer) {
//             for (uint32_t i = 0; i<result->bcount; ++i) {
//                 if (result->buffers[i]) mem_free(result->buffers[i]);
//             }
//             mem_free(result->buffers);
//         }
//     }
//     
//     mem_free(result);
// }
// 
// // MARK: -
// 
// // https://database.guide/2-sample-databases-sqlite/
// // https://embeddedartistry.com/blog/2017/07/05/printf-a-limited-number-of-characters-from-a-string/
// // https://stackoverflow.com/questions/1809399/how-to-format-strings-using-printf-to-get-equal-length-in-the-output
// 
// // SET DATABASE mediastore.sqlite
// // SELECT * FROM Artist LIMIT 10;
// 
// static bool SQCloudRowsetSanityCheck (SQCloudResult *result, uint32_t row, uint32_t col) {
//     if (!result || result->tag != RESULT_ROWSET) return false;
//     if ((row >= result->nrows) || (col >= result->ncols)) return false;
//     return true;
// }
// 
// SQCloudValueType SQCloudRowsetValueType (SQCloudResult *result, uint32_t row, uint32_t col) {
//     if (!SQCloudRowsetSanityCheck(result, row, col)) return VALUE_NULL;
//     return internal_type(result->data[row*result->ncols+col]);
// }
// 
// uint32_t SQCloudRowsetRowsMaxColumnLength (SQCloudResult *result, uint32_t col) {
//     return (result) ? result->clen[ col ] : 0;
// }
// 
// char *SQCloudRowsetColumnName (SQCloudResult *result, uint32_t col, uint32_t *len) {
//     if (!result || result->tag != RESULT_ROWSET) return NULL;
//     if (col >= result->ncols) return NULL;
//     *len = result->blen - (uint32_t)(result->name[col] - result->rawbuffer);
//     return internal_parse_value(result->name[col], len, NULL);
// }
// 
// uint32_t SQCloudRowsetRows (SQCloudResult *result) {
//     if (!SQCloudRowsetSanityCheck(result, 0, 0)) return 0;
//     return result->nrows;
// }
// 
// uint32_t SQCloudRowsetCols (SQCloudResult *result) {
//     if (!SQCloudRowsetSanityCheck(result, 0, 0)) return 0;
//     return result->ncols;
// }
// 
// uint32_t SQCloudRowsetMaxLen (SQCloudResult *result) {
//     if (!SQCloudRowsetSanityCheck(result, 0, 0)) return 0;
//     return result->maxlen;
// }
// 
// char *SQCloudRowsetValue (SQCloudResult *result, uint32_t row, uint32_t col, uint32_t *len) {
//     if (!SQCloudRowsetSanityCheck(result, row, col)) return NULL;
//     
//     // The *len var must contain the remaining length of the buffer pointed by
//     // result->data[row*result->ncols+col]. The caller should not be aware of the
//     // internal implementation of this buffer, so it must be set here.
//     *len = result->blen - (uint32_t)(result->data[row*result->ncols+col] - result->rawbuffer);
//     return internal_parse_value(result->data[row*result->ncols+col], len, NULL);
// }
// 
// int32_t SQCloudRowsetInt32Value (SQCloudResult *result, uint32_t row, uint32_t col) {
//     if (!SQCloudRowsetSanityCheck(result, row, col)) return 0;
//     uint32_t len = 0;
//     char *value = internal_parse_value(result->data[row*result->ncols+col], &len, NULL);
//     
//     char buffer[256];
//     snprintf(buffer, sizeof(buffer), "%.*s", len, value);
//     return (int32_t)strtol(buffer, NULL, 0);
// }
// 
// int64_t SQCloudRowsetInt64Value (SQCloudResult *result, uint32_t row, uint32_t col) {
//     if (!SQCloudRowsetSanityCheck(result, row, col)) return 0;
//     uint32_t len = 0;
//     char *value = internal_parse_value(result->data[row*result->ncols+col], &len, NULL);
//     
//     char buffer[256];
//     snprintf(buffer, sizeof(buffer), "%.*s", len, value);
//     return (int64_t)strtoll(buffer, NULL, 0);
// }
// 
// float SQCloudRowsetFloatValue (SQCloudResult *result, uint32_t row, uint32_t col) {
//     if (!SQCloudRowsetSanityCheck(result, row, col)) return 0.0;
//     uint32_t len = 0;
//     char *value = internal_parse_value(result->data[row*result->ncols+col], &len, NULL);
//     
//     char buffer[256];
//     snprintf(buffer, sizeof(buffer), "%.*s", len, value);
//     return (float)strtof(buffer, NULL);
// }
// 
// double SQCloudRowsetDoubleValue (SQCloudResult *result, uint32_t row, uint32_t col) {
//     if (!SQCloudRowsetSanityCheck(result, row, col)) return 0.0;
//     uint32_t len = 0;
//     char *value = internal_parse_value(result->data[row*result->ncols+col], &len, NULL);
//     
//     char buffer[256];
//     snprintf(buffer, sizeof(buffer), "%.*s", len, value);
//     return (double)strtod(buffer, NULL);
// }
// 
// void SQCloudRowsetDump (SQCloudResult *result, uint32_t maxline) {
//     uint32_t nrows = result->nrows;
//     uint32_t ncols = result->ncols;
//     uint32_t blen = result->blen;
//     
//     // if user specify a maxline then do not print more than maxline characters for every column
//     if (maxline > 0) {
//         for (uint32_t i=0; i<ncols; ++i) {
//             if (result->clen[i] > maxline) result->clen[i] = maxline;
//         }
//     }
//     
//     // print separator header
//     for (uint32_t i=0; i<ncols; ++i) {
//         for (uint32_t j=0; j<result->clen[i]+2; ++j) putchar('-');
//         putchar('|');
//     }
//     printf("\n");
//     
//     // print column names
//     for (uint32_t i=0; i<ncols; ++i) {
//         uint32_t len = blen;
//         uint32_t delta = 0;
//         char *value = internal_parse_value(result->name[i], &len, NULL);
//         
//         // UTF-8 strings need special adjustments
//         uint32_t utf8len = utf8_len(value, len);
//         if (utf8len != len) delta = len - utf8len;
//         printf(" %-*.*s |", result->clen[i] + delta, (maxline && len > maxline) ? maxline : len, value);
//         blen -= len;
//     }
//     printf("\n");
//     
//     // print separator header
//     for (uint32_t i=0; i<ncols; ++i) {
//         for (uint32_t j=0; j<result->clen[i]+2; ++j) putchar('-');
//         putchar('|');
//     }
//     printf("\n");
//     
//     // print result
//     for (uint32_t i=0; i<nrows * ncols; ++i) {
//         uint32_t len = blen;
//         uint32_t delta = 0;
//         char *value = internal_parse_value(result->data[i], &len, NULL);
//         blen -= len;
//         
//         // UTF-8 strings need special adjustments
//         if (!value) {value = "NULL"; len = 4;}
//         uint32_t utf8len = utf8_len(value, len);
//         if (utf8len != len) delta = len - utf8len;
//         printf(" %-*.*s |", (result->clen[i % ncols]) + delta, (maxline && len > maxline) ? maxline : len, value);
//         
//         bool newline = (((i+1) % ncols == 0) || (ncols == 1));
//         if (newline) printf("\n");
//     }
//     
//     // print footer
//     for (uint32_t i=0; i<ncols; ++i) {
//         for (uint32_t j=0; j<result->clen[i]+2; ++j) putchar('-');
//         putchar('|');
//     }
//     printf("\n");
//     
//     printf("Rows: %d - Cols: %d - Bytes: %d Time: %f secs", result->nrows, result->ncols, result->blen, result->time);
// }
import "C"
import "unsafe"

// SQCloudConnection *SQCloudConnect (const char *hostname, int port, SQCloudConfig *config);
func CConnect( Host string, Port int, Username string, Password string, Database string, Timeout int, Family int ) *C.struct_SQCloudConnection {
  conf := C.struct_SQCloudConfigStruct{}
  conf.username = C.CString( Username )
  conf.password = C.CString( Password )
  conf.database = C.CString( Database )
  conf.timeout  = C.int( Timeout )
  conf.family   = C.int( Family )

  cHost := C.CString( Host )

  cConnection := C.SQCloudConnect( cHost, C.int( Port ), &conf )
  
  C.free( unsafe.Pointer( cHost ) )
  C.free( unsafe.Pointer( conf.database ) )
  C.free( unsafe.Pointer( conf.password ) )
  C.free( unsafe.Pointer( conf.username ) )

  return cConnection
}

// SQCloudConnection *SQCloudConnectWithString (const char *s);
func CConnectWithString( ConnectionString string ) *SQCloud {
  cConString := C.CString( ConnectionString )
  connection := SQCloud{ connection: C.SQCloudConnectWithString( cConString ) }
  C.free( unsafe.Pointer( cConString ) )

  if connection.connection == nil {
    return nil
  }

  return &connection
}

// void SQCloudDisconnect (SQCloudConnection *connection);
func (this *SQCloud ) CDisconnect() {
  if this.connection != nil {
    C.SQCloudDisconnect( this.connection )
    this.connection = nil
  }
}
// char *SQCloudUUID (SQCloudConnection *connection);
func (this *SQCloud ) CGetCloudUUID() string {
   return C.GoString( C.SQCloudUUID( this.connection ) )
}

//bool SQCloudIsError (SQCloudConnection *connection);
func (this *SQCloud ) CIsError() bool {
  return bool( C.SQCloudIsError( this.connection ) )
}
//int SQCloudErrorCode (SQCloudConnection *connection);
func (this *SQCloud ) CGetErrorCode() int {
  return int( C.SQCloudErrorCode( this.connection ) )
}
//const char *SQCloudErrorMsg (SQCloudConnection *connection);
func (this *SQCloud ) CGetErrorMessage() string {
  return C.GoString( C.SQCloudErrorMsg( this.connection ) )
}

// SQCloudResult *SQCloudExec (SQCloudConnection *connection, const char *command);
func (this *SQCloud ) CExec( Command string ) *SQCloudResult {
  cCommand := C.CString( Command )
  defer C.free( unsafe.Pointer( cCommand ) )

  // println( "exec ("+Command+").." )

  result := SQCloudResult{ result: C.SQCloudExec( this.connection, cCommand ) }
  if result.result == nil {
    return nil
  }
  return &result
}
// SQCloudResult *SQCloudSetPubSubOnly (SQCloudConnection *connection);
func (this *SQCloud ) CSetPubSubOnly() *SQCloudResult {
  result := SQCloudResult{ result: C.SQCloudSetPubSubOnly( this.connection ) }
  
  if result.result == nil {
    return nil
  }

  return &result
}
// SQCloudResType SQCloudResultType (SQCloudResult *result);
func (this *SQCloudResult ) CGetResultType() uint {
  return uint( C.SQCloudResultType( this.result ) )
}
// uint32_t SQCloudResultLen (SQCloudResult *result);
func (this *SQCloudResult ) CGetResultLen() uint {
  return uint( C.SQCloudResultLen( this.result ) )
}
// char *SQCloudResultBuffer (SQCloudResult *result);
func (this *SQCloudResult ) CGetResultBuffer() string {
  return C.GoString( C.SQCloudResultBuffer( this.result ) )
}
// void SQCloudResultFree (SQCloudResult *result);
func (this *SQCloudResult ) CFree() {
  C.SQCloudResultFree( this.result )
}
// bool SQCloudResultIsOK (SQCloudResult *result);
func (this *SQCloudResult ) CIsOK() bool {
  return bool( C.SQCloudResultIsOK( this.result ) )
}
// SQCloudValueType SQCloudRowsetValueType (SQCloudResult *result, uint32_t row, uint32_t col);
func (this *SQCloudResult ) CGetValueType( Row uint, Column uint ) int {
  return int( C.SQCloudRowsetValueType( this.result, C.uint( Row ), C.uint( Column ) ) )
}
// uint32_t SQCloudResultMaxColumnLenght (SQCloudResult *result, uint32_t col) ;
func (this *SQCloudResult ) CGetMaxColumnLenght( Column uint ) uint {
  return uint( C.SQCloudRowsetRowsMaxColumnLength( this.result, C.uint( Column ) ) )
}
// char *SQCloudRowsetColumnName (SQCloudResult *result, uint32_t col, uint32_t *len);
func (this *SQCloudResult ) CGetColumnName( Column uint ) string {
  var len C.uint32_t = 0
  return C.GoStringN( C.SQCloudRowsetColumnName( this.result, C.uint( Column ), &len ), C.int( len ) )
}
// uint32_t SQCloudRowsetRows (SQCloudResult *result);
func (this *SQCloudResult ) CGetRows() uint {
  return uint( C.SQCloudRowsetRows( this.result ) )
}
// uint32_t SQCloudRowsetCols (SQCloudResult *result);
func (this *SQCloudResult ) CGetColumns() uint {
  return uint( C.SQCloudRowsetCols( this.result ) )
}
// uint32_t SQCloudRowsetMaxLen (SQCloudResult *result);
func (this *SQCloudResult ) CGetMaxLen() uint32 {
  return uint32( C.SQCloudRowsetMaxLen( this.result ) )
}
// char *SQCloudRowsetValue (SQCloudResult *result, uint32_t row, uint32_t col, uint32_t *len);
func (this *SQCloudResult ) CGetStringValue( Row uint, Column uint ) string {
  var len C.uint32_t = 0
  return C.GoStringN( C.SQCloudRowsetValue( this.result, C.uint32_t( Row ), C.uint32_t( Column ), &len ), C.int( len ) ) // Problem: NULL Pointer in return
}
// int32_t SQCloudRowsetInt32Value (SQCloudResult *result, uint32_t row, uint32_t col);
func (this *SQCloudResult ) CGetInt32Value( Row uint, Column uint ) int32 {
  return int32( C.SQCloudRowsetInt32Value( this.result, C.uint( Row ), C.uint( Column ) ) )
}
// int64_t SQCloudRowsetInt64Value (SQCloudResult *result, uint32_t row, uint32_t col);
func (this *SQCloudResult ) CGetInt64Value( Row uint, Column uint ) int64 {
  return int64( C.SQCloudRowsetInt64Value( this.result, C.uint( Row ), C.uint( Column ) ) )
}
// float SQCloudRowsetFloatValue (SQCloudResult *result, uint32_t row, uint32_t col);
func (this *SQCloudResult ) CGetFloat32Value( Row uint, Column uint ) float32 {
  return float32( C.SQCloudRowsetFloatValue( this.result, C.uint( Row ), C.uint( Column ) ) )
}
// double SQCloudRowsetDoubleValue (SQCloudResult *result, uint32_t row, uint32_t col);
func (this *SQCloudResult ) CGetFloat64Value( Row uint, Column uint ) float64 {
  return float64( C.SQCloudRowsetDoubleValue( this.result, C.uint( Row ), C.uint( Column ) ) )
}
// void SQCloudRowsetDump (SQCloudResult *result, uint32_t maxline);
func (this *SQCloudResult ) CDump( MaxLine uint ) {
   C.SQCloudRowsetDump( this.result, C.uint( MaxLine ) )
}

// Reserverd (internal) functions - will never be exported

// bool SQCloudForwardExec(SQCloudConnection *connection, const char *command, bool (*forward_cb) (char *buffer, size_t blen, void *xdata), void *xdata) {
// SQCloudResult *SQCloudSetUUID (SQCloudConnection *connection, const char *UUID) 

// Will be implemented in GO

// void SQCloudSetPubSubCallback (SQCloudConnection *connection, SQCloudPubSubCB callback, void *data);