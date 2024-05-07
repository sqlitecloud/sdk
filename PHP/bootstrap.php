<?php

require __DIR__ . '/vendor/autoload.php';

use Dotenv\Dotenv;

$dotenv = Dotenv::createImmutable(__DIR__);

foreach ($dotenv->load() as $key => $value) {
    putenv("$key=$value");
}
