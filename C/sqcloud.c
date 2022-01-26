//
//  sqcloud.c
//
//  Created by Marco Bambini on 08/02/21.
//

#include "lz4.h"
#include "base64.h"
#include "sqcloud.h"

#include <ctype.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdarg.h>
#include <assert.h>
#include <sys/time.h>

#ifdef _WIN32
#include <winsock2.h>
#include <ws2tcpip.h>
#pragma comment(lib, "Ws2_32.lib")
#include <Shlwapi.h>
#include <io.h>
#include <float.h>
#include "pthread.h"
#else
#include <errno.h>
#include <netdb.h>
#include <signal.h>
#include <unistd.h>
#include <netinet/in.h>
#include <sys/socket.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <arpa/inet.h>
#include <netinet/tcp.h>
#include <sys/ioctl.h>
#include <pthread.h>
#endif

#ifndef SQLITECLOUD_DISABLE_TSL
#include "tls.h"
#endif

// MARK: MACROS -
#ifdef _WIN32
#pragma warning (disable: 4005)
#pragma warning (disable: 4068)
#define readsocket(a,b,c)                   recv((a), (b), (c), 0L)
#define writesocket(a,b,c)                  send((a), (b), (c), 0L)
#else
#define readsocket                          read
#define writesocket                         write
#define closesocket(s)                      close(s)
#endif

#ifndef mem_alloc
#define mem_realloc                         realloc
#define mem_zeroalloc(_s)                   calloc(1,_s)
#define mem_alloc(_s)                       malloc(_s)
#define mem_free(_s)                        free(_s)
#define string_dup(_s)                      strdup(_s)
#endif
#define MIN(a,b)                            (((a)<(b))?(a):(b))

#define MAX_SOCK_LIST                       6           // maximum number of socket descriptor to try to connect to
                                                        // this change is required to support IPv4/IPv6 connections
#define DEFAULT_TIMEOUT                     12          // default connection timeout in seconds

#define REPLY_OK                            "+2 OK"     // default OK reply
#define REPLY_OK_LEN                        5           // default OK reply string length

// https://levelup.gitconnected.com/8-ways-to-measure-execution-time-in-c-c-48634458d0f9
#define TIME_GET(_t1)                       struct timeval _t1; gettimeofday(&_t1, NULL)
#define TIME_VAL(_t1, _t2)                  ((double)(_t2.tv_sec - _t1.tv_sec) + (double)((_t2.tv_usec - _t1.tv_usec)*1e-6))

#define CMD_STRING                          '+'
#define CMD_ZEROSTRING                      '!'
#define CMD_ERROR                           '-'
#define CMD_INT                             ':'
#define CMD_FLOAT                           ','
#define CMD_ROWSET                          '*'
#define CMD_ROWSET_CHUNK                    '/'
#define CMD_JSON                            '#'
#define CMD_RAWJSON                         '{'
#define CMD_NULL                            '_'
#define CMD_BLOB                            '$'
#define CMD_COMPRESSED                      '%'
#define CMD_PUBSUB                          '|'
#define CMD_COMMAND                         '^'
#define CMD_RECONNECT                       '@'
#define CMD_ARRAY                           '='

#define CMD_MINLEN                          2

#define CONNSTRING_KEYVALUE_SEPARATOR       '='
#define CONNSTRING_TOKEN_SEPARATOR          ';'

#define DEFAULT_CHUCK_NBUFFERS              20
#define DEFAULT_CHUNK_MINROWS               2000

#define COMPUTE_BASE64_SIZE(_len)           (((_len + 3 - (_len % 3)) / 3) * 4)

// MARK: - PROTOTYPES -

static SQCloudResult *internal_socket_read (SQCloudConnection *connection, bool mainfd);
static bool internal_socket_write (SQCloudConnection *connection, const char *buffer, size_t len, bool mainfd);
static uint32_t internal_parse_number (char *buffer, uint32_t blen, uint32_t *cstart, uint32_t *extcode);
static SQCloudResult *internal_parse_buffer (SQCloudConnection *connection, char *buffer, uint32_t blen, uint32_t cstart, bool isstatic, bool externalbuffer);
static bool internal_connect (SQCloudConnection *connection, const char *hostname, int port, SQCloudConfig *config, bool mainfd);
static bool internal_set_error (SQCloudConnection *connection, int errcode, const char *format, ...);

// MARK: -

struct SQCloudResult {
    SQCloudResType  tag;                    // RESULT_OK, RESULT_ERROR, RESULT_STRING, RESULT_INTEGER, RESULT_FLOAT, RESULT_ROWSET, RESULT_NULL
    
    bool            ischunk;                // flag used to correctly access the union below
    union {
        struct {
            char        *buffer;            // buffer used by the user (it could be a ptr inside rawbuffer)
            char        *rawbuffer;         // ptr to the buffer to be freed
            uint32_t    balloc;             // buffer allocation size
        };
        struct {
            char        **buffers;          // array of buffers used by rowset sent in chunk
            uint32_t    bcount;             // number of buffers in the array
            uint32_t    bnum;               // number of pre-allocated buffers
            uint32_t    brows;              // number of pre-allocated rows
        };
    };
    
    // common
    uint32_t        blen;                   // total buffer length (also the sum of buffers)
    double          time;                   // full execution time (latency + server side time)
    bool            externalbuffer;         // true if the buffer is managed by the caller code
                                            // false if the buffer can be freed by the SQCloudResultFree func
    uint32_t        nheader;                // number of character in the first part of the header (which is usually skipped)
    
    // used in TYPE_ROWSET only
    uint32_t        nrows;                  // number of rows
    uint32_t        ncols;                  // number of columns
    uint32_t        ndata;                  // number of items stores in data
    char            **data;                 // data contained in the rowset
    char            **name;                 // column names
    char            **decltype;             // column declared types (sqlite mode only)
    char            **dbname;               // column database names (sqlite mode only)
    char            **tblname;              // column table names (sqlite mode only)
    char            **origname;             // column origin names (sqlite mode only)
    uint32_t        *clen;                  // max len for each column (used to display result)
    uint32_t        maxlen;                 // max len for each row/column
} _SQCloudResult;

struct SQCloudConnection {
    int             fd;
    char            errmsg[1024];
    int             errcode;                // error code
    int             xerrcode;               // extended error code
    SQCloudResult   *_chunk;
    SQCloudConfig   *_config;
    
    // pub/sub
    char            *uuid;
    int             pubsubfd;
    SQCloudPubSubCB callback;
    void            *data;
    char            *hostname;
    int             port;
    pthread_t       tid;
    
    #ifndef SQLITECLOUD_DISABLE_TSL
    struct tls      *tls_context;
    struct tls      *tls_pubsub_context;
    #endif
} _SQCloudConnection;

static SQCloudResult SQCloudResultOK = {RESULT_OK, NULL, 0, 0, 0};
static SQCloudResult SQCloudResultNULL = {RESULT_NULL, NULL, 0, 0, 0};

// MARK: - UTILS -

static uint32_t utf8_charbytes (const char *s, uint32_t i) {
    unsigned char c = (unsigned char)s[i];
    
    // determine bytes needed for character, based on RFC 3629
    if ((c > 0) && (c <= 127)) return 1;
    if ((c >= 194) && (c <= 223)) return 2;
    if ((c >= 224) && (c <= 239)) return 3;
    if ((c >= 240) && (c <= 244)) return 4;
    
    // means error
    return 0;
}

static uint32_t utf8_len (const char *s, uint32_t nbytes) {
    uint32_t pos = 0;
    uint32_t len = 0;
    
    while (pos < nbytes) {
        ++len;
        uint32_t n = utf8_charbytes(s, pos);
        if (n == 0) return 0; // means error
        pos += n;
    }
    
    return len;
}

#if 0
static char *extract_connection_token (const char *s, char *key, char buffer[256]) {
    char *target = strstr(s, key);
    if (!target) return NULL;
    
    // find out = separator
    char *p = target;
    while (p[0]) {
        if (p[0] == CONNSTRING_KEYVALUE_SEPARATOR) break;
        ++p;
    }
    
    // skip =
    ++p;
    
    // skip spaces (if any)
    while (p[0]) {
        if (!isspace(p[0])) break;
        ++p;
    }
    
    // copy value to buffer
    int len = 0;
    while (p[0] && len < 255) {
        if (isspace(p[0])) break;
        if (p[0] == CONNSTRING_TOKEN_SEPARATOR) break;
        buffer[len] = p[0];
        ++len;
        ++p;
    }
    
    // null-terminate returning value
    buffer[len] = 0;
    p = &buffer[0];
    
    return p;
}
#endif

// MARK: - PRIVATE -

static int socket_geterror (int fd) {
    int err;
    socklen_t errlen = sizeof(err);
    
    int sockerr = getsockopt(fd, SOL_SOCKET, SO_ERROR, (void *)&err, &errlen);
    if (sockerr < 0) return -1;
    
    return ((err == 0 || err == EINTR || err == EAGAIN || err == EINPROGRESS)) ? 0 : err;
}

