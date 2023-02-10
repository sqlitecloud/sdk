--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.1
--     //             ///   ///  ///    Date        : 2022/11/28
--    ///             ///   ///  ///    Author      : Andrea Donetti
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : apply the indexes proposed by ANALYZER SUGGEST ID <query_id> APPLY NODE <nodeid>
--   ////                ///  ///                      
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Structure with the sqlite3 expert report
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/dashboard/v1/{projectID}/node/{nodeID}/analyzer/query/{queryID}/suggest

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local nodeID,    err, msg = checkNodeID( nodeID )                        if err ~= 0 then return error( err, msg )                     end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false )                                then return error( 401, "Project Disabled" )           end
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg ) end  
end

local machineNodeID, err, msg = verifyNodeID( userID, projectID, nodeID )    if err ~= 0 then return error( err, msg )                 end

command = "ANALYZER SUGGEST ID ? APPLY NODE ?"
commandargs = {queryID, machineNodeID}

result = executeSQL( projectID, command, commandargs )
if not result                                then return error( 404, "ProjectID not found" ) end
if result.ErrorNumber                  ~= 0  then return error( 502, result.ErrorMessage )   end
if result.NumberOfColumns              ~= 3  then return error( 502, "Bad Gateway" )         end
if result.NumberOfRows                 <  1  then return error( 200, "OK" )                  end

error( 200, "OK" )