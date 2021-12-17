json = require "api.json"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

request = json.decode( body )
print( request )

result = queryNode( "SELECT * FROM Dummy" )

Response = {
  Request = request,
  Parameter = {
    First  = args[ 1 ],
    Second = args[ 2 ]
  },

  ResponseID = request[ 'RequestID' ],
  Status = 0,
  Message = "OK",

  QueryResult = result
}

Write( json.encode( Response ) )
SetStatus( 200 )