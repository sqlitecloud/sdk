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
  status            = 200,                       -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  node              = {
    id              = nodeID,                    -- NodeID - It is not good to have a simple int number!!!!!! 
    stats           = {},
  },
}

local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                     end
local machineNodeID,    err, msg = verifyNodeID( userID, projectID, nodeID )    if err ~= 0 then return error( err, msg )                     end

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
    Response.node.stats[ #Response.node.stats + 1 ] = row
    row = nil
  end
end

if #Response.node.stats == 0 then Response.node.stats = nil end

SetStatus( 200 )
Write( jsonEncode( Response ) )