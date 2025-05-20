package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const SEPARATOR = "|"

type Person struct {
	Name string
	Job string
	DOB int
}

func NewPerson(name string, job string, dob int) *Person {
	return &Person{
		Name: name,
		Job: job,
		DOB: dob,
	}
}

func main() {
	file := "./sampleFile.txt"
	persons := readPersonsFromFile(file)

	persons = insertPerson(persons, "Tri Nguyen", "SWE", 2001) // insert a new person
	fmt.Println("Before update:")
	printPersons(persons) // expectation: must have the new added data 


	updatePerson(persons, 2, "Job", "Actress") // update person data at index
	fmt.Println("After update:")
	printPersons(persons) // expectation: must have the new added data 

	writePersonsToFile(persons, file) // write back data to file
}

func readPersonsFromFile(filePath string) []Person {
	f, err := os.Open(filePath)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	persons := make([]Person, 0)

	for scanner.Scan() {
		person := parseTextToPerson(scanner.Text())
		persons = append(persons, *person)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return persons
}

func parseTextToPerson(text string) *Person {
	parsedStringSlice := strings.Split(text, SEPARATOR)

	if len(parsedStringSlice) != reflect.TypeOf(Person{}).NumField() {
		panic("Invalid number of fields in text line. Expected 3 fields separated by '|'")
	}

	name := parsedStringSlice[0]
	job := parsedStringSlice[1]
	dob, err := strconv.Atoi(parsedStringSlice[2])
	
	if err != nil {
		panic(err)
	}

	return NewPerson(name, job, dob)
}

func insertPerson(persons []Person, name string, job string, dob int) []Person {
	return append(persons, *NewPerson(name, job, dob))
}

func printPersons(persons []Person) {
	for i, person := range persons {
		fmt.Printf("Person %d: name: %s, occupation: %s, date of birth: %d\n",
			i + 1, // 1-indexed
			person.Name,
			person.Job,
			person.DOB)
	}
}

func updatePerson(persons []Person, index int, field string, value interface{}) (*Person, error) {
	if index >= len(persons) {
		panic("Attempt to update person failed: Index out of range")
	}

	personAtIndex := &persons[index]
	personValue := reflect.ValueOf(personAtIndex).Elem()
	fieldValue := personValue.FieldByName(field)

	if !fieldValue.IsValid() {
		return nil, fmt.Errorf("Field %s does not exist", field)
	}

	valueToSet := reflect.ValueOf(value)
	if !valueToSet.Type().ConvertibleTo(fieldValue.Type()) {
		return nil, fmt.Errorf("Value type %v cannot be converted to field type %v", valueToSet.Type(), fieldValue.Type())
	}

	fieldValue.Set(valueToSet.Convert(fieldValue.Type()))
	return personAtIndex, nil
}


func writePersonsToFile(persons []Person, targetedFile string) {
	f, err := os.Create(targetedFile)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	for _, person := range persons {
		f.WriteString(parsePersonToText(person))
		f.WriteString("\n")
	}
}

func parsePersonToText(person Person) string {
	return fmt.Sprintf("%s%s%s%s%d", 
		person.Name, 
		SEPARATOR, 
		person.Job, 
		SEPARATOR, 
		person.DOB)
}