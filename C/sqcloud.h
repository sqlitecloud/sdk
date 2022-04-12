//
//  sqcloud.h
//
//  Created by Marco Bambini on 08/02/21.
//

#ifndef __SQCLOUD_CLI__
#define __SQCLOUD_CLI__

#include <stdio.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

#define SQCLOUD_SDK_VERSION         "0.6.0"
#define SQCLOUD_SDK_VERSION_NUM     0x000600
#define SQCLOUD_DEFAULT_PORT        8860
#define SQCLOUD_DEFAULT_TIMEOUT     12
#define SQCLOUD_DEFAULT_UPLOAD_SIZE 512*1024

#define SQCLOUD_IPany               0
#define SQCLOUD_IPv4                2
#define SQCLOUD_IPv6                30

// MARK: -

// opaque datatypes
typedef struct SQCloudConnection    SQCloudConnection;
typedef struct SQCloudResult        SQCloudResult;
typedef void (*SQCloudPubSubCB)     (SQCloudConnection *connection, SQCloudResult *result, void *data);

// configuration struct to be passed to the connect function (currently unused)
typedef struct SQCloudConfigStruct {
    const char      *username;
    const char      *password;
    const char      *database;
    int             timeout;
    int             family;                 // can be: SQCLOUD_IPv4, SQCLOUD_IPv6 or SQCLOUD_IPany
    bool            compression;            // compression flag
    bool            sqlite_mode;            // special sqlite compatibility mode
    bool            zero_text;              // flag to tell the server to zero-terminate strings
    bool            password_hashed;        // private flag
    bool            nonlinearizable;        // flag to request for immediate responses from the server node without waiting for linerizability guarantees
    #ifndef SQLITECLOUD_DISABLE_TSL
    const char      *tls_root_certificate;
    const char      *tls_certificate;
    const char      *tls_certificate_key;
    bool            insecure;               // flag to disable TLS
    #endif
} SQCloudConfig;

// convenient struct to be used in SQCloudDownloadDatabase
typedef struct {
    void            *ptr;
    int             fd;
} SQCloudData;

typedef enum {
    RESULT_OK,
    RESULT_ERROR,
    RESULT_STRING,
    RESULT_INTEGER,
    RESULT_FLOAT,
    RESULT_ROWSET,
    RESULT_ARRAY,
    RESULT_NULL,
    RESULT_JSON,
    RESULT_BLOB
} SQCloudResType;

typedef enum {
    VALUE_INTEGER = 1,
    VALUE_FLOAT = 2,
    VALUE_TEXT = 3,
    VALUE_BLOB = 4,
    VALUE_NULL = 5
} SQCloudValueType;

typedef enum {
    ARRAY_TYPE_SQLITE_EXEC = 10,            // used in SQLITE_MODE only when a write statement is executed (instead of the OK reply)
    ARRAY_TYPE_DB_STATUS = 11,
    
    ARRAY_TYPE_VM_STEP = 20,                // used in VM_STEP (when SQLITE_DONE is returned)
    ARRAY_TYPE_VM_PREPARE = 21,             // used in VM_PREPARE
    ARRAY_TYPE_VM_STEP_ONE = 22,            // unused in this version (will be used to step in a server-side rowset)
    ARRAY_TYPE_VM_SQL = 23,
    
    ARRAY_TYPE_BLOB_OPEN = 30,              // used in BLOB_OPEN
    
    ARRAY_TYPE_BACKUP_INIT = 40,            // used in BACKUP_INIT
    ARRAY_TYPE_BACKUP_STEP = 41,            // used in backupWrite (VFS)
    ARRAY_TYPE_BACKUP_END = 42              // used in backupClose (VFS)
} SQCloudArrayType;

typedef enum {
    INTERNAL_ERRCODE_GENERIC = 100000,
    INTERNAL_ERRCODE_PUBSUB = 100001,
    INTERNAL_ERRCODE_TLS = 100002,
    INTERNAL_ERRCODE_URL = 100003,
    INTERNAL_ERRCODE_MEMORY = 100004,
    INTERNAL_ERRCODE_NETWORK = 100005
} INTERNAL_ERRCODE;

// from SQLiteCloud
typedef enum {
    CLOUD_ERRCODE_MEM = 10000,
    CLOUD_ERRCODE_NOTFOUND = 10001,
    CLOUD_ERRCODE_COMMAND = 10002,
    CLOUD_ERRCODE_INTERNAL = 10003,
    CLOUD_ERRCODE_AUTH = 10004,
    CLOUD_ERRCODE_GENERIC = 10005,
    CLOUD_ERRCODE_RAFT = 10006
} CLOUD_ERRCODE;

