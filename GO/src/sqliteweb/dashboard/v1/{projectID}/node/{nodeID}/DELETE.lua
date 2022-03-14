-- Delete Node (and all the settings for this node), also kill the virtual machine...
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/{nodeID}

-- TODO / MISSING. Kill the virtual machine

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                     end
local nodeID,    err, msg = checkNodeID( nodeID )                        if err ~= 0 then return error( err, msg )                     end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                     end

if userID == 0 then         
  if not getINIBoolean( projectID, "enabled", false )                                then return error( 401, "Project Disabled" )      end
                                                                                          return error( 501, "Not Implemented" )
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg )                     end
  local nodeID,    err, msg = verifyNodeID( userID, projectID, nodeID )  if err ~= 0 then return error( err, msg )                     end

  result = executeSQL( "auth", string.format( "BEGIN TRANSACTION; DELETE FROM NODE_SETTINGS WHERE node_id = %d; DELETE FROM NODE WHERE id = %d; END TRANSACTION;", nodeID, nodeID ) )
  if not result                                                                      then return error( 504, "Gateway Timeout" )       end
  if result.ErrorNumber ~= 0                                                         then return error( 502, result.ErrorMessage )     end
  if result.Value ~= "OK"                                                            then return error( 502, "Bad Gateway" )           end
end

error( 200, "OK" )