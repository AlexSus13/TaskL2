/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/
package main

import "fmt"

//Интерфейс который реализует Объект
type ICar interface {
	GetName() string
	GetSpeed() int
}

//Объект
type Car struct {
	Name  string
	Speed int
}

//Методы Объекта
func (c *Car) GetName() string {
	return c.Name
}

func (c *Car) GetSpeed() int {
	return c.Speed
}

//Создаем новую структуру которая будет наследовать методы
//Объекта ==> структура будет реализовывать интерфейс ICar
type Audi struct {
	Car
}

//Поскольку структура Audi реализует интерфейс ICar
//мы можем ее вернуть как объект типа ICar
func NewAudi() ICar {
	return &Audi{
		Car: Car{
			Name:  "Audi",
			Speed: 210,
		},
	}
}

type BMW struct {
	Car
}

func NewBMW() ICar {
	return &BMW{
		Car: Car{
			Name:  "BMW",
			Speed: 230,
		},
	}
}

//Мы можем создавать струкуры в зависимости от входных данных
func GetCar(s string) (ICar, error) {
	switch s {
	case "Audi":
		return NewAudi(), nil
	case "BMW":
		return NewBMW(), nil
	default:
		return nil, fmt.Errorf("No Car")
	}
}

func Details(Car ICar) string {
	return fmt.Sprintf("Car: %s\nSpeed: %d\n", Car.GetName(), Car.GetSpeed())
}

func main() {
	Cars := []string{"Audi", "BMW", "Alfa Romeo"}

	for _, c := range Cars {
		Car, err := GetCar(c)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(Details(Car))
	}
}
