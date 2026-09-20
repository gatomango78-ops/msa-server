package managers

import (
	"os"

	"msattack/config"
	"github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
)

type FileData struct {
    FileName string `json:"file_name"`
    FileSize int64  `json:"file_size"`
    Hash     string `json:"hash"`
    Target   int    `json:"target"`
    URL      string `json:"url"`
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
    if len(files) > 0 {
        log.Info().Str("sample_url", files[0].URL).Msg("checking first file url")
    }
	return files
}