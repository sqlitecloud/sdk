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

local name,         err, msg = getBodyValue( "name", 0 )                 if err ~= 0 then return error( err, msg )                                end
local description,  err, msg = getBodyValue( "description", 0 )          if err ~= 0 then return error( err, msg )                                end
local adminUsername,     err, msg = getBodyValue( "admin_username", 0 )  if err ~= 0 then return error( err, msg )                                end
local adminPassword,     err, msg = getBodyValue( "admin_password", 0 )  if err ~= 0 then return error( err, msg )                                end

if userID == 0 then 
  if not getINIBoolean( projectID, "enabled", false )                                then return error( 401, "Project Disabled" )                 end
                                                                                          return error( 501, "Not Implemented" )
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg )                                end

  command = string.format( "UPDATE Project SET ")
  commandargs = {}
  local separator = ""
  
  if name and string.len(name)>0 then   
    command = string.format( "%s%s name = ?", command, separator)
    commandargs[#commandargs+1] = name
    separator = ","
  end
  
  if description and string.len(description)>0 then   
    command = string.format( "%s%s description = ?", command, separator)
    commandargs[#commandargs+1] = description
    separator = ","
  end

  if adminUsername and string.len(adminUsername)>0 then   
    command = string.format( "%s%s admin_username = ?", command, separator)
    commandargs[#commandargs+1] = adminUsername
    separator = ","
  end

  if adminPassword and string.len(adminPassword)>0 then   
    command = string.format( "%s%s admin_password = ?", command, separator)
    commandargs[#commandargs+1] = adminPassword
    separator = ","
  end

  -- only if there are any changes
  if #commandargs > 0 then
    command = string.format( "%s WHERE uuid = ?; SELECT changes() AS success;", command )
    commandargs[#commandargs+1] = projectID
    result = executeSQL( "auth", command, commandargs )
    if not result                                                                      then return error( 504, "Gateway Timeout" )                  end
    if result.ErrorMessage      ~= ""                                                  then return error( 502, result.ErrorMessage )                end
    if result.ErrorNumber       ~= 0                                                   then return error( 502, "Bad Gateway" )                      end
    if result.NumberOfRows      ~= 1                                                   then return error( 502, "Bad Gateway" )                      end
    if result.NumberOfColumns   ~= 1                                                   then return error( 502, "Bad Gateway" )                      end
    if result.Rows[ 1 ].success ~= 1                                                   then return error( 500, "ProjectID not found" )              end
  end
end

error( 200, "OK" )