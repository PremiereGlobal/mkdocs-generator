package main

import (
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func makeNav(docsDir string) {

	// Read in the mkdocs file
	mkdocsContents, err := ioutil.ReadFile(filepath.Join(docsDir, "mkdocs.yml"))
	if err != nil {
		log.Fatal("Unable to open mkdocs file: ", err)
	}

	// Ensure the file is in a good format
	var mkdocs map[string]interface{}
	err = yaml.Unmarshal([]byte(mkdocsContents), &mkdocs)
	if err != nil {
		log.Fatal("mkdocs file not valid yaml: ", err)
	}

	// Add the current timestamp to the mkdocs file
	// addTimestamp(mkdocs)

	// Creates the project index.md file and adds it to the nav
	createProjectIndex(mkdocs)

	// There is absolutely a better way to do these next two steps but this was
	// the culmination of many iterations of methods and I don't have time
	// to fix it now

	// Generate a nested map of the files
	navMap := make(map[interface{}]interface{})
	masterFileList.Range(func(key, value interface{}) bool {
		return generateNavMap(navMap, value.(*document))
	})

	// Generate the funky slice of map of slices structure that the nav requires
	nav := make(map[interface{}]interface{})
	generateNav("root", navMap, nav)
	for _, y := range nav["root"].([]map[interface{}]interface{}) {
		mkdocs["nav"] = append(mkdocs["nav"].([]interface{}), y)
	}

	// Marshal the mkdocs file back to yaml
	navData, err := yaml.Marshal(mkdocs)
	if err != nil {
		log.Fatal("Unable to generate mkdocs file: ", err)
	}

	// Write the mkdocs file to the build directory
	err = ioutil.WriteFile(filepath.Join(Args.GetString("build-dir"), "mkdocs.yml"), navData, 0644)
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
	var newMap map[interface{}][]interface{}

	// See if the pathPart exists in the current source
	found := false
	for i, v := range source {
		switch t := v.(type) {
		case map[interface{}]interface{}:
			for j, k := range t {
				if j.(string) == pathPart {
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
			result = append(source, item)
		}
	}

	return result
}

func generateNavMap(rootNav map[interface{}]interface{}, document *document) bool {

	var currentNav map[interface{}]interface{}
	currentNav = rootNav
	filePath := document.scmFilePath()

	parts := strings.Split(filePath, string(os.PathSeparator))

	// if this document is at the root of the project (path length 6), add it
	if len(parts) == 6 && document.docType == markdownType {
		repoPath := fmt.Sprintf("%s/%s", parts[1], parts[3])
		if _, ok := currentNav[repoPath]; !ok {
			currentNav[repoPath] = make(map[interface{}]interface{})
		}
		nextItem := currentNav[repoPath].(map[interface{}]interface{})
		nextItem[parts[5]] = filePath
	}

	// if this file is in the docs/ directory, add it
	if len(parts) == 7 && parts[5] == "docs" && document.docType == markdownType {
		repoPath := fmt.Sprintf("%s/%s", parts[1], parts[3])
		if _, ok := currentNav[repoPath]; !ok {
			currentNav[repoPath] = make(map[interface{}]interface{})
		}
		nextItem := currentNav[repoPath].(map[interface{}]interface{})
		if _, ok := nextItem["docs"]; !ok {
			nextItem["docs"] = make(map[interface{}]interface{})
		}
		nextItem2 := nextItem["docs"].(map[interface{}]interface{})
		nextItem2[parts[6]] = filePath
	}

	return true
}

func generateNav(name string, children interface{}, nav map[interface{}]interface{}) {

	// Assert the type of child(ren)
	switch c := children.(type) {

	// If the child is string, we've reached the end.  This is a file, add it as the value
	case string:
		basename := name[0:len(name)-len(filepath.Ext(name))]
		nav[basename] = c

	// If the child is a map, we've got more children to process
	case map[interface{}]interface{}:

		// Make our slice that will hold the child items
		var navChildren []map[interface{}]interface{}

		// Loop through each child
		for k, v := range c {

			// Make our map that will hold any items the child has
			childItems := make(map[interface{}]interface{})

			// Generate the nav for the child item
			generateNav(k.(string), v, childItems)

			// Add the processed child to the slice
			if len(childItems) > 0 {
				navChildren = append(navChildren, childItems)
			}
		}

		// After processing all the children, add it to the nav (if the children exist)
		if len(navChildren) > 0 {
			nav[name] = navChildren
		}
	}
}

func createProjectIndex(mkdocs map[string]interface{}) {
	// This should be done via some sort of markdown library?
	repoList := make(map[string]map[string]bool)
	var sb strings.Builder
	sb.WriteString("# Projects/Repos")
	sb.WriteString("\n")
	masterFileList.Range(func(key, value interface{}) bool {

		doc := value.(*document)

		// If the project doesn't exist, add it
		if _, ok := repoList[doc.project]; !ok {
			repoList[doc.project] = make(map[string]bool)
		}
		repoList[doc.project][doc.repo] = true

		return true
	})

	keys := make([]string, 0, len(repoList))
	for k := range repoList {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, projectName := range keys {
		sb.WriteString(fmt.Sprintf("##%s\n\n", projectName))
		repoKeys := make([]string, 0, len(repoList[projectName]))
		for k := range repoList[projectName] {
			repoKeys = append(repoKeys, k)
		}
		sort.Strings(repoKeys)
		for _, repoName := range repoKeys {
			sb.WriteString(fmt.Sprintf("[%s](/projects/%s/repos/%s/raw/)\n\n", repoName, projectName, repoName))
		}
	}

	// Write the repo index file
	err := ioutil.WriteFile(filepath.Join(Args.GetString("build-dir"), "docs/projects/index.md"), []byte(sb.String()), 0644)
	if err != nil {
		log.Fatal("Unable to write repo index file: ", err)
	}

	indexNav := make(map[interface{}]interface{})
	indexNav["Projects"] = "projects/index.md"
	mkdocs["nav"] = addNavItemToPath(indexNav, []string{"Projects"}, mkdocs["nav"].([]interface{}))
}

// addTimestamp adds the current time to theme.features.timestamp
func addTimestamp(mkdocs interface{}) {
	switch m := mkdocs.(type) {
	case map[string]interface{}:
		switch t := m["theme"].(type) {
		case map[interface{}]interface{}:
			switch f := t["feature"].(type) {
			case map[interface{}]interface{}:
				f["timestamp"] = time.Now().Format(time.RFC3339)
			}
		}
	}
}
