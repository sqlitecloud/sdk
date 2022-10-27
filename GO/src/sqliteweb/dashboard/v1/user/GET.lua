--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Get all data and settings 
--   ////                ///  ///                     for logged in user
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Structure with user settings
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                     end

Setting = {
  key   = "",
  value = ""
}

Response = {
  status           = 200,                              -- status code: 0 = no error, error otherwise
  message          = "OK",                             -- "OK" or error message
  value            = {
    id                    = tonumber( userid ),        -- UserID, 0 = static user defined in .ini file
    enabled               = false,                     -- Whether this user account is enabled or disabled
    first_name            = "",                        -- User name
    last_name             = "",                        -- Last name
    company               = "",                        -- User company
    email                 = "",                        -- User email - also used as login
    creation_date         = "1970-01-01 00:00:00",     -- Date and time when this user account was created
    last_recovery_request = "1970-01-01 00:00:00",     -- Last date and time when this user has tried to recover his password
    settings              = nil,
  }
}

if userID == 0 then    
  Response.value.enabled               = getINIBoolean( "dashboard", "enabled", false )
  Response.value.first_name            = getINIString( "dashboard", "first_name", "unknown" )
  Response.value.last_name             = getINIString( "dashboard", "last_name", "unknown" )
  Response.value.company               = getINIString( "dashboard", "company", "unknown" )
  Response.value.email                 = getINIString( "dashboard", "email", "unknown" )
  Response.value.creation_date         = getINIString( "dashboard", "modified", "1970-01-01 00:00:00" )
  Response.value.last_recovery_request = getINIString( "dashboard", "modified", "1970-01-01 00:00:00" )
  
else
  data = executeSQL( "auth", "SELECT User.id, User.enabled, first_name, last_name, Company.name as company, email, creation_date, last_recovery_request FROM User JOIN Company ON User.company_id = Company.id WHERE User.id = ?;", {userID} )
  if data.ErrorNumber == 0 and data.NumberOfRows == 1 then
    Response.value = data.Rows[ 1 ]
    Response.value.enabled = bool( Response.value.enabled )

    data = executeSQL( "auth", "SELECT key, value FROM UserSettings WHERE user_id = ?;", {userID} )
    if data.ErrorNumber == 0 and data.NumberOfRows > 0 then
      Response.value.settings = data.Rows
    end
  else
    return error( 500, "Internal Server Error (actual user does not exist)!?!" )
  end
end

SetStatus( 200 )
Write( jsonEncode( Response ) )