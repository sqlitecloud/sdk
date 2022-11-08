--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.1
--     //             ///   ///  ///    Date        : 2022/10/14
--    ///             ///   ///  ///    Author      : Andrea Donetti
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : LIST APIKEYS  
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Structure with apikey settings
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/apikeys

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false )                                then return error( 401, "Project Disabled" )           end
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg ) end  
end

query = "LIST APIKEYS"
result = executeSQL( projectID, query )
if not result                                then return error( 404, "ProjectID not found" ) end
if result.ErrorNumber                  ~= 0  then return error( 502, "Bad Gateway" )         end
if result.NumberOfColumns              ~= 6  then return error( 502, "Bad Gateway" )         end
if result.NumberOfRows                 <  1  then return error( 200, "OK" )                  end

fresult = filter( result.Rows, {[ "username"        ] = "username", 
                                [ "key"             ] = "key", 
                                [ "name"            ] = "name",
                                [ "expiration_date" ] = "expiration_date",
                                [ "restriction"     ] = "restriction",
                             } )

hierarchyresult = nil

if #fresult == 0 then 
    fresult = nil  
else 
  hierarchyresult = {}

  for i = 1, #fresult do 
    local r = fresult[ i ]
    local username = r.username
    r.username = nil
    -- print("username: " .. username)
    local userapikeys = hierarchyresult[username]

    if not userapikeys then 
      userapikeys = {}
      hierarchyresult[username] = userapikeys
    end  

    userapikeys[ #userapikeys+1 ] = r
  end
end


Response = {
  status            = 200,                        -- status code: 0 = no error, error otherwise
  message           = "OK",                       -- "OK" or error message
  value             = hierarchyresult,            -- Array with apikeys for each user
}

SetStatus( 200 )
Write( jsonEncode( Response ) )