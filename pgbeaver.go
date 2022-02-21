package main

import (
    "log"
    "os"
    "net/http"
    "flag"
    "fmt"
    "path"
    "time"
    "os/exec"
    "runtime"
    "strings"
    "math"
    "strconv"
    "io/ioutil"
    "path/filepath"
)

const version = "0.0.2"

const pgbadgerimgurl = "https://pgbadger.darold.net/logo_pgbadger.png"

var configfile string
var allocationfile string
var walletfile string
var debug bool
var allowdelete bool

var username string
var password string

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func showfilesize(bytes int64) string {
	if(bytes < 1000) {
		return(fmt.Sprintf("%d bytes", bytes)) 
	}
	fbytes := float64(bytes)
	if(bytes < 1000000) {
		return(fmt.Sprintf("%0.1f KB", float64(fbytes/1000))) 
	}
	if(bytes < 1000000000) {
		return(fmt.Sprintf("%0.2f MB", float64(fbytes/1000000))) 
	}			
	return(fmt.Sprintf("%0.3f GB", float64(fbytes/1000000000))) 
}

func microTime() float64 {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	micSeconds := float64(now.Nanosecond()) / 1000000000
	return float64(now.Unix()) + micSeconds
}

func getTmpPath() string {
	tmp := fmt.Sprintf("%d", time.Now().UnixNano())
    return (tmp)
}
	
func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB\t", bToMb(m.Alloc))
	fmt.Printf("TotalAlloc = %v MiB\t", bToMb(m.TotalAlloc))
	fmt.Printf("Sys = %v MiB\t", bToMb(m.Sys))
	fmt.Printf("NumGC = %v\n", m.NumGC)
}
	
func WalkMatch(root, pattern string) ([]string, error) {
    var matches []string
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
            return err
        } else if matched {
            matches = append(matches, path)
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return matches, nil
}

func isNumeric(s string) bool {
    _, err := strconv.ParseFloat(s, 64)
    return err == nil
}			



