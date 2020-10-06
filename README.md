# TCP_Proxy
tcp proxy server

Pre-requisites:
Need arrangodb and gorillamux golang projects.

Usage:
go run main.go <IP:PORT>

use the above given IP:PORT for Rest APIs to control the tcp proxy server.

Available APIs

GET   showall     --> shows all proxy connections

GET   show/<ID>   --> shows proxy connection of that ID
  
POST  update/<ID> --> updates proxy connection of that ID
  
POST  delete/<ID> --> deletes proxy connection of that ID
  
