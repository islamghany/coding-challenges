# Takeaways

## OSI Model (open Systems Interconnection Model)

OSI model is a conceptual framework that standardizes the functions of a telecommunication or computing system into seven abstraction layers. A layer serves the layer above it and is served by the layer below it. For example, a layer that provides error-free communications across a network provides the path needed by applications above it, while it calls the next lower layer to send and receive packets that make up the contents of that path. Two instances at one layer are connected by a horizontal connection on that layer.

### Layer 1: Physical Layer
examples: cables, pins, signals, voltages, etc.

This layer deals with the physical connection of the devices. It defines the cable, the pin, the signal, the voltage, etc. It is responsible for the transmission and reception of the unstructured raw data over a physical medium.
it's handled by the hardware and the operating system.

### Layer 2: Data Link Layer
examples: Ethernet, Wi-Fi, etc.

This layer is responsible for framing packets, error detection, and MAC addressing. It is divided into two sub-layers: Logical Link Control (LLC) and Media Access Control (MAC). It is responsible for the transmission of data between two devices on the same network.
the data link layer is managed by the operating system.

### Layer 3: Network Layer
examples: IP, ICMP, ARP, etc.

This layer is responsible for routing packets from the source to the destination. It is responsible for logical addressing and routing. It is responsible for the transmission of data between two devices on different networks.
we can interct with that layer when dealing with IP addresses and routing. for example, you might configure which IP address and port your server should listen to.

### Layer 4: Transport Layer
examples: TCP, UDP, etc.

This layer ensures reliable data transfer between two devices. It is responsible for end-to-end communication. It is responsible for the transmission of data between two devices on the same or different networks.
the transport layer provides services such as Flow Control, Error Detection, Error Recovery, and Segmentation.

**Segmentation**: The transport layer divides the message into smaller packets to be sent over the network. this segmentation is necessary because the network layer has a maximum packet size that it can handle.

**Flow Control**: The transport layer is responsible for controlling the flow of data between two devices. it makes sure that the sender does not overwhelm the receiver with data.

**Error Detection**: The transport layer is responsible for detecting errors in the data. it uses checksums to detect errors in the data.

**Error Recovery**: The transport layer is responsible for recovering from errors in the data. it uses retransmission to recover from errors in the data.

We can interact with the transport layer when dealing with TCP and UDP protocols in our code.

```go

ln, err := net.Listen("tcp", ":8080")
...
for {
	conn, err := ln.Accept()
	...
	go handleConnection(conn)
}
```

### Layer 5: Session Layer
examples: establishing, maintaining, and terminating connections between two devices.

This layer is responsible for establishing, maintaining, and terminating connections between two devices. It is responsible for the synchronization and dialog control between two devices. It is responsible for the transmission of data between two devices on the same or different networks.

we can manage sessions in our code by using cookies and sessions.

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		cookie = &http.Cookie{
			Name:  "session",
			Value: "some value",
		}
		http.SetCookie(w, cookie)
	}
	fmt.Fprintln(w, "cookie:", cookie)
})
```

### Layer 6: Presentation Layer
examples: encryption, compression, data translation, etc.

This layer is responsible for data translation, encryption, and compression. It is responsible for the translation of data between two devices. It is responsible for the transmission of data between two devices on the same or different networks.

we can interact with the presentation layer when dealing with data encoding and decoding.

```go
func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, TLS!")
    })

    server := &http.Server{
        Addr: ":8443",
        TLSConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
        },
    }

    err := server.ListenAndServeTLS("cert.pem", "key.pem")
    if err != nil {
        panic(err)
    }
}
```
above code is an example of how we can interact with the presentation layer by using TLS encryption.

### Layer 7: Application Layer
exmaples: HTTP, FTP, SMTP, etc.

This layer us where the actual application resides. It is responsible for providing services to the user. It is responsible for the transmission of data between two devices on the same or different networks.

we can interact with the application layer by writing our application code.

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
})
fmt.Println("Server is listening on port 8080...")
http.ListenAndServe(":8080", nil)
```


## Load Balancer

A load balancer is a case of a reverse proxy. It distributes incoming network traffic across multiple servers. It is responsible for balancing the load between multiple servers. It is responsible for the transmission of data between two devices on the same or different networks.

There are several Algorithms that can be used to distribute the load between servers:
- Round Robin
- Least Connections
- IP Hash
- Content-Based

### Layer 4 Load Balancer

A Layer 4 load balancer operates at the transport layer. It forwards the traffic based on network information such as IP address and TCP port.

**Pros**:
- Faster than Layer 7 load balancer.
- Can handle more connections than Layer 7 load balancer.
- More secure
- Works well with UDP and TCP protocols.

**Cons**:
- Cannot inspect the content of the packets.
- Cannot make decisions based on the content of the packets.
- No caching.
- No SSL termination.

### Layer 7 Load Balancer

A Layer 7 load balancer operates at the application layer. It forwards the traffic based on the content of the packets.

**Pros**:
- Can inspect the content of the packets.
- Can make decisions based on the content of the packets.
- Caching.
- Great for microservices architecture.
- Authentication and authorization.

**Cons**:
- Slower than Layer 4 load balancer.
- Decryption and encryption overhead.
- Cannot handle as many connections as Layer 4 load balancer.
- Needs to buffer the entire request before forwarding it.
- More complex to configure.
