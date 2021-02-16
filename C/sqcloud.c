//
//  sqcloud.c
//
//  Created by Marco Bambini on 08/02/21.
//

#include "sqcloud.h"
#include <stdlib.h>
#include <string.h>
#include <stdarg.h>
#include <assert.h>

#ifdef _WIN32
#include <winsock2.h>
#include <ws2tcpip.h>
#pragma comment(lib, "Ws2_32.lib")
#include <Shlwapi.h>
#include <io.h>
#include <float.h>
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
#include <sys/time.h>
#include <arpa/inet.h>
#include <netinet/tcp.h>
#include <sys/ioctl.h>
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
#define mem_alloc(_s)                       calloc(1,_s)
#define mem_malloc(_s)                      malloc(_s)
#define mem_free(_s)                        free(_s)
#define MIN(a,b)                            (((a)<(b))?(a):(b))

#define MAX_SOCK_LIST                       6           // maximum number of socket descriptor to try to connect to
                                                        // this change is required to support IPv4/IPv6 connections
#define DEFAULT_TIMEOUT                     12          // default connection timeout in seconds

#define REPLY_OK                            "+2 OK"     // default OK reply
#define REPLY_OK_LEN                        5           // default OK reply string length

// MARK: -

struct SQCloudResult {
    SQCloudResType  tag;                    // TYPE_STRING, TYPE_ROWSET
    
    char            *buffer;                // buffer used by the user (it could be a ptr inside rawbuffer)
    char            *rawbuffer;             // ptr to the buffer to be freed
    uint32_t        balloc;                 // buffer allocation size
    uint32_t        blen;                   // buffer real length
    
    // used in TYPE_ROWSET
    uint32_t        nrows;
    uint32_t        ncols;
    char            **data;
    char            **name;
} _SQCloudResult;

struct SQCloudConnection {
    int             fd;
    char            errmsg[1024];
    int             errcode;
} _SQCloudConnection;

static SQCloudResult SQCloudResultOK = {RESULT_OK, NULL, 0, 0, 0};

// MARK: - PRIVATE -

static int socket_geterror (int fd) {
    int err;
    socklen_t errlen = sizeof(err);
    
    int sockerr = getsockopt(fd, SOL_SOCKET, SO_ERROR, (void *)&err, &errlen);
    if (sockerr < 0) return -1;
    
    return ((err == 0 || err == EINTR || err == EAGAIN || err == EINPROGRESS)) ? 0 : err;
}

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

static void internal_clear_error (SQCloudConnection *connection) {
    connection->errcode = 0;
    connection->errmsg[0] = 0;
}

static bool internal_auth (SQCloudConnection *connection, SQCloudConfig *config) {
    return true;
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
    if (!buffer || buffer[0] == '_') {
        *len = 0;
        if (cellsize) *cellsize = 2;
        return NULL;
    }
    
    uint32_t cstart = 0;
    uint32_t blen = internal_parse_number(&buffer[1], 12, &cstart);
    
    // handle decimal/float cases
    if ((buffer[0] == ':') || (buffer[0] == ',')) {
        *len = cstart - 1;
        if (cellsize) *cellsize = cstart + 1;
        return &buffer[1];
    }
    
    *len = blen;
    if (cellsize) *cellsize = cstart + blen + 1;
    return &buffer[1+cstart];
}

static SQCloudResult *internal_rowset_string(SQCloudConnection *connection, char *buffer, uint32_t blen, uint32_t bstart) {
    SQCloudResult *rowset = (SQCloudResult *)mem_alloc(sizeof(SQCloudResult));
    if (!rowset) {
        internal_set_error(connection, 1, "Unable to allocate memory for SQCloudResult: %d.", sizeof(SQCloudResult));
        return NULL;
    }
    
    rowset->tag = RESULT_STRING;
    rowset->buffer = &buffer[bstart];
    rowset->rawbuffer = buffer;
    rowset->blen = blen;
    rowset->balloc = blen;
    
    return rowset;
}

