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
}
