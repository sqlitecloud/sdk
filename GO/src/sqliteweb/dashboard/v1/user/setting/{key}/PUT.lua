-- Change value for setting key for logged in user
-- https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/user/setting/{key}

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

do return error( 200, "OK" ) end