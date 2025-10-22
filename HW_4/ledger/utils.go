package ledger

import "log"

func CheckValid(v Validatable) error {
	err := v.Validate()
	if err != nil {
		log.Printf("Validation failed: %v", err)
		return err
	}
	log.Println("Validation successful")
	return nil
}
