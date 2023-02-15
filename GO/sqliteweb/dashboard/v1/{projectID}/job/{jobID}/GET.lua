--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2023/02/14
--    ///             ///   ///  ///    Author      : Andrea Donetti
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Get the status of a job 
--   ////                ///  ///                     
--     ////     //////////   ///                      
--        ////            ////          Requires    : Authentication
--          ////     /////              Output      : Structure with job info
--             ///                      Copyright   : 2023 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/{projectID}/job/{jobID}

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

Response = {
  status            = 200,                       -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  value             = {
    uuid            = jobID,
    name            = "",
    modified        = "",
    status          = "",
    steps           = 0,
    progress        = 0,
    error           = 0,
    node_id         = 0,
    node_name       = "",
    hostname        = "",
  },
}

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
  return error( 501, "Not Implemented" )

else
  
  local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                     end

  command = "SELECT Jobs.uuid, Jobs.name, Jobs.stamp AS modified, Jobs.status, Jobs.steps, Jobs.progress, Jobs.error, Node.id as node_id, Node.name AS node_name, Node.hostname AS hostname FROM Jobs JOIN Node ON Jobs.node_id = Node.id WHERE Jobs.uuid = ? AND Jobs.user_id = ? AND Node.project_uuid = ?;"
  jobs = executeSQL( "auth", command, {jobID, userID, projectID} )

  if not jobs                            then return error( 404, "Job not found" ) end
  if jobs.ErrorNumber              ~= 0  then return error( 502, "Bad Gateway" )   end
  if jobs.NumberOfColumns          ~= 10  then return error( 502, "Bad Gateway" )   end
  if jobs.NumberOfRows             ~= 1  then return error( 404, "Job not found" ) end

  Response.value = jobs.Rows[ 1 ]  
end

SetStatus( 200 )
Write( jsonEncode( Response ) )