func main() {
		
		var logpath string
		var logprefix string

		logpath = "./"
		logprefix = "postgresql-"
		
		http.HandleFunc("/version/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Version %s", version)
		})		
		
		http.HandleFunc("/info/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Version %s\n", version)
			fmt.Fprintf(w, "Select Report above, output will appear here\n")
		})	

		http.HandleFunc("/message/", func(w http.ResponseWriter, r *http.Request) {
			msg, ok := r.URL.Query()["message"]	
			if(ok) {
				fmt.Fprintf(w, "%s", msg[0])
			}
		})	
		
		
		
		
		
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {	 

			var u string
			var p string

			if(len(username)>0) && (len(password)>0) {
				var ok bool
			    ua, uok := r.URL.Query()["username"]		    
				pa, pok := r.URL.Query()["password"]
				if(uok && pok) {
					u = ua[0]
					p = pa[0]				
				} else {
					u, p, ok = r.BasicAuth()
					if !ok {
						fmt.Println("Error parsing basic auth")
						w.WriteHeader(401)
						return
					}
				}
				if u != username {
					fmt.Printf("Username provided is correct: %s\n", u)
					w.WriteHeader(401)
					return
				}
				if p != password {
					fmt.Printf("Password provided is correct: %s\n", p)
					w.WriteHeader(401)
					return
				}
			}

			var whiteliststr string
		    pwhitelists, ok := r.URL.Query()["whitelist"]		    
		    if ok && len(pwhitelists[0]) >0 {
				pwhitelist := pwhitelists[0]
				fmt.Println("whitelist: " + string(pwhitelist))			
				if(len(pwhitelist)>0) {
					whiteliststr=pwhitelist
				}
			}

			var blackliststr string			
			blackliststr = "pg_catalog"
		    pblacklists, ok := r.URL.Query()["blacklist"]		    
		    if ok && len(pblacklists[0]) >0 {
				pblacklist := pblacklists[0]
				fmt.Println("blacklist: " + string(pblacklist))			
				if(len(pblacklist)>0) {
					blackliststr=pblacklist
				}
			}
						
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			
					
			fmt.Fprintf(w, "<!doctype html><html lang='en'>")
			fmt.Fprintf(w, "<head>")
			fmt.Fprintf(w, "<meta charset='utf-8'>")
			fmt.Fprintf(w, "<meta name='viewport' content='width=device-width, initial-scale=1, shrink-to-fit=no'>")
			fmt.Fprintf(w, "<link rel='icon' href='data:,'>")
			fmt.Fprintf(w, "<link rel='stylesheet' href='https://cdn.jsdelivr.net/npm/bootstrap@4.5.3/dist/css/bootstrap.min.css' integrity='sha384-TX8t27EcRE3e/ihU7zmQxVncDAy5uIKz4rEkgIXeMed4M0jlfIDPvg6uqKI2xXr2' crossorigin='anonymous'>")
			fmt.Fprintf(w, "<link rel='stylesheet' href='https://cdn.jsdelivr.net/npm/bootstrap-icons@1.8.0/font/bootstrap-icons.css'>")

			fmt.Fprintf(w, "<style>.btn-badger { background-image: url('https://pgbadger.darold.net/logo_pgbadger.png') !important; background-size: contain; }</style>" )
			fmt.Fprintf(w, "<style>.btn-beaver { background-image: url('https://zcdn.uk/wp-content/uploads/2022/02/beaver.png') !important; background-size: contain; }</style>" )
			fmt.Fprintf(w, "<style>.i { font-size: 0.75rem; }</style>" )
			
			fmt.Fprintf(w, "</head>")
			fmt.Fprintf(w, "<body>")


			fmt.Fprintf(w, "<nav class='navbar sticky-top navbar-light bg-light'>")
			fmt.Fprintf(w, "<h3>pgBeaver</h3><br>")
			fmt.Fprintf(w, "</nav>")

			fmt.Fprintf(w, "<main class='container-fluid d-flex h-100 flex-column'>")
			
			// PARENT FILTER CRITERIA
			
			fmt.Fprintf(w, "<div class='row'>")			
			fmt.Fprintf(w, "<div class='input-group col-sm-3'><div class='input-group-prepend'><span class='input-group-text'>whitelist</span></div><input type='text' id='whitelist' class='form-control' value='%s' aria-label='Whitelist' aria-describedby='Whitelist'></div>", whiteliststr)
			fmt.Fprintf(w, "<div class='input-group col-sm-3'><div class='input-group-prepend'><span class='input-group-text'>blacklist</span></div><input type='text' id='blacklist' class='form-control' value='%s' aria-label='Blacklist' aria-describedby='Blacklist'></div>", blackliststr)
			fmt.Fprintf(w, "</div>")
			
			
			// LOG MENU
			
			fmt.Fprintf(w, "<menu class='row' style='max-height: 120px;overflow-y: scroll;'>")
    


			// LOG ENTRIES
			
			files, err := WalkMatch(logpath, logprefix+"*.log")
		    if err != nil {
				fmt.Fprintf(w, "No Log files found in %s<br>", logpath)
		    } else {
							
				
				for _, path := range files {
					logfile := strings.Replace(path, logpath, "", -1)
					htmfile := strings.Replace(logfile, ".log", ".htm", -1)
					logname := strings.Replace(logfile, ".log", "", -1)
					logid := strings.Replace(logname, logprefix, "", -1)

					fmt.Fprintf(w, "<span class='badge badge-info'>"+logid+"&nbsp;")				

					logF, err := os.Stat(logpath+logfile)
					if err != nil {
						//fmt.Println("File does not exist")
					}
															
					if(fileExists(logpath+htmfile)) {
						// Show existing log
						
						htmF, err := os.Stat(logpath+htmfile)
						if err != nil {
							//fmt.Println("File does not exist")
						}
						if( logF.ModTime().Unix()>htmF.ModTime().Unix() ) {
							// log updated since
							fmt.Fprintf(w, "<button title='View (Outdated) Already Generated Log' class='btn btn-badger btn-warning btn-sm' onclick='ifr.src=\"%s\";'><i class='bi-eye'></i></button>", "/show/"+htmfile)
							fmt.Fprintf(w, "<button title='Re-Generate and View Log' class='btn btn-badger btn-danger btn-sm' onclick='ifr.src=\"%s\";'><i class='bi-recycle'></i></button>", "/gen/"+htmfile)					
						} else {
							fmt.Fprintf(w, "<button title='View Already Generated Log' class='btn btn-badger btn-success btn-sm' onclick='ifr.src=\"%s\";'><i class='bi-eye'></i></button>", "/show/"+htmfile)
						} 						
					} else {
						// No log htm generated yet
						if( time.Now().Unix()-logF.ModTime().Unix() < 10) {
							// log might still being generated
							fmt.Fprintf(w, "<button title='Generate and View Log (Current Snapshot)' class='btn btn-badger btn-warning btn-sm' onclick='ifr.src=\"%s\";'><i class='bi-lightning-charge'></i></button>", "/gen/"+htmfile)
						} else {
							// log not recently updated
							fmt.Fprintf(w, "<button title='Generate and View Log' class='btn btn-badger btn-primary btn-sm' onclick='ifr.src=\"%s\";'><i class='bi-lightning-charge'></i></button>", "/gen/"+htmfile)
						}
					}
					fmt.Fprintf(w, "<button title='Download Log' class='btn btn-light btn-sm' onclick='ifr.src=\"%s\";'><i class='bi-file-earmark-arrow-down'></i></button>", "/log/"+logfile)
					fmt.Fprintf(w, "<button title='Scan Log' class='btn btn-beaver btn-success btn-sm' onclick='ifr.src=\"/message/?message=Generating_Report..\"; setTimeout(function(){ var w = document.getElementById(\"whitelist\").value; var b = document.getElementById(\"blacklist\").value; ifr.src=\"%s?whitelist=\"+w+\"&blacklist=\"+b;},100);'><i class='bi-list-check'></i></button>", "/scan/"+logfile)
					if( allowdelete || ((len(u)>0) && (len(p)>0)) ) {
						fmt.Fprintf(w, "<button title='Delete Log' class='btn btn-danger btn-sm' onclick='ifr.src=\"%s\"; setTimeout(() => { window.location.reload() }, 500);'><i class='bi-folder-x'></i></button>", "/del/"+htmfile)
					}	

					fmt.Fprintf(w, "</span>&nbsp;")
				}

			}

			fmt.Fprintf(w, "</menu>")
		
			
			// IFRAME FOR RESULTS
			
			fmt.Fprintf(w, "<iframe class='row flex-grow-1' id='ifr' src='/info/' title='Badger'></iframe>")
			
			fmt.Fprintf(w, "</main>")
			
			// RESIZE IFRAME ONLOAD
			
			fmt.Fprintf(w, "<script>var frame = document.getElementById('ifr'); ifr.onload = function(){ var h=ifr.contentWindow.document.body.scrollHeight; h=h+100; ifr.style.height=h+'px';}</script>")		
			
			fmt.Fprintf(w, "</body>")
			fmt.Fprintf(w, "<script src='https://code.jquery.com/jquery-3.5.1.min.js'></script>")
			fmt.Fprintf(w, "<script src='https://cdn.jsdelivr.net/npm/popper.js@1.16.1/dist/umd/popper.min.js'></script>")
			fmt.Fprintf(w, "<script src='https://cdn.jsdelivr.net/npm/bootstrap@4.5.3/dist/js/bootstrap.min.js'></script>")
			fmt.Fprintf(w, "</html>")
		
			return

		})





		http.HandleFunc("/scan/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")

			start := time.Now()
	
			fmt.Fprintf(w, "<!doctype html><html lang='en'>")
			fmt.Fprintf(w, "<head>")
			fmt.Fprintf(w, "<meta charset='utf-8'>")
			fmt.Fprintf(w, "<meta name='viewport' content='width=device-width, initial-scale=1, shrink-to-fit=no'>")
			fmt.Fprintf(w, "<link rel='icon' href='data:,'>")

			fmt.Fprintf(w, "<link rel='stylesheet' href='https://cdn.jsdelivr.net/npm/bootstrap@4.5.3/dist/css/bootstrap.min.css' integrity='sha384-TX8t27EcRE3e/ihU7zmQxVncDAy5uIKz4rEkgIXeMed4M0jlfIDPvg6uqKI2xXr2' crossorigin='anonymous'>")
			fmt.Fprintf(w, "<link rel='stylesheet' href='https://cdn.jsdelivr.net/npm/bootstrap-icons@1.8.0/font/bootstrap-icons.css'>")

			fmt.Fprintf(w, "<link href='https://cdnjs.cloudflare.com/ajax/libs/c3/0.4.24/c3.min.css' rel='stylesheet'>")
			fmt.Fprintf(w, "<script src='https://cdnjs.cloudflare.com/ajax/libs/d3/3.5.17/d3.min.js'></script>")			
			fmt.Fprintf(w, "<script src='https://cdnjs.cloudflare.com/ajax/libs/c3/0.4.24/c3.min.js'></script>")

			fmt.Fprintf(w, "<style>.c3-ygrid-line.majorline line { stroke: black; }</style>" )
			fmt.Fprintf(w, "<style>pre { color: blue; }</style>" )

			fmt.Fprintf(w, "</head>")
			fmt.Fprintf(w, "<body>")

			fmt.Fprintf(w, "<main class='container-fluid d-flex h-100 flex-column'>")
			fmt.Fprintf(w, "<h3>Log Scan</h3>")

		    logfile := strings.TrimPrefix(r.URL.Path, "/scan/")

			fmt.Printf("Scanning %s\n", logfile)

			if(debug) {
				PrintMemUsage()
			}
    
			// Needs optimising to save memory !!
			
		    log, _ := ioutil.ReadFile(logpath+logfile)		
							
			lines := strings.Split(string(log),"\n")
			
			lines = append(lines, "")

			// Destroy log variable	to free memory (> garbage collection)
			log = nil
			
			fmt.Printf("Scanning %s, %d lines\n", logfile, len(lines))

			if(debug) {
				PrintMemUsage()
			}
			
			const MINFACTOR = 1.3
			const MAXMS = 10000
			const STARTMS = 0.001
			const MAXSHOW = 1000
						
			var msv float64
			var durpos int
			var msi float64
			var durdesc1 string
			var durdesc2 string
			type dur struct {
			     desc string
			     count int
			     from float64
			     to float64
			}

			type logentry struct {
			     dur string
			     str string
			}

			var mydur dur	
			var durations = []dur{}
			var factor float64		
			var i int
	
			var sqlquery string
			var whitelist []string
			var blacklist []string
			
			l := 0
			phase := 0
			
			// Higher factors reduce resolution (number of bars on chart)
			//factor = 10						
			//factor = 2						
			//factor = math.Sqrt(10) // square root
			//factor = math.Cbrt(10) // cube root
			factor = math.Sqrt(math.Sqrt(10)) // 4th root

			blackliststr := ""
			whiteliststr := ""			
			showabovems := 1000.0
			showbelowms := 100000.0
			
			// PROCESS PARAMETERS
			
		    pfactors, ok := r.URL.Query()["factor"]		    
		    if ok && len(pfactors[0]) >0 {
				pfactor := pfactors[0]
				fmt.Println("factor: " + string(pfactor))
				if(isNumeric(pfactor)) {
					pfactorv , _ := strconv.ParseFloat(pfactor, 64)
					if (pfactorv > MINFACTOR) {
						factor = pfactorv
					}
				}
			}

		    pwhitelists, ok := r.URL.Query()["whitelist"]		    
		    if ok && len(pwhitelists[0]) >0 {
				pwhitelist := pwhitelists[0]
				fmt.Println("whitelist: " + string(pwhitelist))			
				if(len(pwhitelist)>0) {
					whiteliststr=pwhitelist
				}
			}
			whitelist = strings.Split(whiteliststr, ",")

		    pblacklists, ok := r.URL.Query()["blacklist"]		    
		    if ok && len(pblacklists[0]) >0 {
				pblacklist := pblacklists[0]
				fmt.Println("blacklist: " + string(pblacklist))			
				if(len(pblacklist)>0) {
					blackliststr=pblacklist
				}
			}
			blacklist = strings.Split(blackliststr, ",")

		    pshowabovemss, ok := r.URL.Query()["showabovems"]		    
		    if ok && len(pshowabovemss[0]) >0 {
				pshowabovems := pshowabovemss[0]
				fmt.Println("showabovems: " + string(pshowabovems))			
				if(isNumeric(pshowabovems)) {
					pshowabovemsv , _ := strconv.ParseFloat(pshowabovems, 64)
					if (pshowabovemsv > 0.00001) {
						showabovems=pshowabovemsv
					}
				}
			}
			
		    pshowbelowmss, ok := r.URL.Query()["showbelowms"]		    
		    if ok && len(pshowbelowmss[0]) >0 {
				pshowbelowms := pshowbelowmss[0]
				fmt.Println("showbelowms: " + string(pshowbelowms))			
				if(isNumeric(pshowbelowms)) {
					pshowbelowmsv , _ := strconv.ParseFloat(pshowbelowms, 64)
					if (pshowbelowmsv < 1000000) {
						showbelowms=pshowbelowmsv
					}
				}
			}			

			// SET UP ARRAY FOR TIME RANGES
			
			msi = STARTMS
			i = 0
			maxedout := false
			prevmsi := 0.0
			for ( maxedout == false ) {
				
				if( msi <= STARTMS ) {
					durdesc1 = "0"									
				} else {
					durdesc1 = fmt.Sprintf("%0.4f", prevmsi)									
				}		
				if( prevmsi <= MAXMS ) {
					durdesc2 = fmt.Sprintf("%0.4f", msi)									
				} else {
					durdesc2 = "infinity"									
				}				
				
				mydur.desc = durdesc1+"-"+durdesc2+"ms"
				mydur.count = 0
				mydur.from = prevmsi
				mydur.to = msi
				durations = append(durations, mydur)

				if( prevmsi > MAXMS )  {
					// +1 includes overflow
					maxedout = true
				}				
				prevmsi = msi
				msi = msi * factor
				i++
			}
			fmt.Printf("%d durations created\n", i)

			
			// CHECK EACH LOG ENTRY
						
			omitted := 0
			blackcount := 0
			whitecount := 0
			
			var entries []logentry
			var myentry logentry
			
			o := 0
			
			for l < len(lines) {
				line := lines[l]
				
				if(debug) {
					if((l % 1000) == 0) {
						// Show Progress on CLI
						// fmt.Printf(".")
					}
				}

				durpos = strings.Index(line, "duration:")
				if(durpos > 0 || (l == len(lines)-1)) {
					
					blackened := false
					whitened := false
					for _ , blackitem := range blacklist {
						if(strings.Index(sqlquery, blackitem)>0) {
							blackened = true
						}
					}
					for _ , whiteitem := range whitelist {
						if(strings.Index(sqlquery, whiteitem)>0) {
							whitened = true
						}
					}

					docount := false
					if(blackened) {
						blackcount++
						omitted++
					} else {
						if(i >= len(durations)) {
							// Probably bogus duration so just put to CLI for now - requires attention! 
							fmt.Printf("."+sqlquery+".")
						}
						if(whitened) {
							whitecount++
							docount = true
						}
					}
					
					if(docount) {
						if(i<len(durations)) {
							durations[i].count++
						}
						if(msv > showabovems) {
							if(msv < showbelowms) {
								// Enter into array
								if(len(entries)<MAXSHOW) {
									// Only add within limit to save memory
									myentry.dur = fmt.Sprintf("%0.3f ms",msv)
									myentry.str = fmt.Sprintf("%s", sqlquery)
									entries = append(entries, myentry)
								}
								o++
							}
						}
					}
						
					sqlquery = ""
					phase = 0
				}



				if(phase == 0) {
				// PHASE 0

					
					if(len(line)>26) {
						// Must be longer than date stamp
						
						// Should really check date but not sure about different formats, so leaving off for now
						// chkdate := line[0:26]
						// _ , err := time.Parse(chkdate,"2022-02-10 17:05:08.423 UTC")
						
						//if err == nil { // date check
							logdata := line[26:]
							durpos = strings.Index(logdata, "duration:")
							if(durpos > 0) {
								durstr := logdata[durpos+10:]
								ms := strings.Replace(durstr, "ms", "", -1)
								ms = strings.Replace(ms, " ", "", -1)
								msv , _ = strconv.ParseFloat(ms, 64)
								msi = STARTMS
								i = 0
								for ( msi < MAXMS ) && ( msi < msv ) && (i < len(durations))  {
									msi = msi * factor
									i++
								}

								phase++
								sqlquery = ""
							}
						//} // omitted date check
					}	
				} else if(phase == 1) {
				// PHASE 1

					sqlquery = sqlquery+" "+line

				}
				
				l++
			}
			
			// Destroy lines			
			lines = nil			
						
			if(debug) {
				PrintMemUsage()
			}
			
			fmt.Fprintf(w, "<br><br>")


			// CHART
			
			fmt.Fprintf(w, "<div id='chart' class='c3' style='max-height: 280px; position: relative;'>")
			
			cols := ""
			msss := ""
			vals := ""
			logvals := ""
			froms := ""
			tos := ""
			labels := ""

			for iss , thisdur := range durations {
				msss = msss+fmt.Sprintf(", %d", iss)
				vals = vals+fmt.Sprintf(", %d", thisdur.count)
				if(len(labels)>0) { labels = labels+"," }
				labels = labels + fmt.Sprintf("%d", thisdur.count)
				if(thisdur.count == 0) {
					logvals = logvals+", 0"
				} else {
					logvals = logvals+fmt.Sprintf(", %f", math.Log10(float64(thisdur.count))+1)
				}
				if(len(cols)>0) { cols = cols+"," }
				cols = cols + fmt.Sprintf("'%s'", thisdur.desc)
				if(len(froms)>0) { froms = froms+"," }
				froms = froms+fmt.Sprintf("%f", thisdur.from)
				if(len(tos)>0) { tos = tos+"," }
				tos = tos+fmt.Sprintf("%f", thisdur.to)
			}

			fmt.Fprintf(w, "<script>")

			fmt.Fprintf(w, "var cols = ["+cols+"];")
			fmt.Fprintf(w, "var froms = ["+froms+"];")
			fmt.Fprintf(w, "var tos = ["+tos+"];")
			fmt.Fprintf(w, "var labels = ["+labels+"];")

			fmt.Fprintf(w, "var chart = c3.generate({")
			fmt.Fprintf(w, "data: {")
			
			
			
			fmt.Fprintf(w, "labels: {")
			fmt.Fprintf(w, "format: {")

			fmt.Fprintf(w, "count: function (v, id, i, j) { return labels[i] }")

			fmt.Fprintf(w, "}")
			fmt.Fprintf(w, "},")
			
  
    
        
 
			fmt.Fprintf(w, "type: 'bar',")
			fmt.Fprintf(w, " columns: [")

			fmt.Fprintf(w, "['count'"+logvals+"]")

			fmt.Fprintf(w, "]")
			
			fmt.Fprintf(w, ",onclick: function (d) { var f = document.getElementById(\"showabovems\"); f.value = froms[d.index]; var domax = document.getElementById('domax').checked; if(domax) { var t = document.getElementById(\"showbelowms\"); t.value = tos[d.index]; } }")
			
			fmt.Fprintf(w, "}")


			fmt.Fprintf(w, ",axis: {")
			fmt.Fprintf(w, "x: {")
			fmt.Fprintf(w, "type: 'category',")
			fmt.Fprintf(w, "categories: ["+cols+"]")
			fmt.Fprintf(w, " }")
			
			fmt.Fprintf(w, ", y: {")
