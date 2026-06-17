package src

const (
	ChunkSize      int    = 1 << 20 // 1 mb
	DirName        string = "portions"
	PortionName    string = "chunk_"
	PortName       string = ":8080"
	MetadataName   string = "metadata.json"
)

var Endpoints = "\nENDPOINTS\n\n" +
	"/upload_file\n" +
	"/get_file:{file name}    (without braces)\n" +
	"/delete_file:{file name} (without braces)\n" +
	"/delete_all_files\n"
