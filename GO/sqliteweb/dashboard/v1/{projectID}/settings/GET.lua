--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/05/31
--    ///             ///   ///  ///    Author      : Andrea Donetti
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : List all project settings
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Structure with user settings
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/settings

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                     end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                     end

Setting = {
  key   = "",
  value = ""
}

Response = {
  status            = 200,                       -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message
  value             = {},                        -- Array with key value pairs
}

local projectID, err, msg = verifyProjectID( userID, projectID )                if err ~= 0 then return error( err, msg )                     end

command = "LIST KEYS DETAILED"
settings = executeSQL( projectID, command )

if not settings                          then return error( 404, "ProjectID OR NodeID not found" ) end
if settings.ErrorNumber            ~= 0  then return error( 502, "Bad Gateway" )                   end
if settings.NumberOfColumns        ~= 5  then return error( 502, "Bad Gateway" )                   end

if settings.NumberOfRows           ~= 0  then 
  Response.value = settings.Rows 
else 
  Response.value = nil                                                                        
end


SetStatus( 200 )
Write( jsonEncode( Response ) )