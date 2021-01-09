
package proxima_db_client_go



//


//load()
func Load(proxima *proxima.ProximaServiceClient, fileLocation string, fileName string)  (*ProximaDatabase, error) {
  config, err := ConfigFromFile(fileLocation, fileName)
  if err != nil {
    return nil, err
  }
  return DBFromConfig(config) //sync
}
