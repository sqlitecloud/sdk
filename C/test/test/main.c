//
//  main.c
//  test
//
//  Created by Marco Bambini on 10/09/22.
//

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include "sqcloud.h"

#define CONNECTION_STRING   "sqlitecloud://admin:admin@localhost/compression=1&root_certificate=%2FUsers%2Fmarco%2FDesktop%2FSQLiteCloud%2FPHP%2Fca.pem"
#define BLOB_FILENAME       "/Users/marco/Desktop/test.jpg"
#define BLOB_LEN            445922

// MARK: -

static bool do_print (SQCloudConnection *conn, SQCloudResult *res) {
    // res NULL means to read error message and error code from conn
    SQCloudResType type = SQCloudResultType(res);
    bool result = true;
    
    switch (type) {
        case RESULT_OK:
            printf("OK");
            break;
            
        case RESULT_ERROR:
            printf("ERROR: %s (%d)", SQCloudErrorMsg(conn), SQCloudErrorCode(conn));
            result = false;
            break;
            
        case RESULT_NULL:
            printf("NULL");
            break;
            
        case RESULT_STRING:
            (SQCloudResultLen(res)) ? printf("%.*s", SQCloudResultLen(res), SQCloudResultBuffer(res)) : printf("");
            break;
            
        case RESULT_JSON:
        case RESULT_INTEGER:
        case RESULT_FLOAT:
            printf("%.*s", SQCloudResultLen(res), SQCloudResultBuffer(res));
            break;
            
        case RESULT_ARRAY:
            SQCloudArrayDump(res);
            break;
            
        case RESULT_ROWSET:
            SQCloudRowsetDump(res, 0, false);
            break;
            
        case RESULT_BLOB:
            printf("BLOB data with len: %d", SQCloudResultLen(res));
            break;
    }
    
    printf("\n\n");
    return result;
}

static bool do_command (SQCloudConnection *conn, char *command, int32_t *int_value, char **string_value, void **blob_value, uint32_t *len) {
    printf("%s\n", command);
    
    SQCloudResult *res = SQCloudExec(conn, command);
    if ((int_value) && (SQCloudResultType(res) == RESULT_INTEGER)) *int_value = SQCloudResultInt32(res);
    if ((string_value) && SQCloudResultType(res) == RESULT_STRING) *string_value = SQCloudResultBuffer(res);
    if ((blob_value) && SQCloudResultType(res) == RESULT_BLOB) *blob_value = SQCloudResultBuffer(res);
    if (len) *len = SQCloudResultLen(res);
    
    bool result = do_print(conn, res);
    SQCloudResultFree(res);
    return result;
}

// MARK: -

static bool test_read_blob (SQCloudConnection *conn) {
    const char *filename = BLOB_FILENAME;
    unlink(filename);
    FILE *f = fopen(filename, "w");
    if (!f) {perror("Error creating file in test_read_blob"); return false;}
    
    int32_t blob_index = 0;
    
    // BLOB OPEN <database_name> TABLE <table_name> COLUMN <column_name> ROWID <rowid> RWFLAG <rwflag>
    if (!do_command(conn, "BLOB OPEN main TABLE images COLUMN picture ROWID 2 RWFLAG 0;", &blob_index, NULL, NULL, NULL)) return false;
    
    // SELECT length(picture) FROM images WHERE rowid=2; => 445922
    int32_t len = 0, lblob = BLOB_LEN;
    
    // BLOB BYTES <index>
    char command[256];
    snprintf(command, sizeof(command), "BLOB BYTES %d", blob_index);
    if (!do_command(conn, command, &len, NULL, NULL, NULL)) return false;
    if (len != lblob) {printf("BLOB size is wrong %d %d (test_read_blob)\n", len, lblob); return false;}
    
    // PERFORM 3 reads
    char buffer[1024*100];
    int offset = 0;
    int blen = sizeof(buffer);
    
    while (1) {
        void *data = NULL;
        uint32_t len = 0;
        
        // BLOB READ <index> SIZE <size> OFFSET <offset>
        snprintf(command, sizeof(command), "BLOB READ %d SIZE %d OFFSET %d", blob_index, blen, offset);
        if (!do_command(conn, command, NULL, NULL, &data, &len)) return false;
        
        if (blen != len) {printf("BLOB read returned a wrong len %d != %d\n", blen, len); return false;}
        
        size_t fwrote = fwrite(data, blen, 1, f);
        if (fwrote != 1) {perror("Error writing BLOB data"); return false;}
        
        offset += blen;
        if (offset == lblob) break;
        if (lblob - offset < blen) blen = lblob - offset;
    }
    
    // close test file
    fclose(f);
    
    // finalize BLOB vm
    snprintf(command, sizeof(command), "BLOB CLOSE %d", blob_index);
    if (!do_command(conn, command, NULL, NULL, NULL, NULL)) return false;
    
    return true;
}

