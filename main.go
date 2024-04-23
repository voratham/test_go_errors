package main

import (
	"database/sql"
	"errors"
	"fmt"
)

func firstFunc() error {
	return fmt.Errorf("original error: firstFunc")
}

func secondFuncNewOldVersion() error {
	firstErr := firstFunc()
	if firstErr != nil {
		return fmt.Errorf("failed secondFuncNewOldVersion %+w", firstErr)
	}
	return nil
	// NOTE out
	/*
		failed secondFuncNewOldVersion original error: firstFunc
	*/
}

func secondFuncNewGoVersion() error {
	firstErr := firstFunc()
	if firstErr != nil {
		//  NOTE: in case go.120
		secondErr := errors.New("failed in secondFuncNewGoVersion")
		return errors.Join(firstErr, secondErr)
		// NOTE: output
		/**
		original error: firstFunc
		failed in second function
		*/
	}
	return nil
}

// NOTE: let to try implement custom error
type CustomError struct {
	Message string
	Wrapped error
}

func (e CustomError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Wrapped)
}

func (e CustomError) Unwrap() error {
	return e.Wrapped
}

func SomeFunction() error {
	return CustomError{
		Message: "original error: something went wrong",
		Wrapped: errors.New("wrapped error"),
	}
}

// NOTE: defined enum
var (
	ErrUserNotFound = errors.New("could not find user")
)

func GetFakeUser() (bool, error) {
	// something error
	testErr := sql.ErrNoRows
	if errors.Is(testErr, sql.ErrNoRows) {
		return false, errors.Join(testErr, ErrUserNotFound)
		// return false, fmt.Errorf("%w %w", testErr, ErrUserNotFound)
	}

	return true, nil

}

func main() {

	err := secondFuncNewOldVersion()
	fmt.Println(err)
	fmt.Println("---------------")

	err2 := secondFuncNewGoVersion()
	fmt.Println(err2)

	fmt.Println("---------------")

	innerErr := errors.Unwrap(err) // it will get error first
	fmt.Println("inner error :", innerErr)

	fmt.Println("---------------")

	innerErr2 := errors.Unwrap(err2) // it will nil because error joined
	fmt.Println("inner error 2 :", innerErr2)

	fmt.Println("custom error ---------------")
	customErr := SomeFunction()
	fmt.Printf("customErr common :: %+v\n ", customErr)

	customErrUnwarp := errors.Unwrap(customErr)

	customerErrorInner := fmt.Errorf("inner customer error: %w", customErrUnwarp)

	fmt.Println("customerErrorInner final:", customerErrorInner)

	_, userErr := GetFakeUser()
	if err != nil {
		if errors.Is(userErr, ErrUserNotFound) {
			fmt.Println("GetFakeUser", userErr.Error())
		}
	}

}
