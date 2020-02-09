package gogm

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type transaction struct {
	neo4jTransaction neo4j.Transaction
	session          neo4j.Session
	close            transactionEnder
}

func newTransaction(driver neo4j.Driver, transactionEnder transactionEnder, accessMode neo4j.AccessMode) (*transaction, error) {

	var (
		err     error
		session neo4j.Session
	)

	if session, err = driver.Session(accessMode); err != nil {
		return nil, err
	}

	var neo4jtransaction neo4j.Transaction
	if neo4jtransaction, err = session.BeginTransaction(); err != nil {
		return nil, err
	}

	return &transaction{
		neo4jTransaction: neo4jtransaction,
		session:          session,
		close:            transactionEnder}, nil
}

func (t *transaction) run(cql string, params map[string]interface{}) (neo4j.Result, error) {
	return t.neo4jTransaction.Run(cql, params)
}

func (t *transaction) Commit() error {
	return t.neo4jTransaction.Commit()
}

func (t *transaction) RollBack() error {
	return t.neo4jTransaction.Rollback()
}

func (t *transaction) Close() error {
	return t.close()
}
