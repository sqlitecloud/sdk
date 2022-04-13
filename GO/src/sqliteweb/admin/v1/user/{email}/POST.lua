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

local email,    err, msg = checkParameter( email, 3 )                   if err ~= 0 then return error( err, string.format( msg, "email" ) )  end

local name,     err, msg = getBodyValue( "name", 2 )                    if err ~= 0 then return error( err, msg )                            end
local company,  err, msg = getBodyValue( "company", 0 )                 if err ~= 0 then return error( err, msg )                            end
local password, err, msg = getBodyValue( "password", 5 )                if err ~= 0 then return error( err, msg )                            end
local enabled,  err, msg = getBodyValue( "enabled", 1 )                 if err ~= 0 then return error( err, msg )                            end

enabled = bool( enabled )

query  = string.format( "SELECT id FROM USER WHERE email = '%s';", enquoteSQL( email ) )
userID = executeSQL( "auth", query )
if not userID                                                                        then return error( 504, "Gateway Timeout" )              end
if userID.ErrorNumber     ~= 0                                                       then return error( 502, "Bad Gateway" )                  end
if userID.NumberOfColumns ~= 1                                                       then return error( 502, "Bad Gateway" )                  end
if userID.NumberOfRows    ~= 1                                                       then return error( 404, "User not found" )               end

userID = userID.Rows[ 1 ].id


query  = string.format( "UPDATE OR FAIL USER_SETTINGS SET value = '%s' WHERE user_id = %d AND key = '%s';", enquoteSQL( value ), userID, enquoteSQL( key ) )
result = executeSQL( "auth", query )

if not result                                                                        then return error( 504, "Gateway Timeout" )              end
if result.ErrorNumber     ~= 0                                                       then return error( 502, "Bad Gateway" )                  end
if result.Value           ~= "OK"                                                    then return error( 502, result.Value )                   end

error( 200, "OK" )