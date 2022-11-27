# Configuration Implementation

## Contents: 
- `types.go`: defines the gossiper and monitor structs which will be be populated by config files (json files)
- `config_loaders.go`: contains functions that load configuration settings from json files
- `config/test`: contains the json files that are used for testing

## config_loaders.go:
- Loads the config files for both monitor and gossiper 
- Populates both monitor and gossiper configs from given file paths

## config_test.go
This tests that the information is populated correctly. Given the json files inside the folder /test, it will read and populate the respective structs
