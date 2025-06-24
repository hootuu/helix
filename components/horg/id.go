package horg

import "github.com/hootuu/helix/components/htree"

var gOrgIdTree *htree.Tree

func initOrgIdTree() error {
	var err error
	gOrgIdTree, err = htree.NewTree("helix_org_tree", 1,
		[]uint{3, 4, 5, 6})
	if err != nil {
		return err
	}
	return nil
}
