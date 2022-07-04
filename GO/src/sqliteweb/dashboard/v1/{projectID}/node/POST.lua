--
--                    ////              SQLite Cloud
--        ////////////  ///
--      ///             ///  ///        Product     : SQLite Cloud Web Server
--     ///             ///  ///         Version     : 1.0.0
--     //             ///   ///  ///    Date        : 2022/03/26
--    ///             ///   ///  ///    Author      : Andreas Pfeil
--   ///             ///   ///  ///
--   ///     //////////   ///  ///      Description : Add a new node
--   ////                ///  ///                     
--     ////     //////////   ///        Requires    : Authentication
--        ////            ////          Output      : Structure with user settings
--          ////     /////              
--             ///                      Copyright   : 2022 by SQLite Cloud Inc.
--
-- -----------------------------------------------------------------------TAB=2

-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

local userID,    err, msg = checkUserID( userid )                        if err ~= 0 then return error( err, msg )                          end
local nodeID,    err, msg = checkNodeID( nodeID )                        if err ~= 0 then return error( err, msg )                          end
local projectID, err, msg = checkProjectID( projectID )                  if err ~= 0 then return error( err, msg )                          end

local name,      err, msg = getBodyValue( "name", 1 )                    if err ~= 0 then return error( err, msg )                          end -- Dev1 Server
local type,      err, msg = getBodyValue( "hardware", 1 )                if err ~= 0 then return error( err, msg )                          end -- 1VCPU/1GB/25GB
local provider,  err, msg = getBodyValue( "region", 1 )                  if err ~= 0 then return error( err, msg )                          end -- NYC3/US
local image,     err, msg = getBodyValue( "type", 1 )                    if err ~= 0 then return error( err, msg )                          end -- worker
local region,    err, msg = getBodyValue( "counter", 0 )                 if err ~= 0 then return error( err, msg )                          end -- 1

do return error( 200, "OK" ) end