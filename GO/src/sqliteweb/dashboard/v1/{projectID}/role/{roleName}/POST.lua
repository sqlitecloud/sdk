--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : CREATE ROLE % [PRIVILEGE %] 
--   ////                ///  ///                     [DATABASE %] [TABLE %] 
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
local privilege, err, msg = getBodyValue( "privilege", 1 )               if err ~= 0 then return error( err, msg )                               end
local database,  err, msg = getBodyValue( "database", 1 )                if err ~= 0 then return error( err, msg )                               end
local table,     err, msg = getBodyValue( "table", 1 )                   if err ~= 0 then return error( err, msg )                               end

                                        query = string.format( "CREATE ROLE '%s'", enquoteSQL( roleName ) )
if string.len( privilege )  > 0    then query = string.format( "%s PRIVILEGE '%s'",            query, enquoteSQL( privilege ) ) end
if string.len( database )   > 0    then query = string.format( "%s DATABASE '%s'",             query, enquoteSQL( database  ) ) end
if string.len( table )      > 0    then query = string.format( "%s TABLE '%s'",                query, enquoteSQL( table     ) ) end
                                        query = string.format( "%s ;",                         query )

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
else
  check_access = string.format( "SELECT COUNT( id ) AS granted FROM USER JOIN PROJECT ON USER.id = user_id WHERE USER.enabled = 1 AND USER.id= %d AND uuid = '%s';", userID, enquoteSQL( projectID ) )
  check_access = executeSQL( "auth", check_access )

  if not check_access                     then return error( 504, "Gateway Timeout" )     end
  if check_access.ErrorNumber       ~= 0  then return error( 502, "Bad Gateway" )         end
  if check_access.NumberOfColumns   ~= 1  then return error( 502, "Bad Gateway" )         end 
  if check_access.NumberOfRows      ~= 1  then return error( 502, "Bad Gateway" )         end
  if check_access.Rows[ 1 ].granted ~= 1  then return error( 401, "Unauthorized" )        end
end

result = executeSQL( projectID, query )
if not result                             then return error( 404, "ProjectID not found" ) end
if result.ErrorNumber       ~= 0          then return error( 404, "Database not found" )  end
if result.NumberOfColumns   ~= 0          then return error( 502, "Bad Gateway" )         end
if result.NumberOfRows      ~= 0          then return error( 502, "Bad Gateway" )         end
if result.Value             ~= "OK"       then return error( 502, "Bad Gateway" )         end

error( 200, "OK" )