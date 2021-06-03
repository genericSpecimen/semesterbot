#!/bin/sh
curl https://rlacollege.edu.in/view-all-details.php > new.php && diff -u old.php new.php | grep '^+.*href' | sed 's/+//' > diff.html