-- Get a JSON with all providers, regions and size parameters

userid = 0

if tonumber( userid ) < 0 then return error( 416, "Invalid User" ) end

Project = {
  id               = "00000000-0000-0000-0000-000000000000",  -- UUID of the project
 
  name             = "",                                      -- Project name
  description      = ""                                       -- Project description
}

Response = {
  status           = 0,                                       -- status code: 0 = no error, error otherwise
  message          = "OK",                                    -- "OK" or error message

  projects         = nil                                      -- Array with project objects
}

if tonumber( userid ) == 0 then 
  -- get list of projects in ini file

print( "los")

  projects = listINIProjects()

  
  if #projects > 0 then
    Response.projects = {}
    
    for i = 1, #projects do  
      if getINIBoolean( projects[ i ], "enabled", false ) then
        Project = {}
        
        Project.id          = projects[ i ]
        Project.name        = getINIString( Project.id, "name", string.format( "SQLiteCloud CORE Server node [%d]", 1 + #Response.projects ) ) -- Falsch!
        Project.description = getINIString( Project.id, "description", "unknown" )

        Response.projects[ 1 + #Response.projects ] = Project
      end
    end
  end

else
  data = executeSQL( "auth", string.format( "SELECT uuid AS id, PROJECT.name, description FROM USER JOIN PROJECT ON USER.id = PROJECT.user_id WHERE USER.enabled = 1 AND user.id = %d;", userid ) )
  if data.ErrorNumber == 0 and data.NumberOfRows == 1 then
    Response.projects = data.Rows[ 1 ]
  end
end

SetStatus( 200 )
Write( jsonEncode( Response ) )