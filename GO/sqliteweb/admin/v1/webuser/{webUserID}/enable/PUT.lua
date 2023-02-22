--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/04/11
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : enable the specified webuser
--   ////                ///  ///                     
--     ////     //////////   ///                      
--        ////            ////          Requires    : Authentication
--          ////     /////              Output      : 
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/admin/v1/webuser/{webUserID}/enable

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local webUserID,    err, msg = checkParameter( webUserID, 1 )       if err ~= 0 then return error( err, string.format( msg, "email" ) )  end
webUserID = tonumber( webUserID )

command = "SELECT first_name as FirstName, last_name as LastName, company as Company, email as Email FROM WebUser WHERE id = ? AND enabled = 0;"
local webuser = executeSQL( "auth", command, { webUserID } )
if not webuser                                                      then return error( 504, "Gateway Timeout" )              end
if webuser.ErrorNumber       ~= 0                                   then return error( 502, "Bad Gateway" )                  end
if webuser.NumberOfColumns   ~= 4                                   then return error( 502, "Bad Gateway" )                  end 
if webuser.NumberOfRows      == 0                                   then return error( 406, "Not Acceptable")                end
webuser = webuser.Rows[1]

local password = createPassword(12, 2, 3, false, true)
if not password or password == ""                                   then return error( 500, "Internal Server Error" )        end
webuser.Password = password

command = "INSERT INTO Company (name) VALUES (?) RETURNING id"
companyID = executeSQL( "auth", command, { webuser.Company } )
if not companyID                                                    then return error( 504, "Gateway Timeout" )              end
if companyID.ErrorNumber       ~= 0                                 then return error( 502, "Bad Gateway" )                  end
if companyID.NumberOfColumns   ~= 1                                 then return error( 502, "Bad Gateway" )                  end 
if companyID.NumberOfRows      ~= 1                                 then return error( 502, "Bad Gateway" )                  end
companyID = companyID.Rows[ 1 ].id

command = "INSERT INTO User (company_id, first_name, last_name, email, password) VALUES (?, ?, ?, ?, ?) RETURNING id;"
result = executeSQL( "auth", command, { companyID, webuser.FirstName, webuser.LastName, webuser.Email, hash(password) } )
if not result                                                       then return error( 504, "Gateway Timeout" )              end
if result.ErrorNumber       == 19                                   then return error( 409, "Email already exists" )         end
if result.ErrorNumber       ~= 0                                    then return error( 403, "Could not create user" )        end
if result.ErrorMessage      ~= ""                                   then return error( 502, result.ErrorMessage )            end
if result.NumberOfRows      ~= 1                                    then return error( 502, "Bad Gateway" )                  end
userID = result.Rows[ 1 ].id

-- use the default sender
local from = "" 
local to = webuser.Email
local subject = "Welcome to SQLiteCloud"
local templateName = "welcome_web"
local language = "en"
result = mail(from, to, subject, templateName, language, webuser)  

if not result then
    command = "DELETE FROM Company WHERE id = ?; DELETE FROM User WHERE id = ?"
    result = executeSQL( "auth", command, { companyID, userID } )
    return error( 500, "Mail not sent" )        
end

command = "UPDATE WebUser SET enabled = 1 WHERE id = ?;"
executeSQL( "auth", command, { webUserID } )

error( 200, "OK" )