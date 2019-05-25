package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
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

	spew.Dump(mkdocs)

	// tree := make(map[string]interface{})
	// tree := makeKey("key1", nil, mkdocs)
	// tree2 := makeKey("key2", nil, tree)
	// tree3 := makeKey("key3", nil, tree2)
	// p := makeKey("key1", &mkdocs)
	// p = makeKey("key2", &p)

	var positionHolder interface{}
	// positionHolder = mkdocs
	positionHolder = makeKey("nav", nil, mkdocs)
	keyParts := strings.Split(mkdocsKey, ".")
	for _, k := range keyParts {
		positionHolder = makeKey(k, nil, positionHolder)
		fmt.Println()
		spew.Dump(positionHolder)
		fmt.Println()
		// if i == len(keyParts)-1 {
		// 	rootNavKey = k
		// } else {
		// 	mkdocsPointer = final(k, mkdocsPointer)
		// 	// final(k, &foo)
		// }
	}

	spew.Dump(mkdocs)

	// Get the item in the path that we want

	mkdocsPointer := make(map[interface{}]interface{})

	// mkdocsPointer = mkdocs
	rootNavKey := ""

	// navRoot := mkdocs["nav"]

	// foo := navRoot

	// spew.Dump(foo)
	// x, _ := navRoot.([]interface{})
	//     element := navAddElement("test", &navRoot)
	//     element["test"] = "foo"
	//
	//
	//     // element2 := getNavElement("Projects", &navRoot)
	//     // g := element2.([]interface {})
	//     // g[0] = "bar"
	//     // element2 = "foo2"
	//
	//     element3 := getNavElement("test1", &navRoot)
	//     // h := element3.(map[interface{}]interface{})
	//     element3["asdf"] = "asdf"
	//
	//     navRoot = append(navRoot.([]interface{}), "blah")
	//
	// fmt.Println()
	// fmt.Println()
	// fmt.Println()

	// mkdocs["nav"] = navRoot
	// spew.Dump(mkdocs)

	// xmkdocsPointer := mkdocs["nav"]
	// var navPointer interface{}
	// navPointer = mkdocs["nav"]
	// g := navPointer
	//
	// // Loop through our keyParts (strings) and ensure they exist
	// for _, k := range keyParts {
	//
	//   // Ensure this part of the tree is in the right format
	//   element, ok := g.([]interface{})
	//   if !ok {
	//     log.Warn("mkdocs file located at ", mkdocsFilePath, " is in a bad format")
	//     continue
	//   }
	//
	//   // Loop through the elements in this position and see if we already have our
	//   // key
	//   found := false
	//
	//   // Here we're looping through each ???
	//   for _, y := range element {
	//
	//     foo := y.(map[interface {}]interface {})
	//     for z := range foo {
	//       // If the key is the one we're looking for, link and continue
	//       key := z.(string)
	//       if key == k {
	//         log.Debug("Found existing key part in mkdocs ", k)
	//         found = true
	//       }
	//     }
	//   }
	//
	//   if !found {
	//     log.Debug("Didn't find key part in mkdocs ", k, " creating")
	//
	//     // Slice that will hold our maps
	//     // var q []interface{}
	//     // q = make([]interface{}, 1)
	//     //
	//     var newMap map[interface {}]interface {}
	//     // newMap = make(map[interface {}]interface {})
	//     // newMap[k] = make([]interface {}, 1)
	//     // var newElement map[interface {}]interface {}
	//     //
	//     //
	//     //
	//     // w := make(map[interface {}]interface {})
	//     // w[k] = q
	//
	//
	//
	//     // newElement = make(map[interface {}]interface {})
	//     // newElement[k] = c
	//     // b := make([]interface{}, 1)
	//     // b[0] = make(map[interface {}]interface {})
	//     // newElement[k] = b
	//
	//     // spew.Dump(newMap)
	//
	//     element = append(element, newMap)
	//     g = newMap[k]
	//     // navPointer = append(navPointer.([]interface{}), newMap)
	//     // navPointer = x
	//     // spew.Dump(navPointer)
	//     // Set our nav pointer to point to the newly added element
	//     // navPointer = newElement[k]
	//   }
	//
	//   // spew.Dump("")
	//   // spew.Dump(mkdocs)
	//
	//
	//   // element := make(map[string]interface{})
	//   //
	//   // generateNav(name string, children interface{}, nav map[string]interface{})
	//   //
	//   // // Ensure this element is of the right type
	//   // switch c := xmkdocsPointer.(type) {
	//   //   case []map[string]interface{} :
	//   //   // good
	//   //     log.Fatal("good", c)
	//   //   default:
	//   //   // bad
	//   //     log.Fatal("bad", c)
	//   //     // log.Fatal("mkdocs file located at ", mkdocsFilePath, " in a bad format")
	//   // }
	//
	//   // generateNav(v, nil, xmkdocsPointer.(map[string]interface{}))
	//   // _, okExists := mkdocsPointer[v]
	//   // if i == len(keyParts) - 1 {
	//   //   rootNavKey = v
	//   // } else {
	//   //
	//   //   if !okExists {
	//   //     log.Debug("Nav item ", v, " doesn't exist or isn't the right type, overwriting")
	//   //     mkdocsPointer[v] = make(map[string]interface{})
	//   //   }
	//   //   nextItem := mkdocsPointer[v].(map[string]interface{})
	//   //   mkdocsPointer = nextItem
	//   // }
	// }
	//
	//
	// mkdocs["nav"] = navPointer
	// spew.Dump(mkdocs)

	// var rootNav map[string]interface{}
	// rootNav = make(map[string]interface{})
	navMap := make(map[string]interface{})
	masterFileList.Range(func(key, value interface{}) bool {
		return generateNavMap(navMap, key.(string))
	})
	// spew.Dump(navMap)
	// Finally, loop through and make the project items into arrays
	// var generatedArrayNav map[string][]map[string]interface{}
	// generatedArrayNav := make(map[string][]map[string]interface{})
	nav := make(map[string]interface{})
	// generatedArrayNav["start"] =
	// i := 0

	generateNav(rootNavKey, navMap, nav)
	// for k, _ := range generatedNav {
	//   // generatedArrayNav[0] = make(m)
	//   generateArrayNav(item map[string]interface, depth int)
	//   generatedArrayNav = append(generatedArrayNav, k)
	//   // i += 1
	// }

	// spew.Dump(nav)
	// spew.Dump(nav)
	mkdocsPointer[rootNavKey] = nav[rootNavKey]
	// spew.Dump(nav[rootNavKey])
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

