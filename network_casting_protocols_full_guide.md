# Network Communication Models: Unicast, Broadcast, Multicast, and Anycast

## 1. Why This Topic Matters

When an application sends data over a network, an important question is:

> **Who should receive this packet?**

Depending on the requirement, communication can be:

- **Unicast** — one sender to one specific receiver.
- **Broadcast** — one sender to every host in a local broadcast domain.
- **Multicast** — one sender to a selected group of receivers.
- **Anycast** — one sender to one preferred destination among multiple instances advertising the same service address.

These are **delivery/addressing models**, not alternatives to TCP and UDP.

TCP and UDP operate at the **transport layer**. Unicast, broadcast, multicast, and anycast describe how traffic is addressed and/or routed to receivers.

---

## 2. First Understand the Layers

A simplified TCP/IP stack:

```text
+--------------------------------------------------+
| Application                                      |
| HTTP, HTTPS, DNS, SSH, DHCP, mDNS, etc.          |
+--------------------------------------------------+
| Transport                                        |
| TCP / UDP                                        |
+--------------------------------------------------+
| Internet / Network                               |
| IPv4 / IPv6 / ICMP / routing                     |
+--------------------------------------------------+
| Data Link                                        |
| Ethernet / Wi-Fi / MAC addressing                |
+--------------------------------------------------+
| Physical                                         |
| Cable / Fiber / Radio                            |
+--------------------------------------------------+
```

A common mistake is to put these concepts in one list:

```text
TCP
UDP
Unicast
Broadcast
Multicast
Anycast
```

They represent different concepts.

A better mental model is:

```text
Transport protocols
├── TCP
└── UDP

Delivery/addressing/routing models
├── Unicast
├── Broadcast
├── Multicast
└── Anycast
```

---

# 3. TCP and UDP

## 3.1 TCP

TCP stands for **Transmission Control Protocol**.

TCP is:

- Connection-oriented
- Reliable
- Ordered
- Byte-stream based
- Uses acknowledgements and retransmissions
- Uses flow control
- Uses congestion control

Before application data is normally exchanged, TCP establishes a connection using the three-way handshake:

```text
Client                                  Server

   -------- SYN ------------------------>
   <------- SYN + ACK -------------------
   -------- ACK ------------------------>

              Connection established
```

Common TCP applications:

- HTTP/1.1 and HTTP/2
- HTTPS
- SSH
- Database connections
- SMTP
- IMAP
- Many REST APIs

TCP fundamentally models communication between **two endpoints**, so normal TCP communication is unicast.

---

## 3.2 UDP

UDP stands for **User Datagram Protocol**.

UDP is:

- Connectionless
- Datagram/message oriented
- No transport-level delivery guarantee
- No transport-level retransmission
- No guaranteed ordering
- Lightweight compared with TCP

Conceptually:

```text
Application
    |
    v
+------------------+
| UDP Header       |
| Source Port      |
| Destination Port |
+------------------+
| Application Data |
+------------------+
    |
    v
IP
```

UDP can be used with:

- Unicast
- IPv4 broadcast
- IP multicast
- Anycast destinations

Examples include DNS, DHCP, mDNS, RTP-based media, and many custom real-time protocols.

> Note: QUIC uses UDP as its underlying transport but implements reliability, streams, encryption, and congestion control above UDP. HTTP/3 runs over QUIC.

---

# 4. Unicast

## Definition

**One sender communicates with one specific destination.**

```text
A ----------------------> B

       ONE TO ONE
```

Example:

```text
Laptop                           Web Server
10.0.0.10 ---------------------> 203.0.113.20
```

The destination address identifies a specific endpoint.

## TCP and UDP

Both TCP and UDP support ordinary unicast communication.

### TCP unicast

```text
Browser ----- TCP/HTTPS -----> Web Server
```

### UDP unicast

```text
DNS Client ----- UDP -------> DNS Resolver
```

## Common examples

