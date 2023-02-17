--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2023/02/14
--    ///             ///   ///  ///    Author      : Andrea Donetti
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Get the status node-related active jobs
--   ////                ///  ///                     
--     ////     //////////   ///                      
--        ////            ////          Requires    : Authentication
--          ////     /////              Output      : Structure with job info
--             ///                      Copyright   : 2023 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/{projectID}/jobs/nodes

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

Response = {
  status            = 200,                       -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  value             = {},
}

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
  return error( 501, "Not Implemented" )

else
  
  local projectID, err, msg = verifyProjectID( userID, projectID )         if err ~= 0 then return error( err, msg )                     end

  command = "SELECT Jobs.uuid, Jobs.name, Jobs.stamp AS modified, Jobs.status, Jobs.steps, Jobs.progress, Jobs.error, Node.id as node_id, Node.name AS node_name, Node.hostname AS hostname FROM Jobs JOIN Node ON Jobs.node_id = Node.id WHERE Node.project_uuid = ? AND (Jobs.archived = 0 OR (Jobs.error = 0 AND Jobs.Progress < Jobs.Steps));"
  jobs = executeSQL( "auth", command, {projectID} )

  if not jobs                            then return error( 404, "Job not found" ) end
  if jobs.ErrorNumber              ~= 0  then return error( 502, "Bad Gateway" )   end
  if jobs.NumberOfColumns          ~= 10  then return error( 502, "Bad Gateway" )   end

  Response.value = jobs.Rows
end

SetStatus( 200 )
Write( jsonEncode( Response ) )