package server

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"serialization/internal/logger"
	logger_pkg "serialization/internal/logger"
	"serialization/serializers"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/ipv4"
)

const (
	BUFFER_SIZE = 1024
)

type Server struct {
	port          string
	dataFormat    string
	isMulticast   bool
	udpAddr       *net.UDPAddr
	groupUdpAddr  *net.UDPAddr
	listeningConn *net.UDPConn
	serializer    serializers.Serializer
	logger        *slog.Logger
}

var portByServerType = map[string]string{
	"proxy":        "2000",
	"native":       "2001",
	"xml":          "2002",
	"json":         "2003",
	"protobuf":     "2004",
	"avro":         "2005",
	"yaml":         "2006",
	"message_pack": "2007",
}

func GetSerializer(serverType string) serializers.Serializer {
	switch serverType {
	case "native":
		return &serializers.NativeSerializer{}
	case "xml":
		return &serializers.XmlSerializer{}
	case "json":
		return &serializers.JsonSerializer{}
	case "protobuf":
		return &serializers.ProtobufSerializer{}
	case "avro":
		return &serializers.AvroSerializer{}
	case "yaml":
		return &serializers.YamlSerializer{}
	case "message_pack":
		return &serializers.MessagePackSerializer{}
	case "proxy":
		return nil
	}
	err_msg := fmt.Sprintf("Неизвестный формат данных '%s'.\n", serverType)
	err_msg += "Доступные форматы: native; xml; json; protobuf; avro; yaml; message_pack."
	log.Fatal(err_msg)
	return nil
}

func (s *Server) resolve(address string) *net.UDPAddr {
	s.logger.Debug(fmt.Sprintf("want to resolve address %s", address))
	result, err := net.ResolveUDPAddr("udp", address)
	logger_pkg.FailOnError("failed to resolve udp addr", err)
	return result
}

func (s *Server) runBenchmark(destination *net.UDPAddr) string {
	conn, err := net.DialUDP("udp", nil, destination)
	logger.FailOnError("failed to set listen conn", err)
	defer conn.Close()

	request := "get_result"
	_, err = conn.Write([]byte(request))
	logger.FailOnError("failed to write data", err)
	s.logger.Info(fmt.Sprintf("[%s -> %s] send msg: '%s'", s.udpAddr.String(), destination.String(), request))

	buf := make([]byte, BUFFER_SIZE)
	s.logger.Debug("init buffer for benchmark results")

	benchmarkResults := ""
	if destination.String() == s.groupUdpAddr.String() {
		for i := 1; i < len(portByServerType); i++ {
			n, _, err := s.listeningConn.ReadFromUDP(buf)
			logger.FailOnError("failed to read data", err)
			benchmarkResults += strings.TrimSpace(string(buf[:n])) + "\n"
		}
	} else {
		n, _, err := s.listeningConn.ReadFromUDP(buf)
		logger.FailOnError("failed to read data", err)
		benchmarkResults += strings.TrimSpace(string(buf[:n])) + "\n"
	}

	s.logger.Debug(fmt.Sprintf("benchmark results: '%s'", benchmarkResults))
	return benchmarkResults
}

func (s *Server) IsProxy() bool {
	return s.port == "2000" && s.serializer == nil
}

func (s *Server) handleProxyMsg(dataFormat string, sender *net.UDPAddr) {
	s.logger.Debug(fmt.Sprintf("handle proxy dataFormat='%s'", dataFormat))

	var destination *net.UDPAddr
	if dataFormat == "all" {
		destination = s.groupUdpAddr
	} else {
		port, ok := portByServerType[dataFormat]
		if !ok {
			log.Fatalf("unknown data format '%s'", dataFormat)
		}
		destination = s.resolve(dataFormat + ":" + port)
	}
	s.logger.Debug(fmt.Sprintf("%s benchmark request to addr %s", dataFormat, destination.String()))

	benchmarkResults := s.runBenchmark(destination)
	s.listeningConn.WriteToUDP([]byte(benchmarkResults), sender)
}

func (s *Server) calculateBenchmark() string {
	result := s.dataFormat + ":\n"

	s.serializer.PrepareData()

	startSerialization := time.Now()
	data := s.serializer.Serialize()
	serializationDuration := time.Since(startSerialization)

	// result += fmt.Sprintf("\tsize of source data      \t%d bytes\n", serializers.GetDataSize())
	result += fmt.Sprintf("\tsize of serialized data  \t%d bytes\n", len(data))

	result += fmt.Sprintf("\tserialization duration   \t%dµs\n", serializationDuration.Microseconds())

	startDeserialization := time.Now()
	s.serializer.Deserialize(data)
	deserializationDuration := time.Since(startDeserialization)

	result += fmt.Sprintf("\tdeserialization duration \t%dµs\n\n", deserializationDuration.Microseconds())

	s.serializer.CheckResult()
	return result
}

