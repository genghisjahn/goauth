goauth
======


Decided to write this after reading:  [http://www.thebuzzmedia.com/designing-a-secure-rest-api-without-oauth-authentication/](http://www.thebuzzmedia.com/designing-a-secure-rest-api-without-oauth-authentication/).

To run locally, first:

**Assuming OSX**  
**To run all three pieces locally**  
1. Make sure you have Go Installed and setup correctly: http://golang.org/doc/install  
2. Git clone to a src directory `git clone git@github.com:genghisjahn/goauth.git`  
3. Get your local IP Address from the terminal run `ifconfig` and looks for the IP address next to `inet` near the top.  
4. Create a virtual IP address on your local machine.  You can find instructions for that [here](http://gerrydevstory.com/2012/08/20/how-to-create-virtual-network-interface-on-mac-os-x/) (There isn't much to it).  
5. *Sudo* edit your hosts file located at /private/etc/hosts and add a line that points `www.order-demo.com` to the virtual IP address from step 4.  
6. cd to the `gaserver` directory and run `go run gaserver.go -http [IP_From_Step_3]:8090`  
7. cd to the `mim` directory and run `go run main.go -outhttp [IP_From_Step_4]:8090`  
8. cd to the `gademo` director and run `go run main.go`.  If the previous steps have been done correctly, you won't need to specify any command line arguments for this command.  
9. Open a browser and naviagte to http://localhost:8080.
