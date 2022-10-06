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

#define CONNECTION_STRING               "sqlitecloud://admin:admin@localhost/compression=1&root_certificate=%2FUsers%2Fmarco%2FDesktop%2FSQLiteCloud%2FPHP%2Fca.pem"
#define CONNECTION_STRING_INSECURE      "sqlitecloud://admin:admin@localhost/compression=1&insecure=1"
#define BLOB_FILENAME                   "/Users/marco/Desktop/test.jpg"
#define BLOB_LEN                        445922
#define BACKUP_FILENAME                 "/Users/marco/Desktop/test.sqlite"
#define VERBOSE_OUTPUT                  0

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

static bool do_command (SQCloudConnection *conn, char *command, int32_t *int_value, char **string_value, void **blob_value, uint32_t *len, SQCloudResult **result_value) {
    printf("%s\n", command);
    
    SQCloudResult *res = SQCloudExec(conn, command);
    if ((int_value) && (SQCloudResultType(res) == RESULT_INTEGER)) *int_value = SQCloudResultInt32(res);
    else if ((string_value) && SQCloudResultType(res) == RESULT_STRING) *string_value = SQCloudResultBuffer(res);
    else if ((blob_value) && SQCloudResultType(res) == RESULT_BLOB) *blob_value = SQCloudResultBuffer(res);
    else if (result_value) *result_value = res;
    if (len) *len = SQCloudResultLen(res);
    
    bool result = do_print(conn, res);
    if (!result_value) SQCloudResultFree(res);
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
    if (!do_command(conn, "BLOB OPEN main TABLE images COLUMN picture ROWID 2 RWFLAG 0;", &blob_index, NULL, NULL, NULL, NULL)) return false;
    
    // SELECT length(picture) FROM images WHERE rowid=2; => 445922
    int32_t len = 0, lblob = BLOB_LEN;
    
    // BLOB BYTES <index>
    char command[256];
    snprintf(command, sizeof(command), "BLOB BYTES %d", blob_index);
    if (!do_command(conn, command, &len, NULL, NULL, NULL, NULL)) return false;
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
        if (!do_command(conn, command, NULL, NULL, &data, &len, NULL)) return false;
        
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
    if (!do_command(conn, command, NULL, NULL, NULL, NULL, NULL)) return false;
    
    return true;
}

static bool test_write_blob (SQCloudConnection *conn) {
    const char *filename = BLOB_FILENAME;
    FILE *f = fopen(filename, "rb");
    if (!f) {perror("Error reading file in test_write_blob"); return false;}
    
    char command[256];
    snprintf(command, sizeof(command), "INSERT INTO images (picture) VALUES (zeroblob(%d));", BLOB_LEN);
    if (!do_command(conn, command, NULL, NULL, NULL, NULL, NULL)) return false;
    
    int32_t rowid = 0;
    if (!do_command(conn, "DATABASE GET ROWID;", &rowid, NULL, NULL, NULL, NULL)) return false;
    
    int32_t blob_index = 0;
    
    // BLOB OPEN <database_name> TABLE <table_name> COLUMN <column_name> ROWID <rowid> RWFLAG <rwflag>
    snprintf(command, sizeof(command), "BLOB OPEN main TABLE images COLUMN picture ROWID %d RWFLAG 1;", rowid);
    if (!do_command(conn, command, &blob_index, NULL, NULL, NULL, NULL)) return false;
    
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
    if (!do_command(conn, command, &len, NULL, NULL, NULL, NULL)) return false;
    if (len != BLOB_LEN) {printf("BLOB size is wrong %d %d (test_write_blob 2)\n", len, BLOB_LEN); return false;}
    
    // finalize BLOB vm
    snprintf(command, sizeof(command), "BLOB CLOSE %d", blob_index);
    if (!do_command(conn, command, NULL, NULL, NULL, NULL, NULL)) return false;
    
    return true;
}

