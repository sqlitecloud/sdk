//
//  main.c
//  sqlitecloud-cli
//
//  Created by Marco Bambini on 08/02/21.
//

#if defined(WIN32) || defined(_WIN32) || defined(__WIN32__) || defined(__NT__)
#define CLI_WINDOWS             1
#endif

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <stdbool.h>
#include "sqcloud.h"
#include "linenoise.h"

#if CLI_WINDOWS
#include <Windows.h>
#else
// Linux only macro necessary to include non standard functions (like strcasestr)
#define _GNU_SOURCE
#include <sys/fcntl.h>
#endif

#define CLI_HISTORY_FILENAME    ".sqlitecloud_history.txt"
#define CLI_VERSION             "1.2"
#define CLI_BUILD_DATE          __DATE__

#ifndef MAXPATH
#define MAXPATH                 4096
#endif

// MARK: -

bool file_exists (const char *path) {
    #if CLI_WINDOWS
    if (GetFileAttributesA(path) != INVALID_FILE_ATTRIBUTES) return true;
    #else
    if (access(path, F_OK) == 0) return true;
    #endif
    
    return false;
}

int file_create (const char *path) {
    // RW for owner, R for group, R for others
    #if CLI_WINDOWS
    mode_t mode = _S_IWRITE;
    #else
    mode_t mode = S_IRUSR | S_IWUSR | S_IRGRP | S_IROTH;
    #endif
    
    return open(path, O_WRONLY | O_CREAT | O_TRUNC, mode);
}

int file_open_read (const char *path) {
    #if CLI_WINDOWS
    mode_t mode = _S_IREAD;
    #else
    mode_t mode = S_IRUSR | S_IRGRP;
    #endif
    return open(path, O_RDONLY, mode);
}

bool file_delete (const char *path) {
    #if CLI_WINDOWS
    return DeleteFileA(path);
    #else
    return (unlink(path) == 0);
    #endif
}

int64_t file_size (int fd) {
    int64_t fsize = 0;
    #if CLI_WINDOWS
    fsize = (int64_t)_lseek(fd, 0, SEEK_END);
    _lseek(fd, 0, SEEK_SET);
    #else
    fsize = (int64_t)lseek(fd, 0, SEEK_END);
    lseek(fd, 0, SEEK_SET);
    #endif
    return fsize;
}

bool path_combine (char path[MAXPATH], const char dirpath[MAXPATH], const char name[512]) {
    #if CLI_WINDOWS
    return (PathCombineA(path, dirpath, name) != NULL);
    #else
    size_t len = strlen(dirpath);
    int n;
    if ((len) && (dirpath[len-1] != '/')) {
        n = snprintf(path, MAXPATH, "%s%s%s", dirpath, "/", name);
    } else {
        n = snprintf(path, MAXPATH, "%s%s", dirpath, name);
    }
    
    return (n > 0) ? true : false;
    #endif
}

// MARK: -

static bool skip_ok = false;
static bool quiet = false;

static void do_print_usage (void) {
    printf("Usage: sqlitecloud-cli [options]\n");
    printf("Options:\n");
    printf("  -v                    print usage and exit\n");
    printf("  -h HOSTNAME           hostname to connect to (default localhost)\n");
    printf("  -p PORT               port to connect to (default %d)\n", SQCLOUD_DEFAULT_PORT);
    printf("  -f FILEPATH           file path with commands to execute\n");
    printf("  -d DATABASE           database name\n");
    printf("  -s CONNECTIONSTRING   connection string\n");
    printf("  -r ROOT_CERTIFICATE   path to root certificate for TLS connection\n");
    printf("  -t CLI_CERTIFICATE    path to client certificate for TLS connection\n");
    printf("  -k CLI_KEY            path to client key certificate for TLS connection\n");
    printf("  -u TIMEOUT            connection timeout in seconds (default no timeout)\n");
    printf("  -y IP                 connection type (IPv4, IPv6 or IPany, default IPv4)\n");
    printf("  -n USERNAME           authentication username\n");
    printf("  -m PASSWORD           authentication password\n");
    printf("  -c                    activate compression\n");
    printf("  -i                    activate insecure mode (non TLS connection)\n");
    printf("  -j                    disable certificate verification\n");
    printf("  -q                    activate quiet mode (disable output print)\n");
    printf("  -z                    request zero-terminated strings in all replies\n");
    printf("  -w                    in case of -f file to execute, skip the line by line processing and send the whole file\n");
}

