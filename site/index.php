<!doctype html>
<?php include 'includes/header.php';?>

<html class="no-js" lang="en">

<body>

    <?php include 'includes/menu.php';?>

    <div id="slider" class="container">
        <div class="row d-flex align-items-center">
            <div class="col-md-5">
                <img class="centrar" src="img/b-icon.svg">
                <h1><span>One app with the all right tools.</span></h1>
            </div>
            <div class="col-md-8 push-1">
                <div class="owl-carousel owl-theme">
                    <img src="img/img-slider.jpg" alt="texto alterno">
                    <img src="img/img-slider.jpg" alt="texto alterno">
                    <img src="img/img-slider.jpg" alt="texto alterno">
                </div>
            </div>
        </div>
    </div>

    <div id="bw-input">
        <div class="container">
            <div class="row no-gutters d-flex justify-content-center">
                <div class="col-10 col-md-12">
                    <h2 class="text-center mb-4">Join the up coming alpha test:</h2>
                </div>
                <div class="col-10 col-md-6 cont-input mt-4">
                    <form class="form-inline">
                        <div class="row w-100 no-gutters">
                            <div class="col-9">
                                <div class="form-group">
                                    <label for="emailform" class="sr-only">Email</label>
                                    <input type="text" class="form-control-plaintext" id="emailform" placeholder="email@example.com">
                                </div>
                            </div>
                            <div class="col text-right d-flex justify-content-end">
                                <button type="submit" class="btn align-self-center"><img src="img/send-arrow.svg"></button>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>

    <?php include 'includes/footer.php';?>

</body>

</html>
