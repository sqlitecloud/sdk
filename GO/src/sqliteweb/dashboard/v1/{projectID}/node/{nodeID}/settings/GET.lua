-- List all nodes
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/1/settings

userid = 1

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

tmp    = nodeID
nodeID = tonumber( nodeID )                                                                 -- Is string and comes from URL. Could be anything!
userid = tonumber( userid )                                                                 -- Is string and comes from JWT. Contents is a number.

if projectID                     == "auth"  then return error( 404, "Forbidden ProjectID" ) end -- fbf94289-64b0-4fc6-9c20-84083f82ee64
if string.len( projectID )       ~= 36      then return error( 400, "Invalid ProjectID" )   end 
if string.format( "%d", nodeID ) ~= tmp     then return error( 400, "Invalid NodeID" )      end
if nodeID                        < 0        then return error( 400, "Invalid NodeID" )      end

Setting = {
  key   = "",
  value = ""
}

Response = {
  status            = 0,                         -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  settings          = nil,                        -- Array with key value pairs
}

if userid == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end

  nodes = getINIArray( projectID, "nodes", "" )
  if not nodes                            then return error( 501, "Internal Server error" ) end
  if #nodes == 0                          then return error( 404, "ProjectID OR NodeID not found" ) end
  if nodeID >= #nodes                     then return error( 404, "ProjectID OR NodeID not found" ) end

else

  query = string.format( "SELECT key, value FROM USER JOIN PROJECT ON USER.id = PROJECT.user_id JOIN NODE ON PROJECT.uuid = NODE.project_uuid JOIN NODE_SETTINGS ON NODE.id = node_id WHERE USER.enabled = 1 AND USER.id = %d AND NODE.id = %d AND uuid='%s';", userid, nodeID, enquoteSQL( projectID ) )
  settings = executeSQL( "auth", query )

  if not settings                          then return error( 404, "ProjectID OR NodeID not found" ) end
  if settings.ErrorNumber            ~= 0  then return error( 502, "Bad Gateway" )                   end
  if settings.NumberOfColumns        ~= 2  then return error( 502, "Bad Gateway" )                   end
  if settings.NumberOfRows           ~= 1  then return error( 404, "ProjectID OR NodeID not found" ) end

  Response.settings = settings.Rows

end

SetStatus( 200 )
Write( jsonEncode( Response ) )