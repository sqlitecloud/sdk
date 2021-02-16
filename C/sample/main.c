//
//  main.c
//  sqlitecloud-cli
//
//  Created by Marco Bambini on 08/02/21.
//

#include <stdio.h>
#include <string.h>
#include "sqcloud.h"

int main(int argc, const char * argv[]) {
    const char *hostname = "localhost";
    if (argc > 1) hostname = argv[1];
    
    SQCloudConnection *conn = SQCloudConnect(hostname, SQCLOUD_DEFAULT_PORT, NULL);
    if (SQCloudIsError(conn)) {
        printf("ERROR connecting to %s: %s (%d)\n", hostname, SQCloudErrorMsg(conn), SQCloudErrorCode(conn));
        return -1;
    } else {
        printf("Connection to %s OK...\n\n", hostname);
    }
    
    // REPL
    while (1) {
        printf(">> ");
        char command[256];
        fgets(command, sizeof(command), stdin);
        size_t len = strlen(command);
        command[len-1] = 0;
        --len;
        
        if (strncmp(command, "QUIT", len) == 0) break;
        SQCloudResult *res = SQCloudExec(conn, command);
        
        SQCloudResType type = SQCloudResultType(res);
        switch (type) {
            case RESULT_OK:
                printf("OK");
                break;
                
            case RESULT_ERROR:
                printf("ERROR: %s (%d)", SQCloudErrorMsg(conn), SQCloudErrorCode(conn));
                break;
                
            case RESULT_STRING:
                printf("%.*s", SQCloudResultLen(res), SQCloudResultBuffer(res));
                break;
                
            case RESULT_ROWSET:
                SQCloudRowSetDump(res);
                break;
        }
        
        printf("\n\n");
        SQCloudResultFree(res);
    }
    
    SQCloudFree(conn);
    return 0;
}
