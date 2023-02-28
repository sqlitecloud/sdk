<?php include_once('include/header.php'); ?>

    <main role="main" class="col-md-9 ml-sm-auto col-lg-10 px-md-4">
    
    <h2 class="mt-3">Logs</h2>
    <div class="table-responsive">
      <?php
        function on_column($value, $row, $col) {
            if ($col == 1) {
                // log type
                switch ($value) {
                    case 1: return 'INTERNAL';
                    case 2: return 'SECURITY';
                    case 3: return 'SQL';
                    case 4: return 'COMMAND';
                    case 5: return 'RAFT';
                    case 6: return 'CLUSTER';
                    case 7: return 'PLUGIN';
                    case 8: return 'CLIENT';
                    default: return $value;
                }
            }
            if ($col == 2) {
                // log level
                switch ($value) {
                    case 0: return '<span class="text-danger">PANIC</span>';
                    case 1: return '<span class="text-danger">FATAL</span>';
                    case 2: return '<span class="text-danger">ERROR</span>';
                    case 3: return '<span class="text-warning">WARNING</span>';
                    case 4: return 'INFO';
                    case 5: return '<span class="text-secondary">DEBUG</span>';
                    default: return $value;
                }
            }
            
            return false;
        }
        $rs = query_listlogs(50);
        render_table($rs, 'on_column');
      ?>
    </div>

    </main>

<?php include_once('include/footer.php'); ?>
