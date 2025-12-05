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
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	divide := strings.Split(data, ",")

	if len(divide) != 3 {
		return 0, "", 0, fmt.Errorf("длина слайса не равен: %d", len(divide))
	}

	steps, err := strconv.Atoi(divide[0])
	if err != nil {
		return 0, "", 0, errors.New("не верный формат количества шагов")
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов равен 0")
	}

	duration, err := time.ParseDuration(divide[2])

	if err != nil {
		return 0, "", 0, errors.New("не верный формат времени")
	}

	if duration <= 0 {
		return 0, "", 0, errors.New("продолжительность времени равен 0")
	}

	return steps, divide[1], duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient

	distance := (float64(steps) * stepLength) / float64(mInKm)

	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distance := distance(steps, height)

	return distance / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, typeTraining, duration, err := parseTraining(data)

	if err != nil {
		return "", err
	}

	switch typeTraining {
	case "Бег":
		distance := distance(steps, height)

		meanSpeed := meanSpeed(steps, height, duration)

		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", errors.New("ошибка при вызове RunningSpentCalories")
		}

		result := fmt.Sprintf("Тип тренировки: %s\n", typeTraining)
		result += fmt.Sprintf("Длительность: %.2f ч.\n", duration.Hours())
		result += fmt.Sprintf("Дистанция: %.2f км.\n", distance)
		result += fmt.Sprintf("Скорость: %.2f км/ч\n", meanSpeed)
		result += fmt.Sprintf("Сожгли калорий: %.2f\n", calories)

		return result, nil
	case "Ходьба":
		distance := distance(steps, height)

		meanSpeed := meanSpeed(steps, height, duration)

		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", errors.New("ошибка при вызове RunningSpentCalories")
		}

		result := fmt.Sprintf("Тип тренировки: %s\n", typeTraining)
		result += fmt.Sprintf("Длительность: %.2f ч.\n", duration.Hours())
		result += fmt.Sprintf("Дистанция: %.2f км.\n", distance)
		result += fmt.Sprintf("Скорость: %.2f км/ч\n", meanSpeed)
		result += fmt.Sprintf("Сожгли калорий: %.2f\n", calories)

		return result, nil
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов равен 0")
	}

	if weight <= 0 || height <= 0 {
		return 0, errors.New("вес и рост должны быть больше 0")
	}

	if duration <= 0 {
		return 0, errors.New("время пути равен 0")
	}

	meanSpeed := meanSpeed(steps, height, duration)

	return ((weight * meanSpeed * duration.Minutes()) / minInH), nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов равен 0")
	}

	if weight <= 0 || height <= 0 {
		return 0, errors.New("вес и рост должны быть больше 0")
	}

	if duration <= 0 {
		return 0, errors.New("время пути равен 0")
	}

	meanSpeed := meanSpeed(steps, height, duration)

	return (((weight * meanSpeed * duration.Minutes()) / minInH) * walkingCaloriesCoefficient), nil
}
