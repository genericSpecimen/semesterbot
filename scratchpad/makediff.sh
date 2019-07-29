#!/bin/sh
wget -o log https://rlacollege.edu.in/view-all-details.php -O new.php
diff -u old.php new.php | grep '^+.*href' | sed 's/+//' > diff.html