--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Modify setting with key 
--   ////                ///  ///                     to value
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/{nodeID}/setting/{key}

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local nodeID,    err, msg = checkNodeID( nodeID )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

local key,       err, msg = checkParameter( key, 3 )                     if err ~= 0 then return error( err, string.format( msg, "key" ) )  end
local value,     err, msg = getBodyValue( "value", 0 )                   if err ~= 0 then return error( err, msg )                          end

if userID == 0 then         
  if not getINIBoolean( projectID, "enabled", false )                                then return error( 401, "Project Disabled" )           end
                                                                                          return error( 501, "Not Implemented" )      
else      
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg )                          end
  local machineNodeID, err, msg = verifyNodeID( userID, projectID, nodeID )  if err ~= 0 then return error( err, msg )                      end

  local settingID, err, msg = getNodeSettingsID( userID, projectID, nodeID, key ) 

  if err ~= 0 then
    query = "INSERT OR REPLACE INTO NodeSettings ( node_id, key, value ) VALUES( ?, ?, ? )"
    queryargs = {nodeID, key, value} 
  else 
    query = "UPDATE OR REPLACE NodeSettings SET value = ? WHERE id = ?;"
    queryargs = {body.value, settingID}
  end

  result = executeSQL( "auth", query, queryargs )
  if not result                                                                      then return error( 504, "Gateway Timeout" )            end
  if result.ErrorNumber ~= 0                                                         then return error( 502, result.ErrorMessage )          end
end

error( 200, "OK" )