-- LIST DATABASE CONNECTIONS [ID] %
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/Dummy/connections


SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

userid = tonumber( userid )                                                                     -- Is string and comes from JWT. Contents is a number.

if projectID               == "auth"      then return error( 404, "Forbidden ProjectID" )   end -- fbf94289-64b0-4fc6-9c20-84083f82ee64
if string.len( projectID ) ~= 36          then return error( 400, "Invalid ProjectID" )     end 
if string.len( databaseName )      == 0   then return error( 500, "Internal Server Error" ) end

query       = string.format( "LIST DATABASE CONNECTIONS '%s';", enquoteSQL( databaseName ) )
connections = nil

if userid == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end

  connections = executeSQL( projectID, query )
else
  check_access = string.format( "SELECT COUNT( id ) AS granted FROM USER JOIN PROJECT ON USER.id = user_id WHERE USER.enabled = 1 AND User.id= %d AND uuid = '%s';", userid, enquoteSQL( projectID ) )
  check_access = executeSQL( "auth", check_access )

  if not check_access                     then return error( 504, "Gateway Timeout" )     end
  if check_access.ErrorNumber       ~= 0  then return error( 502, "Bad Gateway" )         end
  if check_access.NumberOfColumns   ~= 1  then return error( 502, "Bad Gateway" )         end 
  if check_access.NumberOfRows      ~= 1  then return error( 502, "Bad Gateway" )         end
  if check_access.Rows[ 1 ].granted ~= 1  then return error( 401, "Unauthorized" )        end

  connections = executeSQL( projectID, query )
end

if not connections                        then return error( 404, "ProjectID not found" ) end
if connections.ErrorNumber          ~= 0  then return error( 502, "Bad Gateway" )         end
if connections.NumberOfColumns      ~= 2  then return error( 502, "Bad Gateway" )         end
if connections.NumberOfRows         <  1  then return error( 200, "OK" )                  end

all = executeSQL( projectID, "LIST CONNECTIONS;" )

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
  status            = 0,                         -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  connections       = c,                         -- Array with Connection objects
}

SetStatus( 200 )
Write( jsonEncode( Response ) )