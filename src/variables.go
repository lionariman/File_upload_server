package src

var (
	ChunkSize         int    = 1 << 20 // 1 mb
	FileName          string = "file.txt"
	DirName           string = "portions"
	PortionDirName    string = "1mb_"
	PortionDirNameTmp string = "1mb_"
	PortionName       string = "1mb_"
	PortName          string = ":8080"
	Endpoint          string = ""
	Endpoints         string = "\nENDPOINTS\n\n" +
		"/upload_file\n" +
		"/get_file:{file name}    (without braces)\n" +
		"/delete_file:{file name} (without braces)\n" +
		"/delete_all_files\n"
)
