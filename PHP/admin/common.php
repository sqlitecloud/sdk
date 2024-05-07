<?php
	ini_set('display_startup_errors', 1);
	ini_set('display_errors', 1);
	error_reporting(-1);

	if (session_status() === PHP_SESSION_NONE) {
		session_set_cookie_params(0);
		session_start();
	}

	include_once('../src/sqcloud.php');
	
	function do_real_connect($hostname, $port, $username, $password) {
		global $sqlitecloud;
		$sqlitecloud = new SQLiteCloud();
		$sqlitecloud->username = $username;
		$sqlitecloud->password = $password;
		if (file_exists('assets/ca.pem')) {$sqlitecloud->tls_root_certificate = 'assets/ca.pem';}
		$sqlitecloud->compression = false;
		
		try {
			if ($sqlitecloud->connect($hostname, $port) == false) {
				$msg = $sqlitecloud->errmsg;
				$sqlitecloud = NULL;
				return $msg;
			}
		} catch (Exception $e) {
			return $e->getMessage();
		}
		
		return true;
	}
	
	function do_check_connect() {
		global $sqlitecloud;
		if ($sqlitecloud == NULL) {
			$hostname = $_SESSION['hostname'];
			$username = $_SESSION['username'];
			$password = $_SESSION['password'];
			$port = $_SESSION['port'];
			do_real_connect($hostname, $port, $username, $password);
		}
		
		return $sqlitecloud;
	}
	
	function do_disconnect() {
		global $sqlitecloud;
		if ($sqlitecloud != NULL) $sqlitecloud->disconnect();
	}
	
	function query_nodeinfo() {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->execute("LIST INFO;");
	}
	
	function query_listnodes() {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->execute("LIST NODES;");
	}
	
	function query_listcommands() {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->execute("LIST COMMANDS;");
	}
	
	function query_listdatabases($detailed = false) {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->execute(($detailed) ? "LIST DATABASES DETAILED;" : "LIST DATABASES;");
	}
	
	function query_listconnections() {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->execute("LIST CONNECTIONS;");
	}
	
	function query_listplugins() {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->execute("LIST PLUGINS;");
	}
	
	function query_listlogs($n = 50) {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->execute("LIST LOG LIMIT {$n};");
	}
	
	function query_listbackups() {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->execute("LIST BACKUPS;");
	}
	
	function query_listusers($withrules = false) {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->execute(($withrules) ? "LIST USERS WITH ROLES" : "LIST USERS;");
	}
	
	function query_liststats() {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->execute("LIST STATS;");
	}
	
	function query_listlatency() {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->execute("LIST LATENCY;");
	}
	
	function query_listkeys() {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->execute("LIST KEYS DETAILED;");
	}
	
	function exec_sql($dbname, $sql) {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->execute("SWITCH DATABASE {$dbname};{$sql}");
	}
	
	function exec_downloaddatabase($dbname) {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->execute("DOWNLOAD DATABASE {$dbname};");
	}
	
	function exec_downloadstep() {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->execute("DOWNLOAD STEP;");
	}
	
	function exec_uploaddatabase($dbname, $key) {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		$command = ($key) ? "UPLOAD DATABASE '${dbname}' KEY '{$key}';" : "UPLOAD DATABASE '{$dbname}';";
		return $sqlitecloud->execute($command);
	}
	
	function exec_uploadabort() {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->execute("DOWNLOAD ABORT;");
	}
	
	function exec_uploadblob($blob) {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->sendblob($blob);
	}
	
	function exec_lasterror($detailed = false) {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		if ($detailed) return $sqlitecloud->errmsg . ' (' . $sqlitecloud->errcode . ' - ' . $sqlitecloud->xerrcode . ')';
		return $sqlitecloud->errmsg;
	}
	
	function exec_lasterror_code() {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->errcode;
	}
	
	function exec_lasterror_xcode() {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		return $sqlitecloud->xerrcode;
	}
	
	// MARK: -
		
	function current_page() {
		return pathinfo($_SERVER['PHP_SELF'], PATHINFO_FILENAME);
	}
	
	function render_current_page($name) {
		if ($name == current_page()) print 'active';
	}
	
	function render_error($err) {
		// bootstrap 4 error alert: https://getbootstrap.com/docs/4.6/components/alerts/
		print('<div class="alert alert-danger" role="alert">');
		print($err);
		print('</div>');
	}
	
	function render_listdatabases() {
		global $sqlitecloud;
		$sqlitecloud = do_check_connect();
		$rs = $sqlitecloud->execute("LIST DATABASES;");
		if ($rs == false) return;
		
		for ($i=0; $i < $rs->nrows; ++$i) {
			$dbname = $rs->value($i, 0);
			echo "<option>{$dbname}</option>\n";
		}
	}
	
	function render_console_table($rs) {
		ob_start();
		render_table($rs);
		return ob_get_clean();
	}
		
	function render_table($rs, $oncolumn = NULL) {
		// check for error first
		if ($rs == false) {
			global $sqlitecloud;
			$sqlitecloud = do_check_connect();
			render_error($sqlitecloud->errmsg);
			return;
		}
		
		// table
		print('<table class="table table-striped table-sm">');
		
		// build header
		print("<thead><tr>\n");
		for ($i=0; $i < $rs->ncols; ++$i) {
		  $name = $rs->name($i);
		  print("<th>{$name}</th>");
		}
		print("</tr></thead>\n");
		
		// build values
		print("<tbody>\n");
		for ($i=0; $i < $rs->nrows; ++$i) {
			print("<tr>");
			for ($j=0; $j < $rs->ncols; ++$j) {
				$v = $rs->value($i, $j);
				$v2 = false;
				$value = NULL;
				if (!is_null($oncolumn)) $v2 = $oncolumn($v, $i, $j);
				if ($v2 === false) $value = (is_null($v)) ? 'N/A' : htmlentities($v);
				else $value = $v2;
				print("<td>{$value}</td>");
			}
			print("</tr>\n");
		  }
		print("</tbody>\n");
		print("</table>\n");
	}
	
	function render_chart_js($rs, $divID, $kvalue, $debug = false) {
		$script = "var ctx = document.getElementById('{$divID}');";
		$script .= "var myChart = new Chart(ctx, {type: 'line',";
		
		$labels = "data: {labels: [";
		$dataset = "datasets: [{data: [";
		
		// date, key, value
		//$cols = $rs->ncols;
		//$count = $rs->nrows / $rs->ncols;
		//if ($count > 15) $count = 15 * $rs->ncols;
		if (!$rs == false) {
			$count = $rs->nrows;
			for ($i=0; $i<$count; $i+=1) {
				$key = $rs->value($i, 1);
				if ($key == $kvalue) {
					$vdate = $rs->value($i, 0);
					$value = $rs->value($i, 2);
					$labels .= "'{$vdate}',";
					$dataset .= "{$value},";
				}
			}
		}
		
		$labels .= "],";
		$dataset .= "],";
		
		$script .= $labels;
		$script .= $dataset;
		
		$script .= "lineTension: 0, backgroundColor: 'transparent', borderColor: '#007bff', borderWidth: 4, pointBackgroundColor: '#007bff'}]";
		$script .= "},
			options: {
			  scales: {
				yAxes: [{
				  ticks: {
					beginAtZero: false
				  }
				}]
			  },
			  legend: {
				display: false
			  }
			}
		  });";
		  
		return $script;
	}
?>
