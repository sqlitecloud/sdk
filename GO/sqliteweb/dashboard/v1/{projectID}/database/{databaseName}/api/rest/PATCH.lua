--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/12/21
--    ///             ///   ///  ///    Author      : Andrea Donetti
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Update tables REST API settings
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/{projectID}/{databaseName}/api/rest

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,        err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                                   end
local projectID,     err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                                   end
local databaseName,  err, msg = checkParameter( databaseName, 1 )            if err ~= 0 then return error( err, string.format( msg, "databaseName" ) )  end

local jbody = jsonDecode( body )
if not jbody or type(jbody) ~= "table" or #jbody == 0 then return error( 400, "Invalid body") end

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
else
  local projectID,     err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg )                                    end  
end

local sql = ""
local sqlargs = {}

local GET_BIT    = bit(1)
local POST_BIT   = bit(2)
local PATCH_BIT  = bit(3)
local DELETE_BIT = bit(4)

for i,v in ipairs(jbody) do 
  if v.tableName == nil then return error( 400, "Missing tableName field" ) end

  sql = sql .. "REPLACE INTO RestApiSettings(project_uuid, database_name, table_name, methods_mask) VALUES (?,?,?,?);"
  sqlargs[#sqlargs+1] = projectID
  sqlargs[#sqlargs+1] = databaseName
  sqlargs[#sqlargs+1] = v.tableName
  sqlargs[#sqlargs+1] = GET_BIT * bool_to_number(v.GET) + POST_BIT * bool_to_number(v.POST) + PATCH_BIT * bool_to_number(v.PATCH) + DELETE_BIT * bool_to_number(v.DELETE)
end

result = executeSQL( "auth", sql, sqlargs)
if not result                             then return error( 502, "Bad Gateway" )         end
if result.ErrorNumber       ~= 0          then return error( 502, result.ErrorMessage )   end
if result.NumberOfColumns   ~= 0          then return error( 502, "Bad Gateway" )         end
if result.NumberOfRows      ~= 0          then return error( 502, "Bad Gateway" )         end
if result.Value             ~= "OK"       then return error( 404, result.Value )          end

error( 200, "OK" )