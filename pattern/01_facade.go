package main

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/
import (
	"fmt"
)

//Платежная структура фасад
type PaymentFacade struct {
	Bank         *Bank
	Account      *Account
	Card         *Card
	User         *User
	Notification *Notification
}

//Метод фасада по списанию средств
func (pf *PaymentFacade) DebitiOperation(debit int, Password, Name string) error {
	//Проверка пароля от аккаунта
	err := pf.Account.CheckPasswod(Password)
	if err != nil {
		return err
	}
	//Получаем данные карты
	CVV := pf.Account.Card.ReturnCVVandOwnerCard()

	//Проверка данных карты в банке
	NCardInBank, err := pf.Bank.СheckСardByBank(CVV)
	if err != nil {
		return err
	}
	//Списание средств со счета
	err = pf.Bank.DebitBalance(NCardInBank, debit)
	if err != nil {
		return err
	}
	//Отправка информационного сообщения о списании
	pf.Notification.NotificationOfDebit(Name)

	return nil
}

//Метод фасада для пополнения средств
func (pf *PaymentFacade) AddOperation(debit int, Password, Name string) error {

	//Проверка пароля от аккаунта
	err := pf.Account.CheckPasswod(Password)
	if err != nil {
		return err
	}
	//Получаем данные карты
	CVV := pf.Account.Card.ReturnCVVandOwnerCard()

	//Проверка данных карты в банке
	NCardInBank, err := pf.Bank.СheckСardByBank(CVV)
	if err != nil {
		return err
	}
	//Пополнение счета
	err = pf.Bank.AddBalance(NCardInBank, debit)
	if err != nil {
		return err
	}
	//Отправка информационного сообщения о пополнении
	pf.Notification.NotificationOfAdd(Name)

	return nil
}

//Создание структуры фасад
func NewPaymentFacade(Password, Name, CVV string, Balance int) *PaymentFacade {
	var Cards []*Card
	var Users []*User
	User := NewUser(Name, Password)
	Card := NewCard(CVV, Balance)
	Cards = append(Cards, Card)
	Users = append(Users, User)
	Account := NewAccount(Password, Card)
	Bank := NewBank(Cards)
	Notification := NewNotification(Users)

	PaymentFacade := &PaymentFacade{
		Bank:         Bank,
		Account:      Account,
		Card:         Card,
		User:         User,
		Notification: Notification,
	}
	return PaymentFacade
}

//Структра Bank
type Bank struct {
	Cards []*Card
}

//Метод банка для проверки подлинности карты по CVV
func (b *Bank) СheckСardByBank(CVV string) (int, error) {
	error := fmt.Errorf("Неправильные данные карты")

	for NCardInBank, card := range b.Cards {
		if card.CVV == CVV {
			fmt.Println("Проверка карты по CVV в банке")
			return NCardInBank, nil
		}
	}

	return 0, error
}

//Метод банка для списания средств
func (b *Bank) DebitBalance(NCardInBank, debit int) error {

	if (b.Cards[NCardInBank].Balance - debit) > 0 {
		b.Cards[NCardInBank].Balance -= debit
		fmt.Println("Банк выполняет операцию списания")
		return nil
	} else {
		return fmt.Errorf("Недостаточно средств для списания")
	}
}

//Метод банка для добаления средств
func (b *Bank) AddBalance(NCardInBank, debit int) error {

	b.Cards[NCardInBank].Balance += debit
	fmt.Println("Банк выполняет операцию пополнения")
	return nil
}

//Создать новый банк
func NewBank(Cards []*Card) *Bank {
	Bank := &Bank{
		Cards: Cards,
	}
	return Bank
}

//Структура Account
type Account struct {
	Password string
	Card     *Card
}

//Метод структуры для проверки пароля
func (ac *Account) CheckPasswod(Password string) error {
	if ac.Password != Password {
		return fmt.Errorf("Неверный пароль от аккаунта")
	}
	fmt.Println("Вход в Аккаунт выполнен")
	return nil
}

//Создать новый аккаунт
func NewAccount(Password string, Card *Card) *Account {
	Account := &Account{
		Password: Password,
		Card:     Card,
	}
	return Account
}

//Структура Card
type Card struct {
	CVV     string
	Balance int
}

//Метод карты для возврата CVV
func (c *Card) ReturnCVVandOwnerCard() string {
	fmt.Println("Получение CVV карты")
	return c.CVV
}

//Создание новой карты
func NewCard(CVV string, Balance int) *Card {
	Card := &Card{
		CVV:     CVV,
		Balance: Balance,
	}
	return Card
}

//Структура User
type User struct {
	Name     string
	SMS      []string
	Password string
}

//Создать нового пользователя
func NewUser(Name, Password string) *User {
	User := &User{
		Name:     Name,
		Password: Password,
	}
	return User
}

//Структура Notification
type Notification struct {
	Users []*User
}

//Метод для отправки сообщения о списании
func (n *Notification) NotificationOfDebit(Name string) {
	for _, user := range n.Users {
		if user.Name == Name {
			user.SMS = append(user.SMS, "Успешное списание средств,")
		}
	}
}

//Метод для отправки сообщения о пополнении
func (n *Notification) NotificationOfAdd(Name string) {
	for _, user := range n.Users {
		if user.Name == Name {
			user.SMS = append(user.SMS, "Баланс успешно пополнен,")
		}
	}
}
func NewNotification(Users []*User) *Notification {
	Notification := &Notification{
		Users: Users,
	}
	return Notification
}

func main() {
	PaymentFacade := NewPaymentFacade("Password", "Alex", "123", 1000)
	fmt.Println("Сообщения о переводах с карты и баланс: ", PaymentFacade.User.Name, PaymentFacade.User.SMS, PaymentFacade.Card.Balance)

	err := PaymentFacade.DebitiOperation(500, "Password", "Alex")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("После списания: ", PaymentFacade.User.Name, PaymentFacade.User.SMS, PaymentFacade.Card.Balance)

	err = PaymentFacade.AddOperation(250, "Password", "Alex")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("После пополнения: ", PaymentFacade.User.Name, PaymentFacade.User.SMS, PaymentFacade.Card.Balance)
}

