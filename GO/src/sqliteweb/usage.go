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
  sqliteweb --address=0.0.0.0 --port 8433 --stubs=api
  sqliteweb --version
  sqliteweb -?

General Options:
  --stubs <PATH>           Use PATH for dummy responses [default::api]
  --www <PATH>             Use PATH for REACT Web Sites [default::www]
  -?, --help               Show this screen
  --version                Display version information

Connection Options:
  -a, --address IP         Use IP address [default::0.0.0.0]
  -p, --port PORT          Use Port [default::8860]
`