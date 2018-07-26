# Question
```
The original fizz-buzz consists in writing all numbers from 1 to 100, 
and just replacing all multiples of 3 by “fizz”, all multiples of 5 by “buzz”, 
and all multiples of 15 by “fizzbuzz”. The output would look like this:
“1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,fizz,...”
```

For this exemple, the output  must be 
"1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,..."
And this is what I did in my code.

Unless, if N is multiples by int1 and int2, you want to do 3 time replacing
ex: 15 is multiples by 3 and 5, so 15 is replaced to "fizzbuzz","fizz","buzz"

Using the same example above, the output with corresponding original value is:  
```
1 | 2 | 3    | 4 | 5    | 6    | 7 | 8 | 9    | 10   | 11 | 12   | 13 | 14 | 15                 | 16 | ...  
1 | 2 | fizz | 4 | buzz | fizz | 7 | 8 | fizz | buzz | 11 | fizz | 13 | 14 | fizzbuzz fizz buzz | 16 | ...  
```

I realize it with the func processRequestBis() 

# ENDPOINTS 
* endpoint client /run for replcacing fizz-buzz  
* endpoint server /heartbeats to check if server is running  
* endpoint ADMIN /stat to display server's stat information (start time, the number of received client request etc.)  
* endpoint ADMIN /update_conf to update server's configuration  

# FEATURE
* AUTH for endpoint ADMIN  
* manage http code to return to client (400 Bad Request, 401 Unauthorized, 403 Forbidden etc.)  
* manage server's log and the message to return to client  
* cunstom parameters (flag) when start the programe. ex: port  
* update server's configuration. ex: password for ADMIN, option of gzip for the response to client etc.  
* encoding gzip for huge response  
* manage clients' parallel requests by sync.Mutex  
* unitary test   
* check client's request. ex: if necessary fileds are missing  
* check and validate configuration information when update it

# RUN 
* build   
    ```$> go build```
* execute  
    ```$> ./fizzbuzz```
* test  
    ```$> go test```
* example client   
    ```
    $> curl -H "Content-Type: application/json" -X POST -d '{"string1":"yes","string2":"no","int1":5,"int2":8,"limit":100,"response_gzip":false}' http://localhost:8080/run
    ```
* example update configuration   
    ```
    $> curl -H "Content-Type: application/json" -X POST -d '{"password":"SUNXIN","force_gzip":false,"force_gzip_num":0}' http://localhost:8080/update_conf
    ```
* example display server's stat
    ```
    $> curl http://localhost:8080/stat?key=SUNXIN
    ```

# USER TEST
 * path: /client_demo
 * restart the server in order to use default parameters
 * test
    ```
    $> go build
    $> ./fizzbuzz
    $> cd client_demo
    $> go run main.go
    ```
 
# TODO
* When the size of HistoryCall reaches its limit, we should push the information to a DB.  
Or we can create a synchronisation system, update to a DB for each update (each new client's call).  

* I use a simple security to manage ACCESS ONLY ADMIN  
It must be more secure, especially when update server's configuration, we must not send PASSWORD clearly.  
ex: public_key, private_key
 