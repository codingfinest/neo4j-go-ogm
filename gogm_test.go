// MIT License
//
// Copyright (c) 2020 codingfinest
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package gogm_test

import (
	"testing"
	"time"

	"github.com/neo4j/neo4j-go-driver/neo4j"

	"github.com/codingfinest/neo4j-go-ogm/gogm"
	. "github.com/codingfinest/neo4j-go-ogm/gogm/tests/models"
	. "github.com/onsi/gomega"
)

var ogm = gogm.New("bolt://localhost:7687", "neo4j", "Pass1234")
var session, err = ogm.NewSession(true)

const deletedID int64 = -1

var eventListener = &TestEventListener{}

//Context: when simple node object is saved
//Spec: should have metadata properties updated
func TestNodeSave(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	//Context
	simpleNode := &SimpleNode{}
	g.Expect(session.Save(&simpleNode, nil)).NotTo(HaveOccurred())

	//Spec
	g.Expect(simpleNode.CreatedAt).ToNot(BeZero())
	g.Expect(simpleNode.DeletedAt).To(BeZero())
	g.Expect(simpleNode.UpdatedAt).To(BeZero())
	g.Expect(*simpleNode.ID > -1).To(BeTrue())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

//Context:when simple realtionship object is saved
//Spec:should have metadata properties updated
func TestRelationshipSave(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	//Context
	simpleRelationship := &SimpleRelationship{}
	n4 := &Node4{}
	n5 := &Node5{}
	simpleRelationship.N4 = n4
	simpleRelationship.N5 = n5

	g.Expect(session.Save(&simpleRelationship, nil)).NotTo(HaveOccurred())

	//Spec
	g.Expect(simpleRelationship.CreatedAt).ToNot(BeZero())
	g.Expect(simpleRelationship.DeletedAt).To(BeZero())
	g.Expect(simpleRelationship.UpdatedAt).To(BeZero())
	g.Expect(*simpleRelationship.ID > -1).To(BeTrue())
	g.Expect(*n4.ID > -1).To(BeTrue())
	g.Expect(*n5.ID > -1).To(BeTrue())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestNodeSaveWithoutUpdate(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	//Context
	simpleNode := &SimpleNode{}
	g.Expect(session.Save(&simpleNode, nil)).NotTo(HaveOccurred())
	g.Expect(session.Save(&simpleNode, nil)).NotTo(HaveOccurred())

	//Spec
	g.Expect(*simpleNode.ID > -1).To(BeTrue())
	g.Expect(simpleNode.UpdatedAt).To(BeZero())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestRelationshipSaveWithoutUpdate(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	//Context
	simpleRelationship := &SimpleRelationship{}
	n5 := &Node5{}
	n4 := &Node4{}
	simpleRelationship.N4 = n4
	simpleRelationship.N5 = n5
	g.Expect(session.Save(&simpleRelationship, nil)).NotTo(HaveOccurred())
	g.Expect(session.Save(&simpleRelationship, nil)).NotTo(HaveOccurred())

	//Spec
	g.Expect(*simpleRelationship.ID > -1).To(BeTrue())
	g.Expect(*n4.ID > -1).To(BeTrue())
	g.Expect(*n5.ID > -1).To(BeTrue())

	g.Expect(simpleRelationship.UpdatedAt).To(BeZero())
	g.Expect(n4.UpdatedAt).To(BeZero())
	g.Expect(n5.UpdatedAt).To(BeZero())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestNodeSaveWithUpdate(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	simpleNode := &SimpleNode{}
	g.Expect(session.Save(&simpleNode, nil)).NotTo(HaveOccurred())
	simpleNode.Prop1 = "test Prop"
	g.Expect(session.Save(&simpleNode, nil)).NotTo(HaveOccurred())

	g.Expect(simpleNode.UpdatedAt).ToNot(BeZero())
	g.Expect(*simpleNode.ID > -1).To(BeTrue())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestRelationshipSaveWithUpdate(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	simpleRelationship := &SimpleRelationship{}
	n5 := &Node5{}
	n4 := &Node4{}
	simpleRelationship.N4 = n4
	simpleRelationship.N5 = n5
	g.Expect(session.Save(&simpleRelationship, nil)).NotTo(HaveOccurred())
	simpleRelationship.Name = "test Prop"
	g.Expect(session.Save(&simpleRelationship, nil)).NotTo(HaveOccurred())

	g.Expect(simpleRelationship.UpdatedAt).ToNot(BeZero())
	g.Expect(*simpleRelationship.ID > -1).To(BeTrue())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestSaveSliceOfNode(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	simpleNode1 := SimpleNode{}
	simpleNode2 := SimpleNode{}

	simpleNodes := []*SimpleNode{&simpleNode1, &simpleNode2}
	g.Expect(session.Save(&simpleNodes, nil)).NotTo(HaveOccurred())

	g.Expect(simpleNode1.CreatedAt).NotTo(BeZero())
	g.Expect(simpleNode1.DeletedAt).To(BeZero())
	g.Expect(simpleNode1.UpdatedAt).To(BeZero())
	g.Expect(*simpleNode1.ID > -1).To(BeTrue())
	g.Expect(simpleNode2.CreatedAt).NotTo(BeZero())
	g.Expect(simpleNode2.DeletedAt).To(BeZero())
	g.Expect(simpleNode2.UpdatedAt).To(BeZero())
	g.Expect(*simpleNode2.ID > -1).To(BeTrue())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestSaveSliceOfRelationship(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	simpleRelationship1 := SimpleRelationship{}
	n5 := &Node5{}
	n4 := &Node4{}
	simpleRelationship1.N4 = n4
	simpleRelationship1.N5 = n5

	simpleRelationship2 := SimpleRelationship{}
	simpleRelationship2.N4 = n4
	simpleRelationship2.N5 = n5

	simpleRelationships := []*SimpleRelationship{&simpleRelationship1, &simpleRelationship2}

	g.Expect(session.Save(&simpleRelationships, nil)).NotTo(HaveOccurred())

	g.Expect(simpleRelationship1.CreatedAt).NotTo(BeZero())
	g.Expect(simpleRelationship1.DeletedAt).To(BeZero())
	g.Expect(simpleRelationship1.UpdatedAt).To(BeZero())
	g.Expect(*simpleRelationship1.ID > -1).To(BeTrue())

	g.Expect(simpleRelationship2.CreatedAt).NotTo(BeZero())
	g.Expect(simpleRelationship2.DeletedAt).To(BeZero())
	g.Expect(simpleRelationship2.UpdatedAt).To(BeZero())
	g.Expect(*simpleRelationship2.ID > -1).To(BeTrue())

	g.Expect(n4.CreatedAt).NotTo(BeZero())
	g.Expect(n4.DeletedAt).To(BeZero())
	g.Expect(n4.UpdatedAt).To(BeZero())
	g.Expect(*n4.ID > -1).To(BeTrue())

	g.Expect(n5.CreatedAt).NotTo(BeZero())
	g.Expect(n5.DeletedAt).To(BeZero())
	g.Expect(n5.UpdatedAt).To(BeZero())
	g.Expect(*n5.ID > -1).To(BeTrue())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestNodeDelete(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	simpleNode := &SimpleNode{}
	g.Expect(session.Save(&simpleNode, nil)).NotTo(HaveOccurred())
	g.Expect(session.Delete(&simpleNode)).NotTo(HaveOccurred())
	simpleNode.Prop1 = "test"
	g.Expect(session.Save(&simpleNode, nil)).NotTo(HaveOccurred())

	g.Expect(simpleNode.CreatedAt).ToNot(BeZero())
	g.Expect(*simpleNode.ID).To(Equal(deletedID), "Deleted node isn't re-saved")

	g.Expect(simpleNode.DeletedAt).ToNot(BeZero())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestRelationshipDelete(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	simpleRelationship := &SimpleRelationship{}
	n5 := &Node5{}
	n4 := &Node4{}
	simpleRelationship.N4 = n4
	simpleRelationship.N5 = n5
	g.Expect(session.Save(&simpleRelationship, nil)).NotTo(HaveOccurred())
	g.Expect(session.Delete(&simpleRelationship)).NotTo(HaveOccurred())
	simpleRelationship.Name = "test"

	g.Expect(*simpleRelationship.ID).To(Equal(deletedID), "Deleted relationship isn't re-saved")
	g.Expect(*n4.ID > deletedID).To(BeTrue())
	g.Expect(*n5.ID > deletedID).To(BeTrue())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

//Context: When a node object is removed from a parent node entity
//Spec: Corresponding relationship should be removed
func TestRemoveNode(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())

	lo := gogm.NewLoadOptions()
	lo.Depth = -1

	n3 := &Node3{}
	n4 := &Node4{}
	n4.N3 = n3

	n4.Name = "N4"
	n3.Name = "N3"
	g.Expect(session.Save(&n4, nil)).NotTo(HaveOccurred())
	var loadedN4 *Node4
	var loadedN3 *Node3

	n4.N3 = nil
	g.Expect(session.Save(&n4, nil)).NotTo(HaveOccurred())

	g.Expect(session.Clear()).NotTo(HaveOccurred())

	g.Expect(session.Load(&loadedN4, *n4.ID, lo)).NotTo(HaveOccurred())
	g.Expect(session.Load(&loadedN3, *n3.ID, lo)).NotTo(HaveOccurred())

	g.Expect(loadedN4).To(Equal(n4))
	g.Expect(loadedN3).To(Equal(n3))
	g.Expect(*loadedN4).To(Equal(*n4))
	g.Expect(*loadedN3).To(Equal(*n3))

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestDeleteRelationshipEndpoint(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	simpleRelationship := &SimpleRelationship{}

	n3 := &Node3{}
	n4 := &Node4{}
	n5 := &Node5{}

	n4.N3 = n3

	simpleRelationship.N4 = n4
	simpleRelationship.N5 = n5

	g.Expect(session.Save(&simpleRelationship, nil)).NotTo(HaveOccurred())
	g.Expect(session.Delete(&n4)).NotTo(HaveOccurred())

	//Spec
	g.Expect(n4.DeletedAt).NotTo(BeZero())
	g.Expect(simpleRelationship.DeletedAt).NotTo(BeZero())
	g.Expect(*n4.ID).To(Equal(deletedID))
	g.Expect(*simpleRelationship.ID).To(Equal(deletedID))
	g.Expect(*n5.ID > -1).To(BeTrue())
	g.Expect(*n3.ID > -1).To(BeTrue())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestDeleteRelationship(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	simpleRelationship := &SimpleRelationship{}

	n3 := &Node3{}
	n4 := &Node4{}
	n5 := &Node5{}

	n4.N3 = n3

	simpleRelationship.N4 = n4
	simpleRelationship.N5 = n5

	g.Expect(session.Save(&simpleRelationship, nil)).NotTo(HaveOccurred())
	g.Expect(session.Delete(&simpleRelationship)).NotTo(HaveOccurred())

	//Spec
	g.Expect(n4.DeletedAt).To(BeZero())
	g.Expect(n4.UpdatedAt).ToNot(BeZero())
	g.Expect(n5.DeletedAt).To(BeZero())
	g.Expect(n5.UpdatedAt).ToNot(BeZero())
	g.Expect(simpleRelationship.DeletedAt).NotTo(BeZero())
	g.Expect(*n4.ID > deletedID).To(BeTrue())
	g.Expect(*simpleRelationship.ID).To(Equal(deletedID))
	g.Expect(*n5.ID > -1).To(BeTrue())
	g.Expect(*n3.ID > -1).To(BeTrue())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestDeleteAllNodes(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	n4_1 := &Node4{}
	n4_2 := &Node4{}
	n4_3 := &Node4{}
	n4_4 := &Node4{}
	n4_5 := &Node4{}
	n5 := &Node5{}
	simpleRelationship := &SimpleRelationship{}
	simpleRelationship.N4 = n4_5
	simpleRelationship.N5 = n5
	n4s := [4]*Node4{n4_1, n4_2, n4_3, n4_4}
	g.Expect(session.Save(&n4s, nil)).NotTo(HaveOccurred())
	g.Expect(session.Save(&simpleRelationship, nil)).NotTo(HaveOccurred())

	n4Ref := &Node4{}
	countOfN4, err := session.CountEntitiesOfType(&n4Ref)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(session.DeleteAll(&n4Ref, nil)).NotTo(HaveOccurred())
	postDeleteCountOfN5, err := session.CountEntitiesOfType(&n4Ref)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(countOfN4).To(Equal(int64(5)))
	g.Expect(postDeleteCountOfN5).To(Equal(int64(0)))
	for _, n4 := range n4s {
		g.Expect(*n4.ID).To(Equal(deletedID))
	}
	g.Expect(*simpleRelationship.ID).To(Equal(deletedID))
	g.Expect(simpleRelationship.DeletedAt).NotTo(BeZero(), "Deleting n4_5 should delete this related relationship")
	g.Expect(*n5.ID > deletedID).To(BeTrue())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestDeleteAllRelationships(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	simpleRelationship0 := &SimpleRelationship{}
	simpleRelationship1 := &SimpleRelationship{}
	simpleRelationships := [2]*SimpleRelationship{simpleRelationship0, simpleRelationship1}

	n4_0 := &Node4{}
	n5_0 := &Node5{}
	simpleRelationship0.N4 = n4_0
	simpleRelationship0.N5 = n5_0

	n4_1 := &Node4{}
	n5_1 := &Node5{}
	simpleRelationship1.N4 = n4_1
	simpleRelationship1.N5 = n5_1

	g.Expect(session.Save(&simpleRelationships, nil)).NotTo(HaveOccurred())
	g.Expect(session.Clear()).NotTo(HaveOccurred())

	simpleRelationshipRef := &SimpleRelationship{}
	countOfSimpleRelationships, err := session.CountEntitiesOfType(&simpleRelationshipRef)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(session.DeleteAll(&simpleRelationshipRef, nil)).NotTo(HaveOccurred())
	countOfSimpleRelationshipsPostDelete, err := session.CountEntitiesOfType(&simpleRelationshipRef)

	g.Expect(countOfSimpleRelationships).To(Equal(int64(2)))
	g.Expect(countOfSimpleRelationshipsPostDelete).To(BeZero())

	g.Expect(*n4_0.ID > -1).To(BeTrue())
	g.Expect(*n5_0.ID > -1).To(BeTrue())
	g.Expect(*n4_1.ID > -1).To(BeTrue())
	g.Expect(*n5_1.ID > -1).To(BeTrue())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestSaveMultiSourceRelationship(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	theMatrix := &Movie{}
	theMatrix.Title = "The Matrix"
	theMatrix.Released = 1999
	theMatrix.Tagline = "Welcome to the Real World"

	carrieAnne := &Actor{}
	carrieAnne.Name = "Carrie-Anne Moss"
	carrieAnne.Born = 1967

	keanu := &Actor{}
	keanu.Name = "Keanu Reeves"
	keanu.Born = 1964

	carrieAnneMatrixCharacter := &Character{Movie: theMatrix, Actor: carrieAnne, Roles: []string{"Trinity"}, Name: "carrieAnneMatrixCharacter"}
	keanuReevesMatrixCharacter := &Character{Movie: theMatrix, Actor: keanu, Roles: []string{"Neo"}, Name: "keanuReevesMatrixCharacter"}

	theMatrix.AddCharacter(carrieAnneMatrixCharacter)
	g.Expect(session.Save(&keanuReevesMatrixCharacter, nil)).NotTo(HaveOccurred())

	g.Expect(session.Clear()).NotTo(HaveOccurred())

	lo := gogm.NewLoadOptions()
	lo.Depth = -1
	var loadedTheMatrix *Movie
	g.Expect(session.Load(&loadedTheMatrix, *theMatrix.ID, lo)).NotTo(HaveOccurred())

	g.Expect(len(loadedTheMatrix.Characters)).To(Equal(2))

	g.Expect(*loadedTheMatrix.Characters[0].ID).To(Equal(*carrieAnneMatrixCharacter.ID))
	g.Expect(*loadedTheMatrix.Characters[1].ID).To(Equal(*keanuReevesMatrixCharacter.ID))

	g.Expect(*loadedTheMatrix.Characters[0].Actor.ID).To(Equal(*carrieAnne.ID))
	g.Expect(*loadedTheMatrix.Characters[1].Actor.ID).To(Equal(*keanu.ID))

	g.Expect(len(loadedTheMatrix.Characters[0].Actor.Characters)).To(Equal(1))
	g.Expect(*loadedTheMatrix.Characters[0].Actor.Characters[0].ID).To(Equal(*carrieAnneMatrixCharacter.ID))

	g.Expect(len(loadedTheMatrix.Characters[1].Actor.Characters)).To(Equal(1))
	g.Expect(*loadedTheMatrix.Characters[1].Actor.Characters[0].ID).To(Equal(*keanuReevesMatrixCharacter.ID))

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestLoadingSameNodeTypesNotNavigableInBothDir(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	jamesThompson := &Person{}
	jamesThompson.Name = "James Thompson"
	jamesThompson.Tags = []string{"James", "Followee"}

	jessicaThompson := &Person{}
	jessicaThompson.Name = "Jessica Thompson"
	jessicaThompson.Tags = []string{"Jessica", "Follower"}

	angelaScope := &Person{}
	angelaScope.Name = "Angela Scope"
	angelaScope.Tags = []string{"Angela", "Followee"}

	jessicaThompson.Follows = append(jessicaThompson.Follows, jamesThompson, angelaScope)

	g.Expect(session.Save(&jessicaThompson, nil)).NotTo(HaveOccurred())

	g.Expect(session.Clear()).NotTo(HaveOccurred())
	var loadedJessicaThompson, loadedJamesThompson, loadedAngelaScope *Person
	lo := gogm.NewLoadOptions()
	lo.Depth = -1
	g.Expect(session.Load(&loadedJessicaThompson, *jessicaThompson.ID, lo)).NotTo(HaveOccurred())
	g.Expect(session.Load(&loadedJamesThompson, *jamesThompson.ID, lo)).NotTo(HaveOccurred())
	g.Expect(session.Load(&loadedAngelaScope, *angelaScope.ID, lo)).NotTo(HaveOccurred())

	g.Expect(loadedJessicaThompson).NotTo(BeNil())
	g.Expect(loadedJamesThompson).NotTo(BeNil())
	g.Expect(loadedAngelaScope).NotTo(BeNil())

	g.Expect(len(loadedJessicaThompson.Follows)).To(Equal(2))
	g.Expect(*loadedJessicaThompson.Follows[0].ID).To(Equal(*angelaScope.ID))
	g.Expect(*loadedJessicaThompson.Follows[1].ID).To(Equal(*jamesThompson.ID))

	g.Expect(len(loadedJamesThompson.Follows)).To(BeZero())
	g.Expect(len(loadedAngelaScope.Follows)).To(BeZero())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestLoadingSameNodeTypesNavigableInBothDir(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	jamesThompson := &Person2{}
	jamesThompson.Name = "James Thompson"

	jessicaThompson := &Person2{}
	jessicaThompson.Name = "Jessica Thompson"

	angelaScope := &Person2{}
	angelaScope.Name = "Angela Scope"

	jessicaThompson.Follows = append(jessicaThompson.Follows, jamesThompson, angelaScope)

	g.Expect(session.Save(&jessicaThompson, nil)).NotTo(HaveOccurred())

	g.Expect(session.Clear()).NotTo(HaveOccurred())
	var loadedJessicaThompson, loadedJamesThompson, loadedAngelaScope *Person2
	lo := gogm.NewLoadOptions()
	lo.Depth = -1
	g.Expect(session.Load(&loadedJessicaThompson, *jessicaThompson.ID, lo)).NotTo(HaveOccurred())
	g.Expect(session.Load(&loadedJamesThompson, *jamesThompson.ID, lo)).NotTo(HaveOccurred())
	g.Expect(session.Load(&loadedAngelaScope, *angelaScope.ID, lo)).NotTo(HaveOccurred())

	g.Expect(loadedJessicaThompson).NotTo(BeNil())
	g.Expect(loadedJamesThompson).NotTo(BeNil())
	g.Expect(loadedAngelaScope).NotTo(BeNil())

	g.Expect(len(loadedJessicaThompson.Follows)).To(Equal(2))
	g.Expect(*loadedJessicaThompson.Follows[0].ID).To(Equal(*angelaScope.ID))
	g.Expect(*loadedJessicaThompson.Follows[1].ID).To(Equal(*jamesThompson.ID))

	g.Expect(len(loadedJamesThompson.Follows)).To(Equal(1))
	g.Expect(*loadedJamesThompson.Follows[0].ID).To(Equal(*loadedJessicaThompson.ID))

	g.Expect(len(loadedAngelaScope.Follows)).To(Equal(1))
	g.Expect(*loadedAngelaScope.Follows[0].ID).To(Equal(*loadedJessicaThompson.ID))

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

//Context:when path is saved
//Spec:should be able to log full path
func TestFullPathSaveByOGMIsLoadable(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())

	//(node0)-->(node1)-->(node2)<--(node3)-->(node4)
	n0 := &Node0{}
	n1 := &Node1{}
	n2 := &Node2{}
	n3 := &Node3{}
	n4 := &Node4{}

	n0.Name = "0"
	n1.Name = "1"
	n2.Name = "2"
	n3.Name = "3"
	n4.Name = "4"

	n0.N1 = n1
	n1.N2 = n2
	n2.N3 = n3
	n3.N4 = n4

	var loadedN0 *Node0
	g.Expect(session.Save(&n0, nil)).NotTo(HaveOccurred())
	g.Expect(session.Clear()).NotTo(HaveOccurred())

	lo := gogm.NewLoadOptions()
	lo.Depth = -1
	g.Expect(session.Load(&loadedN0, *n0.ID, lo)).NotTo(HaveOccurred())

	g.Expect(loadedN0).ToNot(BeNil())
	g.Expect(*loadedN0.ID).To(Equal(*n0.ID))

	g.Expect(*loadedN0.N1.ID).To(Equal(*n1.ID))

	g.Expect(*loadedN0.N1.N0.ID).To(Equal(*n0.ID))
	g.Expect(*loadedN0.N1.N2.ID).To(Equal(*n2.ID))

	g.Expect(*loadedN0.N1.N2.N1.ID).To(Equal(*n1.ID))
	g.Expect(*loadedN0.N1.N2.N3.ID).To(Equal(*n3.ID))

	g.Expect(*loadedN0.N1.N2.N3.N2.ID).To(Equal(*n2.ID))
	g.Expect(*loadedN0.N1.N2.N3.N4.ID).To(Equal(*n4.ID))

	g.Expect(*loadedN0.N1.N2.N3.N4.N3.ID).To(Equal(*n3.ID))

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

//Context: when path is saved
//Spec: should be able to load path up to depth x
func TestPathSaveByOGMIsLoadable(t *testing.T) {

	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())

	//(node0)-->(node1)-->(node2)<--(node3)-->(node4)
	n0 := &Node0{}
	n1 := &Node1{}
	n2 := &Node2{}
	n3 := &Node3{}
	n4 := &Node4{}

	n0.Name = "0"
	n1.Name = "1"
	n2.Name = "2"
	n3.Name = "3"
	n4.Name = "4"

	n0.N1 = n1
	n1.N2 = n2
	n2.N3 = n3
	n3.N4 = n4

	var loadedN0 *Node0
	g.Expect(session.Save(&n0, nil)).NotTo(HaveOccurred())
	g.Expect(session.Clear()).NotTo(HaveOccurred())

	lo := gogm.NewLoadOptions()
	lo.Depth = 2
	g.Expect(session.Load(&loadedN0, *n0.ID, lo)).NotTo(HaveOccurred())

	g.Expect(loadedN0).ToNot(BeNil())
	g.Expect(*loadedN0.ID).To(Equal(*n0.ID))

	g.Expect(*loadedN0.N1.ID).To(Equal(*n1.ID))

	g.Expect(*loadedN0.N1.N0.ID).To(Equal(*n0.ID))
	g.Expect(*loadedN0.N1.N2.ID).To(Equal(*n2.ID))

	g.Expect(*loadedN0.N1.N2.N1.ID).To(Equal(*n1.ID))
	g.Expect(loadedN0.N1.N2.N3).To(BeNil())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())

}

func TestLoadFromLocalStore(t *testing.T) {

	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())

	//(node0)-->(node1)-->(node2)<--(node3)-->(node4)
	n0 := &Node0{}
	n1 := &Node1{}
	n2 := &Node2{}
	n3 := &Node3{}
	n4 := &Node4{}

	n0.Name = "0"
	n1.Name = "1"
	n2.Name = "2"
	n3.Name = "3"
	n4.Name = "4"

	n0.N1 = n1
	n1.N2 = n2
	n2.N3 = n3
	n3.N4 = n4

	g.Expect(session.Save(&n0, nil)).NotTo(HaveOccurred())
	g.Expect(session.Clear()).NotTo(HaveOccurred())

	var loadedN1 *Node1
	lo := gogm.NewLoadOptions()
	lo.Depth = 2
	g.Expect(session.Load(&loadedN1, *n1.ID, lo)).NotTo(HaveOccurred())

	g.Expect(loadedN1).ToNot(BeNil())
	g.Expect(*loadedN1.ID).To(Equal(*n1.ID))

	g.Expect(*loadedN1.N0.ID).To(Equal(*n0.ID))
	g.Expect(*loadedN1.N2.ID).To(Equal(*n2.ID))

	g.Expect(*loadedN1.N2.N3.ID).To(Equal(*n3.ID))

	lo.Depth = 1
	var loadedN1_1 *Node1
	g.Expect(session.Load(&loadedN1_1, *n1.ID, lo)).NotTo(HaveOccurred())

	g.Expect(loadedN1_1).To(Equal(loadedN1))

	lo.Depth = 2
	g.Expect(session.Load(&loadedN1_1, *n1.ID, lo)).NotTo(HaveOccurred())
	g.Expect(loadedN1_1).To(Equal(loadedN1))

	lo.Depth = 5
	g.Expect(session.Load(&loadedN1_1, *n1.ID, lo)).NotTo(HaveOccurred())
	g.Expect(loadedN1_1).To(Equal(loadedN1))

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())

}

func TestSaveToDepthFromNode(t *testing.T) {

	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())

	//(node0)-->(node1)-->(node2)<--(node3)-->(node4)
	n0 := &Node0{}
	n1 := &Node1{}
	n2 := &Node2{}
	n3 := &Node3{}
	n4 := &Node4{}

	n0.Name = "0"
	n1.Name = "1"
	n2.Name = "2"
	n3.Name = "3"
	n4.Name = "4"

	n0.N1 = n1
	n1.N2 = n2
	n2.N3 = n3
	n3.N4 = n4

	so := gogm.NewSaveOptions()
	so.Depth = 0
	g.Expect(session.Save(&n0, so)).NotTo(HaveOccurred())

	so.Depth = 2
	g.Expect(session.Save(&n0, so)).NotTo(HaveOccurred())

	//TODO load indefitely and verify save length

	so.Depth = 0
	n0.Name = "31"
	g.Expect(session.Save(&n0, so)).NotTo(HaveOccurred())

	so.Depth = 4
	g.Expect(session.Save(&n0, so)).NotTo(HaveOccurred())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestSaveToDepthFromRelationship(t *testing.T) {
	//TODO add verifications
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())

	//Context
	simpleRelationship := &SimpleRelationship{}
	n4 := &Node4{}
	n5 := &Node5{}
	simpleRelationship.N4 = n4
	simpleRelationship.N5 = n5

	n5.N4 = n4

	so := gogm.NewSaveOptions()
	so.Depth = 0
	g.Expect(session.Save(&simpleRelationship, so)).NotTo(HaveOccurred())

	so.Depth = 2
	g.Expect(session.Save(&simpleRelationship, so)).NotTo(HaveOccurred())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestSaveNodeWithCustomID(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())

	n9 := &Node9{}
	n9.TestId = "r"
	g.Expect(session.Save(&n9, nil)).NotTo(HaveOccurred())

	g.Expect(session.Clear()).NotTo(HaveOccurred())

	var loadedN9 *Node9
	lo := gogm.NewLoadOptions()
	lo.Depth = 2
	g.Expect(session.Load(&loadedN9, n9.TestId, lo)).NotTo(HaveOccurred())
	g.Expect(*loadedN9.ID).To(Equal(*n9.ID))

	var loadedN9_1 *Node9
	lo.Depth = 0
	g.Expect(session.Load(&loadedN9_1, loadedN9.TestId, lo)).NotTo(HaveOccurred())
	g.Expect(loadedN9_1 == loadedN9).To(BeTrue())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestSaveRelationshipWithCustomID(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())

	//Context
	r := &SimpleRelationship{}
	r.TestID = "TestID"
	n4 := &Node4{}
	n5 := &Node5{}
	r.N4 = n4
	r.N5 = n5
	g.Expect(session.Save(&r, nil)).NotTo(HaveOccurred())

	g.Expect(session.Clear()).NotTo(HaveOccurred())

	var loadedR *SimpleRelationship
	lo := gogm.NewLoadOptions()
	lo.Depth = 2
	g.Expect(session.Load(&loadedR, r.TestID, lo)).NotTo(HaveOccurred())
	g.Expect(*loadedR.ID).To(Equal(*r.ID))

	var loadedR_1 *SimpleRelationship
	lo.Depth = 0
	g.Expect(session.Load(&loadedR_1, loadedR.TestID, lo)).NotTo(HaveOccurred())
	g.Expect(loadedR_1 == loadedR).To(BeTrue())

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestTransactions(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())

	//(node0)-->(node1)-->(node2)<--(node3)-->(node4)
	n0 := &Node0{}
	n1 := &Node1{}
	n2 := &Node2{}
	n3 := &Node3{}
	n4 := &Node4{}

	n0.Name = "0"
	n1.Name = "1"
	n2.Name = "2"
	n3.Name = "3"
	n4.Name = "4"

	n0.N1 = n1
	n1.N2 = n2
	n2.N3 = n3
	n3.N4 = n4

	so := gogm.NewSaveOptions()
	so.Depth = 2
	g.Expect(session.Save(&n0, so)).NotTo(HaveOccurred())

	tx, err := session.BeginTransaction()
	g.Expect(err).NotTo(HaveOccurred())

	n0.Name = "0Update"
	n2.Name = "2Update"

	//Test Rolling back
	so.Depth = 3
	g.Expect(session.Save(&n0, so)).NotTo(HaveOccurred())
	g.Expect(tx.RollBack()).NotTo(HaveOccurred())
	g.Expect(tx.Close()).NotTo(HaveOccurred())

	var loadedN0 *Node0
	lo := gogm.NewLoadOptions()
	lo.Depth = 0
	g.Expect(session.Load(&loadedN0, *n0.ID, lo)).NotTo(HaveOccurred())

	g.Expect(n0.Name).To(Equal("0Update"))
	g.Expect(loadedN0.Name).To(Equal(n0.Name))
	g.Expect(*loadedN0.N1.N2.N3.ID).To(Equal(*n3.ID), "Store cache still holds state of rolled back transcation")

	g.Expect(session.Reload(&n0)).NotTo(HaveOccurred(), "Reload to sycn runtime objects with backend")
	g.Expect(n0.Name).To(Equal("0"))
	g.Expect(*n3.ID).To(Equal(deletedID), "n3 gets deleted as it was rolled back. n3 can't ever be saved again. Must create new instance to save")

	//Testing committing
	n3 = &Node3{}
	n2.N3 = n3
	n3.Name = "3"
	n3.N4 = n4
	so.Depth = 3

	tx, err = session.BeginTransaction()
	g.Expect(err).NotTo(HaveOccurred())
	n0.Name = "0Update"
	n2.Name = "2Update"
	g.Expect(session.Save(&n0, so)).NotTo(HaveOccurred())
	g.Expect(tx.Commit()).NotTo(HaveOccurred())
	g.Expect(tx.Close()).NotTo(HaveOccurred())

	g.Expect(session.Reload(&n0)).NotTo(HaveOccurred())
	g.Expect(n0.Name).To(Equal("0Update"))
	g.Expect(n2.Name).To(Equal("2Update"))
	g.Expect(*n3.ID).NotTo(Equal(deletedID))

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

// func TestQueryForObject_s(t *testing.T) {
// 	g := NewGomegaWithT(t)
// 	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
// 	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())

// 	var person *Person
// 	g.Expect(session.QueryForObject(&person, "MATCH (person) RETURN person", nil)).NotTo(HaveOccurred())

// 	g.Expect(person).To(BeNil())

// 	jamesThompson := &Person{}
// 	jamesThompson.Name = "James Thompson"
// 	jamesThompson.Tags = []string{"James", "Followee"}

// 	jessicaThompson := &Person{}
// 	jessicaThompson.Name = "Jessica Thompson"
// 	jessicaThompson.Tags = []string{"Jessica", "Follower"}

// 	angelaScope := &Person{}
// 	angelaScope.Name = "Angela Scope"
// 	angelaScope.Tags = []string{"Angela", "Followee"}

// 	jessicaThompson.Follows = append(jessicaThompson.Follows, jamesThompson, angelaScope)

// 	g.Expect(session.Save(&jessicaThompson, nil)).NotTo(HaveOccurred())

// 	g.Expect(session.QueryForObject(&person, "MATCH (person:PERSON) RETURN person", nil)).To(HaveOccurred())
// 	g.Expect(person).To(BeNil())

// 	g.Expect(session.QueryForObject(&person, "MATCH (person:PERSON) WHERE person.name = $name RETURN person", map[string]interface{}{"name": "Angela Scope"})).ToNot(HaveOccurred())
// 	g.Expect(person).To(Equal(angelaScope))

// 	var persons []*Person
// 	g.Expect(session.QueryForObjects(&persons, "MATCH (person:PERSON) RETURN person", nil)).ToNot(HaveOccurred())
// 	g.Expect(len(persons)).To(Equal(3))

// 	sort.SliceStable(persons, func(i, j int) bool { return persons[i].Name < persons[j].Name })
// 	g.Expect(persons[0]).To(Equal(angelaScope))
// 	g.Expect(persons[1]).To(Equal(jamesThompson))

// 	//Note, just for comparison
// 	jessicaThompson.Follows = nil
// 	g.Expect(persons[2]).To(Equal(jessicaThompson))

// 	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
// 	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
// }

// func TestCount(t *testing.T) {
// 	g := NewGomegaWithT(t)
// 	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
// 	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())

// 	count, err := session.Count("MATCH (n:INVALID) RETURN COUNT(n)", nil)
// 	g.Expect(err).NotTo(HaveOccurred())
// 	g.Expect(count).To(BeZero())

// 	simpleRelationship1 := SimpleRelationship{}
// 	n5 := &Node5{}
// 	n4 := &Node4{}
// 	simpleRelationship1.N4 = n4
// 	simpleRelationship1.N5 = n5

// 	simpleRelationship2 := SimpleRelationship{}
// 	simpleRelationship2.N4 = n4
// 	simpleRelationship2.N5 = n5

// 	simpleRelationships := []*SimpleRelationship{&simpleRelationship1, &simpleRelationship2}

// 	g.Expect(session.Save(&simpleRelationships, nil)).NotTo(HaveOccurred())

// 	count, err = session.Count("MATCH (n:NODE4) RETURN COUNT(n)", nil)
// 	g.Expect(err).NotTo(HaveOccurred())
// 	g.Expect(count).To(Equal(int64(1)))

// 	count, err = session.Count("MATCH (n:NODE5) RETURN COUNT(n)", nil)
// 	g.Expect(err).NotTo(HaveOccurred())
// 	g.Expect(count).To(Equal(int64(1)))

// 	n4_1 := &Node4{}
// 	n4_2 := &Node4{}
// 	n4_3 := &Node4{}
// 	n4_4 := &Node4{}
// 	n4s := [4]*Node4{n4_1, n4_2, n4_3, n4_4}
// 	g.Expect(session.Save(&n4s, nil)).NotTo(HaveOccurred())

// 	count, err = session.Count("MATCH (n:NODE4) RETURN COUNT(n)", nil)
// 	g.Expect(err).NotTo(HaveOccurred())
// 	g.Expect(count).To(Equal(int64(5)))

// 	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
// 	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
// }

func TestMappedProperties(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())

	n0 := &Node0{}
	n0.MapProps = map[string]int{"hello": 8}
	n0.InvalidIDMapProp = map[int]int{2: 8}
	n0.MapToInterfaces = map[string]interface{}{"0": true, "1": 2.3, "2": "hello"}

	e := "hello"
	y := ",world"
	n0.AliasedMapProps1 = map[string]*string{"hello": &e, "world": &y}

	g.Expect(session.Save(&n0, nil)).NotTo(HaveOccurred())
	g.Expect(session.Save(&n0, nil)).NotTo(HaveOccurred())

	g.Expect(session.Clear()).NotTo(HaveOccurred())

	var loadedN0 *Node0
	g.Expect(session.Load(&loadedN0, *n0.ID, nil)).NotTo(HaveOccurred())

	g.Expect(n0.UpdatedAt).To(BeZero())
	g.Expect(loadedN0.MapProps).To(Equal(n0.MapProps))
	g.Expect(loadedN0.InvalidIDMapProp).To(BeNil())
	g.Expect(loadedN0.MapToInterfaces).To(Equal(n0.MapToInterfaces))

	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

func TestSavingAndLoadingTime(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

	n10 := &Node10{}
	location, _ := time.LoadLocation("")
	tValue := time.Now().In(location)
	durationValue := neo4j.DurationOf(2, 3, 6, 7)

	n10.Time = tValue
	n10.Time1 = &tValue
	n10.Duration = durationValue
	n10.Duration1 = &durationValue
	g.Expect(session.Save(&n10, nil)).NotTo(HaveOccurred())
	g.Expect(session.Save(&n10, nil)).NotTo(HaveOccurred())
	g.Expect(n10.UpdatedAt).To(BeZero())

	n10.Duration = neo4j.DurationOf(4, 3, 6, 7)
	g.Expect(session.Save(&n10, nil)).NotTo(HaveOccurred())
	g.Expect(n10.UpdatedAt).NotTo(BeZero())

	var loadedN10 *Node10
	g.Expect(session.Load(&loadedN10, *n10.ID, nil)).NotTo(HaveOccurred())
	g.Expect(loadedN10).To(Equal(n10))
	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
}

//Error cases
// func TestSaveNodeWithInvalidCustomID(t *testing.T) {
// 	g := NewGomegaWithT(t)
// 	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
// 	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())

// 	invalidID := &InvalidID{}
// 	testID := "r"
// 	invalidID.TestId = &testID
// 	g.Expect(session.Save(invalidID, nil)).To(HaveOccurred())

// 	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
// 	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
// }

// func TestForbiddenLabel(t *testing.T) {
// 	g := NewGomegaWithT(t)
// 	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
// 	g.Expect(session.RegisterEventListener(eventListener)).NotTo(HaveOccurred())

// 	angelaScope := &Person{}
// 	angelaScope.Name = "Angela Scope"
// 	angelaScope.Tags = []string{"Angela"}

// 	g.Expect(session.Save(&angelaScope, nil)).ToNot(HaveOccurred())
// 	g.Expect(angelaScope.UpdatedAt).To(BeZero())

// 	angelaScope.Tags = []string{"Ana"}
// 	g.Expect(session.Save(&angelaScope, nil)).ToNot(HaveOccurred())
// 	g.Expect(angelaScope.UpdatedAt).NotTo(BeZero())

// 	angelaScope.Tags = []string{"Ana", "PERSON"}
// 	g.Expect(session.Save(&angelaScope, nil)).To(HaveOccurred())

// 	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
// 	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
// }

// func TestQueryForObject_fail(t *testing.T) {
// 	g := NewGomegaWithT(t)
// 	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
// 	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())

// 	jamesThompson := &Person{}
// 	jamesThompson.Name = "James Thompson"
// 	jamesThompson.Tags = []string{"James", "Followee"}

// 	jessicaThompson := &Person{}
// 	jessicaThompson.Name = "Jessica Thompson"
// 	jessicaThompson.Tags = []string{"Jessica", "Follower"}

// 	angelaScope := &Person{}
// 	angelaScope.Name = "Angela Scope"
// 	angelaScope.Tags = []string{"Angela", "Followee"}

// 	jessicaThompson.Follows = append(jessicaThompson.Follows, jamesThompson, angelaScope)

// 	g.Expect(session.Save(&jessicaThompson, nil)).NotTo(HaveOccurred())

// 	var relationships []*SimpleRelationship
// 	g.Expect(session.QueryForObjects(&relationships, "MATCH (person:PERSON) RETURN person", nil)).To(HaveOccurred())

// 	g.Expect(session.PurgeDatabase()).NotTo(HaveOccurred())
// 	g.Expect(session.DisposeEventListener(eventListener)).NotTo(HaveOccurred())
// }

//Test that loading with ID nil loads all the type
