--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/30
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Get a JSON with all 
--   ////                ///  ///                     providers, regions and 
--     ////     //////////   ///                      size parameters
--        ////            ////          Requires    : Authentication
--          ////     /////              Output      : Structure with node info
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/node/{nodeID}

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local nodeID,    err, msg = checkNodeID( nodeID )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

Response = {
  status            = 200,                       -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  value             = {
    id              = nodeID,                    -- NodeID - It is not good to have a simple int number!!!!!!
    name            = "",                        -- Name of this node
    type            = "",                        -- Type fo this node, for example: Leader, Worker
    provider        = "",                        -- Provider of this node
    details         = "?/?/?",                   -- "SFO1/1GB/25GB disk
    region          = "",                        -- Region data for this node
    size            = "",                        -- Size info for this node
    address         = "",                        -- IPv[4,6] address or host name of this node
    port            = 0,                         -- Port this node is listening on
    latitude        = 44.931,                    -- coordinates of the machine
    longitude       = 10.533,
 
    stats           = {},

    status          = "unknown",                 -- Replicating
    
    raft            = { 0, 0 },                  -- array 8960, 8960
  },
}

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end

-- TODO: Verify NodeID

  nodes = getINIArray( projectID, "nodes", "" )
  for i = 1, #nodes do
    url = parseConnectionString( nodes[ i ] )
    if url then 
      if url.Port == 0 then url.Port = 8860 end
      
      if nodeID == i - 1 then

        Response.value.id          = i
        Response.value.name        = string.format( "SQLiteCloud CORE Server node [%d]", i )
        Response.value.type        = getINIString( projectID, "type",      "Worker"        )
        Response.value.provider    = getINIString( projectID, "provider",  "On Premise"    )
        Response.value.details     = getINIString( projectID, "image",     "?/?/?"         )
        Response.value.region      = getINIString( projectID, "region",    "On Premise"    )
        Response.value.size        = getINIString( projectID, "type",      "Unknown"       )
        Response.value.address     = url.Host
        Response.value.port        = url.Port
        Response.value.latitude    = 44.931
        Response.value.longitude   = 10.533

        Response.value.stats       = {}

        -- get special values from the node
        Response.value.status      = "Unknown"
        
        Response.value.raft        = { 0, 0 }

        goto done
      end
    end
  end
                                          do  return error( 404, "NodeID not found" )               end

else
  
  local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                     end
  local machineNodeID,    err, msg = verifyNodeID( userID, projectID, nodeID )    if err ~= 0 then return error( err, msg )              end

  query = string.format( "SELECT Node.id, Node.node_id, Node.name, type, provider, image AS details, region, size, IIF( addr4, addr4, '' ) || IIF( addr4 AND addr6, ',', '' ) || IIF( addr6, addr6, '' ) AS address, port, latitude, longitude FROM User JOIN Company ON User.company_id = Company.id JOIN Project ON Company.id = Project.company_id JOIN Node ON Project.uuid = Node.project_uuid WHERE User.enabled = 1 AND User.id = %d AND Node.id = %d AND uuid='%s';", userID, nodeID, enquoteSQL( projectID ) )
  nodes = executeSQL( "auth", query )

  if not nodes                            then return error( 404, "ProjectID OR NodeID not found" ) end
  if nodes.ErrorNumber              ~= 0  then return error( 502, "Bad Gateway" )                   end
  if nodes.NumberOfColumns          ~= 12 then return error( 502, "Bad Gateway" )                   end
  if nodes.NumberOfRows             ~= 1  then return error( 404, "ProjectID OR NodeID not found" ) end

  Response.value             = nodes.Rows[ 1 ]  -- id, node_id, name, type, provider, image->details, region, size, address, port

  Response.value.stats       = {}

  Response.value.status      = "Unknown"
  Response.value.raft        = { 0, 0 }

  status = executeSQL( projectID, "LIST NODES;" )

  for i = 1, status.NumberOfRows do
    if status.Rows[ i ].status == "Leader" then Response.value.raft[ 2 ] = status.Rows[ i ].match end
    if i == machineNodeID then
      Response.value.status    = status.Rows[ i ].progress
      Response.value.raft[ 1 ] = status.Rows[ i ].match
      Response.value.type      = status.Rows[ i ].status
    end
  end

  query = string.format( "GET INFO LOAD,NUM_CLIENTS,DISK_USAGE_PERC NODE %d;", machineNodeID ) -- server_load, num_clients, cpu_time, mem_current, mem_max
  load = executeSQL( projectID, query )
  -- print("query:", query)

  Response.value.load = {
    load.Rows[ 2 ].ARRAY, -- num_clients
    load.Rows[ 1 ].ARRAY  -- server_load
  }

  query = string.format( "LIST STATS NODE %d;", machineNodeID )
  stats = executeSQL( projectID, query )
  
  for i = 1, stats.NumberOfRows do
    if not row                                  then row = { memory = { current = 0, max = 0 }, cpu = { sys = 0, user = 0 }, clients = { current = 0, max = 0 }, commands = 0, io = { reads = 0, writes = 0 }, bytes = { reads = 0, writes = 0 }, sampletime = "0000-00-00 00:00:00" } end
    if stats.Rows[ i ].key == "CPU_USAGE_SYS"   then row.cpu.user         = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "CPU_USAGE_USER"  then row.cpu.sys          = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "CURRENT_MEMORY"  then row.memory.current   = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "MAX_MEMORY"      then row.memory.max       = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "CURRENT_CLIENTS" then row.clients.current  = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "MAX_CLIENTS"     then row.clients.max      = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "NUM_COMMANDS"    then row.commands         = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "NUM_READS"       then row.io.reads         = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "NUM_WRITES"      then row.io.writes        = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "BYTES_IN"        then row.bytes.reads      = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "BYTES_OUT"       then row.bytes.writes     = stats.Rows[ i ].value
                                                     row.sampletime       = stats.Rows[ i ].datetime
      Response.value.stats[ #Response.value.stats + 1 ] = row
      row = nil
    end
  end

  if #Response.value.stats == 0 then Response.value.stats = nil end
end

::done::

SetStatus( 200 )
Write( jsonEncode( Response ) )