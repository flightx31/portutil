package portutil

import (
	"errors"
	"github.com/flightx31/exception"
	"net"
	"strconv"
)

type Logger interface {
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
	Trace(args ...interface{})
	Print(args ...interface{})
}

var log Logger

func SetLogger(l Logger) {
	log = l
}

// FindOpenPort starts at starting port defined in config, and tries all ports sequentially up to ports to try quantity
// value in config. It will return net.Listener, and net.PacketConn both bound to the open port found.
// Note: returning the listener, and packet connection so that we can hold onto the port until we are ready to use it.
// They will have to be closed in the calling functions.
func FindOpenPort(startingPort int, portsToTry int) (PortConnection, error) {
	// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go

	// TODO learn how to gain access to the network.
	// Heres how it's done in windows: https://www.programmersought.com/article/8831210069/

	connectionDetails := PortConnection{}
	connectionDetails.OurIpAddress = GetOutboundIP()
	//env.OurIpAddress = connectionDetails.OurIpAddress

	//startingPort := env.Config.PortRangeBegin
	//portsToTry := env.Config.PortTryQuantity

	// look for available port
	for c := 0; c < portsToTry; c++ {
		port := startingPort + c
		portString := strconv.Itoa(port)

		tcpListener, err := net.Listen("tcp", ":"+portString)

		if err == nil {

			udpListener, err := net.ListenPacket("udp4", ":"+portString)

			if err == nil {
				log.Info("Binding port: ", port)
				connectionDetails.TCPConnection = tcpListener
				connectionDetails.UDPConnection = udpListener
				connectionDetails.Port = port
				connectionDetails.PortString = portString
				return connectionDetails, nil
			} else {
				log.Trace("Need both udp and tcp to be available on the port. Skipping this port...")
				exception.PanicOnError(tcpListener.Close())
				continue
			}
		}
	}
	return PortConnection{}, errors.New("cannot find port to bind for both tcp and udp")
}

// GetOutboundIP - Gets preferred outbound ip of this machine obtained this from: https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer exception.PanicOnErrorFunc(conn.Close)

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

type PortConnection struct {
	TCPConnection net.Listener
	UDPConnection net.PacketConn
	Port          int
	PortString    string
	OurIpAddress  net.IP
}
