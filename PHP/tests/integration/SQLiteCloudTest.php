<?php

declare(strict_types=1);

include_once 'src/sqcloud.php';

use PHPUnit\Framework\TestCase;

class SQLiteCloudTest extends TestCase
{
    /** Will warn if a query or other basic operation is slower than this */
    const WARN_SPEED_MS = 500;
    /** Will except queries to be quicker than this */
    const EXPECT_SPEED_MS = 6 * 1000;

    private $sqlite;

    protected function tearDown()
    {
        if ($this->sqlite) {
            $this->sqlite->disconnect();
        }
    }

    /**
     * @return SQLiteCloud
     */
    private function getSQLiteConnection()
    {
        $this->sqlite = new SQLiteCloud();
        $this->sqlite->database = getenv('SQLITE_DB');
        $this->sqlite->apikey = getenv('SQLITE_API_KEY');

        $result = $this->sqlite->connect(getenv('SQLITE_HOST'));
        $this->assertTrue($result);

        return $this->sqlite;
    }

    public function testConnectWithoutCredentialsAndApikey()
    {
        $sqlite = new SQLiteCloud();
        $sqlite->username = '';
        $sqlite->password = '';
        $sqlite->apikey = '';
        $sqlite->database = getenv('SQLITE_DB');

        $result = $sqlite->connect(getenv('SQLITE_HOST'), getenv('SQLITE_PORT'));

        $this->assertFalse($result);
    }

    public function testConnect(): void
    {
        $sqlite = new SQLiteCloud();
        $sqlite->username = getenv('SQLITE_USER');
        $sqlite->password = getenv('SQLITE_PASSWORD');
        $sqlite->database = getenv('SQLITE_DB');

        $result = $sqlite->connect(getenv('SQLITE_HOST'), getenv('SQLITE_PORT'));

        $this->assertTrue($result, "Please, verify the connection parameters and the node is running. Message: {$sqlite->errmsg}");

        $this->assertSame(0, $sqlite->errcode);
        $this->assertEmpty($sqlite->errmsg);

        $sqlite->disconnect();
    }

    public function testConnectWithStringWithoutCredentialsAndApikey()
    {
        $sqlite = new SQLiteCloud();
        $sqlite->username = '';
        $sqlite->password = '';
        $sqlite->apikey = '';

        $result = $sqlite->connectWithString('sqlitecloud://' . getenv('SQLITE_HOST') . '/' . getenv('SQLITE_DB'));

        $this->assertFalse($result);

        $sqlite->disconnect();
    }

    public function testConnectWithStringWithCredentials(): void
    {
        $sqlite = new SQLiteCloud();

        $result = $sqlite->connectWithString(getenv('SQLITE_CONNECTION_STRING'));

        $this->assertTrue($result, "Please, verify the connection parameters and the node is running. Message: {$sqlite->errmsg}");

        $this->assertSame(0, $sqlite->errcode);
        $this->assertEmpty($sqlite->errmsg);

        $sqlite->disconnect();
    }

    public function testConnectWithStringWithApiKey(): void
    {
        $sqlite = new SQLiteCloud();

        $result = $sqlite->connectWithString(getenv('SQLITE_CONNECTION_STRING') . '?apikey=' . getenv('SQLITE_API_KEY'));

        $this->assertTrue($result, "Please, verify the connection parameters and the node is running. Message: {$sqlite->errmsg}");

        $this->assertSame(0, $sqlite->errcode);
        $this->assertEmpty($sqlite->errmsg);

        $sqlite->disconnect();
    }

    public function testRowsetData()
    {
        $sqlite = $this->getSQLiteConnection();

        /** @var SQLiteCloudRowset */
        $rowset = $sqlite->execute('SELECT AlbumId FROM albums LIMIT 2');

        $this->assertSame(2, $rowset->nrows);
        $this->assertSame(1, $rowset->ncols);
        $this->assertSame(2, $rowset->version);
    }

    public function testGetValue()
    {
        $sqlite = $this->getSQLiteConnection();

        /** @var SQLiteCloudRowset */
        $rowset = $sqlite->execute('SELECT * FROM albums');

        $this->assertSame('1', $rowset->value(0, 0));
        $this->assertSame('For Those About To Rock We Salute You', $rowset->value(0, 1));
        $this->assertSame('2', $rowset->value(1, 0));
    }