static void *pubsub_thread (void *arg) {
    SQCloudConnection *connection = (SQCloudConnection *)arg;
    
    int fd = connection->pubsubfd;
    #ifndef SQLITECLOUD_DISABLE_TSL
    struct tls *tls = connection->tls_pubsub_context;
    #endif
    
    size_t blen = 2048;
    char *buffer = mem_alloc(blen);
    if (buffer == NULL) return NULL;
    
    char *original = buffer;
    uint32_t tread = 0;

    while (1) {
        fd_set set;
        FD_ZERO(&set);
        FD_SET(fd, &set);
        
        // wait for read event
        int rc = select(fd + 1, &set, NULL, NULL, NULL);
        if (rc <= 0) continue;
        
        //  read payload string
        #ifndef SQLITECLOUD_DISABLE_TSL
        ssize_t nread = (tls) ? tls_read(tls, buffer, blen) : readsocket(fd, buffer, blen);
        if ((tls) && (nread == TLS_WANT_POLLIN || nread == TLS_WANT_POLLOUT)) continue;
        #else
        ssize_t nread = readsocket(fd, buffer, blen);
        #endif
        
        if (nread < 0) {
            const char *msg = "";
            #ifndef SQLITECLOUD_DISABLE_TSL
            if (tls) msg = tls_error(tls);
            #endif
            
            internal_set_error(connection, INTERNAL_ERRCODE_NETWORK, "An error occurred while reading data: %s (%s).", strerror(errno), msg);
            connection->callback(connection, NULL, connection->data);
            break;
        }
        
        if (nread == 0) {
            const char *msg = "";
            #ifndef SQLITECLOUD_DISABLE_TSL
            if (tls) msg = tls_error(tls);
            #endif
            
            internal_set_error(connection, INTERNAL_ERRCODE_NETWORK, "An error occurred while reading data: %s (%s).", strerror(errno), msg);
            connection->callback(connection, NULL, connection->data);
            break;
        }
        
        tread += (uint32_t)nread;
        blen -= (uint32_t)nread;
        buffer += nread;
        
        uint32_t cstart = 0;
        uint32_t clen = internal_parse_number (&original[1], tread-1, &cstart, NULL);
        if (clen == 0) continue;
        
        // check if read is complete
        // clen is the lenght parsed in the buffer
        // cstart is the index of the first space
        // +1 because we skipped the first character in the internal_parse_number function
        if (clen + cstart + 1 != tread) {
            // check buffer allocation and continue reading
            if (clen + cstart > blen) {
                char *clone = mem_alloc(clen + cstart + 1);
                if (!clone) {
                    internal_set_error(connection, INTERNAL_ERRCODE_MEMORY, "Unable to allocate memory: %d.", clen + cstart + 1);
                    connection->callback(connection, NULL, connection->data);
                    break;
                }
                memcpy(clone, original, tread);
                buffer = original = clone;
                blen = (clen + cstart + 1) - tread;
                buffer += tread;
            }
            
            continue;
        }
        
        SQCloudResult *result = internal_parse_buffer(connection, original, tread, (clen) ? cstart : 0, false, false);
        if (result->tag == RESULT_STRING) result->tag = RESULT_JSON;
        
        connection->callback(connection, result, connection->data);
        
        blen = 2048;
        buffer = mem_alloc(blen);
        if (!buffer) break;
        
        original = buffer;
        tread = 0;
    }
    
    return NULL;
}

// MARK: -

static bool internal_init (void) {
    static bool inited = false;
    if (inited) return true;
    
    #ifdef _WIN32
    WSADATA wsaData;
    WSAStartup(MAKEWORD(2,2), &wsaData);
    #else
    // IGNORE SIGPIPE and SIGABORT
    struct sigaction act;
    act.sa_handler = SIG_IGN;
    sigemptyset(&act.sa_mask);
    act.sa_flags = 0;
    sigaction(SIGPIPE, &act, (struct sigaction *)NULL);
    sigaction(SIGABRT, &act, (struct sigaction *)NULL);
    #endif
    
    inited = true;
    return true;
}

static bool internal_set_error (SQCloudConnection *connection, int errcode, const char *format, ...) {
    connection->errcode = errcode;
    
    va_list arg;
    va_start (arg, format);
    vsnprintf(connection->errmsg, sizeof(connection->errmsg), format, arg);
    va_end (arg);
    
    return false;
}

static void internal_parse_uuid (SQCloudConnection *connection, const char *buffer, size_t blen) {
    // sanity check
    if (!buffer || blen == 0) return;
    
    // expected buffer is PAUTH uuid secret
    // PUATH -> 5
    // uuid -> 36
    // secret -> 36
    // spaces -> 2
    if (blen != (5 + 36 + 36 + 2)) return;
    
    if (strncmp(buffer, "PAUTH ", 6) != 0) return;
    
    // allocate 36 (UUID) + 1 (null-terminated) zero-bytes
    char *uuid = mem_zeroalloc(37);
    if (!uuid) return;
    
    memcpy(uuid, &buffer[6], 36);
    connection->uuid = uuid;
}

static void internal_clear_error (SQCloudConnection *connection) {
    connection->errcode = 0;
    connection->xerrcode = 0;
    connection->errmsg[0] = 0;
}

static bool internal_setup_tls (SQCloudConnection *connection, SQCloudConfig *config, bool mainfd) {
    #ifndef SQLITECLOUD_DISABLE_TSL
    if (config && config->insecure) return true;
    
    int rc = 0;
    
    if (tls_init() < 0) {
        return internal_set_error(connection, INTERNAL_ERRCODE_TLS, "Error while initializing TLS library.");
    }
    
    struct tls_config *tls_conf = tls_config_new();
    if (!tls_conf) {
        return internal_set_error(connection, INTERNAL_ERRCODE_TLS, "Error while initializing a new TLS configuration.");
    }
    
    // loads a file containing the root certificates
    if (config && config->tls_root_certificate) {
        rc = tls_config_set_ca_file(tls_conf, config->tls_root_certificate);
        if (rc < 0) {internal_set_error(connection, INTERNAL_ERRCODE_TLS, "Error in tls_config_set_ca_file: %s.", tls_config_error(tls_conf));}
    }
    
    // loads a file containing the server certificate
    if (config && config->tls_certificate) {
        rc = tls_config_set_cert_file(tls_conf, config->tls_certificate);
        if (rc < 0) {internal_set_error(connection, INTERNAL_ERRCODE_TLS, "Error in tls_config_set_cert_file: %s.", tls_config_error(tls_conf));}
    }
    
    // loads a file containing the private key
    if (config && config->tls_certificate_key) {
        rc = tls_config_set_key_file(tls_conf, config->tls_certificate_key);
        if (rc < 0) {internal_set_error(connection, INTERNAL_ERRCODE_TLS, "Error in tls_config_set_key_file: %s.", tls_config_error(tls_conf));}
    }
    
    struct tls *tls_context = tls_client();
    if (!tls_context) {
        return internal_set_error(connection, INTERNAL_ERRCODE_TLS, "Error while initializing a new TLS client.");
    }
    
    // apply configuration to context
    rc = tls_configure(tls_context, tls_conf);
    if (rc < 0) {
        return internal_set_error(connection, INTERNAL_ERRCODE_TLS, "Error in tls_configure: %s.", tls_error(tls_context));
    }
    
    // save context
    if (mainfd) connection->tls_context = tls_context;
    else connection->tls_pubsub_context = tls_context;
    
    #endif
    return true;
}

static SQCloudValueType internal_type (char *buffer) {
    // for VALUE_NULL values we don't return _ but the NULL value itself, so check for this special case
    // internal_parse_value is used both internally to set the value inside a Rowset (1)
    // and also from the public API SQCloudRowsetValue (2)
    // to fix this misbehaviour, (1) should return _ while (2) should return NULL
    // this is really just a convention so it is much more easier to just return NULL everytime
    if (!buffer) return VALUE_NULL;
    
    switch (buffer[0]) {
        case '+': return VALUE_TEXT;
        case ':': return VALUE_INTEGER;
        case ',': return VALUE_FLOAT;
        case '_': return VALUE_NULL;
        case '$': return VALUE_BLOB;
    }
    return VALUE_NULL;
}

static bool internal_has_commandlen (int c) {
    return ((c == CMD_INT) || (c == CMD_FLOAT) || (c == CMD_NULL)) ? false : true;
}

static uint32_t internal_parse_number (char *buffer, uint32_t blen, uint32_t *cstart, uint32_t *extcode) {
    uint32_t value = 0;
    uint32_t extvalue = 0;
    bool isext = false;
    
    for (uint32_t i=0; i<blen; ++i) {
        int c = buffer[i];
        
        // check for optional extended error code (ERRCODE:EXTERRCODE)
        if (c == ':') {
            isext = true;
            continue;
        }
        
        // check for end of value
        if (c == ' ') {
            *cstart = i+1;
            if (extcode) *extcode = extvalue;
            return value;
        }
        
        // compute numeric value
        if (isext) extvalue = (extvalue * 10) + (buffer[i] - '0');
        else value = (value * 10) + (buffer[i] - '0');
    }
    
    return 0;
}

static char *internal_parse_value (char *buffer, uint32_t *len, uint32_t *cellsize) {
    // handle special NULL value case
    if (!buffer || buffer[0] == CMD_NULL) {
        *len = 0;
        if (cellsize) *cellsize = 2;
        return NULL;
    }
    
    // blen originally was hard coded to 24 because the max 64bit value is 20 characters long
    uint32_t cstart = 0;
    uint32_t blen = *len;
    blen = internal_parse_number(&buffer[1], blen, &cstart, NULL);
    
    // handle decimal/float cases
    if ((buffer[0] == CMD_INT) || (buffer[0] == CMD_FLOAT)) {
        *len = cstart - 1;
        if (cellsize) *cellsize = cstart + 1;
        return &buffer[1];
    }
    
    *len = (buffer[0] == CMD_ZEROSTRING) ? blen - 1 : blen;
    if (cellsize) *cellsize = cstart + blen + 1;
    return &buffer[1+cstart];
}

static SQCloudResult *internal_run_command (SQCloudConnection *connection, const char *buffer, size_t blen, bool mainfd) {
    internal_clear_error(connection);
    
    if (!buffer || blen < CMD_MINLEN) return NULL;
    
    TIME_GET(tstart);
    if (!internal_socket_write(connection, buffer, blen, mainfd)) return false;
    SQCloudResult *result = internal_socket_read(connection, mainfd);
    TIME_GET(tend);
    if (result) result->time = TIME_VAL(tstart, tend);
    return result;
}

static SQCloudResult *internal_setup_pubsub (SQCloudConnection *connection, const char *buffer, size_t blen) {
    // check if pubsub was already setup
    if (connection->pubsubfd != 0) return &SQCloudResultOK;
    
    #ifndef SQLITECLOUD_DISABLE_TSL
    if (!internal_setup_tls(connection, connection->_config, false)) return NULL;
    #endif
    
    if (internal_connect(connection, connection->hostname, connection->port, connection->_config, false)) {
        SQCloudResult *result = internal_run_command(connection, buffer, blen, false);
        if (!SQCloudResultIsOK(result)) return result;
        internal_parse_uuid(connection, buffer, blen);
        pthread_create(&connection->tid, NULL, pubsub_thread, (void *)connection);
    } else {
        return NULL;
    }
    
    return &SQCloudResultOK;
}

static SQCloudResult *internal_reconnect (SQCloudConnection *connection, const char *buffer, size_t blen) {
    // DO RE-CONNECT HERE
    return NULL;
}