static bool test_write_blob (SQCloudConnection *conn) {
    const char *filename = BLOB_FILENAME;
    FILE *f = fopen(filename, "rb");
    if (!f) {perror("Error reading file in test_write_blob"); return false;}
    
    char command[256];
    snprintf(command, sizeof(command), "INSERT INTO images (picture) VALUES (zeroblob(%d));", BLOB_LEN);
    if (!do_command(conn, command, NULL, NULL, NULL, NULL)) return false;
    
    int32_t rowid = 0;
    if (!do_command(conn, "DATABASE GET ROWID;", &rowid, NULL, NULL, NULL)) return false;
    
    int32_t blob_index = 0;
    
    // BLOB OPEN <database_name> TABLE <table_name> COLUMN <column_name> ROWID <rowid> RWFLAG <rwflag>
    snprintf(command, sizeof(command), "BLOB OPEN main TABLE images COLUMN picture ROWID %d RWFLAG 1;", rowid);
    if (!do_command(conn, command, &blob_index, NULL, NULL, NULL)) return false;
    
    // prepare buffer
    char buffer[1024*100];
    int blen = sizeof(buffer);
    int offset = 0;
    
    // loop to read input file and write BLOB
    while (!feof(f)) {
        size_t nbytes = fread(buffer, 1, blen, f);
        
        const char *values[] = {buffer};
        uint32_t len[] = {(uint32_t)nbytes};
        SQCloudValueType types[] = {VALUE_BLOB};
        
        // BLOB WRITE <index> OFFSET <offset> DATA <data>
        snprintf(command, sizeof(command), "BLOB WRITE %d OFFSET %d DATA ?;", blob_index, offset);
        
        printf("%s\n", command);
        SQCloudResult *result = SQCloudExecArray(conn, command, values, len, types, 1);
        SQCloudResultDump(conn, result);
        
        offset += nbytes;
    }
    
    if (offset != BLOB_LEN) {printf("BLOB size is wrong %d %d (test_write_blob)\n", offset, BLOB_LEN); return false;}
    
    // close test file
    fclose(f);
    
    // BLOB BYTES <index>
    int32_t len = 0;
    snprintf(command, sizeof(command), "BLOB BYTES %d", blob_index);
    if (!do_command(conn, command, &len, NULL, NULL, NULL)) return false;
    if (len != BLOB_LEN) {printf("BLOB size is wrong %d %d (test_write_blob 2)\n", len, BLOB_LEN); return false;}
    
    // finalize BLOB vm
    snprintf(command, sizeof(command), "BLOB CLOSE %d", blob_index);
    if (!do_command(conn, command, NULL, NULL, NULL, NULL)) return false;
    
    return true;
}

static int test_blob (SQCloudConnection *conn) {
    if (!do_command(conn, "USE DATABASE images.sqlite;", NULL, NULL, NULL, NULL)) goto abort_test;
    if (!test_read_blob(conn)) goto abort_test;
    if (!test_write_blob(conn)) goto abort_test;
    
    return 0;
    
abort_test:
    exit(-1);
    return -1;
}

// MARK: -

