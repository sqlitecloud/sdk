<?php
	include_once('auth.php');
	include_once('common.php');
	$page = current_page();
?>

<!doctype html>
<html lang="en">
  <head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	<title>SQLite Cloud Dashboard</title>

	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@4.6.1/dist/css/bootstrap.min.css" integrity="sha384-zCbKRCUGaJDkqS1kPbPd7TveP5iyJE0EjAuZQTgFLD2ylzuqKfdKlfG/eSrtxUkn" crossorigin="anonymous">

	<style>
	  .bd-placeholder-img {
		font-size: 1.125rem;
		text-anchor: middle;
		-webkit-user-select: none;
		-moz-user-select: none;
		-ms-user-select: none;
		user-select: none;
	  }

	  @media (min-width: 768px) {
		.bd-placeholder-img-lg {
		  font-size: 3.5rem;
		}
	  }
	</style>

	
	<!-- Custom styles for this template -->
	<link href="dashboard.css" rel="stylesheet">
  </head>
  <body>
	
<nav class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
  <a class="navbar-brand col-md-3 col-lg-2 mr-0 px-3" href="http://sqlitecloud.io" target="_blank">SQLite Cloud</a>
  <button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-toggle="collapse" data-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
	<span class="navbar-toggler-icon"></span>
  </button>
  <!--<input class="form-control form-control-dark w-100" type="text" placeholder="Search" aria-label="Search">-->
  <ul class="navbar-nav px-3">
	<li class="nav-item text-nowrap">
	  <a class="nav-link" href="logout.php">Sign out</a>
	</li>
  </ul>
</nav>

<!-- icons from: https://feathericons.com -->
<div class="container-fluid">
  <div class="row">
	<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
	  <div class="sidebar-sticky pt-3">
		<ul class="nav flex-column">
		  <li class="nav-item">
			<a class="nav-link <?php render_current_page('index'); ?>" href="index.php">
			  <span data-feather="home"></span>
			  Nodes
			</a>
		  </li>
		  <li class="nav-item">
			<a class="nav-link <?php render_current_page('databases'); ?>" href="databases.php">
			  <span data-feather="file"></span>
			  Databases
			</a>
		  </li>
		  <li class="nav-item">
			<a class="nav-link <?php render_current_page('commands'); ?>" href="commands.php">
			  <span data-feather="archive"></span>
			  Commands
			</a>
		  </li>
		  <li class="nav-item">
			<a class="nav-link <?php render_current_page('users'); ?>" href="users.php">
			  <span data-feather="users"></span>
			  Users
			</a>
		  </li>
		  <li class="nav-item">
			<a class="nav-link <?php render_current_page('connections'); ?>" href="connections.php">
			  <span data-feather="layers"></span>
			  Connections
			</a>
		  </li>
		  <li class="nav-item">
			<a class="nav-link <?php render_current_page('logs'); ?>" href="logs.php">
			  <span data-feather="align-left"></span>
			  Logs
			</a>
		  </li>
		  <li class="nav-item">
			  <a class="nav-link <?php render_current_page('plugins'); ?>" href="plugins.php">
				<span data-feather="codesandbox"></span>
				Plugins
			  </a>
			</li>
		  <li class="nav-item">
			  <a class="nav-link <?php render_current_page('stats'); ?>" href="stats.php">
				<span data-feather="pie-chart"></span>
				Stats
			  </a>
		  </li>
		  <li class="nav-item">
				<a class="nav-link <?php render_current_page('settings'); ?>" href="settings.php">
				  <span data-feather="sliders"></span>
				  Settings
				</a>
			</li>
		  <li class="nav-item">
			<a class="nav-link <?php render_current_page('latency'); ?>" href="latency.php">
				<span data-feather="activity"></span>
				Latency
			</a>
		  </li>
		  <li class="nav-item">
			<a class="nav-link <?php render_current_page('backups'); ?>" href="backups.php">
			  <span data-feather="archive"></span>
			  Backup
			</a>
		  </li>
		  <li class="nav-item">
			  <a class="nav-link <?php render_current_page('console'); ?>" href="console.php">
				<span data-feather="terminal"></span>
				Console
			  </a>
			</li>
		</ul>

		<!--
		<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted">
		  <span>Saved reports</span>
		  <a class="d-flex align-items-center text-muted" href="#" aria-label="Add a new report">
			<span data-feather="plus-circle"></span>
		  </a>
		</h6>
		<ul class="nav flex-column mb-2">
		  <li class="nav-item">
			<a class="nav-link" href="#">
			  <span data-feather="file-text"></span>
			  Current month
			</a>
		  </li>
		  <li class="nav-item">
			<a class="nav-link" href="#">
			  <span data-feather="file-text"></span>
			  Last quarter
			</a>
		  </li>
		  <li class="nav-item">
			<a class="nav-link" href="#">
			  <span data-feather="file-text"></span>
			  Social engagement
			</a>
		  </li>
		  <li class="nav-item">
			<a class="nav-link" href="#">
			  <span data-feather="file-text"></span>
			  Year-end sale
			</a>
		  </li>
		-->
		</ul>
	  </div>
	</nav>
	