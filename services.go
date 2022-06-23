package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	config "github.com/eqto/config"
	"github.com/spf13/cobra"
)

type Service struct {
	rootCmd  *cobra.Command
	appName  string
	pid      int
	stopFunc ServiceHandler
	context  ServiceContext
}

type ServiceHandler func(context ServiceContext)

const ServiceRuntimePath = "runtime"

func NewService() Service {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}

	name := filepath.Base(path)
	return Service{
		rootCmd: &cobra.Command{Use: name},
		appName: name,
	}
}
func (s *Service) Start(startFunc ServiceHandler) {
	s.context = ServiceContext{}

	var cmdRun = &cobra.Command{
		Use:   "run",
		Short: "Run service and directly show the output",
		Long:  `Use this command if you want to see your application output directly.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			s.prepareStart()
			s.initDb()
			s.signalListener()

			startFunc(s.context)
		},
	}

	cmdRun.Flags().IntVarP(&s.context.Node, "node", "n", 1, "Node number")

	var cmdStart = &cobra.Command{
		Use:   "start",
		Short: "Start service in background.",
		Long:  `Start service in background.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			command := exec.Command(`./`+s.appName, "run", "--node="+strconv.Itoa(s.context.Node))
			outfile, err := os.OpenFile(s.LogFilename(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}
			defer outfile.Close()
			command.Stdout = outfile

			if e := command.Start(); e != nil {
				log.Fatal(e)
			}

			pid := command.Process.Pid

			f, err := os.Create(s.PidFilename())
			if err != nil {
				log.Fatal(err)
			}

			defer f.Close()
			_, err2 := f.WriteString(strconv.Itoa(pid))
			if err2 != nil {
				log.Fatal(err2)
			}
		},
	}
	cmdStart.Flags().IntVarP(&s.context.Node, "node", "n", 1, "Node number")

	s.rootCmd.AddCommand(cmdRun, cmdStart)
}
func (s *Service) Stop(stopFunc ServiceHandler) {
	s.stopFunc = stopFunc

	var cmdStop = &cobra.Command{
		Use:   "stop",
		Short: "Stop service",
		Long:  `Stop service`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			s.shutdown()
		},
	}

	cmdStop.Flags().IntVarP(&s.context.Node, "node", "n", 1, "Node number")

	s.rootCmd.AddCommand(cmdStop)
}

func (s *Service) SetVersion(version string) {
	s.rootCmd.Version = version
}
func (s *Service) Run() {
	s.rootCmd.Execute()
}
func (s *Service) LogFilename() string {
	return ServiceRuntimePath + "/logs/" + s.appName + ".log"
}
func (s *Service) PidFilename() string {
	return ServiceRuntimePath + "/" + s.appName + "-" + strconv.Itoa(s.context.Node) + ".pid"
}
func (s *Service) prepareStart() {
	if _, err := os.Stat(ServiceRuntimePath); os.IsNotExist(err) {
		if err := os.Mkdir(ServiceRuntimePath, os.ModePerm); err != nil {
			fmt.Println(err)
		}
	}

	if _, err := os.Stat(ServiceRuntimePath + "/logs"); os.IsNotExist(err) {
		if err := os.Mkdir(ServiceRuntimePath+"/logs", os.ModePerm); err != nil {
			fmt.Println(err)
		}
	}

	if e := config.Open(`config/main.conf`); e != nil {
		log.Fatal(e)
	}
}
func (s *Service) shutdown() {
	pid, err := ioutil.ReadFile(s.PidFilename())

	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("kill", string(pid))
	cmd.Run()

	// delete pid file
	e := os.Remove(s.PidFilename())
	if e != nil {
		log.Fatal(e)
	}
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
		fmt.Println()
		fmt.Println(sig)
		s.stopFunc(s.context)
		done <- true
	}()
}
func (s *Service) initDb() {
	OpenDb()
}
