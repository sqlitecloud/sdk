<?php include_once('common.php'); ?>
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
    		<h3>Dashboard Admin</h3>
  		</div>
	
		<div class="row">
			<div class="col-md-4 offset-md-4">
			
			<form id="login-form" method="post" action="#" role="form">
				
				<div id="message" class="alert alert-danger" role="alert" style="display: none;"></div>
				
				<div class="row" style="margin-bottom: 0.75em;">
					<label for="hostname" style="font-weight: 600;">Hostname</label>
					  <input type="text" class="form-control" name="hostname" id="hostname" required>
				</div>
				
				<div class="row" style="margin-bottom: 0.75em;">
					<label for="port" style="font-weight: 600;">Port</label>
					  <input type="text" class="form-control" name="port" id="port" value ="8860" required>
				</div>
				
				<div class="row" style="margin-bottom: 0.75em;">
					<label for="username" style="font-weight: 600;">Username</label>
          			<input type="text" class="form-control" name="username" id="username" value="admin" required>
        		</div>
					
				<div class="row" style="margin-bottom: 0.75em;">
					<label for="password" style="font-weight: 600;">Password</label>
          			<input type="password" class="form-control" name="password" id="password" required>
        		</div>
				
				<hr>
				
				<div class="row">
        			<button class="btn btn-primary btn-lg btn-block" type="submit">Login</button>
        		</div>
				
			</form>
			</div>
		</div>
		
    <script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.6/umd/popper.min.js" integrity="sha384-wHAiFfRlMFy6i5SRaxvfOCifBUQy1xHdJ/yoi7FRNXMRBu5WHdZYu1hA6ZOblgut" crossorigin="anonymous"></script>
    <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.2.1/js/bootstrap.min.js" integrity="sha384-B0UglyR+jN6CkvvICOB2joaf5I4l3gm9GU6Hc1og6Ls7i6U/mkkaduKaBhlAXv9k" crossorigin="anonymous"></script>
	<script src="login.js"></script>
  
  </body>
</html>