bool do_print (SQCloudConnection *conn, SQCloudResult *res) {
    // res NULL means to read error message and error code from conn
    SQCLOUD_RESULT_TYPE type = SQCloudResultType(res);
    bool result = true;
    
    switch (type) {
        case RESULT_OK:
            if (skip_ok) return true;
            printf("OK");
            break;
            
        case RESULT_ERROR:
            printf("ERROR: %s (%d - %d)", SQCloudErrorMsg(conn), SQCloudErrorCode(conn), SQCloudErrorOffset(conn));
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
            SQCloudRowsetDump(res, 0, quiet);
            break;
            
        case RESULT_BLOB:
            printf("BLOB data with len: %d", SQCloudResultLen(res));
            break;
    }
    
    printf("\n\n");
    return result;
}

bool do_command (SQCloudConnection *conn, char *command) {
    SQCloudResult *res = SQCloudExec(conn, command);
    bool result = do_print(conn, res);
    SQCloudResultFree(res);
    return result;
}

bool do_command_without_ok_reply (SQCloudConnection *conn, char *command) {
    skip_ok = true;
    bool result = do_command(conn, command);
    skip_ok = false;
    return result;
}

bool do_internal_command (SQCloudConnection *conn, char *command);

bool do_process_file (SQCloudConnection *conn, const char *filename, bool linebyline) {
    // should continue flag set to false by default
    bool should_continue = false;
    
    FILE *file = fopen(filename, "r");
    if (!file) {
        printf("Unable to open file %s.\n", filename);
        return false;
    }
    
    if (linebyline) {
        char line[512];
        while (fgets(line, sizeof(line), file)) {
            line[strcspn(line, "\n")] = 0;
            if (strcasecmp(line, ".PROMPT")==0) {should_continue = true; break;}
            printf(">> %s\n", line);
            (line[0] == '.') ? do_internal_command(conn, line) : do_command(conn, line);
        }
    } else {
        // get file size
        fseek(file, 0, SEEK_END);
        long size = ftell(file);
        fseek(file, 0, SEEK_SET);
        
        char *buffer = malloc(size + 1);
        if (!buffer) {printf("Unable to allocate %ld buffer size.\n", size); goto cleanup;}
        
        size_t nread = fread(buffer, 1, size, file);
        if (nread != size) {printf("An error occurred while reading file %s (%ld - %zu).\n", filename, size, nread); free(buffer); goto cleanup;}
        buffer[size] = 0;
        
        printf(">> Executing file: %s (%ld bytes)\n\n", filename, size);
        
        do_command(conn, buffer);
        free(buffer);
    }
    
cleanup:
    fclose(file);
    return should_continue;
}

// MARK: -

int do_internal_download_cb (void *xdata, const void *buffer, uint32_t blen, int64_t ntot, int64_t nprogress) {
    if (blen) {
        // retrieve file descriptor
        int fd = ((SQCloudData *)xdata)->fd;

        // write data
        if (write(fd, buffer, (size_t)blen) != blen) {
            printf("\nError while writing data to file.\n");
            return -1;
        }
    }
    
    // display a simple text progress
    printf("%.2f%% ", ((double)nprogress / (double)ntot) * 100.0);
    
    // check if it is final step
    if (ntot == nprogress) printf("\n\n");
    
    // means no error and continue the loop
    return 0;
}

bool do_internal_download (SQCloudConnection *conn, char *command) {
    // .download dbname path
    
    // skip command name part
    command += strlen(".download ");
    
    // extract parameters
    char dbname[512];
    char dbpath[MAXPATH];
    if (sscanf(command, "%s %s", (char *)&dbname, (char *)&dbpath) != 2) {
        return false;
    }
    
    // check if path exists
    if (!file_exists(dbpath)) {
        printf("Output path %s does not exist.\n", dbpath);
        return false;
    }
    
    // generate full path to output database
    char path[MAXPATH];
    path_combine(path, dbpath, dbname);
    
    // create file
    int fd = file_create(path);
    if (fd < 0) {
        printf("Unable to create output file %s\n", path);
        return false;
    }
    
    printf("    ");
    SQCloudData data = {.ptr = NULL, .fd = fd};
    bool result = SQCloudDownloadDatabase(conn, dbname, (void *)&data, do_internal_download_cb);
    if (!result) {
        printf("\n");
        do_print(conn, NULL);
    }
    
    close(fd);
    if (!result) file_delete(path);
    
    return result;
}

int do_internal_read_cb (void *xdata, void *buffer, uint32_t *blen, int64_t ntot, int64_t nprogress) {
    int fd = ((SQCloudData *)xdata)->fd;
    
    ssize_t nread = read(fd, buffer, (size_t)*blen);
    if (nread == -1) return -1;
    
    if (nread == 0) printf("UPLOAD COMPLETE\n\n");
    else printf("%.2f%% ", ((double)(nprogress+nread) / (double)ntot) * 100.0);
    
    *blen = (uint32_t)nread;
    return 0;
}

