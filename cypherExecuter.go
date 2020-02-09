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
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type transactionExecuter func(work neo4j.TransactionWork, configurers ...func(*neo4j.TransactionConfig)) (interface{}, error)

type cypherExecuter struct {
	driver      neo4j.Driver
	accessMode  neo4j.AccessMode
	transaction *transaction
}

func newCypherExecuter(driver neo4j.Driver, accessMode neo4j.AccessMode, t *transaction) *cypherExecuter {
	return &cypherExecuter{driver, accessMode, nil}
}

func (c *cypherExecuter) execTransaction(te transactionExecuter, cql string, params map[string]interface{}) (neo4j.Result, error) {
	var (
		err    error
		result neo4j.Result
	)

	if _, err = te(func(tx neo4j.Transaction) (interface{}, error) {
		if result, err = tx.Run(cql, params); err != nil {
			return nil, err
		}
		return result, nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *cypherExecuter) exec(cql string, params map[string]interface{}) (neo4j.Result, error) {
	var (
		result  neo4j.Result
		session neo4j.Session
		err     error
	)
	if c.transaction != nil {
		if result, err = c.transaction.run(cql, params); err != nil {
			return nil, err
		}
		return result, nil
	}

	if session, err = c.driver.Session(c.accessMode); err != nil {
		return nil, err
	}
	defer session.Close()

	transactionMode := session.ReadTransaction
	if c.accessMode == neo4j.AccessModeWrite {
		transactionMode = session.WriteTransaction
	}

	return c.execTransaction(transactionMode, cql, params)
}

func (c *cypherExecuter) setTransaction(transaction *transaction) {
	c.transaction = transaction
}
