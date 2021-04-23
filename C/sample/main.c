//
//  main.c
//  sqlitecloud-cli
//
//  Created by Marco Bambini on 08/02/21.
//

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <stdbool.h>
#include "sqcloud.h"
#include "linenoise.h"

#define CLI_HISTORY_FILENAME    ".sqlitecloud_history.txt"
#define CLI_VERSION             "1.0a4"
#define CLI_BUILD_DATE          __DATE__

static bool skip_ok = false;

void do_command (SQCloudConnection *conn, char *command) {
    SQCloudResult *res = SQCloudExec(conn, command);
    
    SQCloudResType type = SQCloudResultType(res);
    switch (type) {
        case RESULT_OK:
            if (skip_ok) return;
            printf("OK");
            break;
            
        case RESULT_ERROR:
            printf("ERROR: %s (%d)", SQCloudErrorMsg(conn), SQCloudErrorCode(conn));
            break;
            
        case RESULT_NULL:
            printf("NULL");
            break;
            
        case RESULT_STRING:
        case RESULT_INTEGER:
        case RESULT_FLOAT:
            printf("%.*s", SQCloudResultLen(res), SQCloudResultBuffer(res));
            break;
            
        case RESULT_ROWSET:
            SQCloudRowSetDump(res, 0);
            break;
    }
    
    printf("\n\n");
    SQCloudResultFree(res);
}

void do_command_without_ok_reply (SQCloudConnection *conn, char *command) {
    skip_ok = true;
    do_command(conn, command);
    skip_ok = false;
}

bool do_process_file (SQCloudConnection *conn, const char *filename) {
    // should continue flag set to false by default
    bool result = false;
    
    FILE *file = fopen(filename, "r");
    if (!file) {
        printf("Unable to open file %s.\n", filename);
        return false;
    }
    
    char line[512];
    while (fgets(line, sizeof(line), file)) {
        if (strcasecmp(line, ".PROMPT")) {result = true; break;}
        do_command_without_ok_reply(conn, line);
    }
    
    fclose(file);
    return result;
}

// usage:
// % sqlitecloud-cli -h HOST -p PORT -f FILE -c
int main(int argc, char * argv[]) {
    const char *hostname = "localhost";
    const char *filename = NULL;
    int port = SQCLOUD_DEFAULT_PORT;
    bool compression = false;
    int c;
    
    while ((c = getopt (argc, argv, "h:p:f:c")) != -1) {
        switch (c) {
            case 'h': hostname = optarg; break;
            case 'p': port = atoi(optarg); break;
            case 'f': filename = optarg; break;
            case 'c': compression = true; break;
        }
    }
    
    printf("sqlitecloud-cli version %s (build date %s)\n", CLI_VERSION, CLI_BUILD_DATE);
    
    SQCloudConnection *conn = SQCloudConnect(hostname, port, NULL);
    if (SQCloudIsError(conn)) {
        printf("ERROR connecting to %s: %s (%d)\n", hostname, SQCloudErrorMsg(conn), SQCloudErrorCode(conn));
        return -1;
    } else {
        printf("Connection to %s OK...\n\n", hostname);
    }
    
    // load history file
    linenoiseHistoryLoad(CLI_HISTORY_FILENAME);
    
    // activate compression
    if (compression) do_command_without_ok_reply(conn, "SET KEY CLIENT_COMPRESSION TO 1");
    
    if (filename) {
        bool should_continue = do_process_file(conn, filename);
        if (should_continue == false) return 0;
    }
    
    // REPL
    char *command = NULL;
    while((command = linenoise(">> ")) != NULL) {
        if (command[0] != '\0') {
            linenoiseHistoryAdd(command);
            linenoiseHistorySave(CLI_HISTORY_FILENAME);
        }
        if (strncmp(command, "EXIT", 4) == 0) break;
        do_command(conn, command);
    }
    if (command) free(command);
    
    SQCloudDisconnect(conn);
    return 0;
}