static SQCloudResult *internal_parse_array (SQCloudConnection *connection, char *buffer, uint32_t blen, uint32_t bstart) {
    
    SQCloudResult *rowset = (SQCloudResult *)mem_zeroalloc(sizeof(SQCloudResult));
    if (!rowset) {
        internal_set_error(connection, INTERNAL_ERRCODE_MEMORY, "Unable to allocate memory for SQCloudResult: %d.", sizeof(SQCloudResult));
        return NULL;
    }
    
    rowset->tag = RESULT_ARRAY;
    rowset->rawbuffer = buffer;
    rowset->blen = blen;
    rowset->nheader = bstart;
    
    // =LEN N VALUE1 VALUE2 ... VALUEN
    uint32_t start1 = 0;
    uint32_t n = internal_parse_number(&buffer[bstart], blen-1, &start1, NULL);
    
    rowset->ndata = n;
    rowset->data = (char **) mem_alloc(rowset->ndata * sizeof(char *));
    if (!rowset->data) {
        internal_set_error(connection, INTERNAL_ERRCODE_MEMORY, "Unable to allocate memory for SQCloudResult: %d.", rowset->ndata * sizeof(char *));
        mem_free(rowset);
        return NULL;
    }
    
    // loop from i to n to parse each data
    buffer += bstart + start1;
    for (uint32_t i=0; i<n; ++i) {
        uint32_t len = blen, cellsize;
        char *value = internal_parse_value(buffer, &len, &cellsize);
        rowset->data[i] = (value) ? buffer : NULL;
        buffer += cellsize;
        blen -= cellsize;
    }
    
    return rowset;
}

static SQCloudResult *internal_rowset_type (SQCloudConnection *connection, char *buffer, uint32_t blen, uint32_t bstart, SQCloudResType type) {
    SQCloudResult *rowset = (SQCloudResult *)mem_zeroalloc(sizeof(SQCloudResult));
    if (!rowset) {
        internal_set_error(connection, INTERNAL_ERRCODE_MEMORY, "Unable to allocate memory for SQCloudResult: %d.", sizeof(SQCloudResult));
        return NULL;
    }
    
    rowset->tag = type;
    rowset->buffer = &buffer[bstart];
    rowset->rawbuffer = buffer;
    rowset->blen = blen;
    rowset->balloc = blen;
    
    return rowset;
}

static bool internal_parse_rowset_header (SQCloudResult *rowset, char **pbuffer, uint32_t *pblen, uint32_t ncols, bool is_sqlite) {
    char *buffer = *pbuffer;
    uint32_t blen = *pblen;
    
    // header is guarantee to contains column names (1st)
    for (uint32_t i=0; i<ncols; ++i) {
        uint32_t cstart = 0;
        uint32_t len = internal_parse_number(&buffer[1], blen, &cstart, NULL);
        rowset->name[i] = buffer;
        buffer += cstart + len + 1;
        blen -= cstart + len + 1;
        if (rowset->clen[i] < len) rowset->clen[i] = len;
        if (rowset->maxlen < len) rowset->maxlen = len;
    }
    
    if (is_sqlite) {
        rowset->decltype = (char **) mem_alloc(ncols * sizeof(char *));
        if (!rowset->decltype) return false;
        rowset->dbname = (char **) mem_alloc(ncols * sizeof(char *));
        if (!rowset->dbname) return false;
        rowset->tblname = (char **) mem_alloc(ncols * sizeof(char *));
        if (!rowset->tblname) return false;
        rowset->origname = (char **) mem_alloc(ncols * sizeof(char *));
        if (!rowset->origname) return false;
        
        // in sqlite mode header contains column declared types (2nd)
        for (uint32_t i=0; i<ncols; ++i) {
            uint32_t cstart = 0;
            uint32_t len = internal_parse_number(&buffer[1], blen, &cstart, NULL);
            rowset->decltype[i] = buffer;
            buffer += cstart + len + 1;
            blen -= cstart + len + 1;
        }
        
        // in sqlite mode header contains column database names (3rd)
        for (uint32_t i=0; i<ncols; ++i) {
            uint32_t cstart = 0;
            uint32_t len = internal_parse_number(&buffer[1], blen, &cstart, NULL);
            rowset->dbname[i] = buffer;
            buffer += cstart + len + 1;
            blen -= cstart + len + 1;
        }
        
        // in sqlite mode header contains column table names (4th)
        for (uint32_t i=0; i<ncols; ++i) {
            uint32_t cstart = 0;
            uint32_t len = internal_parse_number(&buffer[1], blen, &cstart, NULL);
            rowset->tblname[i] = buffer;
            buffer += cstart + len + 1;
            blen -= cstart + len + 1;
        }
        
        // in sqlite mode header contains column origin names (5th)
        for (uint32_t i=0; i<ncols; ++i) {
            uint32_t cstart = 0;
            uint32_t len = internal_parse_number(&buffer[1], blen, &cstart, NULL);
            rowset->origname[i] = buffer;
            buffer += cstart + len + 1;
            blen -= cstart + len + 1;
        }
    }
    
    *pbuffer = buffer;
    *pblen = blen;
    
    return true;
}

static bool internal_parse_rowset_values (SQCloudResult *rowset, char **pbuffer, uint32_t *pblen, uint32_t index, uint32_t bound, uint32_t ncols) {
    char *buffer = *pbuffer;
    uint32_t blen = *pblen;
    
    for (uint32_t i=index; i<bound; ++i) {
        uint32_t len = blen, cellsize;
        char *value = internal_parse_value(buffer, &len, &cellsize);
        rowset->data[i] = (value) ? buffer : NULL;
        buffer += cellsize;
        blen -= cellsize;
        ++rowset->ndata;
        if (rowset->clen[i % ncols] < len) rowset->clen[i % ncols] = len;
        if (rowset->maxlen < len) rowset->maxlen = len;
    }
    
    return true;
}

static SQCloudResult *internal_parse_rowset (SQCloudConnection *connection, char *buffer, uint32_t blen, uint32_t bstart, uint32_t nrows, uint32_t ncols) {
    
    SQCloudResult *rowset = (SQCloudResult *)mem_zeroalloc(sizeof(SQCloudResult));
    if (!rowset) {
        internal_set_error(connection, INTERNAL_ERRCODE_MEMORY, "Unable to allocate memory for SQCloudResult: %d.", sizeof(SQCloudResult));
        return NULL;
    }
    
    rowset->tag = RESULT_ROWSET;
    rowset->buffer = buffer;
    rowset->rawbuffer = buffer;
    rowset->blen = blen;
    rowset->balloc = blen;
    rowset->nheader = bstart;
    
    rowset->nrows = nrows;
    rowset->ncols = ncols;
    rowset->data = (char **) mem_alloc(nrows * ncols * sizeof(char *));
    rowset->name = (char **) mem_alloc(ncols * sizeof(char *));
    rowset->clen = (uint32_t *) mem_zeroalloc(ncols * sizeof(uint32_t));
    if (!rowset->data || !rowset->name || !rowset->clen) goto abort_rowset;
    
    buffer += bstart;
    
    // parse rowset header
    bool is_sqlite = (connection->_config) && (connection->_config->sqlite_mode);
    if (!internal_parse_rowset_header(rowset, &buffer, &blen, ncols, is_sqlite)) goto abort_rowset;
    
    // parse values
    if (!internal_parse_rowset_values(rowset, &buffer, &blen, 0, nrows * ncols, ncols)) goto abort_rowset;
    
    return rowset;
    
abort_rowset:
    if (rowset->data) mem_free(rowset->data);
    if (rowset->name) mem_free(rowset->name);
    if (rowset->clen) mem_free(rowset->clen);
    if (rowset) mem_free(rowset);
    
    internal_set_error(connection, INTERNAL_ERRCODE_MEMORY, "Unable to allocate internal memory for SQCloudResult.");
    return NULL;
}

static SQCloudResult *internal_parse_rowset_chunck (SQCloudConnection *connection, char *buffer, uint32_t blen, uint32_t bstart, uint32_t idx, uint32_t nrows, uint32_t ncols) {
    SQCloudResult *rowset = connection->_chunk;
    bool first_chunk = false;
    
    // sanity check
    if (idx == 1 && connection->_chunk) {
        // something bad happened here because a first chunk is received while a saved one has not been fully processed
        // lets try to restart the whole process
        SQCloudResultFree(connection->_chunk);
        connection->_chunk = NULL;
        rowset = NULL;
    }
    
    if (!rowset) {
        // this should never happen
        if (idx != 1) return NULL;
        
        // allocate a new rowset
        rowset = (SQCloudResult *)mem_zeroalloc(sizeof(SQCloudResult));
        if (!rowset) {
            internal_set_error(connection, INTERNAL_ERRCODE_MEMORY, "Unable to allocate memory for SQCloudResult: %d.", sizeof(SQCloudResult));
            return NULL;
        }
        first_chunk = true;
        connection->_chunk = rowset;
    }
    
    if (first_chunk) {
        rowset->tag = RESULT_ROWSET;
        rowset->ischunk = true;
        
        rowset->buffers = (char **)mem_zeroalloc((sizeof(char *) * DEFAULT_CHUCK_NBUFFERS));
        if (!rowset->buffers) goto abort_rowset;
        
        rowset->bnum = DEFAULT_CHUCK_NBUFFERS;
        rowset->buffers[0] = buffer;
        rowset->bcount = 1;
        rowset->nheader = bstart;
        
        rowset->brows = nrows + DEFAULT_CHUNK_MINROWS;
        rowset->nrows = nrows;
        rowset->ncols = ncols;
        rowset->data = (char **) mem_alloc(rowset->brows * ncols * sizeof(char *));
        rowset->name = (char **) mem_alloc(ncols * sizeof(char *));
        rowset->clen = (uint32_t *) mem_zeroalloc(ncols * sizeof(uint32_t));
        if (!rowset->data || !rowset->name || !rowset->clen) goto abort_rowset;
        
        buffer += bstart;
        
        // parse rowset header
        bool is_sqlite = (connection->_config) && (connection->_config->sqlite_mode);
        if (!internal_parse_rowset_header(rowset, &buffer, &blen, ncols, is_sqlite)) goto abort_rowset;
    }
    
    // update total buffer size
    rowset->blen += blen;
    
    // check end-chunk condition
    if (idx == 0 && nrows == 0 && ncols == 0) {
        connection->_chunk = NULL;
        if (!rowset->externalbuffer) mem_free(buffer);
        return rowset;
    }
    
    // check if a resize is needed in the array of buffers
    if (rowset->bnum <= rowset->bcount + 1) {
        uint32_t n = rowset->bnum * 2;
        char **temp = (char **)mem_realloc(rowset->buffers, (sizeof(char *) * n));
        if (!temp) goto abort_rowset;
        rowset->buffers = temp;
        rowset->bnum = n;
    }
    
    // check if a resize is needed in the ptr data array
    if (rowset->brows <= rowset->nrows + nrows) {
        uint32_t n = rowset->brows * 2;
        char **temp = (char **)mem_realloc(rowset->data, n * ncols * (sizeof(char *)));
        if (!temp) goto abort_rowset;
        rowset->data = temp;
        rowset->brows = n;
    }
    
    // adjust internal fields
    if (!first_chunk) {
        rowset->buffers[rowset->bcount++] = buffer;
        rowset->nrows += nrows;
        buffer += bstart;
    }
    
    // parse values
    uint32_t index = rowset->ndata;
    uint32_t bound = rowset->ndata + (nrows * ncols);
    
    // parse values
    if (!internal_parse_rowset_values(rowset, &buffer, &blen, index, bound, ncols)) goto abort_rowset;
    
    // this check is for internal usage only
    if (connection->fd == 0) return rowset;
    
    // normal usage
    // send ACK
    if (!internal_socket_write(connection, "OK", 2, true)) goto abort_rowset;
        
    // read next chunk
    return internal_socket_read (connection, true);
    
abort_rowset:
    SQCloudResultFree(rowset);
    connection->_chunk = NULL;
    return NULL;
}