static int test_array (SQCloudConnection *conn) {
    // built-in with 1 binding
    {
        // USE DATABASE <database_name>
        const char *dbname = "mediastore.sqlite";
        const char *values[] = {dbname};
        uint32_t len[] = {(uint32_t)strlen(dbname)};
        SQCloudValueType types[] = {VALUE_TEXT};
        
        const char *command = "USE DATABASE ?";
        printf("%s\n", command);
        SQCloudResult *result = SQCloudExecArray(conn, command, values, len, types, 1);
        SQCloudResultDump(conn, result);
    }
    
    // SQLite with 1 binding
    {
        // SELECT * FROM Artist
        const char *value = "100";
        const char *values[] = {value};
        uint32_t len[] = {(uint32_t)strlen(value)};
        SQCloudValueType types[] = {VALUE_INTEGER};
        
        const char *command = "SELECT * FROM Artist WHERE ArtistId >= ?";
        printf("%s\n", command);
        SQCloudResult *result = SQCloudExecArray(conn, command, values, len, types, 1);
        SQCloudResultDump(conn, result);
    }
    
    // SQLite with 2 bindings
    {
        // SELECT * FROM Artist
        const char *value1 = "100";
        const char *value2 = "200";
        const char *values[] = {value1, value2};
        uint32_t len[] = {(uint32_t)strlen(value1), (uint32_t)strlen(value2)};
        SQCloudValueType types[] = {VALUE_INTEGER, VALUE_INTEGER};
        
        const char *command = "SELECT * FROM Artist WHERE ArtistId >= ? AND ArtistId <= ?";
        printf("%s\n", command);
        SQCloudResult *result = SQCloudExecArray(conn, command, values, len, types, 2);
        SQCloudResultDump(conn, result);
    }
    
    // built-in with 4 bindings
    {
        // SET APIKEY <key> [NAME <key_name>] [RESTRICTION <restriction_type>] [EXPIRATION <expiration_date>]
        const char *key = "aJcdAL6P1JwJHquTP5iK1ahk7b3tAicBBufPSmnkIb4";
        const char *name = "Bind Test";
        const char *restriction = "1";
        const char *expiration = "2022-09-21 18:27:29";
        const char *values[] = {key, name, restriction, expiration};
        uint32_t len[] = {(uint32_t)strlen(key), (uint32_t)strlen(name), (uint32_t)strlen(restriction), (uint32_t)strlen(expiration)};
        SQCloudValueType types[] = {VALUE_TEXT, VALUE_TEXT, VALUE_TEXT, VALUE_TEXT};
        
        const char *command = "SET APIKEY ? NAME ? RESTRICTION ? EXPIRATION ?";
        printf("%s\n", command);
        SQCloudResult *result = SQCloudExecArray(conn, command, values, len, types, 4);
        SQCloudResultDump(conn, result);
    }
    
    // VM
    {
        const char *command = "SELECT * FROM Artist WHERE ArtistId >= ? AND ArtistId <= ?";
        printf("%s\n", command);
        SQCloudVM *vm = SQCloudVMCompile (conn, command, -1, NULL);
        
        bool result = SQCloudVMBindInt (vm, 1, 100);
        result = SQCloudVMBindInt (vm, 2, 105);
        
        /*SQCloudResType type = */SQCloudVMStep(vm);
        SQCloudResult *r = SQCloudVMResult(vm);
        SQCloudResultDump(conn, r);
        
        SQCloudVMClose(vm);
    }
    
    return 0;
}

// MARK: -

static int test_database (SQCloudConnection *conn) {
    
    return 0;
}

// MARK: -

int main (int argc, const char * argv[]) {
    SQCloudConnection *conn = SQCloudConnectWithString(CONNECTION_STRING);
    if (SQCloudIsError(conn)) {
        printf("ERROR connecting: %s (%d)\n", SQCloudErrorMsg(conn), SQCloudErrorCode(conn));
        return -1;
    } else {
        printf("Connection to host OK...\n\n");
    }
    
    test_blob(conn);
    test_array(conn);
    test_database(conn);
    
    SQCloudDisconnect(conn);
    return 0;
}
