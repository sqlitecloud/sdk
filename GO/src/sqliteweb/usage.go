//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud CLI Application
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

var usage        = long_name + ` 

Usage:
  sqliteweb options
  sqliteweb -?|--help|--version

Examples:
  sqliteweb --config=../etc/sqliteweb.ini
  sqliteweb --version
  sqliteweb -?

General Options:
  --config=<PATH>          Use config file in <PATH> [default: /etc/sqliteweb/sqliteweb.ini]
  -?, --help               Show this screen
  --version                Display version information

Connection Options:
  -a, --address IP         Use IP address [default::0.0.0.0]
  -p, --port PORT          Use Port [default::8433]
  -c, --cert <FILE>        Use certificate chain in <FILE>        
  -k, --key <FILE>         Use private certificate key in <FILE>

Server Options:
  --www=<PATH>             Server static web sites from <PATH>
  --api=<PATH>             Server dummy REST stubs from <PATH>
`