| Protocol/Application | Typical Transport | Purpose |
|---|---|---|
| HTTP/1.1, HTTP/2 | TCP | Web applications |
| HTTPS | TCP (HTTP/1.1 or HTTP/2) | Secure web communication |
| SSH | TCP | Remote administration |
| DNS | UDP/TCP | Name resolution |
| Database protocols | Usually TCP | Application-to-database communication |
| ICMP Echo | Neither TCP nor UDP | `ping` / diagnostics |

## Use cases

- Browser to web server
- Microservice to microservice
- Application to database
- SSH to server
- Client to DNS resolver
- API calls

---

# 5. Broadcast

## Definition

Broadcast means:

> **One sender sends traffic to all hosts in a broadcast domain.**

```text
                  +----> Host B
                  |
Host A -----------+----> Host C
                  |
                  +----> Host D

             ONE TO ALL
```

Broadcast is especially associated with **IPv4 and Ethernet LANs**.

---

## 5.1 Ethernet Broadcast

Ethernet has a special broadcast MAC address:

```text
FF:FF:FF:FF:FF:FF
```

A frame sent to this destination is delivered throughout the Layer-2 broadcast domain.

Example:

```text
                   Switch
                 /   |   \
                /    |    \
               v     v     v
             PC-B  PC-C  PC-D
               ^     ^     ^
                \    |    /
                 \   |   /
                   PC-A

Destination MAC:
FF:FF:FF:FF:FF:FF
```

---

## 5.2 ARP and Broadcast

ARP is used in IPv4 LANs to resolve an IPv4 address to a MAC address.

Suppose:

```text
PC-A IP: 192.168.1.10
PC-B IP: 192.168.1.20
```

PC-A wants to communicate with PC-B but does not know PC-B's MAC address.

PC-A sends an ARP request as an Ethernet broadcast:

```text
Who has 192.168.1.20?
Tell 192.168.1.10
```

All devices in the Layer-2 broadcast domain see the request, but the owner of the requested IPv4 address responds.

Important:

> **ARP is not UDP.** ARP is its own protocol carried directly in Ethernet frames.

---

## 5.3 DHCP and UDP Broadcast

DHCP is a classic example involving UDP and IPv4 broadcast.

A newly connected machine may initially know neither:

- Its own IPv4 address
- The DHCP server's address

So DHCP discovery can use broadcast.

```text
New Client
    |
    | DHCPDISCOVER
    | UDP
    v
255.255.255.255
    |
    +------> Host A
    |
    +------> Host B
    |
    +------> DHCP Server
```

Common DHCPv4 ports:

```text
Server: UDP 67
Client: UDP 68
```

A simplified DHCP sequence is often remembered as DORA:

```text
Client                         DHCP Server

   ---- Discover -------------->
   <--- Offer ------------------
   ---- Request ---------------->
   <--- Acknowledge -------------
```

The exact use of broadcast/unicast during DHCP can vary by stage and client/server state.

---

## 5.4 Broadcast Domain

Routers normally separate Layer-2 broadcast domains.

```text
LAN A                    Router                    LAN B

A ----\
B ----- Switch ---------- R ---------- Switch ----- D
C ----/                                          \-- E
```

A normal Layer-2 broadcast originating in LAN A is not simply flooded by the router into LAN B.

This is one reason routers and VLANs are important for controlling broadcast scope.

---

## Broadcast use cases

- ARP requests in IPv4 Ethernet LANs
- DHCPv4 discovery
- Some legacy LAN discovery protocols

IPv6 does **not** use IPv4-style IP broadcast. IPv6 relies heavily on multicast instead.

---

# 6. Multicast

## Definition

Multicast means:

> **One sender sends traffic to a specific group of interested receivers.**

```text
                    +----> Receiver B ✓
                    |
Sender -------------+----> Receiver C ✓
                    |
                    +----> Receiver D ✓

                    X----> Host E
                          Not in group

               ONE TO GROUP
```

Instead of sending an independent copy from the source to every receiver, multicast-capable networks can replicate traffic where paths branch.

---

## 6.1 IPv4 Multicast Addresses

IPv4 multicast uses addresses in:

```text
224.0.0.0/4
```

which corresponds to:

```text
224.0.0.0 - 239.255.255.255
```

