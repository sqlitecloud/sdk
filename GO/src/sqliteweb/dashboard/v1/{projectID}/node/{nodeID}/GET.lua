--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/30
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Get a JSON with all 
--   ////                ///  ///                     providers, regions and 
--     ////     //////////   ///                      size parameters
--        ////            ////          Requires    : Authentication
--          ////     /////              Output      : Structure with node info
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/nodes

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local nodeID,    err, msg = checkNodeID( nodeID )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

Response = {
  status            = 0,                         -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  node              = {
    id              = 0,                         -- NodeID - It is not good to have a simple int number!!!!!!
    name            = "",                        -- Name of this node
    type            = "",                        -- Type fo this node, for example: Leader, Worker
    provider        = "",                        -- Provider of this node
    image           = "",                        -- Image data for this node
    region          = "",                        -- Regin data for this node
    size            = "",                        -- Size info for this node
    address         = "",                        -- IPv[4,6] address or host name of this node
    port            = 0,                         -- Port this node is listening on
 
    status          = "unknown",                 -- Replicating
    details         = "?/?/?",                   -- "SFO1/1GB/25GB disk
    raft            = { 0, 0 },                  -- array 8960, 8960
    load            = { 0, 0 },                  -- some load info
    cpu             = { Used = 0, Total = 0 },   -- some cpu info
    ram             = { Used = 0, Total = 0 },   -- some ram info
    coordinates     = { Lat  = 0, Lng   = 0 },   -- coordinates of the machine
  },
}

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end

  nodes = getINIArray( projectID, "nodes", "" )
  for i = 1, #nodes do
    url = parseConnectionString( nodes[ i ] )
    if url then 
      if url.Port == 0 then url.Port = 8860 end
      
      if nodeID == i - 1 then

        Response.node.id          = i
        Response.node.name        = string.format( "SQLiteCloud CORE Server node [%d]", i )
        Response.node.type        = getINIString( projectID, "type",      "Worker" )
        Response.node.provider    = getINIString( projectID, "provider",  "On Premise" )
        Response.node.image       = getINIString( projectID, "image",     "Unknown" )
        Response.node.region      = getINIString( projectID, "type",      "On Premise" )
        Response.node.size        = getINIString( projectID, "type",      "Unknown" )
        Response.node.address     = url.Host
        Response.node.port        = url.Port

        -- get special values from the node
        Response.node.status      = "unknown"
        Response.node.details     = "?/?/?"
        Response.node.raft        = { 0, 0 }
        Response.node.load        = { 0, 0 }
        Response.node.cpu         = { Used = 0, Total = 0 }
        Response.node.ram         = { Used = 0, Total = 0 }
        Response.node.coordinates = { Lat = 0,  Lng = 0   }
        
        goto done
      end
    end
  end
                                          do  return error( 404, "NodeID not found" )               end

else

  query = string.format( "SELECT NODE.id, NODE.name, type, provider, image, region, size, IIF( addr4, addr4, '' ) || IIF( addr4 AND addr6, ',', '' ) || IIF( addr6, addr6, '' ) AS address, port FROM USER JOIN PROJECT ON USER.id = PROJECT.user_id JOIN NODE ON PROJECT.uuid = NODE.project_uuid WHERE USER.enabled = 1 AND USER.id = %d AND NODE.id = %d AND uuid='%s';", userID, nodeID, enquoteSQL( projectID ) )
  nodes = executeSQL( "auth", query )

  if not nodes                            then return error( 404, "ProjectID OR NodeID not found" ) end
  if nodes.ErrorNumber              ~= 0  then return error( 502, "Bad Gateway" )                   end
  if nodes.NumberOfColumns          ~= 9  then return error( 502, "Bad Gateway" )                   end
  if nodes.NumberOfRows             ~= 1  then return error( 404, "ProjectID OR NodeID not found" ) end

  Response.node = nodes.Rows[ 1 ]
  
  Response.node.status      = "unknown"
  Response.node.details     = "?/?/?"
  Response.node.raft        = { 0, 0 }
  Response.node.load        = { 0, 0 }
  Response.node.cpu         = { Used = 0, Total = 0 }
  Response.node.ram         = { Used = 0, Total = 0 }
  Response.node.coordinates = { Lat = 0,  Lng = 0   }

end

::done::

SetStatus( 200 )
Write( jsonEncode( Response ) )