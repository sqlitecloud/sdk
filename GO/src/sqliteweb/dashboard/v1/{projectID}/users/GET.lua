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
if query.database       then query = query .. string.format( " DATABASE '%s'", query.database )     end
if query.table          then query = query .. string.format( " TABLE '%s'", query.table )           end

users = nil

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false )                                then return error( 401, "Project Disabled" )           end
else
  check_access = string.format( "SELECT COUNT( User.id ) AS granted FROM User JOIN Company ON User.company_id = Company.id JOIN Project ON Company.id = Project.company_id WHERE User.enabled = 1 AND Company.enabled = 1 AND User.id= %d AND uuid = '%s';", userID, enquoteSQL( projectID ) )
  check_access = executeSQL( "auth", check_access )

  if not check_access                       then return error( 504, "Gateway Timeout" )     end
  if check_access.ErrorNumber       ~= 0    then return error( 502, "Bad Gateway" )         end
  if check_access.NumberOfColumns   ~= 1    then return error( 502, "Bad Gateway" )         end 
  if check_access.NumberOfRows      ~= 1    then return error( 502, "Bad Gateway" )         end
  if check_access.Rows[ 1 ].granted ~= 1    then return error( 401, "Unauthorized" )        end
end

users = executeSQL( projectID, query )
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
  roles             = "ADMIN",                    -- Comma separated list of roles
  database          = "",                         -- Database
  table             = ""                          -- Table
}

Response = {
  status            = 200,                        -- status code: 0 = no error, error otherwise
  message           = "OK",                       -- "OK" or error message
  value             = fusers,                     -- Array with user info
}

SetStatus( 200 )
Write( jsonEncode( Response ) )