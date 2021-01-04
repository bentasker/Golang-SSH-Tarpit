/** Golang SSH Tarpit
 * 
 * Accept connections from a downstream client, and send them an infinitely long
 * SSH banner, pausing between each line for various times.
 *
 * The idea being that they're tied up with your tarpit and can't go bother others 
 * 
 * Copyright (C) 2021 B Tasker
 * Released under GNU GPL V3 - See LICENSE
 * 
 */
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
    MAX_SLEEP = 5 // Don't set too high or the client will timeout, suggest < 30
    MIN_LENGTH = 10
    MAX_LENGTH = 120 // This must not be set higher than 253 - SSH spec says 255 max, including the CRLF
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
    var start = time.Now().Unix()
    rand.Seed(start)
    
    
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
            var delta = int(time.Now().Unix() - start)
            fmt.Println("Coward disconnected:", conn.RemoteAddr().String(), "after", delta, "seconds")
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
