




//

//load()
func Load(proxima *proxima.ProximaServiceClient, fileLocation string, fileName string)  (*ProximaDatabase, error) {
  config, err := ConfigFromFile(fileLocation, fileName)
  if err != nil {
    return nil, err
  }
  return DBFromConfig(config) //sync
}

//Load from file
//UGH ok, the yaml file needs to be correctly structured
//Save to file
func Save() {
//save files to db
}

func Update() {
  
}

func Sync(proximaApp *proximaApp) (error) {
  //compare app with
//compare the table version //time-based
//compare the table header
//select the correct table and update
//write the table to file
//write the table to application database
//if incorrect, get the correct data
}


func Remove() {

}


func Compare() {

}

func ConfigFromFile(fileLocation string, fileName string) (map[string]interface{}, error) {

}

func Start(proxima *proxima.ProximaServiceClient, fileLocation string, fileName string) {

}

func Stop(proxima *proxima.ProximaServiceClient, fileLocation string, fileName string) {

}
