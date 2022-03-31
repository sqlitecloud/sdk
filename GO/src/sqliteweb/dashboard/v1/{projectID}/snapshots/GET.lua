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
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                          end

sql     = "LIST BACKUPS;"
backups = executeSQL( projectID, sql )

if not backups                                                                   then return error( 504, "Gateway Timeout" )            end
if backups.ErrorNumber     ~= 0                                                  then return error( 502, result.ErrorMessage )          end
if backups.NumberOfColumns ~= 1                                                  then return error( 502, "Bad Gateway" )                end

snaps = {}
for i = 1, backups.NumberOfRows do
  sql = string.format( "LIST BACKUPS DATABASE '%s';", enquoteSQL( backups.Rows[ i ].name ) )

  snapshots = executeSQL( projectID, sql )
  
  ss = { x = 1 }
  for s = 1, snapshots.NumberOfRows do
    ss[ #ss + 1 ] = { 
      created = snapshots.Rows[ s ].created, 
      id = string.format( "%s:%d", snapshots.Rows[ s ].generation, snapshots.Rows[ s ].index ),
      offset = snapshots.Rows[ s ].offset, 
      replica = snapshots.Rows[ s ].replica, 
      size = snapshots.Rows[ s ].size, 
      type = snapshots.Rows[ s ].type, 
    }
  end


  snaps[ #snaps + 1 ] = { database = backups.Rows[ i ].name, 
  snapshots = ss,
  --snapshots = snapshots.Rows 
  }

end

Response = {
  status            = 200,           -- status code: 200 = no error, error otherwise
  message           = "OK",          -- "OK" or error message

  snapshots         = snaps           -- Array with key value pairs
}

SetStatus( 200 )
Write( jsonEncode( Response ) )
