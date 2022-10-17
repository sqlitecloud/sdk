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

-- https://localhost:8443/admin/v1/user/{email}/setting/{key}

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local email,     err, msg = checkParameter( email, 3 )                    if err ~= 0 then return error( err, string.format( msg, "email" ) ) end
local key,       err, msg = checkParameter( key, 1 )                      if err ~= 0 then return error( err, string.format( msg, "key" ) )   end
 
local value,     err, msg = getBodyValue( "value", 0 )                    if err ~= 0 then return error( err, msg )                           end

query  = "INSERT OR REPLACE INTO UserSettings ( user_id, key, value ) SELECT id, ?, ? FROM User WHERE email = ?;"
-- print( query )
result = executeSQL( "auth", query, {key, value, email} )

if not result                                                                         then return error( 504, "Gateway Timeout" )             end
if result.ErrorNumber     ~= 0                                                        then return error( 502, "Bad Gateway" )                 end
if result.Value           ~= "OK"                                                     then return error( 502, result.Value )                  end

error( 200, "OK" )