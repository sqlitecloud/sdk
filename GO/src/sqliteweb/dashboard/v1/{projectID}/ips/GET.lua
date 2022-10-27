--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : LIST ALLOWED IP [ROLE %] [USER %]
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Structure with IP-info
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/ips

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )       if err ~= 0 then return error( err, msg )  end
local projectID, err, msg = checkProjectID( projectID ) if err ~= 0 then return error( err, msg )  end

if not query.role then role = "" else role = query.role user = ""               end
if not query.user then user = "" else role = ""         user = query.user       end

query = "LIST ALLOWED IP "

if role ~= "" then 
  query = query .. " ROLE ?" 
  queryargs = {role}
elseif user ~= "" then 
  query = query .. " USER ?"
  queryargs = {user}
end

ips   = nil

if userID == 0 then
  if not getINIBoolean( projectID, "enabled", false ) then return error( 401, "Disabled project" ) end
else
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg ) end
end

ips = executeSQL( projectID, query, queryargs )

if not ips                                then return error( 404, "ProjectID not found" ) end
if ips.ErrorNumber                  ~= 0  then return error( 502, "Bad Gateway" )         end
if ips.NumberOfColumns              ~= 3  then return error( 502, "Bad Gateway" )         end
if ips.NumberOfRows                 <  1  then return error( 200, "OK" )                  end

IP = {
  address = "127.0.0.1",                         -- IPv[4/6]
  name    = "name",                              -- Name
  type    = "type";                              -- Type
}

Response = {
  status            = 200,                       -- status code: 200 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message
  value             = ips.Rows,                  -- Array with allowed IP's for this role or user
}

SetStatus( 200 )
Write( jsonEncode( Response ) )