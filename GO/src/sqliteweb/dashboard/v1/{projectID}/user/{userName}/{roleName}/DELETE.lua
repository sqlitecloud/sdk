--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : REVOKE ROLE % USER % 
--   ////                ///  ///                     [DATABASE %] [TABLE %] 
--     ////     //////////   ///                      
--        ////            ////          Requires    : Authentication
--          ////     /////              Output      : status + message
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                              end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                              end
local userName,  err, msg = checkParameter( userName, 3 )                if err ~= 0 then return error( err, string.format( msg, "userName" ) ) end
local roleName,  err, msg = checkParameter( roleName, 3 )                if err ~= 0 then return error( err, string.format( msg, "roleName" ) ) end

local database,  err, msg = getBodyValue( "database", 0 )                if err ~= 0 then return error( err, msg )                              end
local table,     err, msg = getBodyValue( "table", 0 )                   if err ~= 0 then return error( err, msg )                              end

                                          query = string.format( "REVOKE ROLE '%s' USER '%s'" , enquoteSQL( roleName ), enquoteSQL( userName ) )
if string.len( database )  > 0       then query = string.format( "%s DATABASE '%s'"          , query, enquoteSQL( database  ) ) end
if string.len( table )     > 0       then query = string.format( "%s TABLE '%s'"             , query, enquoteSQL( table     ) ) end
                                          query = string.format( "%s;"                       , query )
result = nil

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
else
  check_access = string.format( "SELECT COUNT( User.id ) AS granted FROM User JOIN Company ON User.company_id = Company.id JOIN Project ON Company.id = Project.company_id WHERE User.enabled = 1 AND Company.enabled = 1 AND User.id= %d AND uuid = '%s';", userID, enquoteSQL( projectID ) )
  check_access = executeSQL( "auth", check_access )

  if not check_access                                                                then return error( 504, "Gateway Timeout" )     end
  if check_access.ErrorNumber       ~= 0                                             then return error( 502, "Bad Gateway" )         end
  if check_access.NumberOfColumns   ~= 1                                             then return error( 502, "Bad Gateway" )         end 
  if check_access.NumberOfRows      ~= 1                                             then return error( 502, "Bad Gateway" )         end
  if check_access.Rows[ 1 ].granted ~= 1                                             then return error( 401, "Unauthorized" )        end
end                                          
                               
result = executeSQL( projectID, query )                                          
if not result                                                                        then return error( 404, "ProjectID not found" ) end
if result.ErrorNumber               ~= 0                                             then return error( 502, result.ErrorMessage )   end
if result.NumberOfColumns           ~= 0                                             then return error( 502, "Bad Gateway" )         end
if result.NumberOfRows              ~= 0                                             then return error( 502, "Bad Gateway" )         end
if result.Value                     ~= "OK"                                          then return error( 502, result.Value )          end

error( 200, "OK" )