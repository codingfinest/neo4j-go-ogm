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
