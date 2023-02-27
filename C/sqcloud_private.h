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

#define SQCloudExecBuffer               _reserved0
#define SQCloudForwardExec              _reserved1
#define SQCloudSetClientUUID            _reserved2
#define sqcloud_parse_buffer            _reserved3
#define sqcloud_parse_number            _reserved4
#define sqcloud_result_is_chunk         _reserved5
#define SQCloudTransferDatabase         _reserved8
#define sqcloud_parse_type              _reserved9
#define sqcloud_parse_value             _reserved10
#define sqcloud_parse_array_value       _reserved11
#define sqcloud_parse_array_count       _reserved12
#define SQCloudDownloadDatabaseInternal _reserved13

SQCloudResult *SQCloudExecBuffer (SQCloudConnection *connection, const char *buffer, size_t blen, bool compute_header);
bool SQCloudForwardExec(SQCloudConnection *connection, const char *command, size_t len, bool compute_header, bool (*forward_cb) (char *buffer, size_t blen, void *xdata, void *xdata2), void *xdata, void *xdata2);
SQCloudResult *SQCloudSetClientUUID (SQCloudConnection *connection, const char *UUID);
bool SQCloudTransferDatabase (SQCloudConnection *connection, const char *dbname, const char *key, uint64_t snapshotid, bool isinternaldb, void *xdata, int64_t dbsize, int (*xCallback)(void *xdata, void *buffer, uint32_t *blen, int64_t ntot, int64_t nprogress));
bool SQCloudDownloadDatabaseInternal (SQCloudConnection *connection, const char *dbname, void *xdata,
                              int (*xCallback)(void *xdata, const void *buffer, uint32_t blen, int64_t ntot, int64_t nprogress), uint64_t *raft_index, bool ifexists);

SQCloudResult *sqcloud_parse_buffer (char *buffer, uint32_t blen, uint32_t cstart, SQCloudResult *chunk);
uint32_t sqcloud_parse_number (char *buffer, uint32_t blen, uint32_t *cstart);
char *sqcloud_parse_value (char *buffer, uint32_t *len, uint32_t *cellsize);
char *sqcloud_parse_array_value(char *buffer, uint32_t blen, uint32_t index, uint32_t *len, uint32_t *cellsize, uint32_t *pos, int *type, SQCLOUD_INTERNAL_ERRCODE *errcode);
int32_t sqcloud_parse_array_count(char *buffer, uint32_t blen);
bool sqcloud_result_is_chunk (SQCloudResult *res);
int sqcloud_parse_type (char *buffer);

#endif
