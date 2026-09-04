# Ubuntu Networking Casts — Complete Hands-On Demo Lab

> Tested-style commands for **Ubuntu 24.04**.  
> Run these demos on a normal Ubuntu VM or physical Ubuntu machine with `sudo` access.
>
> The lab covers:
>
> - TCP unicast
> - UDP unicast
> - ICMP
> - ARP / Ethernet broadcast
> - UDP IPv4 broadcast
> - DHCP observation
> - UDP multicast
> - IGMP observation
> - mDNS multicast
> - DNS over UDP and TCP
> - Linux network namespaces
> - Linux bridge
> - Simple HTTP load balancing
> - Anycast concepts and a local routing simulation
>
> **Important:** Commands that depend on your real LAN, DHCP server, or Internet access are marked accordingly.

---

# 1. Install Required Packages

```bash
sudo apt update

sudo apt install -y \
  iproute2 \
  iputils-ping \
  net-tools \
  tcpdump \
  traceroute \
  dnsutils \
  netcat-openbsd \
  curl \
  python3 \
  avahi-daemon \
  avahi-utils
```

Verify:

```bash
ip -V
ping -V
tcpdump --version
nc -h 2>&1 | head
dig -v
curl --version
```

---

# 2. Identify Your Interfaces

```bash
ip -br addr
```

Find the default route:

```bash
ip route
```

A typical result:

```text
default via 192.168.1.1 dev enp0s3
192.168.1.0/24 dev enp0s3 proto kernel scope link src 192.168.1.10
```

Your interface may be named:

```text
eth0
ens33
enp0s3
enp1s0
```

Throughout this guide, use `ip -br addr` to identify the correct interface instead of assuming `eth0`.

---

# 3. Port Numbers

A port identifies an application endpoint on a host.

Examples:

| Port | Typical Service |
|---:|---|
| 22 | SSH |
| 53 | DNS |
| 67 | DHCP server |
| 68 | DHCP client |
| 80 | HTTP |
| 443 | HTTPS |
| 5353 | mDNS |
| 5432 | PostgreSQL |
| 9000 | Arbitrary demo port |

In this lab, `9000`, `9001`, `9999`, etc. are simply convenient unused demo ports.

Check whether port 9000 is free:

```bash
sudo ss -lntup | grep ':9000'
```

No output normally means nothing is listening on that port.

---

# 4. Demo — ICMP Unicast

Terminal 1:

```bash
sudo tcpdump -i any -nn icmp
```

Terminal 2:

```bash
ping -c 4 8.8.8.8
```

If Internet ICMP is blocked, use your gateway:

```bash
ip route | grep default
```

Then:

```bash
ping -c 4 <gateway-ip>
```

Concept:

```text
Ubuntu                         Destination

       ICMP Echo Request
------------------------------------>

       ICMP Echo Reply
<------------------------------------
```

ICMP is neither TCP nor UDP.

---

# 5. Demo — TCP Unicast

## Terminal 1 — TCP server

```bash
nc -l 9000
```

## Terminal 2 — packet capture

```bash
sudo tcpdump -i lo -nn 'tcp port 9000'
```

## Terminal 3 — TCP client

```bash
nc 127.0.0.1 9000
```

Type:

```text
Hello TCP
```

You should see the text in Terminal 1.

The capture shows the TCP connection:

```text
Client                         Server

       SYN
-------------------------------->

       SYN-ACK
<--------------------------------

       ACK
-------------------------------->

       DATA
-------------------------------->
```

Check the listening socket before connecting:

```bash
ss -lntp | grep ':9000'
```

---

# 6. Demo — UDP Unicast

## Terminal 1

```bash
nc -u -l 9001
```

## Terminal 2

```bash
sudo tcpdump -i lo -nn 'udp port 9001'
```

## Terminal 3

```bash
echo "Hello UDP" | nc -u -w1 127.0.0.1 9001
```

Terminal 1 should display:

```text
Hello UDP
```

Notice that UDP has no TCP three-way handshake.

```text
Client                 Server

     UDP datagram
------------------------>
```

Check the UDP socket:

```bash
ss -lunp | grep ':9001'
```

---

# 7. Build a Safe Virtual LAN with Linux Network Namespaces

For broadcast and ARP demonstrations, a virtual LAN is more predictable than using the physical LAN.

