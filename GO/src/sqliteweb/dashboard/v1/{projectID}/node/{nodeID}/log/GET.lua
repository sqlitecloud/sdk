-- Filter log

userid = 1


function getNode

-- LIST LOG FROM % TO % [LEVEL %] [TYPE %] [ORDER DESC]    
-- LIST % ROWS FROM LOG [LEVEL %] [TYPE %] [ORDER DESC]


if tonumber( userid ) < 0 then return error( 416, "Invalid User" ) end
if projectID == "auth" then return error( 500, "Internal server error" ) end
-- TODO: if nodeID == bad 


level = query.level
type  = query.type
order = query.order
if not level  then level  = "0"                       end
if not type   then type   = "4"                       end
if not order  then order  = "ORDER DESC"              end

rows  = query.rows
if rows then
  sql = string.format( "LIST %d ROWS FROM LOG LEVEL %d TYPE %d %s", rows, level, type, order )
else
  from  = query.from
  to    = query.to
  -- dates = executeSQL( "auth", "SELECT DATETIME() AS now, DATETIME( 'now', '-1 day' ) AS yesterday;" )
  -- if not from   then from   = dates.Rows[ 1 ].yesterday end
  -- if not to     then to     = dates.Rows[ 1 ].now       end
  if not to     then to     = "DATETIME()"                  end
  if not from   then from   = "DATETIME( 'now', '-1 day' )" end

  sql = string.format( "LIST LOG FROM %s TO %s LEVEL %d TYPE %d ORDER %s;", from, to, tonumber( level ), tonumber( type ), order ) 
end

print( sql )




log = executeSQL( "auth", sql )
flog = filter( log.Rows, { [ "datetime"    ] = "date", 
                           [ "log_type"    ] = "type", 
                           [ "log_level"   ] = "level",
                           [ "description" ] = "description",
                           [ "username"    ] = "username",
                           [ "database"    ] = "database",
                           [ "ip_address"  ] = "address",
                         } )

Response = {
  status            = 0,                         -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  logs              = flog                       -- Array with key value pairs
}



goto done





database = query.database



print( type )
print( userid )
print( projectID )
print( nodeID )



if tonumber( userid ) < 0 then return error( 416, "Invalid User" ) end
if projectID == "auth" then return error( 500, "Internal server error" ) end
-- TODO: if nodeID == bad 



Response = {
  status            = 0,                         -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  logs              = nil,                        -- Array with key value pairs
}

if tonumber( userid ) == 0 then 

  if getINIBoolean( projectID, "enabled", false ) then
    goto done
  end

  do return error( 500, "NodeID not found" ) end -- TODO: Other error code

else

  query = string.format( "SELECT key, value FROM USER JOIN PROJECT ON USER.id = PROJECT.user_id JOIN NODE ON PROJECT.uuid = NODE.project_uuid JOIN NODE_SETTINGS ON NODE.id = node_id WHERE USER.enabled = 1 AND USER.id = %d AND NODE.id = %d AND uuid='%s';", userid, nodeID, enquoteSQL( projectID ) )
  
  print( query )
  settings = executeSQL( "auth", query )

  if settings.ErrorNumber == 0 then
    if settings.NumberOfRows > 0 then
      Response.settings = settings.Rows
    end
  else
    return error( 500, "NodeID not found" )
  end

end

::done::

SetStatus( 200 )
Write( jsonEncode( Response ) )

