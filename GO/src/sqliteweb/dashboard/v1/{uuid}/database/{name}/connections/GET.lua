json = require "json"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

result = sqlc_query( "auth", "LIST CONNECTIONS" )

Response = {
  Status = 0,
  Message = "Connections List",
  Connections = result.Rows
}

Write( json.encode( Response ) )
SetStatus( 200 )

-- {
--   "ResponseID": 0,
--   "Message": "Connections List",
--   "Connections": [
--       {
--           "Id": 2,
--           "Address": "192.168.1.23",
--           "Username": "admin",
--           "Database": "db1",
--           "ConnectionDate": "January 1, 1970 00:00:00 UTC",
--           "LastActivity": "January 1, 1970 00:00:00 UTC"
--       },
--       {
--           "Id": 4,
--           "Address": "192.168.1.23",
--           "Username": "admin",
--           "Database": "db1",
--           "ConnectionDate": "January 1, 1970 00:00:00 UTC",
--           "LastActivity": "January 1, 1970 00:00:00 UTC"
--       },
--       {
--           "Id": 7,
--           "Address": "192.168.1.23",
--           "Username": "admin",
--           "Database": "db1",
--           "ConnectionDate": "January 1, 1970 00:00:00 UTC",
--           "LastActivity": "January 1, 1970 00:00:00 UTC"
--       }
--   ]
-- }
