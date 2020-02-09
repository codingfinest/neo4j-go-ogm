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
