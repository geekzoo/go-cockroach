/*
 * Heavy run load 11K INSERTS 100K SELECTS per second (Stability / Schema run).
 * curl -s http://ipORdomain:8000/cal_prep #creates db and table
 * Scraper v_2 go-bench-cal_insert 
 * ab -c 100 -n 1000 -k http://domain.com/cal_insert
 * ab -c 100 -n 1000 -k http://domain.com/cal_all
 * in a web browser goto 
 * http://domain.com/cal_prep
 * http://domain.com/cal_insert
 * http://domain.com/cal_all
*/
package main

import (
    "fmt"
    "math"
    "math/rand"
    "time"
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/gorilla/mux"
    "net/http"
    _ "net/http/pprof"
    "log"
    crand  "crypto/rand"
    "io"
    "os"
    "sync/atomic"
    "github.com/Pallinder/go-randomdata"
)


const (
    host            = "127.0.0.1"		//PG HOST <GEO DNS round robbin to HAProxy> -> HAProxy -> LRU/RR cockroachdb nodes
    port            = 26258			//PG PORT
    user            = "root"			//DB USER NAME
    password        = ""			//DB PASSWORD
    dbname          = "uuid"			//Database Name
    fallback_application_name = "🦄Te🦄st🦄"	//Work around for db reported Application_name
    connect_timeout = 10				//DB Time out
    influxdb_host   = "insight.domain.com"	//TODO
    influxdb_port   = 6669			//TODO
    carbon_host     = "127.0.0.1"		//Carbon IP/Hostname <Graphite>
    carbon_port     = "2000"			//Carbon port
    carbon_link     = "US.GF.TESTING.TEST."	//TODO
    carbon_enabled  = false			//Enable=true Disable=false
    irc_host        = "irc.domain.com"		//TODO
    irc_port        = 6666			//TODO
    elastic_host    = "insight.domain.com"	//TODO
    elastic_port    = 8888			//TODO
    srv_host        = "0.0.0.0"			//IP/HostName for HTTP service
    srv_port        = 8000			//HTTP Port
    srv_w_timeout   = 15			//Seconds <Write Timeout>
    srv_r_timeout   = 15			//Seconds <Read Timeout>
    srv_idle_timeout    = 60			//Seconds <Idle Timeout>
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var ops, rps uint64 				//ops sql operations // rps http operations
var s_insert, s_select, s_update, s_delete uint64 //TODO count each type per ops breal down
var blow_out bool 				// circute for extermal stats
var lat, old_lat time.Duration 			//test to zero stats in no new data

func main() {

  if carbon_enabled == true {
    blow_out = true 				//circut breaker
  }
//s_sys()
//time.Sleep(1*time.Second)

fmt.Printf("\033c")

go func() {
    for {
      opsFinal := atomic.LoadUint64(&ops)
      if opsFinal == 18446744073709551615 { 
          fmt.Printf("\033[0;0H\033[2K\rOP/s: \033[33m%d\033[0m", 0)
          time.Sleep(time.Second)
      }else{
	fmt.Printf("\033[0;0H\033[2K\rOP/s:\033[32m %v\033[0m\033[C", opsFinal)
	atomic.AddUint64(&ops, ^opsFinal)
	time.Sleep(time.Second)
     }
    }
}()

go func() {
    for {
      rpsFinal := atomic.LoadUint64(&rps)
      if rpsFinal == 18446744073709551615 {
	  fmt.Printf("\033[2;0H\033[2K\rRP/s: \033[33m%d\033[0m", 0)
	  time.Sleep(time.Second)
      }else{
	  fmt.Printf("\033[2;0H\033[2K\rRP/s:\033[32m %v\033[0m\033[C", rpsFinal)
	  atomic.AddUint64(&rps, ^rpsFinal)
	  time.Sleep(time.Second)
      }
    }
}()

go func() {

    for {
        if blow_out == true {
	  fmt.Printf("\033[3;0H\033[2K\rStats Active:\033[32m %v\033[0m\033[C", blow_out)
	  time.Sleep(time.Second)
        }else{
	  fmt.Printf("\033[3;0H\033[2K\rStats Active:\033[93m %v\033[0m\033[C", blow_out)
	  time.Sleep(time.Second)
        }
    }
}()

go func() {

    for {
        if blow_out == true {
	  fmt.Printf("\033[4;0H\033[2K\rHealthy:\033[32m %v\033[0m\033[C", blow_out)
	  time.Sleep(time.Second)
        }else{
	  fmt.Printf("\033[4;0H\033[2K\rHealthy:\033[93m %v\033[0m\033[C", blow_out)
	  time.Sleep(time.Second)
        }
    }
}()

go func() {

    for {
        old_lat = lat
        if lat < 30000000 {
//		 23317344
	  fmt.Printf("\033[5;0H\033[2K\rLat:\033[32m %v\033[0m\033[C", lat)
	  fmt.Printf("\033[6;0H\033[2K\rLat RAW:\033[32m %d\033[0m\033[C", lat)
fmt.Printf("\033[14;0H\033[2K\rRAW:\033[95m %v, %v\033[0m\033[C", lat, old_lat)
	  time.Sleep(time.Second)
        }else{
	  fmt.Printf("\033[5;0H\033[2K\rLat:\033[93m %v\033[0m\033[C", lat)
	  fmt.Printf("\033[6;0H\033[2K\rLat RAW:\033[91m %d\033[0m\033[C", lat)
fmt.Printf("\033[14;0H\033[2K\rRAW:\033[91m %v, %v\033[0m\033[C", lat, old_lat)
	  time.Sleep(time.Second)
        }
    }
}()

go func() {
    for {
      s_selectFinal := atomic.LoadUint64(&s_select)
      if s_selectFinal == 18446744073709551615 { 
          fmt.Printf("\033[9;0H\033[2K\rSELECTS/s: \033[33m%d\033[0m", 0)
          time.Sleep(time.Second)
      }else{
	fmt.Printf("\033[9;0H\033[2K\rSELECTS/s:\033[32m %v\033[0m\033[C", s_selectFinal)
	atomic.AddUint64(&s_select, ^s_selectFinal)
	time.Sleep(time.Second)
     }
    }
}()

go func() {
    for {
      s_updateFinal := atomic.LoadUint64(&s_update)
      if s_updateFinal == 18446744073709551615 { 
          fmt.Printf("\033[10;0H\033[2K\rUPDATES/s: \033[33m%d\033[0m", 0)
          time.Sleep(time.Second)
      }else{
	fmt.Printf("\033[10;0H\033[2K\rUPDATES/s:\033[32m %v\033[0m\033[C", s_updateFinal)
	atomic.AddUint64(&s_update, ^s_updateFinal)
	time.Sleep(time.Second)
     }
    }
}()

go func() {
    for {
      s_insertFinal := atomic.LoadUint64(&s_insert)
      if s_insertFinal == 18446744073709551615 { 
          fmt.Printf("\033[11;0H\033[2K\rINSERTS/s: \033[33m%d\033[0m", 0)
          time.Sleep(time.Second)
      }else{
	fmt.Printf("\033[11;0H\033[2K\rINSERTS/s:\033[32m %v\033[0m\033[C", s_insertFinal)
	atomic.AddUint64(&s_insert, ^s_insertFinal)
	time.Sleep(time.Second)
     }
    }
}()

go func() {
    for {
      s_deleteFinal := atomic.LoadUint64(&s_delete)
      if s_deleteFinal == 18446744073709551615 { 
          fmt.Printf("\033[12;0H\033[2K\rDELETES/s: \033[33m%d\033[0m", 0)
          time.Sleep(time.Second)
      }else{
	fmt.Printf("\033[12;0H\033[2K\rDELETES/s:\033[32m %v\033[0m\033[C", s_deleteFinal)
	atomic.AddUint64(&s_delete, ^s_deleteFinal)
	time.Sleep(time.Second)
     }
    }
}()

/* Reset blow_out */
go func() {
    for {
      if carbon_enabled == true {
        blow_out = true
      }
        time.Sleep(10*time.Second)
     }
}()
  
    r := mux.NewRouter()
    srv := &http.Server{
        Addr:           fmt.Sprintf("%v:%v",srv_host,srv_port),
        WriteTimeout:   time.Second * srv_w_timeout,
        ReadTimeout:    time.Second * srv_r_timeout,
        IdleTimeout:    time.Second * srv_idle_timeout,
        Handler: r,
    }
    

    r.HandleFunc("/cal_insert", cal_insert).Methods("GET")
    r.HandleFunc("/cal_all", cal_all).Methods("GET")
    r.HandleFunc("/cal_prep", cal_prep).Methods("GET")
    r.HandleFunc("/cal_truncate", cal_truncate).Methods("GET")
    r.HandleFunc("/show_sessions", show_sessions).Methods("GET")
    
    log.Fatal(srv.ListenAndServe())
}
// END MAIN

        func cal_insert(w http.ResponseWriter, r *http.Request) {
            start_init := time.Now()
            w.Header().Set("X-ENGINE", "V2")
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("X-TEST", "cal_insert")
            
            psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s fallback_application_name=%s connect_timeout=%d sslmode=disable", host, port, user, password, dbname, fallback_application_name, connect_timeout)
            db, err := sql.Open("postgres", psqlInfo)
                if err != nil {
            panic(err)
                }
            defer db.Close()
            err = db.Ping()
                if err != nil {
            panic(err)
                }

            w.Write([]byte("<body style=background-color:grey;>"))
            w.Write([]byte("<h1> cal_insert </h1>"))

            w.Write([]byte("<svg width=35 height=35> <circle cx=20 cy=20 r=10 stroke=green stroke-width=4 fill=yellow /> </svg> <font size=1px>Warm up.</font>"))
            w.Write([]byte("<svg width=35 height=35> <circle cx=20 cy=20 r=10 stroke=green stroke-width=4 fill=green /> </svg> <font size=1px>Done. </font>"))

	for run_c := 1; run_c <= 10; run_c++ {
            res_id, err := newUUID()
            calendar_id, err := newUUID()
            reservation_id, err := newUUID()
            company_id := rand.Intn(1156)
            title := randomdata.Paragraph()
            location := randomdata.ProvinceForCountry("US")
            organizer_email := randomdata.Email()
            reservation_begin, reservation_end := gen_meeting()
            t1 := time.Unix(reservation_begin, 0).UTC()
            t2 := time.Unix(reservation_end, 0).UTC()
            
            start2 := time.Now()
            err = db.QueryRow("INSERT INTO uuid.cal_insert (calendar_id, reservation_id, company_id, title, location, organizer_email, reservation_begin, reservation_end) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING reservation_id", calendar_id, res_id, company_id, title, location, organizer_email, t1, t2).Scan(&reservation_id)
// OP/S counter
atomic.AddUint64(&ops, 1)
atomic.AddUint64(&s_insert, 1)
            elapsed2 := time.Since(start2)
            tt2 := fmt.Sprintf("Time SQL:= %s", elapsed2)
            w.Write([]byte("<b><pre>MULTIUUID " + " " + calendar_id + " <font color=green>" + reservation_id + "</font> " + fmt.Sprintf(" co:[%d] ",company_id) + "<font size=0.01px><p>" + title + "</font></p>" + "<font color=yellow> " + tt2 + "</font></b></pre>"))

            //Carbon Format
            now := time.Now()
            epoc_now := now.Unix()
            hostname, err := os.Hostname()
lat = elapsed2


	    if carbon_enabled == true && blow_out == true {
	    Tcc(fmt.Sprintf("GF.TEST.%s.CAL-INSERT.SQL-FUNC %d %d", hostname, elapsed2, epoc_now))
//	    Tcc(fmt.Sprintf("GF.TEST.%s.CAL-INSERT.OPS %d %d", hostname, s_insertFinal, epoc_now))
	    }
            
            if err != nil {}

        }
        
    db.Close()
    w.Write([]byte("DONE"))
    now := time.Now()
    epoc_now := now.Unix()
    hostname, err := os.Hostname()
    elapsed := time.Since(start_init)
    
    if carbon_enabled == true && blow_out == true {
    Tcc(fmt.Sprintf("GF.TEST.%s.CAL-BLK.INSERT.SQL-FUNC %d %d", hostname, elapsed, epoc_now))
    }
    
// RP/S counter 
atomic.AddUint64(&rps, 1)
    return 
    }
    
    /*
     */
    
        func cal_all(w http.ResponseWriter, r *http.Request) {
            start_init := time.Now()     
            w.Header().Set("X-ENGINE", "V2")
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("X-TEST", "cal_all")
            
            psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s fallback_application_name=%s connect_timeout=%d sslmode=disable", host, port, user, password, dbname, fallback_application_name, connect_timeout)
            db, err := sql.Open("postgres", psqlInfo)
                if err != nil {
                panic(err)
                }
            defer db.Close()
            err = db.Ping()
                if err != nil {
                panic(err)
                }
                    
            w.Write([]byte("<body style=background-color:grey;>"))
            ti := time.Now()
            w.Write([]byte("<h1> cal_all </h1>" + "<h2>" + fmt.Sprintf("%v", ti) + "</h2>" ))
            if err != nil {
                log.Fatal(err)
            }
            if err != nil {
                panic(err)
            }

            w.Write([]byte("<svg width=35 height=35> <circle cx=20 cy=20 r=10 stroke=green stroke-width=4 fill=yellow /> </svg> <font size=1px>Warm up.</font>"))
            w.Write([]byte("<svg width=35 height=35> <circle cx=20 cy=20 r=10 stroke=green stroke-width=4 fill=green /> </svg> <font size=1px>Done. </font>"))
        
        for run_c := 1; run_c <= 1; run_c++ {
            reservation_id := ""
            calendar_id := ""
            title := ""
            company_id := 0
            rand_0 := rand.Intn(1156)
            start2 := time.Now()
            err = db.QueryRow("SELECT calendar_id, reservation_id, title, company_id FROM uuid.cal_insert WHERE company_id = $1 LIMIT 1", rand_0).Scan(&calendar_id, &reservation_id, &title, &company_id)

atomic.AddUint64(&ops, 1)
atomic.AddUint64(&s_select, 1)

            elapsed2 := time.Since(start2)
lat = elapsed2	    
            tt2 := fmt.Sprintf("Time SQL:= %s", elapsed2)
            w.Write([]byte("<b><pre>cal_select " + fmt.Sprintf("%v", company_id) + " " + title + "<font color=yellow> " + tt2 + "</font></b></pre>"))
            now := time.Now()
            epoc_now := now.Unix()
            hostname, err := os.Hostname()
	    
	    if carbon_enabled == true {
	    Tcc(fmt.Sprintf("GF.TEST.%s.CAL-SELECT.SQL-FUNC %d %d", hostname, elapsed2, epoc_now))
	    }
            
            title = randomdata.Paragraph()
            start3 := time.Now()
            err = db.QueryRow("UPDATE uuid.cal_insert SET title=$4 WHERE company_id=$1 AND calendar_id=$2 AND reservation_id=$3 RETURNING title", company_id, calendar_id, reservation_id, title).Scan(&title)
atomic.AddUint64(&ops, 1)
atomic.AddUint64(&s_update, 1)
            elapsed3 := time.Since(start3)
lat = elapsed3	    
            tt3 := fmt.Sprintf("Time SQL:= %s", elapsed3)
            
            w.Write([]byte("<b><pre>cal_update " + fmt.Sprintf("%v", company_id) + " " + title + "<font color=yellow> " + tt3 + "</font></b></pre>"))
            //Carbon Format
            now = time.Now()
            epoc_now = now.Unix()
            hostname, err = os.Hostname()
	    
	    if carbon_enabled == true {
	    Tcc(fmt.Sprintf("GF.TEST.%s.CAL-UPDATE.SQL-FUNC %d %d", hostname, elapsed3, epoc_now))
	    }
/*
 */
            rand_1 := rand.Intn(1156)
            start4 := time.Now()
	    err = db.QueryRow("SELECT calendar_id, reservation_id, title, company_id FROM uuid.cal_insert WHERE company_id = $1 LIMIT 1", rand_1).Scan(&calendar_id, &reservation_id, &title, &company_id)
atomic.AddUint64(&s_select, 1)	    
atomic.AddUint64(&ops, 1)	    
            elapsed4 := time.Since(start4)
lat = elapsed4	    
            tt4 := fmt.Sprintf("Time SQL:= %s", elapsed4)
            w.Write([]byte("<b><pre>cal_select " + fmt.Sprintf("%v", company_id) + " " + title + "<font color=yellow> " + tt4 + "</font></b></pre>")) 
            now = time.Now()
            epoc_now = now.Unix()
            hostname, err = os.Hostname()

	    if carbon_enabled == true {
	    Tcc(fmt.Sprintf("GF.TEST.%s.CAL-INSERT.SQL-FUNC %d %d", hostname, elapsed2, epoc_now))
	    }
            //fmt.Printf("GF.TEST.%s.CAL-SELECT %d %d\n", hostname, elapsed4, epoc_now) 
/*
 */
            rand_2 := rand.Intn(1156)
            start5 := time.Now()
            err = db.QueryRow("SELECT calendar_id, reservation_id, company_id FROM uuid.cal_insert WHERE company_id = $1 LIMIT 1", rand_2).Scan(&calendar_id, &reservation_id, &company_id)
atomic.AddUint64(&s_select, 1)	    
atomic.AddUint64(&ops, 1)	    
            elapsed5 := time.Since(start5)
lat = elapsed5	    
            tt5 := fmt.Sprintf("Time SQL:= %s", elapsed5)
            w.Write([]byte("<b><pre>cal_select " + fmt.Sprintf("%v", company_id) + " " + reservation_id + "<font color=yellow> " + tt5 + "</font></b></pre>")) 
            now = time.Now()
            epoc_now = now.Unix()
            hostname, err = os.Hostname()
            
	    if carbon_enabled == true {
	    Tcc(fmt.Sprintf("GF.TEST.%s.CAL-SELECT.SQL-FUNC %d %d", hostname, elapsed5, epoc_now))
	    }
	    //fmt.Printf("GF.TEST.%s.CAL-SELECT %d %d\n", hostname, elapsed5, epoc_now) 

            start6 := time.Now()
            db.Exec(`DELETE FROM uuid.cal_insert WHERE company_id=$1 AND reservation_id=$2 AND calendar_id=$3;`, company_id, reservation_id, calendar_id)
atomic.AddUint64(&ops, 1)
atomic.AddUint64(&s_delete, 1)
            elapsed6 := time.Since(start6)
lat = elapsed6	    
            tt6 := fmt.Sprintf("Time SQL:= %s", elapsed6)
            w.Write([]byte("<b><pre>cal_delete " + fmt.Sprintf("%v", company_id) + " " + reservation_id + "<font color=yellow> " + tt6 + "</font></b></pre>")) 
            now = time.Now()
            epoc_now = now.Unix()
            hostname, err = os.Hostname()
	    
	    if carbon_enabled == true {
	    Tcc(fmt.Sprintf("GF.TEST.%s.CAL-DELETE.SQL-FUNC %d %d", hostname, elapsed6, epoc_now))
	    }
            //fmt.Printf("GF.TEST.%s.CAL-DELETE %d %d\n", hostname, elapsed6, epoc_now)
            
            
            rand_3 := rand.Intn(1156)
            res_id, err := newUUID()
            calendar_id, err = newUUID()
            reservation_id, err = newUUID()
            company_id = rand_3
            title = randomdata.Paragraph()
            location := randomdata.ProvinceForCountry("US")
            organizer_email := randomdata.Email()
            reservation_begin, reservation_end := gen_meeting()
            t1 := time.Unix(reservation_begin, 0).UTC()
            t2 := time.Unix(reservation_end, 0).UTC()
            
            start7 := time.Now()
            err = db.QueryRow("INSERT INTO uuid.cal_insert (calendar_id, reservation_id, company_id, title, location, organizer_email, reservation_begin, reservation_end) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING reservation_id", calendar_id, res_id, company_id, title, location, organizer_email, t1, t2).Scan(&reservation_id)

atomic.AddUint64(&ops, 1)
atomic.AddUint64(&s_insert, 1)

            elapsed7 := time.Since(start7)
lat = elapsed7	    
            tt7 := fmt.Sprintf("Time SQL:= %s", elapsed7)
            w.Write([]byte("<b><pre>cal_insert " + fmt.Sprintf("%v", company_id) + " " + reservation_id + "<font color=yellow> " + tt7 + "</font></b></pre>"))
            now = time.Now()
            epoc_now = now.Unix()
            hostname, err = os.Hostname()
	    
	    if carbon_enabled == true {
	    Tcc(fmt.Sprintf("GF.TEST.%s.CAL-INSERT.SQL-FUNC %d %d", hostname, elapsed7, epoc_now))
	    }
            //fmt.Printf("GF.TEST.%s.CAL-INSERT %d %d\n", hostname, elapsed7, epoc_now) //Write Carbon

            
    for sel_c := 1; sel_c <= 40; sel_c++ {
            rand_0 := rand.Intn(1156)
            start := time.Now()
            err = db.QueryRow("SELECT calendar_id, reservation_id, company_id FROM uuid.cal_insert WHERE company_id = $1 LIMIT 1", rand_0).Scan(&calendar_id, &reservation_id, &company_id)
atomic.AddUint64(&s_select, 1)
atomic.AddUint64(&ops, 1)
            elapsed := time.Since(start)
            tt1 := fmt.Sprintf("Time SQL:= %s", elapsed)
            w.Write([]byte("<b><pre>cal_select_loop " + fmt.Sprintf("%v", company_id) + " " + reservation_id + "<font color=yellow> " + tt1 + "</font></b></pre>"))
            now = time.Now()
            epoc_now = now.Unix()
            hostname, err = os.Hostname()
	    
	    if carbon_enabled == true {
	    Tcc(fmt.Sprintf("GF.TEST.%s.CAL-SELECT.SQL-FUNC %d %d", hostname, elapsed, epoc_now))
	    }
            //fmt.Printf("GF.TEST.%s.CAL-SELECT %d %d\n", hostname, elapsed, epoc_now)
    }
            
            run_c++
            if err != nil {}

        }
    db.Close()
    w.Write([]byte("DONE WHHEEEEE!"))
    now := time.Now()
    epoc_now := now.Unix()
    hostname, err := os.Hostname()
    elapsed := time.Since(start_init)

atomic.AddUint64(&rps, 1)

    if carbon_enabled == true {
    Tcc(fmt.Sprintf("GF.TEST.%s.CAL-BLK.SQL-FUNC %d %d", hostname, elapsed, epoc_now))
    }
    //fmt.Printf("GF.TEST.%s.CAL-BLK.SQL-FUNC %d %d\n", hostname, elapsed, epoc_now) //Write Carbon
    
    }
    
    
func cal_prep(w http.ResponseWriter, r *http.Request) {
        
            w.Header().Set("X-ENGINE", "V2")
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("X-TEST", "cal_prep")
            w.Header().Set("Content-Type", "text/html")
            
            psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s fallback_application_name=%s connect_timeout=%d sslmode=disable", host, port, user, password, dbname, fallback_application_name, connect_timeout)
            db, err := sql.Open("postgres", psqlInfo)
                if err != nil {
            panic(err)
                }
            defer db.Close()
            err = db.Ping()
                if err != nil {
            panic(err)
                }

		fmt.Fprintf(w, "<pre>CREATE DATABASE uuid<br>")
		db.Exec(`CREATE DATABASE uuid;`)
		fmt.Fprintf(w, "CREATE TABLE IF NOT EXISTS uuid.cal_insert<br>")
                db.Exec(`CREATE TABLE IF NOT EXISTS uuid.cal_insert (
                    calendar_id                              UUID NOT NULL,
                    reservation_id                           UUID NOT NULL,
                    company_id                               INT NOT NULL,
                    title                                    STRING NOT NULL,
                    location                                 STRING NOT NULL,
                    organizer_email                          STRING NOT NULL,
                    time_now                                 TIMESTAMPTZ NOT NULL DEFAULT now(),
                    reservation_begin                        TIMESTAMP WITHOUT TIME ZONE NOT NULL,
                    reservation_end                          TIMESTAMP WITHOUT TIME ZONE NOT NULL,
                    PRIMARY KEY                              (calendar_id, reservation_id),
                    INDEX company_id_idx                     (company_id),
                    INDEX reservation_begin_idx              (reservation_begin),
                    INDEX reservation_end_idx                (reservation_end)
                );`)
		fmt.Fprintf(w, "DONE")
		
}


func cal_truncate(w http.ResponseWriter, r *http.Request) {
        
            w.Header().Set("X-ENGINE", "V2")
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("X-TEST", "cal_truncate")
            w.Header().Set("Content-Type", "text/html")
            
            psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s fallback_application_name=%s connect_timeout=%d sslmode=disable", host, port, user, password, dbname, fallback_application_name, connect_timeout)
            db, err := sql.Open("postgres", psqlInfo)
                if err != nil {
            panic(err)
                }
            defer db.Close()
            err = db.Ping()
                if err != nil {
            panic(err)
                }
		fmt.Fprintf(w, "<pre>TRUNCATE uuid.cal_insert<br>")
		db.Exec(`TRUNCATE TABLE uuid.cal_insert;`)
		fmt.Fprintf(w, "DONE")
}


 func show_sessions(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("X-ENGINE", "V8")
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("X-NIKI", "FOREST GODES")
            psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s fallback_application_name=%s connect_timeout=%d sslmode=disable", host, port, user, password, dbname, fallback_application_name, connect_timeout)
            db, err := sql.Open("postgres", psqlInfo)
                if err != nil {
            panic(err)
                }
            defer db.Close()
            err = db.Ping()
                if err != nil {
            panic(err)
                }
//            fmt.Printf("\033[92mSuccessfully connected!\033[0m\n")
            w.Write([]byte("<body style=background-color:grey;>"))
            w.Write([]byte("<h1> DB SESSIONS </h1>"))
            
            node_id := ""
            session_id := ""
            user_name := ""
            client_address := ""
            application_name := ""
            active_queries := ""
            last_active_query := ""
            session_start := ""
            oldest_query_start := ""
            
        rows, err := db.Query("SELECT * FROM [SHOW CLUSTER SESSIONS]") //.Scan( &node_id, &session_id, &user_name, &client_address, &application_name, &active_queries, &last_active_query, &session_start, &oldest_query_start)
                if err != nil {
                    panic(err)
                }
            defer rows.Close()
        for rows.Next() {
            //var node_id string
            //NEED to ADD condition for nil when a cockroach node is down!!!
            var oldest_query_start sql.NullString
            err = rows.Scan(&node_id, &session_id, &user_name, &client_address, &application_name, &active_queries, &last_active_query, &session_start, &oldest_query_start)
                if err != nil {
                    panic(err)
                }
            w.Write([]byte("<b><pre><font color=cyan> " + "|" + node_id + "|" + session_id + "|" + user_name + "|" + client_address + "|" + application_name + "|" + active_queries + "|" + last_active_query + "|" + session_start + "|" + fmt.Sprintf( "%v", oldest_query_start ) + "</b></pre></font>"))
//            fmt.Println("\033[91m" + node_id, session_id, user_name, client_address, application_name, active_queries, last_active_query, session_start, oldest_query_start)
        }
            w.Write([]byte("<b><pre><font color=yellow> " + "|" + node_id + "|" + session_id + "|" + user_name + "|" + client_address + "|" + application_name + "|" + active_queries + "|" + last_active_query + "|" + session_start + "|" + fmt.Sprintf( "%v", oldest_query_start ) + "</b></pre></font>"))
//            fmt.Println(node_id, session_id, user_name, client_address, application_name, active_queries, last_active_query, session_start, oldest_query_start)
   }

    func newUUID() (string, error) {
        uuid := make([]byte, 16)
        n, err := io.ReadFull(crand.Reader, uuid)
        if n != len(uuid) || err != nil {
                return "", err
        }
        // variant bits; see section 4.1.1
        uuid[8] = uuid[8]&^0xc0 | 0x80
        // version 4 (pseudo-random); see section 4.1.3
        uuid[6] = uuid[6]&^0xf0 | 0x40
        return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil 
    }

func gen_meeting() (int64, int64) {

        var log_max float64
        var random float64
        var log_random float64
        var six_months float64 = 15768000
        var one_day float64 = 86400
        var min_meeting float64 = 300
    rand.Seed(time.Now().UnixNano())
    rand_s1 := rand.Intn(32768)
    log_max = math.Log(32768)
    random = float64(rand_s1)
    log_random = math.Log(random)
    random_ratio := log_random / log_max
    random_future := random_ratio * six_months
    random_time := random_ratio * one_day
    date_epoc := time.Now().Unix()
    fdate_epoc := float64(date_epoc)
    real_future := fdate_epoc + six_months - random_future / 1
    real_apt := real_future + one_day - random_time + min_meeting //HAck job min_meeting
    scale_m := int64(real_future)
    scale_d := int64(real_apt)
        return scale_m, scale_d

}

func RandStringRunes(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}
