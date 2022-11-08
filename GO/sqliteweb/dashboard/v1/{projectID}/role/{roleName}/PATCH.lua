--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : RENAME ROLE % TO % 
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/role/{roleName}

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                               end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                               end
local roleName,  err, msg = checkParameter( roleName, 3 )                if err ~= 0 then return error( err, string.format( msg, "roleName" ) )  end
local name,      err, msg = getBodyValue( "name", 0 )                    if err ~= 0 then return error( err, msg )                               end

-- get privilege, it's differet if it is null or empty string
-- local privilege, err, msg = getBodyValue( "privilege", 0 )               if err ~= 0 then return error( err, msg )                               end
if not body                                then return error(400, "Missing body")                                                                end
if string.len( body ) == 0                 then return error(400, "Empty body")                                                                  end 
local jbody = jsonDecode( body )
if not jbody                               then return error(400, "Invalid body")                                                                end 
local privilege = jbody[ "privilege" ]

local database,  err, msg = getBodyValue( "database", 0 )                if err ~= 0 then return error( err, msg )                               end
local table,     err, msg = getBodyValue( "table", 0 )                   if err ~= 0 then return error( err, msg )                               end

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg )                               end
end

-- update name, if needed
if string.len(name) > 0 then
  query  = "RENAME ROLE ? TO ?;"
  result = executeSQL( projectID, query, {roleName, name} )
  if not result                             then return error( 404, "ProjectID not found" ) end
  if result.ErrorNumber       ~= 0          then return error( 404, result.ErrorMessage )  end
  if result.NumberOfColumns   ~= 0          then return error( 502, "Bad Gateway" )         end
  if result.NumberOfRows      ~= 0          then return error( 502, "Bad Gateway" )         end
  if result.Value             ~= "OK"       then return error( 502, "Bad Gateway" )         end

  roleName = name
end

-- update privilege, if needed
if privilege then
  if string.len(privilege) == 0             then return error( 404, string.format( "Invalid privilege", value ))                  end
  
  query = "SET PRIVILEGE ? ROLE ?"
  queryargs = {privilege, roleName}
  if string.len( database )   > 0    then 
    query = query .. " DATABASE ?"
    queryargs[#queryargs+1] = database
  end
  if string.len( table )      > 0    then 
    query = query .. " TABLE ?"
    queryargs[#queryargs+1] = table
  end

  result = executeSQL( projectID, query, queryargs )
  if not result                             then return error( 404, "ProjectID not found" ) end
  if result.ErrorNumber       ~= 0          then return error( 404, result.ErrorMessage )  end
  if result.NumberOfColumns   ~= 0          then return error( 502, "Bad Gateway" )         end
  if result.NumberOfRows      ~= 0          then return error( 502, "Bad Gateway" )         end
  if result.Value             ~= "OK"       then return error( 502, "Bad Gateway" )         end
end

error( 200, "OK" )