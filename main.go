package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type AudioDevice struct {
	Name string
	ID   string
}

func main() {
	// Get current default sink
	currentSink, err := getCurrentSink()
	if err != nil {
		fmt.Printf("Error getting current sink: %v\n", err)
		return
	}

	// Get available sinks
	sinks, err := getAvailableSinks()
	if err != nil {
		fmt.Printf("Error getting available sinks: %v\n", err)
		return
	}

	// Find devices by pattern matching (more resilient to ID changes)
	var headphones, speakers *AudioDevice
	for _, sink := range sinks {
		sinkNameLower := strings.ToLower(sink.Name)
		sinkIDLower := strings.ToLower(sink.ID)

		// Match Logitech headphones by multiple criteria
		if (strings.Contains(sinkIDLower, "logitech") && strings.Contains(sinkIDLower, "pro_x")) ||
			(strings.Contains(sinkNameLower, "logitech") && strings.Contains(sinkNameLower, "pro x")) ||
			strings.Contains(sinkIDLower, "usb-logitech_pro_x") {
			headphones = &sink
			// Match desktop speakers - now also check for renamed description
		} else if strings.Contains(sinkIDLower, "pci-0000_18_00.6") ||
			strings.Contains(sinkNameLower, "desktop speakers") ||
			(strings.Contains(sinkIDLower, "pci-") && strings.Contains(sinkIDLower, "analog-stereo") &&
				!strings.Contains(sinkIDLower, "hdmi") && !strings.Contains(sinkIDLower, "usb")) {
			speakers = &sink
		}
	}

	// Check if both devices are available
	if headphones == nil {
		fmt.Println("Logitech headphones not available")
		return
	}
	if speakers == nil {
		fmt.Println("Desktop speakers not available")
		return
	}

	// Determine which device to switch to
	var targetDevice *AudioDevice
	var targetName string

	if currentSink == headphones.ID {
		// Currently using headphones, switch to speakers
		targetDevice = speakers
		targetName = "Desktop Speakers"
	} else {
		// Currently using speakers (or something else), switch to headphones
		targetDevice = headphones
		targetName = "Logitech Headphones"
	}

	// Switch to target device
	err = setSink(targetDevice.ID)
	if err != nil {
		fmt.Printf("Error switching to %s: %v\n", targetName, err)
		return
	}

	fmt.Printf("Switched audio to: %s (%s)\n", targetDevice.Name, targetName)
}

func getCurrentSink() (string, error) {
	cmd := exec.Command("pactl", "get-default-sink")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getAvailableSinks() ([]AudioDevice, error) {
	cmd := exec.Command("pactl", "list", "short", "sinks")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var sinks []AudioDevice
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Parse pactl output: ID\tNAME\tDRIVER\tSAMPLE_SPEC\tSTATE
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			sinks = append(sinks, AudioDevice{
				ID:   fields[1], // Name field is actually the ID we need
				Name: fields[1],
			})
		}
	}

	// Get human-readable names
	for i := range sinks {
		name, err := getSinkDescription(sinks[i].ID)
		if err == nil {
			sinks[i].Name = name
		}
	}

	return sinks, nil
}

func getSinkDescription(sinkID string) (string, error) {
	cmd := exec.Command("pactl", "list", "sinks")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse the detailed sink list to find our sink and get its description
	lines := strings.Split(string(output), "\n")
	var inTargetSink bool

	for _, line := range lines {
		if strings.Contains(line, "Name: "+sinkID) {
			inTargetSink = true
			continue
		}

		if inTargetSink && strings.HasPrefix(line, "Sink #") {
			// We've moved to the next sink
			break
		}

		if inTargetSink && strings.Contains(line, "Description:") {
			// Extract description
			re := regexp.MustCompile(`Description:\s*(.+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				return strings.TrimSpace(matches[1]), nil
			}
		}
	}

	return sinkID, nil // Fallback to ID if description not found
}

func setSink(sinkID string) error {
	cmd := exec.Command("pactl", "set-default-sink", sinkID)
	return cmd.Run()
}