Not every address in this range is interchangeable; portions have specific scopes and assignments.

---

## 6.2 UDP and Multicast

UDP works naturally with multicast because UDP does not require a point-to-point connection.

```text
                      Multicast Group
                        239.x.x.x

Sender ----- UDP ------+----> Receiver A
                       |
                       +----> Receiver B
                       |
                       +----> Receiver C
```

Common applications include:

- IPTV in managed networks
- Market-data distribution
- Real-time media in controlled networks
- Service discovery

---

## 6.3 IGMP

**IGMP — Internet Group Management Protocol**

IGMP is used between IPv4 hosts and nearby multicast routers to manage multicast group membership.

Conceptually:

```text
Host
 |
 | "I want multicast group X"
 v
Router
```

IGMP is **not UDP**. It is an IP-layer control protocol.

---

## 6.4 mDNS

mDNS stands for **Multicast DNS**.

It uses multicast to resolve names and discover services on a local network without requiring a conventional centralized DNS server for those local names.

Typical environments:

- Printers
- Local device discovery
- Development machines
- Apple Bonjour-compatible discovery

---

## 6.5 OSPF

OSPF is a routing protocol and uses multicast on applicable IPv4 network types.

Well-known IPv4 multicast addresses include:

```text
224.0.0.5   AllSPFRouters
224.0.0.6   AllDRouters
```

OSPF is **not UDP or TCP**; it runs directly over IP.

---

# 7. Anycast

Anycast is especially important in **system design and global infrastructure**.

## Definition

Anycast means:

> **Multiple service locations advertise reachability for the same IP prefix, and routing delivers a packet to one selected location according to routing policy.**

Conceptually:

```text
                     Same service IP/prefix
                              |
                       Internet Routing
                         /    |    \
                        /     |     \
                       v      v      v
                    India   Europe   USA
                     Site    Site    Site
```

A client sends traffic to one destination IP.

The Internet routing system determines which advertised route is preferred.

---

## 7.1 Anycast is NOT multicast

Multicast:

```text
A ----------+----> B
            |
            +----> C
            |
            +----> D

One packet flow intended for multiple group members
```

Anycast:

```text
             +---- B
             |
A -----------+---- C ✓ selected route
             |
             +---- D

Traffic is routed toward one selected site
```

---

## 7.2 Anycast and BGP

BGP stands for **Border Gateway Protocol**.

BGP is the inter-domain routing protocol used between autonomous systems on the Internet.

With Anycast, multiple sites can advertise the same IP prefix.

```text
             Site A
                \
                 \ advertises prefix P
                  \
                   Internet / BGP
                  /
                 / advertises prefix P
                /
             Site B
```

Routers evaluate available BGP routes and select a preferred route according to routing policy.

Important:

> BGP is not an "Anycast protocol." Anycast is an addressing/routing architecture commonly implemented by advertising the same prefix from multiple locations using BGP.

---

## 7.3 "Nearest" is a Simplification

People often say:

> Anycast sends traffic to the nearest server.

A more accurate statement is:

> Anycast sends traffic along the route selected by the routing system.

The selected site is often topologically or geographically favorable, but geographic distance alone does not determine BGP path selection.

---

# 8. Anycast with Load Balancers

Suppose a global application has users in India, Europe, and the United States.

A simplistic architecture might send everybody toward one load-balancing site:

```text
India Users ----\
Europe Users ----+----> One LB ----> Servers
US Users --------/
```

Problems can include:

- Increased latency
- Capacity bottlenecks
- Large failure domain
- Long network paths

A distributed architecture can use multiple regional/edge sites.

```text
                        Anycast Service IP
                               |
                         Internet / BGP
                       /       |       \
                      /        |        \
                     v         v         v
                 India LB   Europe LB   US LB
                    |          |          |
                  / | \      / | \      / | \
                 S1 S2 S3   S4 S5 S6   S7 S8 S9
```

Anycast can help select a **site**.

The load balancer at that site then selects a **backend server**.

```text
User
 |
 v
Anycast routing
 |
 v
Regional Load Balancer
 |
 +----> Backend 1
 |
 +----> Backend 2
 |
 +----> Backend 3
```

