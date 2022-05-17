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


-- https://localhost:8443/admin/v1/users

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )


-- users = executeSQL( "auth", string.format( "SELECT COUNT( id ) AS granted FROM User JOIN Company ON User.company_id = Company.id JOIN Project ON Company.id = Project.company_id WHERE User.enabled = 1 AND Company.enabled = 1 AND User.id= %d AND uuid = '%s';", userID, enquoteSQL( projectID ) ) )
users = executeSQL( "auth", string.format( "SELECT User.id, email, first_name, last_name, Company.name as company FROM User JOIN Company ON User.company_id = Company.id WHERE User.enabled = 1 AND Company.enabled = 1;" ) )

if not users                   then return error( 504, "Gateway Timeout" )     end
if users.ErrorNumber       ~= 0  then return error( 502, "Bad Gateway" )       end
if users.NumberOfColumns   ~= 5  then return error( 502, "Bad Gateway" )       end 

fUsers = nil
if users.NumberOfRows > 0 then
  fUsers = filter( users.Rows, {
    id            = "id",
    email         = "email",
    first_name    = "first_name",
    last_name     = "last_name",
    company       = "company",
    
    -- creation_date = "created",
  } )
end

Response = {
  status = 200,
  message = "OK",

  users = fUsers
}

SetStatus( 200 )
Write( jsonEncode( Response ) )