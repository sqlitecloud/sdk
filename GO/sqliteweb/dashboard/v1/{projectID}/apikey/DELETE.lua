--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/10/14
--    ///             ///   ///  ///    Author      : Andrea Donetti
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : DROP APIKEY %
--   ////                ///  ///                     
--     ////     //////////   ///                      
--        ////            ////          Requires    : Authentication
--          ////     /////              Output      : status + message
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/apikey

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                   end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                   end
local key,       err, msg = getBodyValue( "key", 3 )                     if err ~= 0 then return error( err, msg )                   end

local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg )                   end

command  = "DROP APIKEY ?"
result = executeSQL( projectID, command, {key} )           
            
if not result                                                            then return error( 404, "ProjectID not found" )             end
if result.ErrorNumber       ~= 0                                         then return error( 404, result.ErrorMessage )               end
if result.NumberOfColumns   ~= 0                                         then return error( 502, "Bad Gateway" )                     end
if result.NumberOfRows      ~= 0                                         then return error( 502, "Bad Gateway" )                     end
if result.Value             ~= "OK"                                      then return error( 502, result.Value )                      end

error( 200, "OK" )