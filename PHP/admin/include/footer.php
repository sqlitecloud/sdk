</div>
  </div>

	  <script src="https://cdn.jsdelivr.net/npm/jquery@3.5.1/dist/jquery.slim.min.js" integrity="sha384-DfXdz2htPH0lsSSs5nCTpuj/zy4C+OGpamoFVy38MVBnE+IbbVYUew+OrCXaRkfj" crossorigin="anonymous"></script>
	  
	  <script src="https://cdn.jsdelivr.net/npm/bootstrap@4.6.1/dist/js/bootstrap.bundle.min.js" integrity="sha384-fQybjgWLrvvRgtW6bFlB7jaZrFsaBXjsOMm/tB9LTS58ONXgqbR9W8oWht/amnpF" crossorigin="anonymous"></script>

	  <script src="https://cdn.jsdelivr.net/npm/feather-icons@4.28.0/dist/feather.min.js"></script>
	  <script src="https://cdn.jsdelivr.net/npm/chart.js@2.9.4/dist/Chart.min.js"></script>
	  <script src="dashboard.js"></script>
	  <?php
	  		global $jscript;
			global $jsinclude;
			
			if (isset($jscript)) echo '<script>' . $jscript . '</script>';
			if (isset($jsinclude)) echo '<script src="' . $jsinclude . '"></script>';
			echo "\n";
	  ?>
  </body>
</html>

<?php do_disconnect(); ?>