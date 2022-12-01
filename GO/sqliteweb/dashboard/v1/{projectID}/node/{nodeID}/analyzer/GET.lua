--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.1
--     //             ///   ///  ///    Date        : 2022/11/28
--    ///             ///   ///  ///    Author      : Andrea Donetti
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : ANALYZER LIST GROUPED
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : List of queries slower than threshold 
--          ////     /////                            grouped by database and sql (normalized_sql)
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/{projectID}/node/6/analyzer

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

command = "ANALYZER LIST GROUPED NODE ?"
commandargs = {machineNodeID}

-- if query.id       then 
--   command = command .. " ID ?"
--   commandargs[#commandargs+1] = query.id
-- if query.groupid       then 
--     command = command .. " ID ?"
--     commandargs[#commandargs+1] = query.groupid
-- elseif query.database       then 
--   command = command .. " DATABASE ?"
--   commandargs[#commandargs+1] = query.database
-- end

queries = executeSQL( projectID, command, commandargs )
if not queries                                then return error( 404, "ProjectID not found" ) end
if queries.ErrorNumber                  ~= 0  then return error( 502, "Bad Gateway" )         end
if queries.NumberOfColumns              ~= 6  then return error( 502, "Bad Gateway" )         end
if queries.NumberOfRows                 <  1  then return error( 200, "OK" )                  end

Response = {
  status            = 200,                        -- status code: 0 = no error, error otherwise
  message           = "OK",                       -- "OK" or error message
  value             = queries.Rows,               -- Array with queries info
}

SetStatus( 200 )
Write( jsonEncode( Response ) )