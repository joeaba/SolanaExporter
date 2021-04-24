# solana_exporter

solana_exporter exports basic monitoring data from a Solana node. For getting the data from the node, You can use ip addess of any machine that is running Solana RPC API.
The node gives the data in JSON format. In client.go file, Marshaling the go request to json. 

# Overview
We are just requesting the data from the node in JSON format and response will also be shown in JSON. Whatever the data in the response body comes we are converting it to prometheus format.

For each of the api requests we introduced different .go files each of having a function, which are called from exporter.go & slots.go(which are the part of main file).

# We are converting the JSON response to prometheous in two ways
1. exporter.go 
This file is having different collect functions which will be called internally. Each of the Collect functions is getting the response and send the body to the "must**Metric" function which is having a channel making unnderstandable to prometheus. All the collectors have to register to prometheus.  
2. slots.go
This file is containing the different metrics. All the metrics have to register in prometheus. Many of the api methods are called within this file and setting the data in prometheus format. We are not using channels as in exporter.go.

# Commands to start the project
1. go run exporter.go slots.go -rpcURI=http://<ip-address> -v=2 (From this command we can see all the api responses in JSON Format).
2. http://localhost:8080/metrics (From this command we can see all the data from apis in prometheous format).
