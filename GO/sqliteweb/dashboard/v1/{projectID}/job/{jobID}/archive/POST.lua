--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Archive a job. Archived jobs are no more returned by 
--   ////                ///  ///                     the dashboard/v1/{projectID}/jobs/nodes endpoint
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Structure with user settings
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/{projectID}/job/{jobID}/archive

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )         end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )         end

if userID == 0 then         
  if not getINIBoolean( projectID, "enabled", false )                    then return error( 401, "Project Disabled" )      end
                                                                         return error( 501, "Not Implemented" )
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg )         end

  result = executeSQL( "auth", "UPDATE Jobs SET archived = 1 WHERE uuid = ? AND (SELECT project_uuid FROM Node WHERE id = Jobs.node_id) = ?", {jobID, projectID} )
  if not result                                                          then return error( 504, "Gateway Timeout" )       end
  if result.ErrorNumber     ~= 0                                         then return error( 403, "Could not update the project" )     end
end

error( 200, "OK" )