//
//  sqcloud.c
//
//  Created by Marco Bambini on 08/02/21.
//

#include "lz4.h"
#include "sqcloud.h"
#include <ctype.h>
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

#define mem_realloc                         realloc
#define mem_zeroalloc(_s)                   calloc(1,_s)
#define mem_alloc(_s)                       malloc(_s)
#define mem_free(_s)                        free(_s)
#define string_dup(_s)                      strdup(_s)
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

#define CMD_MINLEN                          2

#define CONNSTRING_KEYVALUE_SEPARATOR       '='
#define CONNSTRING_TOKEN_SEPARATOR          ';'

#define DEFAULT_CHUCK_NBUFFERS              20
#define DEFAULT_CHUNK_MINROWS               2000

// MARK: - PROTOTYPES -

static SQCloudResult *internal_socket_read (SQCloudConnection *connection, bool mainfd);
static bool internal_socket_write (SQCloudConnection *connection, const char *buffer, size_t len, bool mainfd);
static uint32_t internal_parse_number (char *buffer, uint32_t blen, uint32_t *cstart);
static SQCloudResult *internal_parse_buffer (SQCloudConnection *connection, char *buffer, uint32_t blen, uint32_t cstart, bool isstatic, bool externalbuffer);
static bool internal_connect (SQCloudConnection *connection, const char *hostname, int port, SQCloudConfig *config, bool mainfd);

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
    
    // used in TYPE_ROWSET only
    uint32_t        nrows;                  // number of rows
    uint32_t        ncols;                  // number of columns
    uint32_t        ndata;                  // number of items stores in data
    char            **data;                 // data contained in the rowset
    char            **name;                 // column names
    uint32_t        *clen;                  // max len for each column (used to display result)
    uint32_t        maxlen;                 // max len for each row/column
} _SQCloudResult;

struct SQCloudConnection {
    int             fd;
    char            errmsg[1024];
    int             errcode;
    SQCloudResult   *_chunk;
    
    // pub/sub
    char            *uuid;
    int             pubsubfd;
    SQCloudPubSubCB callback;
    void            *data;
    char            *hostname;
    int             port;
    pthread_t       tid;
    
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
        ssize_t nread = readsocket(fd, buffer, blen);
        
        if (nread < 0) {
            printf("Handle error here %s.", strerror(errno));
            break;
            // internal_set_error(connection, 1, "An error occurred while reading data: %s.", strerror(errno));
            // goto abort_read;
        }
        
        if (nread == 0) {
            printf("Handle error here %s.", strerror(errno));
            break;
            // internal_set_error(connection, 1, "Unexpected EOF found while reading data: %s.", strerror(errno));
            // goto abort_read;
        }
        
        tread += (uint32_t)nread;
        blen -= (uint32_t)nread;
        buffer += nread;
        
        uint32_t cstart = 0;
        uint32_t clen = internal_parse_number (&original[1], tread-1, &cstart);
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
                    printf("Handle memory error here %s.", strerror(errno));
                    break;
                    // internal_set_error(connection, 1, "Unable to allocate memory: %d.", clen + cstart + 1);
                    // goto abort_read;
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
    connection->errmsg[0] = 0;
}

static bool internal_setup_ssl (SQCloudConnection *connection, SQCloudConfig *config) {
    return true;
}

