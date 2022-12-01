--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/04/05
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : CREATE USER % PASSWORD % 
--   ////                ///  ///                     [ROLE %] [DATABASE %] [TABLE %]
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2
 
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/user/{userName}

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                                    end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                                    end
local userName,  err, msg = checkParameter( userName, 3 )                if err ~= 0 then return error( err, string.format( msg, "userName" ) )       end
local password,  err, msg = getBodyValue( "password", 3 )                if err ~= 0 then return error( err, msg )                                    end

local rolename,  err, msg = getBodyValue( "rolename", 0 )                if err ~= 0 then return error( err, msg )                                    end
local database,  err, msg = getBodyValue( "database", 0 )                if err ~= 0 then return error( err, msg )                                    end
local table,     err, msg = getBodyValue( "table", 0 )                   if err ~= 0 then return error( err, msg )                                    end

local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                                    end

command = "CREATE USER ? PASSWORD ?"
commandargs = {userName, password}
if string.len( rolename )  > 0     then 
    command = command .. " ROLE ?"
    commandargs[#commandargs+1] = rolename  
end
if string.len( database )  > 0     then 
    command = command .. " DATABASE ?"
    commandargs[#commandargs+1] = database
end
if string.len( table )     > 0     then 
    command = command .. " TABLE ?"
    commandargs[#commandargs+1] = table                            
end

result = executeSQL( projectID, command, commandargs )

if not result                                                                        then return error( 404, "ProjectID not found" )                  end
if result.ErrorNumber       ~= 0                                                     then return error( 404, "Database not found" )                   end
if result.NumberOfColumns   ~= 0                                                     then return error( 502, "Bad Gateway" )                          end
if result.NumberOfRows      ~= 0                                                     then return error( 502, "Bad Gateway" )                          end
if result.Value             ~= "OK"                                                  then return error( 502, "Bad Gateway" )                          end

error( 200, "OK" )