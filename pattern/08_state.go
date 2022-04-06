/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/
package main

import "fmt"

//Паттерн State позволяет объекту изменять свое поведение в зависимости от внутреннего состояния.

// MobileAlertStater  предоставляет общий интерфейс для различных состояний.
type MobileAlertStater interface {
	Alert() string
}

// Mobile реализует оповещение в зависимости от его состояния.
type Mobile struct {
	state MobileAlertStater
}

// Alert возвращает строку предупреждения
func (m *Mobile) Alert() string {
	return m.state.Alert()
}

// SetState изменяет состояние
func (m *Mobile) SetState(state MobileAlertStater) {
	m.state = state
}

// NewMobileAlert - это конструктор Mobile Alert.
func NewMobile() *Mobile {
	return &Mobile{state: &MobileAlertSong{}}
}

// MobileAlertVibration реализует вибросигнализацию
type MobileAlertVibration struct {
}

// Alert возвращает строку предупреждения
func (a *MobileAlertVibration) Alert() string {
	return "Vrrr... Brrr... Vrrr..."
}

// MobileAlertSong реализует звуковое оповещение
type MobileAlertSong struct {
}

// Alert возвращает строку предупреждения
func (a *MobileAlertSong) Alert() string {
	return "Ringtone is playing"
}

func main() {
	Mobile := NewMobile()
	alter := Mobile.Alert()
	fmt.Println(alter)

	MobileAlertVibration := &MobileAlertVibration{}

	Mobile.SetState(MobileAlertVibration)
	alter2 := Mobile.Alert()
	fmt.Println(alter2)
}
