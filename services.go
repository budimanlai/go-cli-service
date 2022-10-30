package services

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	goargs "github.com/budimanlai/go-args"
	goconfig "github.com/budimanlai/go-config"
	"github.com/eqto/dbm"
	_ "github.com/eqto/dbm/driver/mysql"
)

type Service struct {
	Version   string
	AppName   string
	IsStopped bool
	starFunc  ServiceHandler
	stopFunc  ServiceHandler
	Args      *goargs.Args
	Config    *goconfig.Config
	Db        *dbm.Connection

	configFile []string
}

type ServiceHandler func(ctx *Service)

const (
	YYYYMMDDHHMMSS = "2006-01-02 15:04:05"
)

func NewService(configFile ...string) *Service {
	return &Service{
		configFile: configFile,
	}
}

// Start the service
func (s *Service) Start() error {
	s.Args = &goargs.Args{}
	s.Args.Parse()

	switch s.Args.Command {
	case "v":
	case "version":
		fmt.Println(s.AppName, "\nVersi", s.Version)
		break
	case "run":
		e := s.run()
		if e != nil {
			return e
		}
		break
	case "start":
		fmt.Println("Start")
		break
	case "stop":
		fmt.Println("Stop")
		break
	}

	return nil
}

func (s *Service) openDatabase() error {
	cn, e := dbm.Connect("mysql", s.Config.GetString("database.hostname"), s.Config.GetInt("database.port"),
		s.Config.GetString("database.username"), s.Config.GetString("database.password"), s.Config.GetString("database.name"))
	if e != nil {
		return e
	}
	s.Db = cn
	return nil
}
func (s *Service) run() error {
	if len(s.configFile) != 0 {
		s.Config = &goconfig.Config{}
		e := s.Config.Open(s.configFile...)
		if e != nil {
			return e
		}

		e1 := s.openDatabase()
		if e1 != nil {
			return e1
		}
	}

	s.signalListener()
	s.IsStopped = false
	s.starFunc(s)

	return nil
}

func (s *Service) StartHandler(f ServiceHandler) {
	s.starFunc = f
}

func (s *Service) StopHandler(f ServiceHandler) {
	s.stopFunc = f
}

func (s *Service) Log(a ...interface{}) {
	now := time.Now().UTC()
	date := now.Format(YYYYMMDDHHMMSS)
	fmt.Print("[" + date + "] ")
	fmt.Println(a...)
}

func (s *Service) signalListener() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		s.Log(sig)
		s.stopFunc(s)
		done <- true
	}()
}
