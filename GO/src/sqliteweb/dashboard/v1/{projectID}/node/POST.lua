--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Add a new node
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Structure with user settings
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )         end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )         end

local name,      err, msg = getBodyValue( "name", 1 )                    if err ~= 0 then return error( err, msg )         end -- Dev1 Server
local hardware,  err, msg = getBodyValue( "hardware", 1 )                if err ~= 0 then return error( err, msg )         end -- 1VCPU/1GB/25GB
local region,    err, msg = getBodyValue( "region", 1 )                  if err ~= 0 then return error( err, msg )         end -- NYC3/US
local type,      err, msg = getBodyValue( "type", 1 )                    if err ~= 0 then return error( err, msg )         end -- worker
local counter,   err, msg = getBodyValue( "counter", 0 )                 if err ~= 0 then return error( err, msg )         end -- 1
if not counter then counter = 1 end

local provider = 'DigitalOcean'
local size = "small"

if userID == 0 then         
  if not getINIBoolean( projectID, "enabled", false )                    then return error( 401, "Project Disabled" )      end
                                                                         return error( 501, "Not Implemented" )
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg )         end

  local nodeID = 0
  result = executeSQL( "auth", string.format( "SELECT MAX(node_id) AS max_node_id FROM Node WHERE project_uuid = '%s'", enquoteSQL(projectID) ))
  if not result                                                          then return error( 504, "Gateway Timeout" )       end
  if result.NumberOfColumns  ~= 1                                        then return error( 502, "Bad Gateway" )           end
  if result.NumberOfRows     == 1 and result.Rows[ 1 ].max_node_id       then nodeID = result.Rows[ 1 ].max_node_id + 1    end

  for i = 1, counter do
    -- TODO: Create the virtual machine
    -- TODO: Set following columns with data from the provider: addr4, port, latitude, longitude
    local addr4 = "64.227.11.116"
    local port = 9960
    local latitude = "40.8054"
    local longitude = "-74.0241"

    -- TODO: Remove the virtual machine if the INSERT query returned an error
    result = executeSQL( "auth", string.format( "INSERT INTO Node (project_uuid, node_id, name, type, provider, image, region, size, addr4, port, latitude, longitude) VALUES ('%s', %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, '%s', '%s')", enquoteSQL(projectID), nodeID, enquoteSQL(name), enquoteSQL(type), enquoteSQL(provider), enquoteSQL(hardware), enquoteSQL(region), enquoteSQL(size), enquoteSQL(addr4), port, enquoteSQL(latitude), enquoteSQL(longitude) ))
    if not result                                                        then return error( 504, "Gateway Timeout" )       end
    if result.ErrorNumber     ~= 0                                       then return error( 502, result.ErrorMessage )     end
    if result.ErrorNumber ~= 0                                           then return error( 502, result.ErrorMessage )     end
    if result.Value ~= "OK"                                              then return error( 502, "Bad Gateway" )           end
 

    nodeID = nodeID + 1
  end

  reloadNodes(projectID)
end

error( 200, "OK" )
