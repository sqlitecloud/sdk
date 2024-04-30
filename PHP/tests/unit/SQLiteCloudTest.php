<?php

declare(strict_types=1);

include_once 'sqcloud.php';

use PHPUnit\Framework\MockObject\MockObject;
use PHPUnit\Framework\TestCase;


class SQLiteCloudTest extends TestCase {

    public function testConnectWithStringWithPort(): void
    {
         /** @var MockObject|SQLiteCloud */
         $sqlite = $this->getMockBuilder(SQLiteCloud::class)
         ->onlyMethods(['connect'])
         ->getMock();
 
         $sqlite->expects($this->once())
             ->method('connect')
             ->with('disney.sqlite.cloud', 9972)
             ->willReturn(true);
 
         $connectionString = 'sqlitecloud://disney.sqlite.cloud:9972';
 
         $sqlite->connectWithString($connectionString);
    }

    public function testConnectWithStringWithBothApiKeyAndCredentials(): void
    {
        /** @var MockObject|SQLiteCloud */
        $sqlite = $this->getMockBuilder(SQLiteCloud::class)
        ->onlyMethods(['connect'])
        ->getMock();

        $sqlite->expects($this->once())
            ->method('connect')
            ->willReturn(true);

        $connectionString = 'sqlitecloud://pippo:pluto@disney.sqlite.cloud:8860?apikey=abc12345';

        $sqlite->connectWithString($connectionString);

        $this->assertEmpty($sqlite->username);
        $this->assertEmpty($sqlite->password);
        $this->assertSame('abc12345', $sqlite->apikey);
    }

    public function testConnectWithStringWithOptions(): void
    {
         /** @var MockObject|SQLiteCloud */
         $sqlite = $this->getMockBuilder(SQLiteCloud::class)
         ->onlyMethods(['connect'])
         ->getMock();
 
         $sqlite->expects($this->once())
             ->method('connect')
             ->with('disney.sqlite.cloud')
             ->willReturn(true);
 
         $connectionString = 'sqlitecloud://disney.sqlite.cloud/mydb?apikey=abc12345&insecure=true&timeout=100';
 
         $sqlite->connectWithString($connectionString);

         $this->assertSame('mydb', $sqlite->database);
         $this->assertSame('abc12345', $sqlite->apikey);
         $this->assertSame(true, $sqlite->insecure);
         $this->assertSame(100, $sqlite->timeout);
    }

    public function testConnectWithStringWithoutOptionals(): void
    {
         /** @var MockObject|SQLiteCloud */
         $sqlite = $this->getMockBuilder(SQLiteCloud::class)
         ->onlyMethods(['connect'])
         ->getMock();
 
         $sqlite->expects($this->once())
             ->method('connect')
             ->with('disney.sqlite.cloud')
             ->willReturn(true);
 
         $connectionString = 'sqlitecloud://disney.sqlite.cloud';
 
         $sqlite->connectWithString($connectionString);

         $this->assertEmpty($sqlite->username);
         $this->assertEmpty($sqlite->password);
         $this->assertEmpty($sqlite->database);
    }
}