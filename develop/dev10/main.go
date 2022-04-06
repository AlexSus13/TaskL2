package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"context"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

type args struct {
	host    string
	port    string
	timeout time.Duration
}

func main() {

	args, err := getArgs()
	if err != nil {
		log.Fatal(err)
	}

	err = telnetClient(args)
	if err != nil {
		log.Fatal(err)
	}
}

func telnetClient(args *args) error {

	ctx, cancel := context.WithCancel(context.Background())

	address := fmt.Sprintf("%s:%s", args.host, args.port)

	conn, err := net.DialTimeout("tcp", address, args.timeout)
	if err != nil {
		return fmt.Errorf("func net.DialTimeout: %v", err)
	}
	defer conn.Close()

	fmt.Println("Client is start: ", address)

	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		stopSignal := <-signalChan
		log.Println(stopSignal.String())
		cancel()
	}()

	go readFromSocket(conn, cancel)
	go writeToSocket(conn, cancel)

	<-ctx.Done()
	return nil
}

func readFromSocket(conn net.Conn, cancel context.CancelFunc) {

	scanner := bufio.NewScanner(conn)
	for {
		if !scanner.Scan() {
			log.Printf("соединение было прервано")
			cancel()
			return
		}
		text := scanner.Text()
		fmt.Printf("%s\n", text)
	}
}

func writeToSocket(conn net.Conn, cancel context.CancelFunc) {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			log.Println("CANNOT STDIN SCAN")
			cancel()
			return
		}
		str := scanner.Text()

		sl := strings.Split(fmt.Sprintf("% x", str), " ")
		for _, u := range sl {
			if u == "04" {
				log.Println("вы нажали ctrl+D. telnet клиент будет закрыт!")
				cancel()
				return
			}

		}
		_, err := conn.Write([]byte(fmt.Sprintln(str)))
		if err != nil {
			log.Println("ошибка при отправке на сервер", err)
			cancel()
			return
		}
	}
}

func getArgs() (*args, error) {

	if len(os.Args) < 3 {
		return nil, errors.New("No host or port specified")
	}

	if len(os.Args) > 4 {
		return nil, errors.New("Too many command line arguments")
	}

	var host string
	var port string
	var timeout time.Duration

	switch len(os.Args) {
	case 3:
		host = os.Args[1]
		port = os.Args[2]
		timeout = time.Second * 10
	case 4:
		if strings.Contains(os.Args[1], "--timeout=") {

			value := os.Args[1][len(os.Args[1])-1]

			index := strings.Index(os.Args[1], "=")
			timeValue, err := strconv.Atoi(os.Args[1][index+1 : len(os.Args[1])-1])

			if err != nil || value != 's' || timeValue < 1 {
				return nil, errors.New("Incorrect Args, Correct Args format --timeout=Ns, where N >= 1, int")
			}
			host = os.Args[2]
			port = os.Args[3]
			timeout = time.Duration(timeValue) * time.Second
		} else {
			return nil, errors.New("1st command line arguments are not equal to --timeout=")
		}
	}

	return &args{
		host:    host,
		port:    port,
		timeout: timeout,
	}, nil
}
