package main

import (
	"encoding/csv"
	"os"
)

func reportCSV() {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	// header
	w.Write([]string{
		"Zone ID",
		"Name",
		"Type",
		"Results",
	})

	for id, recs := range DB {
		for _, rec := range recs {
			if len(rec.Active) == 0 {
				rec.Active = append(rec.Active, "No open ports found")
			}
			t := rec.Type
			if rec.Alias {
				t = "Alias"
			}
			w.Write([]string{
				id,
				rec.Name,
				t,
				rec.Active[0],
			})
			if len(rec.Values) > 1 {
				for i := 1; i < len(rec.Values); i++ {
					w.Write([]string{
						"",
						"",
						"",
						rec.Active[i],
					})
				}
			}
		}
	}
}
