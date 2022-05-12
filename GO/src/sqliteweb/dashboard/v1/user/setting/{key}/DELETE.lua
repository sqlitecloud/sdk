--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Delete sessing key for 
--   ////                ///  ///                     logged in user
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/user/setting/key

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local key,       err, msg = checkParameter( key, 3 )                     if err ~= 0 then return error( err, string.format( msg, "key" ) )  end

if userID == 0 then         
  if not getINIBoolean( projectID, "enabled", false )                                then return error( 401, "Project Disabled" )           end
                                                                                          return error( 501, "Not Implemented" )
else
  local projectID, err, msg = verifyUserID( userID )                     if err ~= 0 then return error( err, msg )                          end

  result = executeSQL( "auth", string.format( "DELETE FROM UserSettings WHERE user_id = %d AND key = '%s'; SELECT changes() AS success;", userID, enquoteSQL( key ) ) )
  if not result                                                                      then return error( 504, "Gateway Timeout" )            end
  if result.ErrorMessage ~= ""                                                       then return error( 502, result.ErrorMessage )          end
  if result.ErrorNumber  ~= 0                                                        then return error( 502, "Bad Gateway" )                end
  if result.NumberOfRows ~= 1                                                        then return error( 502, "Bad Gateway" )                end
  if result.NumberOfColumns ~= 1                                                     then return error( 502, "Bad Gateway" )                end
  if result.Rows[ 1 ].success ~= 1                                                   then return error( 500, "Key not found" )              end

end

error( 200, "OK" )