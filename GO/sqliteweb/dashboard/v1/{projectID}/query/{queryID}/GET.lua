--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.1
--     //             ///   ///  ///    Date        : 2022/11/28
--    ///             ///   ///  ///    Author      : Andrea Donetti
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : QUERY SUGGEST ID <query_id> 
--   ////                ///  ///                     [PERCENTAGE <percentage>] [APPLY]  
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Structure with query expert analysis info
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/users

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

command = "QUERY SUGGEST ID ?"
commandargs = {queryID}

if query.percentage       then 
  command = command .. " PERCENTAGE ?"
  commandargs[#commandargs+1] = query.percentage
end
if query.apply        then 
  command = command .. " APPLY"
end

res = nil

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false )                                then return error( 401, "Project Disabled" )           end
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg ) end  
end

res = executeSQL( projectID, command, commandargs )
if not res                                then return error( 404, "ProjectID not found" ) end
if res.ErrorNumber                  ~= 0  then return error( 502, "Bad Gateway" )         end
if res.NumberOfColumns              ~= 4  then return error( 502, "Bad Gateway" )         end
if res.NumberOfRows                 <  1  then return error( 200, "OK" )                  end

hierarchyres = nil

if #res == 0 then 
  res = nil  
else 
  hierarchyres = {}

  for i = 1, #res do 
    local row = res[ i ]
    local statement = row.statement
    local report = hierarchyres[statement]

    if not usermap then 
      usermap = {}
      hierarchyusers[username] = usermap
    end  

    usermap[ "enabled" ] = fusers[ i ].enabled
    local roles = usermap[ "roles" ]

    if not roles then 
      roles = {}
      usermap[ "roles" ] = roles
    end

    u.user = nil
    u.enabled = nil

    roles[ #roles + 1 ] = u


  end
end


Response = {
  status            = 200,                        -- status code: 0 = no error, error otherwise
  message           = "OK",                       -- "OK" or error message
  value             = hierarchyusers,             -- Array with user info
}

SetStatus( 200 )
Write( jsonEncode( Response ) )

Response = {
  status            = 200,                        -- status code: 0 = no error, error otherwise
  message           = "OK",                       -- "OK" or error message
  value             = fres,                       -- Array with queries info
}

SetStatus( 200 )
Write( jsonEncode( Response ) )