We will create:

```text
                    Linux Bridge
                       br-demo
                     /         \
                    /           \
                   v             v
              client-ns       server-ns
              10.10.10.11     10.10.10.12
```

## Create namespaces

```bash
sudo ip netns add client-ns
sudo ip netns add server-ns
```

## Create bridge

```bash
sudo ip link add br-demo type bridge
sudo ip link set br-demo up
```

## Create client veth pair

```bash
sudo ip link add veth-client type veth peer name veth-client-br
sudo ip link set veth-client netns client-ns
sudo ip link set veth-client-br master br-demo
sudo ip link set veth-client-br up
```

## Create server veth pair

```bash
sudo ip link add veth-server type veth peer name veth-server-br
sudo ip link set veth-server netns server-ns
sudo ip link set veth-server-br master br-demo
sudo ip link set veth-server-br up
```

## Configure client

```bash
sudo ip netns exec client-ns ip link set lo up
sudo ip netns exec client-ns ip addr add 10.10.10.11/24 dev veth-client
sudo ip netns exec client-ns ip link set veth-client up
```

## Configure server

```bash
sudo ip netns exec server-ns ip link set lo up
sudo ip netns exec server-ns ip addr add 10.10.10.12/24 dev veth-server
sudo ip netns exec server-ns ip link set veth-server up
```

Verify:

```bash
sudo ip netns exec client-ns ip -br addr
sudo ip netns exec server-ns ip -br addr
```

Test:

```bash
sudo ip netns exec client-ns ping -c 3 10.10.10.12
```

---

# 8. Demo — ARP Broadcast

Clear the client's neighbor cache:

```bash
sudo ip netns exec client-ns ip neigh flush all
```

Terminal 1:

```bash
sudo tcpdump -i br-demo -nne arp
```

Terminal 2:

```bash
sudo ip netns exec client-ns ping -c 1 10.10.10.12
```

You should see an ARP request similar to:

```text
ARP, Request who-has 10.10.10.12 tell 10.10.10.11
```

The Ethernet destination for the request is:

```text
ff:ff:ff:ff:ff:ff
```

Concept:

```text
client-ns

"Who has 10.10.10.12?"

       |
       | Ethernet broadcast
       v

FF:FF:FF:FF:FF:FF

       |
       v

server-ns replies with its MAC
```

Check the learned neighbor:

```bash
sudo ip netns exec client-ns ip neigh
```

Important:

**ARP is not UDP.**

---

# 9. Demo — UDP IPv4 Broadcast

We will add a second receiver so the broadcast can be observed by multiple hosts.

```text
                      br-demo
                  /      |      \
                 /       |       \
                v        v        v
             client    server   server2
             .11       .12      .13
```

Create the third namespace:

```bash
sudo ip netns add server2-ns

sudo ip link add veth-server2 type veth peer name veth-server2-br
sudo ip link set veth-server2 netns server2-ns

sudo ip link set veth-server2-br master br-demo
sudo ip link set veth-server2-br up

sudo ip netns exec server2-ns ip link set lo up
sudo ip netns exec server2-ns ip addr add 10.10.10.13/24 dev veth-server2
sudo ip netns exec server2-ns ip link set veth-server2 up
```

### Receiver 1

Terminal 1:

```bash
sudo ip netns exec server-ns nc -u -l 9999
```

### Receiver 2

Terminal 2:

```bash
sudo ip netns exec server2-ns nc -u -l 9999
```

### Capture

Terminal 3:

```bash
sudo tcpdump -i br-demo -nne 'udp port 9999'
```

### Sender

For reliable broadcast behavior, use this small Python sender from the client namespace:

```bash
cat > /tmp/udp-broadcast.py <<'PY'
import socket

s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
s.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)

s.sendto(
    b"Hello UDP Broadcast\n",
    ("10.10.10.255", 9999)
)

s.close()
PY
```

Run:

```bash
sudo ip netns exec client-ns python3 /tmp/udp-broadcast.py
```

Both listeners should receive:

```text
Hello UDP Broadcast
```

This demonstrates:

```text
UDP
 +
IPv4 Broadcast
 =
One sender -> all hosts listening on that LAN/port
```

---

# 10. Demo — DHCP Traffic

