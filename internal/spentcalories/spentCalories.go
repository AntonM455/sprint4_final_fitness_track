package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

// parseTraining: принимает строку с даннымикоторая содержит кол-во шагов, вид активности
// и продолжительность активности. Функция возвращает кол-во шагов, вид и продолжительность
// активности, ошибку. При этом все преобразуется в нужные типы.
func parseTraining(data string) (int, string, time.Duration, error) {
	// Делим входящую строку на части, делаем слайсы.
	separatePart := strings.Split(data, ",")
	if len(separatePart) != 3 {
		return 0, "", 0, errors.New("Ошибка: к нам не пришло одно из значений")
	}
	// Преобразуем строку в число, возвращаем int кол-ва шагов.
	countSteps, err := strconv.Atoi(separatePart[0])
	if err != nil {
		return 0, "", 0, errors.New("ошибка: неверный формат данных для шагов")
	}
	// Преобразуем строку в time.Duration, возвращаем.
	duration, err := time.ParseDuration(separatePart[2])
	if err != nil {
		return 0, "", 0, errors.New("Ошибка: неверный формат данных для продолжительности прогулки")
	}
	return countSteps, separatePart[1], duration, nil

}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	distanceInKm := float64(steps) * lenStep / float64(mInKm)
	return distanceInKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration <= 0 { // Проверяем, чтобы продолжительность не была равна 0
		return 0
	}
	return distance(steps) / duration.Hours() // Возвращаем среднюю скорость
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	// Парсим входные данные с помощью parseTraining
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		return fmt.Sprintf("Ошибка: %v", err)
	}

	// Рассчитываем дистанцию и среднюю скорость.
	dist := distance(steps)
	speed := meanSpeed(steps, duration)

	var calories float64

	// Определяем тип тренировки и считаем калории для каждого из них.
	switch activityType {
	case "Бег":
		calories = RunningSpentCalories(steps, weight, duration)
	case "Ходьба":
		calories = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "Ошибка: неизвестный тип тренировки"
	}

	// Формируем строку результата
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activityType, duration.Hours(), dist, speed, calories)
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	// вызываем функцию meanSpeed, чтобы посчитать среднюю скорость
	runMeanSpead := meanSpeed(steps, duration)
	// Возвращаем кол-во потраченных калорий при беге
	return ((runningCaloriesMeanSpeedMultiplier * runMeanSpead) - runningCaloriesMeanSpeedShift) * weight
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	// вызываем функцию meanSpeed, чтобы посчитать среднюю скорость.
	walkMeanSpead := meanSpeed(steps, duration)
	// Возвращаем кол-во потраченных калорий при ходьбе.
	return ((walkingCaloriesWeightMultiplier * weight) + (walkMeanSpead*walkMeanSpead/height)*walkingSpeedHeightMultiplier) * duration.Hours() * minInH

}
