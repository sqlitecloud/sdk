--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/10/14
--    ///             ///   ///  ///    Author      : Andrea Donetti
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : CREATE APIKEY USER <username> NAME <key_name>  
--   ////                ///  ///                     [RESTRICTION <restriction_type>] [EXPIRATION <expiration_date>]
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2
 
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/apikey

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,      err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                                    end
local projectID,   err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                                    end
local username,    err, msg = getBodyValue( "username", 3 )                if err ~= 0 then return error( err, msg )                                    end
local name,        err, msg = getBodyValue( "name", 1 )                    if err ~= 0 then return error( err, msg )                                    end
local restriction, err, msg = getBodyValue( "restriction", 0 )             if err ~= 0 then return error( err, msg )                                    end
local expiration,  err, msg = getBodyValue( "expiration", 0 )              if err ~= 0 then return error( err, msg )                                    end

local projectID, err, msg = verifyProjectID( userID, projectID )           if err ~= 0 then return error( err, msg )                                    end

command = "CREATE APIKEY USER ? NAME ?"
commandargs = {username, name}
if string.len( restriction ) > 0   then 
    command = command .. " RESTRICTION ?"       
    commandargs[#commandargs+1] = restriction
end
if string.len( expiration )  > 0   then 
    command = command .. " EXPIRATION ?"  
    commandargs[#commandargs+1] = expiration     
end

result = executeSQL( projectID, command, commandargs)

if not result                                                                        then return error( 404, "ProjectID not found" )                  end
if result.ErrorNumber       ~= 0                                                     then return error( 404, result.ErrorMessage )                    end
if result.NumberOfColumns   ~= 0                                                     then return error( 502, "Bad Gateway" )                          end
if result.NumberOfRows      ~= 0                                                     then return error( 502, "Bad Gateway" )                          end

error( 200, "OK" )