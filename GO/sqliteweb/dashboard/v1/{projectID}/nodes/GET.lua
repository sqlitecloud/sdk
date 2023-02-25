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

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/nodes

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                     end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                     end

Node = {
  id            = 0,                                -- NodeID - It is not good to have a simple int number!!!!!!
  name          = "",                               -- Name of this node
  type          = "",                               -- Type fo this node, for example: Leader, Worker
  provider      = "",                               -- Provider of this node
  image         = "",                               -- Image data for this node
  region        = "",                               -- Region data for this node
  size          = "",                               -- Size info for this node
  address       = "",                               -- IPv[4,6] address or host name of this node
  port          = "",                               -- Port this node is listening on
  latitude      = 44.931,
  longitude     = 10.533,
  node_id       = 0,                                -- id of the node inside de cluster
  type          = "",                               -- Type fo this node, for example: Leader, Follower, Worker
  status        = "",                               -- progress status of the node, for example: Probe, Replicate, Snapshot (cluster) or Running (nocluster).
  match         = 0,                                -- is the index of the highest known matched raft entry (LIST NODES)
  last_activity = "",                               -- date and time of the last contact with a follower. Leader has NULL. (LIST NODES)
}

Response = {
  status           = 200,                           -- status code: 200 = no error, error otherwise
  message          = "OK",                          -- "OK" or error message
  value            = {},                            -- Array with node objects
}

if userID == 0 then         
  if not getINIBoolean( projectID, "enabled", false )                                then return error( 401, "Project Disabled" )      end

  nodes = getINIArray( projectID, "nodes", "" )
  for i = 1, #nodes do
    url = parseConnectionString( nodes[ i ] )
    if url then 
      if url.Port == 0 then url.Port = 8860 end

      Node = {}
      Node.id       = #Response.value -- TODO: This is not good
      Node.name     = string.format( "SQLiteCloud CORE Server node [%d]", #Response.value )
      Node.type     = getINIString( projectID, "type",      "Worker"     )
      Node.provider = getINIString( projectID, "provider",  "On Premise" )
      Node.image    = getINIString( projectID, "image",     "Unknown"    )
      Node.region   = getINIString( projectID, "type",      "On Premise" )
      Node.size     = getINIString( projectID, "type",      "Unknown"    )
  
      Node.address  = url.Host
      Node.port     = url.Port
  
      Response.value[ #Response.value + 1 ] = Node
    end
  end

else
  
  local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                     end
  local command     = "SELECT Node.id, Node.name, provider, image AS details, region, size, IIF( hostname NOT NULL AND hostname != '', hostname, IIF(addr6 NOT NULL AND addr6 != '', addr6, addr4)) AS address, port, latitude, longitude, node_id FROM User JOIN Company ON User.company_id = Company.id JOIN Project ON Company.id = Project.company_id JOIN Node ON Project.uuid = Node.project_uuid WHERE User.enabled = 1 AND User.id = ? AND uuid=? AND Node.active = 1;"

  local databases = nil

  nodes = executeSQL( "auth", command, {userID, projectID} )
  if not nodes                            then return error( 404, "ProjectID not found" ) end
  if nodes.ErrorNumber              ~= 0  then return error( 502, "Bad Gateway" )         end
  if nodes.NumberOfColumns          ~= 11 then return error( 502, "Bad Gateway" )         end
  if nodes.NumberOfRows             <  1  then return error( 200, "OK" )                  end

  listNodes = executeSQL( projectID,  "LIST NODES" )
  if not listNodes                            then return error( 404, "Can't list nodes" ) end
  if listNodes.ErrorNumber              ~= 0  then return error( 502, "Bad Gateway" )         end
  if listNodes.NumberOfColumns          ~= 7  then return error( 502, "Bad Gateway" )         end

  -- find the match value of the leader
  match_leader = 0
  for j = 1, listNodes.NumberOfRows do
    if listNodes.Rows[ j ].status == "Leader" then match_leader = listNodes.Rows[ j ].match end
  end

  for i = 1, nodes.NumberOfRows do
    cluster_node_id = nodes.Rows[ i ].node_id

    -- lookup which row from LIST NODES matches the cluster_node_id of the current row
    row = 0
    for j = 1, listNodes.NumberOfRows do
      if listNodes.Rows[ j ].id == cluster_node_id or nodes.NumberOfRows == 1 then row = j break end
    end

    --print("i, row: " .. i .. ", " .. row .. " (" .. cluster_node_id .. ", " .. nodes.Rows[ i ].name .. ", " .. listNodes.Rows[ row ].node ..  ")")

    nodes.Rows[ i ].match_leader = match_leader

    if row > 0 then
      nodes.Rows[ i ].type          = listNodes.Rows[ row ].status 
      nodes.Rows[ i ].status        = listNodes.Rows[ row ].progress 
      nodes.Rows[ i ].match         = listNodes.Rows[ row ].match
      nodes.Rows[ i ].last_activity = listNodes.Rows[ row ].last_activity 
    end
  end  

  Response.value = nodes.Rows
  
end

if #Response.value == 0 then Response.value = nil end

SetStatus( 200 )
Write( jsonEncode( Response ) )