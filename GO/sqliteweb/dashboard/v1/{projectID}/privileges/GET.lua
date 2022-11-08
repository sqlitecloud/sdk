--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : LIST PRIVILEGES / List all 
--   ////                ///  ///                     privileges for this project
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : List with all privileges
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/privileges

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
else
  local projectID, err, msg = verifyProjectID( userID, projectID )                if err ~= 0 then return error( err, msg )                     end
end

privileges = executeSQL( projectID, "LIST PRIVILEGES;" )
if not privileges                          then return error( 404, "ProjectID not found" ) end
if privileges.ErrorNumber            ~= 0  then return error( 502, "Bad Gateway" )         end
if privileges.NumberOfColumns        ~= 1  then return error( 502, "Bad Gateway" )         end
if privileges.NumberOfRows           <  1  then return error( 200, "OK" )                  end

p = {}
for i = 1, privileges.NumberOfRows do 
  p[ #p + 1 ] = privileges.Rows[ i ].name 
end
if #p == 0 then p = nil end

Response = {
  status            = 200,                       -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message
  value        = p,                              -- Array of privileges 
}

SetStatus( 200 )
Write( jsonEncode( Response ) )