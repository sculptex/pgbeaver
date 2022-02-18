# pgbeaver
Postgres Log Analyser Tool

* Allows browsing of log files via browser
* Includes Native tool to graphing and drilling down by query times 
* Incorporates browsing log files and generating **pgbadger** (https://github.com/darold/pgbadger) log analysis via browser

![Screen Shot](https://github.com/sculptex/pgbeaver/blob/main/pgbeaver_screenshot1.png)

## Command Line Usage

###  -certpath string
        path to certificate (if https proto specified) (default "/etc/ssl/certs/certificate.crt")
###  -debug
        Debug true/false
###  -keypath string
        path to key (if https proto specified) (default "/etc/ssl/private/certificate.key")
###  -logpath string
        Log Path (default "/var/log/postgresql/")
###  -password string
        Password
###  -port string
        Port Number (default 5045) (default "5045")
###  -proto string
        http | https (default http) (default "http")
###  -username string
        Username

## Browser Useage
  http://localhost:5045/?username=user&password=pass&whitelist=reference_objects
