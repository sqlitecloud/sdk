--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : GET DATABASE [ID]
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Database ID
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/database/{databaseName}/id

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                                    end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                                    end
local databaseName,  err, msg = checkParameter( databaseName, 1 )        if err ~= 0 then return error( err, string.format( msg, "databaseName" ) )  end

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
else
  local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                                   end
end

id = executeSQL( projectID, "SWITCH DATABASE ?; GET DATABASE ID;", {databaseName} )
if not id                                 then return error( 404, "ProjectID not found" ) end
if id.ErrorNumber                   ~= 0  then return error( 502, "Bad Gateway" )         end
if id.NumberOfColumns               ~= 0  then return error( 502, "Bad Gateway" )         end
if id.NumberOfRows                  ~= 0  then return error( 200, "OK" )                  end

Response = {
  status            = 200,                       -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message
  value             = id.Value,                  -- The database ID
}

SetStatus( 200 )
Write( jsonEncode( Response ) )