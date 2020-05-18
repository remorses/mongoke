package firestore

import (
	"context"
	"errors"
	"reflect"
	"time"

	firestore "cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	mongoke "github.com/remorses/mongoke/src"
)

const (
	TIMEOUT_CONNECT = 5
	TIMEOUT_FIND    = 10
)

type FirestoreDatabaseFunctions struct {
	mongoke.DatabaseInterface
	db *firestore.Client
}

func (self FirestoreDatabaseFunctions) FindMany(ctx context.Context, p mongoke.FindManyParams) ([]mongoke.Map, error) {
	db, err := self.Init(ctx, p.DatabaseUri)
	if err != nil {
		return nil, err
	}
	ctx, _ = context.WithTimeout(ctx, TIMEOUT_FIND*time.Second)
	var query firestore.Query = db.Collection(p.Collection).Query
	if p.Limit != 0 {
		query = query.Limit(p.Limit)
	}
	if p.Offset != 0 {
		query = query.Offset(p.Offset)
	}
	query, err = applyWhereQuery(p.Where, query)
	if err != nil {
		return nil, err
	}
	query = applyOrderByQuery(p.OrderBy, query)
	iter := query.Documents(ctx)
	defer iter.Stop()
	var nodes []mongoke.Map
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var node mongoke.Map
		if err := doc.DataTo(&node); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func applyWhereQuery(where map[string]mongoke.Filter, q firestore.Query) (firestore.Query, error) {
	for k, v := range where {
		if !isZero(v.Eq) {
			q = q.Where(k, "==", v.Eq)
		}
		if !isZero(v.Neq) {
			return q, errors.New("firestore cannot use the `Neq` operator")
		}
		if !isZero(v.Gt) {
			q = q.Where(k, ">", v.Gt)
		}
		if !isZero(v.Gte) {
			q = q.Where(k, ">=", v.Gte)
		}
		if !isZero(v.Lt) {
			q = q.Where(k, "<", v.Lt)
		}
		if !isZero(v.Lte) {
			q = q.Where(k, "<=", v.Lte)
		}
		if !isZero(v.In) {
			q = q.Where(k, "in", v.In)
		}
		if !isZero(v.Nin) {
			return q, errors.New("firestore cannot use the `Nin` operator")
		}
	}
	return q, nil
}

func applyOrderByQuery(orderBy map[string]int, q firestore.Query) firestore.Query {
	for k, v := range orderBy {
		if v == mongoke.ASC {
			q = q.OrderBy(k, firestore.Asc)
		}
		if v == mongoke.DESC {
			q = q.OrderBy(k, firestore.Desc)
		}
	}
	return q
}

func isZero(v interface{}) bool {
	if v == nil || v == false || v == 0 || v == "" {
		return true
	}
	t := reflect.TypeOf(v)
	if !t.Comparable() { // TODO what types are not comparable?
		return true
	}
	return v == reflect.Zero(t).Interface()
}

func (self *FirestoreDatabaseFunctions) Init(ctx context.Context, projectID string) (*firestore.Client, error) {
	if self.db != nil {
		return self.db, nil
	}
	ctx, _ = context.WithTimeout(ctx, TIMEOUT_CONNECT*time.Second)
	// option.WithCredentialsJSON()
	db, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// defer client.Close()
	self.db = db
	return db, nil
}
