/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/
package main

import "fmt"

//Шаблон Builder состоит из четырех компонентов: интерфейса Builder, конкретного Builder, директора и продукта.
//интерфейс Мойщик машин (Builder)
type СarWasher interface {
	setBodyType()
	setNumWheels()
	setInteriorType()
	getCleanCar() car
}

//В зависимости от типа машины разные автомойщики
func getСarWasher(СarWasherType string) СarWasher {
	if СarWasherType == "truck" {
		return &truckСarWasher{}
	}
	if СarWasherType == "passenger" {
		return &passengerСarWasher{}
	}
	return nil
}

//Структура мойщика грузовых машин
type truckСarWasher struct {
	BodyType     string
	InteriorType string
	Wheels       int
}

func newtruckСarWasher() *truckСarWasher {
	return &truckСarWasher{}
}

//Методы помывки грузовика
func (t *truckСarWasher) setBodyType() {
	t.BodyType = "Большой корпус"
}

func (t *truckСarWasher) setInteriorType() {
	t.InteriorType = "Грубая кожа"
}

func (t *truckСarWasher) setNumWheels() {
	t.Wheels = 10
}

//Получаем чистый грузовик
func (t *truckСarWasher) getCleanCar() car {
	return car{
		BodyType:     t.BodyType,
		InteriorType: t.InteriorType,
		Wheels:       t.Wheels,
	}
}

//Структура мойщика легковых
type passengerСarWasher struct {
	BodyType     string
	InteriorType string
	Wheels       int
}

//создание нового мойщика легковых автомобилей
func newPassengerСarWasher() *passengerСarWasher {
	return &passengerСarWasher{}
}

//Методы помывки лековой машины
func newpassengerСarWasher() *passengerСarWasher {
	return &passengerСarWasher{}
}

func (p *passengerСarWasher) setBodyType() {
	p.BodyType = "Маленький корпус"
}

func (p *passengerСarWasher) setInteriorType() {
	p.InteriorType = "Мягкая кожа"
}

func (p *passengerСarWasher) setNumWheels() {
	p.Wheels = 4
}

//Получаем чистый легковой автомобиль
func (p *passengerСarWasher) getCleanCar() car {
	return car{
		BodyType:     p.BodyType,
		InteriorType: p.InteriorType,
		Wheels:       p.Wheels,
	}
}

//Структура машины
type car struct {
	BodyType     string
	InteriorType string
	Wheels       int
}

//Стуктура директора
type director struct {
	washer СarWasher
}

//Создать директора
func newDirector(c СarWasher) *director {
	return &director{
		washer: c,
	}
}

//Изменить работника директору
func (d *director) setСarWasher(cw СarWasher) {
	d.washer = cw
}

//Выполнение работником задания
func (d *director) CleanCar() car {
	d.washer.setBodyType()
	d.washer.setInteriorType()
	d.washer.setNumWheels()
	return d.washer.getCleanCar()
}

func main() {
	truckСarWasher := getСarWasher("truck")
	passengerСarWasher := getСarWasher("passenger")

	director := newDirector(truckСarWasher)
	truckСar := director.CleanCar()

	fmt.Printf("Truck Сar Body Type: %s\n", truckСar.BodyType)
	fmt.Printf("Truck Сar Interior Type: %s\n", truckСar.InteriorType)
	fmt.Printf("Truck Сar Num Wheels: %d\n", truckСar.Wheels)

	director.setСarWasher(passengerСarWasher)
	passengerСar := director.CleanCar()

	fmt.Printf("Passenger Сar Body Type: %s\n", passengerСar.BodyType)
	fmt.Printf("Passenger Сar Interior Type: %s\n", passengerСar.InteriorType)
	fmt.Printf("Passenger Сar Num Wheels: %d\n", passengerСar.Wheels)
}
