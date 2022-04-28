--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : LIST ROLES
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Structure with user roles
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/roles


require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                     end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                     end

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
else
  check_access = string.format( "SELECT COUNT( id ) AS granted FROM USER JOIN PROJECT ON USER.id = user_id WHERE USER.enabled = 1 AND User.id= %d AND uuid = '%s';", userID, enquoteSQL( projectID ) )
  check_access = executeSQL( "auth", check_access )

  if not check_access                     then return error( 504, "Gateway Timeout" )     end
  if check_access.ErrorNumber       ~= 0  then return error( 502, "Bad Gateway" )         end
  if check_access.NumberOfColumns   ~= 1  then return error( 502, "Bad Gateway" )         end 
  if check_access.NumberOfRows      ~= 1  then return error( 502, "Bad Gateway" )         end
  if check_access.Rows[ 1 ].granted ~= 1  then return error( 401, "Unauthorized" )        end
end

roles = executeSQL( projectID, "LIST ROLES;" )
if not roles                              then return error( 404, "ProjectID not found" ) end
if roles.ErrorNumber                ~= 0  then return error( 502, "Bad Gateway" )         end
if roles.NumberOfColumns            ~= 5  then return error( 502, "Bad Gateway" )         end
if roles.NumberOfRows               <  1  then return error( 200, "OK" )                  end

froles = filter( roles.Rows, { [ "rolename"     ] = "name", 
                               [ "privileges"   ] = "privileges", 
                               [ "databasename" ] = "database",
                               [ "tablename"    ] = "table",
                               [ "builtin"      ] = "builtin",
                             } )
if #froles == 0 then froles = nil end
            
Role = {
  name        = "READ",                          -- Role name
  privileges  = "READ",                          -- Role privileges
  database    = "*",                             -- Related database
  table       = "*",                             -- Related table
  builtin     = 1                                -- 1 = build in role, 0 = user defined role set
}

Response = {
  status            = 200,                       -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  roles             = froles,                    -- Array with roles
}

SetStatus( 200 )
Write( jsonEncode( Response ) )