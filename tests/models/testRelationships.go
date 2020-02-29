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

package models

import gogm "github.com/codingfinest/neo4j-go-ogm"

type TestRelationshipEntity struct {
	gogm.Relationship
	TestEntity
}

//Relationships

type SimpleRelationship struct {
	TestRelationshipEntity
	N5     *Node5 `gogm:"startNode"`
	N4     *Node4 `gogm:"endNode"`
	Name   string
	TestID string `gogm:"id"`
}

type SimpleRelationship2 struct {
	TestRelationshipEntity
	N5   *Node5 `gogm:"startNode"`
	N4   *Node4 `gogm:"endNode"`
	Name string
}
