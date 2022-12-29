--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/12/21
--    ///             ///   ///  ///    Author      : Andrea Donetti
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : LIST TABLES with REST API settings
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Tables REST API settings
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/{projectID}/{databaseName}/api/rest

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

tables = executeSQL( projectID, "SWITCH DATABASE ?; LIST TABLES;", {databaseName} )
if not tables                                 then return error( 404, "ProjectID not found" ) end
if tables.ErrorMessage                  ~= "" then return error( 500, tables.ErrorMessage )   end
if tables.ErrorNumber                   ~= 0  then return error( 502, "Bad Gateway" )         end
if tables.NumberOfColumns               ~= 6  then return error( 502, "Bad Gateway" )         end

settings = {}
settingsmap = {}
for i = 1, tables.NumberOfRows do 
  tsettings                     = {}
  tsettings.tableName           = tables.Rows[ i ].name  
  tsettings.GET                 = false
  tsettings.POST                = false
  tsettings.PATCH               = false
  tsettings.DELETE              = false
  settings[ #settings + 1 ]     = tsettings
  settingsmap[ tsettings.tableName ] = tsettings
end

local GET_BIT    = bit(1)
local POST_BIT   = bit(2)
local PATCH_BIT  = bit(3)
local DELETE_BIT = bit(4)

if #settings == 0 then 
    settings = nil 
else
    restsettings = executeSQL( "auth", "SELECT * FROM RestApiSettings WHERE database_name = ?", {databaseName} )
    if not restsettings                                 then return error( 404, "Settings not found" )  end
    if restsettings.ErrorMessage                  ~= "" then return error( 500, tables.ErrorMessage )   end
    if restsettings.ErrorNumber                   ~= 0  then return error( 502, "Bad Gateway" )         end
    if restsettings.NumberOfColumns               ~= 5  then return error( 502, "Bad Gateway" )         end

    for i = 1, restsettings.NumberOfRows do 
        row = restsettings.Rows[ i ]
        tsettings = settingsmap[row.table_name]
        if tsettings then
            tsettings.GET = hasbit(row.methods_mask, GET_BIT)
            tsettings.POST = hasbit(row.methods_mask, POST_BIT)
            tsettings.PATCH = hasbit(row.methods_mask, PATCH_BIT)
            tsettings.DELETE = hasbit(row.methods_mask, DELETE_BIT)
        end
    end
end

Response.value = settings

SetStatus( 200 )
Write( jsonEncode( Response ) )