-- LIST USERS [WITH ROLES] [DATABASE %] [TABLE %]
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/users

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

userid = tonumber( userid )                                                                 -- Is string and comes from JWT. Contents is a number.

if projectID               == "auth"      then return error( 404, "Forbidden ProjectID" ) end -- fbf94289-64b0-4fc6-9c20-84083f82ee64
if string.len( projectID ) ~= 36          then return error( 400, "Invalid ProjectID" )   end 

if not query.database then database = "*" else database = query.database end
if not query.table    then table    = "*" else table    = query.table    end

query = string.format( "LIST USERS WITH ROLES DATABASE '%s' TABLE '%s';", enquoteSQL( database ), enquoteSQL( table ) )
users = nil

if userid == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end

  users = executeSQL( projectID, query )
else
  check_access = string.format( "SELECT COUNT( id ) AS granted FROM USER JOIN PROJECT ON USER.id = user_id WHERE USER.enabled = 1 AND User.id= %d AND uuid = '%s';", userid, enquoteSQL( projectID ) )
  check_access = executeSQL( "auth", check_access )

  if not check_access                     then return error( 504, "Gateway Timeout" )     end
  if check_access.ErrorNumber       ~= 0  then return error( 502, "Bad Gateway" )         end
  if check_access.NumberOfColumns   ~= 1  then return error( 502, "Bad Gateway" )         end 
  if check_access.NumberOfRows      ~= 1  then return error( 502, "Bad Gateway" )         end
  if check_access.Rows[ 1 ].granted ~= 1  then return error( 401, "Unauthorized" )        end

  users = executeSQL( projectID, query )
end

if not users                                then return error( 404, "ProjectID not found" ) end
if users.ErrorNumber                  ~= 0  then return error( 502, "Bad Gateway" )         end
if users.NumberOfColumns              ~= 5  then return error( 502, "Bad Gateway" )         end
if users.NumberOfRows                 <  1  then return error( 200, "OK" )                  end

fusers = filter( users.Rows, { [ "username"     ] = "user", 
                               [ "enabled"      ] = "enabled", 
                               [ "roles"        ] = "roles",
                               [ "databasename" ] = "database",
                               [ "tablename"    ] = "table",
                             } )
if #fusers == 0 then fusers = nil end

User = {
  user              = "admin",                    -- Username
  enabled           = 1,                          -- 1 = enabled, 0 = disabled
  roles             = "ADMIN",                    -- Comma seperated list of roles
  database          = "*",                        -- Database
  table             = "*"                         -- Table
}

Response = {
  status            = 0,                         -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  users             = fusers,                    -- Array with user info
}

SetStatus( 200 )
Write( jsonEncode( Response ) )