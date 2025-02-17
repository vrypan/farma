package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/utils"
)

var frameLsCmd = &cobra.Command{
	Use:     "list [frame name",
	Short:   "List frames",
	Aliases: []string{"ls"},
	Run:     lsFrame,
}

func init() {
	frameCmd.AddCommand(frameLsCmd)
	frameLsCmd.Flags().BoolP("with-subscriptions", "", false, "Show subscriptions for each frame")
}

func lsFrame(cmd *cobra.Command, args []string) {
	subscriptionsFlag, _ := cmd.Flags().GetBool("with-subscriptions")

	db.Open()
	defer db.Close()

	frames := utils.AllFrames()

	/*for _, frame := range frames {
		fmt.Printf("%04d %-32s %s %s\n", frame.Id, frame.Name, frame.Endpoint, frame.Domain)
		if subscriptionsFlag {
			subscriptions, _, err := utils.SubscriptionsByFrame(frame.Id, 100)
			if err != nil {
				fmt.Printf("Error fetching subscriptions: %v\n", err)
				continue
			}
			for _, s := range subscriptions {
				fmt.Printf("  > %s\n", s.NiceString())
			}
		}
	}
	{
	*/
	// Create a struct to hold the output data
	type Output struct {
		FrameId       int      `json:"frame_id"`
		FrameName     string   `json:"frame_name"`
		FrameEndpoint string   `json:"frame_endpoint"`
		FrameDomain   string   `json:"frame_domain"`
		Subscriptions []string `json:"subscriptions,omitempty"`
	}

	// create a slice to hold all output data
	outputs := []Output{}

	for _, frame := range frames {
		frameJson, err := json.Marshal(frame)
		if err != nil {
			fmt.Printf("Error converting frame to json: %v\n", err)
			continue
		}
		fmt.Println(string(frameJson))
		output := Output{
			FrameId:       int(frame.Id),
			FrameName:     frame.Name,
			FrameEndpoint: frame.Endpoint,
			FrameDomain:   frame.Domain,
		}

		// if subscriptions flag is set, fetch and add subscriptions to the output
		if subscriptionsFlag {
			subscriptions, _, err := utils.SubscriptionsByFrame(frame.Id, 100)
			if err != nil {
				fmt.Printf("Error fetching subscriptions: %v\n", err)
				continue
			}
			for _, s := range subscriptions {
				output.Subscriptions = append(output.Subscriptions, s.NiceString())
			}
		}

		outputs = append(outputs, output)
	}

	// convert the outputs slice to json
	outputJson, err := json.Marshal(outputs)
	if err != nil {
		fmt.Printf("Error converting output to json: %v\n", err)
		return
	}

	// print the json data
	fmt.Println(string(outputJson))
}

/*
	fmt.Println("--- DEBUG ---")
	kk, err := db.GetPrefix([]byte(""), 1000)
	fmt.Println("GetPrefix() returned", err)
	fmt.Println()
	for _, k := range kk {
		fmt.Println(string(k))
	}
*/
//}
