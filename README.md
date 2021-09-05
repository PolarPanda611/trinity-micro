# trinity-micro

```
cd example/
go run main.go api
```



```
% wrk -t12 -c100 -d30s http://127.0.0.1:3000/test/user2
Running 30s test @ http://127.0.0.1:3000/test/user2
  12 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.01ms    1.94ms  36.12ms   87.43%
    Req/Sec     4.62k   607.41     6.44k    71.31%
  1656779 requests in 30.01s, 232.26MB read
Requests/sec:  55200.45
Transfer/sec:      7.74MB

% wrk -t12 -c100 -d30s http://127.0.0.1:3000/test/user 
Running 30s test @ http://127.0.0.1:3000/test/user
  12 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.59ms    3.26ms  59.83ms   86.39%
    Req/Sec     4.87k   577.65     6.74k    69.33%
  1746373 requests in 30.01s, 229.84MB read
Requests/sec:  58188.78
Transfer/sec:      7.66MB
```


```
% wrk -t12 -c100 -d30s http://127.0.0.1:3000/test/user2/1
Running 30s test @ http://127.0.0.1:3000/test/user2/1
  12 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   824.07us  471.40us  20.96ms   81.33%
    Req/Sec     9.97k   830.33    13.31k    69.35%
  3584082 requests in 30.10s, 533.22MB read
Requests/sec: 119068.61
Transfer/sec:     17.71MB

% wrk -t12 -c100 -d30s http://127.0.0.1:3000/test/user/1 
Running 30s test @ http://127.0.0.1:3000/test/user/1
  12 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     0.93ms  797.68us  20.86ms   84.70%
    Req/Sec     9.61k   633.72    12.09k    70.32%
  3452473 requests in 30.10s, 484.00MB read
Requests/sec: 114694.49
Transfer/sec:     16.08MB
```