<?php include_once('include/header.php'); ?>

    <main role="main" class="col-md-9 ml-sm-auto col-lg-10 px-md-4">
    <?php $rs = query_liststats(); ?>
    
    <h2 class="mt-3">Memory Usage</h2>
    <canvas class="my-4 w-100" id="memoryChart" width="900" height="380"></canvas>
    
    <h2 class="mt-3">Bytes In</h2>
    <canvas class="my-4 w-100" id="bytesInChart" width="900" height="380"></canvas>
    
    <h2 class="mt-3">Bytes Out</h2>
    <canvas class="my-4 w-100" id="bytesOutChart" width="900" height="380"></canvas>
    
    <h2 class="mt-3">CPU Usage</h2>
    <canvas class="my-4 w-100" id="cpuLoad" width="900" height="380"></canvas>
    
    <!--
    <h2 class="mt-3">Stats</h2>
    <div class="table-responsive">
      <? /*php render_table($rs);*/ ?>
    </div>
    -->

    </main>
    
    <?php
        global $jscript;
        $script1 = render_chart_js($rs, 'memoryChart', 'CURRENT_MEMORY');
        $script2 = render_chart_js($rs, 'bytesInChart', 'BYTES_IN');
        $script3 = render_chart_js($rs, 'bytesOutChart', 'BYTES_OUT', true);
        $script4 = render_chart_js($rs, 'cpuLoad', 'CPU_LOAD');
        
        $jscript = $script1 . "\n" . $script2 . "\n" . $script3 . "\n" . $script4 . "\n";
    ?>

<?php include_once('include/footer.php'); ?>