This is an important distinction:

```text
Anycast/BGP
    |
    +---- chooses a reachable/preferred SITE
                 |
                 v
          Load Balancer
                 |
                 +---- chooses a BACKEND
```

---

# 9. Anycast for DNS

DNS is one of the classic Anycast use cases.

Imagine a DNS service available from many sites.

```text
                        DNS Service IP
                              |
                       Internet / BGP
                      /       |       \
                     v        v        v
                  India    Europe     USA
                   DNS       DNS       DNS
```

A client sends a normal DNS request to the service IP.

The client does not need to manually choose a region. Routing directs the traffic toward one of the sites advertising the service prefix.

Benefits can include:

- Lower latency
- Geographic distribution
- Failure isolation
- Large aggregate capacity
- DDoS resilience when engineered appropriately

---

# 10. Anycast for CDNs and Edge Networks

A CDN may operate many edge locations.

```text
User
 |
 v
Anycast / global routing
 |
 v
Edge Site
 |
 +--> Cached content
 |
 +--> Regional service
 |
 +--> Origin, if required
```

The objective is to bring service processing/content closer in network terms to users and distribute traffic across infrastructure.

Anycast is one possible component of global traffic engineering; DNS-based steering and other mechanisms are also widely used.

---

# 11. Comparison Table

| Model | Meaning | Receiver Model | TCP | UDP | Typical Examples |
|---|---|---|---|---|---|
| Unicast | One-to-one | One specific endpoint | Yes | Yes | HTTPS, SSH, DNS |
| Broadcast | One-to-all in local scope | All hosts in broadcast domain | No normal TCP broadcast | Yes for IPv4 broadcast | DHCPv4; ARP uses L2 broadcast |
| Multicast | One-to-group | Joined/interested group | Traditional TCP: No | Yes | IPTV, mDNS; OSPF uses IP multicast |
| Anycast | One-to-one-of-many | One site selected by routing | Yes | Yes | Global DNS, edge/CDN, global LB entry |

---

# 12. Protocol and Technology Reference

| Protocol / Technology | Cast/Delivery Model | TCP/UDP/Other | Purpose |
|---|---|---|---|
| HTTP/1.1 / HTTP/2 | Usually Unicast | TCP | Web/API communication |
| HTTP/3 | Usually Unicast; can reach Anycast service IP | QUIC over UDP | Modern web transport |
| SSH | Unicast | TCP | Remote login/administration |
| DNS | Usually Unicast; services may use Anycast | UDP and TCP | Name resolution |
| DHCPv4 | Broadcast and Unicast depending on phase | UDP | Automatic IPv4 configuration |
| ARP | Ethernet Broadcast request | Other | IPv4-to-MAC resolution on LAN |
| mDNS | Multicast | UDP | Local name/service discovery |
| IGMP | Multicast control | Other/IP | IPv4 multicast group management |
| OSPF | Unicast/Multicast depending on context | Other/IP | Interior routing |
| BGP | Unicast sessions | TCP | Exchange routing information |
| Anycast architecture | Anycast | Can carry TCP, UDP, etc. | Route clients to one of multiple service sites |
| IPTV | Multicast in many managed deployments | Commonly UDP-based | Efficient one-to-many media distribution |

---

# 13. TCP vs UDP vs Cast Type

The following question is misleading:

> "Does UDP do multicast?"

A better explanation is:

> UDP can be transported in an IP packet whose destination is a multicast address.

Similarly:

```text
UDP + Unicast
UDP + IPv4 Broadcast
UDP + Multicast
UDP + Anycast destination

TCP + Unicast
TCP + Anycast destination
```

Traditional TCP does not provide ordinary broadcast or multicast because a TCP connection is defined between two endpoints.

---

# 14. Example: UDP Unicast DNS

```text
Application
   |
   | DNS Query
   v
UDP
Source Port: ephemeral
Destination Port: 53
   |
   v
IP
Source: Client IP
Destination: DNS Server IP
   |
   v
Network
```

Diagram:

