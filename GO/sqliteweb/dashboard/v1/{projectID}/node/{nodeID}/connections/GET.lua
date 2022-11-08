--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/30
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Filter log
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

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local nodeID,    err, msg = checkNodeID( nodeID )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                          end
local machineNodeID,    err, msg = verifyNodeID( userID, projectID, nodeID )    if err ~= 0 then return error( err, msg )                   end

sql         = "LIST CONNECTIONS NODE ?;"
connections = executeSQL( projectID, sql, {machineNodeID} )

if not connections                                                                   then return error( 504, "Gateway Timeout" )            end
if connections.ErrorNumber     ~= 0                                                  then return error( 502, result.ErrorMessage )          end
if connections.NumberOfColumns ~= 6                                                  then return error( 502, "Bad Gateway" )                end

fcon = nil
if connections.NumberOfRows > 0 then
  fcon = filter( connections.Rows, { [ "id"              ] = "id", 
                                     [ "address"         ] = "address", 
                                     [ "username"        ] = "username",
                                     [ "database"        ] = "database",
                                     [ "connection_date" ] = "connection_date",
                                     [ "last_activity"   ] = "last_activity",
                                   } )
end

Response = {
  status            = 200,           -- status code: 200 = no error, error otherwise
  message           = "OK",          -- "OK" or error message
  value             = fcon           -- Array with key value pairs
}

SetStatus( 200 )
Write( jsonEncode( Response ) )