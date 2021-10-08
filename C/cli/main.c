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

// Linux only macro necessary to include non standard functions (like strcasestr)
#define _GNU_SOURCE

#define CLI_HISTORY_FILENAME    ".sqlitecloud_history.txt"
#define CLI_VERSION             "1.0"
#define CLI_BUILD_DATE          __DATE__

// MARK: -

static bool skip_ok = false;
static bool quiet = false;

void do_print (SQCloudConnection *conn, SQCloudResult *res) {
    // res NULL means to read error message and error code from conn
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
            
        case RESULT_JSON:
        case RESULT_STRING:
        case RESULT_INTEGER:
        case RESULT_FLOAT:
            printf("%.*s", SQCloudResultLen(res), SQCloudResultBuffer(res));
            break;
            
        case RESULT_ARRAY:
            SQCloudArrayDump(res);
            break;
            
        case RESULT_ROWSET:
            SQCloudRowsetDump(res, 0, quiet);
            break;
            
        case RESULT_BLOB:
            printf("BLOB data with len: %d", SQCloudResultLen(res));
            break;
    }
    
    printf("\n\n");
}

void do_command (SQCloudConnection *conn, char *command) {
    SQCloudResult *res = SQCloudExec(conn, command);
    do_print(conn, res);
    SQCloudResultFree(res);
}

void do_command_without_ok_reply (SQCloudConnection *conn, char *command) {
    skip_ok = true;
    do_command(conn, command);
    skip_ok = false;
}

bool do_process_file (SQCloudConnection *conn, const char *filename) {
    // should continue flag set to false by default
    bool should_continue = false;
    
    FILE *file = fopen(filename, "r");
    if (!file) {
        printf("Unable to open file %s.\n", filename);
        return false;
    }
    
    char line[512];
    while (fgets(line, sizeof(line), file)) {
        line[strcspn(line, "\n")] = 0;
        if (strcasecmp(line, ".PROMPT")==0) {should_continue = true; break;}
        printf(">> %s\n", line);
        do_command(conn, line);
    }
    
    fclose(file);
    return should_continue;
}

// MARK: -

void pubsub_callback (SQCloudConnection *connection, SQCloudResult *result, void *data) {
    // THIS CALLBACK IS EXECUTED IN ANOTHER THREAD
    do_print(connection, result);
}

// MARK: -

// usage:
// % sqlitecloud-cli -h HOST -p PORT -f FILE -c -i -r root_certificate -s client_certificate -t client_certificate_key
int main(int argc, char * argv[]) {
    const char *hostname = "localhost";
    const char *filename = NULL;
    const char *database = NULL;
    const char *root_certificate_path = NULL;
    const char *client_certificate_path = NULL;
    const char *client_certificate_key_path = NULL;
    
    int port = SQCLOUD_DEFAULT_PORT;
    bool compression = false;
    bool insecure = false;
    bool sqlite = false;
    bool zerotext = false;
    int family = SQCLOUD_IPv4;
    
    int c;
    while ((c = getopt (argc, argv, "h:p:f:ciqxzr:s:t:d:y:")) != -1) {
        switch (c) {
            case 'h': hostname = optarg; break;
            case 'p': port = atoi(optarg); break;
            case 'f': filename = optarg; break;
            case 'c': compression = true; break;
            case 'i': insecure = true; break;
            case 'q': quiet = true; break;
            case 'x': sqlite = true; break;
            case 'z': zerotext = true; break;
            case 'd': database = optarg; break;
            case 'r': root_certificate_path = optarg; break;
            case 's': client_certificate_path = optarg; break;
            case 't': client_certificate_key_path = optarg; break;
            case 'y':
                if (strcasestr(optarg, "IPv6") != 0) {family = SQCLOUD_IPv6;}
                else if (strcasestr(optarg, "IPany") != 0) {family = SQCLOUD_IPany;}
                break;
        }
    }
    
    if (!quiet) printf("sqlitecloud-cli version %s (build date %s)\n", CLI_VERSION, CLI_BUILD_DATE);
    
    // setup config
    SQCloudConfig config = {0};
    config.family = family;
    
    // setup TLS config parameter
    #ifndef SQLITECLOUD_DISABLE_TSL
    if (insecure) config.insecure = true;
    if (root_certificate_path) config.tls_root_certificate = root_certificate_path;
    if (client_certificate_path) config.tls_certificate = client_certificate_path;
    if (client_certificate_key_path) config.tls_certificate_key = client_certificate_key_path;
    #endif
    
    if (sqlite) config.sqlite_mode = true;
    if (zerotext) config.zero_text = true;
    if (compression) config.compression = true;
    if (database) config.database = database;
    
    // try to connect to hostname:port
    SQCloudConnection *conn = SQCloudConnect(hostname, port, &config);
    if (SQCloudIsError(conn)) {
        printf("ERROR connecting to %s: %s (%d)\n", hostname, SQCloudErrorMsg(conn), SQCloudErrorCode(conn));
        return -1;
    } else {
        printf("Connection to %s OK...\n\n", hostname);
    }
    
    // load history file
    linenoiseHistoryLoad(CLI_HISTORY_FILENAME);
    
    if (filename) {
        bool should_continue = do_process_file(conn, filename);
        if (should_continue == false) return 0;
    }
    
    // SET Pub/Sub callback
    SQCloudSetPubSubCallback(conn, pubsub_callback, NULL);
    
    // REPL
    char *command = NULL;
    while((command = linenoise(">> ")) != NULL) {
        if (command[0] != '\0') {
            linenoiseHistoryAdd(command);
            linenoiseHistorySave(CLI_HISTORY_FILENAME);
        }
        if (strncmp(command, "EXIT", 4) == 0) break;
        do_command(conn, command);
        if (strncmp(command, "QUIT", 4) == 0) break;
    }
    if (command) free(command);
    
    SQCloudDisconnect(conn);
    return 0;
}
