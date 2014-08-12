goauth
======

##What is this?

Goauth is a **for-learning-purposes-only** project that demonstrates a public/private key authentication system, various Man in the Middle attacks, and how the `gaserver` can determine when a request/order is legitimate or has been tampered with in some way.  The system is designed to (hopefully) authenticate a `user` and guard against replay attacks, delay attacks and tampering attacks(changing the order, the URL the http method/verb). 


##Install & Setup  
**Assuming OSX**  

**To run all three pieces locally**  
1. Make sure you have Go Installed and setup correctly: http://golang.org/doc/install  
2. Git clone to a src directory `git clone git@github.com:genghisjahn/goauth.git`  
3. Get your local IP Address from the terminal run `ifconfig` and looks for the IP address next to `inet` near the top.  
4. Create a virtual IP address on your local machine.  You can find instructions for that [here](http://gerrydevstory.com/2012/08/20/how-to-create-virtual-network-interface-on-mac-os-x/) (There isn't much to it).  
5. *Sudo* edit your hosts file located at /private/etc/hosts and add a line that points `www.order-demo.com` to the virtual IP address from step 4.  
6. cd to the `gaserver` directory and run `go run gaserver.go -http [IP_From_Step_3]:8090`  
7. cd to the `mim` directory and run `go run main.go -inhttp [IP_From_Step_4]:8090 -outhttp [IP_From_Step_3]:8090`  
8. cd to the `gademo` director and run `go run main.go`.  If the previous steps have been done correctly, you won't need to specify any command line arguments for this command.  
9. Open a browser and naviagte to [http://localhost:8080](http://localhost:8080).  
10. Enter an integer for `No. Shares` and `Max Price` and then click `Save`
11. 

## Attack!  

1. From the `mim` directory, running `go run main.go -inhttp IPAddr:Port -outhttp:Port` will run without any attacks.  The Man in the Middle portion just passes the request to the real server and returns it without any tampering.  This should result with a Success/200
2. Running `go run main.go -inhttp IPAddr:Port -outhttp:Port -attack repeat` will attempt to process the same order twice.  The web site [http://localhost:8080](http://localhost:8080) will report a success, because the first order will go through but the `mim.go` and the `gaserver.go` will log `Duplicate Nonce` errors.  The second order is ignored.
3. Valid `attack` values are:  
  1.  `none`  The default
  2.  `changeorder` Changes the value of the NumShares and MaxPrice properties.
  3.  `repeat` Attempts to run the same order twice.
  4.  `delay` Attemps to delay the order and the process at a later time.
  5.  `delayrepeat` Attemps to run an order and then run the same order after a small delay
  6.  `changeurl` Attemps to run the order against a different URL, say from `order-demo.com` to `order-live.com`
  7.  `invalid` Passes invalid json.




Decided to write this after reading:  [http://www.thebuzzmedia.com/designing-a-secure-rest-api-without-oauth-authentication/](http://www.thebuzzmedia.com/designing-a-secure-rest-api-without-oauth-authentication/).
