//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.2.1
//     //             ///   ///  ///    Date        : 2022/02/15
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description :
//   ////                ///  ///
//     ////     //////////   ///
//        ////            ////
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package main

import (
  "encoding/json"
  "fmt"
  "net"
  "net/http"
  "sqlitecloud"
  "strings"
  "time"
  "strconv"

  "github.com/golang-jwt/jwt"
)

type Credentials struct {
  Login     string
  Password  string
}

type AuthRequest struct {
  RequestID int
  Credentials
}

type Response struct {
  Status     int
  Message    string
}

type AuthServer struct {
  Realm         string
  JWTSecret     []byte
  JWTTTL        int64

  db            *sqlitecloud.SQCloud
  host          string
  port          int
  login         string
  password      string
  cert          string
}

func init() {
  initializeSQLiteWeb()

  SQLiteWeb.router.HandleFunc( "/dashboard/v1/auth", SQLiteWeb.Auth.auth ).Methods( "POST" )
  SQLiteWeb.router.HandleFunc( "/dashboard/v1/auth", SQLiteWeb.Auth.JWTAuth( SQLiteWeb.Auth.reAuth ) ).Methods( "GET" )
}

/*
 * return >0 = success: UserID
 * return -1 = invalid credentials
 * return -2 = wrong credentials
 * return -3 = Internal error: could not create connection to auth server
 * return -4 = Internal error: could not connect to auth server
 * return -5 = Internal error: could not send auth query
 * return -6 = Internal error: invalid response from db
*/
func (this *AuthServer) lookupUserID( Login string, Password string ) int64 {
  var res *sqlitecloud.Result = nil
  var err error

  Login    = strings.TrimSpace( Login )
  Password = strings.TrimSpace( Password )

  if Login    == "" { return -1 }
  if Password == "" { return -1 }

  Login    = sqlitecloud.SQCloudEnquoteString( Login )
  Password = sqlitecloud.SQCloudEnquoteString( Password )

  query   := fmt.Sprintf( "SELECT id FROM User WHERE email = '%s' AND password = '%s' AND enabled = 1 LIMIT 1;", Login, Password )

  if this.db != nil {
    if res, err = this.db.Select( query ); err != nil || res == nil {
      this.db.Close()
      this.db = nil
      res     = nil
    }
  }

  if this.db == nil {
    if this.db = sqlitecloud.New( this.cert, 10 ) ; this.db == nil {
      return -3
    }

    if err := this.db.Connect( this.host, this.port, this.login, this.password, "users.sqlite", 10, "NO", 0 ); err != nil {
      this.db.Close()
      this.db = nil
      return -4
    }
  }

  if res == nil {
    if res, err = this.db.Select( query ); err != nil {
      if res != nil { res.Free() }
      res = nil
    }
  }

  if res == nil {
    this.db.Close()
    this.db = nil
    return -5
  }

  defer res.Free()
  if res.GetNumberOfRows() != 1 || res.GetNumberOfColumns() != 1 {
    return -6
  }
  return res.GetInt64Value_( 0, 0 )
}

func (this *AuthServer) getAuthorization( Header http.Header ) ( string, error ) {
  switch {
    case      Header[ "Authorization" ] == nil: fallthrough
    case len( Header[ "Authorization" ] ) < 1:  return "", fmt.Errorf( "Authorization header not found" )
    default:
      for _, header := range Header[ "Authorization" ] {
        if strings.HasPrefix( header, "Bearer " ) { return header[ 7: ], nil }
      }
      return Header[ "Authorization" ][ 0 ], nil
  }
}
func (this *AuthServer) decodeClaims( tokenString string ) ( *jwt.StandardClaims, error ) {
  //SigningMethodHS256
  token, err := jwt.ParseWithClaims( tokenString, &jwt.StandardClaims{}, func( token *jwt.Token ) ( interface{}, error ) {
    if _, ok := token.Method.( *jwt.SigningMethodHMAC ); !ok  {
      return nil, fmt.Errorf( "Unexpected signing method: '%v'", token.Header[ "alg" ] )
    }
    return this.JWTSecret, nil
  } )

  switch {
  case err   != nil:    return nil, err
  case token == nil:    return nil, fmt.Errorf( "Could not parse token" )
  case !token.Valid:    return nil, fmt.Errorf( "Invalid token" )
  default:
    switch claims, ok := token.Claims.(*jwt.StandardClaims); {
    case !ok:           return nil, fmt.Errorf( "No claims found" )
    case claims == nil: return nil, fmt.Errorf( "Could not parse token claims" )
    default:            return claims, nil
  } }
}
func (this *AuthServer) verifyClaims( claims *jwt.StandardClaims, reader *http.Request ) error {
  now        := time.Now().Unix()
  ip, _, err := net.SplitHostPort( reader.RemoteAddr )
  uip        := net.ParseIP( ip )

  switch {
    case err    != nil                                : return err
    case uip    == nil                                : return fmt.Errorf( "Invalid ClientIP" )
    case claims == nil                                : return fmt.Errorf( "Nil Claims" )
    case !claims.VerifyAudience( uip.String(), true ) : return fmt.Errorf( "Invalid Audience" )
    case !claims.VerifyExpiresAt( now, true )         : return fmt.Errorf( "Clain has expired" )
    case !claims.VerifyIssuedAt( now, true )          : return fmt.Errorf( "Invalid Issue Date" )
    case !claims.VerifyIssuer( long_name, true )      : return fmt.Errorf( "Invalidf Issuer" )
    case !claims.VerifyNotBefore( now, true )         : return fmt.Errorf( "Claim from the future" )
    case claims.Subject != this.Realm                 : return fmt.Errorf( "Invalid SUbject" )
    default                                           : return nil
  }
}

