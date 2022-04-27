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

-- https://localhost:8443/admin/v1/user/{email}/recover

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local email,    err, msg = checkParameter( email, 3 )                   if err ~= 0 then return error( err, string.format( msg, "email" ) )  end

template_data = {
  From     = "<will.be@overwritten.com>",
  To       = email,
  Subject  = "Password recovery",
  Password = "secret",
}

query = string.format( "SELECT password FROM USER WHERE email = '%s';", enquoteSQL( email ) )
result = executeSQL( "auth", query )

if not result                                                                       then return error( 504, "Gateway Timeout" )              end
if result.ErrorMessage      ~= ""                                                   then return error( 502, result.ErrorMessage )            end
if result.ErrorNumber       ~= 0                                                    then return error( 403, "Could fetch user data" )        end
if result.NumberOfRows      ~= 1                                                    then return error( 404, "User not found" )               end
if result.NumberOfColumns   ~= 1                                                    then return error( 502, "Bad Gateway" )                  end                                

template_data.Password = result.Rows[ 1 ].password

mail( "recover.eml", "de", template_data )

error( 200, "OK" )