```text
Client                           DNS Resolver

DNS query
   |
   +--------- UDP/IP ----------->

              ONE TO ONE
```

---

# 15. Example: DHCP Broadcast

```text
                  DHCPDISCOVER
                       |
                       v
Client ----------> Broadcast
                       |
             +---------+---------+
             |         |         |
             v         v         v
           Host A    Host B    DHCP Server
```

Purpose:

The client needs network configuration and may not initially know which DHCP server is available.

---

# 16. Example: Multicast Video

Without multicast, a sender/application may need separate unicast flows:

```text
Server ---- Stream ----> User A
Server ---- Stream ----> User B
Server ---- Stream ----> User C
```

With multicast in a multicast-enabled network:

```text
                       +----> User A
                       |
Server ---> Multicast -+----> User B
                       |
                       +----> User C
```

The network can replicate packets at appropriate branching points.

This can be efficient for one-to-many delivery in managed multicast networks.

---

# 17. Example: Global Load Balancing with Anycast

```text
                           Global Service IP
                                  |
                                  v
                           Internet Routing
                     _____________|_____________
                    /             |             \
                   /              |              \
                  v               v               v
             Asia Site       Europe Site       US Site
                 |                |                |
                 v                v                v
              L4/L7 LB         L4/L7 LB         L4/L7 LB
              / |  \            / | \            / | \
             S1 S2 S3          S4 S5 S6          S7 S8 S9
```

Two decisions occur:

### Decision 1 — Global routing

```text
Which SITE should receive this traffic?
```

Anycast/BGP can participate here.

### Decision 2 — Local load balancing

```text
Which BACKEND should process the request?
```

The regional load balancer makes this decision.

---

# 18. L4 vs L7 Load Balancing

Anycast should not be confused with L4 or L7 load balancing.

## Layer 4 Load Balancer

Makes decisions primarily using transport/network information such as:

```text
Source IP
Destination IP
Source Port
Destination Port
Protocol
```

Conceptually:

```text
Client
  |
  v
L4 Load Balancer
  |
  +----> Server A
  +----> Server B
```

## Layer 7 Load Balancer

Understands application protocols such as HTTP.

It may route based on:

```text
Hostname
URL path
HTTP headers
Cookies
Other application metadata
```

Example:

```text
                 L7 Load Balancer
                       |
          +------------+------------+
          |                         |
          v                         v
      /users/*                  /orders/*
          |                         |
          v                         v
    User Service              Order Service
```

A global system can combine all of these:

```text
Anycast / DNS steering
          |
          v
Regional L4/L7 Load Balancer
          |
          v
Application Servers
```

---

# 19. Broadcast vs Multicast

## Broadcast

```text
Sender
   |
   +----> A
   +----> B
   +----> C
   +----> D

Everyone in the broadcast domain
```

## Multicast

```text
Sender
   |
   +----> A ✓
   +----> B ✓
   X----> C
   +----> D ✓

Only the targeted multicast group
```

Multicast is therefore more selective than broadcast.

---

# 20. Multicast vs Anycast

This is one of the most important distinctions.

## Multicast

```text
One sender
    |
    +----> Receiver A
    +----> Receiver B
    +----> Receiver C
```

**One to many group members.**

## Anycast

```text
                    +---- Site A
                    |
Client -------------+---- Site B ✓
                    |
                    +---- Site C
```

**One to one selected from many possible sites.**

---

# 21. Broadcast vs Anycast

Broadcast:

```text
One --> Everyone in local broadcast domain
```

Anycast:

```text
One --> One selected service location
```

They solve completely different problems.

Broadcast is mainly a local-network delivery mechanism.

Anycast is commonly used for distributed routing and global service architectures.

---

# 22. IPv6 Important Difference

IPv6 does not use broadcast in the same way IPv4 does.

Instead, IPv6 relies heavily on multicast.

For example, IPv6 Neighbor Discovery replaces the role that ARP serves for IPv4 neighbor resolution.

```text
IPv4:
ARP + Ethernet broadcast

IPv6:
Neighbor Discovery + ICMPv6 multicast
```