static SQCloudValueType internal_type (char *buffer) {
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

static uint32_t internal_parse_number (char *buffer, uint32_t blen, uint32_t *cstart) {
    uint32_t value = 0;
    
    for (uint32_t i=0; i<blen; ++i) {
        if (buffer[i] == ' ') {
            *cstart = i+1;
            return value;
        }
        value = (value * 10) + (buffer[i] - '0');
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
    blen = internal_parse_number(&buffer[1], blen, &cstart);
    
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
    
    if (internal_connect(connection, connection->hostname, connection->port, NULL, false)) {
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

static SQCloudResult *internal_rowset_type (SQCloudConnection *connection, char *buffer, uint32_t blen, uint32_t bstart, SQCloudResType type) {
    SQCloudResult *rowset = (SQCloudResult *)mem_zeroalloc(sizeof(SQCloudResult));
    if (!rowset) {
        internal_set_error(connection, 1, "Unable to allocate memory for SQCloudResult: %d.", sizeof(SQCloudResult));
        return NULL;
    }
    
    rowset->tag = type;
    rowset->buffer = &buffer[bstart];
    rowset->rawbuffer = buffer;
    rowset->blen = blen;
    rowset->balloc = blen;
    
    return rowset;
}

static SQCloudResult *internal_parse_rowset (SQCloudConnection *connection, char *buffer, uint32_t blen, uint32_t bstart, uint32_t nrows, uint32_t ncols) {    
    SQCloudResult *rowset = (SQCloudResult *)mem_zeroalloc(sizeof(SQCloudResult));
    if (!rowset) {
        internal_set_error(connection, 1, "Unable to allocate memory for SQCloudResult: %d.", sizeof(SQCloudResult));
        return NULL;
    }
    
    rowset->tag = RESULT_ROWSET;
    rowset->buffer = buffer;
    rowset->rawbuffer = buffer;
    rowset->blen = blen;
    rowset->balloc = blen;
    
    rowset->nrows = nrows;
    rowset->ncols = ncols;
    rowset->data = (char **) mem_alloc(nrows * ncols * sizeof(char *));
    rowset->name = (char **) mem_alloc(ncols * sizeof(char *));
    rowset->clen = (uint32_t *) mem_zeroalloc(ncols * sizeof(uint32_t));
    if (!rowset->data || !rowset->name || !rowset->clen) goto abort_rowset;
    
    buffer += bstart;
    
    // the first column contains names
    for (uint32_t i=0; i<ncols; ++i) {
        uint32_t cstart = 0;
        uint32_t len = internal_parse_number(&buffer[1], blen, &cstart);
        rowset->name[i] = buffer;
        buffer += cstart + len + 1;
        blen -= cstart + len + 1;
        if (rowset->clen[i] < len) rowset->clen[i] = len;
        if (rowset->maxlen < len) rowset->maxlen = len;
    }
    
    // parse values
    for (uint32_t i=0; i<nrows * ncols; ++i) {
        uint32_t len = blen, cellsize;
        char *value = internal_parse_value(buffer, &len, &cellsize);
        rowset->data[i] = (value) ? buffer : NULL;
        buffer += cellsize;
        blen -= cellsize;
        ++rowset->ndata;
        if (rowset->clen[i % ncols] < len) rowset->clen[i % ncols] = len;
        if (rowset->maxlen < len) rowset->maxlen = len;
    }
    
    return rowset;
    
abort_rowset:
    if (rowset->data) mem_free(rowset->data);
    if (rowset->name) mem_free(rowset->name);
    if (rowset->clen) mem_free(rowset->clen);
    if (rowset) mem_free(rowset);
    
    internal_set_error(connection, 1, "Unable to allocate internal memory for SQCloudResult.");
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
            internal_set_error(connection, 1, "Unable to allocate memory for SQCloudResult: %d.", sizeof(SQCloudResult));
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
        
        rowset->brows = nrows + DEFAULT_CHUNK_MINROWS;
        rowset->nrows = nrows;
        rowset->ncols = ncols;
        rowset->data = (char **) mem_alloc(rowset->brows * ncols * sizeof(char *));
        rowset->name = (char **) mem_alloc(ncols * sizeof(char *));
        rowset->clen = (uint32_t *) mem_zeroalloc(ncols * sizeof(uint32_t));
        if (!rowset->data || !rowset->name || !rowset->clen) goto abort_rowset;
        
        buffer += bstart;
        
        // first buffer is guarantee to contains column names
        for (uint32_t i=0; i<ncols; ++i) {
            uint32_t cstart = 0;
            uint32_t len = internal_parse_number(&buffer[1], blen, &cstart);
            rowset->name[i] = buffer;
            buffer += cstart + len + 1;
            blen -= cstart + len + 1;
            if (rowset->clen[i] < len) rowset->clen[i] = len;
            if (rowset->maxlen < len) rowset->maxlen = len;
        }
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
            internal_set_error(connection, 1, "Unable to allocate memory: %d.", blen);
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
        uint32_t tlen = internal_parse_number(&buffer[1], blen-1, &cstart1);
        uint32_t clen = internal_parse_number(&buffer[cstart1 + 1], blen-1, &cstart2);
        uint32_t ulen = internal_parse_number(&buffer[cstart1 + cstart2 + 1], blen-1, &cstart3);
        
        // start of compressed buffer
        zdata = &buffer[tlen - clen + cstart1 + 1];
        
        // start of raw uncompressed header
        char *hstart = &buffer[cstart1 + cstart2 + cstart3 + 1];
        
        // try to allocate a buffer big enough to hold uncompressed data + raw header
        long clonelen = ulen + (zdata - hstart) + 1;
        char *clone = mem_alloc (clonelen);
        if (!clone) {
            internal_set_error(connection, 1, "Unable to allocate memory to uncompress buffer: %d.", clonelen);
            if (!isstatic && !externalbuffer) mem_free(buffer);
            return NULL;
        }
        
        // copy raw buffer
        memcpy(clone, hstart, zdata - hstart);
        
        // uncompress buffer and sanity check the result
        uint32_t rc = LZ4_decompress_safe(zdata, clone + (zdata - hstart), clen, ulen);
        if (rc <= 0 || rc != ulen) {
            internal_set_error(connection, 1, "Unable to decompress buffer (err code: %d).", rc);
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
        case CMD_JSON: {
            // +LEN string
            uint32_t cstart = 0;
            uint32_t len = internal_parse_number(&buffer[1], blen-1, &cstart);
            SQCloudResType type = (buffer[0] == CMD_JSON) ? RESULT_JSON : RESULT_STRING;
            if (buffer[0] == CMD_ZEROSTRING) --len;
            else if (buffer[0] == CMD_COMMAND) return internal_run_command(connection, &buffer[cstart+1], len, true);
            else if (buffer[0] == CMD_PUBSUB) return internal_setup_pubsub(connection, &buffer[cstart+1], len);
            else if (buffer[0] == CMD_RECONNECT) return internal_reconnect(connection, &buffer[cstart+1], len);
            SQCloudResult *res = internal_rowset_type(connection, buffer, len, cstart+1, type);
            if (res) res->externalbuffer = externalbuffer;
            return res;
        }
            
        case CMD_ERROR: {
            // -LEN ERRCODE ERRMSG
            uint32_t cstart = 0, cstart2 = 0;
            uint32_t len = internal_parse_number(&buffer[1], blen-1, &cstart);
            
            uint32_t errcode = internal_parse_number(&buffer[cstart + 1], blen-1, &cstart2);
            connection->errcode = (int)errcode;
            
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
            
            internal_parse_number(&buffer[1], blen-1, &cstart1); // parse len (already parsed in blen parameter)
            uint32_t idx = (buffer[0] == CMD_ROWSET) ? 0 : internal_parse_number(&buffer[cstart1 + 1], blen-1, &cstart2);
            uint32_t nrows = internal_parse_number(&buffer[cstart1 + cstart2 + 1], blen-1, &cstart3);
            uint32_t ncols = internal_parse_number(&buffer[cstart1 + cstart2 + + cstart3 + 1], blen-1, &cstart4);
            
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
            return &SQCloudResultNULL;
        }
    }
    
    if (!isstatic && !externalbuffer) mem_free(buffer);
    return NULL;
}

static bool internal_socker_forward_read (SQCloudConnection *connection, bool (*forward_cb) (char *buffer, size_t blen, void *xdata), void *xdata) {
    char sbuffer[8129];
    uint32_t blen = sizeof(sbuffer);
    uint32_t cstart = 0;
    uint32_t tread = 0;
    uint32_t clen = 0;
    
    char *buffer = sbuffer;
    char *original = buffer;
    int fd = connection->fd;
    
    while (1) {
        // perform read operation
        ssize_t nread = readsocket(fd, buffer, blen);
        
        // sanity check read
        if (nread < 0) {
            internal_set_error(connection, 1, "An error occurred while reading data: %s.", strerror(errno));
            goto abort_read;
        }
        
        if (nread == 0) {
            internal_set_error(connection, 1, "Unexpected EOF found while reading data: %s.", strerror(errno));
            goto abort_read;
        }
        
        // forward read to callback
        bool result = forward_cb(buffer, nread, xdata);
        if (!result) goto abort_read;
        
        // update internal counter
        tread += (uint32_t)nread;
        
        // determine command length
        if (clen == 0) {
            clen = internal_parse_number (&original[1], tread-1, &cstart);
            
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
    char header[1024];
    char *buffer = (char *)&header;
    uint32_t blen = sizeof(header);
    uint32_t tread = 0;

    int fd = (mainfd) ? connection->fd : connection->pubsubfd;
    char *original = buffer;
    while (1) {
        ssize_t nread = readsocket(fd, buffer, blen);
        
        if (nread < 0) {
            internal_set_error(connection, 1, "An error occurred while reading data: %s.", strerror(errno));
            goto abort_read;
        }
        
        if (nread == 0) {
            internal_set_error(connection, 1, "Unexpected EOF found while reading data: %s.", strerror(errno));
            goto abort_read;
        }
        
        tread += (uint32_t)nread;
        blen -= (uint32_t)nread;
        buffer += nread;
        
        // parse buffer looking for command length
        uint32_t cstart = 0;
        uint32_t clen = 0;
        
        if (internal_has_commandlen(original[0])) {
            clen = internal_parse_number (&original[1], tread-1, &cstart);
            if (clen == 0) continue;
        }
        
        // check if read is complete
        // clen is the lenght parsed in the buffer
        // cstart is the index of the first space
        // +1 because we skipped the first character in the internal_parse_number function
        if (clen + cstart + 1 != tread) {
            // check buffer allocation and continue reading
            if (clen + cstart > blen) {
                char *clone = mem_alloc(clen + cstart + 1);
                if (!clone) {
                    internal_set_error(connection, 1, "Unable to allocate memory: %d.", clen + cstart + 1);
                    goto abort_read;
                }
                memcpy(clone, original, tread);
                buffer = original = clone;
                blen = (clen + cstart + 1) - tread;
                buffer += tread;
            }
            
            continue;
        }
        
        // command is complete so parse it
        return internal_parse_buffer(connection, original, tread, (clen) ? cstart : 0, (original == header), false);
    }
    
abort_read:
    if (original != (char *)&header) mem_free(original);
    return NULL;
}

static bool internal_socket_write (SQCloudConnection *connection, const char *buffer, size_t len, bool mainfd) {
    size_t written = 0;
    
    int fd = (mainfd) ? connection->fd : connection->pubsubfd;
    
    // write header
    char header[32];
    int hlen = snprintf(header, sizeof(header), "+%zu ", len);
    ssize_t n = writesocket(fd, header, hlen);
    if (n != hlen) return internal_set_error(connection, 1, "An error occurred while writing header data: %s.", strerror(errno));
    
    // write buffer
    while (written < len) {
        ssize_t nwrote = writesocket(fd, buffer, len);
        //printf("writesocket connfd:%d nwrote:%d", fd, nwrote);
        
        if (nwrote < 0) {
            return internal_set_error(connection, 1, "An error occurred while writing data: %s.", strerror(errno));
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
    
    if (config->username && config->password) {
        char buffer[1024];
        snprintf(buffer, sizeof(buffer), "AUTH USER %s PASS %s", config->username, config->password);
        SQCloudResult *res = internal_run_command(connection, buffer, strlen(buffer), true);
        if (res != &SQCloudResultOK) return false;
    }
    
    if (config->database) {
        char buffer[1024];
        snprintf(buffer, sizeof(buffer), "USE DATABASE %s", config->database);
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
    hints.ai_family = (config) ? config->family : AF_UNSPEC;
    hints.ai_socktype = SOCK_STREAM;
    
    // get the address information for the server using getaddrinfo()
    char port_string[256];
    snprintf(port_string, sizeof(port_string), "%d", port);
    int rc = getaddrinfo(hostname, port_string, &hints, &addr_list);
    if (rc != 0 || addr_list == NULL) {
        return internal_set_error(connection, 1, "Error while resolving getaddrinfo (host %s not found).", hostname);
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
        return internal_set_error(connection, 1, "An error occurred while trying to connect: %s.", strerror(errno));
    }
    
    // bail if there was a timeout
    if ((time(NULL) - start) >= connect_timeout) {
        return internal_set_error(connection, 1, "Connection timeout while trying to connect (%d).", connect_timeout);
    }
    
    // turn off non-blocking
    int ioctl_blocking = 0;    /* ~0; //TRUE; */
    ioctl(sockfd, FIONBIO, &ioctl_blocking);
    
    // SSL on sockfd
    if (!internal_setup_ssl(connection, config)) return false;
    
    // finalize connection
    if (mainfd) {
        connection->fd = sockfd;
        connection->port = port;
        connection->hostname = strdup(hostname);
    } else {
        connection->pubsubfd = sockfd;
    }
    return true;
}

// MARK: - URL -

static int url_extract_username_password (const char *s, char b1[512], char b2[512]) {
    // user:pass@host.com:port/dbname?timeout=10&key2=value2&key3=value3
    
    // lookup username (if any)
    char *username = strchr(s, ':');
    if (!username) return 0;
    size_t len = username - s;
    if (len > 511) return -1;
    memcpy(b1, s, len);
    b1[len] = 0;
    
    // lookup username (if any)
    char *password = strchr(s, '@');
    if (!password) return 0;
    len = password - username - 1;
    if (len > 511) return -1;
    memcpy(b2, username+1, len);
    b2[len] = 0;
    
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
    
    // lookup username (if any)
    char *value = strchr(s, '&');
    if (!value) value = strchr(s, 0);
    if (!value) return 0;
    len = value - key - 1;
    if (len > 511) return -1;
    memcpy(b2, key+1, len);
    b2[len] = 0;
    
    return (int)(value - s) + 1;
}

// MARK: - RESERVED -

bool _reserved1 (SQCloudConnection *connection, const char *command, bool (*forward_cb) (char *buffer, size_t blen, void *xdata), void *xdata) {
    if (!forward_cb) return false;
    if (!internal_socket_write(connection, command, strlen(command), true)) return false;
    if (!internal_socker_forward_read(connection, forward_cb, xdata)) return false;
    return true;
}

SQCloudResult *_reserved2 (SQCloudConnection *connection, const char *UUID) {
    char command[512];
    snprintf(command, sizeof(command), "SET CLIENT UUID TO %s", UUID);
    return internal_run_command(connection, command, strlen(command), true);
}

SQCloudResult *_reserved3 (char *buffer, uint32_t blen, uint32_t cstart, SQCloudResult *chunk) {
    SQCloudConnection connection = {0};
    connection._chunk = chunk;
    SQCloudResult *res = internal_parse_buffer(&connection, buffer, blen, cstart, false, true);
    return res;
}

uint32_t _reserved4 (char *buffer, uint32_t blen, uint32_t *cstart) {
    return internal_parse_number(buffer, blen, cstart);
}

// MARK: - PUBLIC -

SQCloudConnection *SQCloudConnect (const char *hostname, int port, SQCloudConfig *config) {
    internal_init();
    
    SQCloudConnection *connection = mem_zeroalloc(sizeof(SQCloudConnection));
    if (!connection) return NULL;
    
    if (internal_connect(connection, hostname, port, config, true)) {
        if (config) internal_connect_apply_config(connection, config);
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
    SQCloudConfig *config = NULL;
    SQCloudConfig sconfig;
    
    // lookup for optional username/password
    char username[512];
    char password[512];
    int rc = url_extract_username_password(&s[n], username, password);
    if (rc == -1) return NULL;
    if (rc) {
        sconfig.username = string_dup(username);
        sconfig.password = string_dup(password);
        config = &sconfig;
    }
    
    // lookup for mandatory hostname
    n += rc;
    char hostname[512];
    char port_s[512];
    rc = url_extract_hostname_port(&s[n], hostname, port_s);
    if (rc <= 0) return NULL;
    int port = (int)strtol(port_s, NULL, 0);
    if (port == 0) port = SQCLOUD_DEFAULT_PORT;
    
    // lookup for optional database
    n += rc;
    char database[512];
    rc = url_extract_database(&s[n], database);
    if (rc == -1) return NULL;
    if (rc > 0) {
        sconfig.database = string_dup(database);
        config = &sconfig;
    }
    
    // lookup for optional key(s)/value(s)
    n += rc;
    char key[512];
    char value[512];
    while ((rc = url_extract_keyvalue(&s[n], key, value)) > 0) {
        if (strcasecmp(key, "timeout")) {
            int timeout = (int)strtol(value, NULL, 0);
            sconfig.timeout = (timeout) ? timeout : 0;
            config = &sconfig;
        }
        n += rc;
    }
    
    return SQCloudConnect(hostname, port, config);
}

SQCloudResult *SQCloudExec (SQCloudConnection *connection, const char *command) {
    return internal_run_command(connection, command, strlen(command), true);
}

void SQCloudDisconnect (SQCloudConnection *connection) {
    if (!connection) return;
    
    // free SSL
    
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
        internal_set_error(connection, 1, "A PubSub callback must be set before executing a PUBSUB ONLY command.");
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

int SQCloudErrorCode (SQCloudConnection *connection) {
    return (connection) ? connection->errcode : 666;
}

const char *SQCloudErrorMsg (SQCloudConnection *connection) {
    return (connection) ? connection->errmsg : "Not enoght memory to allocate a SQCloudConnection.";
}

// MARK: -

SQCloudResType SQCloudResultType (SQCloudResult *result) {
    return (result) ? result->tag : RESULT_ERROR;
}

bool SQCloudResultIsOK (SQCloudResult *result) {
    return (result == &SQCloudResultOK);
}

uint32_t SQCloudResultLen (SQCloudResult *result) {
    return (result) ? result->blen : 0;
}

char *SQCloudResultBuffer (SQCloudResult *result) {
    return (result) ? result->buffer : NULL;
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
        
        if (result->ischunk && !result->externalbuffer) {
            for (uint32_t i = 0; i<result->bcount; ++i) {
                if (result->buffers[i]) mem_free(result->buffers[i]);
            }
            mem_free(result->buffers);
        }
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
    *len = result->blen - (uint32_t)(result->data[row*result->ncols+col] - result->rawbuffer);
    return internal_parse_value(result->data[row*result->ncols+col], len, NULL);
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

void SQCloudRowsetDump (SQCloudResult *result, uint32_t maxline) {
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
    
    printf("Rows: %d - Cols: %d - Bytes: %d Time: %f secs", result->nrows, result->ncols, result->blen, result->time);
}
