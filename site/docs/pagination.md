---
route: /docs/pagination
title: pagination
---

Mongoke generates queries that can be paginated via the fileds `first` and `after` or `last` and `before` .

By default Mongoke sorts the documents in descending order, this means that using the query `Users(first: 10)` fill fetch the 10 users most recently created.

### Paging forward

You page forward using forst and after arguments.

After executing this query with `{first: 10}` you can get the next page with `{first: 10, after: pageInfo.endCursor}` .

``` 
query Users($first: Int, $after: AnyScalar) {
  Users(first: $first, after: $after) {
    nodes {
      name
      _id
    }
    pageInfo {
      endCursor
      hasNextPage
    }
  }
}
```

### Paging backwards

You page backwards using last and before arguments.

After executing this query with `{last: 10}` you can get the next page with `{last: 10, before: pageInfo.startCursor}` .

To use this query in you frontend remember that you will have to stack the fetched items one over the over, as if you are paging backwards.

``` 
query Users($last: Int, $before: AnyScalar) {
  Users(first: $last, after: $before) {
    nodes {
      name
      _id
    }
    pageInfo {
      endCursor
      hasNextPage
    }
  }
}
```

## Change paging direction

By default Mongoke sorts the documents in descending order, this means that using the query `Users(first: 10)` fill fetch the 10 users most recently created.

If you want the opposite behviour you can set the `direction` argument to `ASC` 

The following query will fetch the 10 first created users.

``` 
query Users($first: Int, $after: AnyScalar) {
  Users(first: $first, after: $after) {
    nodes {
      name
      _id
    }
    pageInfo {
      endCursor
      hasNextPage
    }
  }
}
```

## Change the ordering field ( `cursorFiled` )

By default Mongoke uses the `_id` documents field for ordering and paginating results, this means that the returned `pageInfo.endCursor` and `pageInfo.startCursor` are the simply the last and first `_id` document field.

You can change this behaviour using the `cursorField` argument, this can be any of the document scalar field.

Remember that to make queries efficent you should add MongoDb indexes to the `cusrsorFiled` you use often.

!!! note

    The `cursorFiled` argument is not a string but an enum, you don't need to add the `"` string quotes around the argument.

``` 
query Users($first: Int, $after: AnyScalar, cursorField: name) {
  Users(first: $first, after: $after) {
    nodes {
      name
      _id
    }
    pageInfo {
      endCursor # this is simply the last document name field
      hasNextPage
    }
  }
}
```

