package title

import (
    "fmt"
    "os"
	"time"
    

	"msattack/config"
	"msattack/managers"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

const PackInfoURL = "https://%s%s/pack/%d/"

func GetPackInfo(c *fiber.Ctx) error {
    log.Info().Str("body", string(c.Body())).Msg("payload from client")
    host := c.Hostname()
    proto := c.Protocol()
    actualPackInfoURL := fmt.Sprintf("%s://%s/snkp/msatk/prod/pack/6120000/pack_info_list.txt", proto, host)
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "version":     6120000,
        "url":         actualPackInfoURL,
        "status":      0,
        "response":    0,
        "server_time": time.Now().Unix(),
    })
}

func GetFileList(c *fiber.Ctx) error {
	log.Info().Msg("POST /title/get_file_list")
	log.Info().Any("file_list", managers.GenerateFileList()).Msg("checking file list")
	configuration := config.GlobalConfig

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"master_ver":                 configuration.MasterVersion,
		"max_dl_stream_num_for_mtbl": 1,
		"max_dl_stream_num_for_dlc":  4,
		"status":                     0,
		"response":                   0,
		"file_list":                  managers.GenerateFileList(),
	})
}

func GetMasterTable(c *fiber.Ctx) error {
	tableNames := c.Context().QueryArgs().PeekMulti("table[]")
	requestedTables := make([]string, 0, len(tableNames))
	for _, tableName := range tableNames {
		name := string(tableName)
		requestedTables = append(requestedTables, name)
	}

	configuration := config.GlobalConfig
	rawBytes, readErr := os.ReadFile(configuration.MasterTableFilename)
	if readErr != nil {
		log.Error().Err(readErr).Msg("Failed to read master table file.")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":   1,
			"response": 1,
		})
	}

	var fullTable map[string]sonic.NoCopyRawMessage
	if unmarshalErr := sonic.Unmarshal(rawBytes, &fullTable); unmarshalErr != nil {
		log.Error().Err(unmarshalErr).Msg("Failed to unmarshal master table file.")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":   1,
			"response": 1,
		})
	}

	responseTable := fullTable
	if len(requestedTables) > 0 {
		responseTable = make(map[string]sonic.NoCopyRawMessage, len(requestedTables))
		for _, name := range requestedTables {
			if data, ok := fullTable[name]; ok {
				responseTable[name] = data
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"master_ver": configuration.MasterVersion,
		"status":     0,
		"response":   0,
		"table":      responseTable,
	})
}