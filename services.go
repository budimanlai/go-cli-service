package services

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	goargs "github.com/budimanlai/go-args"
	goconfig "github.com/budimanlai/go-config"
	goping "github.com/budimanlai/go-services-ping"
	"github.com/eqto/dbm"
	_ "github.com/eqto/dbm/driver/mysql"
)

type Service struct {
	Version    string
	AppName    string
	IsStopped  bool
	starFunc   ServiceHandler
	stopFunc   ServiceHandler
	Args       *goargs.Args
	Config     *goconfig.Config
	Db         *dbm.Connection
	LogService *LogService
	Ping       *goping.ServicePing

	configFile []string
}

type ServiceHandler func(ctx *Service)

func NewService(configFile ...string) *Service {
	srv := &Service{
		configFile: configFile,
	}
	srv.Config = &goconfig.Config{}
	e := srv.Config.Open(srv.configFile...)

	if e != nil {
		panic(e)
	}

	srv.Args = &goargs.Args{}
	srv.Args.Parse()

	srv.LogService = NewLogService(srv.Args.ScriptName)
	e1 := srv.LogService.Init()
	if e1 != nil {
		panic(e1)
	}

	return srv
}

// Start the service
func (s *Service) Start() error {

	switch s.Args.Command {
	case "v":
	case "version":
		fmt.Println(s.AppName, "\nVersi", s.Version)
		break
	case "run":
		e1 := s.run()
		if e1 != nil {
			return e1
		}
		break
	case "start":

		a := s.Args.GetRawArgs()
		a[0] = "run"
		command := exec.Command(`./`+s.Args.ScriptName, a...)
		outfile, err := os.OpenFile(s.LogService.GetLogFile(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer outfile.Close()
		command.Stdout = outfile

		if e := command.Start(); e != nil {
			return e
		}

		pid := command.Process.Pid

		f, err := os.Create(s.LogService.GetPidFile())
		if err != nil {
			return err
		}

		defer f.Close()
		_, err2 := f.WriteString(strconv.Itoa(pid))
		if err2 != nil {
			return err2
		}

		break
	case "stop":
		pid, err := ioutil.ReadFile(s.LogService.GetPidFile())

		if err != nil {
			return err
		}

		cmd := exec.Command("kill", string(pid))
		cmd.Run()

		// delete pid file
		e1 := os.Remove(s.LogService.GetPidFile())
		if e1 != nil {
			return e1
		}
		break
	}

	return nil
}

func (s *Service) SetDatabase(db *dbm.Connection) {
	s.Db = db
}

func (s *Service) openDatabase() error {
	if s.Db != nil {
		return nil
	}
	cn, e := dbm.Connect("mysql", s.Config.GetString("database.hostname"), s.Config.GetInt("database.port"),
		s.Config.GetString("database.username"), s.Config.GetString("database.password"), s.Config.GetString("database.database"))
	if e != nil {
		return e
	}
	s.Db = cn
	return nil
}
func (s *Service) run() error {
	e1 := s.openDatabase()
	if e1 != nil {
		return e1
	}

	s.Ping = &goping.ServicePing{
		Config:      s.Config,
		Indentifier: s.Args.ScriptName,
	}
	s.Ping.OpenDatabase()
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
	s.LogService.Log(a...)
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
