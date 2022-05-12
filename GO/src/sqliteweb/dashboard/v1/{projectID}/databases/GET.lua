--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : LIST DATABASES
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Database Infos
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/databases

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )       if err ~= 0 then return error( err, msg ) end
local projectID, err, msg = checkProjectID( projectID ) if err ~= 0 then return error( err, msg ) end

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg ) end  
end

databases = executeSQL( projectID, "LIST DATABASES DETAILED;" )
if not databases                          then return error( 404, "ProjectID not found" ) end
if databases.ErrorNumber            ~= 0  then return error( 502, "Bad Gateway" )         end
if databases.NumberOfColumns        < 10  then return error( 502, "Bad Gateway" )         end
if databases.NumberOfRows           <  1  then return error( 200, "OK" )                  end

db = {}
for i = 1, databases.NumberOfRows do 
  database                = {}
  database.name           = databases.Rows[ i ].name
  database.size           = databases.Rows[ i ].size
  database.connections    = databases.Rows[ i ].connections
  database.encryption     = databases.Rows[ i ].encryption
  database.backup         = databases.Rows[ i ].backup
  database.fragmentation  = databases.Rows[ i ].fragmentation
  database.stats          = { databases.Rows[ i ].nread,   databases.Rows[ i ].nwrite   }
  database.bytes          = { databases.Rows[ i ].inbytes, databases.Rows[ i ].outbytes }
  db[ #db + 1 ]           = database
end
if #db == 0 then db = nil end

Database = {
  name              = "Db1",
  size              = 18000000000,
  connections       = 5,
  encryption        = nil,
  backup            = "Daily",
  stats             = { 521, 12 },
  bytes             = { 8700000, 712 },
  fragmentation     = { Used = 2400000, total = 712000 }
}

Response = {
  status            = 200,                       -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  databases         = db,                        -- Array with Database objects
}

SetStatus( 200 )
Write( jsonEncode( Response ) )