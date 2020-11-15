package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func report() {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	for k, v := range DB {
		switch v.Type {
		case "A", "AAAA", "CNAME":
			line := []string{
				k,
				v.Type,
				v.Values[0],
			}
			w.Write(line)
			if len(v.Values) > 1 {
				for i := 1; i < len(v.Values); i++ {
					line := []string{
						"",
						"",
						v.Values[i],
					}
					w.Write(line)
				}
			}
		default:
			fmt.Printf("Skipping %v (%v)\n", k, v.Type)
		}
	}
}
