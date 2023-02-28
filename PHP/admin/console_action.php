<?php
	include_once('auth.php');
	include_once('common.php');
	
	$data = json_decode(file_get_contents('php://input'), true);
	$database = $data["database"];
	$sql = $data["sql"];
	
	$r = NULL;
	$rc = exec_sql($database, $sql);
	
	if ($rc === false) {
		$r = array('result' => 0, 'msg' => exec_lasterror(true));
	} else if ($rc instanceof SQLiteCloudRowset) {
		$r = array('result' => 2, 'msg' => render_console_table($rc));
	} else if ($rc === true) {
		$r = array('result' => 1, 'msg' => 'Query succesfully executed.');
	} else if ($rc === null) {
		$r = array('result' => 1, 'msg' => 'NULL');
	} else {
		$r = array('result' => 1, 'msg' => $rc);
	}
	
	echo json_encode($r);
?>