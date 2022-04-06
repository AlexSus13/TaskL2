/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/
package main
//Паттерн разработки стратегии — это поведенческий паттерн проектирования. 
//Этот шаблон проектирования позволяет изменять поведение объекта во время выполнения без каких-либо изменений в классе этого объекта.
import "fmt"

//Интерфейс отчищения кэша
//который принимает Объект
type evictionAlgo interface {
	evict(c *cache)
}

//Структуры которые имеют метод evict
//==> реализуют интефейс evictionAlgo
type fifo struct {
}

func (l *fifo) evict(c *cache) {
	fmt.Println("Evicting by fifo strtegy")
}

type lru struct {
}

func (l *lru) evict(c *cache) {
	fmt.Println("Evicting by lru strtegy")
}

type lfu struct {
}

func (l *lfu) evict(c *cache) {
	fmt.Println("Evicting by lfu strtegy")
}

//Структура cache - Объект
type cache struct {
	storage      map[string]string
	evictionAlgo evictionAlgo
	capacity     int
	maxCapacity  int
}

//Функция для создания Объекта
func initCache(e evictionAlgo) *cache {
	storage := make(map[string]string)
	return &cache{
		storage:      storage,
		evictionAlgo: e,
		capacity:     0,
		maxCapacity:  2,
	}
}

//Метод Объекта который устанавливает метод отчищения кэша
//Выбирается структура со своим методом evict
func (c *cache) setEvictionAlgo(e evictionAlgo) {
	c.evictionAlgo = e
}

//Метод Объекта который добавляет даные в кэш
//и удаляет выбранным методом если достигнута максимальная память
func (c *cache) add(key, value string) {
	if c.capacity == c.maxCapacity {
		c.evict()
	}
	c.capacity++
	c.storage[key] = value
}

func (c *cache) get(key string) {
	delete(c.storage, key)
}

//Метод evit Объекта
func (c *cache) evict() {
	c.evictionAlgo.evict(c) //У переданной структуры вызываем метод evict
	c.capacity--
}

func main() {
	lfu := &lfu{}
	cache := initCache(lfu) //Меняем мотод отчищения
	cache.add("a", "1")
	cache.add("b", "2")
	cache.add("c", "3")

	lru := &lru{}
	cache.setEvictionAlgo(lru)
	cache.add("d", "4")

	fifo := &fifo{}
	cache.setEvictionAlgo(fifo)
	cache.add("e", "5")
}