bool do_internal_upload (SQCloudConnection *conn, char *command) {
    // .upload dbname path [key]
    
    // skip command name part
    command += strlen(".upload ");
    
    // extract parameters
    char *key = NULL;
    char dbkey[512];
    char dbname[512];
    char dbpath[MAXPATH];
    
    // parse command
    int count = sscanf(command, "%s %s %s", (char *)&dbname, (char *)&dbpath, (char *)&dbkey);
    
    // dbname and path parameters are mandatory
    if (count < 2) {
        printf("Database name and path are mandatory in the .upload command\n");
        return false;
    }
    
    // third parameter is optional
    if (count == 3) key = dbkey;
    
    // check if path exists
    if (!file_exists(dbpath)) {
        printf("Database %s does not exist.\n", dbpath);
        return false;
    }
    
    // open file in read-only mode (database must not be in use)
    int fd = file_open_read(dbpath);
    if (fd < 0) {
        printf("Unable to open database file %s\n", dbpath);
        return false;
    }
    
    // get file size (to have a nice progress stat)
    int64_t dbsize = file_size(fd);
    if (dbsize < 0) dbsize = 0;
    
    printf("    ");
    SQCloudData data = {.ptr = NULL, .fd = fd};
    bool result = SQCloudUploadDatabase(conn, dbname, key, (void *)&data, dbsize, do_internal_read_cb);
    close(fd);
    
    if (!result) {
        printf("\n");
        do_print(conn, NULL);
    }
    
    return result;
}

bool do_internal_file (SQCloudConnection *conn, char *command) {
    // .file path_to_file
    
    // skip command name part
    command += strlen(".file ");
    
    return do_process_file(conn, command, false);
}

bool do_internal_prepare (SQCloudConnection *conn, char *command) {
    // .prepare sql
    
    // skip command name part
    command += strlen(".prepare ");
    
    // sanity check
    if (strlen(command) >= 1024) {
        printf("SQL statement too long (%lu, max 1024)\n", strlen(command));
        return false;
    }
    
    // build and send PREPARE statement
    char sql[1024];
    snprintf(sql, sizeof(sql), "SELECT VM_PREPARE('%s');", command);
    
    /*
     [0] 21     (ARRAY_TYPE_PREPARE_VM, see ARRAY_TYPE enum)
     [1] 0      (VM index)
     [2] 0      (number of SQL params)
     [3] 1      (statement is read-only)
     [4] 1      (number of columms in the SQL statement, if any)
     [5] 0      (statement is explain)
     [6] 0      (number of character to skip in original command to get next one, or 0 if none)
     [7] 0      (statement is finalyzed)
     */
    
    return do_command(conn, sql);
}

bool do_internal_step (SQCloudConnection *conn, char *command) {
    // .step vm_index
    
    // skip command name part
    command += strlen(".step ");
    
    // build and send PREPARE statement
    char sql[128];
    snprintf(sql, sizeof(sql), "SELECT VM_STEP(%s);", command);
    
    /*
     [0] 21     (ARRAY_TYPE_SQLITE_EXEC or ARRAY_TYPE_DONE_VM, see ARRAY_TYPE enum)
     [1] 0      (VM index, in ARRAY_TYPE_SQLITE_EXEC is always 0)
     [2] 0      (sqlite3_last_insert_rowid)
     [3] 1      (sqlite3_changes)
     [4] 1      (sqlite3_total_changes)
     [5] 0      (statement is finalyzed)
     */
    
    return do_command(conn, sql);
}

bool do_internal_vmcommand (SQCloudConnection *conn, const char *cname, char *command) {
    // sanity check
    char sql[1024];
    if (strlen(command) >= sizeof(sql)) {
        printf("SQL statement too long (%lu, max 1024)\n", strlen(command));
        return false;
    }
    
    // build statement
    char vmcommand[64];
    if (strcmp(cname, ".prepare") == 0) snprintf(vmcommand, sizeof(vmcommand), "VM_PREPARE");
    else if (strcmp(cname, ".step") == 0) snprintf(vmcommand, sizeof(vmcommand), "VM_STEP");
    else if (strcmp(cname, ".clear") == 0) snprintf(vmcommand, sizeof(vmcommand), "VM_CLEAR");
    else if (strcmp(cname, ".reset") == 0) snprintf(vmcommand, sizeof(vmcommand), "VM_RESET");
    else if (strcmp(cname, ".finalize") == 0) snprintf(vmcommand, sizeof(vmcommand), "VM_FINALIZE");
    
    
    // build and execute statement
    snprintf(sql, sizeof(sql), "SELECT %s('%s');", vmcommand, command);
    return do_command(conn, sql);
}

