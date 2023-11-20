
var express = require("express");
var router = express.Router();
require('dotenv').config()

var sqlitecloud = require("sqlitecloud-js");

var CHINOOK_URL = process.env.CHINOOK_URL
console.assert(CHINOOK_URL, "CHINOOK_URL environment variable not set in .env");

/* GET chinook tracks as json */
router.get("/", async function (req, res, next) {
  var database = new sqlitecloud.Database(CHINOOK_URL);
  var tracks = await database.sql`SELECT * FROM tracks LIMIT 20`;
  res.send({ data: tracks });
});

module.exports = router;
