--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/04/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Reset the password for the specified user
--   ////                ///  ///                     if the token is valid
--     ////     //////////   ///                      
--        ////            ////          Requires    : Authentication
--          ////     /////              Output      : 
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/admin/v1/user/{email}/recover

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local email,    err, msg = checkParameter( email, 3 )             if err ~= 0 then return error( err, string.format( msg, "email" ) )  end

local token,    err, msg = getBodyValue( "token", 0 )             if err ~= 0 then return error( err, msg )                                end
local password, err, msg = getBodyValue( "password", 6 )          if err ~= 0 then return error( err, msg )                                end



-- check the token
local token_timeout_minutes = 30
command = string.format("SELECT user_id FROM PasswordResetToken WHERE token = ? AND used = 0 AND julianday(DATETIME('now')) < julianday(creation_date, '+%d minutes')", token_timeout_minutes)
result = executeSQL( "auth", command, {hash(token)} )

if not result                                                                       then return error( 504, "Gateway Timeout" )              end
if result.ErrorMessage      ~= ""                                                   then return error( 502, result.ErrorMessage )            end
if result.ErrorNumber       ~= 0                                                    then return error( 403, "Could fetch user data" )        end
if result.NumberOfRows      ~= 1                                                    then return error( 401, "Invalid token" )                end
if result.NumberOfColumns   ~= 1                                                    then return error( 502, "Bad Gateway" )                  end                                

userid = result.Rows[ 1 ].user_id

command = "UPDATE User SET password = ? WHERE id = ?; SELECT changes() AS success;"
result = executeSQL( "auth", command, {hash(password), userid} )

if not result                                                                      then return error( 504, "Gateway Timeout" )                  end
if result.ErrorMessage      ~= ""                                                  then return error( 502, "Bad Gateway" )                end
if result.ErrorNumber       ~= 0                                                   then return error( 502, "Bad Gateway" )                      end
if result.NumberOfRows      ~= 1                                                   then return error( 502, "Bad Gateway" )                      end
if result.NumberOfColumns   ~= 1                                                   then return error( 502, "Bad Gateway" )                      end
if result.Rows[ 1 ].success ~= 1                                                   then return error( 400, "User not found" )              end

command = "UPDATE PasswordResetToken SET used = 1 WHERE token = ?"
result = executeSQL( "auth", command, {hash(token)} )

error( 200, "OK" )