static SQCloudResult *internal_parse_buffer (SQCloudConnection *connection, char *buffer, uint32_t blen, uint32_t cstart, bool isstatic, bool externalbuffer) {
    if (blen <= 1) return false;
    
    // try to check if it is a OK reply: +2 OK
    if ((blen == REPLY_OK_LEN) && (strncmp(buffer, REPLY_OK, REPLY_OK_LEN) == 0)) {
        return &SQCloudResultOK;
    }
    
    // if buffer is static (stack based allocation) then it must be duplicated
    if (buffer[0] != CMD_ERROR && isstatic) {
        char *clone = mem_alloc(blen);
        if (!clone) {
            internal_set_error(connection, INTERNAL_ERRCODE_MEMORY, "Unable to allocate memory: %d.", blen);
            return NULL;
        }
        memcpy(clone, buffer, blen);
        buffer = clone;
        isstatic = false;
    }
    
    // check for compressed reply before the parse step
    char *zdata = NULL;
    if (buffer[0] == CMD_COMPRESSED) {
        // %TLEN CLEN ULEN *0 NROWS NCOLS DATA
        uint32_t cstart1 = 0;
        uint32_t cstart2 = 0;
        uint32_t cstart3 = 0;
        uint32_t tlen = internal_parse_number(&buffer[1], blen-1, &cstart1, NULL);
        uint32_t clen = internal_parse_number(&buffer[cstart1 + 1], blen-1, &cstart2, NULL);
        uint32_t ulen = internal_parse_number(&buffer[cstart1 + cstart2 + 1], blen-1, &cstart3, NULL);
        
        // start of compressed buffer
        zdata = &buffer[tlen - clen + cstart1 + 1];
        
        // start of raw uncompressed header
        char *hstart = &buffer[cstart1 + cstart2 + cstart3 + 1];
        
        // try to allocate a buffer big enough to hold uncompressed data + raw header
        long clonelen = ulen + (zdata - hstart) + 1;
        char *clone = mem_alloc (clonelen);
        if (!clone) {
            internal_set_error(connection, INTERNAL_ERRCODE_MEMORY, "Unable to allocate memory to uncompress buffer: %d.", clonelen);
            if (!isstatic && !externalbuffer) mem_free(buffer);
            return NULL;
        }
        
        // copy raw buffer
        memcpy(clone, hstart, zdata - hstart);
        
        // uncompress buffer and sanity check the result
        uint32_t rc = LZ4_decompress_safe(zdata, clone + (zdata - hstart), clen, ulen);
        if (rc <= 0 || rc != ulen) {
            internal_set_error(connection, INTERNAL_ERRCODE_GENERIC, "Unable to decompress buffer (err code: %d).", rc);
            if (!isstatic && !externalbuffer) mem_free(buffer);
            return NULL;
        }
        
        // decompression is OK so replace buffer
        if (!isstatic && !externalbuffer) mem_free(buffer);
        
        isstatic = false;
        buffer = clone;
        blen = ulen;
        
        // at this point the buffer used in the SQCloudResult is a newly allocated one (clone)
        // so externalbuffer flag must be set to false
        externalbuffer = false;
    }
    
    // parse reply
    switch (buffer[0]) {
        case CMD_ZEROSTRING:
        case CMD_RECONNECT:
        case CMD_PUBSUB:
        case CMD_COMMAND:
        case CMD_STRING:
        case CMD_ARRAY:
        case CMD_BLOB:
        case CMD_JSON: {
            // +LEN string
            uint32_t cstart = 0;
            uint32_t len = internal_parse_number(&buffer[1], blen-1, &cstart, NULL);
            SQCloudResType type = (buffer[0] == CMD_JSON) ? RESULT_JSON : RESULT_STRING;
            if (buffer[0] == CMD_ZEROSTRING) --len;
            else if (buffer[0] == CMD_COMMAND) return internal_run_command(connection, &buffer[cstart+1], len, true);
            else if (buffer[0] == CMD_PUBSUB) return internal_setup_pubsub(connection, &buffer[cstart+1], len);
            else if (buffer[0] == CMD_RECONNECT) return internal_reconnect(connection, &buffer[cstart+1], len);
            else if (buffer[0] == CMD_ARRAY) return internal_parse_array(connection, buffer, len, cstart+1);
            else if (buffer[0] == CMD_BLOB) type = RESULT_BLOB;
            SQCloudResult *res = internal_rowset_type(connection, buffer, len, cstart+1, type);
            if (res) res->externalbuffer = externalbuffer;
            return res;
        }
            
        case CMD_ERROR: {
            // -LEN ERRCODE ERRMSG
            uint32_t cstart = 0, cstart2 = 0;
            uint32_t len = internal_parse_number(&buffer[1], blen-1, &cstart, NULL);
            
            uint32_t excode = 0;
            uint32_t errcode = internal_parse_number(&buffer[cstart + 1], blen-1, &cstart2, &excode);
            connection->errcode = (int)errcode;
            connection->xerrcode = (int)excode;
            
            len -= cstart2;
            memcpy(connection->errmsg, &buffer[cstart + cstart2 + 1], MIN(len, sizeof(connection->errmsg)));
            connection->errmsg[len] = 0;
            
            // check free buffer
            if (!isstatic && !externalbuffer) mem_free(buffer);
            return NULL;
        }
        
        case CMD_ROWSET:
        case CMD_ROWSET_CHUNK: {
            // CMD_ROWSET:          *LEN ROWS COLS DATA
            // CMD_ROWSET_CHUNK:    /LEN IDX ROWS COLS DATA
            uint32_t cstart1 = 0, cstart2 = 0, cstart3 = 0, cstart4 = 0;
            
            internal_parse_number(&buffer[1], blen-1, &cstart1, NULL); // parse len (already parsed in blen parameter)
            uint32_t idx = (buffer[0] == CMD_ROWSET) ? 0 : internal_parse_number(&buffer[cstart1 + 1], blen-1, &cstart2, NULL);
            uint32_t nrows = internal_parse_number(&buffer[cstart1 + cstart2 + 1], blen-1, &cstart3, NULL);
            uint32_t ncols = internal_parse_number(&buffer[cstart1 + cstart2 + + cstart3 + 1], blen-1, &cstart4, NULL);
            
            uint32_t bstart = cstart1 + cstart2 + cstart3 + cstart4 + 1;
            SQCloudResult *res = NULL;
            if (buffer[0] == CMD_ROWSET) res = internal_parse_rowset(connection, buffer, blen, bstart, nrows, ncols);
            else res = internal_parse_rowset_chunck(connection, buffer, blen, bstart, idx, nrows, ncols);
            if (res) res->externalbuffer = externalbuffer;
            
            // check free buffer
            if (!res && !isstatic && !externalbuffer) mem_free(buffer);
            return res;
        }
        
        case CMD_NULL:
            if (!isstatic && !externalbuffer) mem_free(buffer);
            return &SQCloudResultNULL;
            
        case CMD_INT:
        case CMD_FLOAT: {
            // NUMBER case
            internal_parse_value(buffer, &blen, NULL);
            SQCloudResult *res = internal_rowset_type(connection, buffer, blen, 1, (buffer[0] == CMD_INT) ? RESULT_INTEGER : RESULT_FLOAT);
            if (res) res->externalbuffer = externalbuffer;
            
            if (!res && !isstatic && !externalbuffer) mem_free(buffer);
            return res;
        }
            
        case CMD_RAWJSON: {
            // handle JSON here
            // a JSON parser must process raw buffer
            return &SQCloudResultNULL;
        }
    }
    
    if (!isstatic && !externalbuffer) mem_free(buffer);
    return NULL;
}

static bool internal_socket_forward_read (SQCloudConnection *connection, bool (*forward_cb) (char *buffer, size_t blen, void *xdata, void *xdata2), void *xdata, void *xdata2) {
    char sbuffer[8129];
    uint32_t blen = sizeof(sbuffer);
    uint32_t cstart = 0;
    uint32_t tread = 0;
    uint32_t clen = 0;
    
    char *buffer = sbuffer;
    char *original = buffer;
    int fd = connection->fd;
    #ifndef SQLITECLOUD_DISABLE_TSL
    struct tls *tls = connection->tls_context;
    #endif
    
    while (1) {
        // perform read operation
        #ifndef SQLITECLOUD_DISABLE_TSL
        ssize_t nread = (tls) ? tls_read(tls, buffer, blen) : readsocket(fd, buffer, blen);
        if ((tls) && (nread == TLS_WANT_POLLIN || nread == TLS_WANT_POLLOUT)) continue;
        #else
        ssize_t nread = readsocket(fd, buffer, blen);
        #endif
        
        // sanity check read
        if (nread < 0) {
            const char *msg = "";
            #ifndef SQLITECLOUD_DISABLE_TSL
            if (tls) msg = tls_error(tls);
            #endif
            
            internal_set_error(connection, INTERNAL_ERRCODE_NETWORK, "An error occurred while reading data: %s (%s).", strerror(errno), msg);
            goto abort_read;
        }
        
        if (nread == 0) {
            const char *msg = "";
            #ifndef SQLITECLOUD_DISABLE_TSL
            if (tls) msg = tls_error(tls);
            #endif
            
            internal_set_error(connection, INTERNAL_ERRCODE_NETWORK, "Unexpected EOF found while reading data: %s (%s).", strerror(errno), msg);
            goto abort_read;
        }
        
        // forward read to callback
        bool result = forward_cb(buffer, nread, xdata, xdata2);
        if (!result) goto abort_read;
        
        // update internal counter
        tread += (uint32_t)nread;
        
        // determine command length
        if (clen == 0) {
            clen = internal_parse_number (&original[1], tread-1, &cstart, NULL);
            
            // handle special cases
            if ((original[0] == CMD_INT) || (original[0] == CMD_FLOAT) || (original[0] == CMD_NULL)) clen = 0;
            else if (clen == 0) continue;
        }
        
        // check if read is complete
        if (clen + cstart + 1 == tread) break;
    }
    
    return true;
    
abort_read:
    return false;
}