static int test_blob (SQCloudConnection *conn) {
    if (!do_command(conn, "USE DATABASE images.sqlite;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
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

static int test_backup (SQCloudConnection *conn) {
    if (!do_command(conn, "USE DATABASE mediastore.sqlite;", NULL, NULL, NULL, NULL, NULL)) return -1;
    
    const char *filename = BACKUP_FILENAME;
    unlink(filename);
    FILE *f = fopen(filename, "w");
    if (!f) {perror("Error creating file in test_backup"); return -1;}
    
    // BACKUP INIT [<dest_name>] [SOURCE <source_name>]
    SQCloudResult *result = NULL;
    if (!do_command(conn, "BACKUP INIT", NULL, NULL, NULL, NULL, &result)) return -1;
    
    // sanity check
    if (SQCloudResultType(result) != RESULT_ARRAY) {
        printf("Wrong result type\n");
        return -1;
    }
    
    // extract information
    int32_t index = SQCloudArrayInt32Value(result, 1);
    int32_t page_size = SQCloudArrayInt32Value(result, 2);
    
    #if VERBOSE_OUTPUT
    int32_t page_count = SQCloudArrayInt32Value(result, 3);
    int counter = 0;
    printf("Page Size: %d\n", page_size);
    printf("Page Count: %d\n", page_count);
    #endif
    
    SQCloudResultFree(result);
    
    char command[512];
    int64_t current_offset = 0;
    int npages = 40;
    int32_t total = 0;
    
    while (1) {
        // BACKUP STEP <index> PAGES <npages>
        snprintf(command, sizeof(command), "BACKUP STEP %d PAGES %d;", index, npages);
        if (!do_command(conn, command, NULL, NULL, NULL, NULL, &result)) goto abort_backup;
        
        SQCloudResType rtype = SQCloudResultType(result);
        if (rtype != RESULT_ARRAY) goto abort_backup;
        
        /*
         [0] -> TYPE
         [1] -> INDEX
         [2] -> PAGE_TOTALE
         [3] -> PAGE_REMAINING
         [4] -> PAGE_COUNTER
         [5] -> BLOB
         */
        
        #if VERBOSE_OUTPUT
        int32_t page_totals = SQCloudArrayInt32Value(result, 2);
        #endif
        int32_t page_remaining = SQCloudArrayInt32Value(result, 3);
        int32_t page_counter = SQCloudArrayInt32Value(result, 4);
        
        #if VERBOSE_OUTPUT
        printf("Counter: %d\n", ++counter);
        printf("Page Total: %d\n", page_totals);
        printf("Page Remaining: %d\n", page_remaining);
        printf("Page Counter: %d\n", page_counter);
        #endif
        
        uint32_t blen = 0;
        char *data = SQCloudArrayValue (result, 5, &blen);
        #if VERBOSE_OUTPUT
        printf("Len: %d\n", blen);
        #endif
        
        // sanity check
        if (blen != ((page_size * page_counter) + (page_counter * sizeof(int64_t)))) {
            printf("Block size error\n");
            goto abort_backup;
        }
        
        char *p = data;
        for (int i=0; i<page_counter; ++i) {
            uint64_t *poff = (uint64_t *)p;
            int64_t offset = (int64_t)ntohll(*poff);
            p += sizeof(int64_t);
            
            if (current_offset != offset) {
                if (fseek(f, offset, SEEK_SET) != 0) {perror("Error seeking in the outfile file"); goto abort_backup;}
                current_offset = offset;
            }
            
            size_t fwrote = fwrite(p, page_size, 1, f);
            if (fwrote != 1) {perror("Error writing BLOB data"); goto abort_backup;}
            
            total += page_size;
            p += page_size;
        }
        
        #if VERBOSE_OUTPUT
        printf("Total: %d\n", total);
        #endif
        if (page_remaining == 0) break;
    }
     
    // extract some information
    snprintf(command, sizeof(command), "BACKUP REMAINING %d", index);
    if (!do_command(conn, command, NULL, NULL, NULL, NULL, &result)) return -1;
    
    snprintf(command, sizeof(command), "BACKUP PAGECOUNT %d", index);
    if (!do_command(conn, command, NULL, NULL, NULL, NULL, &result)) return -1;
    
    // close backup
    snprintf(command, sizeof(command), "BACKUP FINISH %d", index);
    if (!do_command(conn, command, NULL, NULL, NULL, NULL, &result)) return -1;
    
    fclose(f);
    return 0;
    
abort_backup:
    snprintf(command, sizeof(command), "BACKUP FINISH %d", index);
    do_command(conn, command, NULL, NULL, NULL, NULL, NULL);
    fclose(f);
    exit(-1);
    return -1;
}

// MARK: -

static int test_database (SQCloudConnection *conn) {
    if (!do_command(conn, "USE DATABASE mediastore.sqlite;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    
    if (!do_command(conn, "DATABASE FILENAME main;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE READONLY main;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE ERRNO;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE TXNSTATE main;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    
    /*
     #define SQLITE_DBSTATUS_LOOKASIDE_USED       0
     #define SQLITE_DBSTATUS_CACHE_USED           1
     #define SQLITE_DBSTATUS_SCHEMA_USED          2
     #define SQLITE_DBSTATUS_STMT_USED            3
     #define SQLITE_DBSTATUS_LOOKASIDE_HIT        4
     #define SQLITE_DBSTATUS_LOOKASIDE_MISS_SIZE  5
     #define SQLITE_DBSTATUS_LOOKASIDE_MISS_FULL  6
     #define SQLITE_DBSTATUS_CACHE_HIT            7
     #define SQLITE_DBSTATUS_CACHE_MISS           8
     #define SQLITE_DBSTATUS_CACHE_WRITE          9
     #define SQLITE_DBSTATUS_DEFERRED_FKS        10
     #define SQLITE_DBSTATUS_CACHE_USED_SHARED   11
     #define SQLITE_DBSTATUS_CACHE_SPILL         12
     */
    
    if (!do_command(conn, "DATABASE STATUS 0 RESET 0;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE STATUS 1 RESET 0;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE STATUS 2 RESET 0;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE STATUS 3 RESET 0;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE STATUS 4 RESET 0;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE STATUS 5 RESET 0;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE STATUS 6 RESET 0;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE STATUS 7 RESET 0;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE STATUS 8 RESET 0;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE STATUS 9 RESET 0;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE STATUS 10 RESET 0;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE STATUS 11 RESET 0;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE STATUS 12 RESET 0;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    
    /*
     #define SQLITE_LIMIT_LENGTH                    0
     #define SQLITE_LIMIT_SQL_LENGTH                1
     #define SQLITE_LIMIT_COLUMN                    2
     #define SQLITE_LIMIT_EXPR_DEPTH                3
     #define SQLITE_LIMIT_COMPOUND_SELECT           4
     #define SQLITE_LIMIT_VDBE_OP                   5
     #define SQLITE_LIMIT_FUNCTION_ARG              6
     #define SQLITE_LIMIT_ATTACHED                  7
     #define SQLITE_LIMIT_LIKE_PATTERN_LENGTH       8
     #define SQLITE_LIMIT_VARIABLE_NUMBER           9
     #define SQLITE_LIMIT_TRIGGER_DEPTH            10
     #define SQLITE_LIMIT_WORKER_THREADS           11
     */
    
    if (!do_command(conn, "DATABASE LIMIT 0 VALUE -1;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE LIMIT 1 VALUE -1;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE LIMIT 2 VALUE -1;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE LIMIT 3 VALUE -1;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE LIMIT 4 VALUE -1;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE LIMIT 5 VALUE -1;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE LIMIT 6 VALUE -1;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE LIMIT 7 VALUE -1;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE LIMIT 8 VALUE -1;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE LIMIT 9 VALUE -1;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE LIMIT 10 VALUE -1;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    if (!do_command(conn, "DATABASE LIMIT 11 VALUE -1;", NULL, NULL, NULL, NULL, NULL)) goto abort_test;
    
    return 0;
    
abort_test:
    exit(-1);
    return -1;
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
    test_backup(conn);
    test_database(conn);
    
    SQCloudDisconnect(conn);
    return 0;
}
