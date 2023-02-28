<?php include_once('include/header.php'); ?>

    <main role="main" class="col-md-9 ml-sm-auto col-lg-10 px-md-4">
    
    <h2 class="mt-3">Databases</h2>
    <div class="table-responsive">
      <?php
        function on_column($value, $row, $col) {
            if ($col != 0) return false;
            return '<a href="download.php?db=' . $value . '">' . $value . '</a>';
        }
        $rs = query_listdatabases(true);
        render_table($rs, "on_column");
      ?>
    </div>
    
    <!-- Upload disabled in this version
    <hr />
    
    <h2 class="mt-3">Database Upload</h2>
    <form id="upload-form" method="post" action="#">
        <p>Please select an SQLite database and click "Upload" to continue.</p>
        <div class="form-group">
            <progress id="upload_progress" value="0" max="100" class="col-4" style="display: none;"></progress>
            <div id="upload_error" class="alert alert-danger" role="alert" style="display: none;"></div>
            <div id="upload_message" class="alert alert-success col-4" role="alert" style="display: none;"></div>
        </div>
        
        <div class="form-group">
            <label for="enckey">Encryption key (optional)</label>
            <input type="text" class="form-control col-4" id="enckey">
        </div>
        <div class="form-group">
            <input type="file" class="form-control-file" id="datafile">
        </div>
        <button id="upload_button" class="btn btn-primary" type="submit">Upload</button>
    </form>
    -->

    </main>

<?php
    global $jsinclude;
    $jsinclude = "upload.js";
?>
<?php include_once('include/footer.php'); ?>
