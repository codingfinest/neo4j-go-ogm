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

type TestObject interface {
	SetCreatedTime(int64)
	SetLoadedTime(int64)
	SetUpdatedTime(int64)
	SetDeletedTime(int64)
	ClearMetaTimestamps()
}

type TestEntity struct {
	CreatedAt int64
	LoadedAt  int64
	UpdatedAt int64
	DeletedAt int64
}

func (n *TestEntity) SetCreatedTime(time int64) {
	n.CreatedAt = time
}
func (n *TestEntity) SetLoadedTime(time int64) {
	n.LoadedAt = time
}
func (n *TestEntity) SetUpdatedTime(time int64) {
	n.UpdatedAt = time
}
func (n *TestEntity) SetDeletedTime(time int64) {
	n.DeletedAt = time
}
func (n *TestEntity) ClearMetaTimestamps() {
	n.CreatedAt = 0
	n.LoadedAt = 0
	n.UpdatedAt = 0
	n.DeletedAt = 0
}
