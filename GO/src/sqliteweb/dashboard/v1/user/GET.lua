-- Get all data and settings for logged in user

if tonumber( userid ) < 0 then return error( 416, "Invalid User" ) end

Setting = {
  key   = "",
  value = ""
}

Response = {
  status           = 0,                         -- status code: 0 = no error, error otherwise
  message          = "OK",                      -- "OK" or error message

  id               = tonumber( userid ),        -- UserID, 0 = static user defined in .ini file
  enabled          = false,                     -- Whether this user account is enabled or disabled
  name             = "",                        -- User name
  company          = "",                        -- User company
  email            = "",                        -- User email - also used as login
  password         = "*******",                 -- User password - this fiels is always 7 stars
  creationDate     = "1970-01-01 00:00:00",     -- Date and time when this user account was created
  lastRecoveryTime = "1970-01-01 00:00:00",     -- Last date and time when this user has tried to recover his password

  settings         = nil,
}

if tonumber( userid ) == 0 then 
  Response.enabled          = bool( getINIString( "dashboard", "enabled", "false" ) )
  Response.name             = getINIString( "dashboard", "name", "unknown" )
  Response.company          = getINIString( "dashboard", "company", "unknown" )
  Response.email            = getINIString( "dashboard", "email", "unknown" )
  Response.creationDate     = getINIString( "dashboard", "modified", "1970-01-01 00:00:00" )
  Response.lastRecoveryTime = getINIString( "dashboard", "modified", "1970-01-01 00:00:00" )
else
  data = executeSQL( "auth", string.format( "SELECT 0 AS status, 'OK' AS message, id, enabled, name, company, email, '*******' AS password, creation_date AS creationDate, last_recovery_request AS lastRecoveryTime from USER WHERE id = %d ;", userid ) )
    if data.ErrorNumber == 0 and data.NumberOfRows == 1 then
      Response = data.Rows[ 1 ]
    Response.enabled = bool( Response.enabled )

    data = executeSQL( "auth", string.format( "SELECT key, value FROM USER_SETTINGS WHERE user_id = %d ;", userid ) )
    if data.ErrorNumber == 0 and data.NumberOfRows > 0 then
      Response.settings = data.Rows
    end
  else
    return error( 500, "Internal Server Error (actual user does not exist)!?!" )
  end
end

SetStatus( 200 )
Write( jsonEncode( Response ) )