static SQCloudResult *internal_parse_rowset (SQCloudConnection *connection, char *buffer, uint32_t blen, uint32_t bstart, uint32_t nrows, uint32_t ncols) {
    SQCloudResult *rowset = (SQCloudResult *)mem_alloc(sizeof(SQCloudResult));
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
    if (!rowset->data) return NULL;
    
    buffer += bstart;
    
    // the first column contains names
    for (uint32_t i=0; i<ncols; ++i) {
        uint32_t cstart = 0;
        uint32_t len = internal_parse_number(&buffer[1], blen, &cstart);
        rowset->name[i] = buffer;
        buffer += cstart + len + 1;
    }
    
    for (uint32_t i=0; i<nrows * ncols; ++i) {
        uint32_t len, cellsize;
        char *value = internal_parse_value(buffer, &len, &cellsize);
        rowset->data[i] = (value) ? buffer : NULL;
        buffer += cellsize;
        // printf("%d) %.*s\n", i, len, value);
    }
    
    return rowset;
}

static SQCloudResult *internal_parse_buffer (SQCloudConnection *connection, char *buffer, uint32_t blen, uint32_t cstart, bool isstatic) {
    if (blen <= 1) return false;
    
    // try to check if it is a OK reply: +2 OK
    if ((blen == REPLY_OK_LEN) && (strncmp(buffer, REPLY_OK, REPLY_OK_LEN) == 0)) {
        return &SQCloudResultOK;
    }
    
    // if buffer is static (stack based allocation) then it must be duplicated
    if (buffer[0] != '-' && isstatic) {
        char *clone = mem_alloc(blen);
        if (!clone) {
            internal_set_error(connection, 1, "Unable to allocate memory: %d.", blen);
            return NULL;
        }
        memcpy(clone, buffer, blen);
        buffer = clone;
    }
    
    // parse reply
    switch (buffer[0]) {
        case '+': {
            // +LEN string
            uint32_t cstart = 0;
            uint32_t len = internal_parse_number(&buffer[1], blen-1, &cstart);
            return internal_rowset_string(connection, buffer, len, cstart+1);
        }
            break;
            
        case '-': {
            // -LEN ERRCODE ERRMSG
            uint32_t cstart = 0, cstart2 = 0;
            uint32_t len = internal_parse_number(&buffer[1], blen-1, &cstart);
            
            uint32_t errcode = internal_parse_number(&buffer[cstart + 1], blen-1, &cstart2);
            connection->errcode = (int)errcode;
            
            len -= cstart2;
            memcpy(connection->errmsg, &buffer[cstart + cstart2 + 1], MIN(len, sizeof(connection->errmsg)));
            connection->errmsg[len] = 0;
            return NULL;
        }
            break;
            
        case '*': {
            // *LEN ROWS COLS DATA
            uint32_t cstart1 = 0, cstart2 = 0, cstart3 = 0;
            internal_parse_number(&buffer[1], blen-1, &cstart1);
            
            uint32_t nrows = internal_parse_number(&buffer[cstart1 + 1], blen-1, &cstart2);
            uint32_t ncols = internal_parse_number(&buffer[cstart1 + cstart2 + 1], blen-1, &cstart3);
            
            uint32_t bstart = cstart1 + cstart2 + cstart3 + 1;
            // printf("%.*s\n\n", blen, buffer);
            return internal_parse_rowset(connection, buffer, blen, bstart, nrows, ncols);
        }
            break;
    }
    
    return NULL;
}

static SQCloudResult *internal_socket_read (SQCloudConnection *connection) {
    // most of the time one read will be sufficient
    char header[1024];
    char *buffer = (char *)&header;
    uint32_t blen = sizeof(header);
    uint32_t tread = 0;

    char *original = buffer;
    while (1) {
        ssize_t nread = readsocket(connection->fd, buffer, blen);
        
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
                char *clone = mem_malloc(clen + cstart + 1);
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
        return internal_parse_buffer(connection, original, tread, cstart, (original == header));
    }
    
abort_read:
    if (original != (char *)&header) mem_free(original);
    return NULL;
}

