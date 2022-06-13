--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/30
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Get the backup settings 
--   ////                ///  ///                     for each database
--     ////     //////////   ///                      
--        ////            ////          Requires    : Authentication
--          ////     /////              Output      : 
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/backups/settings

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                          end

settings = executeSQL( projectID, "LIST BACKUP SETTINGS;" )

dbs = {}
if not settings                                                          then return error( 504, "Gateway Timeout" )                        end
if settings.ErrorNumber     ~= 0                                         then return error( 502, settings.ErrorMessage )                    end
if settings.NumberOfColumns ~= 4                                         then return error( 502, "Bad Gateway" )                            end

if #settings.Rows > 0 then
  dbs = filter( settings.Rows, {
    name              = "name",
    enabled           = "enabled",
    backup_retention  = "backup_retention",
  } )
else 
  dbs = nil
end

Response = {
  status            = 200,           -- status code: 200 = no error, error otherwise
  message           = "OK",          -- "OK" or error message
  value             = dbs            -- Array with key value pairs
}

SetStatus( 200 )
Write( jsonEncode( Response ) )