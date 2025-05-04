package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	dataParts := strings.Split(data, ",")

	if len(dataParts) != 2 {
		err := fmt.Errorf("ожидалась строка из двух частей, получилось частей: %d", len(dataParts))
		log.Print(err.Error())
		return 0, 0, err
	}

	steps, err := strconv.Atoi(dataParts[0])

	if err != nil {
		return 0, 0, err
	}

	if steps <= 0 {
		err := fmt.Errorf("количество шагов !(%d > 0)", steps)
		log.Print(err.Error())
		return 0, 0, err
	}

	duration, err := time.ParseDuration(dataParts[1])

	if err != nil {
		return 0, 0, err
	}

	if duration <= 0 {
		err := fmt.Errorf("длительность !(%d > 0)", duration)
		log.Print(err.Error())
		return 0, 0, err
	}

	return steps, duration, nil

}

func DayActionInfo(data string, weight, height float64) string {
	steps, walkDuration, err := parsePackage(data)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	if walkDuration <= 0 {
		return ""
	}

	walkLengthM := float64(steps) * stepLength

	walkLengthKm := walkLengthM / mInKm

	caloriesSpent, err := spentcalories.WalkingSpentCalories(steps, weight, height, walkDuration)

	if err != nil {
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, walkLengthKm, caloriesSpent)

}
