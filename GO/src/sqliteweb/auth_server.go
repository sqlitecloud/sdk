//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.2.2
//     //             ///   ///  ///    Date        : 2022/02/18
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
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type AuthServer struct {
	Realm     string
	JWTSecret []byte
	JWTTTL    int64

	db       *sqlitecloud.SQCloud
	host     string
	port     int
	login    string
	password string
	cert     string
}

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	IPAddress string `json:"ipa,omitempty"`
	jwt.RegisteredClaims
}

func init() {
	initializeSQLiteWeb()

	SQLiteWeb.router.HandleFunc("/dashboard/v1/auth", SQLiteWeb.Auth.auth).Methods("POST")
	SQLiteWeb.router.HandleFunc("/dashboard/v1/auth", SQLiteWeb.Auth.JWTAuth(SQLiteWeb.Auth.getTokenFromAuthorization, SQLiteWeb.Auth.reAuth)).Methods("GET")
}

/*
 * return  0 = success: root User from .ini file
 * return >0 = success: UserID
 * return -1 = invalid credentials
 * return -2 = wrong credentials
 * return -3 = Internal error: could not create connection to auth server
 * return -4 = Internal error: could not connect to auth server
 * return -5 = Internal error: could not send auth query
 * return -6 = Internal error: invalid response from db
 */
func (this *AuthServer) lookupUserID(Login string, Password string) (int64, string, string) {
	Login = strings.TrimSpace(Login)
	Password = strings.TrimSpace(Password)

	if Login == "" {
		return -1, "", ""
	}
	if Password == "" {
		return -1, "", ""
	}

	if CheckCredentials("dashboard", Login, Password) {
		return 0, "", ""
	}

	Login = sqlitecloud.SQCloudEnquoteString(Login)
	Password = sqlitecloud.SQCloudEnquoteString(Password)

	query := fmt.Sprintf("SELECT id, first_name, last_name FROM User WHERE email = '%s' AND ( password = '%s' OR password = '%s' ) AND enabled = 1 LIMIT 1;", Login, Password, MD5(Password))

	if res, err, _, _ := cm.ExecuteSQL("auth", query); res != nil {
		defer res.Free()
		switch {
		case err != nil:
			return -5, "", ""
		case res.GetNumberOfRows() < 1:
			return -2, "", ""
		case res.GetNumberOfColumns() != 3:
			return -6, "", ""
		default:
			return res.GetInt64Value_(0, 0), res.GetStringValue_(0, 1), res.GetStringValue_(0, 2)
		}
	}
	return -7, "", ""
}

func (this *AuthServer) getTokenFromAuthorization(r *http.Request) (string, error) {
	switch {
	case r.Header["Authorization"] == nil:
		fallthrough
	case len(r.Header["Authorization"]) < 1:
		return "", fmt.Errorf("Authorization header not found")
	default:
		for _, header := range r.Header["Authorization"] {
			if strings.HasPrefix(header, "Bearer ") {
				return header[7:], nil
			}
		}
		return r.Header["Authorization"][0], nil
	}
}

type SqliteCloudToken struct {
	session SqliteCloudTokenSession
}

type SqliteCloudTokenSession struct {
	status    int
	message   string
	createdAt int64
	maxAge    int
}

func (this *AuthServer) getTokenFromCookie(r *http.Request) (string, error) {
	SQLiteWeb.Logger.Debug("getTokenFromCookie: cookie: ", r.Header["Cookie"])

	c, err := r.Cookie("sqlite-cloud-token")
	if err != nil {
		SQLiteWeb.Logger.Debug("getTokenFromCookie: error1: ", err)
		return "", fmt.Errorf("Authorization cookie not found")
	}
	// SQLiteWeb.Logger.Debug("getTokenFromCookie: found cookie ", c)
	
	// TODO: check other fields of the cookie

	return c.Value, nil
}

func (this *AuthServer) decodeClaims(tokenString string) (*Claims, error) {
	//SigningMethodHS256
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: '%v'", token.Header["alg"])
		}
		return this.JWTSecret, nil
	})

	switch {
	case err != nil:
		return nil, err
	case token == nil:
		return nil, fmt.Errorf("Could not parse token")
	case !token.Valid:
		return nil, fmt.Errorf("Invalid token")
	default:
		switch claims, ok := token.Claims.(*Claims); {
		case !ok:
			return nil, fmt.Errorf("No claims found")
		case claims == nil:
			return nil, fmt.Errorf("Could not parse token claims")
		default:
			return claims, nil
		}
	}
}
func (this *AuthServer) verifyClaims(claims *Claims, reader *http.Request) error {
	now := time.Now().Unix()
	ip, _, err := net.SplitHostPort(reader.RemoteAddr)
	// Note: IP validation has been diabled because the new dashboard calls backend endpoints from nodejs, node from the user's browser
	uip := net.ParseIP(ip)

	switch {
	case err != nil:
	case uip == nil:
		err = fmt.Errorf("Invalid ClientIP")
	case claims == nil:
		err = fmt.Errorf("Nil Claims")
	case !claims.VerifyAudience(service_name, true):
		err = fmt.Errorf("Invalid Audience")
	case !claims.VerifyExpiresAt(time.Unix(now, 0), true):
		err = fmt.Errorf("Claim has expired")
	case !claims.VerifyIssuedAt(time.Unix(now, 0), true):
		err = fmt.Errorf("Invalid Issue Date")
	case !claims.VerifyIssuer(jwt_issuer, true):
		err = fmt.Errorf("Invalid Issuer")
	case !claims.VerifyNotBefore(time.Unix(now, 0), true):
		err = fmt.Errorf("Claim from the future")
	// case claims.Subject != this.Realm                             : return fmt.Errorf( "Invalid Subject" )
	default:
		return nil
	}

	SQLiteWeb.Logger.Errorf("verifyClaims error: %s", err)
	return err
}

