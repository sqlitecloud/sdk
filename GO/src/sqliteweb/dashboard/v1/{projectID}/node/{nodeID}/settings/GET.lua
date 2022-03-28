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
  status            = 0,                         -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  settings          = nil,                        -- Array with key value pairs
}

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end

  nodes = getINIArray( projectID, "nodes", "" )
  if not nodes                            then return error( 501, "Internal Server error" ) end
  if #nodes == 0                          then return error( 404, "ProjectID OR NodeID not found" ) end
  if nodeID >= #nodes                     then return error( 404, "ProjectID OR NodeID not found" ) end

else

  query = string.format( "SELECT key, value FROM USER JOIN PROJECT ON USER.id = PROJECT.user_id JOIN NODE ON PROJECT.uuid = NODE.project_uuid JOIN NODE_SETTINGS ON NODE.id = node_id WHERE USER.enabled = 1 AND USER.id = %d AND NODE.id = %d AND uuid='%s';", userID, nodeID, enquoteSQL( projectID ) )
  settings = executeSQL( "auth", query )

  if not settings                          then return error( 404, "ProjectID OR NodeID not found" ) end
  if settings.ErrorNumber            ~= 0  then return error( 502, "Bad Gateway" )                   end
  if settings.NumberOfColumns        ~= 2  then return error( 502, "Bad Gateway" )                   end
  if settings.NumberOfRows           ~= 1  then return error( 404, "ProjectID OR NodeID not found" ) end

  Response.settings = settings.Rows

end

SetStatus( 200 )
Write( jsonEncode( Response ) )