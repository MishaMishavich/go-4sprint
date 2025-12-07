package daysteps

import (
	"errors"
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
	divide := strings.Split(data, ",")

	if len(divide) != 2 {
		return 0, 0, fmt.Errorf("длина слайса не равен 2. Получено: %v", len(divide))
	}

	steps, err := strconv.Atoi(divide[0])
	if err != nil {
		return 0, 0, errors.New("не верный формат количества шагов")
	}

	if steps <= 0 {
		return 0, 0, errors.New("количество шагов равен 0")
	}

	duration, err := time.ParseDuration(divide[1])

	if err != nil {
		return 0, 0, errors.New("не верный формат времени")
	}

	if duration <= 0 {
		return 0, 0, errors.New("продолжительность времени равен 0")
	}

	return steps, duration, nil

}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("ошибка при вызове parsePackage", err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distance := (float64(steps) * stepLength) / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return fmt.Sprintf("ошибка при вызове WalkingSpentCalories: %d", err)
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)
}
