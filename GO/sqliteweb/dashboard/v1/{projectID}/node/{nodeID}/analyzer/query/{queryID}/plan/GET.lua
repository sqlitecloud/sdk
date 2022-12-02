--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.1
--     //             ///   ///  ///    Date        : 2022/11/28
--    ///             ///   ///  ///    Author      : Andrea Donetti
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : ANALYZER PLAN ID <query_id> NODE <nodeid>
--   ////                ///  ///                      
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Structure with the EXPLAIN QUERY PLAN report
--          ////     /////                            Array of records, each record can have a "records" field with the array of inner records, if any
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/dashboard/v1/{projectID}/node/{nodeID}/analyzer/query/{queryID}/plan

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

command = "ANALYZER PLAN ID ? NODE ?"
commandargs = {queryID,machineNodeID}

result = executeSQL( projectID, command, commandargs )
if not result                                then return error( 404, "ProjectID not found" ) end
if result.ErrorNumber                  ~= 0  then return error( 502, "Bad Gateway" )         end
if result.NumberOfColumns              ~= 4  then return error( 502, "Bad Gateway" )         end
if result.NumberOfRows                 <  1  then return error( 200, "OK" )                  end

-- records is the tree structure for records
-- it contains an array of record object, and each record object can
-- contain a "records" field with the array of inner records, if any
records = nil
-- recordsbyid let me directly get a record object from its id 
-- when I look for a parent record
recordsbyid = nil

if #result.Rows == 0 then 
  result = nil  
else 
  records = {}
  recordsbyid = {}

  for i = 1, #result.Rows do 
    local record = result.Rows[ i ]    
    recordsbyid[record.id] = record

    if record.parent == 0 then
      records[#records+1] = record
    else 
      parent = recordsbyid[record.parent]
      if parent then 
        if not parent.records then  parent.records = {} end
        parent.records[#parent.records + 1] = record
      end
    end
  end
end


Response = {
  status            = 200,                        -- status code: 0 = no error, error otherwise
  message           = "OK",                       -- "OK" or error message
  value             = records,                    -- array of records, each record can contain a "records" field with the array of inner records
}

SetStatus( 200 )
Write( jsonEncode( Response ) )