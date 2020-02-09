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

package gogm

import (
	"errors"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type transactionEnder func() error

type transactioner struct {
	transaction *transaction
	accessMode  neo4j.AccessMode
}

func newTransactioner(accessMode neo4j.AccessMode) *transactioner {
	return &transactioner{accessMode: accessMode}
}

func (t *transactioner) beginTransaction(s *sessionImpl) (*transaction, error) {
	if t.transaction != nil {
		return nil, errors.New("Transaction already exists")
	}

	var err error
	if t.transaction, err = newTransaction(s.driver, t.endTransaction(s), t.accessMode); err != nil {
		return nil, err
	}

	s.cypherExecuter.setTransaction(t.transaction)

	return t.transaction, nil
}

func (t *transactioner) endTransaction(s *sessionImpl) func() error {

	return func() error {
		var err error
		if t.transaction == nil {
			return errors.New("No transaction exist")
		}
		if err = t.transaction.neo4jTransaction.Close(); err != nil {
			return err
		}

		if err = t.transaction.session.Close(); err != nil {
			return err
		}

		t.transaction = nil
		s.cypherExecuter.setTransaction(nil)

		return nil
	}
}
