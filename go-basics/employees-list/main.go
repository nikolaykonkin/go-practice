package main

import "fmt"

// Employee — структура для хранения данных о сотруднике
type Employee struct {
	LastName        string
	FirstName       string
	Age             int
	CurrentPosition string
	Salary          int
}

// Displayable — интерфейс для вывода информации о сотруднике
type Displayable interface {
	Display()
}

// Display выводит данные сотрудника в читаемом виде
func (e Employee) Display() {
	fmt.Printf("%-10s %-10s | Должность: %-15s | Возраст: %d | Зарплата: %d руб.\n",
		e.LastName, e.FirstName, e.CurrentPosition, e.Age, e.Salary)
}

// FilterEmployees возвращает сотрудников, у которых возраст >= minAge и зарплата >= minSalary
func FilterEmployees(employees []Employee, minAge int, minSalary int) []Employee {
	var filtered []Employee
	for _, emp := range employees {
		if emp.Age >= minAge && emp.Salary >= minSalary {
			filtered = append(filtered, emp)
		}
	}
	return filtered
}

// AddEmployee добавляет нового сотрудника в срез и возвращает обновленный срез
func AddEmployee(employees []Employee, lastName, firstName string, age int, position string, salary int) []Employee {
	newEmp := Employee{
		LastName:        lastName,
		FirstName:       firstName,
		Age:             age,
		CurrentPosition: position,
		Salary:          salary,
	}
	return append(employees, newEmp)
}

func main() {
	employees := []Employee{
		{"Андреев", "Андрей", 28, "Разработчик", 160000},
		{"Иванов", "Иван", 34, "Дизайнер", 130000},
		{"Аникина", "Анна", 42, "Тимлид", 300000},
		{"Петров", "Петр", 25, "Аналитик", 95000},
		{"Марьина", "Марина", 35, "Стажер", 85000},
	}

	employees = AddEmployee(employees, "Михайлов", "Михаил", 32, "Тестировщик", 105000)

	minAge := 32
	minSalary := 100000

	fmt.Printf("--- Сотрудники (возраст от %d, зарплата от %d руб.) ---\n\n", minAge, minSalary)

	filtered := FilterEmployees(employees, minAge, minSalary)
	if len(filtered) == 0 {
		fmt.Println("Сотрудники, подходящие под критерии, не найдены.")
		return
	}

	for _, emp := range filtered {
		var d Displayable = emp
		d.Display()
	}
}
