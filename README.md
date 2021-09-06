# trinity-micro

```
cd example/
go run main.go api
```



```
% wrk -t12 -c100 -d10s http://127.0.0.1:3000/benchmark/simple_raw
Running 10s test @ http://127.0.0.1:3000/benchmark/simple_raw
  12 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   781.37us  442.87us   7.51ms   77.68%
    Req/Sec    10.53k   709.86    12.45k    65.98%
  1268705 requests in 10.10s, 181.49MB read
Requests/sec: 125611.65
Transfer/sec:     17.97MB


%  wrk -t12 -c100 -d10s http://127.0.0.1:3000/benchmark/simple          
Running 10s test @ http://127.0.0.1:3000/benchmark/simple
  12 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     0.94ms  848.65us  15.40ms   85.68%
    Req/Sec     9.83k     1.12k   12.75k    74.26%
  1184808 requests in 10.10s, 169.49MB read
Requests/sec: 117297.24
Transfer/sec:     16.78MB
```


```
%  wrk -t12 -c100 -d10s http://127.0.0.1:3000/benchmark/path_param_raw/1
Running 10s test @ http://127.0.0.1:3000/benchmark/path_param_raw/1
  12 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   795.47us  487.06us  11.88ms   78.66%
    Req/Sec    10.44k     0.87k   12.67k    64.60%
  1259132 requests in 10.10s, 176.52MB read
Requests/sec: 124669.29
Transfer/sec:     17.48MB

%  wrk -t12 -c100 -d10s http://127.0.0.1:3000/benchmark/path_param/1    
Running 10s test @ http://127.0.0.1:3000/benchmark/path_param/1
  12 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   835.88us  701.02us  22.89ms   90.41%
    Req/Sec    10.11k   646.21    12.74k    73.33%
  1207872 requests in 10.01s, 202.74MB read
  Non-2xx or 3xx responses: 1207872
Requests/sec: 120619.89
Transfer/sec:     20.25MB
```