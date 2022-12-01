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


-- https://localhost:8443/admin/v1/user/{email}/settings

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local email,       err, msg = checkParameter( email, 3 )                     if err ~= 0 then return error( err, string.format( msg, "email" ) )  end

command    = "SELECT key, value FROM User JOIN UserSettings ON USER.id = UserSettings.user_id WHERE User.email = ?;"
settings = executeSQL( "auth", command, {email} )

if not settings                                                                          then return error( 504, "Gateway Timeout" )              end
if settings.ErrorNumber       ~= 0                                                       then return error( 502, "Bad Gateway" )                  end
if settings.NumberOfColumns   ~= 2                                                       then return error( 502, "Bad Gateway" )                  end 
if settings.NumberOfRows > 0                                                             then settings = settings.Rows else settings = nil        end

Response = {
  status    = 200,
  message   = "OK",
  value  = settings,
}

SetStatus( 200 )
Write( jsonEncode( Response ) )