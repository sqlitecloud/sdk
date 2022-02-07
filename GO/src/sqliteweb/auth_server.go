//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.1.1
//     //             ///   ///  ///    Date        : 2021/12/20
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
  "fmt"
  "time"
  "net/http"
  "sqlitecloud"
  "encoding/json"
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
  ResponseID int
  Status     int
  Message    string
}

type CustomClaims struct {
  Credentials
  jwt.StandardClaims
} 

type TokenInfo struct {
  Credentials
  ExpiresAt         int64
  RequestsPerSecond int64
  RequestLeft       int64
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
  Tokens        map[string]TokenInfo
}

func init() {
  initializeSQLiteWeb()

  SQLiteWeb.router.HandleFunc( "/api/v1/auth", SQLiteWeb.Auth.auth ).Methods( "POST" )
  SQLiteWeb.router.HandleFunc( "/api/v1/auth", SQLiteWeb.Auth.JWTAuth( SQLiteWeb.Auth.unAuth ) ).Methods( "DELETE" )
}

func (this *AuthServer) getUserID( Login string, Password string ) int64 {
  var res *sqlitecloud.Result = nil
  var err error

  if this.db != nil {
    if res, err = this.db.Select( fmt.Sprintf( "SELECT id FROM User WHERE email = '%s' AND password = '%s' AND enabled = 1 LIMIT 1;", Login, Password ) ); err != nil || res == nil {
      this.db.Close()
      this.db = nil
      res     = nil
    }
  }

  if this.db == nil { 
    if this.db = sqlitecloud.New( this.cert, 10 ) ; this.db == nil { 
      return -1 
    }

    if err := this.db.Connect( this.host, this.port, this.login, this.password, "users.sqlite", 10, "NO", 0 ); err != nil {
      this.db.Close()
      this.db = nil
      return -2
    }
  }
  
  if res == nil {
    if res, err = this.db.Select( fmt.Sprintf( "SELECT id FROM User WHERE email = '%s' AND password = '%s' AND enabled = 1 LIMIT 1;", Login, Password ) ); err != nil {
      if res != nil { res.Free() }
      res = nil
    }
  }

  if res == nil {
    this.db.Close()
    this.db = nil
    return -3
  }

  defer res.Free()
  if res.GetNumberOfRows() != 1 || res.GetNumberOfColumns() != 1 {
    return 0
  }
  return res.GetInt64Value_( 0, 0 ) 
}

func (this *AuthServer) getAuthorization( Header http.Header ) string {
  switch {
    case      Header[ "Authorization" ] == nil: return ""
    case len( Header[ "Authorization" ] ) < 1:  return ""
    default:                                    return Header[ "Authorization" ][ 0 ] // or better the last?
  }
}

func (this *AuthServer) unAuth( writer http.ResponseWriter, request *http.Request ) {
  this.cors( writer, request )

  response := Response{
    ResponseID: 0,
    Status:     -1,
    Message:    "ERROR: Token not found.",
  }

  if token := SQLiteWeb.Auth.getAuthorization( request.Header ); token != "" {
    delete( this.Tokens, token )
    response.Status  = 0;
    response.Message = ""
  } 

  if jResponse, err := json.Marshal( response ); err == nil {
    writer.Header().Set( "Content-Type", "application/json" )
    writer.Header().Set( "Content-Encoding", "utf-8" )
    writer.Write( jResponse )
  } else {
    http.Error( writer, err.Error(), http.StatusInternalServerError )
  }
}

func (this *AuthServer) tokenExists( Token string ) bool {
  _, exits := this.Tokens[ Token ]
  return exits
}

func (this *AuthServer) JWTAuth( nextHandler http.HandlerFunc ) http.HandlerFunc {
  return func( writer http.ResponseWriter, reader *http.Request ) {
    switch t, ok := this.Tokens[ SQLiteWeb.Auth.getAuthorization( reader.Header ) ]; {
    case !ok:                             fallthrough
    case t.ExpiresAt < time.Now().Unix(): fallthrough
    case t.RequestLeft < 1:
      writer.Header().Set( "WWW-Authenticate", fmt.Sprintf( "Bearer realm=\"%s\"", this.Realm ) )
      writer.WriteHeader( http.StatusUnauthorized )
    default:
      t.RequestLeft--
      nextHandler.ServeHTTP( writer, reader )
    }
  }
}

func (this *AuthServer) auth( writer http.ResponseWriter, request *http.Request ) {
  this.cors( writer, request )

  // Read JSON Packet
  var authRequest AuthRequest
  json.NewDecoder( request.Body ).Decode( &authRequest ); 

  // Read & Overwrite from (old) Token
  token := SQLiteWeb.Auth.getAuthorization( request.Header )
  if t, ok := this.Tokens[ token ]; ok {
    authRequest.Login    = t.Login
    authRequest.Password = t.Password
    delete( this.Tokens, token )
  }
  
  response := Response {
    ResponseID: authRequest.RequestID,
    Status:     1,
    Message:    "Wrong Credentials",
  }

  if authRequest.Login == "" || authRequest.Password == "" {
    writer.WriteHeader( http.StatusBadRequest )
  
  } else {
    now    := time.Now().Unix()
    claims := &jwt.StandardClaims {
      Id:         "0",
      Issuer:     long_name,
      IssuedAt:   now,
      NotBefore:  now,
      ExpiresAt:  now + this.JWTTTL,
      Subject:    this.Realm,
    }
  
    // Check credentials
    if userID := this.getUserID( authRequest.Login, authRequest.Password ); userID < 1 {
      writer.WriteHeader( http.StatusUnauthorized )

    } else {
      claims.Id = fmt.Sprintf( "%d", userID )

      // Delete double logins
      for t, ti := range this.Tokens {
        if ti.Login == authRequest.Login && ti.Password == authRequest.Password {
          delete( this.Tokens, t )
        }
      }

      Token            := jwt.NewWithClaims( jwt.SigningMethodHS256, claims )
      TokenString, err := Token.SignedString( this.JWTSecret ) // = Header, Payload, Signature

      if err != nil {
        response.Status = 2
        response.Message = "Intenal Server Error"
        writer.WriteHeader( http.StatusInternalServerError )

      } else {
        response.Status  = 0
        response.Message = TokenString
      
        TokenString = "Bearer " + TokenString
        this.Tokens[ TokenString ] = TokenInfo {
          Credentials:        authRequest.Credentials,
          ExpiresAt:          now + this.JWTTTL,
          RequestsPerSecond:  1000,
          RequestLeft:        1000,
        }
      }
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