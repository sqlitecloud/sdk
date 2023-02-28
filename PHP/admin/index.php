<?php include_once('include/header.php'); ?>

    <main role="main" class="col-md-9 ml-sm-auto col-lg-10 px-md-4">
    
    <h2 class="mt-3">Nodes</h2>
    <div class="table-responsive">
      <?php
        $rs = query_listnodes();
        render_table($rs);
      ?>
    </div>
     
     <?php $hostname = $_SESSION['hostname']; ?>
     <h2 class="mt-5"><?php echo $hostname; ?></h2>
     <div class="table-responsive">
       <?php
         $rs = query_nodeinfo();
         render_table($rs);
       ?>
     </div>

    </main>

<?php include_once('include/footer.php'); ?>
