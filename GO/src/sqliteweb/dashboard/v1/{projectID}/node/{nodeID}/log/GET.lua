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

-- LIST LOG FROM % TO % [LEVEL %] [TYPE %] [ORDER DESC]    
-- LIST % ROWS FROM LOG [LEVEL %] [TYPE %] [ORDER DESC]

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local nodeID,    err, msg = checkNodeID( nodeID )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

if not query.level  then query.level = "4"    end
if not query.type   then query.type  = "4"    end
if not query.order  then query.order = "DESC" end

local level,     err, msg = checkNumber( query.level, 0, 5 )             if err ~= 0 then return error( err, string.format( msg, "level" ) ) end
local type,      err, msg = checkNumber( query.type, 1, 8 )              if err ~= 0 then return error( err, string.format( msg, "type" ) )  end

if query.order ~= "DESC" and query.order ~= "ASC"                                    then return error( 400, "Order must be ASC or DESC" )   end
local order = string.format( "ORDER %s", query.order )

if query.rows then
  local rows,    err, msg = checkNumber( query.rows, 1, 10000 )          if err ~= 0 then return error( err, string.format( msg, "level" ) ) end

  sql = string.format( "LIST %d ROWS FROM LOG LEVEL %d TYPE %d %s", rows, level, type, order )
else
  if not query.to     then query.to     = now    end
  if not query.from   then query.from   = now_1h end

  local from,    err, msg = checkDateTime( query.from )                  if err ~= 0 then return error( err, string.format( msg, "from" ) ) end
  local to,      err, msg = checkDateTime( query.to )                    if err ~= 0 then return error( err, string.format( msg, "to" ) )   end

  sql = string.format( "LIST LOG FROM '%s' TO '%s' LEVEL %d TYPE %d %s;", from, to, level, type, order ) 
end

local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                          end

log = executeSQL( projectID, sql )
if not log                                                                           then return error( 504, "Gateway Timeout" )            end
if log.ErrorNumber ~= 0                                                              then return error( 502, result.ErrorMessage )          end

flog = nil
if log.NumberOfRows > 0 then
  flog = filter( log.Rows, { [ "datetime"    ] = "date", 
                             [ "log_type"    ] = "type", 
                             [ "log_level"   ] = "level",
                             [ "description" ] = "description",
                             [ "username"    ] = "username",
                             [ "database"    ] = "database",
                             [ "ip_address"  ] = "address",
                           } )
end

Response = {
  status            = 200,                       -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  logs              = flog                       -- Array with key value pairs
}

SetStatus( 200 )
Write( jsonEncode( Response ) )