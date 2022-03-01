-- LIST DATABASES
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/databases

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

function getNumberOfConnections( projectID, databaseName )
  query = string.format( "LIST DATABASE CONNECTIONS '%s'; ", enquoteSQL( databaseName) )
  --print( query )
  if not query                            then  return 0 end
  if query.ErrorNumber ~= 0               then  return 0 end
  if query.NumberOfColumns ~= 2           then  return 0 end
                                                return query.NumberOfRows
end

userid = tonumber( userid )                                                                 -- Is string and comes from JWT. Contents is a number.

if projectID               == "auth"      then return error( 404, "Forbidden ProjectID" ) end -- fbf94289-64b0-4fc6-9c20-84083f82ee64
if string.len( projectID ) ~= 36          then return error( 400, "Invalid ProjectID" )   end 

query     = "LIST DATABASES;"
databases = nil

if userid == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end

  databases = executeSQL( projectID, query )
else
  check_access = string.format( "SELECT COUNT( id ) AS granted FROM USER JOIN PROJECT ON USER.id = user_id WHERE USER.enabled = 1 AND User.id= %d AND uuid = '%s';", userid, enquoteSQL( projectID ) )
  check_access = executeSQL( "auth", check_access )

  if not check_access                     then return error( 504, "Gateway Timeout" )     end
  if check_access.ErrorNumber       ~= 0  then return error( 502, "Bad Gateway" )         end
  if check_access.NumberOfColumns   ~= 1  then return error( 502, "Bad Gateway" )         end 
  if check_access.NumberOfRows      ~= 1  then return error( 502, "Bad Gateway" )         end
  if check_access.Rows[ 1 ].granted ~= 1  then return error( 401, "Unauthorized" )        end

  databases = executeSQL( projectID, query )
end

if not databases                          then return error( 404, "ProjectID not found" ) end
if databases.ErrorNumber            ~= 0  then return error( 502, "Bad Gateway" )         end
if databases.NumberOfColumns        ~= 1  then return error( 502, "Bad Gateway" )         end
if databases.NumberOfRows           <  1  then return error( 200, "OK" )                  end

db = {}
for i = 1, databases.NumberOfRows do 
  database                = {}
  database.name           = databases.Rows[ i ].name
  database.size           = 0
  database.connections    = getNumberOfConnections( projectID, database.name )
  database.encryption     = ""
  database.backup         = "Daily"
  database.stats          = { 521, 12 }
  database.bytes          = { 8700000, 712 }
  database.fragmentation  = { Used = 2400000, total = 712000}
  db[ #db + 1 ]           = database
end
if #db == 0 then db = nil end

Database = {
  name              = "Db1",
  size              = 18000000000,
  connections       = 5,
  encryption        = nil,
  backup            = "Daily",
  stats             = { 521, 12 },
  bytes             = { 8700000, 712 },
  fragmentation     = { Used = 2400000, total = 712000 }
}

Response = {
  status            = 0,                         -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  databases         = db,                        -- Array with Database objects
}

SetStatus( 200 )
Write( jsonEncode( Response ) )