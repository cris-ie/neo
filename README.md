# neo

Chart based application to query a weeks worth ov nasas neo api

# Requirements 
- linux (tested with ubuntu 20.10)
- build-essentials
- docker
- k3s
- helm

# Installation

1. Copy values.yaml to ./chart/
2. Run make deploy in .

The service should now be reachable under http://neo.127.0.0.1.nip.io

# Endpoints

http://neo.127.0.0.1.nip.io/status - reports that the server is ready to serve and a db connection is established
http://neo.127.0.0.1.nip.io/liveness - reports that the server is ready to serve html files (a db connection might not yet be established or guaranteed)
http://neo.127.0.0.1.nip.io/neo/week - reports the number of NEOs for the next week
http://neo.127.0.0.1.nip.io/neo/next - reports the next NEO (Optional Query Parameter: ?hazardous=true - if set to true the next NEO that has the flag IsPotentiallyHazardousAstroid set to true is returned)