static bool internal_socket_write (SQCloudConnection *connection, const char *buffer, size_t len) {
    size_t written = 0;
    
    // write header
    char header[32];
    int hlen = snprintf(header, sizeof(header), "+%zu ", len);
    writesocket(connection->fd, header, hlen);
    
    // write buffer
    while (written < len) {
        ssize_t nwrote = writesocket(connection->fd, buffer, len);
        //printf("writesocket connfd:%d nwrote:%d", connection->fd, nwrote);
        
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

static bool internal_connect (SQCloudConnection *connection, const char *hostname, int port, SQCloudConfig *config) {
    // apparently a listening IPv4 socket can accept incoming connections from only IPv4 clients
    // so I must explicitly connect using IPv4 if I want to be able to connect with older cubeSQL versions
    // https://stackoverflow.com/questions/16480729/connecting-ipv4-client-to-ipv6-server-connection-refused
    
    // ipv4/ipv6 specific variables
    struct sockaddr_storage serveraddr;
    struct addrinfo hints, *addr_list = NULL, *addr;
    
    // ipv6 code from https://www.ibm.com/support/knowledgecenter/ssw_ibm_i_72/rzab6/xip6client.htm
    memset(&hints, 0x00, sizeof(hints));
    hints.ai_flags    = AI_NUMERICSERV;
    hints.ai_family   = AF_UNSPEC;
    hints.ai_socktype = SOCK_STREAM;
    
    // check if we were provided the address of the server using
    // inet_pton() to convert the text form of the address to binary form.
    // If it is numeric then we want to prevent getaddrinfo() from doing any name resolution.
    int rc = inet_pton(AF_INET, hostname, &serveraddr);
    if (rc == 1) { /* valid IPv4 text address? */
        hints.ai_family = AF_INET;
        hints.ai_flags |= AI_NUMERICHOST;
    }
    else {
        rc = inet_pton(AF_INET6, hostname, &serveraddr);
        if (rc == 1) { /* valid IPv6 text address? */
            hints.ai_family = AF_INET6;
            hints.ai_flags |= AI_NUMERICHOST;
        }
    }
    
    // get the address information for the server using getaddrinfo()
    char port_string[256];
    snprintf(port_string, sizeof(port_string), "%d", port);
    rc = getaddrinfo(hostname, port_string, &hints, &addr_list);
    if (rc != 0 || addr_list == NULL) {
        return internal_set_error(connection, 1, "Error while resolving getaddrinfo (host %s not found).", hostname);
    }
    
    // begin non-blocking connection loop
    int sock_index = 0;
    int sock_current = 0;
    int sock_list[MAX_SOCK_LIST] = {0};
    for (addr = addr_list; addr != NULL; addr = addr->ai_next, ++sock_index) {
        if (sock_index >= MAX_SOCK_LIST) break;
        
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
    
    // authorize
    if (!internal_auth(connection, config)) return false;
    
    // finalize connection
    connection->fd = sockfd;
    return true;
}

// MARK: - PUBLIC -

SQCloudConnection *SQCloudConnect(const char *hostname, int port, SQCloudConfig *config) {
    internal_init();
    
    SQCloudConnection *connection = mem_alloc(sizeof(SQCloudConnection));
    if (!connection) return NULL;
    
    internal_connect(connection, hostname, port, config);
    return connection;
}

SQCloudResult *SQCloudExec(SQCloudConnection *connection, const char *command) {
    internal_clear_error(connection);
    
    if (!internal_socket_write(connection, command, strlen(command))) return false;
    return internal_socket_read(connection);
}

void SQCloudFree (SQCloudConnection *connection) {
    if (!connection) return;
    
    // free SSL
    // try to gracefully close connection
    
    if (connection->fd) closesocket(connection->fd);
    mem_free(connection);
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

uint32_t SQCloudResultLen (SQCloudResult *result) {
    return (result) ? result->blen : 0;
}

char *SQCloudResultBuffer (SQCloudResult *result) {
    return (result) ? result->buffer : NULL;
}

void SQCloudResultFree (SQCloudResult *result) {
    if (!result) return;
    if (result == &SQCloudResultOK) return;
    
    mem_free(result->rawbuffer);
    
    if (result->tag == RESULT_ROWSET) {
        mem_free(result->name);
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

static bool SQCloudRowSetSanityCheck (SQCloudResult *result, uint32_t row, uint32_t col) {
    if (!result || result->tag != RESULT_ROWSET) return false;
    if ((row >= result->nrows) || (col >= result->ncols)) return false;
    return true;
}

SQCloudValueType SQCloudRowSetValueType (SQCloudResult *result, uint32_t row, uint32_t col) {
    if (!SQCloudRowSetSanityCheck(result, row, col)) return VALUE_NULL;
    return internal_type(result->data[row*col]);
}

char *SQCloudRowSetColumnName (SQCloudResult *result, uint32_t col, uint32_t *len) {
    if (!result || result->tag != RESULT_ROWSET) return NULL;
    if (col >= result->ncols) return NULL;
    return internal_parse_value(result->name[col], len, NULL);
}

uint32_t SQCloudRowSetRows (SQCloudResult *result) {
    if (!SQCloudRowSetSanityCheck(result, 0, 0)) return 0;
    return result->nrows;
}

uint32_t SQCloudRowSetCols (SQCloudResult *result) {
    if (!SQCloudRowSetSanityCheck(result, 0, 0)) return 0;
    return result->ncols;
}

char *SQCloudRowSetValue (SQCloudResult *result, uint32_t row, uint32_t col, uint32_t *len) {
    if (!SQCloudRowSetSanityCheck(result, row, col)) return NULL;
    return internal_parse_value(result->data[row*col], len, NULL);
}

int32_t SQCloudRowSetInt32Value (SQCloudResult *result, uint32_t row, uint32_t col) {
    if (!SQCloudRowSetSanityCheck(result, row, col)) return 0;
    uint32_t len = 0;
    char *value = internal_parse_value(result->data[row*col], &len, NULL);
    
    char buffer[256];
    snprintf(buffer, sizeof(buffer), "%.*s", len, value);
    return (int32_t)strtol(buffer, NULL, 0);
}

int64_t SQCloudRowSetInt64Value (SQCloudResult *result, uint32_t row, uint32_t col) {
    if (!SQCloudRowSetSanityCheck(result, row, col)) return 0;
    uint32_t len = 0;
    char *value = internal_parse_value(result->data[row*col], &len, NULL);
    
    char buffer[256];
    snprintf(buffer, sizeof(buffer), "%.*s", len, value);
    return (int64_t)strtoll(buffer, NULL, 0);
}

float SQCloudRowSetFloatValue (SQCloudResult *result, uint32_t row, uint32_t col) {
    if (!SQCloudRowSetSanityCheck(result, row, col)) return 0.0;
    uint32_t len = 0;
    char *value = internal_parse_value(result->data[row*col], &len, NULL);
    
    char buffer[256];
    snprintf(buffer, sizeof(buffer), "%.*s", len, value);
    return (float)strtof(buffer, NULL);
}

double SQCloudRowSetDoubleValue (SQCloudResult *result, uint32_t row, uint32_t col) {
    if (!SQCloudRowSetSanityCheck(result, row, col)) return 0.0;
    uint32_t len = 0;
    char *value = internal_parse_value(result->data[row*col], &len, NULL);
    
    char buffer[256];
    snprintf(buffer, sizeof(buffer), "%.*s", len, value);
    return (double)strtod(buffer, NULL);
}

void SQCloudRowSetDump (SQCloudResult *result) {
    uint32_t nrows = result->nrows;
    uint32_t ncols = result->ncols;
    
    for (uint32_t i=0; i<ncols; ++i) {
        uint32_t len = 0;
        char *value = internal_parse_value(result->name[i], &len, NULL);
        printf("%-20.*s", len, value);
        
        bool newline = ((i % ncols == 1) || (ncols == 1));
        if (!newline) printf("|");
    }
    printf("\n");
    
    for (uint32_t i=0; i<ncols; ++i) {
        printf("--------------------");
        bool newline = ((i % ncols == 1) || (ncols == 1));
        if (!newline) printf("|");
    }
    printf("\n");
    
    for (uint32_t i=0; i<nrows * ncols; ++i) {
        uint32_t len = 0;
        char *value = internal_parse_value(result->data[i], &len, NULL);
        printf("%-20.*s", len, (value) ? value : "NULL");
        
        bool newline = ((i % ncols == 1) || (ncols == 1));
        (newline) ? printf("\n") : printf("|");
    }
    
}
