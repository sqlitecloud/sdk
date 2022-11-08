--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Get a JSON with all available hardware codes
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : list of strings (hardware codes)
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                     end

Response = {
  status           = 200,                                     -- status code: 200 = no error, error otherwise
  message          = "OK",                                    -- "OK" or error message
  value            = nil                                      -- Array with project objects
}

if userID == 0 then                                           -- get list of projects in ini file
    return error( 501, "Not yet implemented in on-premise version" )

else
  local uID, companyID, err, msg = verifyUserID( userID )                  if err ~= 0 then return error( err, msg )                                end

  Response.value = {"1VCPU/1GB/25GB", "1VCPU/2GB/50GB", "2VCPU/2GB/60GB", "2VCPU/16GB/300GB"}
end

SetStatus( 200 )
Write( jsonEncode( Response ) )