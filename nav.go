package main

import (
	// "fmt"
	"github.com/davecgh/go-spew/spew"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	// "reflect"
	"strings"
  "fmt"
)

func makeNav(mkdocsFilePath string, mkdocsKey string) {

	// Read in the mkdocs file
	mkdocsContents, err := ioutil.ReadFile(mkdocsFilePath)
	if err != nil {
		log.Fatal("Unable to open mkdocs file: ", err)
	}

	// Ensure the file is in a good format
	var mkdocs map[string]interface{}
	err = yaml.Unmarshal([]byte(mkdocsContents), &mkdocs)
	if err != nil {
		log.Fatal("mkdocs file not valid yaml: ", err)
	}

	// var positionHolder []interface{}

	// Make the root "nav" key if it doesn't exist
	// positionHolder = makeKey("nav", mkdocs["nav"])
  // positionHolder = append(positionHolder.([]interface{}), "test")
  // spew.Dump(mkdocs)
  // positionHolder = mkdocs["nav"].([]interface{})
  // x = mkdocs
  // spew.Dump(&x)
  // spew.Dump(&mkdocs["nav"])
  // mkdocs["nav"] = addNavItemToPath("test", []string{"test1"}, mkdocs["nav"].([]interface{}))
  // mkdocs["nav"] = append(positionHolder, "test")
	// Split the desired key path into its parts
	keyParts := strings.Split(mkdocsKey, ".")
  testMap := make(map[interface{}]interface{})
  testMap["nav"] = "test"

	rootNavKey := ""
	for i, k := range keyParts {
		if i == len(keyParts)-1 {
			rootNavKey = k
		} else {
			// positionHolder = makeKey(k, positionHolder)
      // spew.Dump(positionHolder)
      fmt.Println()
		}
	}

  // spew.Dump(positionHolder)



	navMap := make(map[interface{}]interface{})
	masterFileList.Range(func(key, value interface{}) bool {
		return generateNavMap(navMap, key.(string))
	})
	//
	// Finally, loop through and make the project items into arrays
	// var generatedArrayNav map[string][]map[string]interface{}
	// generatedArrayNav := make(map[string][]map[string]interface{})
	nav := make(map[interface{}]interface{})
	generateNav(rootNavKey, navMap, nav)

  mkdocs["nav"] = addNavItemToPath(nav, keyParts, mkdocs["nav"].([]interface{}))
  spew.Dump(mkdocs)

	// positionHolder.(map[interface{}]interface{})[rootNavKey] = nav[rootNavKey]
	navData, err := yaml.Marshal(mkdocs)
	if err != nil {
		log.Fatal("Unable to generate mkdocs file: ", err)
	}

	// Write the mkdocs file
	err = ioutil.WriteFile("testmkdocs.yml", navData, 0644)
	if err != nil {
		log.Fatal("Unable to write mkdocs file: ", err)
	}
}


// Add the item to the given path with source as the existing structure
func addNavItemToPath(item map[interface{}]interface{}, path []string, source []interface{}) []interface{} {

  // Pop item off path
  pathPart, newPath := path[0], path[1:]

  log.Info("Adding item to ", pathPart)

  var newSource []interface{}
  var result []interface{}

  // Exit condition
  // if len(newPath) > 0 {

    var newMap map[interface{}][]interface{}

    // See if the pathPart exists in the current source
    found := false
    // sourceIndex := 0
    for i, v := range source {
      switch t := v.(type) {
      case map[interface{}]interface{}:
        for j, k := range t {
          // log.Info("Here's the nav item in position ", i, ": ", j.(string))
          if j.(string) == pathPart {

            // newMap := v

            log.Debug("Nav key '", pathPart, "' exists already in index ", i)
            found = true
            result = source
            if len(newPath) > 0 {
              newMap := t
              newMap[j] = addNavItemToPath(item, newPath, k.([]interface{}))
              result[i] = newMap
            } else {
              log.Warn("Replacing existing navigation key '", pathPart, "' with mkdocs structure")
              log.Debug("Adding mkdocs structure to '", pathPart, "'")

              // // Last condition returns the item we're adding
              // newMap = make(map[interface{}][]interface{})
              // newMap[pathPart] =
              // result = append(source, newMap)
              //
              //
              // result = make([]interface{}, 1)
              // result[0] = item
              result[i] = item

            }

          }
        }
      default:
        log.Warn("Nav index ", i, " is the wrong type")
      }
    }

    // If we didn't find the key, make a new one
    if !found {

      if len(newPath) > 0 {
        newMap = make(map[interface{}][]interface{})
        newMap[pathPart] = addNavItemToPath(item, newPath, newSource)
        result = append(source, newMap)
      } else {
        log.Debug("Adding mkdocs structure to '", pathPart, "'")

        // Last condition returns the item we're adding
        // newMap = item
        result = append(source, item)

        // result = make([]interface{}, 1)
        // result[0] = item
      }
    }

  // } else {

    // if len(source) > 0 {
    //   log.Warn("Replacing existing navigation key '", pathPart, "' with mkdocs structure")
    // }
    //
    // log.Debug("Adding mkdocs structure to '", pathPart, "'")
    //
    // // Last condition returns the item we're adding
    // newMap = make(map[interface{}][]interface{})
    // newMap[pathPart] =
    // result = append(source, newMap)
    //
    //
    // result = make([]interface{}, 1)
    // result[0] = item
  // }

  return result
}

