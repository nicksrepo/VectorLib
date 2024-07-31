package vector

import "github.com/genovatix/vectorlib/utils"

type Index map[string]*KDNode

func NewIndexer() Index {

	i := make(map[string]*KDNode)
	return i

}

func (i Index) Add(topic string, interfaceID *KDNode) {
	err := utils.Revert(i.checkIndexExistence(topic), "should not exist")
	if err != nil {
		panic(err)
	}
	i[topic] = interfaceID
}

func (i Index) checkIndexExistence(topic string) bool {
	return i[topic] == nil
}

func (i Index) GC() {}
