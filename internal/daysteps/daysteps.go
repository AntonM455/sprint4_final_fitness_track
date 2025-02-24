package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах.
)

// parsePackage: Мы передаем в функцию строку, в которой всегда должно быть кол-во шагов
// и продолжительность прогулки. Если одного из параметров нет - возвращаем ошибку. На выход функция
// выдает кол-во шагов и продолжительность прогулки в формате int и time, т.е. парсит.
func parsePackage(data string) (int, time.Duration, error) {
	// Делим входящую строку на части, делаем слайсы.
	separatePart := strings.Split(data, ",")
	if len(separatePart) != 2 {
		return 0, 0, errors.New("Ошибка: к нам не пришло одно из значений")
	}
	// Преобразуем строку в число, возвращаем int кол-ва шагов.
	countSteps, err := strconv.Atoi(separatePart[0])
	if err != nil {
		return 0, 0, errors.New("Ошибка: неверный формат данных для шагов")
	}
	// Преобразуем строку в time.Duration, возвращаем.
	duration, err := time.ParseDuration(separatePart[1])
	if err != nil {
		return 0, 0, errors.New("Ошибка: неверный формат данных для продолжительности прогулки")
	}
	return countSteps, duration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	// вызываем функцию parsePackage для того чтобы преобразовать нашу строку в 2 типа: шаги
	// длительность прогулки.
	steps, duration, err := parsePackage(data)
	if err != nil {
		return "Ошибка: что-то пошло не так"
	}
	// проверяем чтобы кол-во щагов не было равно 0, если ошибка - возвращаем пустую строку.
	if steps <= 0 {
		return ""
	}
	// Считаем длину в метрах, переводим в километры, далее вызываем функцию WalkingSpentCalories
	// чтобы посчитать сколько калорий мы сожгли на ходьбе, выводим результат.
	distanceInMt := float64(steps) * StepLength
	distanceInKm := distanceInMt / 1000
	caloriesBurned := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	result := fmt.Sprintf("Количество шагов: %d. \nДистанция составила %.2f км. \nВы сожгли %.2f ккал.",
		steps, distanceInKm, caloriesBurned)
	return result
}
