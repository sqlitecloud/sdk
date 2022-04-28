--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : LIST DATABASE CONNECTIONS [ID] %
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Structure with connection
--          ////     /////                            details
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2
 
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/Dummy/connections

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                                   end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                                   end
local databaseName,  err, msg = checkParameter( databaseName, 1 )        if err ~= 0 then return error( err, string.format( msg, "databaseName" ) )  end

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
else
  check_access = string.format( "SELECT COUNT( id ) AS granted FROM USER JOIN PROJECT ON USER.id = user_id WHERE USER.enabled = 1 AND User.id= %d AND uuid = '%s';", userID, enquoteSQL( projectID ) )
  check_access = executeSQL( "auth", check_access )

  if not check_access                     then return error( 504, "Gateway Timeout" )     end
  if check_access.ErrorNumber       ~= 0  then return error( 502, "Bad Gateway" )         end
  if check_access.NumberOfColumns   ~= 1  then return error( 502, "Bad Gateway" )         end 
  if check_access.NumberOfRows      ~= 1  then return error( 502, "Bad Gateway" )         end
  if check_access.Rows[ 1 ].granted ~= 1  then return error( 401, "Unauthorized" )        end
end

connections = executeSQL( projectID, string.format( "LIST DATABASE CONNECTIONS '%s';", enquoteSQL( databaseName ) ) )
if not connections                        then return error( 404, "ProjectID not found" ) end
if connections.ErrorNumber          ~= 0  then return error( 502, "Bad Gateway" )         end
if connections.NumberOfColumns      ~= 2  then return error( 502, "Bad Gateway" )         end
if connections.NumberOfRows         <  1  then return error( 200, "OK" )                  end


all = executeSQL( projectID, "LIST CONNECTIONS;" )
if not all                                then return error( 404, "ProjectID not found" ) end
if all.ErrorNumber                  ~= 0  then return error( 502, "Bad Gateway" )         end
if all.NumberOfColumns              ~= 6  then return error( 502, "Bad Gateway" )         end

c = {}
for i = 1, connections.NumberOfRows do 
  connection                    = {}
  connection.id                 = connections.Rows[ i ].client_id
  for j = 1, all.NumberOfRows do
    if connection.id == all.Rows[ j ].id then
      connection.address        = all.Rows[ j ].address
      connection.username       = all.Rows[ j ].username
      connection.database       = all.Rows[ j ].database
      connection.connectionDate = all.Rows[ j ].connection_date
      connection.lastActivity   = all.Rows[ j ].last_activity
      break
    end
  end
  c[ #c + 1 ] = connection
end
if #c == 0 then c = nil end

Connection = {
  id              = 0,                          -- Internal connection id
  address         = "127.0.0.1",                -- Clients IPv[4/6]address
  username        = "admin",                    -- Login username
  database        = "Dummy",                    -- Database name in use
  connectionDate  = "1970-01-01 00:00:00",      -- Date of connection in SQL format
  lastActivity    = "1970-01-01 00:00:00"       -- Date of last Activity in SQL format
}

Response = {
  status            = 200,                       -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  connections       = c,                         -- Array with Connection objects
}

SetStatus( 200 )
Write( jsonEncode( Response ) )