package hotstuff

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/sirupsen/logrus"
	"github.com/theodocius271/hotstuff/logging"
	pb "github.com/theodocius271/hotstuff/proto"
	"google.golang.org/grpc"
)

// Server 封装了 HotStuff 服务器的功能
type Server struct {
	service       *HotStuffService
	grpcServer    *grpc.Server
	address       string
	logger        *logrus.Logger
	sigChan       chan os.Signal
	shutdownChan  chan struct{}
	shuttingDown  bool
	shutdownMutex sync.Mutex
}

// Config 包含服务器配置信息
type Config struct {
	Address     string        // 服务地址，格式为 ip:port
	IdleTimeout time.Duration // 空闲连接超时时间
}

// NewServer 创建一个新的 HotStuff 服务器
func NewServer(id uint32, support consensus.ConsenterSupport) *Server {

	server := &Server{
		logger:       logging.GetLogger(),
		sigChan:      make(chan os.Signal, 1),
		shutdownChan: make(chan struct{}),
		shuttingDown: false,
	}

	// 创建 HotStuff 服务
	server.service = NewHotStuffService(NewBasicHotStuff(id, support))
	server.address = server.service.hotStuff.GetSelfInfo().Address

	// 创建 gRPC 服务器，添加拦截器用于日志记录和恢复
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(server.unaryInterceptor),
	}
	server.grpcServer = grpc.NewServer(opts...)

	// 注册 HotStuff 服务
	pb.RegisterHotStuffServiceServer(server.grpcServer, server.service)

	return server
}

// Logging and Recovery
func (s *Server) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	defer func() {
		if r := recover(); r != nil {
			s.logger.Errorf("Panic recovered in gRPC handler: %v", r)
		}
	}()
	resp, err := handler(ctx, req)
	s.logger.Debugf("gRPC method %s took %s", info.FullMethod, time.Since(start))

	return resp, err
}

// Run Start Server and block sig
func (s *Server) Run() error {

	signal.Notify(s.sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go s.handleSignals()

	// lis, err := net.Listen("tcp", s.address)
	port := s.address[strings.Index(s.address, ":"):]
	s.logger.Infof("[HOTSTUFF] Server start at port%s", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", s.address, err)
	}

	s.logger.Infof("HotStuff server starting on %s", s.address)

	err = s.grpcServer.Serve(lis)
	if err != nil && !s.shuttingDown {
		return fmt.Errorf("gRPC server failed: %v", err)
	}
	s.logger.Info("HotStuff server stopped")
	return nil
}

// RunAsync 异步启动服务器并立即返回
func (s *Server) RunAsync() error {

	signal.Notify(s.sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go s.handleSignals()

	//lis, err := net.Listen("tcp", s.address)
	port := s.address[strings.Index(s.address, ":"):]
	s.logger.Infof("[HOTSTUFF] Server start at port%s", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", s.address, err)
	}

	s.logger.Infof("HotStuff server starting on %s", s.address)

	go func() {
		err := s.grpcServer.Serve(lis)
		if err != nil && !s.shuttingDown {
			s.logger.Errorf("gRPC server failed: %v", err)
		}
		s.logger.Info("HotStuff server stopped")
	}()

	return nil
}

func (s *Server) handleSignals() {
	for {
		select {
		case sig := <-s.sigChan:
			s.logger.Infof("Received signal: %v", sig)
			s.Shutdown()
			return
		case <-s.shutdownChan:
			return
		}
	}
}

// Gracefully shutdown server
func (s *Server) Shutdown() {
	s.shutdownMutex.Lock()
	if s.shuttingDown {
		s.shutdownMutex.Unlock()
		return
	}
	s.shuttingDown = true
	s.shutdownMutex.Unlock()

	s.logger.Info("Shutting down HotStuff server...")
	// 创建关闭完成通知通道
	done := make(chan struct{})

	go func() {
		s.grpcServer.GracefulStop()
		s.service.GetImpl().SafeExit()
		close(done)
	}()

	// 等待优雅关闭完成或超时
	select {
	case <-done:
		s.logger.Info("Server shutdown completed gracefully")
	case <-time.After(10 * time.Second):
		s.logger.Warn("Server shutdown timed out, forcing exit")
		s.grpcServer.Stop()
	}

	// 通知信号处理 goroutine 退出
	close(s.shutdownChan)
}

func (s *Server) GetHotStuffImpl() *BasicHotStuff {
	return s.service.GetImpl()
}

func (s *Server) GetAddress() string {
	return s.address
}
