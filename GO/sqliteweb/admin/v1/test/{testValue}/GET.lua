require "sqlitecloud"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

if testValue == "0" then return error(200, "OK") end
error( 418, "I'm a teapot" )
