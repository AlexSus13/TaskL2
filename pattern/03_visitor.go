/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/
package main

import "fmt"

//Шаблон проектирования посетителей — это поведенческий шаблон проектирования,
//который позволяет добавлять поведение к структуре, фактически не изменяя структуру.
type Employee interface {
	FullName()
	Accept(Visitor) //для добавления нового интерфейса
}

type Cat struct {
	Name   string
	Breed  string
	Weight int
	Age    int
}

func (c Cat) FullName() {
	fmt.Println("Cat: ", c.Name, " ", c.Breed)
}

//добавим интерфейс к структуре Cat
func (c Cat) Accept(v Visitor) {
	v.VisitCat(c)
}

type Dog struct {
	Name   string
	Breed  string
	Weight int
	Age    int
}

func (d Dog) FullName() {
	fmt.Println("Dog: ", d.Name, " ", d.Breed)
}

//добавим интерфейс к структуре Dog
func (d Dog) Accept(v Visitor) {
	v.VisitDog(d)
}

//Добавленные интерфейсы
type Visitor interface {
	VisitCat(c Cat)
	VisitDog(d Dog)
}

//Новая структура переноска для животных
type AnimalCarrier struct {
	bonusWeight int
}

//Метод новый структуры
func (ac AnimalCarrier) VisitCat(c Cat) {
	fmt.Println(c.Weight + ac.bonusWeight)
}

//Метод новой структуры
func (ac AnimalCarrier) VisitDog(d Dog) {
	fmt.Println(d.Weight + ac.bonusWeight)
}

func main() {

	cat := Cat{"Marty", "British", 4, 1}
	dog := Dog{"Rex", "Mops", 5, 2}

	cat.FullName()
	cat.Accept(AnimalCarrier{2}) //добавляем к веу животного вес переноски

	dog.FullName()
	dog.Accept(AnimalCarrier{3})
}
