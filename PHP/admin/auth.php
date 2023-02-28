<?php
	session_set_cookie_params(0);
	session_start();
	
	// check if there is no valid session or if session is expired
	$session_max_time = 60 * 600; // 60 minutes
	if ((!isset($_SESSION['username'])) || (isset($_SESSION['start']) && (time() - $_SESSION['start'] >= $session_max_time))) {
		session_unset();
		session_destroy();
    	header("Location: login.php");
		exit();
	}

?>