static SQCloudResult *internal_socket_read (SQCloudConnection *connection, bool mainfd) {
    // most of the time one read will be sufficient
    char header[4096];
    char *buffer = (char *)&header;
    uint32_t blen = sizeof(header);
    uint32_t tread = 0;
    
    uint32_t cstart = 0;
    uint32_t clen = 0;

    int fd = (mainfd) ? connection->fd : connection->pubsubfd;
    #ifndef SQLITECLOUD_DISABLE_TSL
    struct tls *tls = (mainfd) ? connection->tls_context : connection->tls_pubsub_context;
    #endif
    
    char *original = buffer;
    while (1) {
        #ifndef SQLITECLOUD_DISABLE_TSL
        ssize_t nread = (tls) ? tls_read(tls, buffer, blen) : readsocket(fd, buffer, blen);
        if ((tls) && (nread == TLS_WANT_POLLIN || nread == TLS_WANT_POLLOUT)) continue;
        #else
        ssize_t nread = readsocket(fd, buffer, blen);
        #endif
        
        if (nread < 0) {
            const char *msg = "";
            #ifndef SQLITECLOUD_DISABLE_TSL
            if (tls) msg = tls_error(tls);
            #endif
            
            internal_set_error(connection, INTERNAL_ERRCODE_NETWORK, "An error occurred while reading data: %s (%s).", strerror(errno), msg);
            goto abort_read;
        }
        
        if (nread == 0) {
            const char *msg = "";
            #ifndef SQLITECLOUD_DISABLE_TSL
            if (tls) msg = tls_error(tls);
            #endif
            
            internal_set_error(connection, INTERNAL_ERRCODE_NETWORK, "Unexpected EOF found while reading data: %s (%s).", strerror(errno), msg);
            goto abort_read;
        }
        
        tread += (uint32_t)nread;
        blen -= (uint32_t)nread;
        buffer += nread;
        
        if (internal_has_commandlen(original[0])) {
            // parse buffer looking for command length
            if (clen == 0) clen = internal_parse_number (&original[1], tread-1, &cstart, NULL);
            if (clen == 0) continue;
            
            // check if read is complete
            // clen is the lenght parsed in the buffer
            // cstart is the index of the first space
            // +1 because we skipped the first character in the internal_parse_number function
            if (clen + cstart + 1 != tread) {
                // check buffer allocation and continue reading
                if (clen + cstart - tread > blen) {
                    char *clone = mem_alloc(clen + cstart + 1);
                    if (!clone) {
                        internal_set_error(connection, INTERNAL_ERRCODE_MEMORY, "Unable to allocate memory: %d.", clen + cstart + 1);
                        goto abort_read;
                    }
                    memcpy(clone, original, tread);
                    buffer = original = clone;
                    blen = (clen + cstart + 1) - tread;
                    buffer += tread;
                }
                continue;
            }
        } else {
            // it is a command with no explicit len
            // so make sure that the final character is a space
            if (original[tread-1] != ' ') continue;
        }
        
        // command is complete so parse it
        return internal_parse_buffer(connection, original, tread, (clen) ? cstart : 0, (original == header), false);
    }
    
abort_read:
    if (original != (char *)&header) mem_free(original);
    return NULL;
}

static bool internal_socket_raw_write (SQCloudConnection *connection, const char *buffer) {
    // this function is used only to debug possible security issues
    int fd = connection->fd;
    #ifndef SQLITECLOUD_DISABLE_TSL
    struct tls *tls = connection->tls_context;
    #endif
    
    size_t len = strlen(buffer);
    size_t written = 0;
    while (len > 0) {
        #ifndef SQLITECLOUD_DISABLE_TSL
        ssize_t nwrote = (tls) ? tls_write(tls, buffer, len) : writesocket(fd, buffer, len);
        if ((tls) && (nwrote == TLS_WANT_POLLIN || nwrote == TLS_WANT_POLLOUT)) continue;
        #else
        ssize_t nwrote = writesocket(fd, buffer, len);
        #endif
        
        if (nwrote < 0) {
            const char *msg = "";
            #ifndef SQLITECLOUD_DISABLE_TSL
            if (tls) msg = tls_error(tls);
            #endif
            return internal_set_error(connection, INTERNAL_ERRCODE_NETWORK, "An error occurred while writing data: %s (%s).", strerror(errno), msg);
        } else if (nwrote == 0) {
            return true;
        } else {
            written += nwrote;
            buffer += nwrote;
            len -= nwrote;
        }
    }
    
    return true;
}

static bool internal_socket_write (SQCloudConnection *connection, const char *buffer, size_t len, bool mainfd) {
    int fd = (mainfd) ? connection->fd : connection->pubsubfd;
    #ifndef SQLITECLOUD_DISABLE_TSL
    struct tls *tls = (mainfd) ? connection->tls_context : connection->tls_pubsub_context;
    #endif
    
    size_t written = 0;
    
    // write header
    char header[32];
    char *p = header;
    int hlen = snprintf(header, sizeof(header), "+%zu ", len);
    int len1 = hlen;
    while (len1) {
        #ifndef SQLITECLOUD_DISABLE_TSL
        ssize_t nwrote = (tls) ? tls_write(tls, p, len1) : writesocket(fd, p, len1);
        if ((tls) && (nwrote == TLS_WANT_POLLIN || nwrote == TLS_WANT_POLLOUT)) continue;
        #else
        ssize_t nwrote = writesocket(fd, p, len1);
        #endif
        
        if ((nwrote < 0) || (nwrote == 0 && written != hlen)) {
            const char *msg = "";
            #ifndef SQLITECLOUD_DISABLE_TSL
            if (tls) msg = tls_error(tls);
            #endif
            return internal_set_error(connection, INTERNAL_ERRCODE_NETWORK, "An error occurred while writing header data: %s (%s).", strerror(errno), msg);
        } else {
            written += nwrote;
            p += nwrote;
            len1 -= nwrote;
        }
    }
    
    // write buffer
    written = 0;
    while (len > 0) {
        #ifndef SQLITECLOUD_DISABLE_TSL
        ssize_t nwrote = (tls) ? tls_write(tls, buffer, len) : writesocket(fd, buffer, len);
        if ((tls) && (nwrote == TLS_WANT_POLLIN || nwrote == TLS_WANT_POLLOUT)) continue;
        #else
        ssize_t nwrote = writesocket(fd, buffer, len);
        #endif
        
        if (nwrote < 0) {
            const char *msg = "";
            #ifndef SQLITECLOUD_DISABLE_TSL
            if (tls) msg = tls_error(tls);
            #endif
            return internal_set_error(connection, INTERNAL_ERRCODE_NETWORK, "An error occurred while writing data: %s (%s).", strerror(errno), msg);
        } else if (nwrote == 0) {
            return true;
        } else {
            written += nwrote;
            buffer += nwrote;
            len -= nwrote;
        }
    }
    
    return true;
}

static void internal_socket_set_timeout (int sockfd, int timeout_secs) {
    #ifdef _WIN32
    DWORD timeout = timeout_secs * 1000;
    setsockopt(sockfd, SOL_SOCKET, SO_RCVTIMEO, (const char*)&timeout, sizeof timeout);
    setsockopt(sockfd, SOL_SOCKET, SO_SNDTIMEO, (const char*)&timeout, sizeof timeout);
    #else
    struct timeval tv;
    tv.tv_sec = timeout_secs;
    tv.tv_usec = 0;
    setsockopt(sockfd, SOL_SOCKET, SO_RCVTIMEO, (const char*)&tv, sizeof tv);
    setsockopt(sockfd, SOL_SOCKET, SO_SNDTIMEO, (const char*)&tv, sizeof tv);
    #endif
}

static bool internal_connect_apply_config (SQCloudConnection *connection, SQCloudConfig *config) {
    if (config->timeout) {
        internal_socket_set_timeout(connection->fd, config->timeout);
    }

    char buffer[2048];
    int len = 0;
    
    if (config->username && config->password && strlen(config->username) && strlen(config->password)) {
        len += snprintf(&buffer[len], sizeof(buffer) - len, "AUTH USER %s PASSWORD %s;", config->username, config->password);
    }
    
    if (config->database && strlen(config->database)) {
        len += snprintf(&buffer[len], sizeof(buffer) - len, "USE DATABASE %s;", config->database);
    }
    
    if (config->sqlite_mode) {
        len += snprintf(&buffer[len], sizeof(buffer) - len, "SET KEY CLIENT_SQLITE TO 1;");
    }
    
    if (config->compression) {
        len += snprintf(&buffer[len], sizeof(buffer) - len, "SET KEY CLIENT_COMPRESSION TO 1;");
    }
    
    if (config->zero_text) {
        len += snprintf(&buffer[len], sizeof(buffer) - len, "SET KEY CLIENT_ZEROTEXT TO 1;");
    }
    
    if (len > 0) {
        SQCloudResult *res = internal_run_command(connection, buffer, strlen(buffer), true);
        if (res != &SQCloudResultOK) return false;
    }
    
    return true;
}

