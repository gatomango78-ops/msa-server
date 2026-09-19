package title

import (
	"fmt"
	"os"

	"msattack/config"
	"msattack/errors"
	"msattack/managers"
	"msattack/utils"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

const PackInfoURL = `https://%s%s/pack/%d/`

func GetPackInfo(c *fiber.Ctx) error {
    log.Info().Msg("POST /title/get_pack_info")

    configuration := config.GlobalConfig
    actualPackInfoURL := "/snkp/msatk/prod/pack/6120000/pack_info_list.txt"

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "version":  configuration.PackVersion,
        "url":      actualPackInfoURL,
        "response": utils.GenerateErrorCode(errors.SUCCESS),
    })
}

func GetFileList(c *fiber.Ctx) error {
	log.Info().Msg("POST /title/get_file_list")

	configuration := config.GlobalConfig

	err := c.Status(fiber.StatusOK).JSON(fiber.Map{
		"master_ver": configuration.MasterVersion,
		// Not sure what these are for (parallel download limits for master table and DLC?)
		"max_dl_stream_num_for_mtbl": 1,
		"max_dl_stream_num_for_dlc":  4,
		"response":                   utils.GenerateErrorCodeWithTime(errors.SUCCESS),
		"file_list":                  managers.GenerateFileList(),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Error in getting file list.")
	}

	c.Set("Connection", "close")
	c.Set("Server", "Apache")
	c.Set("Vary", "Accept-Encoding")

	return err
}

func GetMasterTable(c *fiber.Ctx) error {
	// Since we can potentially have multiple "table[]" query parameters, we iterate over them to get the full list
	tableNames := c.Context().QueryArgs().PeekMulti("table[]")

	requestedTables := make([]string, 0, len(tableNames))
	for _, tableName := range tableNames {
		name := string(tableName)
		requestedTables = append(requestedTables, name)
		log.Info().Msgf("GET /title/get_master_table?table[]=%s", name)
	}

	configuration := config.GlobalConfig

	rawBytes, readErr := os.ReadFile(configuration.MasterTableFilename)
	if readErr != nil {
		log.Error().Err(readErr).Msg("Failed to read master table file.")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"response": utils.GenerateErrorCode(errors.MAINTENANCE),
		})
	}

	// The master table file is expected to be a JSON object keyed by table name,
	// e.g. {"chara_m": [...], "weapon_m": [...], ...}
	var fullTable map[string]sonic.NoCopyRawMessage
	if unmarshalErr := sonic.Unmarshal(rawBytes, &fullTable); unmarshalErr != nil {
		log.Error().Err(unmarshalErr).Msg("Failed to unmarshal master table file.")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"response": utils.GenerateErrorCode(errors.MAINTENANCE),
		})
	}

	// If specific tables were requested, only return those; otherwise return everything.
	responseTable := fullTable
	if len(requestedTables) > 0 {
		responseTable = make(map[string]sonic.NoCopyRawMessage, len(requestedTables))
		for _, name := range requestedTables {
			if data, ok := fullTable[name]; ok {
				responseTable[name] = data
			} else {
				log.Warn().Msgf("Requested table %q not found in master table file.", name)
			}
		}
	}

	err := c.Status(fiber.StatusOK).JSON(fiber.Map{
		"master_ver": configuration.MasterVersion,
		"response":   utils.GenerateErrorCodeWithTime(errors.SUCCESS),
		"table":      responseTable,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Error in getting master table data.")
	}

	c.Set("Connection", "close")
	c.Set("Server", "Apache")
	c.Set("Vary", "Accept-Encoding")

	return err
}