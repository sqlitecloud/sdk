--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/04/11
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Restore this backup and 
--   ////                ///  ///                     make it actual (restore)
--     ////     //////////   ///                      
--        ////            ////          Requires    : Authentication
--          ////     /////              Output      : status + message
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- Restore this snapshot and make it actual (restore)
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/backup/{datebaseName}/{timeStamp}

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,       err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                                   end
local projectID,    err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                                   end
local databaseName, err, msg = checkParameter( databaseName, 1 )            if err ~= 0 then return error( err, string.format( msg, "databaseName" ) )  end
local timeStamp,    err, msg = checkParameter( timeStamp, 20 )              if err ~= 0 then return error( err, string.format( msg, "timeStamp" ) )     end

local projectID,    err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                                   end

query  = string.format( "RESTORE BACKUP DATABASE '%s' TIMESTAMP '%s';", enquoteSQL( databaseName ), enquoteSQL( timeStamp ) )
result = executeSQL( projectID, query )

if not result                                                                           then return error( 504, "Gateway Timeout" )                     end
if result.ErrorMessage ~= ""                                                            then return error( 502, result.ErrorMessage )                   end
if result.ErrorNumber  ~= 0                                                             then return error( 502, "Bad Gateway" )                         end
error( 200, "OK" )