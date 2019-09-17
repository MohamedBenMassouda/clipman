package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func store(text string, history []string, histfile string, max int, persist bool) error {
	if text == "" {
		return nil
	}

	l := len(history)
	if l > 0 {
		// drop oldest items that exceed max list size
		if l >= max {
			// usually just one item, but more if we suddenly reduce our --max-items
			history = history[l-max+1:]
		}

		// remove duplicates
		history = filter(history, text)
	}

	history = append(history, text)

	// dump history to file so that other apps can query it
	if err := write(history, histfile); err != nil {
		return fmt.Errorf("error writing history: %s", err)
	}

	// make the copy buffer available to all applications,
	// even when the source has disappeared
	if persist {
		if err := exec.Command("wl-copy", []string{"--", text}...).Run(); err != nil {
			log.Printf("Error running wl-copy: %s", err) // don't abort, minor error
		}
	}

	return nil
}

// filter removes all occurrences of text
func filter(slice []string, text string) []string {
	var filtered []string
	for _, s := range slice {
		if s != text {
			filtered = append(filtered, s)
		}
	}

	return filtered
}

// write dumps history to json file
func write(history []string, histfile string) error {
	b, err := json.Marshal(history)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(histfile, b, 0644)
}
