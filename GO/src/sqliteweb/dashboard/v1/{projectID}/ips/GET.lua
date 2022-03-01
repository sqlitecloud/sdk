-- LIST ALLOWED IP [ROLE %] [USER %]
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/ips

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

userid = tonumber( userid )                                                                 -- Is string and comes from JWT. Contents is a number.

if projectID               == "auth"      then return error( 404, "Forbidden ProjectID" ) end -- fbf94289-64b0-4fc6-9c20-84083f82ee64
if string.len( projectID ) ~= 36          then return error( 400, "Invalid ProjectID" )   end 

if not query.role then role = "*" else role = query.role user = ""         end
if not query.user then user = "*" else role = ""         user = query.user end

if role ~= "" then role = string.format( "ROLE '%s'", enquoteSQL( role ) ) end
if user ~= "" then user = string.format( "USER '%s'", enquoteSQL( user ) ) end

query = string.format( "LIST ALLOWED IP %s %s ;", role, user )
ips = nil

if userid == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end

  ips = executeSQL( projectID, query )
else
  check_access = string.format( "SELECT COUNT( id ) AS granted FROM USER JOIN PROJECT ON USER.id = user_id WHERE USER.enabled = 1 AND User.id= %d AND uuid = '%s';", userid, enquoteSQL( projectID ) )
  check_access = executeSQL( "auth", check_access )

  if not check_access                     then return error( 504, "Gateway Timeout" )     end
  if check_access.ErrorNumber       ~= 0  then return error( 502, "Bad Gateway" )         end
  if check_access.NumberOfColumns   ~= 1  then return error( 502, "Bad Gateway" )         end 
  if check_access.NumberOfRows      ~= 1  then return error( 502, "Bad Gateway" )         end
  if check_access.Rows[ 1 ].granted ~= 1  then return error( 401, "Unauthorized" )        end

  ips = executeSQL( projectID, query )
end

if not ips                                then return error( 404, "ProjectID not found" ) end
if ips.ErrorNumber                  ~= 0  then return error( 502, "Bad Gateway" )         end
if ips.NumberOfColumns              ~= 3  then return error( 502, "Bad Gateway" )         end
if ips.NumberOfRows                 <  1  then return error( 200, "OK" )                  end

IP = {
  address = "127.0.0.1",                         -- IPv[4/6]
  name    = "name",                              -- Name
  type    = "type";                              -- Type
}

Response = {
  status            = 0,                         -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  ips               = ips.Rows,                  -- Array with allowded IP's for this role or user
}

SetStatus( 200 )
Write( jsonEncode( Response ) )