//
//  main.c
//  sqlitecloud-cli
//
//  Created by Marco Bambini on 08/02/21.
//

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "sqcloud.h"
#include "linenoise.h"

#define CLI_HISTORY_FILENAME    ".sqlitecloud_history.txt"
#define CLI_VERSION             "1.0a2"
#define CLI_BUILD_DATE          __DATE__

void do_command (SQCloudConnection *conn, char *command) {
    SQCloudResult *res = SQCloudExec(conn, command);
    
    SQCloudResType type = SQCloudResultType(res);
    switch (type) {
        case RESULT_OK:
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

int main(int argc, const char * argv[]) {
    const char *hostname = "localhost";
    int port = SQCLOUD_DEFAULT_PORT;
    
    // a very simple command line parser (atoi not really recommended)
    if (argc > 1) hostname = argv[1];
    if (argc > 2) port = atoi(argv[2]);
    if (hostname == NULL) hostname = "localhost";
    if (port <= 0) port = SQCLOUD_DEFAULT_PORT;
    
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
