/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/
package main
//Цепочка вызывов — это поведенческий паттерн проектирования, который позволяет передавать запросы последовательно по цепочке обработчиков.
import "fmt"

type User struct {
	NameVerificationDone     bool
	PasswordVerificationDone bool
	AuthDone                 bool
}

//Интерфейс регистрации
type AuthService interface {
	Execute(*User)
	SetNext(AuthService)
}

//Структура проверки имени пользователя
type NameVerification struct {
	Next AuthService
}

//Метод проверяющий была ли завершена проверка имени
func (nv *NameVerification) Execute(u *User) {
	if u.NameVerificationDone {
		fmt.Println("Проверка Имени прошла успешно")
		nv.Next.Execute(u) //Запуск следующео метода
	}

	fmt.Println("Проверка имени...")
	u.AuthDone = true
	nv.Next.Execute(u)
}

//Установка следующего метода для проверки
func (nv *NameVerification) SetNext(next AuthService) {
	nv.Next = next
}

//Структура проверки пароля пользователя
type PasswordVerification struct {
	Next AuthService
}

//Метод проверяющий была ли завершена проверка пароля
func (pv *PasswordVerification) Execute(u *User) {
	if u.PasswordVerificationDone {
		fmt.Println("Проверка Пароля прошла успешно")
		pv.Next.Execute(u) //Запуск следующео метода
	}

	fmt.Println("Проверка пароля...")
	u.PasswordVerificationDone = true
	pv.Next.Execute(u)
}

//Установка следующего метода для проверки
func (pv *PasswordVerification) SetNext(next AuthService) {
	pv.Next = next
}

//Структура проверки регистрации пользователя
type Registration struct {
	Next AuthService
}

//Метод проверяющий была ли завершена регистрация
func (r *Registration) Execute(u *User) {
	if u.AuthDone {
		fmt.Println("Проверка валидности данных...")
	}

	fmt.Println("Аутентификация прошла успешно")
}

//Установка следующего метода для проверки
func (r *Registration) SetNext(next AuthService) {
	r.Next = next
}

func main() {
	//Создаем пользователя который проходит аутентификацию
	user := new(User)
	//Создаем структуры которые будут выполнять шаги аутентификации
	NameV := new(NameVerification)
	PasswordV := new(PasswordVerification)
	Reg := new(Registration)

	//У структур назначаем следующую структуру которая будет выполнять
	//следующий шаг аутентификации
	NameV.SetNext(PasswordV)
	PasswordV.SetNext(Reg)

	//Отправляем пользовател на проверку
	NameV.Execute(user)
}

