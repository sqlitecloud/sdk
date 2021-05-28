//
//  sqcloud_private.h
//  sqlitecloud
//
//  Created by Marco Bambini on 28/05/21.
//  Copyright © 2021 SQLite Cloud Inc. All rights reserved.
//

#ifndef __SQCLOUD_PRIVATE__
#define __SQCLOUD_PRIVATE__

#include "sqcloud.h"

bool SQCloudForwardExec(SQCloudConnection *connection, const char *command, bool (*forward_cb) (char *buffer, size_t blen, void *xdata), void *xdata);
SQCloudResult *SQCloudSetUUID (SQCloudConnection *connection, const char *UUID);

#endif
