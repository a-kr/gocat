
package main;

import (
    "log"
    "net"
    "os"
    "io"
)

func dieOnError(err error) {
    if err != nil {
        log.Fatalf("Error: %s", err)
    }
}

func MustResolveTCPAddr(addr string) *net.TCPAddr {
    a, err := net.ResolveTCPAddr("tcp", addr)
    dieOnError(err)
    return a
}

func main() {
    if len(os.Args) < 3 {
        log.Printf("gocat: trivial tcp proxy / port forwarder")
        log.Fatalf("Usage: gocat <address:port-to-listen-on> <address:port-to-forward-to>")
    }

    sAddr := os.Args[1]
    dAddr := os.Args[2]

    ss, err := net.ListenTCP("tcp", MustResolveTCPAddr(sAddr))
    dieOnError(err)
    for {
        client, err := ss.AcceptTCP()
        dieOnError(err)
        go handleConnection(client, dAddr)
    }
}

func handleConnection(client *net.TCPConn, dAddr string) {
    log.Printf("Accepted new connection from %v", client.RemoteAddr())
    conn, err := net.DialTCP("tcp", nil, MustResolveTCPAddr(dAddr))
    dieOnError(err)

    go func() {
        n, err := io.Copy(conn, client)
        log.Printf("Closed client-to-dst connection %v after writing %d bytes: %v", client.RemoteAddr(), n, err)
    }()
    go func() {
        n, err := io.Copy(client, conn)
        log.Printf("Closed dst-to-client connection %v after writing %d bytes: %v", client.RemoteAddr(), n, err)
    }()
}


