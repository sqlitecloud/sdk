--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/05/31
--    ///             ///   ///  ///    Author      : Andrea Donetti
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Enable/disable a plugin
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/setting/{key}

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,     err, msg = checkUserID( userid )                       if err ~= 0 then return error( err, msg )                                  end
local projectID,  err, msg = checkProjectID( projectID )                 if err ~= 0 then return error( err, msg )                                  end
local pluginName, err, msg = checkParameter( pluginName, 1 )             if err ~= 0 then return error( err, string.format( msg, "pluginName" ) )   end

local enabled,    err, msg = getBodyValue( "enabled", 0 )                 if err ~= 0 then return error( err, msg )                          end -- Dev1 Server

local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                          end

if not enabled or enabled == "" or enabled == 0 or enabled == "0" or enabled == "false" then  
    query  = "DISABLE PLUGIN ?"
    queryargs = {pluginName} 
else  
    query  = "ENABLE PLUGIN ?"
    queryargs = {pluginName} 
end 
   
result = executeSQL( projectID, query, queryargs )
if not result                                                                      then return error( 504, "Gateway Timeout" )            end
if result.ErrorNumber ~= 0                                                         then return error( 502, result.ErrorMessage )          end

error( 200, "OK" )