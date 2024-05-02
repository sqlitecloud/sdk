<?php

declare(strict_types=1);

include_once 'sqcloud.php';

use PHPUnit\Framework\TestCase;

class SQLiteCloudTest extends TestCase
{
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

    public function testConnectWithStringWithCredentials(): void
    {
        $sqlite = new SQLiteCloud();

        $result = $sqlite->connectWithString(getenv('SQLITE_CONNECTION_STRING_CREDENTIALS'));

        $this->assertTrue($result, "Please, verify the connection parameters and the node is running. Message: {$sqlite->errmsg}");

        $this->assertSame(0, $sqlite->errcode);
        $this->assertEmpty($sqlite->errmsg);

        $sqlite->disconnect();
    }

    public function testConnectWithStringWithApiKey(): void
    {
        $sqlite = new SQLiteCloud();

        $result = $sqlite->connectWithString(getenv('SQLITE_CONNECTION_STRING_API_KEY'));

        $this->assertTrue($result, "Please, verify the connection parameters and the node is running. Message: {$sqlite->errmsg}");

        $this->assertSame(0, $sqlite->errcode);
        $this->assertEmpty($sqlite->errmsg);

        $sqlite->disconnect();
    }

    public function testRowsetData()
    {
        $sqlite = new SQLiteCloud();

        $result = $sqlite->connectWithString(getenv('SQLITE_CONNECTION_STRING_API_KEY'));

        $this->assertTrue($result);

        /** @var SQLiteCloudRowset */
        $rowset = $sqlite->execute('SELECT AlbumId FROM albums LIMIT 2');

        $this->assertSame(2, $rowset->nrows);
        $this->assertSame(1, $rowset->ncols);
        $this->assertSame(2, $rowset->version);
    }

    public function testGetValue()
    {
        $sqlite = new SQLiteCloud();

        $result = $sqlite->connectWithString(getenv('SQLITE_CONNECTION_STRING_API_KEY'));

        $this->assertTrue($result);

        /** @var SQLiteCloudRowset */
        $rowset = $sqlite->execute('SELECT * FROM albums');

        $this->assertSame('1', $rowset->value(0, 0));
        $this->assertSame('For Those About To Rock We Salute You', $rowset->value(0, 1));
        $this->assertSame('2', $rowset->value(1, 0));
    }

    public function testInvalidRowNumberForValue()
    {
        $sqlite = new SQLiteCloud();

        $result = $sqlite->connectWithString(getenv('SQLITE_CONNECTION_STRING_API_KEY'));

        $this->assertTrue($result);

        /** @var SQLiteCloudRowset */
        $rowset = $sqlite->execute('SELECT * FROM albums LIMIT 1');

        $this->assertNull($rowset->value(1, 1));
    }

    public function testColumnName()
    {
        $sqlite = new SQLiteCloud();

        $result = $sqlite->connectWithString(getenv('SQLITE_CONNECTION_STRING_API_KEY'));

        $this->assertTrue($result);

        /** @var SQLiteCloudRowset */
        $rowset = $sqlite->execute('SELECT * FROM albums');

        $this->assertSame('AlbumId', $rowset->name(0));
        $this->assertSame('Title', $rowset->name(1));
    }

    public function testInvalidRowNumberForColumnName()
    {
        $sqlite = new SQLiteCloud();

        $result = $sqlite->connectWithString(getenv('SQLITE_CONNECTION_STRING_API_KEY'));

        $this->assertTrue($result);

        /** @var SQLiteCloudRowset */
        $rowset = $sqlite->execute('SELECT AlbumId FROM albums');

        $this->assertNull($rowset->name(1));
    }
}
