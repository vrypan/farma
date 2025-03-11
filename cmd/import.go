package cmd

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/models"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import data from json",
	Run:   importData,
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringP("path", "p", "./export", "Import directory")
}

func importData(cmd *cobra.Command, args []string) {
	err := db.Open()
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		fmt.Println()
		fmt.Println("If you are running the farma server, you will have to")
		fmt.Println("shut it down in order to use import/export commands.")
		return
	}
	defer db.Close()

	outDir, _ := cmd.Flags().GetString("path")

	processPbFile("s_id", outDir)
	processPbFile("f_id", outDir)
	processPbFile("l_user", outDir)
	processPbFile("n_id", outDir)
	//processKvFile("f_name", outDir)
	//processKvFile("f_endpoint", outDir)
	//processKvFile("s_url", outDir)
	//processKvFile("s_token", outDir)
	//processSeqFile("FrameId", outDir)

}

func processSeqFile(dataType string, outDir string) {
	filename := fmt.Sprintf("%s/%s.json", outDir, dataType)
	var jsonMap map[string]uint64

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	entitiesBytes, _ := io.ReadAll(file)
	err = json.Unmarshal(entitiesBytes, &jsonMap)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return
	}

	for k, v := range jsonMap {
		fmt.Printf("Importing Key: %s\n", k)

		buffer := new(bytes.Buffer)
		binary.Write(buffer, binary.BigEndian, v)
		err = db.Set([]byte(k), buffer.Bytes())
		if err != nil {
			fmt.Printf("Error saving key: %v\n", err)
		}
	}

}

func processKvFile(dataType string, outDir string) {
	filename := fmt.Sprintf("%s/%s.json", outDir, dataType)
	var jsonMap map[string]string

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	entitiesBytes, _ := io.ReadAll(file)
	err = json.Unmarshal(entitiesBytes, &jsonMap)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return
	}

	for k, v := range jsonMap {
		fmt.Printf("Importing Key: %s\n", k)
		err = db.Set([]byte(k), []byte(v))
		if err != nil {
			fmt.Printf("Error saving key: %v\n", err)
		}
	}

}

func processPbFile(dataType string, outDir string) {
	filename := fmt.Sprintf("%s/%s.json", outDir, dataType)
	var jsonMap map[string]json.RawMessage

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	entitiesBytes, _ := io.ReadAll(file)
	err = json.Unmarshal(entitiesBytes, &jsonMap)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return
	}

	processEntities(jsonMap, dataType)
}

func processEntities(jsonMap map[string]json.RawMessage, dataType string) {
	type dataFunc func() proto.Message
	typeMap := map[string]dataFunc{
		"s_id":   func() proto.Message { return &models.Subscription{} },
		"f_id":   func() proto.Message { return &models.Frame{} },
		"l_user": func() proto.Message { return &models.UserLog{} },
		"n_id":   func() proto.Message { return &models.Notification{} },
	}

	/*
		exportKeys("f_name", []byte("f:name:"), outDir)
		exportKeys("f_endpoint", []byte("f:endpoint:"), outDir)
		exportKeys("s_url", []byte("s:url:"), outDir)
		exportKeys("s_token", []byte("s:token:"), outDir)
		exportSeq("FrameId", []byte("FrameId"), outDir)
	*/
	for k, v := range jsonMap {
		var bytes []byte
		fmt.Println("Importing key:", k)

		dataFunc, ok := typeMap[dataType]
		if !ok {
			fmt.Printf("Invalid datatype: %s\n", dataType)
			continue
		}

		msg := dataFunc()
		err := protojson.Unmarshal(v, msg)
		if err != nil {
			fmt.Printf("Error unmarshalling protojson: %v\n", err)
			continue
		}

		bytes, err = proto.Marshal(msg)
		if err != nil {
			fmt.Printf("Error marshalling proto: %v\n", err)
			continue
		}
		err = db.Set([]byte(k), bytes)
		if err != nil {
			fmt.Printf("Error saving key: %v\n", err)
		}
	}
}
