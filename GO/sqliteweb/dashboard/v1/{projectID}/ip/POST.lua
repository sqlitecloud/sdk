--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : ADD ALLOWED IP % [ROLE %] 
--   ////                ///  ///                     [USER %]  
--     ////     //////////   ///                      
--        ////            ////          Requires    : Authentication
--          ////     /////              Output      : status + message
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/{projectID}/ip/{ip}

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                              end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                              end
local ip,        err, msg = getBodyValue( "ip", 3 )                      if err ~= 0 then return error( err, string.format( msg, "ip" ) )       end
local role,      err, msg = getBodyValue( "role", 0 )                    if err ~= 0 then return error( err, msg )                              end
local user,      err, msg = getBodyValue( "user", 0 )                    if err ~= 0 then return error( err, msg )                              end

if string.len( role ) < 1 and string.len( user ) < 1                                 then return error( 400, "Missing role or user" )           end

query = "ADD ALLOWED IP ?"
queryargs = {ip}
if string.len(role) > 0 then 
  query = query .. " ROLE ?"
  queryargs[#queryargs+1] = role
end
if string.len(user) > 0 then 
  query = query .. " USER ?"
  queryargs[#queryargs+1] = user
end

result    = nil

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg ) end  
end

result = executeSQL( projectID, query, queryargs )
if not result                             then return error( 404, "ProjectID not found" ) end
if result.ErrorNumber       ~= 0          then return error( 404, result.ErrorMessage )   end
if result.NumberOfColumns   ~= 0          then return error( 502, "Bad Gateway" )         end
if result.NumberOfRows      ~= 0          then return error( 502, "Bad Gateway" )         end
if result.Value             ~= "OK"       then return error( 502, result.Value )          end

error( 200, "OK" )