This avoids the need for an IPv6 broadcast address.

---

# 23. Where Does ICMP Fit?

ICMP is neither TCP nor UDP.

It is a network-layer control/diagnostic protocol.

Examples:

```text
ping
traceroute-related responses
Destination Unreachable
Time Exceeded
```

A normal ping is conceptually:

```text
Host A ----- ICMP Echo Request -----> Host B
Host A <---- ICMP Echo Reply -------- Host B
```

Usually this is unicast.

---

# 24. Where Does ARP Fit?

ARP is also neither TCP nor UDP.

ARP operates between Layer 2 and Layer 3 concepts for IPv4 LAN communication.

```text
IP Address
    |
    | ARP
    v
MAC Address
```

Example:

```text
"Who has 192.168.1.20?"
```

The request is normally sent using an Ethernet broadcast destination.

---

# 25. Where Does BGP Fit?

BGP is a routing protocol.

BGP peers establish TCP connections, conventionally using port:

```text
TCP 179
```

BGP exchanges routing information:

```text
Router A <------ TCP/BGP ------> Router B
```

Anycast deployments often depend on BGP advertisements:

```text
Site A ---- advertises prefix ----\
                                   \
                                    Internet
                                   /
Site B ---- advertises prefix ----/
```

The same service prefix can therefore be reachable through multiple sites.

---

# 26. System Design Example

Suppose we are designing:

```text
api.example.com
```

with users worldwide.

One possible architecture:

```text
User
 |
 | DNS lookup
 v
DNS
 |
 | Service address
 v
Anycast / Global Traffic Layer
 |
 +-----------------------------+
 |                             |
 v                             v
Asia Region                Europe Region
 |                             |
 v                             v
Load Balancer              Load Balancer
 |                             |
 +----> API 1                  +----> API 4
 +----> API 2                  +----> API 5
 +----> API 3                  +----> API 6
```

The technologies solve different problems:

| Component | Responsibility |
|---|---|
| DNS | Convert hostname to address and potentially participate in traffic steering |
| Anycast/BGP | Route toward one of multiple service sites |
| Load Balancer | Distribute traffic among backends |
| TCP/UDP/QUIC | Transport application traffic |
| HTTP/HTTPS | Application protocol |
| Backend | Execute application logic |

---

# 27. What If the Load Balancer Becomes Busy?

A single load balancer can itself become:

- A bottleneck
- A failure point
- A capacity limitation

Instead of:

```text
Millions of users
       |
       v
    ONE LB
       |
       v
   Servers
```

large systems can distribute the load-balancing layer:

```text
                         Global Entry
                              |
                    +---------+---------+
                    |         |         |
                    v         v         v
                   LB1       LB2       LB3
                    |         |         |
                    +---- Backend Pool -+
```

At global scale:

```text
                         Anycast/DNS
                             |
              +--------------+--------------+
              |              |              |
              v              v              v
          Asia LBs       Europe LBs       US LBs
              |              |              |
              v              v              v
           Servers        Servers        Servers
```

Techniques include:

- Multiple load balancers
- High availability
- Horizontal scaling
- L4 load balancing
- L7 load balancing
- DNS-based traffic steering
- Anycast
- Regional architectures
- Health checks
- Autoscaling

---

# 28. Failure Handling with Anycast

Suppose:

```text
              Same service prefix
                /      |      \
               v       v       v
            India    Europe    USA
```

If the India site becomes unavailable, an operator may withdraw its route advertisement.

Conceptually:

```text
India   X  no longer advertises route
Europe  ✓  advertises route
USA     ✓  advertises route
```

After routing converges, traffic that would have used the India route can follow another available route.

Actual failover behavior depends on:

- BGP convergence
- Routing policy
- Health-check design
- Route withdrawal mechanisms
- Provider architecture

Anycast therefore contributes to resilience, but it does not automatically make an application healthy.

---

# 29. Security Considerations

Each model has different operational considerations.

## Broadcast

Too much broadcast traffic can consume LAN resources.

Network segmentation using VLANs and routers reduces broadcast-domain size.

## Multicast

