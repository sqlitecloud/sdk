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

#define SQCLOUD_SDK_VERSION         "0.5.0"
#define SQCLOUD_SDK_VERSION_NUM     0x000500
#define SQCLOUD_DEFAULT_PORT        8860
#define SQCLOUD_DEFAULT_TIMEOUT     12

#define SQCLOUD_IPany               0
#define SQCLOUD_IPv4                2
#define SQCLOUD_IPv6                30

// opaque datatypes
typedef struct SQCloudConnection    SQCloudConnection;
typedef struct SQCloudResult        SQCloudResult;
typedef void (*SQCloudPubSubCB)    (SQCloudConnection *connection, SQCloudResult *result, void *data);

// configuration struct to be passed to the connect function (currently unused)
typedef struct SQCloudConfigStruct {
    const char  *username;
    const char  *password;
    const char  *database;
    int         timeout;
    int         family;             // can be: AF_INET, AF_INET6 or AF_UNSPEC
    #ifndef SQLITECLOUD_DISABLE_TSL
    const char  *tls_root_certificate;
    const char  *tls_certificate;
    const char  *tls_certificate_key;
    bool        insecure;           // flag to disable TLS
    #endif
} SQCloudConfig;

typedef enum {
    RESULT_OK,
    RESULT_ERROR,
    RESULT_STRING,
    RESULT_INTEGER,
    RESULT_FLOAT,
    RESULT_ROWSET,
    RESULT_NULL,
    RESULT_JSON
} SQCloudResType;

typedef enum {
    VALUE_INTEGER = 1,
    VALUE_FLOAT = 2,
    VALUE_TEXT = 3,
    VALUE_BLOB = 4,
    VALUE_NULL = 5
} SQCloudValueType;

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

SQCloudConnection *SQCloudConnect (const char *hostname, int port, SQCloudConfig *config);
SQCloudConnection *SQCloudConnectWithString (const char *s);
SQCloudResult *SQCloudExec (SQCloudConnection *connection, const char *command);
char *SQCloudUUID (SQCloudConnection *connection);
void SQCloudDisconnect (SQCloudConnection *connection);
void SQCloudSetPubSubCallback (SQCloudConnection *connection, SQCloudPubSubCB callback, void *data);
SQCloudResult *SQCloudSetPubSubOnly (SQCloudConnection *connection);

bool SQCloudIsError (SQCloudConnection *connection);
int SQCloudErrorCode (SQCloudConnection *connection);
const char *SQCloudErrorMsg (SQCloudConnection *connection);

SQCloudResType SQCloudResultType (SQCloudResult *result);
uint32_t SQCloudResultLen (SQCloudResult *result);
char *SQCloudResultBuffer (SQCloudResult *result);
void SQCloudResultFree (SQCloudResult *result);
bool SQCloudResultIsOK (SQCloudResult *result);

SQCloudValueType SQCloudRowsetValueType (SQCloudResult *result, uint32_t row, uint32_t col);
uint32_t SQCloudRowsetRowsMaxColumnLength (SQCloudResult *result, uint32_t col);
char *SQCloudRowsetColumnName (SQCloudResult *result, uint32_t col, uint32_t *len);
uint32_t SQCloudRowsetRows (SQCloudResult *result);
uint32_t SQCloudRowsetCols (SQCloudResult *result);
uint32_t SQCloudRowsetMaxLen (SQCloudResult *result);
char *SQCloudRowsetValue (SQCloudResult *result, uint32_t row, uint32_t col, uint32_t *len);
int32_t SQCloudRowsetInt32Value (SQCloudResult *result, uint32_t row, uint32_t col);
int64_t SQCloudRowsetInt64Value (SQCloudResult *result, uint32_t row, uint32_t col);
float SQCloudRowsetFloatValue (SQCloudResult *result, uint32_t row, uint32_t col);
double SQCloudRowsetDoubleValue (SQCloudResult *result, uint32_t row, uint32_t col);
void SQCloudRowsetDump (SQCloudResult *result, uint32_t maxline);

#ifdef __cplusplus
}
#endif

#endif
