ServiceCommunicator
===================

A command serializer/deserializer for a remote service, using an adapted version of the Redis protocol.

    package main
    
    import (
      "bufio"
    	"fmt"
    	"github.com/chaosphere2112/ServiceCommunicator/servicecomm"
    	"os"
    )
    
    func main() {
    
    	fmt.Println("Enter your message:")
    
    	read := bufio.NewReader(os.Stdin)
    	str, _ := read.ReadString('\n')
    
    	message := servicecomm.EncodeMessage(str[0 : len(str)-1])
    
    	serv := servicecomm.NewDecoder()
    
    	serv.DecodeMessage(message)
    
    	if serv.LastMessage != nil {
    		fmt.Println(serv.LastMessage)
    	} else {
    		fmt.Println("Message was invalid.")
    	}
    }
