--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/04/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : LIST TABLES
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Database Infos
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/{projectID}/{databaseName}/tables

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

Response = {
  status            = 200,                       -- status code: 200 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message
  value             = nil,                       -- List of table info
}

local userID,        err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                                   end
local projectID,     err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                                   end
local databaseName,  err, msg = checkParameter( databaseName, 1 )            if err ~= 0 then return error( err, string.format( msg, "databaseName" ) )  end

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
else
  local projectID,     err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg )                                    end  
end

tables = executeSQL( projectID, string.format( "SWITCH DATABASE '%s'; LIST TABLES;", enquoteSQL( databaseName ) ) )
if not tables                                 then return error( 404, "ProjectID not found" ) end
if tables.ErrorMessage                  ~= "" then return error( 502, tables.ErrorMessage )   end
if tables.ErrorNumber                   ~= 0  then return error( 502, "Bad Gateway" )         end
if tables.NumberOfColumns               ~= 6  then return error( 502, "Bad Gateway" )         end

if tables.NumberOfRows                  > 0   then 
  Response.value = filter( tables.Rows,  { [ "name"   ] = "name", 
                                           [ "schema" ] = "schema", 
                                           [ "type"   ] = "type", 
                                           [ "ncol"   ] = "columns", 
                                           [ "strict" ] = "strict",
                                           [ "wr"     ] = "wr",  
                                           } )
end

SetStatus( 200 )
Write( jsonEncode( Response ) )