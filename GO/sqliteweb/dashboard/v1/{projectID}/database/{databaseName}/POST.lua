--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : CREATE DATABASE % [KEY %] 
--   ////                ///  ///                     [ENCODING %] [IF NOT EXISTS]
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/{databaseName}

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                                    end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                                    end
local dbName,    err, msg = checkParameter( databaseName, 1 )            if err ~= 0 then return error( err, string.format( msg, "databaseName" ) )   end
local key,       err, msg = getBodyValue( "key", 0 )                     if err ~= 0 then return error( err, msg )                                    end
local encoding,  err, msg = getBodyValue( "encoding", 0 )                if err ~= 0 then return error( err, msg )                                    end

local query = "CREATE DATABASE ?"
local queryargs = {dbName}

if string.len( key )      > 0 then 
  query = string.format( "%s KEY ?", query ) 
  queryargs[#queryargs+1] = key
end

if string.len( encoding ) > 0 then 
  query = string.format( "%s ENCODING ?", query ) 
  queryargs[#queryargs+1] = encoding
end

query = string.format( "%s IF NOT EXISTS;", query )

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg ) end  
end

result = executeSQL( projectID, query, queryargs )
if not result                             then return error( 404, "ProjectID not found" ) end
if result.ErrorNumber       ~= 0          then return error( 502, result.ErrorMessage )   end
if result.NumberOfColumns   ~= 0          then return error( 502, "Bad Gateway" )         end
if result.NumberOfRows      ~= 0          then return error( 502, "Bad Gateway" )         end

if result.Value             ~= "OK"       then return error( 404, result.Value )          end

error( 200, "OK" )