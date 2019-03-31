# gocron environment variables

Docker compose expects an environment file with the following variables:

```
GC_DBFQDN=<your database fqdn or ip address>
GC_DBPORT=5432
GC_DBUSER=gocron
GC_DBPASS=password
GC_DBDATABASE=gocron
GC_INTERVAL=20
GC_SLACKHOOKURL=<your webhook url>
GC_SLACKCHANNEL=<your slack channel for alerts>
```
