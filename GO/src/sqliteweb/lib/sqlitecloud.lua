function checkParameter( parameter, minLength )
  if not parameter                           then return nil, 400, "Missing parameter '%s'"     end
  if string.len( parameter ) < 1             then return nil, 400, "Empty parameter '%s'"       end
  if string.len( parameter ) < minLength     then return nil, 400, "Invalid parameter '%s'"     end
                                                  return parameter, 0, nil
end

function checkNumber( value, minValue, maxValue ) 
  if not value                               then return nil, 400, "Missing '%s'"                                         end
  if string.len( value ) < 1                 then return nil, 400, "Empty '%s'"                                           end

  tmp = tonumber( value )                         
  if string.format( "%d", tmp ) ~= value     then return nil, 400, "'%s' is not a number"                                 end
  if tmp < minValue                          then return nil, 400, string.format( "'%%s' is less than %d", minValue )     end
  if tmp > maxValue                          then return nil, 400, string.format( "'%%s' is greater than %d", maxValue )  end
                                                  return tmp, 0, nil
end

function checkDateTime( value ) 
  if not value                               then return nil, 400, "Missing '%s'"                                         end
  if string.len( value ) < 1                 then return nil, 400, "Empty '%s'"                                           end
  if string.len( value ) ~= 19               then return nil, 400, "Error in format of '%s'"                              end -- 2022-03-01 00:00:00
  if value:sub( 5, 5 ) ~= "-"                then return nil, 400, "Error in format of '%s'"                              end
  if value:sub( 8, 8 ) ~= "-"                then return nil, 400, "Error in format of '%s'"                              end
  if value:sub( 11, 11 ) ~= " "              then return nil, 400, "Error in format of '%s'"                              end
  if value:sub( 14, 14 ) ~= ":"              then return nil, 400, "Error in format of '%s'"                              end
  if value:sub( 17, 17 ) ~= ":"              then return nil, 400, "Error in format of '%s'"                              end

  for i = 1, #value do
    c = value:sub( i, i )
    if c == "0" or c == "1" or c == "2" or c == "3" or c == "4" or c == "5" or c == "6" or c == "7" or c == "8" or c == "9" or c == "-" or c == " " or c == ":" then goto next end
    do return nil, 400, "Error in format of '%s'" end
    ::next::
  end
                                                  return value, 0, nil
end

function getBodyValue( value, minLength )     
  if not body                                then return nil, 400, "Missing body"                                         end
  if string.len( body ) == 0                 then return nil, 400, "Empty body"                                           end

  local jbody = jsonDecode( body )
  if not jbody                                then return nil, 400, "Invalid body"                                         end

  if not jbody[ value ]                       then return nil, 400, string.format( "Missing '%s' in body", value )         end
  if string.len( jbody[ value ] ) < minLength then return nil, 400, string.format( "Invalid data in '%s' in body", value ) end


  print( "Body value " )
  print( value )
  print( jbody[ value ] )
                                              return jbody[ value ], 0, nil 
end

function checkUserID( userid )               -- Is string and comes from JWT. Contents is a number.
  if not userid                              then return -1, 400, "Invalid UserID"              end
  userID = tonumber( userid )                                               
  if string.format( "%d", userID ) ~= userid then return -1, 400, "UserID is Not a Number"      end
  if userID < 0                              then return -1, 400, "Invalid UserID"              end
                                             return userID, 0, nil 
end

function checkProjectID( uuid )               -- fbf94289-64b0-4fc6-9c20-84083f82ee64
  if not uuid                                then return nil, 400, "Invalid ProjectID"          end
  if uuid == "auth"                          then return nil, 404, "Forbidden ProjectID"        end 
  if string.len( uuid ) ~= 36                then return nil, 400, "Invalid ProjectID"          end 
                                                  return uuid, 0, nil
end

function checkNodeID( nodeid )                -- Is string but MUST contains a number
  if not nodeid                              then return -1, 400, "Invalid NodeID"              end
  nodeID = tonumber( nodeid )                                               
  if string.format( "%d", nodeID ) ~= nodeid then return -1, 400, "NodeID is Not a Number"      end
  if nodeID < 0                              then return -1, 400, "Invalid NodeID"              end
                                             return nodeID, 0, nil 
end

------

function verifyUserID( userID )
  result = executeSQL( "auth", string.format( "SELECT enabled FROM USER WHERE id = %d;", userID ) )

  if not result                     then return -1, 503, "Service Unavailable"  end
  if result.ErrorNumber       ~= 0  then return -1, 502, "Bad Gateway"          end
  if result.NumberOfColumns   ~= 1  then return -1, 502, "Bad Gateway"          end 
  if result.NumberOfRows      ~= 1  then return -1, 404, "Not Found"            end
  if result.Rows[ 1 ].enabled ~= 1  then return -1, 401, "Unauthorized"         end
                                         return userID, 0, nil
end


