--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/04/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Send an email with a reset password link
--   ////                ///  ///                    
--     ////     //////////   ///                      
--        ////            ////          Requires    : Authentication
--          ////     /////              Output      : 
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/admin/v1/user/{email}/recover

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local email,    err, msg = checkParameter( email, 3 )                   if err ~= 0 then return error( err, string.format( msg, "email" ) )  end

-- check the email
command = "SELECT id FROM User WHERE email = ?;"
result = executeSQL( "auth", command, {email} )

if not result                                                                       then return error( 504, "Gateway Timeout" )              end
if result.ErrorMessage      ~= ""                                                   then return error( 502, result.ErrorMessage )            end
if result.ErrorNumber       ~= 0                                                    then return error( 403, "Could fetch user data" )        end
if result.NumberOfRows      ~= 1                                                    then return error( 404, "User not found" )               end
if result.NumberOfColumns   ~= 1                                                    then return error( 502, "Bad Gateway" )                  end                                

resetToken = rand64UrlSafeString()
command = "INSERT INTO PasswordResetToken (user_id, token) VALUES (?, ?) RETURNING rowid"
result = executeSQL( "auth", command, {result.Rows[ 1 ].id, hash(resetToken)} )

if not result                                                                       then return error( 504, "Gateway Timeout" )              end
if result.ErrorMessage      ~= ""                                                   then return error( 502, result.ErrorMessage )            end
if result.ErrorNumber       ~= 0                                                    then return error( 403, "Could fetch user data" )        end
if result.NumberOfRows      ~= 1                                                    then return error( 500, "Internal ServerError" )         end
if result.NumberOfColumns   ~= 1                                                    then return error( 502, "Bad Gateway" )                  end                                


template_data = {
  Token = resetToken,
}

-- with HTML template use the following tag otherwise the link would be overwritten 
-- so that any time a customer clicks a link, SendGrid can track those clicks: 
-- <a clicktracking="off" href='https://mysite/auth/'>My Site</a> 

-- fromstr string (use default sender if empty), tostr string, subject string, templateName string, language string, data map[string]string
mail("", email, "Reset your password", "recover.eml", "en", template_data)

error( 200, "OK" )