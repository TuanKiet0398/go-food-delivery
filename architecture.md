# Architecture

Simple layered ("clean-ish") architecture per module. Each feature module is split into
4 layers, each depending only on the layer below it through an interface:

```
transport (gin handler)  →  biz (use case)  →  storage (gorm)  →  database
                                   ↓
                                model (DTO / entity)
```

- **transport**: parses HTTP request, calls a biz use case, writes the HTTP response.
  Knows about `gin.Context`, nothing about SQL.
- **biz**: one struct per use case, holding a narrow store interface (only the methods
  it needs). Validates input and enforces business rules.
- **storage**: implements the store interfaces biz declares, using GORM. Knows about SQL.
- **model**: plain structs shared by all layers (DB row shape, request/filter shape).

Dependency direction is inverted at the biz layer: `biz` defines the interface it needs
(e.g. `CreateRestaurantStore`), and `storage.sqlStore` happens to satisfy it — biz never
imports the storage package's concrete type.

## Directory layout

```
main.go                                  # wiring: DB connection, routes
component/appctx/app_context.go          # AppContext: carries *gorm.DB (and future deps)
common/
  sql_model.go                           # SQLModel: ID, Status, CreatedAt, UpdatedAt (embedded in every entity)
  paging.go                              # Paging: page/limit/total + Fulfill() defaults
  app_response.go                        # successRes envelope: {data, paging, filter}
module/restaurant/
  model/
    restaurant.go                        # Restaurant, RestaurantCreate, RestaurantUpdate
    filter.go                            # Filter (query filters, e.g. owner_id)
  biz/
    create_restaurant.go                 # createRestaurantBiz.CreateRestaurant
    delete_restaurant.go                 # deleteRestaurantBiz.DeleteRestaurant (soft delete + guard)
    list_restaurant.go                   # listRestaurantBiz.ListRestaurant
  storage/
    store.go                             # sqlStore{db *gorm.DB}, NewSQLStore
    create.go                            # sqlStore.Create
    delete.go                            # sqlStore.Delete (soft delete via status=0)
    find.go                              # sqlStore.FindDataWithCondition
    list.go                              # sqlStore.ListDataWithCondition (filter + pagination)
  transport/ginrestaurant/
    create_restaurant.go                 # POST   /v1/restaurants
    delete_restaurant.go                 # DELETE /v1/restaurants/:id
    list_restaurant.go                   # GET    /v1/restaurants
```

`GET /v1/restaurants/:id` and `PATCH /v1/restaurants/:id` are still inline closures in
`main.go` and have not been migrated into this layered structure yet.

## Request flow (example: `GET /v1/restaurants` — list with pagination/filter)

1. **main.go** registers the route: `restaurants.GET("", ginrestaurant.ListRestaurant(appContext))`.
2. **transport/ginrestaurant/list_restaurant.go**
   - Gets the `*gorm.DB` from `appCtx.GetMainDBConnection()`.
   - Binds query params into `common.Paging` and `restaurantmodel.Filter` via `c.ShouldBind`.
   - Calls `pagingData.Fulfill()` to apply defaults (page=1, limit=50).
   - Constructs `restaurantstorage.NewSQLStore(db)` and `restaurantbiz.NewListRestaurantBiz(store)`.
   - Calls `biz.ListRestaurant(ctx, &filter, &pagingData)`.
   - Wraps the result in `common.NewSuccessResponse(result, pagingData, filter)`.
3. **biz/list_restaurant.go**
   - `listRestaurantBiz.ListRestaurant` just forwards to `store.ListDataWithCondition`
     (no extra business rules for listing yet).
4. **storage/list.go**
   - Builds a reusable `query()` closure: base condition `status = 1`, plus
     `owner_id = ?` if `filter.OwnerId > 0`.
   - Runs `query().Count(&paging.Total)` first to get the total row count.
   - Calls `paging.Fulfill()` again (defensive) and computes `offset`.
   - Runs `query().Offset(offset).Limit(limit).Order("id desc").Find(&result)`.
   - `query()` is re-invoked (not reused/chained) each time to avoid GORM's condition
     accumulation across chained calls.
5. Response JSON: `{"data": [...restaurants], "paging": {...}, "filter": {...}}`.

## Request flow (example: `POST /v1/restaurants` — create)

1. `transport/ginrestaurant/create_restaurant.go` binds JSON body into
   `restaurantmodel.RestaurantCreate`.
2. `biz/create_restaurant.go`: `CreateRestaurant` validates `Name != ""`, then calls
   `store.Create`.
3. `storage/create.go`: `s.db.Create(&data)` inserts the row (GORM fills `ID`,
   `CreatedAt`, `UpdatedAt`, default `Status` via `common.SQLModel` tags).
4. Response: `common.SimpleSuccessResponse(data.ID)`.

## Request flow (example: `DELETE /v1/restaurants/:id` — soft delete)

1. `transport/ginrestaurant/delete_restaurant.go` parses `id` from the path.
2. `biz/delete_restaurant.go`: `DeleteRestaurant` first calls
   `store.FindDataWithCondition(ctx, {"id": id})` to load the row, and returns an error
   if `oldData.Status == 0` (already deleted) — guards against double-delete.
3. `storage/delete.go`: `Delete` doesn't actually remove the row; it does
   `UPDATE restaurants SET status = 0 WHERE id = ?` (soft delete).
4. Response: `common.SimpleSuccessResponse(true)`.

## Key shared types

- **`common.SQLModel`** (embedded in every entity model): `ID`, `Status` (1 = active,
  0 = soft-deleted), `CreatedAt`, `UpdatedAt`. Entities only declare their own
  domain fields (`Name`, `Addr`, ...) on top of this.
- **`common.Paging`**: `Page`, `Limit`, `Total`; `Fulfill()` applies defaults
  (page=1, limit=50) when unset.
- **`common.successRes`** (via `NewSuccessResponse`/`SimpleSuccessResponse`): the single
  response envelope used by every handler.
- **`appctx.AppContext`**: currently just wraps `*gorm.DB`; passed into every transport
  constructor so handlers never take `*gorm.DB` directly. Future cross-cutting deps
  (config, logger, S3 client, etc.) get added here, not to individual handlers.

## Known gaps / inconsistencies (as of this writing)

- `GET /v1/restaurants/:id` and `PATCH /v1/restaurants/:id` in `main.go` bypass the
  biz/storage layers entirely (direct `db.Where(...)` calls) and haven't been migrated
  to the `restaurant/biz` + `restaurant/storage` structure yet.
- `ListDataWithCondition`'s `moreKeys ...string` parameter (intended for `Preload`)
  is accepted but never used.
- Soft-delete status uses `int` (`1`/`0`) consistently in code now, but the README's
  wording ("`status = "0"`") is stale relative to that.
