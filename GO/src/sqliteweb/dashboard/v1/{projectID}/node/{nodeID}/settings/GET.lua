--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : List all nodes
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Structure with user settings
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/1/settings

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                     end
local nodeID,    err, msg = checkNodeID( nodeID )                        if err ~= 0 then return error( err, msg )                     end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                     end

Setting = {
  key   = "",
  value = ""
}

Response = {
  status            = 200,                       -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message
  value             = {},                        -- Array with key value pairs
}

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" )  end

  nodes = getINIArray( projectID, "nodes", "" )
  if not nodes                            then return error( 501, "Internal Server error" )         end
  if #nodes == 0                          then return error( 404, "ProjectID OR NodeID not found" ) end
  if nodeID >= #nodes                     then return error( 404, "ProjectID OR NodeID not found" ) end

else

  local projectID, err, msg = verifyProjectID( userID, projectID )                if err ~= 0 then return error( err, msg )                     end
  local machineNodeID,    err, msg = verifyNodeID( userID, projectID, nodeID )    if err ~= 0 then return error( err, msg )                     end

  query = string.format( "SELECT key, value FROM User JOIN Company ON User.company_id = Company.id JOIN Project ON Company.id = Project.company_id JOIN Node ON Project.uuid = Node.project_uuid JOIN NodeSettings ON Node.id = NodeSettings.node_id WHERE User.enabled = 1 AND User.id = %d AND Node.id = %d AND uuid='%s';", userID, nodeID, enquoteSQL( projectID ) )
  settings = executeSQL( "auth", query )

  if not settings                          then return error( 404, "ProjectID OR NodeID not found" ) end
  if settings.ErrorNumber            ~= 0  then return error( 502, "Bad Gateway" )                   end
  if settings.NumberOfColumns        ~= 2  then return error( 502, "Bad Gateway" )                   end
  -- if settings.NumberOfRows           ~= 1  then return error( 404, "ProjectID OR NodeID not found" ) end

  Response.value = settings.Rows

end

if #Response.value == 0 then Response.value = nil end

SetStatus( 200 )
Write( jsonEncode( Response ) )