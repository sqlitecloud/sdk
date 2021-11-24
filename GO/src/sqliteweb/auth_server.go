//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.0.1
//     //             ///   ///  ///    Date        : 2021/11/17
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
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type Credentials struct {
  Login 		string
	Password 	string
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
	ExpiresAt					int64
	RequestsPerSecond	int64
	RequestLeft				int64
}

type AuthServer struct {
	JWTSecret			[]byte
	JWTTTL				int64

	Tokens 				map[string]TokenInfo
}


func init() {
	initializeSQLiteWeb()

	SQLiteWeb.router.HandleFunc( "/api/v1/auth", SQLiteWeb.Auth.auth ).Methods( "POST" )
	SQLiteWeb.router.HandleFunc( "/api/v1/auth", SQLiteWeb.Auth.JWTAuth( SQLiteWeb.Auth.unAuth ) ).Methods( "DELETE" )
}

func (this *AuthServer) login( Login string, Password string ) bool {
	return true
}

func headerHasAuthToken( Header http.Header ) bool {
	switch {
	case Header[ "Token" ] == nil: 			return false
	case len( Header[ "Token" ] ) < 1: 	return false
	default: 														return true
	}
}
func (this *AuthServer) getAuthTokenFromHeader( Header http.Header ) string {
	switch headerHasAuthToken( Header ) {
	case true: return Header[ "Token" ][ 0 ]
	default:   return ""
	}
}

func (this *AuthServer) unAuth( writer http.ResponseWriter, request *http.Request ) {
	response := Response{
		ResponseID: 0,
		Status: 	  -1,
		Message:    "ERROR: Token not found.",
	}

	if token := SQLiteWeb.Auth.getAuthTokenFromHeader( request.Header ); token != "" {
		delete( this.Tokens, token )
		response.Status  = 0;
		response.Message = "OK"
	} 

	if jResponse, err := json.Marshal( response ); err == nil {
		writer.Header().Set( "Content-Type", "application-json" )
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
	  switch t, ok := this.Tokens[ SQLiteWeb.Auth.getAuthTokenFromHeader( reader.Header ) ]; {
	  case !ok: 														fallthrough
	  case t.ExpiresAt < time.Now().Unix(): fallthrough
	  case t.RequestLeft < 1:
	  	writer.Header().Set( "WWW-Authenticate", fmt.Sprintf( "Bearer realm=\"%s\"", "api/v1/" ) )
	  	writer.WriteHeader( http.StatusUnauthorized )
	  default:
	  	t.RequestLeft--
	  	nextHandler.ServeHTTP( writer, reader )
	  }
	}
}

func (this *AuthServer) auth( writer http.ResponseWriter, request *http.Request ) {
	// Read JSON Packet
	var authRequest AuthRequest
	json.NewDecoder( request.Body ).Decode( &authRequest ); 

	// Read & Overwrite from (old) Token
	token := SQLiteWeb.Auth.getAuthTokenFromHeader( request.Header )
	if t, ok := this.Tokens[ token ]; ok {
		authRequest.Login    = t.Login
		authRequest.Password = t.Password
		delete( this.Tokens, token )
	}

	if authRequest.Login == "" || authRequest.Password == "" {
		writer.WriteHeader( http.StatusBadRequest )
		return
	}
	
	// Check credentials
	if !this.login( authRequest.Login, authRequest.Password ) {
		writer.WriteHeader( http.StatusUnauthorized )
		return
	}

	// Delete double logins
	for t, ti := range this.Tokens {
		if ti.Login == authRequest.Login && ti.Password == authRequest.Password {
			delete( this.Tokens, t )
		}
	}

	now    := time.Now().Unix()
	// claims := &CustomClaims{ 
	// 	 Credentials: Credentials{
	// 	 	Login: 			authRequest.Login, 
	// 	 	Password: 	authRequest.Password, 
	// 	 },
	// 	StandardClaims: jwt.StandardClaims {
	// 		Id: 				fmt.Sprintf( "%d", authRequest.RequestID ),
	// 		Issuer:     long_name,
	// 		IssuedAt: 	now,
	// 		NotBefore: 	now,
	// 		ExpiresAt: 	now + this.JWTTTL,
	// 		Subject: 		"api/v1/",
	// }	}

	claims := &jwt.StandardClaims {
			Id: 				fmt.Sprintf( "%d", authRequest.RequestID ),
			// Issuer:     long_name,
			IssuedAt: 	now,
			NotBefore: 	now,
			ExpiresAt: 	now + this.JWTTTL,
			Subject: 		"api/v1/",
	}

  Token            := jwt.NewWithClaims( jwt.SigningMethodHS256, claims )
	TokenString, err := Token.SignedString( this.JWTSecret ) // = Header, Payload, Signature

	if err != nil {
		writer.WriteHeader( http.StatusInternalServerError )
		return
	}

	response := Response{
		ResponseID: authRequest.RequestID,
		Status: 	  0,
		Message:    TokenString,
	}

	if jResponse, err := json.Marshal( response ); err == nil {
		// http.SetCookie( writer, &http.Cookie {
		// 	Name: "Token",
		// 	Expires: expirationTime,
		// 	Value: TokenString,
		// } )

		writer.Header().Set( "Content-Type", "application-json" )
		writer.Header().Set( "Content-Encoding", "utf-8" )
		writer.Write( jResponse )

		this.Tokens[ TokenString ] = TokenInfo {
			//Credentials: 				claims.Credentials,
			ExpiresAt:					now + this.JWTTTL,
			RequestsPerSecond: 	1000,
			RequestLeft: 				1000,
		}

	} else {
		http.Error( writer, err.Error(), http.StatusInternalServerError )
	}
}

// // https://fusionauth.io/blog/2021/02/18/securing-golang-microservice/
// func verifyJWT( Token string ) bool {
// 
// 	t, err := jwt.Parse( Token, func( token *jwt.Token ) ( publicKey interface{}, err error ) {
//     if _, ok := token.Method.(*jwt.SigningMethodHMAC ); !ok {
// 			return nil, fmt.Errorf( "Invalid Signing Method: %v", token.Header[ "alg" ] )
// 		}
// 
// 		if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
// 			return nil, fmt.Errorf( "Expired token" )
// 		}
// 
// 
// 		return publicKey, nil
// 	} )
// 
// 	return err == nil && t.Valid
// }
// 
// func JWTAuth( nextHandler http.HandlerFunc ) http.HandlerFunc {
// 	return func( writer http.ResponseWriter, reader *http.Request ) {
// 		token := SQLiteWeb.Auth.getAuthTokenFromHeader( reader.Header )
// 		if verifyJWT( token ) {
// 			nextHandler.ServeHTTP( writer, reader )
// 			return
// 		}
// 		writer.Header().Set( "WWW-Authenticate", fmt.Sprintf( "Bearer realm=\"%s\"", "api/v1/" ) )
// 		writer.WriteHeader( http.StatusUnauthorized )
// 	}
// }