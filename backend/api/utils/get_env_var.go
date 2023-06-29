package utils

import (
	"log"
	"os"
)

func GetEnvVar(name string) string {
	variable, present := os.LookupEnv(name)
	if !present {
		log.Fatalf("%s env variable is missing", name)
	}
	return variable
}
