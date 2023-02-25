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

userID = executeSQL( "auth", "SELECT id || '' AS id FROM USER WHERE email = ?;", { email } )
if not userID                                                                        then return error( 504, "Gateway Timeout" )              end
if userID.ErrorMessage    ~= ""                                                      then return error( 502, userID.ErrorMessage )            end
if userID.ErrorNumber     ~= 0                                                       then return error( 502, "Bad Gateway" )                  end
if userID.NumberOfColumns ~= 1                                                       then return error( 502, "Bad Gateway" )                  end
if userID.NumberOfRows    ~= 1                                                       then return error( 404, "User not found" )               end

local userID,    err, msg = checkUserID( userID.Rows[ 1 ].id  )          if err ~= 0 then return error( err, msg )                            end

-- save the company id
command = string.format( "SELECT company_id FROM User WHERE id = %d", userID)
company = executeSQL( "auth", command )
if not company                                                                       then return error( 504, "Gateway Timeout" )              end
if company.ErrorMessage    ~= ""                                                     then return error( 502, company.ErrorMessage )            end
if company.ErrorNumber     ~= 0                                                      then return error( 502, "Bad Gateway" )                  end
if company.NumberOfRows    ~= 1                                                      then return error( 502, "Bad Gateway" )                  end
companyID = company.Rows[ 1 ].company_id

command = string.format( "DELETE FROM User WHERE id = %s; DELETE FROM UserSettings WHERE user_id = %d", userID, userID)
result = executeSQL( "auth", command )
if not result                                                                       then return error( 504, "Gateway Timeout" )              end
if result.ErrorMessage    ~= ""                                                     then return error( 502, result.ErrorMessage )            end
if result.ErrorNumber     ~= 0                                                      then return error( 502, "Bad Gateway" )                  end

-- delete the user and the company only if the user was the last user. 
-- use a transaction to guarantee no concurrency issues.
-- If the company has been deleted the query returns 1, otherwise the query returns 0
command = string.format( "BEGIN IMMEDIATE; DELETE FROM Company WHERE id = (SELECT IIF((SELECT COUNT(User.id) FROM User WHERE company_id = %d) = 0, %d, 0)); COMMIT; SELECT changes() as deleted;", companyID, companyID)
result = executeSQL( "auth", command )
if not result                                                                       then return error( 504, "Gateway Timeout" )              end
if result.ErrorMessage    ~= ""                                                     then return error( 502, result.ErrorMessage )            end
if result.ErrorNumber     ~= 0                                                      then return error( 502, "Bad Gateway" )                  end
if result.NumberOfRows    ~= 1                                                      then return error( 502, "Bad Gateway" )                  end
must_delete_projects = result.Rows[ 1 ].deleted == 1

code = 200
message = ""
if must_delete_projects then 
    pcommand = string.format( "SELECT uuid FROM Project WHERE company_id = %s", companyID)
    project = executeSQL( "auth", pcommand )
    if not project                                                                       then return error( 504, "Gateway Timeout" )              end
    if project.ErrorMessage    ~= ""                                                     then return error( 502, userID.ErrorMessage )            end
    if project.ErrorNumber     ~= 0                                                      then return error( 502, "Bad Gateway" )                  end

    command = ""
    for i = 1, project.NumberOfRows do
        projectID = project.Rows[ i ].uuid
        c, m = deleteProject(projectID, userID) 
        if c ~= 200 then 
            code = c
            separator = "; "
            if message == "" then separator = "" end
            message = string.format( "%s%s" , separator, m)
        end
    end
end

if code == 200 then message = "OK" end

error( code, message )