-- Change value for setting key for logged in user
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/user/setting/{key}

-- TODO: Check if UPDATE WAS SUCCESSFULL, remove INSERT OR REPLACE

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
  local projectID, err, msg = verifyUserID( userID )                     if err ~= 0 then return error( err, msg )                          end

  -- result = executeSQL( "auth", string.format( "INSERT OR REPLACE INTO USER_SETTINGS ( user_id, key, value ) VALUES ( %d, '%s', '%s' );", userID, enquoteSQL( key ), enquoteSQL( value ) ) )
  result = executeSQL( "auth", string.format( "UPDATE USER_SETTINGS SET value = '%s' WHERE user_id = %d AND key = '%s';", enquoteSQL( value ), userID, enquoteSQL( key ) ) )
  if not result                                                                      then return error( 504, "Gateway Timeout" )            end
  if result.ErrorNumber ~= 0                                                         then return error( 502, result.ErrorMessage )          end
  if result.Value ~= "OK"                                                            then return error( 502, "Bad Gateway" )                end
end

error( 200, "OK" )