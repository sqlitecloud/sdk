<?php
	$data = json_decode(file_get_contents('php://input'), true);
	$hostname = $data["hostname"];
	$username = $data["username"];
	$password = $data["password"];
	$port = $data["port"];
	
	include_once('common.php');
	$rc = do_real_connect($hostname, $port, $username, $password);
	
	$r = NULL;
	if ($rc === true) {
		// successfully connected
		$r = array('result' => 1);
		$_SESSION['start'] = time();
		$_SESSION['port'] = $port;
		$_SESSION['hostname'] = $hostname;
		$_SESSION['username'] = $username;
		$_SESSION['password'] = $password;
	} else {
		// error
		$r = array('result' => 0, 'msg' => $rc);
	}
	
	echo json_encode($r);
?>