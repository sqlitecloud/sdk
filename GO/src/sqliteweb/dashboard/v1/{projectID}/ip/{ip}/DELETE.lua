-- REMOVE ALLOWED IP % [ROLE %] [USER %]
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/ip/{ip}

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

userid = tonumber( userid )                                                                     -- Is string and comes from JWT. Contents is a number.

if projectID                  == "auth"   then return error( 404, "Forbidden ProjectID" )   end -- fbf94289-64b0-4fc6-9c20-84083f82ee64
if string.len( projectID )    ~= 36       then return error( 400, "Invalid ProjectID" )     end 
if string.len( ip )           < 1         then return error( 400, "Invalid IP" )            end 
if string.len( body )         == 0        then return error( 400, "Missing body" )          end

body = jsonDecode( body ) 

if not body                               then return error( 400, "Invalid body" )          end
if not body.role                          then body.role = ""                               end
if not body.user                          then body.user = ""                               end

if string.len( body.role ) < 1 and string.len( body.user ) < 1 then return error( 400, "Missing role or user" ) end

                                               query = string.format( "REMOVE ALLOWED IP '%s'", enquoteSQL( ip ) )
if string.len( body.role )   > 0          then query = string.format( "%s ROLE '%s'"       , query, enquoteSQL( body.role ) ) end
if string.len( body.user )   > 0          then query = string.format( "%s USER '%s'"       , query, enquoteSQL( body.user ) ) end
                                               query = string.format( "%s ;"               , query )

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
if result.ErrorNumber       ~= 0          then return error( 404, result.ErrorMessage )   end
if result.NumberOfColumns   ~= 0          then return error( 502, "Bad Gateway" )         end
if result.NumberOfRows      ~= 0          then return error( 502, "Bad Gateway" )         end
if result.Value             ~= "OK"       then return error( 502, result.Value )          end

error( 200, "OK" )