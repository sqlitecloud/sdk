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

local name,       err, msg = getBodyValue( "name", 1 )                    if err ~= 0 then return error( err, msg )         end -- Dev1 Server
local hardware,   err, msg = getBodyValue( "hardware", 1 )                if err ~= 0 then return error( err, msg )         end -- 1VCPU/1GB/25GB
local region,     err, msg = getBodyValue( "region", 1 )                  if err ~= 0 then return error( err, msg )         end -- NYC3/US
local servertype, err, msg = getBodyValue( "type", 1 )                    if err ~= 0 then return error( err, msg )         end -- worker
local counter,    err, msg = getBodyValue( "counter", 0 )                 if err ~= 0 then return error( err, msg )         end -- 1
if not counter then counter = 1 end

local provider = 'DigitalOcean'
local size = "small"

Response = {
  status            = 200,                       -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message
  value             = {},                        -- Array with new node objects {uuid: string, name: string, nodeID: int}
}

if userID == 0 then         
  if not getINIBoolean( projectID, "enabled", false )                    then return error( 401, "Project Disabled" )      end
                                                                         return error( 501, "Not Implemented" )
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg )         end

  local nodeID = 1
  result = executeSQL( "auth", "SELECT max_node_id FROM Project WHERE uuid = ?", {projectID} )
  if not result                                                          then return error( 504, "Gateway Timeout" )       end
  if result.ErrorNumber     ~= 0                                         then return error( 403, "Could not get project's info" )     end
  if result.NumberOfColumns  ~= 1                                        then return error( 502, "Bad Gateway" )           end
  if result.NumberOfRows     == 1 and result.Rows[ 1 ].max_node_id       then nodeID = result.Rows[ 1 ].max_node_id + 1    end

  result = executeSQL( "auth", "UPDATE Project SET max_node_id = max_node_id + ? WHERE uuid =  ?", {counter, projectID} )
  if not result                                                          then return error( 504, "Gateway Timeout" )       end
  if result.ErrorNumber     ~= 0                                         then return error( 403, "Could not update the project" )     end
  
  -- extract the numeric suffix of the input name
  local namenumberstr = stringMatch(name, "\\d+$")
  local basename = name
  if namenumberstr ~= nil then basename = name:gsub(namenumberstr.."$", "") end
  local namenumber = tonumber(namenumberstr)

  for i = 1, counter do
    if namenumber ~= nil then 
      -- the input name already ends with a number
      -- so add i-1 to that number and use with the same number of digits of the original number
      -- example "ubuntu-s-1vcpu-1gb-1" -> "ubuntu-s-1vcpu-1gb-1" .. "ubuntu-s-1vcpu-1gb-2"
      local format = string.format("%%s%%0%dd", string.len(namenumberstr) )
      name = string.format(format, basename, namenumber + i - 1) 
    else 
      -- the input name doesn't end with a number
      -- so add concatenate the nodeID
      -- example "ubuntu-s-1vcpu-1gb" -> "ubuntu-s-1vcpu-1gb-01" .. "ubuntu-s-1vcpu-1gb-02"
      name = string.format("%s-%02d", basename, nodeID) 
    end
  
    local jobuuid = createNode(userID, name, region, hardware, servertype, projectID, nodeID)
    if jobuuid ~= nil and jobuuid ~= '' then 
      Response.value[#Response.value + 1] = {uuid = jobuuid, name = name, nodeID = nodeID}
    elseif counter == 1 then
      return error( 500, "Internal Server Error")
    else 
      Response.value[#Response.value + 1] = {name = name, error = "Internal Server Error"}
    end

    nodeID = nodeID + 1
  end

  -- reloadNodes(projectID)
end

SetStatus( 200 )
Write( jsonEncode( Response ) )