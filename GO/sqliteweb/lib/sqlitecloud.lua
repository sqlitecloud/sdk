function checkParameter( parameter, minLength )
  if not parameter                           then return nil, 400, "Missing parameter '%s'"     end
  if string.len( parameter ) < 1             then return nil, 400, "Empty parameter '%s'"       end
  if string.len( parameter ) < minLength     then return nil, 400, "Invalid parameter '%s'"     end
                                                  return parameter, 0, nil
end

function checkNumber( value, minValue, maxValue ) 
  if not value                               then return nil, 400, "Missing '%s'"                                         end
  if string.len( value ) < 1                 then return nil, 400, "Empty '%s'"                                           end

  local tmp = tonumber( value )                         
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
  --print( "Body value " )
  --print( value )
  --print( minLength )
  --print( ".." )

  if not body                                then return nil, 400, "Missing body"                                         end
  if string.len( body ) == 0                 then return nil, 400, "Empty body"                                           end
  
  local jbody = jsonDecode( body )
  if not jbody                               then return nil, 400, "Invalid body"                                         end

  if minLength > 0 then
    if not jbody[ value ]                    then return nil, 400, string.format( "Missing '%s' in body", value )         end
    if string.len( jbody[ value ] ) < minLength then return nil, 400, string.format( "Invalid data in '%s' in body", value ) end
  else
    if not jbody[ value ]                    then return "", 0, nil                                                       end
  end

  -- print( jbody[ value ] )
                                              return jbody[ value ], 0, nil 
end

function checkUserID( userid )               -- Is string and comes from JWT. Contents is a number.
  if not userid                              then return -1, 400, "Invalid UserID"              end
  local uid = tonumber( userid )
  if not uid                                 then return -1, 400, "Invalid UserID"              end                                                                              
  if string.format( "%d", uid ) ~= userid    then return -1, 400, "UserID is Not a Number"      end
  if uid < 0                                 then return -1, 400, "Invalid UserID"              end
                                             return uid, 0, nil 
end

function checkProjectID( uuid )               -- fbf94289-64b0-4fc6-9c20-84083f82ee64
  if not uuid                                then return nil, 400, "Invalid ProjectID"          end
  if uuid == "auth"                          then return nil, 404, "Forbidden ProjectID"        end 
  if string.len( uuid ) ~= 36                then return nil, 400, "Invalid ProjectID"          end 
                                                  return uuid, 0, nil
end

function checkNodeID( nodeid )                -- Is string but MUST contains a number
  if not nodeid                              then return -1, 400, "Invalid NodeID"              end
  local nodeID = tonumber( nodeid ) 
  if not nodeID                              then return -1, 400, "Invalid NodeID"              end                               
  if string.format( "%d", nodeID ) ~= nodeid then return -1, 400, "NodeID is Not a Number"      end
  if nodeID < 0                              then return -1, 400, "Invalid NodeID"              end
                                             return nodeID, 0, nil 
end

------

function verifyUserID( userID )
  local result = executeSQL( "auth", "SELECT User.enabled AND Company.enabled AS enabled, User.company_id FROM User JOIN Company ON User.company_id = Company.id WHERE User.id = ?;", {userID} )

  if not result                     then return -1, -1, 503, "Service Unavailable"  end
  if result.ErrorNumber       ~= 0  then return -1, -1, 502, "Bad Gateway"          end
  if result.NumberOfColumns   ~= 2  then return -1, -1, 502, "Bad Gateway"          end 
  if result.NumberOfRows      ~= 1  then return -1, -1, 404, "Not Found"            end
  if result.Rows[ 1 ].enabled ~= 1  then return -1, -1, 401, "Unauthorized"         end
                                         return userID, result.Rows[ 1 ].company_id, 0, nil
end


function verifyLogin( username, password )
  local query  = "SELECT id, enabled FROM USER WHERE email=? AND password=?;"
  local result = executeSQL( "auth", query, {username, password} )

  if not result                     then return -1, 503, "Service Unavailable"  end
  if result.ErrorNumber       ~= 0  then return -1, 502, "Bad Gateway"          end
  if result.NumberOfColumns   ~= 2  then return -1, 502, "Bad Gateway"          end 
  if result.NumberOfRows      ~= 1  then return -1, 401, "Wrong Credentials"    end
  if result.Rows[ 1 ].enabled ~= 1  then return -1, 401, "Unauthorized"         end
                                         return result.Rows[ 1 ].id, 0, nil
end

function verifyProjectID( userID, projectUUID ) 
  local query  = "SELECT uuid FROM User JOIN Company ON User.company_id = Company.id JOIN Project ON Company.id = Project.company_id WHERE User.enabled=1 AND Company.enabled = 1 AND User.id=? AND Project.uuid = ?;"
  --print( query )
  local result = executeSQL( "auth", query, {userID, projectUUID} )

  if not result                     then return nil, 503, "Service Unavailable" end
  if result.ErrorNumber       ~= 0  then return nil, 502, "Bad Gateway"         end
  if result.NumberOfColumns   ~= 1  then return nil, 502, "Bad Gateway"         end 
  if result.NumberOfRows      < 1   then return nil, 404, "Project Not Found"   end
  if result.NumberOfRows      > 1   then return nil, 502, "Bad Gateway"         end 
                                         return result.Rows[ 1 ].uuid, 0, nil
end

