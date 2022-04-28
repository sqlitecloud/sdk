--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Create a new setting with 
--   ////                ///  ///                     key and value
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/{nodeID}/setting/{key}

-- TODO: Modernize + use INSERT OR UPDATE

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local nodeID,    err, msg = checkNodeID( nodeID )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end
local key,       err, msg = checkParameter( key, 3 )                     if err ~= 0 then return error( err, string.format( msg, "key" ) )  end
local value,     err, msg = getBodyValue( "value", 0 )                   if err ~= 0 then return error( err, msg )                          end

query  = string.format( "INSERT OR REPLACE INTO NODE_SETTINGS ( node_id, key, value ) SELECT NODE.id, '%s', '%s' FROM USER JOIN PROJECT ON USER.id = PROJECT.user_id JOIN NODE ON NODE.project_uuid = PROJECT.uuid WHERE USER.enabled = 1 AND USER_id = %d AND NODE.id = %d;", enquoteSQL( key ), enquoteSQL( value ), userID, nodeID )

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
else

  local projectID, err, msg = verifyProjectID( userID, projectID )      if err ~= 0  then return error( err, msg )                          end
  local machineNodeID, err, msg = verifyNodeID( userID, projectID, nodeID ) if err ~= 0 then return error( err, msg )                       end

  result = executeSQL( "auth", query )

  if not result                                                                     then return error( 404, "ProjectID not found" )         end
  if result.ErrorNumber       ~= 0                                                  then return error( 502, result.ErrorMessage )           end
  if result.NumberOfColumns   ~= 0                                                  then return error( 502, "Bad Gateway" )                 end
  if result.NumberOfRows      ~= 0                                                  then return error( 502, "Bad Gateway" )                 end
  if result.Value             ~= "OK"                                               then return error( 502, result.Value )                  end
end                                       

error( 200, "OK" )