// Take in a map`
// Add the key to a map
// The value is a new []map[string]interface{}
func makeKey(name string, source []interface{}) []interface{} {

  var newMap map[interface{}][]interface{}
	newMap = make(map[interface{}][]interface{})
	// newSlice := make([]map[interface{}]interface{}, 1)
  var newSlice []interface{}
  // newSlice[0] = make(map[interface{}]interface{})
  newMap[name] = newSlice

	// switch s := source.(type) {

	// This case statement should match all maps under, but not including, the
	// root level "nav" key
	// case map[interface{}]interface{}:
	// 	if _, ok := s[name]; ok {
	// 		log.Info("key ", name, " exists already")
	// 		switch t := s[name].(type) {
	// 		case []interface{}:
	// 			log.Info("Adding item to ", name)
	// 			t = append(t, newMap)
	// 			return newMap
	// 		default:
	// 			log.Fatal("key ", name, " is the wrong type ", reflect.TypeOf(s[name]))
	// 		}
	// 	}
	// 	s[name] = newSlice
	// 	switch t := s[name].(type) {
	// 	case []map[interface{}]interface{}:
	// 		return t[0]
	// 	}


  // case []interface{}:
    // spew.Dump(s)
    found := false
    for i, v := range source {
      switch t := v.(type) {
			case map[interface{}]interface{}:
        for j, _ := range t {
          log.Info("Here's a nav item in position ", i, ": ", j.(string))
          if j.(string) == name {
            log.Info("slice key ", name, " exists already")
            found = true
          }
        }
		  default:
		    log.Warn("index ", i, " is the wrong type")
			}
    }

    if !found {

     source = append(source, newMap)
     // spew.Dump(s)
     // fmt.Println()
     // b := newMap[name]
     return newMap[name]
     // return

    }

		// This case statement is only for the root level "nav" key, all other keys
		// under "nav" should be map[interface{}]interface{}
	// case map[string]interface{}:
	// 	// See if the key already exists
	// 	if _, ok := s[name]; ok {
  //     log.Info("key ", name, " exists already")
  //     switch s[name].(type) {
	// 		case []interface{}:
	// 			return s[name]
  //
	// 		// return s[name]
  //
	// 		// switch t := s[name].(type) {
	// 		//   case []interface{}:
	// 		//     log.Info("Adding item to ", name)
	// 		//     t = append(t, newMap)
	// 		//     spew.Dump(t)
	// 		//     return t[len(t)-1]
	// 		  default:
	// 		    log.Fatal("key ", name, " is the wrong type ", reflect.TypeOf(s[name]))
	// 		}
	// 	} else {
	// 		s[name] = newSlice
	// 	}
	// 	switch t := s[name].(type) {
	// 	case []map[interface{}]interface{}:
	// 		return t[0]
	// 	case []interface{}:
	// 		return t[0]
	// 	default:
	// 		log.Fatal("key ", name, " is the wrong type ", reflect.TypeOf(s[name]))
	// 	}
  // default:
  //   log.Fatal("Wrong type adding ", name, ": ", reflect.TypeOf(s))
	// }

	return nil
}

func generateNavMap(rootNav map[interface{}]interface{}, filePath string) bool {

	var currentNav map[interface{}]interface{}
	currentNav = rootNav

	parts := strings.Split(filePath, string(os.PathSeparator))

	for i, p := range parts {

		// File path should be "/projects/<project>/repos/<repos>/<raw/browse>/filepath..."
		if i != 0 && i != 2 && i != 4 {
			if i == len(parts)-1 {
				currentNav[p] = p
			} else {
				if _, ok := currentNav[p]; !ok {
					currentNav[p] = make(map[interface{}]interface{})
				}
				nextItem := currentNav[p].(map[interface{}]interface{})
				currentNav = nextItem
			}
		}
	}

	return true
}

func generateNav(name string, children interface{}, nav map[interface{}]interface{}) {

	// Assert the type of child(ren)
	switch c := children.(type) {

	// If the child is string, we've reached the end.  This is a file, add it as the value
	case string:
		nav[name] = c

	// If the child is a map, we've got more children to process
	case map[interface{}]interface{}:

		// Make our slice that will hold the child items
		navChildren := make([]map[interface{}]interface{}, len(c))

		j := 0
		// Loop through each child
		for k, v := range c {

			// Make our map that will hold any items the child has
			childItems := make(map[interface{}]interface{})

			// Generate the nav for the child item
			generateNav(k.(string), v, childItems)

			// Add the processed child to the slice
			navChildren[j] = childItems

			j = j + 1
		}

		// After processing all the children, add it to the nav
		nav[name] = navChildren
	}

}
