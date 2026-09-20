package managers

import (
	"os"

	"msattack/config"
	"github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
)

type FileData struct {
	FileName string json:"file_name"
	FileSize any    json:"file_size"
	Hash     string json:"hash"
	// "target" is "0" for normal files and "1" for master table files
	Target   any    json:"target"
	URL      string json:"url"
}

func GenerateFileList() []FileData {
	configuration := config.GlobalConfig
	var files []FileData

	fileListBytes, err := os.ReadFile(configuration.FileListFilename)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read file list.")
		return []FileData{}
	}

	err = sonic.Unmarshal(fileListBytes, &files)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal file list.")
		return []FileData{}
	}

	return files
}