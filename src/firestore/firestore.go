package firestore

import (
	"context"
	"errors"
	"math"
	"time"

	firestore "cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/testutil"
)

const (
	TIMEOUT_CONNECT = 5
	TIMEOUT_FIND    = 10
)

type FirestoreDatabaseFunctions struct {
	Config goke.Config
	db     *firestore.Client
}

func (self FirestoreDatabaseFunctions) projectId() string {
	return self.Config.Firestore.ProjectID
}

func (self *FirestoreDatabaseFunctions) FindMany(ctx context.Context, p goke.FindManyParams, hook goke.TransformDocument) ([]goke.Map, error) {
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
	var nodes []goke.Map
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var node goke.Map
		if err := doc.DataTo(&node); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	nodes, err = goke.FilterDocuments(nodes, hook)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (self *FirestoreDatabaseFunctions) InsertMany(ctx context.Context, p goke.InsertManyParams, hook goke.TransformDocument) (goke.NodesMutationPayload, error) {
	db, err := self.Init(ctx)
	payload := goke.NodesMutationPayload{}
	if err != nil {
		return payload, err
	}
	nodes, err := goke.FilterDocuments(p.Data, hook)
	if err != nil {
		return payload, err
	}
	for _, x := range nodes {
		_, _, err := db.Collection(p.Collection).Add(ctx, x)
		if err != nil {
			return payload, err
		}
		// TODO if firestore uses some id i should add it to the returned nodes
		payload.AffectedCount++
		payload.Returning = append(payload.Returning, x)

	}
	return payload, nil
}

func (self *FirestoreDatabaseFunctions) UpdateMany(ctx context.Context, p goke.UpdateParams, hook goke.TransformDocument) (goke.NodesMutationPayload, error) {
	if p.Limit == 0 {
		p.Limit = math.MaxInt16
	}
	db, err := self.Init(ctx)
	payload := goke.NodesMutationPayload{}
	if err != nil {
		return payload, err
	}
	var query firestore.Query = db.Collection(p.Collection).Query
	query, err = applyWhereQuery(p.Where, query)
	if err != nil {
		return payload, err
	}

	iter := query.Documents(ctx)
	defer iter.Stop()

	for payload.AffectedCount < p.Limit {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return payload, err
		}

		// check with hook
		var m goke.Map
		err = doc.DataTo(&m)
		if err != nil {
			return payload, err
		}
		m, err = hook(m)
		if err != nil {
			return payload, err
		}
		if m == nil {
			continue
		}

		// update the doc
		_, err = doc.Ref.Update(ctx, setToUpdates(p.Set))
		if err != nil {
			return payload, err
		}
		payload.AffectedCount++

		// refetch doc with updates
		doc, err = doc.Ref.Get(ctx)
		if err != nil {
			return payload, err
		}
		var node goke.Map
		if err := doc.DataTo(&node); err != nil {
			return payload, err
		}
		payload.Returning = append(payload.Returning, node)
	}
	return payload, nil
}

func (self *FirestoreDatabaseFunctions) DeleteMany(ctx context.Context, p goke.DeleteManyParams, hook goke.TransformDocument) (goke.NodesMutationPayload, error) {
	if p.Limit == 0 {
		p.Limit = math.MaxInt16
	}
	db, err := self.Init(ctx)
	payload := goke.NodesMutationPayload{}
	if err != nil {
		return payload, err
	}
	var query firestore.Query = db.Collection(p.Collection).Query

	query, err = applyWhereQuery(p.Where, query)
	if err != nil {
		return payload, err
	}

	iter := query.Documents(ctx)
	defer iter.Stop()
	// TODO use batching for firestore's updateMany
	for payload.AffectedCount < p.Limit {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return payload, err
		}

		// check with hook
		var m goke.Map
		err = doc.DataTo(&m)
		if err != nil {
			return payload, err
		}
		m, err = hook(m)
		if err != nil {
			return payload, err
		}
		if m == nil {
			continue
		}
		// delete
		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			return payload, err
		}
		payload.AffectedCount++
		payload.Returning = append(payload.Returning, m)
	}
	return payload, nil
}

func setToUpdates(set goke.Map) []firestore.Update {
	var updates []firestore.Update
	for k, v := range set {
		// fieldPath := []string{k}
		updates = append(updates, firestore.Update{Path: k, Value: v})
	}
	return updates
}

func applyWhereQuery(where goke.WhereTree, q firestore.Query) (firestore.Query, error) {
	println(testutil.Pretty("where", where))
	if len(where.Or) != 0 {
		return q, errors.New("cannot use or operator with firestore")
	}
	for k, v := range where.Match {
		if !isZero(v.Neq) {
			return q, errors.New("firestore cannot use the `Neq` operator")
		}
		if !isZero(v.Nin) {
			return q, errors.New("firestore cannot use the `Nin` operator")
		}
		if !isZero(v.Eq) {
			q = q.Where(k, "==", v.Eq)
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
	}
	for _, a := range where.And {
		q, err := applyWhereQuery(a, q)
		if err != nil {
			return q, err
		}
	}
	return q, nil
}

func applyOrderByQuery(orderBy map[string]int, q firestore.Query) firestore.Query {
	for k, v := range orderBy {
		if v == goke.ASC {
			q = q.OrderBy(k, firestore.Asc)
		}
		if v == goke.DESC {
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
