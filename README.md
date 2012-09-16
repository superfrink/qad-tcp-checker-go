qad-tcp-checker-go
==================

Quick and dirty http server connection status checker in Go

== Background ==

Nodes in a group of load balanced web servers were taking turns going down.
They were taking too long to respond to TCP connections.

I didn't have any graphing or monitoring.  I quickly wrote something in Perl
to tell me when a server was taking too long to respond.  Then I could login
and debug.

The Perl code sequentially checked each node and after checking all nodes
slept for a second before trying again.  With 12 nodes there could be more
than 15 seconds between when a host was checked.

I wrote this Go version to remove the sequential limitation in my previous
program.  I thought about using Perl threads but that may not work with the
USR1 signal.  I thought about multiplexing I/O.  Is there a way to multiplex
the connect() call?