This demo observes your real Ubuntu network.

First determine your physical interface:

```bash
ip route | grep default
```

Suppose it is `enp0s3`.

Capture:

```bash
sudo tcpdump -i enp0s3 -nne 'udp port 67 or udp port 68'
```

DHCPv4 uses:

```text
UDP 67 = server
UDP 68 = client
```

A simplified DHCP sequence:

```text
Client                     DHCP Server

 Discover -------------------->
 <----------------------- Offer

 Request --------------------->
 <------------------------- ACK
```

This is commonly remembered as:

```text
D O R A

Discover
Offer
Request
Acknowledge
```

Whether you observe all phases, and whether each is broadcast or unicast, depends on the current client state and network configuration.

Do **not** deliberately disrupt a production/server network just to force DHCP renewal.

---

# 11. Demo — UDP Multicast

Use multicast group:

```text
239.10.10.10
```

and port:

```text
10000
```

Python is used here because its multicast socket behavior is explicit and reproducible.

## Multicast receiver

```bash
cat > /tmp/multicast-receiver.py <<'PY'
import socket
import struct

GROUP = "239.10.10.10"
PORT = 10000

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)

sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
sock.bind(("", PORT))

membership = struct.pack(
    "=4s4s",
    socket.inet_aton(GROUP),
    socket.inet_aton("0.0.0.0")
)

sock.setsockopt(
    socket.IPPROTO_IP,
    socket.IP_ADD_MEMBERSHIP,
    membership
)

print(f"Listening on multicast {GROUP}:{PORT}")

while True:
    data, addr = sock.recvfrom(65535)
    print(f"{addr}: {data.decode(errors='replace')}")
PY
```

Run receiver:

```bash
python3 /tmp/multicast-receiver.py
```

You can open a second terminal and run another receiver:

```bash
python3 /tmp/multicast-receiver.py
```

## Multicast sender

```bash
cat > /tmp/multicast-sender.py <<'PY'
import socket

GROUP = "239.10.10.10"
PORT = 10000

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)

sock.setsockopt(socket.IPPROTO_IP, socket.IP_MULTICAST_TTL, 1)

sock.sendto(
    b"Hello Multicast",
    (GROUP, PORT)
)

sock.close()
PY
```

Capture:

```bash
sudo tcpdump -i any -nn 'host 239.10.10.10 and udp port 10000'
```

Send:

```bash
python3 /tmp/multicast-sender.py
```

Both receiver processes should display the multicast message.

Concept:

```text
                   239.10.10.10
                        |
Sender -----------------+
                        |
                   +----+----+
                   |         |
                   v         v
               Receiver1  Receiver2
```

---

# 12. Demo — Inspect Multicast Membership

Run:

```bash
ip maddr show
```

While the multicast receiver is running, inspect multicast membership.

For IPv4 multicast protocol traffic:

```bash
sudo tcpdump -i any -nn igmp
```

Important:

**IGMP is not UDP.**

IGMP manages IPv4 multicast group membership between hosts and multicast-aware network infrastructure.

---

# 13. Demo — mDNS

mDNS uses:

```text
IPv4 multicast: 224.0.0.251
UDP port:       5353
```

Make sure Avahi is running:

```bash
sudo systemctl enable --now avahi-daemon
```

Browse services:

```bash
avahi-browse -a
```

Capture:

```bash
sudo tcpdump -i any -nn 'udp port 5353'
```

You may see traffic to:

```text
224.0.0.251.5353
```

Concept:

```text
mDNS
 |
 +--> UDP
 |
 +--> Multicast
 |
 +--> 224.0.0.251:5353
```

---

# 14. Demo — DNS over UDP

Capture:

```bash
sudo tcpdump -i any -nn 'udp port 53'
```

Another terminal:

```bash
dig @8.8.8.8 example.com
```

If your network blocks direct DNS to `8.8.8.8`, use the resolver shown by:

```bash
resolvectl status
```

DNS normally uses UDP for many ordinary queries.

Concept:

```text
Client                       DNS Resolver

     UDP DNS Query
------------------------------>

     UDP DNS Response
<------------------------------
```

---

# 15. Demo — DNS over TCP

Capture:

```bash
sudo tcpdump -i any -nn 'tcp port 53'
```

Run:

