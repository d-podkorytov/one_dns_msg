// Speed 112922 qps
//       109598 qps on single core
package main

import (
	"flag"
	"fmt"
	"net"
	"log"
//        "golang.org/x/sys/unix"
//        "syscall/syscall_linux" 
        "syscall"
        "unsafe"
)

var port = flag.Int("p", 53, "The listening port")
var udpPackageBufferSize = flag.Int("l", 1024, "The size of the udp package buffer")

func main() {
	flag.Parse()

	// open socket
	//socket, err := net.ListenUDP("udp4", &net.UDPAddr{
	//	IP:   net.IPv4(0, 0, 0, 0),
	//	Port: *port,
	//})

         fd , err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP) //(domain, typ, proto int) 
	if err != nil { fmt.Println("open socket failed!", err)
		        return }
        
        
        server_addr := syscall.SockaddrInet4{Port: 53}
        copy(server_addr.Addr[:], net.ParseIP("0.0.0.0").To4())
         
        err_bind:=syscall.Bind(fd, &server_addr)
        
        if err_bind != nil { fmt.Println("bind socket failed!", err_bind)
		        return }
        
        log.Println("listen ", syscall.PACKET_HOST)

//	defer socket.Close()

	for {
		// endless loop wait clients asks
		oob := make([]byte, *udpPackageBufferSize)
		p := make([]byte, *udpPackageBufferSize)
		//readn, remoteAddr, _ := socket.ReadFromUDP(ask_data)
                //
                 n, oobn , recvflags , from , err  := syscall.Recvmsg(fd, p, oob, syscall.MSG_WAITALL )  
                 fmt.Println(" recvmsg= ", n, oobn , recvflags , from , err )
                 fmt.Println(" p= ", p )                
                 //for fast working dont control returned error
		//if err != nil { fmt.Println("recvfrom error!", err)
		//	        continue }
		
		go process(fd, p[:n], oob,from)
	}
}

func process(fd int, ask_data []byte,oob []byte, remoteAddr syscall.Sockaddr ) { //syscall.Sockaddr
// here do reply on ask data

fmt.Println("Process data ",ask_data)

                id:=make([]byte,2)

                id[0]=ask_data[0]
                id[1]=ask_data[1]

                rr_record := []byte{0,0,1,0,0,1,0,1,0,0,0,0,4,49,48,48,48,3,100,105,112,0,0,1,0,1,192,12,0,1,0,1,0,0,8,212,0,4,87,118,90,81}

                rr_record[0]=id[0]
                rr_record[1]=id[1]

// conn.WriteToUDP(rr_record, remoteAddr)
// SendmsgN(fd int, p, oob []byte, to Sockaddr, flags int) (n int, err error)
err := Sendmsg(fd, ask_data, oob , remoteAddr, syscall.MSG_WAITALL ) 
fmt.Println("sendmsg err= ",err)

//	_, err := conn.WriteToUDP(rr_record, remoteAddr)
//	if err != nil { fmt.Println("send data error", err) }
}

//func SendmsgN(fd, ask_data, oob , remoteAddr, syscall.MSG_WAITALL ) 

//func Sendmsg(fd int, p, oob []byte, to syscall.Sockaddr, flags int) (err error) {
//	var ptr   uintptr
//	var salen _Socklen
//	if to != nil {
//		var err error
//		ptr, salen, err = to.sockaddr()
//		if err != nil {
//			return err
//		}
//	}
//	var msg syscall.Msghdr
//	msg.Name = (*byte)(unsafe.Pointer(ptr))
//	msg.Namelen = uint32(salen)
//	var iov syscall.Iovec
//	if len(p) > 0 {
//		iov.Base = (*byte)(unsafe.Pointer(&p[0]))
//		iov.SetLen(len(p))
//	}
//	var dummy byte
//	if len(oob) > 0 {
//		// send at least one normal byte
//		if len(p) == 0 {
//			iov.Base = &dummy
//			iov.SetLen(1)
//		}
//		msg.Control = (*byte)(unsafe.Pointer(&oob[0]))
//		msg.SetControllen(len(oob))
//	}
//	msg.Iov = &iov
//	msg.Iovlen = 1
//	if err = sendmsg(fd, &msg, flags); err != nil {
//		return
//	}
//	return
//}
