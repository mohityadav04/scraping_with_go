package utils

type KVStore map[string] string

func GetConfig() KVStore {
	KVStoreInstance := make(KVStore)
	
	const (
		mongoPort = "27017"
		appPort = "5000"
	)

	KVStoreInstance["mongoPort"] = mongoPort
	KVStoreInstance["appPort"] = appPort
	KVStoreInstance["success"] = "SUCCESS"
	KVStoreInstance["failure"] = "FAILURE"
	
	return KVStoreInstance
}