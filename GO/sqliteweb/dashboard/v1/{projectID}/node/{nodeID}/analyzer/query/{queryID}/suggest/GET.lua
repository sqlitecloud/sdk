--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.1
--     //             ///   ///  ///    Date        : 2022/11/28
--    ///             ///   ///  ///    Author      : Andrea Donetti
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : ANALYZER SUGGEST ID <query_id> [PERCENTAGE <percentage>] NODE <nodeid>
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

command = "ANALYZER SUGGEST ID ?"
commandargs = {queryID}

if query.percentage       then 
  command = command .. " PERCENTAGE ?"
  commandargs[#commandargs+1] = query.percentage
end
-- if query.apply        then 
--   command = command .. " APPLY"
-- end

command = command .. " NODE ?"
commandargs[#commandargs+1] = machineNodeID

result = executeSQL( projectID, command, commandargs )
if not result                                then return error( 404, "ProjectID not found" ) end
if result.ErrorNumber                  ~= 0  then return error( 502, "Bad Gateway" )         end
if result.NumberOfColumns              ~= 3  then return error( 502, "Bad Gateway" )         end
if result.NumberOfRows                 <  1  then return error( 200, "OK" )                  end

statements = nil

if #result.Rows == 0 then 
  result = nil  
else 
  statements = {}

  for i = 1, #result.Rows do 
    local record = result.Rows[ i ]
    local statementid = record.statement + 1
    local statement = statements[statementid]
    if not statement then 
      statement = {} 
      statements[statementid] = statement
    end

    local t = nil
    if     record.type == 1 then t = "sql"
    elseif record.type == 2 then t = "indexes"
    elseif record.type == 3 then t = "plan"
    elseif record.type == 4 then t = "candidates"
    end

    if t then statement[t] = record.report end
  end
end


Response = {
  status            = 200,                        -- status code: 0 = no error, error otherwise
  message           = "OK",                       -- "OK" or error message
  value             = statements,                 -- array of statements, each statement is a map with 4 string reports, one for each key (sql, indexes, plan, candidates)
}

SetStatus( 200 )
Write( jsonEncode( Response ) )