```bash
dig +tcp @8.8.8.8 example.com
```

Now you should see a TCP connection, including the handshake.

This demonstrates that DNS can use both:

```text
DNS over UDP
DNS over TCP
```

---

# 16. Demo — HTTP/TCP

Start a web server:

```bash
python3 -m http.server 8080 --bind 127.0.0.1
```

Capture:

```bash
sudo tcpdump -i lo -nn 'tcp port 8080'
```

Request:

```bash
curl http://127.0.0.1:8080/
```

Concept:

```text
curl
 |
 | HTTP
 v
TCP
 |
 v
127.0.0.1:8080
 |
 v
Python HTTP server
```

---

# 17. Demo — Simple Load Balancer

We will run three backend HTTP servers and one simple round-robin load balancer.

Architecture:

```text
                    Client
                      |
                      v
                127.0.0.1:8080
                      |
                Load Balancer
                 /     |     \
                /      |      \
               v       v       v
            :8081    :8082    :8083
           Server1  Server2  Server3
```

Create backend directories:

```bash
mkdir -p /tmp/lb-demo/server1
mkdir -p /tmp/lb-demo/server2
mkdir -p /tmp/lb-demo/server3

echo "Hello from Server-1" > /tmp/lb-demo/server1/index.html
echo "Hello from Server-2" > /tmp/lb-demo/server2/index.html
echo "Hello from Server-3" > /tmp/lb-demo/server3/index.html
```

Run each backend in a separate terminal.

Terminal 1:

```bash
cd /tmp/lb-demo/server1
python3 -m http.server 8081 --bind 127.0.0.1
```

Terminal 2:

```bash
cd /tmp/lb-demo/server2
python3 -m http.server 8082 --bind 127.0.0.1
```

Terminal 3:

```bash
cd /tmp/lb-demo/server3
python3 -m http.server 8083 --bind 127.0.0.1
```

Create the load balancer:

```bash
cat > /tmp/lb-demo/lb.py <<'PY'
from http.server import BaseHTTPRequestHandler, HTTPServer
import itertools
import urllib.request

BACKENDS = [
    "http://127.0.0.1:8081",
    "http://127.0.0.1:8082",
    "http://127.0.0.1:8083",
]

backend_cycle = itertools.cycle(BACKENDS)

class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        backend = next(backend_cycle)

        try:
            with urllib.request.urlopen(backend + self.path) as response:
                body = response.read()

            self.send_response(200)
            self.send_header("Content-Type", "text/plain")
            self.end_headers()
            self.wfile.write(body)

        except Exception as exc:
            self.send_response(502)
            self.end_headers()
            self.wfile.write(str(exc).encode())

    def log_message(self, fmt, *args):
        pass

HTTPServer(("127.0.0.1", 8080), Handler).serve_forever()
PY
```

Terminal 4:

```bash
python3 /tmp/lb-demo/lb.py
```

Test:

```bash
for i in {1..9}; do
    curl -s http://127.0.0.1:8080/
done
```

Expected pattern:

```text
Hello from Server-1
Hello from Server-2
Hello from Server-3
Hello from Server-1
Hello from Server-2
Hello from Server-3
...
```

Explain:

```text
Load balancer:
chooses a backend

Anycast:
routing chooses one service site
```

They are different concepts.

---

# 18. Anycast — What Can Be Demonstrated Locally?

Real Internet Anycast normally involves:

- Multiple routing locations
- The same IP prefix advertised from those locations
- Routing protocols such as BGP
- Routers choosing a preferred route

Concept:

```text
                     SAME SERVICE PREFIX
                           10.100.0.1
                         /            \
                        /              \
                       v                v
                    Site-A           Site-B
                       \                /
                        \              /
                         \            /
                           Router
                             |
                             v
                           Client
```

The client sends to:

```text
10.100.0.1
```

but routing decides which site receives the traffic.

---

# 19. Simple Local Anycast Routing Simulation

This demo illustrates the **one address, multiple possible service locations, routing chooses one path** concept.

For a production-quality BGP Anycast lab, use FRRouting/containerlab. This simpler namespace demo avoids requiring FRR.

Create router and two sites:

```bash
sudo ip netns add router-ns
sudo ip netns add site-a-ns
sudo ip netns add site-b-ns
```

Create router ↔ Site A:

