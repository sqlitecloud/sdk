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

local email,     err, msg = checkParameter( email, 3 )                   if err ~= 0  then return error( err, string.format( msg, "email" ) )   end
local key,       err, msg = checkParameter( key, 1 )                     if err ~= 0  then return error( err, string.format( msg, "key" ) )     end

query  = string.format( "SELECT id FROM USER WHERE email = '%s';", enquoteSQL( email ) )
userID = executeSQL( "auth", query )
if not userID                                                                         then return error( 504, "Gateway Timeout" )               end
if userID.ErrorNumber     ~= 0                                                        then return error( 502, "Bad Gateway" )                   end
if userID.NumberOfColumns ~= 1                                                        then return error( 502, "Bad Gateway" )                   end
if userID.NumberOfRows    ~= 1                                                        then return error( 404, "User not found" )                end

userID = userID.Rows[ 1 ].id        

query  = string.format( "DELETE FROM USER_SETTINGS WHERE user_id = %d AND key = '%s';", userID, enquoteSQL( key ) )       
result = executeSQL( "auth", query )        

if not result                                                                         then return error( 504, "Gateway Timeout" )               end
if result.ErrorNumber     ~= 0                                                        then return error( 502, "Bad Gateway" )                   end
if result.Value           ~= "OK"                                                     then return error( 502, result.Value )                    end

error( 200, "OK" )