package utils

type KVStore map[string] string

func GetConfig() KVStore {
	KVStoreInstance := make(KVStore)
	
	const (
		mongoAddress = "mongo:8888"
		appPort = "5000"
		appHost = "localhost"
		mongoDatabase = "amazon"
	)

	KVStoreInstance["mongoAddress"] = mongoAddress
	KVStoreInstance["appHost"] = appHost
	KVStoreInstance["appPort"] = appPort
	KVStoreInstance["mongoDatabase"] = mongoDatabase

	
	return KVStoreInstance
}