```bash
sudo ip link add r-a type veth peer name a-r
sudo ip link set r-a netns router-ns
sudo ip link set a-r netns site-a-ns
```

Create router ↔ Site B:

```bash
sudo ip link add r-b type veth peer name b-r
sudo ip link set r-b netns router-ns
sudo ip link set b-r netns site-b-ns
```

Configure Router:

```bash
sudo ip netns exec router-ns ip link set lo up

sudo ip netns exec router-ns ip addr add 172.16.1.1/30 dev r-a
sudo ip netns exec router-ns ip link set r-a up

sudo ip netns exec router-ns ip addr add 172.16.2.1/30 dev r-b
sudo ip netns exec router-ns ip link set r-b up
```

Configure Site A:

```bash
sudo ip netns exec site-a-ns ip link set lo up

sudo ip netns exec site-a-ns ip addr add 172.16.1.2/30 dev a-r
sudo ip netns exec site-a-ns ip link set a-r up

sudo ip netns exec site-a-ns ip addr add 10.100.0.1/32 dev lo

sudo ip netns exec site-a-ns ip route add default via 172.16.1.1
```

Configure Site B:

```bash
sudo ip netns exec site-b-ns ip link set lo up

sudo ip netns exec site-b-ns ip addr add 172.16.2.2/30 dev b-r
sudo ip netns exec site-b-ns ip link set b-r up

sudo ip netns exec site-b-ns ip addr add 10.100.0.1/32 dev lo

sudo ip netns exec site-b-ns ip route add default via 172.16.2.1
```

Notice:

```text
Site A loopback = 10.100.0.1
Site B loopback = 10.100.0.1
```

The same service address exists in two isolated network namespaces.

Tell the router to prefer Site A:

```bash
sudo ip netns exec router-ns \
  ip route add 10.100.0.1/32 via 172.16.1.2
```

Check:

```bash
sudo ip netns exec router-ns ip route get 10.100.0.1
```

Expected route is via:

```text
172.16.1.2
```

Ping:

```bash
sudo ip netns exec router-ns ping -c 2 10.100.0.1
```

Now simulate Site A route withdrawal:

```bash
sudo ip netns exec router-ns \
  ip route del 10.100.0.1/32 via 172.16.1.2
```

Add Site B route:

```bash
sudo ip netns exec router-ns \
  ip route add 10.100.0.1/32 via 172.16.2.2
```

Check:

```bash
sudo ip netns exec router-ns ip route get 10.100.0.1
```

Now the same destination:

```text
10.100.0.1
```

is reached through Site B.

Diagram:

```text
Before:

                  10.100.0.1
                       |
                     Site A ✓
                       |
Router ----------------+
                       |
                     Site B
                  10.100.0.1


After route change:

                  10.100.0.1
                       |
                     Site A
                       |
Router ----------------+
                       |
                     Site B ✓
                  10.100.0.1
```

This demonstrates the central Anycast idea, but it is **static routing**, not BGP.

Real Internet Anycast commonly automates route selection/withdrawal through BGP and routing policy.

---

# 20. Compare All Four Cast Models

| Model | Meaning | TCP | UDP | Example |
|---|---|---:|---:|---|
| Unicast | One → One | Yes | Yes | HTTP, SSH, DNS |
| Broadcast | One → all hosts in local IPv4/L2 scope | No normal TCP broadcast | Yes | DHCPv4; ARP uses Ethernet broadcast |
| Multicast | One → subscribed group | Traditional TCP: No | Yes | mDNS, IPTV |
| Anycast | One → one selected site among many | Yes | Yes | Global DNS, CDN/edge, global LB entry |

---

# 21. Protocol Mapping

| Protocol / Technology | Transport | Communication Model | Purpose |
|---|---|---|---|
| HTTP/1.1 / HTTP/2 | TCP | Usually Unicast | Web/API |
| HTTP/3 | QUIC/UDP | Usually Unicast; may target Anycast service | Web/API |
| SSH | TCP | Unicast | Remote login |
| DNS | UDP/TCP | Usually Unicast; service may be Anycast | Name resolution |
| DHCPv4 | UDP | Broadcast + Unicast depending on state | IP configuration |
| ARP | Neither TCP nor UDP | Ethernet broadcast request | IPv4 → MAC resolution |
| mDNS | UDP | Multicast | Local discovery |
| IGMP | Neither TCP nor UDP | Multicast control | Group membership |
| ICMP | Neither TCP nor UDP | Usually Unicast | Diagnostics/control |
| OSPF | Directly over IP | Can use multicast | Routing |
| BGP | TCP 179 | Unicast peer sessions | Route exchange |
| Anycast | Routing architecture | Can carry TCP/UDP/etc. | Select one service site |

