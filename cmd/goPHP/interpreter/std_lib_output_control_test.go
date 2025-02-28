package interpreter

// TODO tests for ob_clean
// TODO tests for ob_flush
// TODO tests for ob_end_clean
// TODO tests for ob_end_flush
// TODO tests for ob_get_clean
// TODO tests for ob_get_flush
// TODO tests for ob_get_contents
// TODO tests for ob_get_level
// TODO tests for ob_start

// TODO test ob_start without explicit closing -> automatically closing and flushing

/*
<?php
echo 0;
    ob_start();
        ob_start();
            ob_start();
                ob_start();
                    echo 1;
                ob_end_flush();
                echo 2;
            $ob = ob_get_clean();
        echo 3;
        ob_flush();
        ob_end_clean();
    echo 4;
    ob_end_flush();
echo '-' . $ob;
?>
--EXPECT--
034-12
*/
