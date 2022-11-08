--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Update project values
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/projects

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,       err, msg = checkUserID( userid )                     if err ~= 0 then return error( err, msg )                                end
local projectID,    err, msg = checkProjectID( projectID )               if err ~= 0 then return error( err, msg )                                end

local name,         err, msg = getBodyValue( "name", 1 )                 if err ~= 0 then return error( err, msg )                                end
local description,  err, msg = getBodyValue( "description", 0 )          if err ~= 0 then return error( err, msg )                                end
local adminUsername,     err, msg = getBodyValue( "admin_username", 4 )  if err ~= 0 then return error( err, msg )                                end
local adminPassword,     err, msg = getBodyValue( "admin_password", 4 )  if err ~= 0 then return error( err, msg )                                end

if userID == 0 then 
  if not getINIBoolean( projectID, "enabled", false )                                then return error( 401, "Project Disabled" )                 end
                                                                                          return error( 501, "Not Implemented" )
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg )                                end

  query = "UPDATE Project SET name = ?, description = ?, admin_username = ?, admin_password = ? WHERE uuid = ?; SELECT changes() AS success;"
  result = executeSQL( "auth", query, {name, description, adminUsername, adminPassword, projectID} )

  if not result                                                                      then return error( 504, "Gateway Timeout" )                  end
  if result.ErrorMessage      ~= ""                                                  then return error( 502, result.ErrorMessage )                end
  if result.ErrorNumber       ~= 0                                                   then return error( 502, "Bad Gateway" )                      end
  if result.NumberOfRows      ~= 1                                                   then return error( 502, "Bad Gateway" )                      end
  if result.NumberOfColumns   ~= 1                                                   then return error( 502, "Bad Gateway" )                      end
  if result.Rows[ 1 ].success ~= 1                                                   then return error( 500, "ProjectID not found" )              end
end

error( 200, "OK" )