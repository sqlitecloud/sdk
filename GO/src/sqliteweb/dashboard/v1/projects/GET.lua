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

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                     end

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

if userID == 0 then                                           -- get list of projects in ini file
  projects = listINIProjects()

  if #projects > 0 then
    Response.value = {}
    
    for i = 1, #projects do  
      if getINIBoolean( projects[ i ], "enabled", false ) then
        Project = {}
        
        Project.id          = projects[ i ]
        Project.name        = getINIString( Project.id, "name", string.format( "SQLiteCloud CORE Server node [%d]", 1 + #Response.value ) ) -- Falsch!
        Project.description = getINIString( Project.id, "description", "unknown" )

        Response.value[ 1 + #Response.value ] = Project
      end
    end
  end

else

  data = executeSQL( "auth", string.format( "SELECT uuid AS id, Project.name, description FROM User JOIN Company ON User.company_id = Company.id JOIN Project ON Company.id = Project.company_id WHERE User.enabled = 1 AND Company.enabled = 1 AND User.id = %d;", userid ) )

  if not data                              then return error( 404, "User not found" )                end
  if data.ErrorNumber                ~= 0  then return error( 502, "Bad Gateway" )                   end
  if data.NumberOfColumns            ~= 3  then return error( 502, "Bad Gateway" )                   end

  if data.NumberOfRows                > 0  then Response.value = data.Rows                        end
end

SetStatus( 200 )
Write( jsonEncode( Response ) )