// MARK: -
SQCloudConnection *SQCloudConnect (const char *hostname, int port, SQCloudConfig *config);
SQCloudConnection *SQCloudConnectWithString (const char *s);
SQCloudResult *SQCloudExec (SQCloudConnection *connection, const char *command);
SQCloudResult *SQCloudRead (SQCloudConnection *connection);
char *SQCloudUUID (SQCloudConnection *connection);
bool SQCloudSendBLOB (SQCloudConnection *connection, void *buffer, uint32_t blen);
void SQCloudDisconnect (SQCloudConnection *connection);
void SQCloudSetPubSubCallback (SQCloudConnection *connection, SQCloudPubSubCB callback, void *data);
SQCloudResult *SQCloudSetPubSubOnly (SQCloudConnection *connection);

// MARK: - Error -
bool SQCloudIsError (SQCloudConnection *connection);
bool SQCloudIsSQLiteError (SQCloudConnection *connection);
int SQCloudErrorCode (SQCloudConnection *connection);
int SQCloudExtendedErrorCode (SQCloudConnection *connection);
const char *SQCloudErrorMsg (SQCloudConnection *connection);
void SQCloudErrorReset (SQCloudConnection *connection);
void SQCloudErrorSetCode (SQCloudConnection *connection, int errcode);
void SQCloudErrorSetMsg (SQCloudConnection *connection, const char *format, ...);

// MARK: - Result -
SQCloudResType SQCloudResultType (SQCloudResult *result);
uint32_t SQCloudResultLen (SQCloudResult *result);
char *SQCloudResultBuffer (SQCloudResult *result);
int32_t SQCloudResultInt32 (SQCloudResult *result);
int64_t SQCloudResultInt64 (SQCloudResult *result);
double SQCloudResultDouble (SQCloudResult *result);
void SQCloudResultFree (SQCloudResult *result);
bool SQCloudResultIsOK (SQCloudResult *result);
bool SQCloudResultIsError (SQCloudResult *result);

// MARK: - Rowset -
SQCloudValueType SQCloudRowsetValueType (SQCloudResult *result, uint32_t row, uint32_t col);
uint32_t SQCloudRowsetRowsMaxColumnLength (SQCloudResult *result, uint32_t col);
char *SQCloudRowsetColumnName (SQCloudResult *result, uint32_t col, uint32_t *len);
uint32_t SQCloudRowsetRows (SQCloudResult *result);
uint32_t SQCloudRowsetCols (SQCloudResult *result);
uint32_t SQCloudRowsetMaxLen (SQCloudResult *result);
char *SQCloudRowsetValue (SQCloudResult *result, uint32_t row, uint32_t col, uint32_t *len);
uint32_t SQCloudRowsetValueLen (SQCloudResult *result, uint32_t row, uint32_t col);
int32_t SQCloudRowsetInt32Value (SQCloudResult *result, uint32_t row, uint32_t col);
int64_t SQCloudRowsetInt64Value (SQCloudResult *result, uint32_t row, uint32_t col);
float SQCloudRowsetFloatValue (SQCloudResult *result, uint32_t row, uint32_t col);
double SQCloudRowsetDoubleValue (SQCloudResult *result, uint32_t row, uint32_t col);
void SQCloudRowsetDump (SQCloudResult *result, uint32_t maxline, bool quiet);

// MARK: - Array -
SQCloudValueType SQCloudArrayValueType (SQCloudResult *result, uint32_t index);
uint32_t SQCloudArrayCount (SQCloudResult *result);
char *SQCloudArrayValue (SQCloudResult *result, uint32_t index, uint32_t *len);
int32_t SQCloudArrayInt32Value (SQCloudResult *result, uint32_t index);
int64_t SQCloudArrayInt64Value (SQCloudResult *result, uint32_t index);
float SQCloudArrayFloatValue (SQCloudResult *result, uint32_t index);
double SQCloudArrayDoubleValue (SQCloudResult *result, uint32_t index);
void SQCloudArrayDump (SQCloudResult *result);

// MARK: - Upload/Download -
bool SQCloudDownloadDatabase (SQCloudConnection *connection, const char *dbname, void *xdata,
                              int (*xCallback)(void *xdata, const void *buffer, uint32_t blen, int64_t ntot, int64_t nprogress));
bool SQCloudUploadDatabase (SQCloudConnection *connection, const char *dbname, const char *key, void *xdata, int64_t dbsize, int (*xCallback)(void *xdata, void *buffer, uint32_t *blen, int64_t ntot, int64_t nprogress));

// MARK: - Base64 -
char *SQCloudBinaryToB64 (char *dest, void const *src, size_t *size);
void *SQCloudB64ToBinary (void *dest, char const *src, size_t *size);
size_t SQCloudComputeB64Size (size_t binarySize);

#ifdef __cplusplus
}
#endif

#endif
