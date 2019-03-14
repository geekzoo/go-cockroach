# go-cockroach
Scrappy load generation tool for cockroach via HTTP
/*
 * Heavy run load 11K INSERTS 100K SELECT Target per transaction (Stability / Schema run).
 * curl -s http://ipORdomain:8000/cal_prep #creates db and table
 * Scraper v_2 go-bench-cal_insert 
 * ab -c 100 -n 1000 -k http://domain.com/cal_insert
 * ab -c 100 -n 1000 -k http://domain.com/cal_all
 * in a web browser goto 
 * http://domain.com/cal_insert
 * http://domain.com/cal_all
*/
