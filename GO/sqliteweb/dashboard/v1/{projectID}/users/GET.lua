--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.1
--     //             ///   ///  ///    Date        : 2022/04/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil, Andrea
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : LIST USERS [WITH ROLES] 
--   ////                ///  ///                     [DATABASE %] [TABLE %]
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Structure with user settings
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

query = "LIST USERS WITH ROLES"
queryargs = {}

if query.database       then 
  query = query .. " DATABASE ?"
  queryargs[#queryargs+1] = query.database
end
if query.table          then 
  query = query .. " TABLE ?"
  queryargs[#queryargs+1] = query.table
end

users = nil

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false )                                then return error( 401, "Project Disabled" )           end
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg ) end  
end

users = executeSQL( projectID, query, queryargs )
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

hierarchyusers = nil

if #fusers == 0 then 
  fusers = nil  
else 
  hierarchyusers = {}

  for i = 1, #fusers do 
    local u = fusers[ i ]
    local username = u.user
    local usermap = hierarchyusers[username]

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