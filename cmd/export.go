package cmd

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/models"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export data to json",
	Run:   exportData,
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP("out", "o", "./export", "Output directory")
}

func exportData(cmd *cobra.Command, args []string) {
	err := db.Open()
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		fmt.Println()
		fmt.Println("If you are running the farma server, you will have to")
		fmt.Println("shut it down in order to use import/export commands.")
		return
	}
	defer db.Close()

	outDir, _ := cmd.Flags().GetString("out")
	if err := os.Mkdir(outDir, 0755); err != nil {
		fmt.Printf("failed to create dir %s: %v\n", outDir, err)
		return
	}

	exportEntities("l_user", []byte("l:user:"), outDir, &models.UserLog{})
	exportEntities("s_id", []byte("s:id:"), outDir, &models.Subscription{})
	exportEntities("f_id", []byte("f:id:"), outDir, &models.Frame{})
	exportEntities("n_id", []byte("n:id:"), outDir, &models.Notification{})
	exportKeys("f_name", []byte("f:name:"), outDir)
	exportKeys("f_endpoint", []byte("f:endpoint:"), outDir)
	exportKeys("s_url", []byte("s:url:"), outDir)
	exportKeys("s_token", []byte("s:token:"), outDir)
	exportKeys("f_pk", []byte("f:pk:"), outDir)
	exportSeq("FrameId", []byte("FrameId"), outDir)
}

func exportEntities(entityName string, prefix []byte, outDir string, entity proto.Message) {
	fmt.Printf("Exporting Entify %s...\n", entityName)

	entities := make(map[string]proto.Message)
	next := prefix
	for {
		keys, nextKey, err := db.GetKeysWithPrefix(prefix, next, 100)
		if err != nil {
			fmt.Printf("Error getting keys: %v\n", err)
			return
		}

		for _, key := range keys {
			value, err := db.Get(key)
			if err != nil {
				fmt.Printf("Error getting key %s: %v\n", key, err)
				return
			}
			err = proto.Unmarshal(value, entity)
			if err != nil {
				fmt.Printf("Error unmarshaling value for key %s: %v. Skipping...\n", key, err)

			} else {
				entities[string(key)] = proto.Clone(entity)
			}

		}

		if nextKey == nil {
			break
		}
		next = nextKey
	}

	exportPbJSON(entities, fmt.Sprintf("%s/%s.json", outDir, entityName))
}

func exportPbJSON(entities map[string]proto.Message, filename string) {
	jsonMap := make(map[string]json.RawMessage, len(entities))
	for k, v := range entities {
		j, err := protojson.Marshal(v)
		if err != nil {
			fmt.Printf("Error marshaling value for key %s: %v. Skipping...\n", k, err)
			continue
		}
		jsonMap[k] = j
	}
	jsonData, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
		return
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing to output file: %v\n", err)
		return
	}
}

func exportKeys(entityName string, prefix []byte, outDir string) {
	fmt.Printf("Exporting Index %s...\n", entityName)

	entities := make(map[string]string)
	next := prefix
	for {
		keys, nextKey, err := db.GetKeysWithPrefix(prefix, next, 100)
		if err != nil {
			fmt.Printf("Error getting keys: %v\n", err)
			return
		}

		for _, key := range keys {
			value, err := db.Get(key)
			if err != nil {
				fmt.Printf("Error getting key %s: %v\n", key, err)
				return
			}

			entities[string(key)] = string(value)
		}

		if nextKey == nil {
			break
		}
		next = nextKey
	}

	exportKeysJSON(entities, fmt.Sprintf("%s/%s.json", outDir, entityName))
}

func exportKeysJSON(entities map[string]string, filename string) {
	jsonData, err := json.Marshal(entities)
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
		return
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing to output file: %v\n", err)
		return
	}
}

func exportSeq(seqName string, key []byte, outDir string) {
	fmt.Printf("Exporting Sequence %s...\n", seqName)

	kv := make(map[string]uint64)
	val, err := db.Get(key)
	if err != nil {
		fmt.Printf("Error getting key %s: %v\n", key, err)
		return
	}

	kv[string(key)] = binary.BigEndian.Uint64(val)

	jsonData, err := json.Marshal(kv)
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
		return
	}

	filename := fmt.Sprintf("%s/%s.json", outDir, seqName)
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing to output file: %v\n", err)
		return
	}
}
