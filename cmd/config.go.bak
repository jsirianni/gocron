package cmd
import (
    "os"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)


// Read in the config file
func GetConfig(verbose bool) Config {
      var config Config
      yamlFile, err := ioutil.ReadFile(cfgFile)
      if err != nil {
           CheckError(err, verbose)
           os.Exit(1)
      }

      err = yaml.Unmarshal(yamlFile, &config)
      if err != nil {
            CheckError(err, verbose)
            os.Exit(1)
      }

      return config
}
