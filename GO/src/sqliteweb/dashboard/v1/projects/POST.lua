--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Create a (empty) project
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/projects

require "sqlitecloud"
uuid = require "uuid"

uuid.seed()

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,       err, msg = checkUserID( userid )                     if err ~= 0 then return error( err, msg )                                end

local name,         err, msg = getBodyValue( "name", 1 )                 if err ~= 0 then return error( err, msg )                                end
local description,  err, msg = getBodyValue( "description", 0 )          if err ~= 0 then return error( err, msg )                                end
local username,     err, msg = getBodyValue( "username", 4 )             if err ~= 0 then return error( err, msg )                                end
local password,     err, msg = getBodyValue( "password", 4 )             if err ~= 0 then return error( err, msg )                                end

Project = {
  id               = "00000000-0000-0000-0000-000000000000",  -- UUID of the project
 
  name             = "",                                      -- Project name
  description      = ""                                       -- Project description
}

Response = {
  status           = 200,                                     -- status code: 200 = no error, error otherwise
  message          = "OK",                                    -- "OK" or error message

  projects         = nil                                      -- Array with project objects
}

if userID == 0 then return error( 501, "Not Implemented" )
else
  for i = 1, 20 do
    Project.id          = uuid() -- create a random uuid (and check if it is not already taken)
    Project.name        = name
    Project.description = description

    query = string.format( "INSERT INTO PROJECT VALUES( '%s', %d, '%s', '%s', '%s', '%s' );", enquoteSQL( Project.id ), userID, enquoteSQL( name ), enquoteSQL( description ), enquoteSQL( username ), enquoteSQL( password ) )

    result = executeSQL( "auth", query )
    if not result              then goto continue end
    if result.ErrorNumber ~= 0 then goto continue end
    if result.Value ~= "OK"    then goto continue end
    
    Response.projects    = { Project }
    
    SetStatus( 200 )
    Write( jsonEncode( Response ) )

    do return end

  ::continue::
  end
end

error( 500, "Could not create Project" )