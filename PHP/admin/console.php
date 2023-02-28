<?php include_once('include/header.php'); ?>

    <main role="main" class="col-md-9 ml-sm-auto col-lg-10 px-md-4">
    
    <h2 class="mt-3">Console</h2>
    
    <div class="col-md-6">
    <form id="console-form" method="post" action="#">
        <div class="form-group">
            <label for="database">Database</label>
            <select class="form-control" id="database">
                <?php render_listdatabases(); ?>
            </select>
        </div>
        <div class="form-group">
            <textarea class="form-control" id="sql" rows="3"></textarea>
        </div>
        
        <div class="form-group">
            <div id="console-error" style="display: none;" class="alert alert-danger col-8" role="alert"></div>
            <div id="console-message" style="display: none;" class="alert alert-success col-8" role="alert"></div>
        </div>
        
        <button id="console-button" class="btn btn-primary" type="submit">Execute</button>
    </form>
    </div>
    
    <div id="console-table" class="mt-3" style="display: none;">
    </div>
    
    </main>

<?php
    global $jsinclude;
    $jsinclude = "console.js";
?>
<?php include_once('include/footer.php'); ?>
