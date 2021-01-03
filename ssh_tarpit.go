package main

import (
    "fmt"
    "math/rand"
    "net"
    "strings"
    "time"
)

const (
    LISTEN_PORT = "2222"
    MIN_SLEEP = 1
    MAX_SLEEP = 5
    MIN_LENGTH = 10
    MAX_LENGTH = 120
)


func main() {
    // Start listening on the specified port
    listener, err := net.Listen("tcp", "0.0.0.0:" + LISTEN_PORT)
    if err != nil {
	panic(err)
    }
    
    // Close the listener when the app closes
    defer listener.Close()
    fmt.Println("Listening on " + LISTEN_PORT)
    
    // Start accepting connections
    for {
        conn, err := listener.Accept()
        if err != nil {
            // trap error and move onto a new connection
            continue
	}        
        
        // Handle the connection in a new goroutine
        fmt.Println("Tarpitting " + conn.RemoteAddr().String())
        go handleConnection(conn)
    }
    
}


func handleConnection(conn net.Conn) {
    // Handle the new connection, and write a long (and in fact, never-ending) banner
    
    // Close the connection when this function ends
    defer conn.Close()
    
    // We're going to be generating psuedo-random numbers, so seed it with the time the connection opened
    rand.Seed(time.Now().Unix())
    
    // Define some bits before we enter the loop
    var delay time.Duration
    var strlength int
    var randstr string
    
    // Main loop - get a random string, write it, sleep then do it again
    for {
        // Calculate a length for the string we should output
        strlength = rand.Intn(MAX_LENGTH - MIN_LENGTH) + MIN_LENGTH
        
        // Generate the string
        randstr = genString(strlength)
        
        // Write it to the socket
        _, err := conn.Write([]byte(randstr + "\r\n"))
        
        // Now check the write worked - if the client went away we'll get an error
        // at that point, we should stop wasting resources and free up the FD
        if err != nil {
            fmt.Println("Coward disconnected: " + conn.RemoteAddr().String())
            conn.Close()
            break
        }
        
        /* Sleep for a period before sending the next
         * We vary the period a bit to tie the client up for varying amounts of time
        */ 
        delay = time.Duration(rand.Intn(MAX_SLEEP - MIN_SLEEP) + MIN_SLEEP)
        time.Sleep(delay * time.Second)
    }
}


func genString(length int) (string){
    // Generate a psuedo-random string
    
    
    // Keep out charset ascii - it is pretending to be a printable banner, after all
    charSet := "abcdedfghijklmnopqrstABCDEDFGHIJKLMNOPQRSTUVWXYZ0123456789=.<>?!#@''"
    
    var output strings.Builder
    
    // Generate the string
    for i := 0; i < length; i++ {
        randnum := rand.Intn(len(charSet))
        randChar := charSet[randnum]
        output.WriteString(string(randChar))
    }
    
    return output.String()
    
}
