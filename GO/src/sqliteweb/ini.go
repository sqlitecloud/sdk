//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.1.0
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

import "fmt"
import "crypto/md5"

func GetINIString( section string, key string, defaultValue string ) string {
  for _, s := range cfg.SectionStrings() {
    switch {
    case s != section                         : continue
    case cfg.Section( section ).HasKey( key ) : return cfg.Section( section ).Key( key ).MustString( defaultValue )
    default                                   : return defaultValue
  } }
  return defaultValue
}

func MD5( data string ) string {
  return fmt.Sprintf( "%x", md5.Sum( []byte( data ) ) )
}

func CheckCredentials( section string, email string, password string ) bool {
  switch {
  case GetINIString( section, "email", "" )    != email           : return false
  case GetINIString( section, "password", "" ) == password        : return true
  case GetINIString( section, "password", "" ) == MD5( password ) : return true
  default                                                         : return false
  }
}