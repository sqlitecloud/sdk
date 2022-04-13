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

  node              = {
    id              = nodeID,                    -- NodeID - It is not good to have a simple int number!!!!!!
    name            = "",                        -- Name of this node
    type            = "",                        -- Type fo this node, for example: Leader, Worker
    provider        = "",                        -- Provider of this node
    -- image           = "",                        -- Image data for this node
    details         = "?/?/?",                   -- "SFO1/1GB/25GB disk
    region          = "",                        -- Regin data for this node
    size            = "",                        -- Size info for this node
    address         = "",                        -- IPv[4,6] address or host name of this node
    port            = 0,                         -- Port this node is listening on
    latitude        = 44.931,                    -- coordinates of the machine
    longitude       = 10.533,
 
    stats           = {},

    status          = "unknown",                 -- Replicating
    
    raft            = { 0, 0 },                  -- array 8960, 8960
    
    --load            = { 0, 0 },                  -- some load info
    --cpu             = { Used = 0, Total = 0 },   -- some cpu info
    --ram             = { Used = 0, Total = 0 },   -- some ram info
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

        Response.node.id          = i
        Response.node.name        = string.format( "SQLiteCloud CORE Server node [%d]", i )
        Response.node.type        = getINIString( projectID, "type",      "Worker"        )
        Response.node.provider    = getINIString( projectID, "provider",  "On Premise"    )
        Response.node.details     = getINIString( projectID, "image",     "?/?/?"         )
        Response.node.region      = getINIString( projectID, "region",    "On Premise"    )
        Response.node.size        = getINIString( projectID, "type",      "Unknown"       )
        Response.node.address     = url.Host
        Response.node.port        = url.Port
        Response.node.latitude    = 44.931
        Response.node.longitude   = 10.533

        Response.node.stats       = {}

        -- get special values from the node
        Response.node.status      = "Unknown"
        
        Response.node.raft        = { 0, 0 }

        -- Response.node.load        = { 0, 0 }
        -- Response.node.cpu         = { Used = 0, Total = 0 }
        -- Response.node.ram         = { Used = 0, Total = 0 }

        goto done
      end
    end
  end
                                          do  return error( 404, "NodeID not found" )               end

else
  
  local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                     end
  local nodeID,    err, msg = verifyNodeID( userID, projectID, nodeID )    if err ~= 0 then return error( err, msg )                     end

  query = string.format( "SELECT NODE.id, NODE.name, type, provider, image AS details, region, size, IIF( addr4, addr4, '' ) || IIF( addr4 AND addr6, ',', '' ) || IIF( addr6, addr6, '' ) AS address, port, latitude, longitude FROM USER JOIN PROJECT ON USER.id = PROJECT.user_id JOIN NODE ON PROJECT.uuid = NODE.project_uuid WHERE USER.enabled = 1 AND USER.id = %d AND NODE.id = %d AND uuid='%s';", userID, nodeID, enquoteSQL( projectID ) )
  nodes = executeSQL( "auth", query )

  if not nodes                            then return error( 404, "ProjectID OR NodeID not found" ) end
  if nodes.ErrorNumber              ~= 0  then return error( 502, "Bad Gateway" )                   end
  if nodes.NumberOfColumns          ~= 11 then return error( 502, "Bad Gateway" )                   end
  if nodes.NumberOfRows             ~= 1  then return error( 404, "ProjectID OR NodeID not found" ) end

  Response.node             = nodes.Rows[ 1 ]  -- id, name, type, provider, image->details, region, size, address, port

  Response.node.stats       = {}

  Response.node.status      = "Unknown"
  Response.node.raft        = { 0, 0 }

  -- Response.node.load        = { 0, 0 }
  -- Response.node.cpu         = {}
  -- Response.node.ram         = {}

  ------

  status = executeSQL( projectID, "LIST NODES;" )

  for i = 1, status.NumberOfRows do
    if status.Rows[ i ].status == "Leader" then Response.node.raft[ 2 ] = status.Rows[ i ].match end
    if i == nodeID then
      Response.node.status    = status.Rows[ i ].progress
      Response.node.raft[ 1 ] = status.Rows[ i ].match
      Response.node.type      = status.Rows[ i ].status
    end
  end

  query = string.format( "GET LOAD DETAILED NODE %d;", nodeID ) -- server_load, num_clients, cpu_time, mem_current, mem_max
  load = executeSQL( projectID, query )
  Response.node.load = {
    load.Rows[ 2 ].ARRAY, -- num_clients
    load.Rows[ 1 ].ARRAY  -- server_load
  }

  query = string.format( "LIST STATS NODE %d;", nodeID )
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
      Response.node.stats[ #Response.node.stats + 1 ] = row
      row = nil
    end
  end

  if #Response.node.stats == 0 then Response.node.stats = nil end

  --query = string.format( "LIST STATS NODE %d;", nodeID )
  --stats = executeSQL( projectID, query )
  --cpu = { Used = 0, Total = 0 }
  --ram = { Used = 0, Total = 0 }
  --for i = 1, stats.NumberOfRows do
  --  if stats.Rows[ i ].key == "CPU_USAGE_SYS"   then cpu.Used  = stats.Rows[ i ].value            end
  --  if stats.Rows[ i ].key == "CPU_USAGE_USER"  then cpu.Total = stats.Rows[ i ].value + cpu.Used end
  --  if stats.Rows[ i ].key == "CURRENT_MEMORY"  then ram.Used  = stats.Rows[ i ].value            end
  --  if stats.Rows[ i ].key == "MAX_MEMORY"      then ram.Total = stats.Rows[ i ].value            end
  --  if stats.Rows[ i ].key == "BYTES_OUT"       then 
  --  
  --    Response.node.cpu[ #Response.node.cpu + 1 ] = cpu
  --    Response.node.ram[ #Response.node.ram + 1 ] = ram
  --    cpu = { Used = 0, Total = 0 }
  --    ram = { Used = 0, Total = 0 }
--
  --  end
  --end

end

::done::

SetStatus( 200 )
Write( jsonEncode( Response ) )