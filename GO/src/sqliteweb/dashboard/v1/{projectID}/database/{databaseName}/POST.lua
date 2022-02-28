-- CREATE DATABASE % [KEY %] [ENCODING %] [IF NOT EXISTS]
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/{databaseName}

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

userid = tonumber( userid )                                                                   -- Is string and comes from JWT. Contents is a number.

if projectID                  == "auth" then return error( 404, "Forbidden ProjectID" )   end -- fbf94289-64b0-4fc6-9c20-84083f82ee64
if string.len( projectID )    ~= 36     then return error( 400, "Invalid ProjectID" )     end 
if string.len( databaseName ) < 1       then return error( 400, "Invalid DatabaseName" )  end
if string.len( body )         == 0      then return error( 400, "Missing body" )          end

body = jsonDecode( body )

if body     == nil                      then return error( 400, "Invalid body" )          end
if not body.key                         then body.key      = ""                           end
if not body.encoding                    then body.encoding = ""                           end

                                             query = string.format( "CREATE DATABASE '%s'", enquoteSQL( databaseName ) )
if string.len( body.key )      > 0      then query = string.format( "%s KEY '%s'",          query, enquoteSQL( body.key      ) ) end
if string.len( body.encoding ) > 0      then query = string.format( "%s ENCODING '%s'",     query, enquoteSQL( body.encoding ) ) end
                                             query = string.format( "%s IF NOT EXISTS;",    query )

result = nil

if userid == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end

  result = executeSQL( projectID, query )
else
  check_access = string.format( "SELECT COUNT( id ) AS granted FROM USER JOIN PROJECT ON USER.id = user_id WHERE USER.enabled = 1 AND USER.id= %d AND uuid = '%s';", userid, enquoteSQL( projectID ) )
  check_access = executeSQL( "auth", check_access )

  if not check_access                     then return error( 504, "Gateway Timeout" )     end
  if check_access.ErrorNumber       ~= 0  then return error( 502, "Bad Gateway" )         end
  if check_access.NumberOfColumns   ~= 1  then return error( 502, "Bad Gateway" )         end 
  if check_access.NumberOfRows      ~= 1  then return error( 502, "Bad Gateway" )         end
  if check_access.Rows[ 1 ].granted ~= 1  then return error( 401, "Unauthorized" )        end

  result = executeSQL( projectID, query )
end

if not result                             then return error( 404, "ProjectID not found" ) end
if result.ErrorNumber       ~= 0          then return error( 502, result.ErrorMessage )   end
if result.NumberOfColumns   ~= 0          then return error( 502, "Bad Gateway" )         end
if result.NumberOfRows      ~= 0          then return error( 502, "Bad Gateway" )         end

if result.Value             ~= "OK"       then return error( 404, result.Value )          end

error( 200, "OK" )