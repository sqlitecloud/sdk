--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/25
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : ping endpoint
--   ////                ///  ///       Requires    : Authentication
--     ////     //////////   ///        Output      : "PONG"
--        ////            ////
--          ////     /////
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

require "sqlitecloud"

do return error( 200, "PONG" ) end