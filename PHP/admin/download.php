<?php
	include_once('common.php');
	$dbName = $_GET['db'];
	$maxRead = 1 * 1024 * 1024; // 1MB
	$dbSize = 0;
	$totRead = 0;
	
	// send DOWNLOAD DATABASE command
	$r = exec_downloaddatabase($dbName);
	$dbSize = $r[0];

	// These headers will force download on browser,
	// and set the custom file name for the download, respectively.
	header('Content-Type: application/octet-stream');
	header('Content-Disposition: attachment; filename="' . $dbName . '"');
	
	while ($totRead < $dbSize) {
		$data = exec_downloadstep();
		echo $data;
		$totRead += strlen($data);
		ob_flush();
	}
	exit;
	
?>