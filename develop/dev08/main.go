/*
=== Взаимодействие с ОС ===
Необходимо реализовать собственный шелл
встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах
Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"os" //Пакет os предоставляет независимый от платформы интерфейс к функциям операционной системы.
	"os/exec" //Пакет os/exec запускает внешние команды.
	"strings"
)

const (
	pwd  = "pwd"
	cd   = "cd"
	echo = "echo"
	ps   = "ps"
	kill = "kill"
	exit = "exit"
)

func main() {

        scanner := bufio.NewScanner(os.Stdin) //Создаем новый сканер
        for scanner.Scan() { //В бесконечном цикле считываем команды из os.Stdin

		fullCommand := scanner.Text() //Получаем отсканированный текст.

		err := runCommand(fullCommand) //Запускаем выполнение команды.
		if err != nil {
			log.Fatalf("Error while run command: %s", err)
		}
	}
	//Возвращает ошибку при сканировании.
	err := scanner.Err()
	if err != nil {
		log.Fatalf("Error while reading Stdin: %s", err)
	}
}

func runCommand(fullCommand string) error {
	//Разделяем полную команду на части по пробелу (получим слайс частей команд).
	PartOfCommand := strings.Split(fullCommand, " ")
	//Если на входе пустая строка заканчиваем выполнение.
	if len(PartOfCommand) == 0 {
		return nil
	}

	switch PartOfCommand[0] { //Проверяем первую часть команды.
	case pwd:
		//Команда возвращает структуру Cmd для выполнения указанной программы с заданными аргументами.
		//Первый параметр - это программа, которую нужно запустить; остальные аргументы - это параметры программы.
		//Cmd представляет внешнюю команду, которая готовится или выполняется.
		//Запускаем bash c опцией -с(команда считывается из строки) и команду которую надо выполнить.
		cmd := exec.Command("bash", "-c", fullCommand)
		//cmd.Stdout и cmd.Stderr методы определяют стандартный вывод процесса и ошибку.
		//Если любой из них равен nil, Run подключает соответствующий файловый дескриптор к нулевому устройству.
		//Все, что вы запишется в /dev/null(нулевое устройство), будет отброшено.
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		//Run запускает указанную команду и ждет ее завершения.
		//Возвращаемая ошибка равна нулю, если команда выполняется, 
		//не имеет проблем с копированием stdin, stdout и stderr и завершается с нулевым статусом выхода.
		return cmd.Run()
	case cd:
		if len(PartOfCommand) < 2 {
			return nil
		}
		//Chdir изменяет текущий рабочий каталог на именованный каталог.
		err := os.Chdir(PartOfCommand[1])
		if err != nil {
			return err
		}
		runCommand("pwd")
		return nil
	case echo:
		cmd := exec.Command("bash", "-c", fullCommand)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	case ps:
		cmd := exec.Command("bash", "-c", fullCommand)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	case kill:
		cmd := exec.Command("bash", "-c", fullCommand)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	case exit:
		os.Exit(0)
	default:
		fmt.Println("Command not found.")
		return nil
	}
	return nil
}
