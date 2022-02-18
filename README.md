# pgbeaver
Postgres Log Analyser Tool

* Allows browsing of log files via browser
* Includes Native tool to graphing and drilling down by query times 
* Incorporates browsing log files and generating **pgbadger** (https://github.com/darold/pgbadger) log analysis via browser (requires pgbadger in same folder as executable)

![Screen Shot](https://github.com/sculptex/pgbeaver/blob/main/pgbeaver_screenshot1.png)

## Command Line Usage

####  -logpath string
        Log Path (default "/var/log/postgresql/")
####  -port string
        Port Number (default 5045) (default "5045")

### FOR CLI DEBUG INFO
####  -debug
        Debug true/false

### FOR USER AUTHENTICATION
####  -username string
        Username
####  -password string
        Password

### FOR SSL
####  -proto string
        http | https (default http) (default "http")
####  -certpath string
        path to certificate (if https proto specified) (default "/etc/ssl/certs/certificate.crt")
####  -keypath string
        path to key (if https proto specified) (default "/etc/ssl/private/certificate.key")

## Browser Usage
  http://localhost:5045/?username=user&password=pass&whitelist=reference_objects