type TokenFunc func(*http.Request) (string, error)

func (this *AuthServer) JWTAuth(tokenFunc TokenFunc, nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, reader *http.Request) {
		switch token, err := tokenFunc(reader); {
		case err != nil:
			fallthrough
		case token == "":
			this.challengeAuth(writer)
		default:
			switch claims, err := SQLiteWeb.Auth.decodeClaims(token); {
			case err != nil:
				fallthrough
			case this.verifyClaims(claims, reader) != nil:
				this.challengeAuth(writer)
			default:
				nextHandler.ServeHTTP(writer, reader)
			}
		}
	}
}

func (this *AuthServer) GetUserID(tokenFunc TokenFunc, request *http.Request) (int64, error) {
	if token, err := tokenFunc(request); err == nil && token != "" {
		if claims, err := SQLiteWeb.Auth.decodeClaims(token); err == nil && claims != nil {
			switch userID, err := strconv.ParseInt(claims.Subject, 10, 64); {
			case err != nil:
				return -1, err
			case userID < 1:
				return userID, fmt.Errorf("Invalid UserID")
			default:
				return userID, nil
			}
		} else {
			return -1, err
		}
	} else {
		return -1, err
	}
}

////

func (this *AuthServer) auth(writer http.ResponseWriter, request *http.Request) {
	this.cors(writer, request)

	var credentials Credentials

	switch err := json.NewDecoder(request.Body).Decode(&credentials); {
	case err != nil:
		sendError(writer, err.Error(), http.StatusBadRequest)
	default:
		uid, fName, lName := this.lookupUserID(credentials.Login, credentials.Password)
		this.authorize(writer, request, uid, fName, lName)
	}
}

func (this *AuthServer) reAuth(writer http.ResponseWriter, request *http.Request) {
	this.cors(writer, request)

	token, _ := SQLiteWeb.Auth.getTokenFromAuthorization(request)
	claims, _ := SQLiteWeb.Auth.decodeClaims(token)

	switch userID, err := strconv.ParseInt(claims.Subject, 10, 64); {
	case err != nil:
		sendError(writer, err.Error(), http.StatusBadRequest)
	case userID < 0:
		sendError(writer, "Invalid UserID", http.StatusBadRequest)
	default:
		this.authorize(writer, request, userID, claims.FirstName, claims.LastName)
	}
}

///

func (this *AuthServer) challengeAuth(writer http.ResponseWriter) {
	// writer.Header().Set( "WWW-Authenticate", fmt.Sprintf( "Bearer realm=\"%s\"", this.Realm ) )
	writer.Header().Set("WWW-Authenticate", fmt.Sprintf("realm=\"%s\", error=\"invalid_token\"", this.Realm))
	writer.WriteHeader(http.StatusUnauthorized)
	writer.Write([]byte(fmt.Sprintf("{\"status\":%d,\"message\":\"%s\"}", http.StatusUnauthorized, "Invalid Token")))
}

/*
 * error codes: 0 = ok
 *              1 = bad request / could not parse json
 *              2 = invalid Cliaten id
 *              3 = wrong credentials (invalid/wrong format)
 *              4 = wrong credentials (not found on auth server)
 *              5 = insternal server error
 */
func (this *AuthServer) authorize(writer http.ResponseWriter, request *http.Request, userID int64, firstName string, lastName string) {
	response := Response{
		Status:  500,
		Message: "Internal Server Error",
	}

	now := time.Now().Unix()
	ip, _, _ := net.SplitHostPort(request.RemoteAddr)
	uip := net.ParseIP(ip)

	switch {
	case uip == nil:
		response.Status = 400
		response.Message = "Invalid ClientIP"
		writer.WriteHeader(http.StatusBadRequest)

	case userID == -1:
		response.Status = 400
		response.Message = "Wrong Credentials"
		writer.WriteHeader(http.StatusBadRequest)

	case userID == -2:
		response.Status = 400
		response.Message = "Wrong Credentials"
		writer.WriteHeader(http.StatusUnauthorized)

	case userID == -3 || userID == -4 || userID == -5:
		writer.WriteHeader(http.StatusInternalServerError)

	default:

		claims := &Claims{
			FirstName: firstName,
			LastName:  lastName,
			IPAddress: uip.String(),

			RegisteredClaims: jwt.RegisteredClaims{
				Audience:  []string{service_name},
				ExpiresAt: jwt.NewNumericDate(time.Unix(now+this.JWTTTL, 0)),
				IssuedAt:  jwt.NewNumericDate(time.Unix(now, 0)),
				Issuer:    jwt_issuer,
				NotBefore: jwt.NewNumericDate(time.Unix(now, 0)),
				Subject:   fmt.Sprintf("%d", userID),
			},
		}

		Token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		TokenString, err := Token.SignedString(this.JWTSecret) // = Header, Payload, Signature

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)

		} else {
			response.Status = 200
			response.Message = TokenString
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Content-Encoding", "utf-8")

	if jResponse, err := json.Marshal(response); err == nil {
		writer.Write(jResponse)
	} else {
		sendError(writer, "Internal Error", http.StatusInternalServerError)
		//http.Error( writer, err.Error(), http.StatusInternalServerError )
	}
}
