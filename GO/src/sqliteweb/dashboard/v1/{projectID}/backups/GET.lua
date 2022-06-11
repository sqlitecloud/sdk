--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/30
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Get a list of databases
--   ////                ///  ///                     that have backup enabled
--     ////     //////////   ///                      
--        ////            ////          Requires    : Authentication
--          ////     /////              Output      : 
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                          end

backups = executeSQL( projectID, "LIST BACKUPS;" )

if not backups                                                                       then return error( 504, "Gateway Timeout" )            end
if backups.ErrorNumber     ~= 0                                                      then return error( 502, backups.ErrorMessage )         end
if backups.NumberOfColumns ~= 1                                                      then return error( 502, "Bad Gateway" )                end

dbs = {}
for i = 1, backups.NumberOfRows do
  dbs[ #dbs + 1 ] = backups.Rows[ i ].name
end

if #dbs == 0 then dbs = nil end

Response = {
  status            = 200,           -- status code: 200 = no error, error otherwise
  message           = "OK",          -- "OK" or error message
  value             = dbs            -- Array with key value pairs
}

SetStatus( 200 )
Write( jsonEncode( Response ) )