package jni

type InitArgs struct {
	DataDir     string `json:"dataDir"`
	Debug       bool   `json:"debug"`
	NoPrefix    bool   `json:"NoPrefix"`
	Dev         bool   `json:"dev"`
	ForceBinDir bool   `json:"forceBinDir"`
	LogStd      bool   `json:"logStd"`
}