// Take in a map`
// Add the key to a map
// The value is a new []map[string]interface{}
func makeKey(name string, value interface{}, source interface{}) interface{} {

	newMap := make(map[string]interface{})
	newSlice := make([]map[string]interface{}, 1)
	newSlice[0] = newMap

	switch s := source.(type) {
	case map[interface{}]interface{}:
		s[name] = newSlice
		// return s[name][0]
	case map[string]interface{}:
		s[name] = newSlice
		switch t := s[name].(type) {
		case []map[string]interface{}:
			return t[0]
		}
	}

	return nil
}

func makeIndex(source interface{}) interface{} {
	switch s := source.(type) {
	case []interface{}:
		item := make(map[interface{}]interface{})
		s = append(s, item)
		return item[len(s)-1]
	}

	return nil
}

// navAddElement adds an element ot the nav
func navAddElement(name string, nav *interface{}) map[interface{}]interface{} {

	element := make(map[interface{}]interface{})

	b := *nav
	*nav = append(b.([]interface{}), element)

	return element

}

func final(name string, nav map[interface{}]interface{}) map[interface{}]interface{} {

	if _, ok := nav[name]; !ok {
		x := make([]map[interface{}]interface{}, 1)
		x[0] = make(map[interface{}]interface{})

		// nav[name] =
		return x[0]
	}

	return nil

}

func generateNavMap(rootNav map[string]interface{}, filePath string) bool {

	var currentNav map[string]interface{}
	currentNav = rootNav

	parts := strings.Split(filePath, string(os.PathSeparator))

	for i, p := range parts {

		// File path should be "/projects/<project>/repos/<repos>/<raw/browse>/filepath..."
		if i != 0 && i != 2 && i != 4 {
			if i == len(parts)-1 {
				currentNav[p] = p
			} else {
				if _, ok := currentNav[p]; !ok {
					currentNav[p] = make(map[string]interface{})
				}
				nextItem := currentNav[p].(map[string]interface{})
				currentNav = nextItem
			}
		}
	}

	return true
}

func generateNav(name string, children interface{}, nav map[string]interface{}) {

	// Assert the type of child(ren)
	switch c := children.(type) {

	// If the child is string, we've reached the end.  This is a file, add it as the value
	case string:
		nav[name] = c

		// If the child is a map, we've got more children to process
	case map[string]interface{}:

		// Make our slice that will hold the child items
		navChildren := make([]map[string]interface{}, len(c))

		j := 0
		// Loop through each child
		for k, v := range c {

			// Make our map that will hold any items the child has
			childItems := make(map[string]interface{})

			// Generate the nav for the child item
			generateNav(k, v, childItems)

			// Add the processed child to the slice
			navChildren[j] = childItems

			j = j + 1
		}

		// After processing all the children, add it to the nav
		nav[name] = navChildren
	}

}
