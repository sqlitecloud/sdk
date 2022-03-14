-- Delete setting with key
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                     end
local nodeID,    err, msg = checkNodeID( nodeID )                        if err ~= 0 then return error( err, msg )                     end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                     end

if not key  or string.len( key  ) <  1                                               then return error( 400, "Missing Key" )           end

if userID == 0 then         
  if not getINIBoolean( projectID, "enabled", false )                                then return error( 401, "Project Disabled" )      end
                                                                                          return error( 501, "Not Implemented" )
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg )                     end
  local nodeID,    err, msg = verifyNodeID( userID, projectID, nodeID )  if err ~= 0 then return error( err, msg )                     end

  local settingID, err, msg = getNodeSettingsID( userID, projectID, nodeID, key ) 
  
  if err == 0 then 
    result = executeSQL( "auth", string.format( "DELETE FROM NODE_SETTINGS WHERE id = %d;", settingID ) )
    if not result                                                                    then return error( 504, "Gateway Timeout" )       end
    if result.ErrorNumber ~= 0                                                       then return error( 502, result.ErrorMessage )     end
  end
end

error( 200, "OK" )