func (this *AuthServer) JWTAuth( nextHandler http.HandlerFunc ) http.HandlerFunc {
  return func( writer http.ResponseWriter, reader *http.Request ) {

    switch token, err := SQLiteWeb.Auth.getAuthorization( reader.Header ); {
    case err   != nil                                   : fallthrough
    case token == ""                                    : this.challengeAuth( writer )
    default:
      switch claims, err := SQLiteWeb.Auth.decodeClaims( token ); {
      case err != nil                                   : fallthrough
      case this.verifyClaims( claims, reader ) != nil   : this.challengeAuth( writer )
      default                                           : nextHandler.ServeHTTP( writer, reader )
  } } }
}

func (this *AuthServer) GetUserID( request *http.Request ) ( int64, error ) {
  token, _  := SQLiteWeb.Auth.getAuthorization( request.Header )
  claims, _ := SQLiteWeb.Auth.decodeClaims( token )

  switch userID, err := strconv.ParseInt( claims.Id, 10, 64 ); {
  case err != nil : return -1, err
  case userID < 1 : return userID, fmt.Errorf( "Invalid UserID" )
  default         : return userID, nil
  }
}

////

func (this *AuthServer) auth( writer http.ResponseWriter, request *http.Request ) {
  this.cors( writer, request )

  var authRequest AuthRequest

  switch err := json.NewDecoder( request.Body ).Decode( &authRequest ); {
  case err != nil : this.sendError( writer, 1, err.Error(), http.StatusBadRequest )
  default         : this.authorize( writer, request, this.lookupUserID( authRequest.Login, authRequest.Password ) )
  }
}

func (this *AuthServer) reAuth( writer http.ResponseWriter, request *http.Request ) {
  this.cors( writer, request )

  token, _  := SQLiteWeb.Auth.getAuthorization( request.Header )
  claims, _ := SQLiteWeb.Auth.decodeClaims( token )

  switch userID, err := strconv.ParseInt( claims.Id, 10, 64 ); {
  case err != nil : this.sendError( writer, 3, err.Error(), http.StatusBadRequest )
  case userID < 0 : this.sendError( writer, 4, "Invalid UserID", http.StatusBadRequest )
  default         : this.authorize( writer, request, userID )
  }
}

///

func (this *AuthServer) challengeAuth( writer http.ResponseWriter ) {
  // writer.Header().Set( "WWW-Authenticate", fmt.Sprintf( "Bearer realm=\"%s\"", this.Realm ) )
  writer.Header().Set( "WWW-Authenticate", fmt.Sprintf( "realm=\"%s\"", this.Realm ) )
  writer.WriteHeader( http.StatusUnauthorized )
}

/*
 * error codes: 0 = ok
 *              1 = bad request / could not parse json
 *              2 = invalid Cliaten id
 *              3 = wrong credentials (invalid/wrong format)
 *              4 = wrong credentials (not found on auth server)
 *              5 = insternal server error
*/
func (this *AuthServer) authorize( writer http.ResponseWriter, request *http.Request, userID int64 ) {
  response := Response {
    Status:     5,
    Message:    "Internal Server Error",
  }

  now      := time.Now().Unix()
  ip, _, _ := net.SplitHostPort( request.RemoteAddr )
  uip      := net.ParseIP( ip )

  switch {
  case uip == nil:
    response.Status  = 2
    response.Message = "Invalid ClientIP"
    writer.WriteHeader( http.StatusBadRequest )

  case userID == -1:
    response.Status  = 3
    response.Message = "Wrong Credentials"
    writer.WriteHeader( http.StatusBadRequest )

  case userID == -2:
    response.Status  = 4
    response.Message = "Wrong Credentials"
    writer.WriteHeader( http.StatusUnauthorized )

  case userID == -3 || userID == -4 || userID == -5:
    writer.WriteHeader( http.StatusInternalServerError )

  default:

    claims := &jwt.StandardClaims {
      Audience:   uip.String(),
      ExpiresAt:  now + this.JWTTTL,
      Id:         fmt.Sprintf( "%d", userID ),
      IssuedAt:   now,
      Issuer:     long_name,
      NotBefore:  now,
      Subject:    this.Realm,
    }

    Token            := jwt.NewWithClaims( jwt.SigningMethodHS256, claims )
    TokenString, err := Token.SignedString( this.JWTSecret ) // = Header, Payload, Signature

    if err != nil {
      writer.WriteHeader( http.StatusInternalServerError )

    } else {
      response.Status  = 0
      response.Message = TokenString
    }
  }

  writer.Header().Set( "Content-Type", "application/json" )
  writer.Header().Set( "Content-Encoding", "utf-8" )

  if jResponse, err := json.Marshal( response ); err == nil {
    writer.Write( jResponse )
  } else {
    http.Error( writer, err.Error(), http.StatusInternalServerError )
  }
}

func (this *AuthServer) sendError( writer http.ResponseWriter, status int, message string, statusCode int ) {
  writer.Header().Set( "Content-Type", "application/json" )
  writer.Header().Set( "Content-Encoding", "utf-8" )
  writer.Write( []byte( fmt.Sprintf( "{\"Status\":%d,\"Message\":\"%s\"}", status, message ) ) )
  writer.WriteHeader( statusCode )
}