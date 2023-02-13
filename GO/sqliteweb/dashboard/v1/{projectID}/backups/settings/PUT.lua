--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Modify the backup settings 
--   ////                ///  ///                     of all the databases
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/backups/settings

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end
local values,    err, msg = getBodyValue( "values", 0 )                  if err ~= 0 then return error( err, msg )                          end 

-- -- get the array of values from the body json
-- if not body                                then error( 400, "Missing body" ) return                                     end
-- if string.len( body ) == 0                 then error( 400, "Empty body" ) return                                       end

-- local jbody = jsonDecode( body )
-- if not jbody                               then error( 400, "Invalid body" ) return                                     end

-- local values = jbody[values" ) 
-- if not values or #values == 0              then error( 400, "Missing or empty values" ) return                          end

if userID == 0 then         
  if not getINIBoolean( projectID, "enabled", false )                                then return error( 401, "Project Disabled" )           end
                                                                                          return error( 501, "Not Implemented" )  
else     
  local projectID, err, msg = verifyProjectID( userID, projectID )       if err ~= 0 then return error( err, msg )                          end

--   print("values:" .. #values)

  local c = ""
  local args = {}
  for i = 1, #values do
    local value = values[ i ]
    local name = value[ "name" ]
    -- print("values i:" .. i .. " name:" .. name)

    if name then
      local enabled = value["enabled"] 
      if enabled == nil or enabled == "" then enabled = 0 end
      if type(enabled) == "string" then 
        enabled,     err, msg = checkNumber( enabled, 0, 1 )           if err ~= 0 then return error( err, string.format( msg, "enabled" ) ) end
      end
      c = c .. "SET DATABASE ? KEY backup TO ?;"
      args[#args+1] = name
      args[#args+1] = enabled

      local retention = value["backup_retention"]
      if retention and string.len(retention) then 
        c = c .. " SET DATABASE ? KEY backup_retention TO ?;"
        args[#args+1] = name
        args[#args+1] = retention
      else 
        c = c .. " REMOVE DATABASE ? KEY backup_retention; "
        args[#args+1] = name
      end 
    end
  end

  if #values > 0 then
    c = c .. " APPLY BACKUP SETTINGS;"
    -- print("command: ".. c)
    result = executeSQL( projectID, c, args )
    if not result                                   then return error( 504, "Gateway Timeout" )            end
    if result.ErrorNumber ~= 0                      then return error( 502, result.ErrorMessage )          end
  end

end

error( 200, "OK" )