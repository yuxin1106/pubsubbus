package main 

import (
		"fmt"
		"os"
		"time"
		"io/ioutil"
		"log"

		"nanomsg.org/go/mangos/v2"
		"nanomsg.org/go/mangos/v2/protocol/pub"
		"nanomsg.org/go/mangos/v2/protocol/sub"
		"nanomsg.org/go/mangos/v2/protocol/bus"

	    // register transports
		_"nanomsg.org/go/mangos/v2/transport/all"
)

// global variable to hold clientMsg received by server 
//var clientMsgRecv

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}


// client to mulitple server communication
func client(url string, clientMsg string) {
	var sock mangos.Socket
	var err error
	//var clientMsg = "Hello World"
	
	if sock, err = pub.NewSocket(); err != nil {
		die("can't get new pub socket: %s", err)
	}
	if err = sock.Listen(url); err != nil {
		die("can't listen on pub socket: %s", err.Error())
	}

	for {
		// Could also use sock.RecvMsg to get header
		fmt.Printf("CLIENT: PUBLISHING clientMsg: " + clientMsg + "\n")
		if err = sock.Send([]byte(clientMsg)); err != nil {
			die("Failed publishing: %s", err.Error())
		}
		time.Sleep(time.Second)
	}

}


func server(url string, name string) {
	var sock mangos.Socket
	var err error
	var msg []byte

	if sock, err = sub.NewSocket(); err != nil {
		die("can't get new sub socket: %s", err.Error())
	}
    if err = sock.Dial(url); err != nil {
    	die("can't dial on sub socket: %s", err.Error())
    }
    // Empty byte array effectively subscribes to everything
    err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
    if err != nil {
    	die("can't subscribe: %s", err.Error())
    }

    for {
    	if msg, err = sock.Recv(); err != nil {
    		die("can't recv: %s", err.Error())
    	}
        // update global variable
        //clientMsgRecv = string(msg)
        err := ioutil.WriteFile("msg/"+name+".txt", msg, 0644)
		if err != nil {
			log.Fatal(err)
		}

    	fmt.Printf("SERVER(%s): RECEIVED %s\n", name, string(msg))
    }
}

func node(args []string) {
	var sock mangos.Socket
	var err error
	var msg []byte
	var x int

	// read clientMsg from txt file
	b, err := ioutil.ReadFile("msg/" + args[1] + ".txt") // just pass the file name
    if err != nil {
        fmt.Print(err)
    }

    clientMsg := string(b) // convert content to a 'string'
    clientMsg =  clientMsg[:len(clientMsg)]

	if sock, err = bus.NewSocket(); err != nil {
		die("bus.NewSocket: %s", err)
	}

	if err = sock.Listen(args[2]); err != nil {
		die("sock.Listen: %s", err.Error())
	}

	// wait for everyone to start listening
	time.Sleep(time.Second)
	for x = 3; x < len(args); x++ {
		if err = sock.Dial(args[x]); err != nil {
			die("socket.Dial: %s", err.Error())
		}
	}

	// wait for everyone to join
	time.Sleep(time.Second)

	fmt.Printf("%s: SENDING '%s' ONTO BUS\n", args[1], args[1])
	if err = sock.Send([]byte(clientMsg + " by " + args[1])); err != nil {
		die("sock.Send: %s", err.Error())
	}

	for {
		if msg, err = sock.Recv(); err != nil {
			die("sock.Recv: %s", err.Error())
		}
		fmt.Printf("%s: RECEIVED \"%s\" FROM BUS\n", args[1], string(msg))
	}

}




func main() {
	if len(os.Args) > 2 && os.Args[1] == "client" {
		client(os.Args[2], os.Args[3])
		os.Exit(0)
	}
	if len(os.Args) > 3 && os.Args[1] == "server" {
		server(os.Args[2], os.Args[3])
		os.Exit(0)
	}
	//fmt.Fprintln(os.Stderr, "Usage: pubsub server|client <URL> <ARG>\n")
	

	if len(os.Args) > 3 && os.Args[1] != "server"{
		node(os.Args)
		os.Exit(0)
	}
	//xfmt.Fprintf(os.Stderr, "Usage: bus <NODENAME> <URL> <URL>... \n")
	os.Exit(1)
}