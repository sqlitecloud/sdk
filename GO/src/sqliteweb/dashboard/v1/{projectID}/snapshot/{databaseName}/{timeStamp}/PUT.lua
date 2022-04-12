--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/04/11
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

-- Resore this snapshot and make it actual (restore)
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/snapshot/{datebaseName}/{snapshotID}

-- {snapshotID} = generation:index

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

-- Main way: RESTORE BACKUP DATABASE 'db1.sqlite' TIMESTAMP '2022-04-01T16:21:55Z';
-- alternative way: RESTORE BACKUP DATABASE 'db1.sqlite' GENERATION 'd97d39a36e816caa' INDEX 0;

local userID,       err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                                   end
local projectID,    err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                                   end
local databaseName, err, msg = checkParameter( databaseName, 1 )            if err ~= 0 then return error( err, string.format( msg, "databaseName" ) )  end
local snapshotID,   err, msg = checkParameter( snapshotID, 18 )             if err ~= 0 then return error( err, string.format( msg, "snapshotID" ) )    end

local projectID,    err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                                   end


{ "generation":"generationName", "index":"IndexNumber", "timestamp":"TimeStamp"}

Response = {
  status            = 200,           -- status code: 200 = no error, error otherwise
  message           = "OK",          -- "OK" or error message

  snapshots         = {}           -- Array with key value pairs
}

snapshots = executeSQL( projectID, string.format( "LIST BACKUPS DATABASE '%s';", enquoteSQL( databaseName ) ) )
if not snapshots                                                                        then return error( 504, "Gateway Timeout" )            end
if snapshots.ErrorNumber     ~= 0                                                       then return error( 502, result.ErrorMessage )          end
if snapshots.NumberOfColumns ~= 7                                                       then return error( 502, "Bad Gateway" )                end

Response.snapshots = snapshots.Rows

if #Response.snapshots == 0 then Response.snapshots = nil end

SetStatus( 200 )
Write( jsonEncode( Response ) )
