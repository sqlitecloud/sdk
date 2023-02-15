SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

print( getINIString( "dashboard", "path", ";-(" ) ) -- /Users/pfeil/GitHub/SqliteCloud/sdk/GO/src/sqliteweb/dashboard))

print( userid ) -- userid from token

-- print( uuid )
result = executeSQL( projectID, "LIST CONNECTIONS" )
nResult = filter( result.Rows, { [ "address"          ] = "Address", 
                                 [ "connection_date"  ] = "ConnectionDate", 
                                 [ "database"         ] = "Database",
                                 [ "id"               ] = "Id",
                                 [ "last_activity"    ] = "LastActivity",
                                 [ "username"         ] = "Username"
                               } )

Response = {
  Args = args,
  Header = header,
  request = body,
  jRequest = jsonDecode( body ),
  Query = query,

  Status = 0,
  Message = "Connections List",
  Connections = nResult,
}

-- do return error( 404, "das wart nix" ) end

template_data = {
  From    = "<will.be@overwritten.com>",
  To      = "andreas.pfeil@web.de",
  Subject = "MySubject"
}

--mail( "welcome.eml", "en", template_data )

Write( jsonEncode( Response ) )
SetStatus( 200 )