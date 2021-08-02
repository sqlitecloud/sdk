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

#define SQCloudForwardExec      _reserved1
#define SQCloudSetUUID          _reserved2
#define sqcloud_parse_buffer    _reserved3
#define sqcloud_parse_number    _reserved4

bool SQCloudForwardExec(SQCloudConnection *connection, const char *command, bool (*forward_cb) (char *buffer, size_t blen, void *xdata), void *xdata);
SQCloudResult *SQCloudSetUUID (SQCloudConnection *connection, const char *UUID);

SQCloudResult *sqcloud_parse_buffer (char *buffer, uint32_t blen, uint32_t cstart);
uint32_t sqcloud_parse_number (char *buffer, uint32_t blen, uint32_t *cstart);

#endif
