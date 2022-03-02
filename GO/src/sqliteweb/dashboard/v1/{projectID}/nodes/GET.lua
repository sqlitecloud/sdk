-- List all nodes
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/nodes

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

userid = tonumber( userid )                                                                 -- Is string and comes from JWT. Contents is a number.

if projectID               == "auth"      then return error( 404, "Forbidden ProjectID" ) end -- fbf94289-64b0-4fc6-9c20-84083f82ee64
if string.len( projectID ) ~= 36          then return error( 400, "Invalid ProjectID" )   end 

Node = {
  id        = 0,                                -- NodeID - It is not good to have a simple int number!!!!!!
  name      = "",                               -- Name of this node
  type      = "",                               -- Type fo this node, for example: Leader, Worker
  provider  = "",                               -- Provider of this node
  image     = "",                               -- Image data for this node
  region    = "",                               -- Regin data for this node
  size      = "",                               -- Size info for this node
  address   = "",                               -- IPv[4,6] address or host name of this node
  port      = ""                                -- Port this node is listening on
}

Response = {
  status           = 0,                         -- status code: 0 = no error, error otherwise
  message          = "OK",                      -- "OK" or error message

  nodes            = {},                        -- Array with node objects
}

query = string.format( "SELECT NODE.id, NODE.name, type, provider, image, region, size, IIF( addr4, addr4, '' ) || IIF( addr4 AND addr6, ',', '' ) || IIF( addr6, addr6, '' ) AS address, port FROM USER JOIN PROJECT ON USER.id = PROJECT.user_id JOIN NODE ON PROJECT.uuid = NODE.project_uuid WHERE USER.enabled = 1 AND USER.id = %d AND uuid='%s';", userid, enquoteSQL( projectID ) )
databases = nil

if userid == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end

  nodes = getINIArray( projectID, "nodes", "" )
  for i = 1, #nodes do
    url = parseConnectionString( nodes[ i ] )
    if url then 
      if url.Port == 0 then url.Port = 8860 end

      Node = {}
      Node.id       = #Response.nodes -- TODO: This is not good
      Node.name     = string.format( "SQLiteCloud CORE Server node [%d]", #Response.nodes )
      Node.type     = getINIString( projectID, "type",      "Worker" )
      Node.provider = getINIString( projectID, "provider",  "On Premise" )
      Node.image    = getINIString( projectID, "image",     "Unknown" )
      Node.region   = getINIString( projectID, "type",      "On Premise" )
      Node.size     = getINIString( projectID, "type",      "Unknown" )
  
      Node.address  = url.Host
      Node.port     = url.Port
  
      Response.nodes[ #Response.nodes + 1 ] = Node
    end
  end

else

  nodes = executeSQL( "auth", query )
  if not nodes                            then return error( 404, "ProjectID not found" ) end
  if nodes.ErrorNumber              ~= 0  then return error( 502, "Bad Gateway" )         end
  if nodes.NumberOfColumns          ~= 9  then return error( 502, "Bad Gateway" )         end
  if nodes.NumberOfRows             <  1  then return error( 200, "OK" )                  end

  Response.nodes = nodes.Rows
  
end

if #Response.nodes == 0 then Response.nodes = nil end

SetStatus( 200 )
Write( jsonEncode( Response ) )