static bool internal_connect (SQCloudConnection *connection, const char *hostname, int port, SQCloudConfig *config, bool mainfd) {
    // ipv4/ipv6 specific variables
    struct addrinfo hints, *addr_list = NULL, *addr;
    
    // ipv6 code from https://www.ibm.com/support/knowledgecenter/ssw_ibm_i_72/rzab6/xip6client.htm
    memset(&hints, 0, sizeof(hints));
    hints.ai_family = (config) ? config->family : AF_INET;
    hints.ai_socktype = SOCK_STREAM;
    
    // get the address information for the server using getaddrinfo()
    char port_string[256];
    snprintf(port_string, sizeof(port_string), "%d", port);
    int rc = getaddrinfo(hostname, port_string, &hints, &addr_list);
    if (rc != 0 || addr_list == NULL) {
        return internal_set_error(connection, INTERNAL_ERRCODE_NETWORK, "Error while resolving getaddrinfo (host %s not found).", hostname);
    }
    
    // begin non-blocking connection loop
    int sock_index = 0;
    int sock_current = 0;
    int sock_list[MAX_SOCK_LIST] = {0};
    for (addr = addr_list; addr != NULL; addr = addr->ai_next, ++sock_index) {
        if (sock_index >= MAX_SOCK_LIST) break;
        if ((addr->ai_family != AF_INET) && (addr->ai_family != AF_INET6)) continue;
        
        sock_current = socket(addr->ai_family, addr->ai_socktype, addr->ai_protocol);
        if (sock_current < 0) continue;
        
        // set socket options
        int len = 1;
        setsockopt(sock_current, SOL_SOCKET, SO_KEEPALIVE, (const char *) &len, sizeof(len));
        len = 1;
        setsockopt(sock_current, IPPROTO_TCP, TCP_NODELAY, (const char *) &len, sizeof(len));
        #ifdef SO_NOSIGPIPE
        len = 1;
        setsockopt(sock_current, SOL_SOCKET, SO_NOSIGPIPE, (const char *) &len, sizeof(len));
        #endif
        
        // by default, an IPv6 socket created on Windows Vista and later only operates over the IPv6 protocol
        // in order to make an IPv6 socket into a dual-stack socket, the setsockopt function must be called
        if (addr->ai_family == AF_INET6) {
            #ifdef _WIN32
            DWORD ipv6only = 0;
            #else
            int   ipv6only = 0;
            #endif
            setsockopt(sock_current, IPPROTO_IPV6, IPV6_V6ONLY, (void *)&ipv6only, sizeof(ipv6only));
        }
        
        // turn on non-blocking
        unsigned long ioctl_blocking = 1;    /* ~0; //TRUE; */
        ioctl(sock_current, FIONBIO, &ioctl_blocking);
        
        // initiate non-blocking connect ignoring return code
        connect(sock_current, addr->ai_addr, addr->ai_addrlen);
        
        // add sock_current to internal list of trying to connect sockets
        sock_list[sock_index] = sock_current;
    }
    
    // free not more needed memory
    freeaddrinfo(addr_list);
    
    // calculate the connection timeout and reset timers
    // if timeout is <= 0 then it is set to SQCLOUD_DEFAULT_TIMEOUT for the connect phase
    int connect_timeout = (config && config->timeout > 0) ? config->timeout : SQCLOUD_DEFAULT_TIMEOUT;
    time_t start = time(NULL);
    time_t now = start;
    rc = 0;
    
    int sockfd = 0;
    fd_set write_fds;
    fd_set except_fds;
    struct timeval tv;
    
    while (rc == 0 && ((now - start) < connect_timeout)) {
        FD_ZERO(&write_fds);
        FD_ZERO(&except_fds);
        
        int nfds = 0;
        for (int i=0; i<MAX_SOCK_LIST; ++i) {
            if (sock_list[i]) {
                FD_SET(sock_list[i], &write_fds);
                FD_SET(sock_list[i], &except_fds);
                if (nfds < sock_list[i]) nfds = sock_list[i];
            }
        }
        
        tv.tv_sec = connect_timeout;
        tv.tv_usec = 0;
        rc = select(nfds + 1, NULL, &write_fds, &except_fds, &tv);
        
        if (rc == 0) break; // timeout
        else if (rc == -1) {
            if (errno == EINTR || errno == EAGAIN || errno == EINPROGRESS) continue;
            break; // handle error
        }
        
        // check for error first
        for (int i=0; i<MAX_SOCK_LIST; ++i) {
            if (sock_list[i] > 0) {
                if (FD_ISSET(sock_list[i], &except_fds)) {
                    closesocket(sock_list[i]);
                    sock_list[i] = 0;
                }
            }
        }
        
        // check which file descriptor is ready (need to check for socket error also)
        for (int i=0; i<MAX_SOCK_LIST; ++i) {
            if (sock_list[i] > 0) {
                if (FD_ISSET(sock_list[i], &write_fds)) {
                    int err = socket_geterror(sock_list[i]);
                    if (err > 0) {
                        closesocket(sock_list[i]);
                        sock_list[i] = 0;
                    } else {
                        sockfd = sock_list[i];
                        break;
                    }
                }
            }
        }
        // check if a valid descriptor has been found
        if (sockfd != 0) break;
        
        // no socket ready yet
        now = time(NULL);
        rc = 0;
    }
    
    // close still opened sockets
    for (int i=0; i<MAX_SOCK_LIST; ++i) {
        if ((sock_list[i] > 0) && (sock_list[i] != sockfd)) closesocket(sock_list[i]);
    }
    
    // bail if there was an error
    if (rc < 0) {
        return internal_set_error(connection, INTERNAL_ERRCODE_NETWORK, "An error occurred while trying to connect: %s.", strerror(errno));
    }
    
    // bail if there was a timeout
    if ((time(NULL) - start) >= connect_timeout) {
        return internal_set_error(connection, INTERNAL_ERRCODE_NETWORK, "Connection timeout while trying to connect (%d).", connect_timeout);
    }
    
    // turn off non-blocking
    int ioctl_blocking = 0;    /* ~0; //TRUE; */
    ioctl(sockfd, FIONBIO, &ioctl_blocking);
    
    // finalize connection
    if (mainfd) {
        connection->fd = sockfd;
        connection->port = port;
        connection->hostname = strdup(hostname);
        #ifndef SQLITECLOUD_DISABLE_TSL
        if (config && !config->insecure) {
            int rc = tls_connect_socket(connection->tls_context, sockfd, hostname);
            if (rc < 0) printf("Error on tls_connect_socket: %s\n", tls_error(connection->tls_context));
        }
        #endif
    } else {
        connection->pubsubfd = sockfd;
        #ifndef SQLITECLOUD_DISABLE_TSL
        if (config && !config->insecure) {
            int rc = tls_connect_socket(connection->tls_pubsub_context, sockfd, hostname);
            if (rc < 0) printf("Error on tls_connect_socket\n");
        }
        #endif
    }
    return true;
}

void internal_rowset_dump (SQCloudResult *result, uint32_t maxline, bool quiet) {
    uint32_t nrows = result->nrows;
    uint32_t ncols = result->ncols;
    uint32_t blen = result->blen;
    
    // if user specify a maxline then do not print more than maxline characters for every column
    if (maxline > 0) {
        for (uint32_t i=0; i<ncols; ++i) {
            if (result->clen[i] > maxline) result->clen[i] = maxline;
        }
    }
    
    // print separator header
    for (uint32_t i=0; i<ncols; ++i) {
        for (uint32_t j=0; j<result->clen[i]+2; ++j) putchar('-');
        putchar('|');
    }
    printf("\n");
    
    // print column names
    for (uint32_t i=0; i<ncols; ++i) {
        uint32_t len = blen;
        uint32_t delta = 0;
        char *value = internal_parse_value(result->name[i], &len, NULL);
        
        // UTF-8 strings need special adjustments
        uint32_t utf8len = utf8_len(value, len);
        if (utf8len != len) delta = len - utf8len;
        printf(" %-*.*s |", result->clen[i] + delta, (maxline && len > maxline) ? maxline : len, value);
        blen -= len;
    }
    printf("\n");
    
    // print separator header
    for (uint32_t i=0; i<ncols; ++i) {
        for (uint32_t j=0; j<result->clen[i]+2; ++j) putchar('-');
        putchar('|');
    }
    printf("\n");
    
    #if 0
    // print types (just for debugging)
    printf("\n");
    for (uint32_t i=0; i<nrows; ++i) {
        for (uint32_t j=0; j<ncols; ++j) {
            SQCloudValueType type = SQCloudRowsetValueType(result, i, j);
            printf("%d ", type);
        }
        printf("\n");
    }
    printf("\n");
    #endif
    
    // print result
    for (uint32_t i=0; i<nrows * ncols; ++i) {
        uint32_t len = blen;
        uint32_t delta = 0;
        char *value = internal_parse_value(result->data[i], &len, NULL);
        blen -= len;

        // UTF-8 strings need special adjustments
        if (!value) {value = "NULL"; len = 4;}
        uint32_t utf8len = utf8_len(value, len);
        if (utf8len != len) delta = len - utf8len;
        printf(" %-*.*s |", (result->clen[i % ncols]) + delta, (maxline && len > maxline) ? maxline : len, value);
        
        bool newline = (((i+1) % ncols == 0) || (ncols == 1));
        if (newline) printf("\n");
    }
    
    // print footer
    for (uint32_t i=0; i<ncols; ++i) {
        for (uint32_t j=0; j<result->clen[i]+2; ++j) putchar('-');
        putchar('|');
    }
    printf("\n");
    
    printf("Rows: %d - Cols: %d - Bytes: %d", result->nrows, result->ncols, result->blen);
    if (!quiet) printf(" Time: %f secs", result->time);
    fflush( stdout );
}

// MARK: - URL -

static int char2hex (int c) {
    if (isdigit(c)) return (c - '0');
    c = toupper(c);
    if (c >='A' && c <='F') return (c - 'A' + 0x0A);
    return -1;
}

static int url_decode (char s[512]) {
    int i = 0;
    int j = 0;
    int len = (int)strlen(s);
    
    while (i < len) {
        int c = s[i];
        if (c == '%') {
            if (i + 2 >= len) return 0;
            c = (char2hex(s[i+1]) * 0x10) + char2hex(s[i+2]);
            if (c < 0) return 0;
            s[j] = c;
            i += 2;
        } else {
            s[j] = c;
        }
        ++i;
        ++j;
    }
    s[j] = 0;
    return j;
}

static int url_extract_username_password (const char *s, char b1[512], char b2[512]) {
    // user:pass@host.com:port/dbname?timeout=10&key2=value2&key3=value3
    
    // lookup username (if any)
    char *username = strchr(s, ':');
    if (!username) return 0;
    size_t len = username - s;
    if (len > 511) return -1;
    memcpy(b1, s, len);
    b1[len] = 0;
    if (url_decode(b1) <= 0) return 0;
    
    // lookup username (if any)
    char *password = strchr(s, '@');
    if (!password) return 0;
    len = password - username - 1;
    if (len > 511) return -1;
    memcpy(b2, username+1, len);
    b2[len] = 0;
    if (url_decode(b2) <= 0) return 0;
    
    return (int)(password - s) + 1;
}

static int url_extract_hostname_port (const char *s, char b1[512], char b2[512]) {
    // host.com:port/dbname?timeout=10&key2=value2&key3=value3
    
    // lookup hostname (if any)
    char *hostname = strchr(s, ':');
    if (!hostname) hostname = strchr(s, '/');
    if (!hostname) hostname = strchr(s, '?');
    if (!hostname) hostname = strchr(s, 0);
    if (!hostname) return -1;
    size_t len = hostname - s;
    if (len > 511) return -1;
    memcpy(b1, s, len);
    b1[len] = 0;
    if (url_decode(b1) <= 0) return 0;
    
    // lookup port (if any)
    char *port = strchr(s, ':');
    if (port) {
        char *p = port + 1;
        ++len;
        
        int i = 0;
        while (p[0]) {
            if ((p[0] == '/') || (p[0] == '?') || (p[0] == 0)) break;
            if (i+1 > 511) return -1;
            b2[i++] = p[0];
            ++len;
            ++p;
        }
        b2[len] = 0;
        if (url_decode(b2) <= 0) return 0;
    }
    
    // adjust returned len
    if (s[len] != 0) ++len;
    
    return (int)len;
}

