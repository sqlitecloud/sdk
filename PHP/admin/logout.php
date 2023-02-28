<?php
  session_start();
  setcookie (session_id(), "start", time() - 3600);
  setcookie (session_id(), "port", time() - 3600);
  setcookie (session_id(), "hostname", time() - 3600);
  setcookie (session_id(), "username", time() - 3600);
  setcookie (session_id(), "password", time() - 3600);
  session_destroy();
  session_write_close();
?>

<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.2.1/css/bootstrap.min.css" integrity="sha384-GJzZqFGwb1QTTN6wy59ffF1BuGJpLSa9DkKMp0DgiMDm4iYMj70gZWKYbI706tWS" crossorigin="anonymous">
    <title>SQLite Cloud Admin</title>
  </head>
  
  <body class="bg-light">
	<div class="container">
	
		<div class="py-5 text-center">
    		<img class="d-block mx-auto mb-4" src="assets/images/logo.png" alt="" width="268" height="54">
    		<h3>Goodbye!</h3>
  		</div>
		
    <script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.6/umd/popper.min.js" integrity="sha384-wHAiFfRlMFy6i5SRaxvfOCifBUQy1xHdJ/yoi7FRNXMRBu5WHdZYu1hA6ZOblgut" crossorigin="anonymous"></script>
    <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.2.1/js/bootstrap.min.js" integrity="sha384-B0UglyR+jN6CkvvICOB2joaf5I4l3gm9GU6Hc1og6Ls7i6U/mkkaduKaBhlAXv9k" crossorigin="anonymous"></script>
	
  </body>
</html>
