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
	log.Info().Str("raw_body", string(c.Body())).Msg("client request body")
	log.Info().Msg("POST /title/get_file_list")
	log.Info().Any("file_list", managers.GenerateFileList()).Msg("checking file list")
	
	dump, _ := sonic.MarshalString(managers.GenerateFileList())
	log.Info().Str("json_dump", dump).Msg("inspect file list")
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"master_ver":                 7130000,
		"dl_mtbl_merge_lim_size":     2097152,
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
		requestedTables = append(requestedTables, string(tableName))
	}

	configuration := config.GlobalConfig
	rawBytes, readErr := os.ReadFile(configuration.MasterTableFilename)
	if readErr != nil {
		log.Error().Err(readErr).Msg("Failed to read master table file.")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": 1, "response": 1})
	}

	var fullTable map[string]sonic.NoCopyRawMessage
	if unmarshalErr := sonic.Unmarshal(rawBytes, &fullTable); unmarshalErr != nil {
		log.Error().Err(unmarshalErr).Msg("Failed to unmarshal master table file.")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": 1, "response": 1})
	}

	responseTable := make(map[string]sonic.NoCopyRawMessage)
	if len(requestedTables) > 0 {
		for _, name := range requestedTables {
			if data, ok := fullTable[name]; ok {
				responseTable[name] = data
			}
		}
	}
	// Fallback de seguridad: si las tablas pedidas vienen vacías o no matchean, devolvemos todo para evitar crash del cliente
	if len(responseTable) == 0 {
		responseTable = fullTable
	}
    responseBytes, _ := sonic.Marshal(responseTable)
	    log.Info().Int("response_table_bytes", len(responseBytes)).Msg("sending master_table payload")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"master_ver": 7130000,
		"status":     0,
		"response":   0,
		"table":      responseTable,
	})
}