static int url_extract_database (const char *s, char b1[512]) {
    // dbname?timeout=10&key2=value2&key3=value3
    
    // lookup database (if any)
    char *database = strchr(s, '?');
    if (database) {
        size_t len = database - s;
        if (len > 511) return -1;
        memcpy(b1, s, len);
        b1[len] = 0;
        if (url_decode(b1) <= 0) return 0;
        
        return (int)(len + 1);
    }
    
    // there is no ? separator character
    // that means that there should be
    // no key/value
    char *guard = strchr(s, '=');
    if (guard) return 0;
    
    // database name is the s string
    size_t len = strlen(s);
    if (len > 511) return -1;
    memcpy(b1, s, len);
    b1[len] = 0;
    if (url_decode(b1) <= 0) return 0;
    
    return (int)len;
}

static int url_extract_keyvalue (const char *s, char b1[512], char b2[512]) {
    // timeout=10&key2=value2&key3=value3
    
    // lookup key (if any)
    char *key = strchr(s, '=');
    if (!key) return 0;
    size_t len = key - s;
    if (len > 511) return -1;
    memcpy(b1, s, len);
    b1[len] = 0;
    if (url_decode(b1) <= 0) return 0;
    
    // lookup value (if any)
    char *value = strchr(s, '&');
    if (!value) value = strchr(s, 0);
    if (!value) return 0;
    len = value - key - 1;
    if (len > 511) return -1;
    memcpy(b2, key+1, len);
    b2[len] = 0;
    if (url_decode(b2) <= 0) return 0;
    
    return (int)(value - s) + 1;
}

// MARK: - RESERVED -

bool _reserved1 (SQCloudConnection *connection, const char *command, bool (*forward_cb) (char *buffer, size_t blen, void *xdata, void *xdata2), void *xdata, void *xdata2) {
    if (!forward_cb) return false;
    if (!internal_socket_write(connection, command, strlen(command), true)) return false;
    if (!internal_socket_forward_read(connection, forward_cb, xdata, xdata2)) return false;
    return true;
}

SQCloudResult *_reserved2 (SQCloudConnection *connection, const char *username, const char *passwordhash, const char *UUID) {
    char buffer[1024];
    int len = 0;
    if (username) {
        len += snprintf(&buffer[len], sizeof(buffer) - len, "AUTH USER %s HASH %s;", username, passwordhash);
    }
    len += snprintf(&buffer[len], sizeof(buffer) - len, "SET CLIENT UUID TO %s;", UUID);
    return internal_run_command(connection, buffer, strlen(buffer), true);
}

SQCloudResult *_reserved3 (char *buffer, uint32_t blen, uint32_t cstart, SQCloudResult *chunk) {
    SQCloudConnection connection = {0};
    connection._chunk = chunk;
    SQCloudResult *res = internal_parse_buffer(&connection, buffer, blen, cstart, false, true);
    return res;
}

uint32_t _reserved4 (char *buffer, uint32_t blen, uint32_t *cstart) {
    return internal_parse_number(buffer, blen, cstart, NULL);
}

bool _reserved5 (SQCloudResult *res) {
    return res->ischunk;
}

bool _reserved6 (SQCloudConnection *connection, const char *buffer) {
    internal_clear_error(connection);
    return internal_socket_raw_write(connection, buffer);
}

SQCloudResult *_reserved7 (SQCloudConnection *connection) {
    return internal_socket_read(connection, true);
}

// MARK: - PUBLIC -

SQCloudConnection *SQCloudConnect (const char *hostname, int port, SQCloudConfig *config) {
    internal_init();
    
    SQCloudConnection *connection = mem_zeroalloc(sizeof(SQCloudConnection));
    if (!connection) return NULL;
    
    #ifndef SQLITECLOUD_DISABLE_TSL
    if (!internal_setup_tls(connection, config, true)) return connection;
    #endif
    
    if (internal_connect(connection, hostname, port, config, true)) {
        if (config) internal_connect_apply_config(connection, config);
        connection->_config = config;
    }
    
    return connection;
}

SQCloudConnection *SQCloudConnectWithString (const char *s) {
    // URL STRING FORMAT
    // sqlitecloud://user:pass@host.com:port/dbname?timeout=10&key2=value2&key3=value3
    
    // sanity check
    const char domain[] = "sqlitecloud://";
    int n = sizeof(domain) - 1;
    if (strncmp(s, domain, n) != 0) return NULL;
    
    // config struct
    SQCloudConfig *config = (SQCloudConfig *)mem_zeroalloc(sizeof(SQCloudConfig));
    if (!config) return NULL;
    
    // default IPv4
    config->family = SQCLOUD_IPv4;
    
    // lookup for optional username/password
    char username[512];
    char password[512];
    int rc = url_extract_username_password(&s[n], username, password);
    if (rc == -1) {
        mem_free(config);
        return NULL;
    }
    if (rc) {
        config->username = string_dup(username);
        config->password = string_dup(password);
    }
    
    // lookup for mandatory hostname
    n += rc;
    char hostname[512];
    char port_s[512];
    rc = url_extract_hostname_port(&s[n], hostname, port_s);
    if (rc <= 0) {
        mem_free(config);
        return NULL;
    }
    int port = (int)strtol(port_s, NULL, 0);
    if (port <= 0) port = SQCLOUD_DEFAULT_PORT;
    
    // lookup for optional database
    n += rc;
    char database[512];
    rc = url_extract_database(&s[n], database);
    if (rc == -1) {
        mem_free(config);
        return NULL;
    }
    if (rc > 0) {
        config->database = string_dup(database);
    }
    
    // lookup for optional key(s)/value(s)
    n += rc;
    char key[512];
    char value[512];
    while ((rc = url_extract_keyvalue(&s[n], key, value)) > 0) {
        if (strcasecmp(key, "timeout") == 0) {
            int timeout = (int)strtol(value, NULL, 0);
            config->timeout = (timeout > 0) ? timeout : 0;
        }
        else if (strcasecmp(key, "compression") == 0) {
            int compression = (int)strtol(value, NULL, 0);
            config->compression = (compression > 0) ? true : false;
        }
        else if (strcasecmp(key, "sqlite") == 0) {
            int sqlite_mode = (int)strtol(value, NULL, 0);
            config->sqlite_mode = (sqlite_mode > 0) ? true : false;
        }
        else if (strcasecmp(key, "zerotext") == 0) {
            int zero_text = (int)strtol(value, NULL, 0);
            config->zero_text = (zero_text > 0) ? true : false;
        }
        else if (strcasecmp(key, "memory") == 0) {
            int in_memory = (int)strtol(value, NULL, 0);
            if (in_memory) config->database = string_dup(":memory:");
        }
        #ifndef SQLITECLOUD_DISABLE_TSL
        else if (strcasecmp(key, "insecure") == 0) {
            int insecure = (int)strtol(value, NULL, 0);
            config->insecure = (insecure > 0) ? true : false;
        }
        else if (strcasecmp(key, "root_certificate") == 0) {
            config->tls_root_certificate = strdup(value);
        }
        else if (strcasecmp(key, "client_certificate") == 0) {
            config->tls_certificate = strdup(value);
        }
        else if (strcasecmp(key, "client_certificate_key") == 0) {
            config->tls_certificate_key = strdup(value);
        }
        #endif
        n += rc;
    }
    
    return SQCloudConnect(hostname, port, config);
}

SQCloudResult *SQCloudExec (SQCloudConnection *connection, const char *command) {
    return internal_run_command(connection, command, strlen(command), true);
}

SQCloudResult *SQCloudRead (SQCloudConnection *connection) {
    return internal_socket_read(connection, true);
}

void SQCloudDisconnect (SQCloudConnection *connection) {
    if (!connection) return;
    
    // free TLS
    #ifndef SQLITECLOUD_DISABLE_TSL
    if (connection->tls_context) {
        tls_close(connection->tls_context);
        tls_free(connection->tls_context);
    }
    
    if (connection->tls_pubsub_context) {
        tls_close(connection->tls_pubsub_context);
        tls_free(connection->tls_pubsub_context);
    }
    #endif
    
    // try to gracefully close connections
    if (connection->fd) {
        closesocket(connection->fd);
    }
    
    if (connection->pubsubfd) {
        closesocket(connection->pubsubfd);
    }
    
    // free memory
    if (connection->hostname) {
        mem_free(connection->hostname);
    }
    
    if (connection->uuid) {
        mem_free(connection->uuid);
    }
    
    mem_free(connection);
}

void SQCloudSetPubSubCallback (SQCloudConnection *connection, SQCloudPubSubCB callback, void *data) {
    connection->callback = callback;
    connection->data = data;
}

SQCloudResult *SQCloudSetPubSubOnly (SQCloudConnection *connection) {
    if (!connection->callback) {
        internal_set_error(connection, INTERNAL_ERRCODE_PUBSUB, "A PubSub callback must be set before executing a PUBSUB ONLY command.");
        return NULL;
    }
    
    const char *command = "PUBSUB ONLY";
    return internal_run_command(connection, command, strlen(command), true);
}

char *SQCloudUUID (SQCloudConnection *connection) {
    return connection->uuid;
}

// MARK: -

bool SQCloudIsError (SQCloudConnection *connection) {
    return (!connection || connection->errcode);
}

bool SQCloudIsSQLiteError (SQCloudConnection *connection) {
    // https://www.sqlite.org/rescode.html
    return (connection && connection->errcode < 10000);
}

int SQCloudErrorCode (SQCloudConnection *connection) {
    return (connection) ? connection->errcode : INTERNAL_ERRCODE_GENERIC;
}

int SQCloudExtendedErrorCode (SQCloudConnection *connection) {
    return (connection) ? connection->xerrcode : 0;
}

const char *SQCloudErrorMsg (SQCloudConnection *connection) {
    return (connection) ? connection->errmsg : "Not enough memory to allocate a SQCloudConnection.";
}

void SQCloudErrorReset (SQCloudConnection *connection) {
    internal_clear_error(connection);
}

void SQCloudErrorSetCode (SQCloudConnection *connection, int errcode) {
    connection->errcode = errcode;
}

void SQCloudErrorSetMsg (SQCloudConnection *connection, const char *format, ...) {
    va_list arg;
    va_start (arg, format);
    vsnprintf(connection->errmsg, sizeof(connection->errmsg), format, arg);
    va_end (arg);
}

// MARK: -

SQCloudResType SQCloudResultType (SQCloudResult *result) {
    return (result) ? result->tag : RESULT_ERROR;
}

bool SQCloudResultIsOK (SQCloudResult *result) {
    return (result == &SQCloudResultOK);
}

bool SQCloudResultIsError (SQCloudResult *result) {
    return (!result);
}

