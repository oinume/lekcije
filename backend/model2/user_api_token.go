// Code generated by SQLBoiler 4.6.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package model2

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// UserAPIToken is an object representing the database table.
type UserAPIToken struct {
	Token     string    `boil:"token" json:"token" toml:"token" yaml:"token"`
	UserID    uint      `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *userAPITokenR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L userAPITokenL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UserAPITokenColumns = struct {
	Token     string
	UserID    string
	CreatedAt string
	UpdatedAt string
}{
	Token:     "token",
	UserID:    "user_id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

var UserAPITokenTableColumns = struct {
	Token     string
	UserID    string
	CreatedAt string
	UpdatedAt string
}{
	Token:     "user_api_token.token",
	UserID:    "user_api_token.user_id",
	CreatedAt: "user_api_token.created_at",
	UpdatedAt: "user_api_token.updated_at",
}

// Generated where

var UserAPITokenWhere = struct {
	Token     whereHelperstring
	UserID    whereHelperuint
	CreatedAt whereHelpertime_Time
	UpdatedAt whereHelpertime_Time
}{
	Token:     whereHelperstring{field: "`user_api_token`.`token`"},
	UserID:    whereHelperuint{field: "`user_api_token`.`user_id`"},
	CreatedAt: whereHelpertime_Time{field: "`user_api_token`.`created_at`"},
	UpdatedAt: whereHelpertime_Time{field: "`user_api_token`.`updated_at`"},
}

// UserAPITokenRels is where relationship names are stored.
var UserAPITokenRels = struct {
}{}

// userAPITokenR is where relationships are stored.
type userAPITokenR struct {
}

// NewStruct creates a new relationship struct
func (*userAPITokenR) NewStruct() *userAPITokenR {
	return &userAPITokenR{}
}

// userAPITokenL is where Load methods for each relationship are stored.
type userAPITokenL struct{}

var (
	userAPITokenAllColumns            = []string{"token", "user_id", "created_at", "updated_at"}
	userAPITokenColumnsWithoutDefault = []string{"token", "user_id", "created_at", "updated_at"}
	userAPITokenColumnsWithDefault    = []string{}
	userAPITokenPrimaryKeyColumns     = []string{"token"}
)

type (
	// UserAPITokenSlice is an alias for a slice of pointers to UserAPIToken.
	// This should almost always be used instead of []UserAPIToken.
	UserAPITokenSlice []*UserAPIToken
	// UserAPITokenHook is the signature for custom UserAPIToken hook methods
	UserAPITokenHook func(context.Context, boil.ContextExecutor, *UserAPIToken) error

	userAPITokenQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	userAPITokenType                 = reflect.TypeOf(&UserAPIToken{})
	userAPITokenMapping              = queries.MakeStructMapping(userAPITokenType)
	userAPITokenPrimaryKeyMapping, _ = queries.BindMapping(userAPITokenType, userAPITokenMapping, userAPITokenPrimaryKeyColumns)
	userAPITokenInsertCacheMut       sync.RWMutex
	userAPITokenInsertCache          = make(map[string]insertCache)
	userAPITokenUpdateCacheMut       sync.RWMutex
	userAPITokenUpdateCache          = make(map[string]updateCache)
	userAPITokenUpsertCacheMut       sync.RWMutex
	userAPITokenUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var userAPITokenBeforeInsertHooks []UserAPITokenHook
var userAPITokenBeforeUpdateHooks []UserAPITokenHook
var userAPITokenBeforeDeleteHooks []UserAPITokenHook
var userAPITokenBeforeUpsertHooks []UserAPITokenHook

var userAPITokenAfterInsertHooks []UserAPITokenHook
var userAPITokenAfterSelectHooks []UserAPITokenHook
var userAPITokenAfterUpdateHooks []UserAPITokenHook
var userAPITokenAfterDeleteHooks []UserAPITokenHook
var userAPITokenAfterUpsertHooks []UserAPITokenHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *UserAPIToken) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userAPITokenBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *UserAPIToken) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userAPITokenBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *UserAPIToken) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userAPITokenBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *UserAPIToken) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userAPITokenBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *UserAPIToken) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userAPITokenAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *UserAPIToken) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userAPITokenAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *UserAPIToken) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userAPITokenAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *UserAPIToken) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userAPITokenAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *UserAPIToken) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userAPITokenAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUserAPITokenHook registers your hook function for all future operations.
func AddUserAPITokenHook(hookPoint boil.HookPoint, userAPITokenHook UserAPITokenHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		userAPITokenBeforeInsertHooks = append(userAPITokenBeforeInsertHooks, userAPITokenHook)
	case boil.BeforeUpdateHook:
		userAPITokenBeforeUpdateHooks = append(userAPITokenBeforeUpdateHooks, userAPITokenHook)
	case boil.BeforeDeleteHook:
		userAPITokenBeforeDeleteHooks = append(userAPITokenBeforeDeleteHooks, userAPITokenHook)
	case boil.BeforeUpsertHook:
		userAPITokenBeforeUpsertHooks = append(userAPITokenBeforeUpsertHooks, userAPITokenHook)
	case boil.AfterInsertHook:
		userAPITokenAfterInsertHooks = append(userAPITokenAfterInsertHooks, userAPITokenHook)
	case boil.AfterSelectHook:
		userAPITokenAfterSelectHooks = append(userAPITokenAfterSelectHooks, userAPITokenHook)
	case boil.AfterUpdateHook:
		userAPITokenAfterUpdateHooks = append(userAPITokenAfterUpdateHooks, userAPITokenHook)
	case boil.AfterDeleteHook:
		userAPITokenAfterDeleteHooks = append(userAPITokenAfterDeleteHooks, userAPITokenHook)
	case boil.AfterUpsertHook:
		userAPITokenAfterUpsertHooks = append(userAPITokenAfterUpsertHooks, userAPITokenHook)
	}
}

// One returns a single userAPIToken record from the query.
func (q userAPITokenQuery) One(ctx context.Context, exec boil.ContextExecutor) (*UserAPIToken, error) {
	o := &UserAPIToken{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model2: failed to execute a one query for user_api_token")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all UserAPIToken records from the query.
func (q userAPITokenQuery) All(ctx context.Context, exec boil.ContextExecutor) (UserAPITokenSlice, error) {
	var o []*UserAPIToken

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model2: failed to assign all query results to UserAPIToken slice")
	}

	if len(userAPITokenAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all UserAPIToken records in the query.
func (q userAPITokenQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to count user_api_token rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q userAPITokenQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model2: failed to check if user_api_token exists")
	}

	return count > 0, nil
}

// UserAPITokens retrieves all the records using an executor.
func UserAPITokens(mods ...qm.QueryMod) userAPITokenQuery {
	mods = append(mods, qm.From("`user_api_token`"))
	return userAPITokenQuery{NewQuery(mods...)}
}

// FindUserAPIToken retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUserAPIToken(ctx context.Context, exec boil.ContextExecutor, token string, selectCols ...string) (*UserAPIToken, error) {
	userAPITokenObj := &UserAPIToken{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `user_api_token` where `token`=?", sel,
	)

	q := queries.Raw(query, token)

	err := q.Bind(ctx, exec, userAPITokenObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model2: unable to select from user_api_token")
	}

	if err = userAPITokenObj.doAfterSelectHooks(ctx, exec); err != nil {
		return userAPITokenObj, err
	}

	return userAPITokenObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *UserAPIToken) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model2: no user_api_token provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userAPITokenColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	userAPITokenInsertCacheMut.RLock()
	cache, cached := userAPITokenInsertCache[key]
	userAPITokenInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			userAPITokenAllColumns,
			userAPITokenColumnsWithDefault,
			userAPITokenColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(userAPITokenType, userAPITokenMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(userAPITokenType, userAPITokenMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `user_api_token` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `user_api_token` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `user_api_token` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, userAPITokenPrimaryKeyColumns))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	_, err = exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "model2: unable to insert into user_api_token")
	}

	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.Token,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "model2: unable to populate default values for user_api_token")
	}

CacheNoHooks:
	if !cached {
		userAPITokenInsertCacheMut.Lock()
		userAPITokenInsertCache[key] = cache
		userAPITokenInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the UserAPIToken.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *UserAPIToken) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	userAPITokenUpdateCacheMut.RLock()
	cache, cached := userAPITokenUpdateCache[key]
	userAPITokenUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			userAPITokenAllColumns,
			userAPITokenPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("model2: unable to update user_api_token, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `user_api_token` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, userAPITokenPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(userAPITokenType, userAPITokenMapping, append(wl, userAPITokenPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to update user_api_token row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by update for user_api_token")
	}

	if !cached {
		userAPITokenUpdateCacheMut.Lock()
		userAPITokenUpdateCache[key] = cache
		userAPITokenUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q userAPITokenQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to update all for user_api_token")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to retrieve rows affected for user_api_token")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UserAPITokenSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("model2: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userAPITokenPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `user_api_token` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, userAPITokenPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to update all in userAPIToken slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to retrieve rows affected all in update all userAPIToken")
	}
	return rowsAff, nil
}

var mySQLUserAPITokenUniqueColumns = []string{
	"token",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *UserAPIToken) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model2: no user_api_token provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userAPITokenColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLUserAPITokenUniqueColumns, o)

	if len(nzUniques) == 0 {
		return errors.New("cannot upsert with a table that cannot conflict on a unique column")
	}

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzUniques {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	userAPITokenUpsertCacheMut.RLock()
	cache, cached := userAPITokenUpsertCache[key]
	userAPITokenUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			userAPITokenAllColumns,
			userAPITokenColumnsWithDefault,
			userAPITokenColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			userAPITokenAllColumns,
			userAPITokenPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("model2: unable to upsert user_api_token, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`user_api_token`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `user_api_token` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(userAPITokenType, userAPITokenMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(userAPITokenType, userAPITokenMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	_, err = exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "model2: unable to upsert for user_api_token")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(userAPITokenType, userAPITokenMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "model2: unable to retrieve unique values for user_api_token")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "model2: unable to populate default values for user_api_token")
	}

CacheNoHooks:
	if !cached {
		userAPITokenUpsertCacheMut.Lock()
		userAPITokenUpsertCache[key] = cache
		userAPITokenUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single UserAPIToken record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UserAPIToken) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("model2: no UserAPIToken provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), userAPITokenPrimaryKeyMapping)
	sql := "DELETE FROM `user_api_token` WHERE `token`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete from user_api_token")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by delete for user_api_token")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q userAPITokenQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model2: no userAPITokenQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete all from user_api_token")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by deleteall for user_api_token")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UserAPITokenSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(userAPITokenBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userAPITokenPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `user_api_token` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, userAPITokenPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete all from userAPIToken slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by deleteall for user_api_token")
	}

	if len(userAPITokenAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *UserAPIToken) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindUserAPIToken(ctx, exec, o.Token)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserAPITokenSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UserAPITokenSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userAPITokenPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `user_api_token`.* FROM `user_api_token` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, userAPITokenPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model2: unable to reload all in UserAPITokenSlice")
	}

	*o = slice

	return nil
}

// UserAPITokenExists checks if the UserAPIToken row exists.
func UserAPITokenExists(ctx context.Context, exec boil.ContextExecutor, token string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `user_api_token` where `token`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, token)
	}
	row := exec.QueryRowContext(ctx, sql, token)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model2: unable to check if user_api_token exists")
	}

	return exists, nil
}