function verifyNodeID( userID, projectUUID, nodeID ) 
  local query  = "SELECT NODE.id, NODE.node_id FROM User JOIN Company ON User.company_id = Company.id JOIN Project ON Company.id = Project.company_id JOIN Node ON Project.uuid = Node.project_uuid WHERE User.enabled = 1 AND User.id = ? AND Project.uuid = ? AND Node.id = ?;"
  --print( query )
  local result = executeSQL( "auth", query, {userID, projectUUID, nodeID} )
  
  if not result                     then return nil, 503, "Service Unavailable" end
  if result.ErrorNumber       ~= 0  then return nil, 502, "Bad Gateway"         end
  if result.NumberOfColumns   ~= 2  then return nil, 502, "Bad Gateway"         end 
  if result.NumberOfRows      < 1   then return nil, 404, "NodeID Not Found"    end
  if result.NumberOfRows      > 1   then return nil, 502, "Bad Gateway"         end 
                                         return result.Rows[ 1 ].node_id, 0, nil
end

function getNodeSettingsID( userID, projectUUID, nodeID, key ) 
  local query  = "SELECT NodeSettings.id FROM User JOIN Company ON User.company_id = Company.id JOIN Project ON Company.id = Project.company_id JOIN Node ON Project.uuid = Node.project_uuid JOIN NodeSettings ON Node.id = NodeSettings.node_id WHERE User.enabled = 1 AND User.id = ? AND Project.uuid = ? AND Node.id = ? AND NodeSettings.key = ?;"
  --print( query )
  local result = executeSQL( "auth", query, {userID, projectUUID, nodeID, key} )

  if not result                     then return nil, 503, "Service Unavailable" end
  if result.ErrorNumber       ~= 0  then return nil, 502, "Bad Gateway"         end
  if result.NumberOfColumns   ~= 1  then return nil, 502, "Bad Gateway"         end 
  if result.NumberOfRows      < 1   then return nil, 404, "Setting Not Found"   end
  if result.NumberOfRows      > 1   then return nil, 502, "Bad Gateway"         end 
                                         return result.Rows[ 1 ].id, 0, nil
end

-- local userID,     errorCode, errorMessage  = verifyLogin( "my.address@domain.com", "password" )     if errorCode ~= 0 then return error( errorCode, errorMessage ) end
-- local uuid,       errorCode, errorMessage  = verifyProjectID( userID, projectID )                     if errorCode ~= 0 then return error( errorCode, errorMessage ) end
-- local nodeID,     errorCode, errorMessage  = verifyNodeID( userID, uuid, 1 )                        if errorCode ~= 0 then return error( errorCode, errorMessage ) end
-- local settingID,  errorCode, errorMessage  = getNodeSettingsID( userID, uuid, 1, "testkeyvalz"  )  if errorCode ~= 0 then return error( errorCode, errorMessage ) end

function deleteProject(projectID, userID)
  -- local command  = "SELECT Node.id AS nodeID, Node.node_id AS clusterNodeID, Node.active FROM Node JOIN Project ON Project.uuid = Node.project_uuid JOIN Company ON Company.id = Project.company_id JOIN User ON User.company_id = Company.id WHERE USER.enabled = 1 AND USER.id = ? AND PROJECT.uuid = ?;"
  local command  = "SELECT id AS nodeID, node_id AS clusterNodeID, active FROM Node WHERE project_uuid = ?;"
  local result = executeSQL( "auth", command, {projectID} ) 
  if not result                                                                       then return 503, "Service Unavailable"  end
  if result.ErrorNumber       ~= 0                                                    then return 502, "Bad Gateway"          end
  if result.NumberOfColumns   ~= 3                                                    then return 502, "Bad Gateway"          end 

  local nodeErrors = 0

  command = ""
  for i = 1, result.NumberOfRows do
    nodeID = result.Rows[ i ].nodeID
    clusterNodeID = result.Rows[ i ].clusterNodeID
    active = result.Rows[ i ].active

    local deleted = true
    if active == 1  then   deleted = deleteNode(userID, nodeID, clusterNodeID, projectID) end
    if deleted      then   command = string.format( "%s DELETE FROM NodeSettings WHERE node_id = %d; DELETE FROM Node WHERE id = %d; DELETE FROM Jobs WHERE node_id = %d;", command, nodeID, nodeID, nodeID ) 
    else                   nodeErrors = nodeErrors + 1    end  
  end  

  if nodeErrors == 0 then command = string.format( "%s DELETE FROM Project WHERE uuid = ?;", command ) end
  
  result = executeSQL( "auth", command, {projectID} )
  if not result                                                                      then return 504, "Gateway Timeout"             end
  if result.ErrorMessage ~= ""                                                       then return 502, result.ErrorMessage           end
  if result.ErrorNumber  ~= 0                                                        then return 502, "Bad Gateway"                 end

  if nodeErrors > 0                                                                  then return 502, string.format( "%d nodes can't be deleted", nodeErrors)             end
  
  return 200, "OK"
end


function error( code, message )
  local result = {
    status  = code,
    message = message
  }
  SetStatus( code )
  SetHeader( "Content-Type", "application/json" )
  SetHeader( "Content-Encoding", "utf-8" )
  Write( jsonEncode( result ) )
end

function bool( data )
  if type( data ) == "boolean"  then return data      end
  if type( data ) == "number"   then return data ~= 0 end
  if type( data ) == "function" then return false     end
  if type( data ) == "nil"      then return false     end

  local  data = string.lower( data )
         
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

function bool_to_number(value)
  return value and 1 or 0
end

-- mask utils

function bit(p)
  return 2 ^ (p - 1)  -- 1-based indexing
end

-- Typical call:  if hasbit(x, bit(3)) then ...
function hasbit(x, p)
  return x % (p + p) >= p       
end

function setbit(x, p)
  return hasbit(x, p) and x or x + p
end

function clearbit(x, p)
  return hasbit(x, p) and x - p or x
end
