package firestore

import (
	"context"
	"errors"
	"math"
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
	Config mongoke.Config
	mongoke.DatabaseInterface
	db *firestore.Client
}

func (self FirestoreDatabaseFunctions) projectId() string {
	return self.Config.Firestore.ProjectID
}

func (self *FirestoreDatabaseFunctions) FindMany(ctx context.Context, p mongoke.FindManyParams) ([]mongoke.Map, error) {
	db, err := self.Init(ctx)
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

func (self *FirestoreDatabaseFunctions) UpdateOne(ctx context.Context, p mongoke.UpdateOneParams) (mongoke.NodeMutationPayload, error) {
	res, err := self.updateMany(ctx, p, 1)
	payload := mongoke.NodeMutationPayload{}
	if err != nil {
		return payload, err
	}
	if len(res.Returning) == 0 {
		return payload, nil
	}
	payload.Returning = res.Returning[0]
	payload.AffectedCount = 1
	return payload, nil
}

func (self *FirestoreDatabaseFunctions) UpdateMany(ctx context.Context, p mongoke.UpdateOneParams) (mongoke.NodesMutationPayload, error) {
	return self.updateMany(ctx, p, math.MaxInt32)
}

func (self *FirestoreDatabaseFunctions) updateMany(ctx context.Context, p mongoke.UpdateOneParams, count int) (mongoke.NodesMutationPayload, error) {
	db, err := self.Init(ctx)
	if err != nil {
		return mongoke.NodesMutationPayload{}, err
	}
	var query firestore.Query = db.Collection(p.Collection).Query

	query, err = applyWhereQuery(p.Where, query)
	if err != nil {
		return mongoke.NodesMutationPayload{}, err
	}

	iter := query.Documents(ctx)
	defer iter.Stop()
	payload := mongoke.NodesMutationPayload{}
	for payload.AffectedCount < count {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return payload, err
		}
		_, err = doc.Ref.Update(ctx, setToUpdates(p.Set))
		if err != nil {
			return payload, err
		}
		payload.AffectedCount++
		doc, err = doc.Ref.Get(ctx)
		if err != nil {
			return payload, err
		}
		var node mongoke.Map
		if err := doc.DataTo(&node); err != nil {
			return payload, err
		}
		payload.Returning = append(payload.Returning, node)
	}
	return payload, nil
}

func setToUpdates(set mongoke.Map) []firestore.Update {
	var updates []firestore.Update
	for k, v := range set {
		// fieldPath := []string{k}
		updates = append(updates, firestore.Update{Path: k, Value: v})
	}
	return updates
}

func applyWhereQuery(where map[string]mongoke.Filter, q firestore.Query) (firestore.Query, error) {
	// println(testutil.Pretty("where", where))
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
	if v == nil {
		return true
	}
	if l, ok := v.([]interface{}); ok { // this way the in and nin operators cannot be given []
		return len(l) == 0
	}
	return false
}

func (self *FirestoreDatabaseFunctions) InsertMany(ctx context.Context, p mongoke.InsertManyParams) ([]mongoke.Map, error) {
	db, err := self.Init(ctx)
	if err != nil {
		return nil, err
	}
	for _, x := range p.Data {
		_, _, err := db.Collection(p.Collection).Add(ctx, x)
		if err != nil {
			return nil, err
		}
		// TODO if firestore uses some id i should add it to the returned nodes
	}
	return p.Data, nil
}

func (self *FirestoreDatabaseFunctions) Init(ctx context.Context) (*firestore.Client, error) {
	if self.db != nil {
		return self.db, nil
	}
	ctx, _ = context.WithTimeout(ctx, TIMEOUT_CONNECT*time.Second)
	// option.WithCredentialsJSON()
	uri := self.projectId()
	if uri == "" {
		return nil, errors.New("firestore projectId is missing")
	}
	db, err := firestore.NewClient(ctx, uri)
	if err != nil {
		return nil, err
	}

	// defer client.Close()
	self.db = db
	return db, nil
}
