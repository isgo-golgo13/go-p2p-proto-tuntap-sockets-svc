## P2P Tunnel Service (Lightweight VPN)

###### p2p_tunelsvc is P2P/UDP peer-to-peer server pair that receives TUN packets and send TUN packets (along with any enclosed private data in the header) and a receiver peer that receives the TUN packets (reads/writes to the packet and resends the TUN packet out again on TUN device layer (virtual Layer 3).  

#### Cloning the Project

`git clone https://github.com/7Tunnels/network-poc/blob/master/p2p_tunnelsvc`

#### Compiling the Project

`make`  # make will automatically execute Go dep esnure -v to downooad all project package dependencies

#### Running the Client

#### Running the Server