Requires correct multicast routing and group-management configuration.

Common mechanisms include:

- IGMP
- IGMP snooping
- Multicast routing protocols

## Anycast

Requires careful routing and service-state design.

Potential concerns include:

- Route leaks
- Route hijacking
- Session behavior when routing changes
- BGP convergence
- State synchronization
- Health-based route withdrawal

For stateful applications, engineers must think carefully about what happens if traffic moves between sites.

---

# 30. Quick Memory Trick

```text
UNICAST
UNI = ONE

A ----------> B

One to One
```

```text
BROADCAST
BROAD = EVERYONE IN LOCAL BROADCAST SCOPE

          +--> B
A --------+--> C
          +--> D

One to All
```

```text
MULTICAST
MULTI = GROUP

          +--> B ✓
A --------+--> C ✓
          X--> D

One to Group
```

```text
ANYCAST
ANY = ONE OF MULTIPLE POSSIBLE INSTANCES

          +--- B
A --------+--> C ✓
          +--- D

One to One-of-Many
```

---

# 31. Interview / Training Questions

### Q1. Can UDP use unicast?

Yes.

Example: a normal DNS request to a specific DNS resolver.

### Q2. Can UDP use broadcast?

Yes, with IPv4 broadcast.

Example: DHCPv4 discovery.

### Q3. Can UDP use multicast?

Yes.

Examples include mDNS and many multicast streaming applications.

### Q4. Can UDP be sent to an Anycast address?

Yes.

A DNS service may expose an Anycast address while clients send ordinary UDP DNS requests to it.

### Q5. Can TCP broadcast?

Traditional TCP does not support broadcast communication. A TCP connection is point-to-point.

### Q6. Can TCP multicast?

Traditional TCP does not provide IP multicast semantics.

### Q7. Can TCP work with Anycast?

Yes, but long-lived/stateful connections require careful Anycast architecture because routing changes can direct new packets/connections differently.

### Q8. Is Anycast a protocol?

No.

It is an addressing/routing technique.

### Q9. Is BGP Anycast?

No.

BGP is a routing protocol commonly used to implement Internet-scale Anycast by advertising the same prefix from multiple sites.

### Q10. Is ARP UDP?

No.

ARP is a separate protocol and its IPv4 LAN request is normally carried in an Ethernet broadcast frame.

### Q11. Is OSPF UDP?

No.

OSPF runs directly over IP.

### Q12. Is IGMP UDP?

No.

IGMP is an IP-layer protocol used for IPv4 multicast group management.

---

# 32. Final Summary

The most important idea is to separate **transport** from **delivery/routing**.

```text
                 APPLICATION
       HTTP / DNS / DHCP / SSH / mDNS
                      |
                      v
                  TRANSPORT
                 TCP / UDP
                      |
                      v
               IP / NETWORKING
                      |
       +--------------+--------------+
       |              |              |
    Unicast       Multicast       Anycast
       |
 IPv4 broadcast is also available in
 appropriate local-network scenarios
```

In practical terms:

```text
Unicast
One sender -> One destination
Examples: HTTPS, SSH, ordinary DNS

Broadcast
One sender -> Everyone in local broadcast domain
Examples: DHCPv4 discovery; ARP uses Ethernet broadcast

Multicast
One sender -> Selected group
Examples: IPTV, mDNS, multicast routing/control

Anycast
One sender -> One selected site among many
Examples: Global DNS, CDN/edge entry, global load-balancing architectures
```

For system design, a particularly useful architecture is:

```text
Users
  |
  v
DNS / Anycast Global Entry
  |
  v
Regional Site
  |
  v
L4/L7 Load Balancer
  |
  v
Backend Services
```

Each component solves a different problem:

- **DNS** — name resolution and sometimes traffic steering.
- **Anycast/BGP** — global route selection toward a service site.
- **Load balancer** — backend selection within a site.
- **TCP/UDP/QUIC** — transport.
- **HTTP/DNS/etc.** — application protocol.
- **Backend service** — business logic.

Understanding these boundaries is essential for networking, Kubernetes, cloud architecture, distributed systems, and system design.
