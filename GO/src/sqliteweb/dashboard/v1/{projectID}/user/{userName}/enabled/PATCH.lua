--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : ENABLE USER % / DISABLE USER %
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/user/{userName}/enabled

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                              end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                              end
local userName,  err, msg = checkParameter( userName, 3 )                if err ~= 0 then return error( err, string.format( msg, "userName" ) ) end
local enabled,   err, msg = getBodyValue( "enabled", 0 )                 if err ~= 0 then return error( err, msg )                              end

if not enabled then enabled = true end

if bool( enabled ) then  
  query = string.format( "ENABLE USER '%s';", enquoteSQL( userName ) )
else
  query = string.format( "DISABLE USER '%s';", enquoteSQL( userName ) )
end
result  = nil

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" )                                              end
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg )                              end
end

result = executeSQL( projectID, query )
if not result                                                                        then return error( 404, "ProjectID not found" )            end
if result.ErrorNumber       ~= 0                                                     then return error( 404, result.ErrorMessage )              end
if result.NumberOfColumns   ~= 0                                                     then return error( 502, "Bad Gateway" )                    end
if result.NumberOfRows      ~= 0                                                     then return error( 502, "Bad Gateway" )                    end
if result.Value             ~= "OK"                                                  then return error( 502, result.Value )                     end

error( 200, "OK" )