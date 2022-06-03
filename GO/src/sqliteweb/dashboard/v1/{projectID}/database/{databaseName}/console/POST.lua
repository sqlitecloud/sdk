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

  value         = nil
}
        
local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                                    end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                                    end
local dbName,    err, msg = checkParameter( databaseName, 1 )            if err ~= 0 then return error( err, string.format( msg, "databaseName" ) )   end
local command,   err, msg = getBodyValue( "command", 0 )                 if err ~= 0 then return error( err, msg )                                    end

local query = string.format( "SWITCH DATABASE '%s'; %s", enquoteSQL( dbName ), command )

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
else
  check_access = string.format( "SELECT COUNT( User.id ) AS granted FROM User JOIN Company ON User.company_id = Company.id JOIN Project ON Company.id = Project.company_id WHERE User.enabled = 1 AND Company.enabled = 1 AND User.id = %d AND uuid = '%s';", userID, enquoteSQL( projectID ) )
  check_access = executeSQL( "auth", check_access )

  if not check_access                     then return error( 504, "Gateway Timeout" )     end
  if check_access.ErrorNumber       ~= 0  then return error( 502, "Bad Gateway" )         end
  if check_access.NumberOfColumns   ~= 1  then return error( 502, "Bad Gateway" )         end 
  if check_access.NumberOfRows      ~= 1  then return error( 502, "Bad Gateway" )         end
  if check_access.Rows[ 1 ].granted ~= 1  then return error( 401, "Unauthorized" )        end
end

result = executeSQL( projectID, query )
if not result                             then return error( 504, "Gateway Timeout" )     end
if result.ErrorNumber ~= 0                then Response.value = string.format( "ERROR: %s (%d)", result.ErrorMessage, result.ErrorNumber )
elseif result.Value                       then Response.value = result.Value
elseif result.Rows                        then Response.value = result.Rows               end 

SetStatus( 200 )
Write( jsonEncode( Response ) )