--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/04/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : 
--   ////                ///  ///                     
--     ////     //////////   ///                      
--        ////            ////          Requires    : Authentication
--          ////     /////              Output      : 
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/admin/v1/user/{email}

-- -- SELECT NODE.id AS nodeID FROM Node JOIN Project ON Project.uuid = Node.project_uuid JOIN Company ON Company.id = Project.company_id JOIN User ON User.company_id = Company.id WHERE USER.email='my.address@domain.com'
-- -- DELETE FROM NodeSettings WHERE node_id = id

-- -- SELECT id FROM USER WHERE email =
-- -- DELETE FROM UserSettings WHERE user_id = id
-- -- DELETE FROM PROJECT WHERE user_id = id
-- DELETE FROM USER WHERE id = id 

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local email,    err, msg = checkParameter( email, 3 )                    if err ~= 0 then return error( err, string.format( msg, "email" ) )  end

query  = string.format( "SELECT id || '' AS id FROM USER WHERE email = '%s';", enquoteSQL( email ) )
userID = executeSQL( "auth", query )
if not userID                                                                        then return error( 504, "Gateway Timeout" )              end
if userID.ErrorMessage    ~= ""                                                      then return error( 502, userID.ErrorMessage )            end
if userID.ErrorNumber     ~= 0                                                       then return error( 502, "Bad Gateway" )                  end
if userID.NumberOfColumns ~= 1                                                       then return error( 502, "Bad Gateway" )                  end
if userID.NumberOfRows    ~= 1                                                       then return error( 404, "User not found" )               end

local userID,    err, msg = checkUserID( userID.Rows[ 1 ].id  )          if err ~= 0 then return error( err, msg )                            end

-- TODO: missing APIs to manage the Customer table and to remove projects and node for deleted companies, 
--       we cannot delete projects and nodes here with the following code because other users can still 
--       exist for the same Company

-- query = string.format( "SELECT NODE.id AS nodeID FROM Node JOIN Project ON Project.uuid = Node.project_uuid JOIN Company ON Company.id = Project.company_id JOIN User ON User.company_id = Company.id WHERE User.id ='%d';", userID )
-- nodes = executeSQL( "auth", query )

-- if not nodes                                                                         then return error( 504, "Gateway Timeout" )              end
-- if nodes.ErrorMessage    ~= ""                                                       then return error( 502, userID.ErrorMessage )            end
-- if nodes.ErrorNumber     ~= 0                                                        then return error( 502, "Bad Gateway" )                  end

transaction = "BEGIN TRANSACTION;"

-- if nodes.NumberOfRows > 0 then
--   if nodes.NumberOfColumns ~= 1                                                      then return error( 502, "Bad Gateway" )                  end

--   for i = 1, nodes.NumberOfRows do
--     transaction = string.format( "%s DELETE FROM NodeSettings WHERE node_id=%d;", transaction, nodes.Rows[ i ].id )
--   end
-- end

transaction = string.format( "%s DELETE FROM UserSettings WHERE user_id=%d;", transaction, userID )
-- transaction = string.format( "%s DELETE FROM Project WHERE user_id=%d;"     , transaction, userID )
transaction = string.format( "%s DELETE FROM User WHERE id=%d;"             , transaction, userID )
transaction = string.format( "%s COMMIT TRANSACTION;"                       , transaction         )

result = executeSQL( "auth", transaction )
if not result                                                                        then return error( 504, "Gateway Timeout" )              end
if result.ErrorMessage ~= ""                                                         then return error( 502, result.ErrorMessage )            end
if result.ErrorNumber  ~= 0                                                          then return error( 502, "Bad Gateway" )                  end
if result.Value        ~= "OK"                                                       then return error( 500, "Internal Server Error" )        end

error( 200, "OK" )