    public function testSelectUTF8ValueAndColumnName()
    {
        $sqlite = $this->getSQLiteConnection();

        /** @var SQLiteCloudRowset */
        $rowset = $sqlite->execute("SELECT 'Minha História'");

        $this->assertSame('Minha História', $rowset->value(0, 0));
        $this->assertSame("'Minha História'", $rowset->name(0));
    }

    public function testColumnNotFound()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute("SELECT not_a_column FROM albums");

        $this->assertFalse($rowset);
        $this->assertSame(1, $sqlite->errcode);
        $this->assertSame("no such column: not_a_column", $sqlite->errmsg);
    }

    public function testInvalidRowNumberForValue()
    {
        $sqlite = $this->getSQLiteConnection();

        /** @var SQLiteCloudRowset */
        $rowset = $sqlite->execute("SELECT 'one row'");

        $this->assertNull($rowset->value(1, 1));
    }

    public function testColumnName()
    {
        $sqlite = $this->getSQLiteConnection();

        /** @var SQLiteCloudRowset */
        $rowset = $sqlite->execute('SELECT * FROM albums');

        $this->assertSame('AlbumId', $rowset->name(0));
        $this->assertSame('Title', $rowset->name(1));
    }

    public function testInvalidRowNumberForColumnName()
    {
        $sqlite = $this->getSQLiteConnection();

        /** @var SQLiteCloudRowset */
        $rowset = $sqlite->execute('SELECT AlbumId FROM albums');

        $this->assertNull($rowset->name(1));
    }

    /**
     * @large
     */
    public function testLongString()
    {
        $sqlite = $this->getSQLiteConnection();

        $size = 1024 * 1024;
        $value = 'LOOOONG';
        while (strlen($value) < $size) {
            $value .= 'a';
        }
        $rowset = $sqlite->execute("SELECT '{$value}' 'VALUE'");

        $this->assertEmpty($sqlite->errmsg);
        $this->assertNotFalse($rowset);
        $this->assertSame(0, $sqlite->errcode);

        $this->assertSame(1, $rowset->nrows);
        $this->assertSame(1, $rowset->ncols);
        $this->assertSame('VALUE', $rowset->name(0));
        $this->assertSame($value, $rowset->value(0, 0));
    }

    public function testInteger()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('TEST INTEGER');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertSame(123456, $rowset);
    }

    public function testFloat()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('TEST FLOAT');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertSame(3.1415926, $rowset);
    }

    public function testString()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('TEST STRING');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertSame('Hello World, this is a test string.', $rowset);
    }

    public function testZeroString()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('TEST ZERO_STRING');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertSame('Hello World, this is a zero-terminated test string.', $rowset);
    }

    public function testEmptyString()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('TEST STRING0');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertSame('', $rowset);
    }

    public function testCommand()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('TEST COMMAND');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertSame('PONG', $rowset);
    }

    public function testJson()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('TEST JSON');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertJson($rowset);
        $this->assertEquals(
            [
                'msg-from' => ['class' => 'soldier', 'name' => 'Wixilav'],
                'msg-to' => ['class' => 'supreme-commander', 'name' => '[Redacted]'],
                'msg-type' => ['0xdeadbeef', 'irc log'],
                'msg-log' => [
                    'soldier: Boss there is a slight problem with the piece offering to humans',
                    'supreme-commander: Explain yourself soldier!',
                    "soldier: Well they don't seem to move anymore...",
                    'supreme-commander: Oh snap, I came here to see them twerk!'
                ]
            ],
            json_decode($rowset, true)
        );
    }

    public function testBlob()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('TEST BLOB');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertSame(1000, strlen($rowset));
    }

    public function testBlob0()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('TEST BLOB0');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertSame(0, strlen($rowset));
    }

    public function testError()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('TEST ERROR');
        $this->assertSame(66666, $sqlite->errcode);
        $this->assertSame('This is a test error message with a devil error code.', $sqlite->errmsg);
    }

    public function testExtError()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('TEST EXTERROR');
        $this->assertSame(66666, $sqlite->errcode);
        $this->assertSame(333, $sqlite->xerrcode);
        $this->assertSame('This is a test error message with an extcode and a devil error code.', $sqlite->errmsg);
    }

    public function testArray()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('TEST ARRAY');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertTrue(is_array($rowset));
        $this->assertCount(5, $rowset);
        $this->assertSame('Hello World', $rowset[0]);
        $this->assertSame('123456', $rowset[1]);
        $this->assertSame('3.1415', $rowset[2]);
        $this->assertNull($rowset[3]);
    }

    public function testRowset()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('TEST ROWSET');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertGreaterThanOrEqual(30, $rowset->nrows);
        $this->assertSame(2, $rowset->ncols);
        $this->assertTrue(in_array($rowset->version, [1, 2]));
        $this->assertSame('key', $rowset->name(0));
        $this->assertSame('value', $rowset->name(1));
    }

    public function testMaxRowsOption()
    {
        $sqlite = new SQLiteCloud();
        $sqlite->database = getenv('SQLITE_DB');
        $sqlite->apikey = getenv('SQLITE_API_KEY');
        $sqlite->maxrows = 1;
        $result = $sqlite->connect(getenv('SQLITE_HOST'));

        $this->assertTrue($result);

        $rowset = $sqlite->execute('SELECT * FROM albums');
        $this->assertNotFalse($rowset);
        $this->assertGreaterThan(100, $rowset->nrows);

        $sqlite->disconnect();
    }

    public function testMaxRowsetOptionToFailWhenRowsetIsBigger()
    {
        $sqlite = new SQLiteCloud();
        $sqlite->database = getenv('SQLITE_DB');
        $sqlite->apikey = getenv('SQLITE_API_KEY');
        $sqlite->maxrowset = 1024;
        $result = $sqlite->connect(getenv('SQLITE_HOST'));

        $this->assertTrue($result);

        $rowset = $sqlite->execute('SELECT * FROM albums');
        $this->assertFalse($rowset);
        $this->assertSame('RowSet too big to be sent (limit set to 1024 bytes).', $sqlite->errmsg);

        $sqlite->disconnect();
    }

    public function testMaxRowsetOptionToSuccedWhenRowsetIsLighter()
    {
        $sqlite = new SQLiteCloud();
        $sqlite->database = getenv('SQLITE_DB');
        $sqlite->apikey = getenv('SQLITE_API_KEY');
        $sqlite->maxrowset = 1024;
        $result = $sqlite->connect(getenv('SQLITE_HOST'));

        $this->assertTrue($result);

        $rowset = $sqlite->execute("SELECT 'hello world'");
        $this->assertNotFalse($rowset);
        $this->assertSame(1, $rowset->nrows);

        $sqlite->disconnect();
    }

    public function testChunckedRowset()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('TEST ROWSET_CHUNK');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertSame(147, $rowset->nrows);
        $this->assertSame(1, $rowset->ncols);
        $this->assertSame('key', $rowset->name(0));
    }

    public function testChunckedRowsetTwice()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('TEST ROWSET_CHUNK');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertSame(147, $rowset->nrows);
        $this->assertSame(1, $rowset->ncols);
        $this->assertSame('key', $rowset->name(0));

        $rowset = $sqlite->execute('TEST ROWSET_CHUNK');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertSame(147, $rowset->nrows);
        $this->assertSame(1, $rowset->ncols);
        $this->assertSame('key', $rowset->name(0));

        $rowset = $sqlite->execute('SELECT 1');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertSame(1, $rowset->nrows);
    }

    public function testSerializedOperations()
    {
        $numQueries = 20;

        $sqlite = $this->getSQLiteConnection();

        for ($i = 0; $i < $numQueries; $i++) {
            $rowset = $sqlite->execute("select {$i} as 'count', 'hello' as 'string'");

            $this->assertEmpty($sqlite->errmsg);
            $this->assertSame(1, $rowset->nrows);
            $this->assertSame(2, $rowset->ncols);
            $this->assertSame('count', $rowset->name(0));
            $this->assertSame('string', $rowset->name(1));
            $this->assertSame("{$i}", $rowset->value(0, 0));
            $this->assertTrue(in_array($rowset->version, [1, 2]));
        }
    }

    public function testQueryTimeout()
    {
        $sqlite = new SQLiteCloud();
        $sqlite->apikey = getenv('SQLITE_API_KEY');
        $sqlite->database = getenv('SQLITE_DB');
        $sqlite->timeout = 1; // 1 sec

        $result = $sqlite->connect(getenv('SQLITE_HOST'));
        $this->assertTrue($result);

        // this operation should take more then 1s
        $rowset = $sqlite->execute(
            // just a long running query
            "WITH RECURSIVE r(i) AS (
                VALUES(0)
                UNION ALL
                SELECT i FROM r
                LIMIT 10000000
            )
            SELECT i FROM r WHERE i = 1;"
        );
        $this->assertFalse($rowset);

        $sqlite->disconnect();
    }

    public function testXXLQuery()
    {
        $xxlQuery = 300000;
        $longSql = '';

        $sqlite  = $this->getSQLiteConnection();

        while (strlen($longSql) < $xxlQuery) {
            for ($i = 0; $i < 5000; $i++) {
                $longSql .= 'SELECT ' . strlen($longSql) . "'HowLargeIsTooMuch'; ";
            }

            $rowset = $sqlite->execute($longSql);
            $this->assertSame(1, $rowset->nrows);
            $this->assertSame(1, $rowset->ncols);
            $this->assertGreaterThanOrEqual(strlen($longSql) - 50, $rowset->value(0, 0));
        }
    }

    /**
     * @large
     */
    public function testSingleXXLQuery()
    {
        $xxlQuery = 200000;
        $longSql = '';

        $sqlite = $this->getSQLiteConnection();

        while (strlen($longSql) < $xxlQuery) {
            $longSql .= strlen($longSql) . "_";
        }
        $selectedValue = "start_{$longSql}end";
        $longSql = "SELECT '{$selectedValue}'";

        $rowset = $sqlite->execute($longSql);

        $this->assertSame(1, $rowset->nrows);
        $this->assertSame(1, $rowset->ncols);
        $this->assertSame("'{$selectedValue}'", $rowset->name(0));
    }

    public function testMetadata()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('LIST METADATA');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertGreaterThanOrEqual(32, $rowset->nrows);
        $this->assertSame(8, $rowset->ncols);
    }

    public function testSelectResultsWithNoColumnName()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute("SELECT 42, 'hello'");
        $this->assertEmpty($sqlite->errmsg);
        $this->assertSame(1, $rowset->nrows);
        $this->assertSame(2, $rowset->ncols);
        $this->assertSame('42', $rowset->name(0));
        $this->assertSame("'hello'", $rowset->name(1));
        $this->assertSame('42', $rowset->value(0, 0));
        $this->assertSame('hello', $rowset->value(0, 1));
    }

    public function testSelectLongFormattedString()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute("USE DATABASE :memory:; SELECT '" . str_repeat('x', 1000) . "' AS DDD");
        $this->assertEmpty($sqlite->errmsg);
        $this->assertSame(1, $rowset->nrows);
        $this->assertSame(1, $rowset->ncols);
        $this->assertStringStartsWith('xxxxxxxx', $rowset->value(0, 0));
        $this->assertSame(1000, strlen($rowset->value(0, 0)));
    }

    public function testSelectDatabase()
    {
        $sqlite = new SQLiteCloud();
        $sqlite->apikey = getenv('SQLITE_API_KEY');
        $sqlite->database = '';

        $sqlite->connect(getenv('SQLITE_HOST'));

        $rowset = $sqlite->execute('USE DATABASE chinook.sqlite');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertEmpty($rowset->nrows);
        $this->assertEmpty($rowset->ncols);
    }

    public function testSelectTracksWithoutLimit()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('SELECT * FROM tracks');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertGreaterThanOrEqual(3000, $rowset->nrows);
        $this->assertSame(9, $rowset->ncols);
    }

    public function testSelectTracksWithLimit()
    {
        $sqlite = $this->getSQLiteConnection();
        $rowset = $sqlite->execute('SELECT * FROM tracks LIMIT 10;');
        $this->assertEmpty($sqlite->errmsg);
        $this->assertSame(10, $rowset->nrows);
        $this->assertSame(9, $rowset->ncols);
    }

    /**
     * @large
     */
    public function testStressTest20xStringSelectIndividual()
    {
        $numQueries = 20;
        $completed = 0;
        $startTime = microtime(true);

        $sqlite = $this->getSQLiteConnection();

        for ($i = 0; $i < $numQueries; $i++) {
            $rowset = $sqlite->execute('TEST STRING');
            $this->assertNull($sqlite->errmsg);
            $this->assertSame('Hello World, this is a test string.', $rowset);

            if (++$completed >= $numQueries) {
                $queryMs = round((microtime(true) - $startTime) * 1000 / $numQueries);
                if ($queryMs > self::WARN_SPEED_MS) {
                    $this->assertLessThan(self::EXPECT_SPEED_MS, $queryMs, "{$numQueries}x test string, {$queryMs}ms per query");
                }
            }
        }
    }

    /**
     * @large
     */
    public function testStressTest20xIndividualSelect()
    {
        $numQueries = 20;
        $completed = 0;
        $startTime = microtime(true);

        $sqlite = $this->getSQLiteConnection();

        for ($i = 0; $i < $numQueries; $i++) {
            $rowset = $sqlite->execute('SELECT * FROM albums ORDER BY RANDOM() LIMIT 4;');
            $this->assertNull($sqlite->errmsg);
            $this->assertSame(4, $rowset->nrows);
            $this->assertSame(3, $rowset->ncols);
            if (++$completed >= $numQueries) {
                $queryMs = round((microtime(true) - $startTime) * 1000 / $numQueries);
                if ($queryMs > self::WARN_SPEED_MS) {
                    $this->assertLessThan(self::EXPECT_SPEED_MS, $queryMs, "{$numQueries}x individual selects, {$queryMs}ms per query");
                }
            }
        }
    }

    /**
     * @long
     */
    public function testStressTest20xBatchedSelects()
    {
        $numQueries = 20;
        $completed = 0;
        $startTime = microtime(true);

        $sqlite = $this->getSQLiteConnection();

        for ($i = 0; $i < $numQueries; $i++) {
            $rowset = $sqlite->execute(
                'SELECT * FROM albums ORDER BY RANDOM() LIMIT 16; SELECT * FROM albums ORDER BY RANDOM() LIMIT 12; SELECT * FROM albums ORDER BY RANDOM() LIMIT 8; SELECT * FROM albums ORDER BY RANDOM() LIMIT 4;'
            );
            $this->assertNull($sqlite->errmsg);
            $this->assertSame(4, $rowset->nrows);
            $this->assertSame(3, $rowset->ncols);
            if (++$completed >= $numQueries) {
                $queryMs = round((microtime(true) - $startTime) * 1000 / $numQueries);
                if ($queryMs > self::WARN_SPEED_MS) {
                    $this->assertLessThan(self::EXPECT_SPEED_MS, $queryMs, "{$numQueries}x batched selects, {$queryMs}ms per query");
                }
            }
        }
    }

    public function testDownloadDatabase()
    {
        $sqlite = $this->getSQLiteConnection();

        $dbInfo = $sqlite->execute('DOWNLOAD DATABASE ' . getenv('SQLITE_DB'));
        $this->assertNotFalse($dbInfo);
        $dbSize = $dbInfo[0];

        $totRead = 0;
        $data = '';
        while ($totRead < $dbSize) {
            $data .= $sqlite->execute("DOWNLOAD STEP;");
            $totRead += strlen($data);
        }
        $tempFile = tempnam(sys_get_temp_dir(), 'chinook');
        file_put_contents($tempFile, $data);

        $db = new SQLite3($tempFile);
        $rowset = $db->query('SELECT * from albums');

        $this->assertNotFalse($rowset);
        $this->assertSame('AlbumId', $rowset->columnName(0));
        $this->assertSame('Title', $rowset->columnName(1));
    }

    public function testCompressionSingleColumn()
    {
        $sqlite = new SQLiteCloud();
        $sqlite->apikey = getenv('SQLITE_API_KEY');
        $sqlite->database = getenv('SQLITE_DB');
        $sqlite->compression = true;

        $result = $sqlite->connect(getenv('SQLITE_HOST'));
        $this->assertTrue($result);

        // min compression size for rowset set by default to 20400 bytes
        $blobSize = 20 * 1024;
        $rowset = $sqlite->execute("SELECT hex(randomblob({$blobSize})) AS 'someColumnName'");
        
        $this->assertEmpty($sqlite->errmsg);
        $this->assertSame(1, $rowset->nrows);
        $this->assertSame(1, $rowset->ncols);
        $this->assertSame("someColumnName", $rowset->name(0));
        $this->assertSame($blobSize * 2, strlen($rowset->value(0, 0)));

        $sqlite->disconnect();
    }

    public function testCompressionMultipleColumns()
    {
        $sqlite = new SQLiteCloud();
        $sqlite->apikey = getenv('SQLITE_API_KEY');
        $sqlite->database = getenv('SQLITE_DB');
        $sqlite->compression = true;

        $result = $sqlite->connect(getenv('SQLITE_HOST'));
        $this->assertTrue($result);

        // min compression size for rowset set by default to 20400 bytes
        $rowset = $sqlite->execute("SELECT * from albums inner join albums a2 on albums.AlbumId = a2.AlbumId");
        
        $this->assertEmpty($sqlite->errmsg);
        $this->assertGreaterThan(0, $rowset->nrows);
        $this->assertGreaterThan(0, $rowset->ncols);
        $this->assertSame("AlbumId", $rowset->name(0));

        $sqlite->disconnect();
    }
}
