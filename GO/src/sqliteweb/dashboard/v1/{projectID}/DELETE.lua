--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Delete Project (and all the 
--   ////                ///  ///                     nodes and node settings for this node)
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                     end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                     end

if userID == 0 then         
  if not getINIBoolean( projectID, "enabled", false )                                then return error( 401, "Project Disabled" )      end
                                                                                          return error( 501, "Not Implemented" )
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg )                     end
  
  local query  = string.format( "SELECT NODE.id AS nodeID FROM NODE JOIN PROJECT ON PROJECT.uuid = NODE.project_uuid JOIN USER ON USER.id = PROJECT.user_id WHERE USER.enabled = 1 AND USER.id = %d AND PROJECT.uuid = '%s';", userID, enquoteSQL( projectID ) )
  local result = executeSQL( "auth", query ) 

  if not result                                                                       then return -1, 503, "Service Unavailable"       end
  if result.ErrorNumber       ~= 0                                                    then return -1, 502, "Bad Gateway"               end
  if result.NumberOfColumns   ~= 1                                                    then return -1, 502, "Bad Gateway"               end 

  query = "BEGIN TRANSACTION;"
  for i = 1, result.NumberOfRows do
    nodeID = result.Rows[ i ].nodeID
    query  = string.format( "%s DELETE FROM NodeSettings WHERE node_id = %d; DELETE FROM Node WHERE id = %d;", query, nodeID, nodeID )
  end  
  query = string.format( "%s END TRANSACTION;", query )
  
  result = executeSQL( "auth", query )
  if not result                                                                      then return error( 504, "Gateway Timeout" )            end
  if result.ErrorMessage ~= ""                                                       then return error( 502, result.ErrorMessage )          end
  if result.ErrorNumber  ~= 0                                                        then return error( 502, "Bad Gateway" )                end
end

error( 200, "OK" )