uint32_t SQCloudResultLen (SQCloudResult *result) {
    return (result) ? result->blen : 0;
}

char *SQCloudResultBuffer (SQCloudResult *result) {
    return (result) ? result->buffer : NULL;
}

int32_t SQCloudResultInt32 (SQCloudResult *result) {
    if ((!result) || (result->tag != RESULT_INTEGER)) return 0;
    
    char *buffer = result->buffer;
    buffer[result->blen] = 0;
    return (int32_t)strtol(buffer, NULL, 0);
}

int64_t SQCloudResultInt64 (SQCloudResult *result) {
    if ((!result) || (result->tag != RESULT_INTEGER)) return 0;
    
    char *buffer = result->buffer;
    buffer[result->blen] = 0;
    return (int64_t)strtoll(buffer, NULL, 0);
}

double SQCloudResultDouble (SQCloudResult *result) {
    if ((!result) || (result->tag != RESULT_FLOAT)) return 0.0;
    
    char *buffer = result->buffer;
    buffer[result->blen] = 0;
    return (double)strtod(buffer, NULL);
}

void SQCloudResultFree (SQCloudResult *result) {
    if (!result || (result == &SQCloudResultOK) || (result == &SQCloudResultNULL)) return;
    
    if (!result->ischunk && !result->externalbuffer) {
        mem_free(result->rawbuffer);
    }
    
    if (result->tag == RESULT_ROWSET) {
        mem_free(result->name);
        mem_free(result->data);
        mem_free(result->clen);
        if (result->decltype) mem_free(result->decltype);
        if (result->dbname) mem_free(result->dbname);
        if (result->tblname) mem_free(result->tblname);
        if (result->origname) mem_free(result->origname);
        
        if (result->ischunk && !result->externalbuffer) {
            for (uint32_t i = 0; i<result->bcount; ++i) {
                if (result->buffers[i]) mem_free(result->buffers[i]);
            }
            mem_free(result->buffers);
        }
    }
    
    if (result->tag == RESULT_ARRAY) {
        mem_free(result->data);
    }
    
    mem_free(result);
}

// MARK: -

// https://database.guide/2-sample-databases-sqlite/
// https://embeddedartistry.com/blog/2017/07/05/printf-a-limited-number-of-characters-from-a-string/
// https://stackoverflow.com/questions/1809399/how-to-format-strings-using-printf-to-get-equal-length-in-the-output

// SET DATABASE mediastore.sqlite
// SELECT * FROM Artist LIMIT 10;

static bool SQCloudRowsetSanityCheck (SQCloudResult *result, uint32_t row, uint32_t col) {
    if (!result || result->tag != RESULT_ROWSET) return false;
    if ((row >= result->nrows) || (col >= result->ncols)) return false;
    return true;
}

SQCloudValueType SQCloudRowsetValueType (SQCloudResult *result, uint32_t row, uint32_t col) {
    if (!SQCloudRowsetSanityCheck(result, row, col)) return VALUE_NULL;
    return internal_type(result->data[row*result->ncols+col]);
}

uint32_t SQCloudRowsetRowsMaxColumnLength (SQCloudResult *result, uint32_t col) {
    return (result) ? result->clen[ col ] : 0;
}

char *SQCloudRowsetColumnName (SQCloudResult *result, uint32_t col, uint32_t *len) {
    if (!result || result->tag != RESULT_ROWSET) return NULL;
    if (col >= result->ncols) return NULL;
    *len = result->blen - (uint32_t)(result->name[col] - result->rawbuffer);
    return internal_parse_value(result->name[col], len, NULL);
}

uint32_t SQCloudRowsetRows (SQCloudResult *result) {
    if (!SQCloudRowsetSanityCheck(result, 0, 0)) return 0;
    return result->nrows;
}

uint32_t SQCloudRowsetCols (SQCloudResult *result) {
    if (!SQCloudRowsetSanityCheck(result, 0, 0)) return 0;
    return result->ncols;
}

uint32_t SQCloudRowsetMaxLen (SQCloudResult *result) {
    if (!SQCloudRowsetSanityCheck(result, 0, 0)) return 0;
    return result->maxlen;
}

char *SQCloudRowsetValue (SQCloudResult *result, uint32_t row, uint32_t col, uint32_t *len) {
    if (!SQCloudRowsetSanityCheck(result, row, col)) return NULL;
    
    // The *len var must contain the remaining length of the buffer pointed by
    // result->data[row*result->ncols+col]. The caller should not be aware of the
    // internal implementation of this buffer, so it must be set here.
    char *value = result->data[row*result->ncols+col];
    *len = (value) ? result->blen - (uint32_t)(value - result->rawbuffer) + result->nheader : 2;
    return internal_parse_value(value, len, NULL);
}

uint32_t SQCloudRowSetValueLen (SQCloudResult *result, uint32_t row, uint32_t col) {
    uint32_t len = 0;
    SQCloudRowsetValue(result, row, col, &len);
    return len;
}

int32_t SQCloudRowsetInt32Value (SQCloudResult *result, uint32_t row, uint32_t col) {
    if (!SQCloudRowsetSanityCheck(result, row, col)) return 0;
    uint32_t len = 0;
    char *value = internal_parse_value(result->data[row*result->ncols+col], &len, NULL);
    
    char buffer[256];
    snprintf(buffer, sizeof(buffer), "%.*s", len, value);
    return (int32_t)strtol(buffer, NULL, 0);
}

int64_t SQCloudRowsetInt64Value (SQCloudResult *result, uint32_t row, uint32_t col) {
    if (!SQCloudRowsetSanityCheck(result, row, col)) return 0;
    uint32_t len = 0;
    char *value = internal_parse_value(result->data[row*result->ncols+col], &len, NULL);
    
    char buffer[256];
    snprintf(buffer, sizeof(buffer), "%.*s", len, value);
    return (int64_t)strtoll(buffer, NULL, 0);
}

float SQCloudRowsetFloatValue (SQCloudResult *result, uint32_t row, uint32_t col) {
    if (!SQCloudRowsetSanityCheck(result, row, col)) return 0.0;
    uint32_t len = 0;
    char *value = internal_parse_value(result->data[row*result->ncols+col], &len, NULL);
    
    char buffer[256];
    snprintf(buffer, sizeof(buffer), "%.*s", len, value);
    return (float)strtof(buffer, NULL);
}

double SQCloudRowsetDoubleValue (SQCloudResult *result, uint32_t row, uint32_t col) {
    if (!SQCloudRowsetSanityCheck(result, row, col)) return 0.0;
    uint32_t len = 0;
    char *value = internal_parse_value(result->data[row*result->ncols+col], &len, NULL);
    
    char buffer[256];
    snprintf(buffer, sizeof(buffer), "%.*s", len, value);
    return (double)strtod(buffer, NULL);
}

void SQCloudRowsetDump (SQCloudResult *result, uint32_t maxline, bool quiet) {
    internal_rowset_dump(result, maxline, quiet);
}

// MARK: -

static bool SQCloudArraySanityCheck (SQCloudResult *result, uint32_t index) {
    if (!result || result->tag != RESULT_ARRAY) return false;
    if (index >= result->ndata) return false;
    return true;
}

SQCloudValueType SQCloudArrayValueType (SQCloudResult *result, uint32_t index) {
    if (!SQCloudArraySanityCheck(result, index)) return VALUE_NULL;
    return internal_type(result->data[index]);
}

uint32_t SQCloudArrayCount (SQCloudResult *result) {
    if (result->tag != RESULT_ARRAY) return 0;
    return result->ndata;
}

char *SQCloudArrayValue (SQCloudResult *result, uint32_t index, uint32_t *len) {
    if (!SQCloudArraySanityCheck(result, index)) return NULL;
    
    // The *len var must contain the remaining length of the buffer pointed by
    // result->data[index]. The caller should not be aware of the
    // internal implementation of this buffer, so it must be set here.
    char *value = result->data[index];
    *len = (value) ? result->blen - (uint32_t)(value - result->rawbuffer) + result->nheader : 2;
    return internal_parse_value(value, len, NULL);
}

int32_t SQCloudArrayInt32Value (SQCloudResult *result, uint32_t index) {
    if (!SQCloudArraySanityCheck(result, index)) return 0;
    uint32_t len = 0;
    char *value = internal_parse_value(result->data[index], &len, NULL);
    
    char buffer[256];
    snprintf(buffer, sizeof(buffer), "%.*s", len, value);
    return (int32_t)strtol(buffer, NULL, 0);
}

int64_t SQCloudArrayInt64Value (SQCloudResult *result, uint32_t index) {
    if (!SQCloudArraySanityCheck(result, index)) return 0;
    uint32_t len = 0;
    char *value = internal_parse_value(result->data[index], &len, NULL);
    
    char buffer[256];
    snprintf(buffer, sizeof(buffer), "%.*s", len, value);
    return (int64_t)strtoll(buffer, NULL, 0);
}

float SQCloudArrayFloatValue (SQCloudResult *result, uint32_t index) {
    if (!SQCloudArraySanityCheck(result, index)) return 0.0;
    uint32_t len = 0;
    char *value = internal_parse_value(result->data[index], &len, NULL);
    
    char buffer[256];
    snprintf(buffer, sizeof(buffer), "%.*s", len, value);
    return (float)strtof(buffer, NULL);
}

double SQCloudArrayDoubleValue (SQCloudResult *result, uint32_t index) {
    if (!SQCloudArraySanityCheck(result, index)) return 0.0;
    uint32_t len = 0;
    char *value = internal_parse_value(result->data[index], &len, NULL);
    
    char buffer[256];
    snprintf(buffer, sizeof(buffer), "%.*s", len, value);
    return (double)strtod(buffer, NULL);
}

void SQCloudArrayDump (SQCloudResult *result) {
    if (result->tag != RESULT_ARRAY) return;
    
    for (uint32_t i=0; i<result->ndata; ++i) {
        uint32_t len;
        char *value = SQCloudArrayValue(result, i, &len);
        if (!value) {value = "NULL"; len = 4;}
        printf("[%d] %.*s\n", i, len, value);
    }
}

// MARK: -

char *SQCloudBinaryToB64 (char *dest, void const *src, size_t *size) {
    char *buffer = bintob64(dest, src, *size);
    *size = (buffer) ? (buffer - dest) : 0;
    return buffer;
}

void *SQCloudB64ToBinary (void *dest, char const *src, size_t *size) {
    void *buffer = b64tobin(dest, src);
    *size = (buffer) ? (buffer - dest) : 0;
    return buffer;
}

size_t SQCloudComputeB64Size (size_t binarySize) {
    return COMPUTE_BASE64_SIZE(binarySize);
}

