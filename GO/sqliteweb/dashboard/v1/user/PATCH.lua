--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Update my dashboard user
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : status + message
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/projects

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,       err, msg = checkUserID( userid )                     if err ~= 0 then return error( err, msg )                                end

local first_name,   err, msg = getBodyValue( "first_name", 0 )           if err ~= 0 then return error( err, msg )                                end
local last_name,    err, msg = getBodyValue( "last_name", 0 )            if err ~= 0 then return error( err, msg )                                end
local email,        err, msg = getBodyValue( "email", 0 )                if err ~= 0 then return error( err, msg )                                end
local password,     err, msg = getBodyValue( "password", 0 )             if err ~= 0 then return error( err, msg )                                end

if userID == 0 then 
  return error( 501, "Not Implemented" )
else
  command = string.format( "UPDATE User SET ")
  commandargs = {}
  local separator = ""

  fields = {first_name = first_name, last_name = last_name, email = email, password = hash(password)}

  for k, v in pairs(fields) do
    if v and string.len(v)>0 then  
      command = string.format( "%s%s %s = ?", command, separator, k)
      commandargs[#commandargs+1] = v
      separator = ","
    end
  end

  -- only if there are any changes
  if #commandargs > 0 then
    command = string.format( "%s WHERE id = %d; SELECT changes() AS success;", command, userID)
    result = executeSQL( "auth", command, commandargs )
    if not result                                                                      then return error( 504, "Gateway Timeout" )                  end
    if result.ErrorMessage      ~= ""                                                  then return error( 502, result.ErrorMessage )                end
    if result.ErrorNumber       ~= 0                                                   then return error( 502, "Bad Gateway" )                      end
    if result.NumberOfRows      ~= 1                                                   then return error( 502, "Bad Gateway" )                      end
    if result.NumberOfColumns   ~= 1                                                   then return error( 502, "Bad Gateway" )                      end
    if result.Rows[ 1 ].success ~= 1                                                   then return error( 500, "User not found" )              end
  end
end

error( 200, "OK" )