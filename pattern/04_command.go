/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/
package main

//Шаблон команды это поведенческий шаблон проектирования.
//Он предлагает инкапсулировать запрос как отдельный объект.
//Созданный объект имеет всю информацию о запросе и,
//таким образом, может выполнять его самостоятельно.
import (
	"fmt"
)

//Объект реализует интерфейс device
type device interface {
	on()
	off()
}

//Объект
type tv struct {
	isRunning bool
}

//Методы Объекта
func (t *tv) on() {
	t.isRunning = true
	fmt.Println("Turning tv on")
}

func (t *tv) off() {
	t.isRunning = false
	fmt.Println("Turning tv off")
}

//Интерфейс выполнения
type command interface {
	execute()
}

//Создаем структуры(команды) которые содержат объект,
//реализующий интерфейс device(содержащий методы Объекта)
//и реализуют интерфейс command с методом execute()
type offCommand struct {
	device device
}

//Методы "активируют" методы Объекта
func (c *offCommand) execute() {
	c.device.off()
}

type onCommand struct {
	device device
}

func (c *onCommand) execute() {
	c.device.on()
}

//Создаем структуру которая принимает объект который
//реализует интерфейс выполнения(execute)
type button struct {
	command command
}

//У структуры есть метод который "активируюет" метод
//выполнения(execute)
func (b *button) press() {
	b.command.execute()
}

func main() {
	tv := &tv{}

	onCommand := &onCommand{
		device: tv,
	}
	offCommand := &offCommand{
		device: tv,
	}

	onButton := &button{
		command: onCommand,
	}
	onButton.press()

	offButton := &button{
		command: offCommand,
	}

	offButton.press()
}
