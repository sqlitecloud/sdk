--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Modify the node info
--   ////                ///  ///                   
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/{nodeID}

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local nodeID,    err, msg = checkNodeID( nodeID )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

local name,      err, msg = getBodyValue( "name", 1 )                    if err ~= 0 then return error( err, msg )                          end -- Dev1 Server
local type,      err, msg = getBodyValue( "type", 1 )                    if err ~= 0 then return error( err, msg )                          end -- worker
local provider,  err, msg = getBodyValue( "provider", 1 )                if err ~= 0 then return error( err, msg )                          end -- DigitalOcean
local image,     err, msg = getBodyValue( "image", 1 )                   if err ~= 0 then return error( err, msg )                          end -- i386/1/1MB/100MB
local region,    err, msg = getBodyValue( "region", 1 )                  if err ~= 0 then return error( err, msg )                          end -- Rome/Italy
local size,      err, msg = getBodyValue( "size", 1 )                    if err ~= 0 then return error( err, msg )                          end -- small
local address,   err, msg = getBodyValue( "address", 1 )                 if err ~= 0 then return error( err, msg )                          end -- 64.227.11.116
local port,      err, msg = getBodyValue( "port", 1 )                    if err ~= 0 then return error( err, msg )                          end -- 9960

command = "UPDATE Node SET name=?, type=?, provider=?, image=?, region=?, size=?,"
commandargs = {name, type, provider, image, region, size}

if contains( address, ":" ) then
  command = command .. " addr6=?,"
  commandargs[#commandargs+1] = address
else
  command = command .. " addr4=?,"
  commandargs[#commandargs+1] = address
end

command = command .. " port=? WHERE project_uuid = ? AND id=?;"
commandargs[#commandargs+1] = port
commandargs[#commandargs+1] = projectID
commandargs[#commandargs+1] = nodeID

if userID == 0 then         
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Project Disabled" )           end
  return error( 501, "Not Implemented" )   
else      
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg )                          end
  local machineNodeID, err, msg = verifyNodeID( userID, projectID, nodeID )  if err ~= 0 then return error( err, msg )                      end                                                                                         

  result = executeSQL( "auth", command, commandargs )
  if not result                                                                      then return error( 504, "Gateway Timeout" )            end
  if result.ErrorNumber ~= 0                                                         then return error( 502, result.ErrorMessage )          end

end

error( 200, "OK" )