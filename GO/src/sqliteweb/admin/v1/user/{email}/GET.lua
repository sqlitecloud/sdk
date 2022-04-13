--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/04/11
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : 
--   ////                ///  ///                     
--     ////     //////////   ///                      
--        ////            ////          Requires    : Authentication
--          ////     /////              Output      : 
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2


-- https://localhost:8443/admin/v1/user/{email}

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local email,       err, msg = checkParameter( email, 3 )                     if err ~= 0 then return error( err, string.format( msg, "email" ) )  end

query = string.format( "SELECT id, enabled, name, company, email, password, creation_date AS created, last_recovery_request AS recoveryRequest FROM USER WHERE email = '%s';", enquoteSQL( email ) )
user  = executeSQL( "auth", query )

if not user                                                                              then return error( 504, "Gateway Timeout" )              end
if user.ErrorNumber       ~= 0                                                           then return error( 502, "Bad Gateway" )                  end
if user.NumberOfColumns   ~= 8                                                           then return error( 502, "Bad Gateway" )                  end 
if user.NumberOfRows > 0                                                                 then user = user.Rows else user = nil                    end

Response = {
  status  = 200,
  message = "OK",
  user    = user
}

SetStatus( 200 )
Write( jsonEncode( Response ) )