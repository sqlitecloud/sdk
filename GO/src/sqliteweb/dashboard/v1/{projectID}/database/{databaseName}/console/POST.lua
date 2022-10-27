--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/06/01
--    ///             ///   ///  ///    Author      : Andrea Donetti
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Execute SQLiteCloud commands
--   ////                ///  ///                    
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message + response
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/{databaseName}

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

Response = {
  status        = 200,                       -- status code: 0 = no error, error otherwise
  message       = "OK",                      -- "OK" or error message

  value         = {
    result      = 0,                         -- 0 -> value is an error string, 1 -> value is a string, 2 -> value is a rowset
    value       = nil,
    columns     = nil
  }
}
        
local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                                    end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                                    end
local dbName,    err, msg = checkParameter( databaseName, 1 )            if err ~= 0 then return error( err, string.format( msg, "databaseName" ) )   end
local command,   err, msg = getBodyValue( "command", 0 )                 if err ~= 0 then return error( err, msg )                                    end

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
else
  local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                                   end
end

local query = string.format( "SWITCH DATABASE ?; %s", command )

result = executeSQL( projectID, query, {dbName} )
if not result                             then return error( 504, "Gateway Timeout" )     end

-- print("executeSQL err:" .. result.ErrorMessage .. " Rows:" .. result.NumberOfRows )

if result.ErrorNumber ~= 0                then 
  Response.value.result = 0
  Response.value.value = string.format( "%s (%d:%d)", result.ErrorMessage, result.ErrorNumber, result.ExtendedErrorNumber ) 
elseif result.Value                       then 
  Response.value.result = 1
  Response.value.value = result.Value
elseif result.NumberOfColumns > 0         then 
  Response.value.result = 2
  Response.value.columns = result.Columns
  Response.value.value = result.Rows               
end 

SetStatus( 200 )
Write( jsonEncode( Response ) )