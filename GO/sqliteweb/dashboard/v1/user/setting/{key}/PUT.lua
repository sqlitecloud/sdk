--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/04/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Change value for setting 
--   ////                ///  ///                     key for logged in user
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end

local key,       err, msg = checkParameter( key, 3 )                     if err ~= 0 then return error( err, string.format( msg, "key" ) )  end
local value,     err, msg = getBodyValue( "value", 0 )                   if err ~= 0 then return error( err, msg )                          end

if userID == 0 then         
  if not getINIBoolean( projectID, "enabled", false )                                then return error( 401, "Project Disabled" )           end
                                                                                          return error( 501, "Not Implemented" )
else
  local uID, companyID, err, msg = verifyUserID( userID )                     if err ~= 0 then return error( err, msg )                          end

  query  = "UPDATE UserSettings SET value = ? WHERE user_id = ? AND key = ?; SELECT changes() AS success;"
  result = executeSQL( "auth", query, {value, userID, key} )
  if not result                                                                      then return error( 504, "Gateway Timeout" )            end
  if result.ErrorMessage ~= ""                                                       then return error( 502, result.ErrorMessage )          end
  if result.ErrorNumber  ~= 0                                                        then return error( 502, "Bad Gateway" )                end
  if result.NumberOfRows ~= 1                                                        then return error( 502, "Bad Gateway" )                end
  if result.NumberOfColumns ~= 1                                                     then return error( 502, "Bad Gateway" )                end
  if result.Rows[ 1 ].success ~= 1                                                   then return error( 500, "Key not found" )              end
end

error( 200, "OK" )