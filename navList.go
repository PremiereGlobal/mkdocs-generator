package main

type NavList []map[interface{}]interface{}

func (n NavList) Len() int {
	return len(n)
}

func (n NavList) Less(i, j int) bool {
	// There should be only one key per map, if there is more than one, well, we're screwed
	a := ""
	for key, _ := range n[i] {
		a = key.(string)
	}

	b := ""
	for key, _ := range n[j] {
		b = key.(string)
	}

	return a < b
}

func (n NavList) Swap(i, j int) {
	a := n[i]
	b := n[j]
	n[i] = b
	n[j] = a
}
