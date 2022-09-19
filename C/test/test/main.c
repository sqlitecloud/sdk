//
//  main.c
//  test
//
//  Created by Marco Bambini on 10/09/22.
//

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include "sqcloud.h"

#define CONNECTION_STRING   "sqlitecloud://admin:admin@localhost/compression=1&root_certificate=%2FUsers%2Fmarco%2FDesktop%2FSQLiteCloud%2FPHP%2Fca.pem"

int main (int argc, const char * argv[]) {
    SQCloudConnection *conn = SQCloudConnectWithString(CONNECTION_STRING);
    if (SQCloudIsError(conn)) {
        printf("ERROR connecting: %s (%d)\n", SQCloudErrorMsg(conn), SQCloudErrorCode(conn));
        return -1;
    } else {
        printf("Connection to host OK...\n\n");
    }
    
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
        
        SQCloudResType type = SQCloudVMStep(vm);
        SQCloudResult *r = SQCloudVMResult(vm);
        SQCloudResultDump(conn, r);
        
        SQCloudVMClose(vm);
    }
    
    SQCloudDisconnect(conn);
}
