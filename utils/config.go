package utils

type KVStore map[string] string

func GetConfig() KVStore {
	KVStoreInstance := make(KVStore)
	
	const (
		mongoAddress = "db:27017"
		appPort = "5000"
	)

	KVStoreInstance["mongoAddress"] = mongoAddress
	KVStoreInstance["appPort"] = appPort
	KVStoreInstance["success"] = "SUCCESS"
	KVStoreInstance["failure"] = "FAILURE"
	
	return KVStoreInstance
}