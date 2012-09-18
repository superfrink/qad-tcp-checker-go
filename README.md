qad-tcp-checker-go
==================

Quick and dirty http server connection status checker in Go

== Background ==

Nodes in a group of load balanced web servers were taking turns going down.
They were taking too long to respond to TCP connections.

I didn't have any graphing or monitoring.  I quickly wrote something in Perl
to tell me when a server was taking too long to respond so I could login and
debug.  See https://github.com/superfrink/qad-tcp-checker

The Perl code sequentially checked each node and after checking all nodes
slept for a second before trying again.  With 12 nodes there could be more
than 15 seconds between when a host was checked.

I wrote this Go version to remove the sequential limitation in my previous
program.  Goroutines seemed to be a natural solution.  I thought about using
Perl threads but that may not work with the USR1 signal.  I thought about
multiplexing I/O but is there a way to multiplex the connect() call?

This program prints a warning when connecting to a server does not complete in
1 second.  Pressing the Enter key while the program is running will dump a
summary of the connection statistics.

== Sample output ==

  $ go run qad-tcp-checker.go  
  .test failure  
  127.0.0.1 failure
  .a failure  
  ..127.0.0.1 failure  
  a failure  
  127.0.0.1 failure  
  a failure  
  127.0.0.1 failure  
    
  map[test:map[failure:1 success:2] a:map[failure:3] superfrink.net:map[success:3] 127.0.0.1:map[failure:4]]  
              test       1 failures       2 successes  
                 a       3 failures       0 successes  
    superfrink.net       0 failures       3 successes  
         127.0.0.1       4 failures       0 successes  
  a failure  
  127.0.0.1 failure  
  a failure  
  127.0.0.1 failure  

== Program archtecture ==

1) The aggregator reads server status data sent to a channel by the
   check_host_list function.

   When a message arrives on a second channel the aggregator prints the server
   statistics to the console when a message is received.

2) The check_host_list is a fan-in that calls check_host on each host.

3) The check_host checks a single host, writes the status to a channel, and
   sleeps for 1 to 2 seconds before checking again.

4) The main function starts the other components and then waits.  If the user
   presses the Enter key a message is sent to the aggregator asking it to show
   the server statistics.