function verifyLogin( username, password )
  query  = string.format( "SELECT id, enabled FROM USER WHERE email='%s' AND password='%s';", enquoteSQL( username ), enquoteSQL( password ) )
  result = executeSQL( "auth", query )

  if not result                     then return -1, 503, "Service Unavailable"  end
  if result.ErrorNumber       ~= 0  then return -1, 502, "Bad Gateway"          end
  if result.NumberOfColumns   ~= 2  then return -1, 502, "Bad Gateway"          end 
  if result.NumberOfRows      ~= 1  then return -1, 401, "Wrong Credentials"    end
  if result.Rows[ 1 ].enabled ~= 1  then return -1, 401, "Unauthorized"         end
                                         return result.Rows[ 1 ].id, 0, nil
end

function verifyProjectID( userID, projectUUID ) 
  query  = string.format( "SELECT uuid FROM USER JOIN PROJECT ON USER.id = PROJECT.user_id WHERE USER.enabled=1 AND USER.id=%d AND PROJECT.uuid = '%s';", userID, enquoteSQL( projectUUID ) )
  print( query )
  result = executeSQL( "auth", query )

  if not result                     then return nil, 503, "Service Unavailable" end
  if result.ErrorNumber       ~= 0  then return nil, 502, "Bad Gateway"         end
  if result.NumberOfColumns   ~= 1  then return nil, 502, "Bad Gateway"         end 
  if result.NumberOfRows      < 1   then return nil, 404, "Project Not Found"   end
  if result.NumberOfRows      > 1   then return nil, 502, "Bad Gateway"         end 
                                         return result.Rows[ 1 ].uuid, 0, nil
end

function verifyNodeID( userID, projectUUID, nodeID ) 
  query  = string.format( "SELECT NODE.id  FROM USER JOIN PROJECT ON USER.id = PROJECT.user_id JOIN NODE ON PROJECT.uuid = NODE.project_uuid WHERE USER.enabled = 1 AND USER.id=%d AND PROJECT.uuid = '%s' AND NODE.id = %d;", userID, enquoteSQL( projectUUID ), nodeID )
  print( query )
  result = executeSQL( "auth", query )

  if not result                     then return nil, 503, "Service Unavailable" end
  if result.ErrorNumber       ~= 0  then return nil, 502, "Bad Gateway"         end
  if result.NumberOfColumns   ~= 1  then return nil, 502, "Bad Gateway"         end 
  if result.NumberOfRows      < 1   then return nil, 404, "NodeID Not Found"    end
  if result.NumberOfRows      > 1   then return nil, 502, "Bad Gateway"         end 
                                         return result.Rows[ 1 ].id, 0, nil
end

function getNodeSettingsID( userID, projectUUID, nodeID, key ) 
  query  = string.format( "SELECT NODE_SETTINGS.id FROM USER JOIN PROJECT ON USER.id = PROJECT.user_id JOIN NODE ON PROJECT.uuid = NODE.project_uuid JOIN NODE_SETTINGS ON NODE.id = NODE_SETTINGS.node_id WHERE USER.enabled = 1 AND USER.id=%d AND PROJECT.uuid = '%s' AND NODE.id = %d AND NODE_SETTINGS.key='%s';", userID, enquoteSQL( projectUUID ), nodeID, enquoteSQL( key ) )
  print( query )
  result = executeSQL( "auth", query )

  if not result                     then return nil, 503, "Service Unavailable" end
  if result.ErrorNumber       ~= 0  then return nil, 502, "Bad Gateway"         end
  if result.NumberOfColumns   ~= 1  then return nil, 502, "Bad Gateway"         end 
  if result.NumberOfRows      < 1   then return nil, 404, "Setting Not Found"   end
  if result.NumberOfRows      > 1   then return nil, 502, "Bad Gateway"         end 
                                         return result.Rows[ 1 ].id, 0, nil
end

-- local userID,     errorCode, errorMessage  = verifyLogin( "my.address@domain.com", "password" )     if errorCode ~= 0 then return error( errorCode, errorMessage ) end
-- local uuid,       errorCode, errorMessage  = verifyProject( userID, projectID )                     if errorCode ~= 0 then return error( errorCode, errorMessage ) end
-- local nodeID,     errorCode, errorMessage  = verifyNodeID( userID, uuid, 1 )                        if errorCode ~= 0 then return error( errorCode, errorMessage ) end
-- local settingID,  errorCode, errorMessage  = getNodeSettingsID( userID, uuid, 1, "testkeyvalz"  )  if errorCode ~= 0 then return error( errorCode, errorMessage ) end



function getNumberOfConnections( projectID, databaseName )
  query = string.format( "LIST DATABASE CONNECTIONS '%s'; ", enquoteSQL( databaseName) )

  if not query                            then  return 0 end
  if query.ErrorNumber ~= 0               then  return 0 end
  if query.NumberOfColumns ~= 2           then  return 0 end
                                                return query.NumberOfRows
end

function error( code, message )
  result = {
    status  = code,
    message = message
  }
  SetStatus( code )
  SetHeader( "Content-Type", "application/json" )
  SetHeader( "Content-Encoding", "utf-8" )
  Write( jsonEncode( result ) )
end

function bool( data )
  data = string.lower( data )
  if     data == "1"        then return true
	elseif data == "true"     then return true
	elseif data == "enable"   then return true
	elseif data == "enabled"  then return true
	else                           return false
  end
end

function contains( str, needle )
  for i = 1, #str do
    if str:sub( i, i ) == needle then return true end
  end
  return false
end