//			fmt.Fprintf(w, "tick: { format: function (d) { return Math.pow(10,d.value).toFixed(2); } }")
			fmt.Fprintf(w, "tick: { format: function (v, id, i, j) { return i; } }")
			fmt.Fprintf(w, " }")
            			
			fmt.Fprintf(w, " }")	

			
			fmt.Fprintf(w, ",grid: {")
			fmt.Fprintf(w, "y: {")
			fmt.Fprintf(w, "lines:[")
			
			// WIP CALCULATE LOG VALS
/*
			var h float64
			var y float64
			h = 0
			y = 1
			for h <(3*(MAXMS / STARTMS)) {
				if(h>0) {
					fmt.Fprintf(w, ",")
				}
				y = 1+((h/3)*math.Cbrt(10))
				fmt.Fprintf(w, "{value: %f, text: '%f'}", y, math.Pow(10,y))				
				h++
			}
*/			

			// MANUAL IMPLEMENTATION !

			// plot as log10 ratios
			// 1 = 10^0 but origin at 1 so as 1 has enough clickable area
			fmt.Fprintf(w, "{value: 1,				text: '1',			class: 'majorline'},")
			fmt.Fprintf(w, "{value: 1.215443469,	text: ''},")
			fmt.Fprintf(w, "{value: 1.464158883,	text: '5'},")
			fmt.Fprintf(w, "{value: 2, 				text: '10',			class: 'majorline'},")
			fmt.Fprintf(w, "{value: 2.215443469,	text: ''},")
			fmt.Fprintf(w, "{value: 2.464158883,	text: '46'},")
			fmt.Fprintf(w, "{value: 3,				text: '100',		class: 'majorline'},")
			fmt.Fprintf(w, "{value: 3.215443469,	text: ''},")
			fmt.Fprintf(w, "{value: 3.464158883,	text: '464'},")
			fmt.Fprintf(w, "{value: 4, 				text: '1,000',		class: 'majorline'},")
			fmt.Fprintf(w, "{value: 4.215443469,	text: ''},")
			fmt.Fprintf(w, "{value: 4.464158883,	text: '4,642'},")
			fmt.Fprintf(w, "{value: 5,				text: '10,000',		class: 'majorline'},")
			fmt.Fprintf(w, "{value: 5.215443469,	text: ''},")
			fmt.Fprintf(w, "{value: 5.464158883,	text: '46,416'},")
			fmt.Fprintf(w, "{value: 6,				text: '100,000',	class: 'majorline'},")
			fmt.Fprintf(w, "{value: 6.215443469,	text: ''},")
			fmt.Fprintf(w, "{value: 6.464158883,	text: '464,159'},")
			fmt.Fprintf(w, "{value: 7,				text: '1,000,000',	class: 'majorline'},")
			fmt.Fprintf(w, "{value: 7.215443469,	text: ''},")
			fmt.Fprintf(w, "{value: 7.464158883,	text: '4,641,589'},")
			fmt.Fprintf(w, "{value: 8,				text: '10,000,000',	class: 'majorline'}")
			
			fmt.Fprintf(w, "]")
			fmt.Fprintf(w, "}")
			fmt.Fprintf(w, "}")
			              
			
			fmt.Fprintf(w, "});")
			fmt.Fprintf(w, "</script>")        
       
            
            
			// MENU
			
			fmt.Fprintf(w, "<div class='row'>")			
			fmt.Fprintf(w, "<div class='input-group col-sm-3'><div class='input-group-prepend'><span class='input-group-text'>Whitelist</span></div><input type='text' id='whitelist' class='form-control' value='"+whiteliststr+"' aria-label='Whitelist' aria-describedby='Whitelist'></div>")
			fmt.Fprintf(w, "<div class='input-group col-sm-3'><div class='input-group-prepend'><span class='input-group-text'>Blacklist</span></div><input type='text' id='blacklist' class='form-control' value='"+blackliststr+"' aria-label='Blacklist' aria-describedby='Blacklist'></div>")
			fmt.Fprintf(w, "<div id='minms' class='input-group col-sm-2'><div class='input-group-prepend'><span class='input-group-text'>Show Abo/ve (ms)</span></div><input type='number' id='showabovems' class='form-control' value='%f' aria-label='Above (ms)' aria-describedby='Above (ms)'></div>", showabovems)
			fmt.Fprintf(w, "<div id='max' class='input-group col-sm-1'><div class='input-group-prepend'><span class='input-group-text'>Change Max?</div><input type='checkbox' id='domax' class='form-control' value='checked' aria-label='Change Max' aria-describedby='Change Max'></div>")
			fmt.Fprintf(w, "<div id='maxms' class='input-group col-sm-2'><div class='input-group-prepend'><span class='input-group-text'>Show Below (ms)</span></div><input type='number' id='showbelowms' class='form-control' value='%f' aria-label='Below (ms)' aria-describedby='Below (ms)'></div>", showbelowms)
			fmt.Fprintf(w, "<button type='submit' class='btn btn-primary sm-1' onclick='var w = document.getElementById(\"whitelist\").value; var b = document.getElementById(\"blacklist\").value; var a = document.getElementById(\"showabovems\").value; var x = document.getElementById(\"showbelowms\").value; window.location.href=\"%s?whitelist=\"+w+\"&blacklist=\"+b+\"&showabovems=\"+a+\"&showbelowms=\"+x;'>Filter Queries</button>", "/scan/"+logfile)
			fmt.Fprintf(w, "</div>")



			// RESULT SUMMARY
			
			
			fmt.Fprintf(w, "<div class='row'>")
			fmt.Fprintf(w, "<div class='input-group col-sm-2'>Filtered Results = %d</div>", o)
			fmt.Fprintf(w, "<div class='input-group col-sm-2'>Total Omitted = %d</div>", omitted)
			fmt.Fprintf(w, "<div class='input-group col-sm-2'>Blacklisted = %d</div>", blackcount)
			fmt.Fprintf(w, "<div class='input-group col-sm-2'>Whitelisted = %d</div>", whitecount)
			fmt.Fprintf(w, "</div>")

			if(o >= MAXSHOW) {
				fmt.Fprintf(w, "<div class='alert alert-danger' role='alert'>Max limit Exceeded! Output Truncated</div>")
			}
 			
   			s := 0
   			for _ , entry := range entries {
				if(s < MAXSHOW) {
					fmt.Fprintf(w, "<div class='row'><div class='input-group col-sm-2'>%s</div><div class='input-group col-sm-10'><pre>%s</pre></div></div>", entry.dur, entry.str)
				}
				s++
			}
 			
			fmt.Fprintf(w, "</main>")
			
			fmt.Fprintf(w, "<script>var frame = document.getElementById('ifr');ifr.onload = function(){ var h=ifr.contentWindow.document.body.scrollHeight; h=h+100; ifr.style.height=h+'px';}</script>")		
			
			fmt.Fprintf(w, "</body>")
			fmt.Fprintf(w, "<script src='https://code.jquery.com/jquery-3.5.1.min.js'></script>")
			fmt.Fprintf(w, "<script src='https://cdn.jsdelivr.net/npm/popper.js@1.16.1/dist/umd/popper.min.js'></script>")
			fmt.Fprintf(w, "<script src='https://cdn.jsdelivr.net/npm/bootstrap@4.5.3/dist/js/bootstrap.min.js'></script>")
			fmt.Fprintf(w, "</html>")
			
			elapsed := time.Since(start)
			fmt.Printf("Took %s\n\n", elapsed)								

			runtime.GC()
			return
		})




		http.HandleFunc("/gen/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")				    
		    htmfile := strings.TrimPrefix(r.URL.Path, "/gen/")
			logfile := strings.Replace(htmfile, ".htm", ".log", -1)
			
			fmt.Fprintf(w, "Generating Output..<br>")

			fmt.Fprintf(w, logfile)
			cmd := exec.Command("./pgbadger", logpath+logfile, "-o", logpath+htmfile)
		    _ = cmd.Run()

			fmt.Fprintf(w, "<br>Refreshing Pane..")

			fmt.Fprintf(w, "<script>setTimeout(() => { window.location.href='http://golem.zcnhosts.com:5045/show/" + htmfile + "' }, 2000);</script>")

			return

		})



		http.HandleFunc("/del/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")				    
		    htmfile := strings.TrimPrefix(r.URL.Path, "/del/")
			logfile := strings.Replace(htmfile, ".htm", ".log", -1)
			
			fmt.Fprintf(w, "Deleting File..<br>")

			os.Remove(logpath+logfile)

			return

		})

		http.HandleFunc("/show/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")				    
		    logfile := strings.TrimPrefix(r.URL.Path, "/show/")
		    htm, err := ioutil.ReadFile(logpath+logfile)		
		    if err != nil {
				fmt.Fprintf(w, "ERR")
		    } else {
				fmt.Fprintf(w, "%s", htm)
			}

			return

		})

		http.HandleFunc("/log/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/octet-stream")				    
		    logfile := strings.TrimPrefix(r.URL.Path, "/log/")
		    log, err := ioutil.ReadFile(logpath+logfile)		
		    if err != nil {
				fmt.Fprintf(w, "ERR")
		    } else {
				fmt.Fprintf(w, "%s", log)
			}

			return

		})
		
		// Allow user to specify port    
		var port string
		flag.StringVar(&port, "port", "5045", "Port Number (default 5045)")
		
		// Allow user to specify http/https protocol    
		var proto string
		flag.StringVar(&proto, "proto", "http", "http | https (default http)")

		// User Auth    
		flag.StringVar(&username, "username", "", "Username")

		// Password    
		flag.StringVar(&password, "password", "", "Password")

		// Log Path    
		flag.StringVar(&logpath, "logpath", "/var/log/postgresql/", "Log Path")

		// Debug    
		debugptr := flag.Bool("debug", false, "Debug true/false")

		// Allow Delete    
		allowdeleteptr := flag.Bool("allowdelete", false, "Allow Log File Deletion true/false")
		
		// If https protocol selected, must specify certpath & keypath
		var certpath string
		var keypath string
		flag.StringVar(&certpath, "certpath", "/etc/ssl/certs/certificate.crt", "path to certificate (if https proto specified)")
		flag.StringVar(&keypath, "keypath", "/etc/ssl/private/certificate.key", "path to key (if https proto specified)")
					
		flag.Parse()
		
		// eval to global var
		debug = *debugptr
		allowdelete = *allowdeleteptr
		
		// Advise listening on port
		fmt.Println("Listening on port: ", port, ", username: ", username, ", password: ", password)
		if(proto=="http") {
			log.Fatal(http.ListenAndServe(":"+port, nil))
		} 
		if(proto=="https") {
			log.Fatal(http.ListenAndServeTLS(":"+port, certpath, keypath, nil))
		}

}

