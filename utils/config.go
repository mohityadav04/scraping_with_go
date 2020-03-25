package utils

type KVStore map[string] string

func GetConfig() KVStore {
	KVStoreInstance := make(KVStore)
	
	const (
		mongoAddress = "localhost:27017"
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