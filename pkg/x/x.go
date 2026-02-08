package x

func SetupAll(logLevel string, aesKey string) {
	SetupLogger(logLevel)
	SetupValidator()
	AESKey = aesKey
}
