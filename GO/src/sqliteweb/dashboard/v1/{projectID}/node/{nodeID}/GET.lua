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

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/node/{nodeID}

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local nodeID,    err, msg = checkNodeID( nodeID )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

Response = {
  status            = 200,                       -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  value             = {
    id              = nodeID,                    -- NodeID - It is not good to have a simple int number!!!!!!
    name            = "",                        -- Name of this node
    provider        = "",                        -- Provider of this node
    details         = "?/?/?",                   -- "SFO1/1GB/25GB disk
    region          = "",                        -- Region data for this node
    size            = "",                        -- Size info for this node
    address         = "",                        -- IPv[4,6] address or host name of this node
    port            = 0,                         -- Port this node is listening on
    latitude        = 44.931,                    -- coordinates of the machine
    longitude       = 10.533,
  },
}

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end

-- TODO: Verify NodeID

  nodes = getINIArray( projectID, "nodes", "" )
  for i = 1, #nodes do
    url = parseConnectionString( nodes[ i ] )
    if url then 
      if url.Port == 0 then url.Port = 8860 end
      
      if nodeID == i - 1 then

        Response.value.id          = i
        Response.value.name        = string.format( "SQLiteCloud CORE Server node [%d]", i )
        Response.value.type        = getINIString( projectID, "type",      "Worker"        )
        Response.value.provider    = getINIString( projectID, "provider",  "On Premise"    )
        Response.value.details     = getINIString( projectID, "image",     "?/?/?"         )
        Response.value.region      = getINIString( projectID, "region",    "On Premise"    )
        Response.value.size        = getINIString( projectID, "type",      "Unknown"       )
        Response.value.address     = url.Host
        Response.value.port        = url.Port
        Response.value.latitude    = 44.931
        Response.value.longitude   = 10.533

        goto done
      end
    end
  end
                                          do  return error( 404, "NodeID not found" )               end

else
  
  local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                     end
  local machineNodeID,    err, msg = verifyNodeID( userID, projectID, nodeID )    if err ~= 0 then return error( err, msg )              end

  query = string.format( "SELECT Node.id, Node.node_id, Node.name, provider, image AS details, region, size, IIF( addr4, addr4, '' ) || IIF( addr4 AND addr6, ',', '' ) || IIF( addr6, addr6, '' ) AS address, port, latitude, longitude FROM User JOIN Company ON User.company_id = Company.id JOIN Project ON Company.id = Project.company_id JOIN Node ON Project.uuid = Node.project_uuid WHERE User.enabled = 1 AND User.id = %d AND Node.id = %d AND uuid='%s';", userID, nodeID, enquoteSQL( projectID ) )
  nodes = executeSQL( "auth", query )

  if not nodes                            then return error( 404, "ProjectID OR NodeID not found" ) end
  if nodes.ErrorNumber              ~= 0  then return error( 502, "Bad Gateway" )                   end
  if nodes.NumberOfColumns          ~= 11 then return error( 502, "Bad Gateway" )                   end
  if nodes.NumberOfRows             ~= 1  then return error( 404, "ProjectID OR NodeID not found" ) end

  Response.value             = nodes.Rows[ 1 ]  -- id, node_id, name, type, provider, image->details, region, size, address, port
end

::done::

SetStatus( 200 )
Write( jsonEncode( Response ) )