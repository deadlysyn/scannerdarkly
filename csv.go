package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func report() {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	w.Write([]string{
		"Name",
		"Zone ID",
		"Type",
		"Values",
	})

	for k, v := range DB {
		switch v.Type {
		case "A", "AAAA", "CNAME":
			w.Write([]string{
				k,
				v.ID,
				v.Type,
				v.Values[0],
			})
			if len(v.Values) > 1 {
				for i := 1; i < len(v.Values); i++ {
					w.Write([]string{
						"",
						"",
						"",
						v.Values[i],
					})
				}
			}
		default:
			fmt.Printf("Skipping %v (%v)\n", k, v.Type)
		}
	}
}
