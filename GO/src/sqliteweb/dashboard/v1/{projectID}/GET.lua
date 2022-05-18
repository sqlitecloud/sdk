--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Get a JSON with all providers, 
--   ////                ///  ///                     regions and size parameters
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Structure with project info
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                         if err ~= 0 then return error( err, msg )                     end
local projectID, err, msg = checkProjectID( projectID )                   if err ~= 0 then return error( err, msg )                     end

Project = {
  id               = "00000000-0000-0000-0000-000000000000",  -- UUID of the project
  name             = "",                                      -- Project name
  description      = ""                                       -- Project description
}

Response = {
  status           = 200,                                     -- status code: 200 = no error, error otherwise
  message          = "OK",                                    -- "OK" or error message
  value            = nil                                      -- Array with project objects
}

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false )                                 then return error( 401, "Project Disabled" )      end
  Project.id          = projectID
  Project.name        = getINIString( Project.id, "name", string.format( "SQLiteCloud CORE Server node [%d]", 1 ) ) -- Falsch!
  Project.description = getINIString( Project.id, "description", "unknown" )
else 
  local projectID, err, msg  = verifyProjectID( userID, projectID )      if err ~= 0  then return error( err, msg )                     end

  project = executeSQL( "auth", string.format( "SELECT uuid AS id, Project.name, description FROM User JOIN Company ON User.company_id = Company.id JOIN Project ON Company.id = Project.company_id WHERE User.enabled = 1 AND user.id = %d AND Project.uuid = '%s';", userID, enquoteSQL( projectID ) ) )

  if not project                                                                      then return error( 404, "User not found" )        end
  if project.ErrorMessage               ~= ""                                         then return error( 502, project.ErrorMessage )    end
  if project.ErrorNumber                ~= 0                                          then return error( 502, "Bad Gateway" )           end
  if project.NumberOfColumns            ~= 3                                          then return error( 502, "Bad Gateway" )           end
  if project.NumberOfRows               == 1                                          then Response.value = project.Rows[ 1 ]           end
end

SetStatus( 200 )
Write( jsonEncode( Response ) )