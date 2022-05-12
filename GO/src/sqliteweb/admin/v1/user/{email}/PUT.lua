--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/04/11
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

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local email,    err, msg = checkParameter( email, 3 )                   if err ~= 0 then return error( err, string.format( msg, "email" ) )  end

local newEmail, err, msg = getBodyValue( "email", 0 )                   if err ~= 0 then return error( err, msg )                            end
local firstName,err, msg = getBodyValue( "first_name", 2 )              if err ~= 0 then return error( err, msg )                            end
local lastName, err, msg = getBodyValue( "last_name", 2 )               if err ~= 0 then return error( err, msg )                            end
-- local company,  err, msg = getBodyValue( "company", 0 )                 if err ~= 0 then return error( err, msg )                            end
local password, err, msg = getBodyValue( "password", 5 )                if err ~= 0 then return error( err, msg )                            end
local enabled,  err, msg = getBodyValue( "enabled", 1 )                 if err ~= 0 then return error( err, msg )                            end

if not newEmail                   then newEmail = email end
if string.len( newEmail ) == 0    then newEmail = email end

if bool( enabled )                then enabled = 1 else enabled = 0 end

query  = string.format( "UPDATE User SET last_name = '%s', first_name = '%s', email = '%s', password = '%s', enabled = %d WHERE email = '%s'; SELECT changes() AS success;", enquoteSQL( lastName ), enquoteSQL( firstName ), enquoteSQL( newEmail ), enquoteSQL( password ), enabled, enquoteSQL( email ) )
result = executeSQL( "auth", query )
if not result                                                                       then return error( 504, "Gateway Timeout" )              end
if result.ErrorMessage      ~= ""                                                   then return error( 502, result.ErrorMessage )            end
if result.ErrorNumber       ~= 0                                                    then return error( 403, "Could not create user" )        end
if result.NumberOfColumns   ~= 1                                                    then return error( 502, "Bad Gateway" )                  end
if result.NumberOfRows      ~= 1                                                    then return error( 502, "Bad Gateway" )                  end
if result.Rows[ 1 ].success == 0                                                    then return error( 404, "User not found" )               end
if result.Rows[ 1 ].success  > 1                                                    then return error( 500, "Internal Server Error" )        end                                       

error( 200, "OK" )