bool do_internal_command (SQCloudConnection *conn, char *command) {
    // extract command name
    char cname[512];
    sscanf(command, "%s ", (char *)&cname);
    
    if (strcmp(cname, ".download") == 0) return do_internal_download(conn, command);
    if (strcmp(cname, ".upload") == 0) return do_internal_upload(conn, command);
    if (strcmp(cname, ".file") == 0) return do_internal_file(conn, command);
    
    if ((strcmp(cname, ".prepare") == 0) || (strcmp(cname, ".step") == 0) || (strcmp(cname, ".clear") == 0) ||
        (strcmp(cname, ".reset") == 0) || (strcmp(cname, ".finalize") == 0)) {
        command += strlen(cname) + 1; // +1 to take in account the space separator
        return do_internal_vmcommand(conn, cname, command);
    }
    
    printf("Unable to recognize internal command: %s\n", command);
    return false;
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
    const char *username = NULL;
    const char *password = NULL;
    const char *filename = NULL;
    const char *database = NULL;
    const char *connstring = NULL;
    const char *root_certificate_path = NULL;
    const char *client_certificate_path = NULL;
    const char *client_certificate_key_path = NULL;
    
    int family = SQCLOUD_IPv4;
    int port = SQCLOUD_DEFAULT_PORT;
    int timeout = 0;
    
    bool compression = false;
    bool insecure = false;
    bool noverifycert = false;
    bool zerotext = false;
    bool linebyline = true;
    
    int c;
    while ((c = getopt (argc, argv, "h:p:f:vcijqxwzr:s:t:d:y:u:n:m:")) != -1) {
        switch (c) {
            case 'v': do_print_usage(); return 0;
            case 'h': hostname = optarg; break;
            case 'p': port = atoi(optarg); break;
            case 'f': filename = optarg; break;
            case 'c': compression = true; break;
            case 'i': insecure = true; break;
            case 'j': noverifycert = true; break;
            case 'q': quiet = true; break;
            case 'z': zerotext = true; break;
            case 'w': linebyline = false; break;
            case 'd': database = optarg; break;
            case 'r': root_certificate_path = optarg; break;
            case 't': client_certificate_path = optarg; break;
            case 'k': client_certificate_key_path = optarg; break;
            case 'u': timeout = atoi(optarg); break;
            case 'y':
                if (strcasestr(optarg, "IPv6") != 0) {family = SQCLOUD_IPv6;}
                else if (strcasestr(optarg, "IPany") != 0) {family = SQCLOUD_IPANY;}
                break;
            case 'n': username = optarg; break;
            case 'm': password = optarg; break;
            case 's': connstring = optarg; break;
        }
    }
    
    if (!quiet) printf("sqlitecloud-cli version %s (build date %s)\n", CLI_VERSION, CLI_BUILD_DATE);
    
    // sanity check username and password (both are required if one is set)
    if ((username != NULL && password == NULL) || (username == NULL && password != NULL)) {
        printf("Please provide both username and password for authentication.\n");
        return -1;
    }
    
    // setup config
    SQCloudConfig config = {0};
    config.family = family;
    config.timeout = timeout;
    
    // try to connect to hostname:port
    SQCloudConnection *conn = NULL;
    if (connstring) {
        conn = SQCloudConnectWithString(connstring, &config);
    } else {
        // setup TLS config parameter
        #ifndef SQLITECLOUD_DISABLE_TLS
        if (insecure) config.insecure = true;
        if (noverifycert) config.no_verify_certificate = true;
        if (root_certificate_path) config.tls_root_certificate = root_certificate_path;
        if (client_certificate_path) config.tls_certificate = client_certificate_path;
        if (client_certificate_key_path) config.tls_certificate_key = client_certificate_key_path;
        #endif
        
        if (zerotext) config.zero_text = true;
        if (compression) config.compression = true;
        if (database) config.database = database;
        if (username) config.username = username;
        if (password) config.password = password;
        
        conn = SQCloudConnect(hostname, port, &config);
    }
    
    if (SQCloudIsError(conn)) {
        printf("ERROR connecting to %s: %s (%d)\n", hostname, SQCloudErrorMsg(conn), SQCloudErrorCode(conn));
        return -1;
    } else {
        if (!quiet) printf("Connection to %s:%d OK...\n\n", hostname, port);
    }
    
    // load history file
    linenoiseHistoryLoad(CLI_HISTORY_FILENAME);
    
    if (filename) {
        bool should_continue = do_process_file(conn, filename, linebyline);
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
        if (strncmp(command, ".exit", 5) == 0) break;
        (command[0] == '.') ? do_internal_command(conn, command) : do_command(conn, command);
        if (strncmp(command, "QUIT SERVER", 11) == 0) break;
    }
    if (command) free(command);
    
    SQCloudDisconnect(conn);
    return 0;
}