func (s *Server) handleSerializerMsg(sender *net.UDPAddr) {
	s.logger.Debug("start handle serializer msg", "is_multicast_addr", s.isMulticast)

	sendMsg := s.calculateBenchmark()
	proxyAddr := s.resolve("proxy:2000")
	_, err := s.listeningConn.WriteToUDP([]byte(sendMsg), proxyAddr)
	logger.FailOnError("failed to write data", err)

	s.logger.Info(fmt.Sprintf("[%s -> %s] send msg: '%s'", s.udpAddr.String(), proxyAddr.String(), sendMsg))
}

func (s *Server) checkAndCorrectServerMsg(msg []byte) string {
	receivedMsg := strings.TrimSpace(string(msg))
	s.logger.Info(fmt.Sprintf("received msg: '%s'", receivedMsg), "is_multicast_addr", s.isMulticast)

	if len(receivedMsg) < 10 || receivedMsg[:10] != "get_result" {
		log.Fatal("Сообщение должно начинаться с 'get_result'")
	}

	if !s.IsProxy() && len(receivedMsg) > 10 {
		log.Fatal("Ожидаем сообщение 'get_result'")
	}

	return receivedMsg
}

func (s *Server) runHandler() {
	buf := make([]byte, BUFFER_SIZE)
	s.logger.Debug("init buffer to read command", "is_multicast_addr", s.isMulticast)

	for {
		s.logger.Debug("Ожидаем получения сообщения...", "is_multicast_addr", s.isMulticast)
		n, addr, err := s.listeningConn.ReadFromUDP(buf)
		logger.FailOnError("failed to read data", err)
		receivedMsg := s.checkAndCorrectServerMsg(buf[:n])

		if s.IsProxy() {
			dataFormat := strings.TrimSpace(receivedMsg[10:])
			s.handleProxyMsg(dataFormat, addr)
		} else {
			s.handleSerializerMsg(addr)
		}
	}
}

func (s *Server) RunServer(wg *sync.WaitGroup) {
	defer wg.Done()

	var err error
	if s.isMulticast {
		// use ListenUDP because ListenMulticastUDP disables IP_MULTICAST_LOOP
		// https://stackoverflow.com/questions/43109552/how-to-set-ip-multicast-loop-on-multicast-udpconn-in-golang
		s.listeningConn, err = net.ListenUDP("udp", s.groupUdpAddr)
		logger.FailOnError("failed to listen multicast udp", err)

		conn := ipv4.NewPacketConn(s.listeningConn)
		val, err := conn.MulticastLoopback()
		logger.FailOnError("failed get loopback for multicast packet", err)
		if !val {
			err := conn.SetMulticastLoopback(true)
			logger.FailOnError("failed set loopback for multicast packet", err)
		}
	} else {
		s.listeningConn, err = net.ListenUDP("udp", s.udpAddr)
		logger.FailOnError("failed to listen udp", err)
	}
	defer s.listeningConn.Close()

	s.runHandler()
}

func (s *Server) Init(serverType *string, logger *slog.Logger, isMulticast bool) {
	s.logger = logger
	port, ok := portByServerType[*serverType]
	if !ok {
		log.Fatalf("unknown server type %s", *serverType)
	}

	groupUdpAddrStr := os.Getenv("GROUP_UDP_ADDR")
	if groupUdpAddrStr == "" {
		log.Fatal("GROUP_UDP_ADDR is empty")
	}
	s.groupUdpAddr = s.resolve(groupUdpAddrStr)
	if !s.groupUdpAddr.IP.IsMulticast() {
		log.Fatal("GROUP_UDP_ADDR does not contain multicast address")
	}
	s.logger.Debug(fmt.Sprintf("group udp addr: %s", s.groupUdpAddr))

	s.port = port

	s.isMulticast = isMulticast
	if isMulticast {
		if *serverType == "proxy" {
			log.Fatal("proxy does not listen to multicast address")
		}

		s.udpAddr = s.groupUdpAddr
	} else {
		s.udpAddr = s.resolve(*serverType + ":" + s.port)
	}

	s.serializer = GetSerializer(*serverType)
	if s.serializer != nil {
		s.dataFormat = *serverType
	}

	logger.Info(fmt.Sprintf("init addr: %s", s.udpAddr.String()), "is_multicast_addr", s.isMulticast)
	logger.Info(fmt.Sprintf("init server_type: %s", *serverType), "is_multicast_addr", s.isMulticast)
}
