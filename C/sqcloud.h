//
//  sqcloud.h
//
//  Created by Marco Bambini on 08/02/21.
//

#ifndef __SQCLOUD_CLI__
#define __SQCLOUD_CLI__

#include <stdio.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

#define SQCLOUD_SDK_VERSION         "0.2.0"
#define SQCLOUD_SDK_VERSION_NUM     0x000200
#define SQCLOUD_DEFAULT_PORT        8860
#define SQCLOUD_DEFAULT_TIMEOUT     12

// opaque datatypes
typedef struct SQCloudConnection    SQCloudConnection;
typedef struct SQCloudResult        SQCloudResult;

// configuration struct to be passed to the connect function (currently unused)
typedef struct {
    const char *username;
    const char *password;
    const char *database;
    int timeout;
} SQCloudConfig;

typedef enum {
    RESULT_OK,
    RESULT_ERROR,
    RESULT_STRING,
    RESULT_INTEGER,
    RESULT_FLOAT,
    RESULT_ROWSET,
    RESULT_NULL
} SQCloudResType;

typedef enum {
    VALUE_INTEGER = 1,
    VALUE_FLOAT = 2,
    VALUE_TEXT = 3,
    VALUE_BLOB = 4,
    VALUE_NULL = 5
} SQCloudValueType;

SQCloudConnection *SQCloudConnect (const char *hostname, int port, SQCloudConfig *config);
SQCloudResult *SQCloudExec (SQCloudConnection *connection, const char *command);
void SQCloudDisconnect (SQCloudConnection *connection);

bool SQCloudIsError (SQCloudConnection *connection);
int SQCloudErrorCode (SQCloudConnection *connection);
const char *SQCloudErrorMsg (SQCloudConnection *connection);

SQCloudResType SQCloudResultType (SQCloudResult *result);
uint32_t SQCloudResultLen (SQCloudResult *result);
char *SQCloudResultBuffer (SQCloudResult *result);
void SQCloudResultFree (SQCloudResult *result);

SQCloudValueType SQCloudRowSetValueType (SQCloudResult *result, uint32_t row, uint32_t col);
char *SQCloudRowSetColumnName (SQCloudResult *result, uint32_t col, uint32_t *len);
uint32_t SQCloudRowSetRows (SQCloudResult *result);
uint32_t SQCloudRowSetCols (SQCloudResult *result);
char *SQCloudRowSetValue (SQCloudResult *result, uint32_t row, uint32_t col, uint32_t *len);
int32_t SQCloudRowSetInt32Value (SQCloudResult *result, uint32_t row, uint32_t col);
int64_t SQCloudRowSetInt64Value (SQCloudResult *result, uint32_t row, uint32_t col);
float SQCloudRowSetFloatValue (SQCloudResult *result, uint32_t row, uint32_t col);
double SQCloudRowSetDoubleValue (SQCloudResult *result, uint32_t row, uint32_t col);
void SQCloudRowSetDump (SQCloudResult *result, uint32_t maxline);

#ifdef __cplusplus
}
#endif

#endif
