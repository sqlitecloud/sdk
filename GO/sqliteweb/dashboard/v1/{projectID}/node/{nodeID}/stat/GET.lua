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

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/nodes/{nodeID}/stat

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                     end
local nodeID,    err, msg = checkNodeID( nodeID )                        if err ~= 0 then return error( err, msg )                     end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                     end

Response = {
  status            = 200,          -- status code: 0 = no error, error otherwise
  message           = "OK",         -- "OK" or error message
  value             = {
    id              = nodeID,       -- Unique node ID 
    type            = "Leader",     -- Type fo this node, for example: Leader, Worker
    status          = "Replicate",  -- progress status of the node, for example: "Replicate", "Probe", "Snapshot" (cluster) or "Running" (nocluster)
    raft            = { 0, 0 },     -- array, index of the last raft entry matched by the node and by the leader, respectively
    load            = nil,          -- Load of the machine: num_clients, server_load, disk_usage_perc, example: [12,0.5,36.52]
    stats           = {},
    memory          = 0,            -- physical memory of the node
  },
}

local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                     end
local machineNodeID, err, msg = verifyNodeID( userID, projectID, nodeID )    if err ~= 0 then return error( err, msg )                 end

status = executeSQL( projectID, "LIST NODES;" )
if not status                                 then return error( 504, "Gateway Timeout" )       end
if status.ErrorNumber        ~= 0             then return error( 502, "Bad Gateway" )           end
if status.NumberOfRows       == 0             then return error( 404, "Empty node list" )       end

for i = 1, status.NumberOfRows do
  if status.Rows[ i ].status == "Leader" then Response.value.raft[ 2 ] = status.Rows[ i ].match end
  if i == machineNodeID then
    Response.value.status    = status.Rows[ i ].progress
    Response.value.raft[ 1 ] = status.Rows[ i ].match
    Response.value.type      = status.Rows[ i ].status
  end
end

command = "GET INFO LOAD,NUM_CLIENTS,DISK_USAGE_PERC NODE ?;" -- server_load, num_clients, cpu_time, mem_current, mem_max
load = executeSQL( projectID, command, {machineNodeID} )
-- print("command:", command)
if not load                                   then return error( 504, "Gateway Timeout" )       end
if load.ErrorNumber        ~= 0               then return error( 502, "Bad Gateway" )           end
if load.NumberOfRows       ~= 3               then return error( 502, "Bad Gateway" )           end

Response.value.load = {
  load.Rows[ 2 ].ARRAY, -- num_clients
  load.Rows[ 1 ].ARRAY, -- server_load
  load.Rows[ 3 ].ARRAY  -- disk_usage_perc
}

command = "LIST STATS NODE ? MEMORY;"
stats = executeSQL( projectID, command, {machineNodeID} )

if not stats                                  then return error( 504, "Gateway Timeout" )       end
if stats.ErrorNumber        ~= 0              then return error( 502, "Bad Gateway" )           end
if stats.NumberOfColumns    ~= 3              then return error( 502, "Bad Gateway" )           end
if stats.NumberOfRows       == 0              then return error( 404, "Stats not found" )       end

count = 0
for i = 1, stats.NumberOfRows do
  if stats.Rows[ i ].key == "PHYSICAL_MEMORY" then
    Response.value.memory = stats.Rows[ i ].value 
  else
    if not row then 
      row = { memory = { current = 0, max = 0 }, cpu = 0, clients = { current = 0, max = 0 }, commands = 0, io = { reads = 0, writes = 0 }, bytes = { reads = 0, writes = 0 }, sampletime = stats.Rows[ i ].datetime } 
    end

    if stats.Rows[ i ].key == "CPU_LOAD"        then row.cpu              = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "CURRENT_MEMORY"  then row.memory.current   = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "MAX_MEMORY"      then row.memory.max       = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "CURRENT_CLIENTS" then row.clients.current  = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "MAX_CLIENTS"     then row.clients.max      = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "NUM_COMMANDS"    then row.commands         = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "NUM_READS"       then row.io.reads         = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "NUM_WRITES"      then row.io.writes        = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "BYTES_IN"        then row.bytes.reads      = stats.Rows[ i ].value end
    if stats.Rows[ i ].key == "BYTES_OUT"       then row.bytes.writes     = stats.Rows[ i ].value end

    count = count + 1
    if count == 10 then 
      Response.value.stats[ #Response.value.stats + 1 ] = row 
      count = 0
      row = nil
    end
  end
end

if #Response.value.stats == 0 then Response.value.stats = nil end

SetStatus( 200 )
Write( jsonEncode( Response ) )