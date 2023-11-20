//
// GET /api/hello - example API route using sqlitecloud-js
//

import { NextRequest, NextResponse } from "next/server";
import { Database } from "sqlitecloud-js";

// connection string is defined in .env file
const CHINOOK_URL = process.env.CHINOOK_URL as string;
console.assert(CHINOOK_URL, "Please configure a .env file with CHINOOK_URL pointing to the chinook database, see .env.example");

// route for /api/hello
export async function GET(request: NextRequest) {
  // connect to database using connection string provided in https://dashboard.sqlitecloud.io/
  const database = new Database(CHINOOK_URL);

  // retrieve rows from chinook database using a plain SQL query
  const tracks = await database.sql`SELECT * FROM tracks LIMIT 20`;

  // return as json response
  return NextResponse.json<{data: any}>({ data: tracks });
}
