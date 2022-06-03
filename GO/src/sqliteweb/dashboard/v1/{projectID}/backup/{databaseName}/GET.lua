--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/30
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : List available backups for a database
--   ////                ///  ///                     
--     ////     //////////   ///                      
--        ////            ////          Requires    : Authentication
--          ////     /////              Output      : 
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,       err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                                   end
local projectID,    err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                                   end
local databaseName, err, msg = checkParameter( databaseName, 1 )            if err ~= 0 then return error( err, string.format( msg, "databaseName" ) )  end

local projectID,    err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                                   end

Response = {
  status            = 200,           -- status code: 200 = no error, error otherwise
  message           = "OK",          -- "OK" or error message
  value             = {}             -- Array with key value pairs
}

backups = executeSQL( projectID, string.format( "LIST BACKUPS DATABASE '%s';", enquoteSQL( databaseName ) ) )
if not backups                                                                        then return error( 504, "Gateway Timeout" )            end
if backups.ErrorNumber     ~= 0                                                       then return error( 502, result.ErrorMessage )          end
if backups.NumberOfColumns ~= 7                                                       then return error( 502, "Bad Gateway" )                end

if #backups.Rows > 0 then
  Response.value = filter( backups.Rows, {
    type    = "type",
    replica = "replica",
    size    = "size",
    created = "time_stamp",
  } )
else 
  Response.value = nil
end

SetStatus( 200 )
Write( jsonEncode( Response ) )
