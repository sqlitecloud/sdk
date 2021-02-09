//
//  sqcloud.h
//
//  Created by Marco Bambini on 08/02/21.
//

#ifndef __SQCLOUD_CLI__
#define __SQCLOUD_CLI__

#include <stdio.h>
#include <stdbool.h>

#define SQCLOUD_SDK_VERSION         "0.1.0"
#define SQCLOUD_SDK_VERSION_NUM     0x000100
#define SQCLOUD_DEFAULT_PORT        8860
#define SQCLOUD_DEFAULT_TIMEOUT     12

typedef struct SQCloudConnection    SQCloudConnection;
typedef struct SQCloudResult        SQCloudResult;

typedef struct {
    const char *username;
    const char *password;
    const char *database;
    int timeout;
} SQCloudConfig;

typedef enum {
    TYPE_OK,
    TYPE_ERROR,
    TYPE_ROWSET
} SQCloudResType;

SQCloudConnection *SQCloudConnect (const char *hostname, int port, SQCloudConfig *config);
SQCloudResult *SQCloudExec (SQCloudConnection *connection, const char *command);
void SQCloudFree (SQCloudConnection *connection);

bool SQCloudIsError (SQCloudConnection *connection);
int SQCloudErrorCode (SQCloudConnection *connection);
const char *SQCloudErrorMsg (SQCloudConnection *connection);

SQCloudResType SQCloudResultType (SQCloudResult *result);
void SQCloudResultFree (SQCloudResult *result);

void SQCloudRowSetDump (SQCloudResult *result);

#endif
