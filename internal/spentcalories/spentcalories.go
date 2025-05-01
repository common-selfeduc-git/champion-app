package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	dataParts := strings.Split(data, ",")

	if len(dataParts) != 3 {
		err := fmt.Errorf("ожидалась строка из трёх частей, получилось частей: %d", len(dataParts))
		log.Print(err.Error())
		return 0, "", 0, err
	}

	steps, err := strconv.Atoi(dataParts[0])
	// уважаемый ревьювер! тут мне не понятно по условию постановки задачи, должен я вернуть вид активности или пустую строку, потому возвращаю вид активности
	if err != nil {
		return 0, dataParts[1], 0, err
	}

	if steps <= 0 {
		err := fmt.Errorf("количество шагов !(%d > 0)", steps)
		log.Print(err.Error())
		return 0, dataParts[1], 0, err
	}

	duration, err := time.ParseDuration(dataParts[2])

	if err != nil {
		return 0, dataParts[1], 0, err
	}

	if duration <= 0 {
		err := fmt.Errorf("длительность !(%d > 0)", duration)
		log.Print(err.Error())
		return 0, dataParts[1], 0, err
	}

	if len(dataParts) <= 0 {
		err := fmt.Errorf("вид тренировки не указан %s", dataParts[1])
		log.Print(err.Error())
		return 0, dataParts[1], 0, err
	}

	return steps, dataParts[1], duration, nil
}

func distance(steps int, height float64) float64 {

	stepLength := height * stepLengthCoefficient

	distanceM := stepLength * float64(steps)

	distanceKm := distanceM / mInKm

	return distanceKm

}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0
	}

	distance := distance(steps, height)

	durationH := distance / duration.Hours()

	return durationH

}

func TrainingInfo(data string, weight, height float64) (string, error) {

	steps, exType, duration, err := parseTraining(data)

	if err != nil {
		return "", err
	}

	distance := distance(steps, height)

	meanSpeed := meanSpeed(steps, height, duration)

	var caloriesSpent float64

	switch exType {
	case "Бег":
		caloriesSpent, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Ходьба":
		caloriesSpent, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки %s", exType)
	}

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", exType, duration.Hours(), distance, meanSpeed, caloriesSpent), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		err := fmt.Errorf("количество шагов %d<=0", steps)
		log.Print(err.Error())
		return 0, err
	}

	if weight <= 0 {
		err := fmt.Errorf("вес %.2f<=0", weight)
		log.Print(err.Error())
		return 0, err
	}

	if height <= 0 {
		err := fmt.Errorf("рост %.2f<=0", height)
		log.Print(err.Error())
		return 0, err
	}

	if duration <= 0 {
		err := fmt.Errorf("некорректная длительность %s", duration)
		log.Print(err.Error())
		return 0, err
	}

	meanSpeed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	calories := (weight * meanSpeed * durationInMinutes) / minInH

	return calories, nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов %d<=0", steps)
	}

	if weight <= 0 {
		return 0, fmt.Errorf("вес %.2f<=0", weight)
	}

	if height <= 0 {
		return 0, fmt.Errorf("рост %.2f<=0", height)
	}

	meanSpeed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	calories := (weight * meanSpeed * durationInMinutes) / minInH

	caloriesSpent := calories * walkingCaloriesCoefficient

	return caloriesSpent, nil

}
