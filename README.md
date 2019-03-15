# go-cockroach
Scrappy load generation tool for cockroach via HTTP.

 * Heavy run load 11K INSERTS 100K SELECT Target per transaction (Stability / Schema run).
 * go get github.com/lib/pq
 * go get github.com/gorilla/mux
 * go get github.com/Pallinder/go-randomdata
 * go get github.com/geekzoo/go-cockroach
 * cd $GOPATH/github.com/geekzoo/go-cockroach/
 * vim cal.go 
 Change host to your cockroachdb host name and port, user name/password 
 I use HAProxy for load balancing to the back end database server, I'll include a skel config at some point.
 * go run cal.go
 * curl -s http://domain.com:8000/cal_prep #creates db and table
 * ab -c 100 -n 1000 -k http://domain.com:8000/cal_insert
 * ab -c 100 -n 1000 -k http://domain.com:8000/cal_all
 * in a web browser goto 
 * http://domain.com:8000/cal_insert
 * http://domain.com:8000/cal_all

<b>TODO:</b> 
  * Add runtime options for http listen port
  * Add connection to carbon and influx