---

# 22. Useful Packet-Capture Filters

TCP:

```bash
sudo tcpdump -i any -nn tcp
```

UDP:

```bash
sudo tcpdump -i any -nn udp
```

ICMP:

```bash
sudo tcpdump -i any -nn icmp
```

ARP:

```bash
sudo tcpdump -i any -nne arp
```

DNS:

```bash
sudo tcpdump -i any -nn 'port 53'
```

DHCP:

```bash
sudo tcpdump -i any -nne 'udp port 67 or udp port 68'
```

mDNS:

```bash
sudo tcpdump -i any -nn 'udp port 5353'
```

IGMP:

```bash
sudo tcpdump -i any -nn igmp
```

Specific TCP port:

```bash
sudo tcpdump -i any -nn 'tcp port 9000'
```

Specific UDP port:

```bash
sudo tcpdump -i any -nn 'udp port 9999'
```

---

# 23. Useful Linux Commands During Training

Interfaces:

```bash
ip -br addr
```

Routes:

```bash
ip route
```

Neighbors / ARP cache:

```bash
ip neigh
```

TCP listening sockets:

```bash
ss -lntp
```

UDP listening sockets:

```bash
ss -lunp
```

All sockets:

```bash
ss -tulpn
```

Network namespaces:

```bash
sudo ip netns list
```

Bridge:

```bash
ip link show br-demo
```

Bridge ports:

```bash
bridge link
```

Multicast addresses:

```bash
ip maddr show
```

---

# 24. Cleanup

After the namespace demos:

```bash
sudo ip netns del client-ns 2>/dev/null || true
sudo ip netns del server-ns 2>/dev/null || true
sudo ip netns del server2-ns 2>/dev/null || true

sudo ip link del br-demo 2>/dev/null || true
```

Anycast simulation cleanup:

```bash
sudo ip netns del router-ns 2>/dev/null || true
sudo ip netns del site-a-ns 2>/dev/null || true
sudo ip netns del site-b-ns 2>/dev/null || true
```

Temporary files:

```bash
rm -f /tmp/udp-broadcast.py
rm -f /tmp/multicast-receiver.py
rm -f /tmp/multicast-sender.py
rm -rf /tmp/lb-demo
```

---

# 25. Recommended Live Training Order

Use this sequence:

```text
1. IP + ports
        ↓
2. ICMP ping
        ↓
3. TCP unicast + handshake
        ↓
4. UDP unicast
        ↓
5. Linux namespaces
        ↓
6. ARP broadcast
        ↓
7. UDP broadcast
        ↓
8. DHCP observation
        ↓
9. UDP multicast
        ↓
10. IGMP
        ↓
11. mDNS
        ↓
12. DNS UDP vs TCP
        ↓
13. HTTP
        ↓
14. Load balancing
        ↓
15. Anycast routing simulation
```

For almost every demo, keep `tcpdump` visible.

That allows students to connect:

```text
Application
     ↓
Socket / Port
     ↓
TCP / UDP
     ↓
IP
     ↓
Ethernet
     ↓
Actual packet
```

---

# 26. Final Mental Model

```text
                        APPLICATION
              HTTP / DNS / DHCP / mDNS
                            |
                            v
                       TCP / UDP
                            |
                            v
                            IP
                            |
          +-----------------+-----------------+
          |                 |                 |
       Unicast          Multicast          Anycast
          |
     IPv4/L2 Broadcast
     where applicable
```

Remember:

```text
Unicast   = One → One

Broadcast = One → Everyone in local broadcast scope

Multicast = One → Group

Anycast   = One → One selected instance/site
```

And:

```text
TCP / UDP
    !=
Unicast / Broadcast / Multicast / Anycast
```

TCP and UDP are transport protocols.

The cast terms describe delivery/addressing/routing behavior.
