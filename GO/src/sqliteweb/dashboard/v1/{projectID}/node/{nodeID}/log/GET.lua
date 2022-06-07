--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/30
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Filter log
--   ////                ///  ///                     
--     ////     //////////   ///                      
--        ////            ////          Requires    : Authentication
--          ////     /////              Output      : 
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- LIST LOG [FROM <start_date>] [TO <end_date>] [LEVEL <log_level>] [TYPE <log_type>] [ID] [ORDER DESC] [LIMIT <count>] [CURSOR <cursorid>] [NODE <nodeid>]", PRIVILEGE_HOSTADMIN, command_list_log_date, true, false, false, BITMASK(COMMAND_FLAG_READ)},
-- COUNT LOG [FROM <start_date>] [TO <end_date>] [LEVEL <log_level>] [TYPE <log_type>] [ID] [ORDER DESC] [NODE <nodeid>]  

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local nodeID,    err, msg = checkNodeID( nodeID )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                          end
local machineNodeID, err, msg = verifyNodeID( userID, projectID, nodeID )if err ~= 0 then return error( err, msg )                          end

Response = {
  status            = 200,          -- status code: 0 = no error, error otherwise
  message           = "OK",         -- "OK" or error message
  value             = {
    count           = nil,          -- Number of logs for the current filters, only returned if the CURSOR arg is empty
    next_cursor     = nil,          -- Value to be used in the next request to get the next page
    logs            = nil,           -- Array of logs
  },
}

if not query.level    then query.level  = ""     end
if not query.type     then query.type   = ""     end
if not query.limit    then query.limit  = "100"  end
if not query.cursor   then query.cursor = ""     end
if not query.from     then query.from   = ""     end
if not query.to       then query.to     = ""     end

local slevel = ""
local stype  = ""
if string.len( query.level ) > 0 then
  local level,     err, msg = checkNumber( query.level, 0, 5 )           if err ~= 0 then return error( err, string.format( msg, "level" ) ) end
  slevel = string.format( "LEVEL %d", level )
end
if string.len( query.type ) > 0 then
  local type,      err, msg = checkNumber( query.type, 1, 8 )            if err ~= 0 then return error( err, string.format( msg, "type" ) )  end
  stype = string.format( "TYPE %d", type )
end

local order = "ORDER DESC"

local sfrom = ""
local sto = ""
if string.len( query.from ) > 0 then  
  local from,    err, msg = checkDateTime( query.from )                  if err ~= 0 then return error( err, string.format( msg, "from" ) ) end
  sfrom = string.format( "FROM '%s'", from)
end
if string.len( query.to ) > 0 then  
  local to,      err, msg = checkDateTime( query.to )                    if err ~= 0 then return error( err, string.format( msg, "to" ) )   end
  sto = string.format( "TO '%s'", to)
end

local slimit = ""
local scursor = ""
if string.len( query.limit ) > 0 then
  local limit,     err, msg = checkNumber( query.limit, 0, 1000 )          if err ~= 0 then return error( err, string.format( msg, "limit" ) ) end
  slimit = string.format( "LIMIT %d", limit )

  if string.len( query.cursor ) > 0 then
    local cursor,     err, msg = checkNumber( query.cursor, 0, 18446744073709551615 )      if err ~= 0 then return error( err, string.format( msg, "cursor" ) ) end
    scursor = string.format( "CURSOR %d", cursor )
  else 
    -- get the total COUNT
    sql = string.format( "COUNT LOG %s %s %s %s ID %s NODE %d;", sfrom, sto, slevel, stype, order, machineNodeID ) 
    countlog = executeSQL( projectID, sql )
    if not countlog                                           then return error( 504, "Gateway Timeout" )            end
    if countlog.ErrorNumber ~= 0                              then return error( 502, countlog.ErrorMessage )        end
    if countlog.NumberOfColumns ~= 2                          then return error( 502, "Bad Gateway" )                end

    Response.value.count = countlog.Rows[1].count
    if Response.value.count > 0 and countlog.Rows[1].next_cursor then 
      scursor = string.format( "CURSOR %d", countlog.Rows[1].next_cursor )
    else 
      SetStatus( 200 )
      Write( jsonEncode( Response ) )
      return
    end
  end
end

sql = string.format( "LIST LOG %s %s %s %s ID %s %s %s NODE %d;", sfrom, sto, slevel, stype, order, slimit, scursor, machineNodeID ) 

log = executeSQL( projectID, sql )
if not log                                                                           then return error( 504, "Gateway Timeout" )            end
if log.ErrorNumber ~= 0                                                              then return error( 502, log.ErrorMessage )             end

flog = nil
if log.NumberOfRows > 0 then
  flog = filter( log.Rows, { [ "datetime"    ] = "date", 
                             [ "log_type"    ] = "type", 
                             [ "log_level"   ] = "level",
                             [ "description" ] = "description",
                             [ "username"    ] = "username",
                             [ "database"    ] = "database",
                             [ "ip_address"  ] = "address",
                             [ "connection_id" ] = "connection_id",
                           } )

  -- get the next cursor 
  if string.len(slimit) > 0  then   Response.value.next_cursor = log.Rows[log.NumberOfRows].id - 1     end
  Response.value.logs = flog                      
else
  Response.value = nil
end

SetStatus( 200 )
Write( jsonEncode( Response ) )