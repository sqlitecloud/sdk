-- Modify setting with key to value
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/{nodeID}/setting/{key}

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

userid = tonumber( userid )                                                                         -- Is string and comes from JWT. Contents is a number.
nodeid = tonumber( nodeID )                                                                         -- Is string but MUST contains a number

if projectID                     == "auth"  then return error( 404, "Forbidden ProjectID" )     end -- fbf94289-64b0-4fc6-9c20-84083f82ee64
if string.len( projectID )       ~= 36      then return error( 400, "Invalid ProjectID" )       end 
if string.format( "%d", nodeid ) ~= nodeID  then return error( 400, "NodeID is not a number" )  end 
if string.len( key  )            <  1       then return error( 400, "Missing Key" )             end
if string.len( body )            == 0       then return error( 400, "Missing body" )            end

body = jsonDecode( body ) 

if not body                                 then return error( 400, "Invalid body" )            end
if not body.value                           then return error( 400, "Missing value in body" )   end

-- TODO: DONT USE INSERT, CHANGE TO UPDATE
query  = string.format( "INSERT OR REPLACE INTO NODE_SETTINGS ( node_id, key, value ) SELECT NODE.id, '%s', '%s' FROM USER JOIN PROJECT ON USER.id = PROJECT.user_id JOIN NODE ON NODE.project_uuid = PROJECT.uuid WHERE USER.enabled = 1 AND USER_id = %d AND NODE.id = %d;", enquoteSQL( key ), enquoteSQL( body.value ), userid, nodeid )
result = nil

if userid == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
else
  check_access = string.format( "SELECT COUNT( id ) AS granted FROM USER JOIN PROJECT ON USER.id = user_id WHERE USER.enabled = 1 AND USER.id= %d AND uuid = '%s';", userid, enquoteSQL( projectID ) )
  check_access = executeSQL( "auth", check_access )

  if not check_access                     then return error( 504, "Gateway Timeout" )     end
  if check_access.ErrorNumber       ~= 0  then return error( 502, "Bad Gateway" )         end
  if check_access.NumberOfColumns   ~= 1  then return error( 502, "Bad Gateway" )         end 
  if check_access.NumberOfRows      ~= 1  then return error( 502, "Bad Gateway" )         end
  if check_access.Rows[ 1 ].granted ~= 1  then return error( 401, "Unauthorized" )        end

  result = executeSQL( projectID, query )




  result = executeSQL( "auth", query )

  if not result                             then return error( 404, "ProjectID not found" ) end
  if result.ErrorNumber       ~= 0          then return error( 502, result.ErrorMessage )   end
  if result.NumberOfColumns   ~= 0          then return error( 502, "Bad Gateway" )         end
  if result.NumberOfRows      ~= 0          then return error( 502, "Bad Gateway" )         end
  if result.Value             ~= "OK"       then return error( 502, result.Value )          end
end

error( 200, "OK" )