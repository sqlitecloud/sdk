-- LIST PRIVILEGES / List all privileges for this project
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/privileges

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

userid = tonumber( userid )                                                                 -- Is string and comes from JWT. Contents is a number.

if projectID               == "auth"      then return error( 404, "Forbidden ProjectID" ) end -- fbf94289-64b0-4fc6-9c20-84083f82ee64
if string.len( projectID ) ~= 36          then return error( 400, "Invalid ProjectID" )   end 

query      = "LIST PRIVILEGES ;"
privileges = nil

if userid == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end

  privileges = executeSQL( projectID, query )
else
  check_access = string.format( "SELECT COUNT( id ) AS granted FROM USER JOIN PROJECT ON USER.id = user_id WHERE USER.enabled = 1 AND User.id= %d AND uuid = '%s';", userid, enquoteSQL( projectID ) )
  check_access = executeSQL( "auth", check_access )

  if not check_access                     then return error( 504, "Gateway Timeout" )     end
  if check_access.ErrorNumber       ~= 0  then return error( 502, "Bad Gateway" )         end
  if check_access.NumberOfColumns   ~= 1  then return error( 502, "Bad Gateway" )         end 
  if check_access.NumberOfRows      ~= 1  then return error( 502, "Bad Gateway" )         end
  if check_access.Rows[ 1 ].granted ~= 1  then return error( 401, "Unauthorized" )        end

  privileges = executeSQL( projectID, query )
end

if not privileges                          then return error( 404, "ProjectID not found" ) end
if privileges.ErrorNumber            ~= 0  then return error( 502, "Bad Gateway" )         end
if privileges.NumberOfColumns        ~= 1  then return error( 502, "Bad Gateway" )         end
if privileges.NumberOfRows           <  1  then return error( 200, "OK" )                  end

p = {}
for i = 1, privileges.NumberOfRows do 
  p[ #p + 1 ] = privileges.Rows[ i ].name 
end
if #p == 0 then p = nil end

Response = {
  status            = 0,                         -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  privileges        = p,                         -- Array of privileges 
}

SetStatus( 200 )
Write( jsonEncode( Response ) )