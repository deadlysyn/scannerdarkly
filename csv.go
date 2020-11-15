package main

import (
	"encoding/csv"
	"os"
)

func reportCSV() {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	w.Write([]string{
		"Name",
		"Zone ID",
		"Type",
		"Values",
	})

	for id, recs := range DB {
		for _, rec := range recs {
			w.Write([]string{
				rec.Name,
				id,
				rec.Type,
				rec.Values[0],
			})
			if len(rec.Values) > 1 {
				for i := 1; i < len(rec.Values); i++ {
					w.Write([]string{
						"",
						"",
						"",
						rec.Values[i],
					})
				}
			}
		}
	}
}
