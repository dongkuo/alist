package jni

type InitArgs struct {
	DataDir     string `json:"data_dir"`
	Debug       bool   `json:"debug"`
	NoPrefix    bool   `json:"no_prefix"`
	Dev         bool   `json:"dev"`
	ForceBinDir bool   `json:"force_bin_dir"`
	LogStd      bool   `json:"log_std"`
}
