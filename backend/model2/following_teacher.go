// Code generated by SQLBoiler 4.14.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

// FollowingTeacher is an object representing the database table.
type FollowingTeacher struct {
	UserID    uint      `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	TeacherID uint      `boil:"teacher_id" json:"teacher_id" toml:"teacher_id" yaml:"teacher_id"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *followingTeacherR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L followingTeacherL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var FollowingTeacherColumns = struct {
	UserID    string
	TeacherID string
	CreatedAt string
	UpdatedAt string
}{
	UserID:    "user_id",
	TeacherID: "teacher_id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

var FollowingTeacherTableColumns = struct {
	UserID    string
	TeacherID string
	CreatedAt string
	UpdatedAt string
}{
	UserID:    "following_teacher.user_id",
	TeacherID: "following_teacher.teacher_id",
	CreatedAt: "following_teacher.created_at",
	UpdatedAt: "following_teacher.updated_at",
}

// Generated where

var FollowingTeacherWhere = struct {
	UserID    whereHelperuint
	TeacherID whereHelperuint
	CreatedAt whereHelpertime_Time
	UpdatedAt whereHelpertime_Time
}{
	UserID:    whereHelperuint{field: "`following_teacher`.`user_id`"},
	TeacherID: whereHelperuint{field: "`following_teacher`.`teacher_id`"},
	CreatedAt: whereHelpertime_Time{field: "`following_teacher`.`created_at`"},
	UpdatedAt: whereHelpertime_Time{field: "`following_teacher`.`updated_at`"},
}

// FollowingTeacherRels is where relationship names are stored.
var FollowingTeacherRels = struct {
}{}

// followingTeacherR is where relationships are stored.
type followingTeacherR struct {
}

// NewStruct creates a new relationship struct
func (*followingTeacherR) NewStruct() *followingTeacherR {
	return &followingTeacherR{}
}

// followingTeacherL is where Load methods for each relationship are stored.
type followingTeacherL struct{}

var (
	followingTeacherAllColumns            = []string{"user_id", "teacher_id", "created_at", "updated_at"}
	followingTeacherColumnsWithoutDefault = []string{"user_id", "teacher_id", "created_at", "updated_at"}
	followingTeacherColumnsWithDefault    = []string{}
	followingTeacherPrimaryKeyColumns     = []string{"user_id", "teacher_id"}
	followingTeacherGeneratedColumns      = []string{}
)

type (
	// FollowingTeacherSlice is an alias for a slice of pointers to FollowingTeacher.
	// This should almost always be used instead of []FollowingTeacher.
	FollowingTeacherSlice []*FollowingTeacher
	// FollowingTeacherHook is the signature for custom FollowingTeacher hook methods
	FollowingTeacherHook func(context.Context, boil.ContextExecutor, *FollowingTeacher) error

	followingTeacherQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	followingTeacherType                 = reflect.TypeOf(&FollowingTeacher{})
	followingTeacherMapping              = queries.MakeStructMapping(followingTeacherType)
	followingTeacherPrimaryKeyMapping, _ = queries.BindMapping(followingTeacherType, followingTeacherMapping, followingTeacherPrimaryKeyColumns)
	followingTeacherInsertCacheMut       sync.RWMutex
	followingTeacherInsertCache          = make(map[string]insertCache)
	followingTeacherUpdateCacheMut       sync.RWMutex
	followingTeacherUpdateCache          = make(map[string]updateCache)
	followingTeacherUpsertCacheMut       sync.RWMutex
	followingTeacherUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var followingTeacherAfterSelectHooks []FollowingTeacherHook

var followingTeacherBeforeInsertHooks []FollowingTeacherHook
var followingTeacherAfterInsertHooks []FollowingTeacherHook

var followingTeacherBeforeUpdateHooks []FollowingTeacherHook
var followingTeacherAfterUpdateHooks []FollowingTeacherHook

var followingTeacherBeforeDeleteHooks []FollowingTeacherHook
var followingTeacherAfterDeleteHooks []FollowingTeacherHook

var followingTeacherBeforeUpsertHooks []FollowingTeacherHook
var followingTeacherAfterUpsertHooks []FollowingTeacherHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *FollowingTeacher) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range followingTeacherAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *FollowingTeacher) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range followingTeacherBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *FollowingTeacher) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range followingTeacherAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *FollowingTeacher) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range followingTeacherBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *FollowingTeacher) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range followingTeacherAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *FollowingTeacher) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range followingTeacherBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *FollowingTeacher) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range followingTeacherAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *FollowingTeacher) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range followingTeacherBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *FollowingTeacher) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range followingTeacherAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddFollowingTeacherHook registers your hook function for all future operations.
func AddFollowingTeacherHook(hookPoint boil.HookPoint, followingTeacherHook FollowingTeacherHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		followingTeacherAfterSelectHooks = append(followingTeacherAfterSelectHooks, followingTeacherHook)
	case boil.BeforeInsertHook:
		followingTeacherBeforeInsertHooks = append(followingTeacherBeforeInsertHooks, followingTeacherHook)
	case boil.AfterInsertHook:
		followingTeacherAfterInsertHooks = append(followingTeacherAfterInsertHooks, followingTeacherHook)
	case boil.BeforeUpdateHook:
		followingTeacherBeforeUpdateHooks = append(followingTeacherBeforeUpdateHooks, followingTeacherHook)
	case boil.AfterUpdateHook:
		followingTeacherAfterUpdateHooks = append(followingTeacherAfterUpdateHooks, followingTeacherHook)
	case boil.BeforeDeleteHook:
		followingTeacherBeforeDeleteHooks = append(followingTeacherBeforeDeleteHooks, followingTeacherHook)
	case boil.AfterDeleteHook:
		followingTeacherAfterDeleteHooks = append(followingTeacherAfterDeleteHooks, followingTeacherHook)
	case boil.BeforeUpsertHook:
		followingTeacherBeforeUpsertHooks = append(followingTeacherBeforeUpsertHooks, followingTeacherHook)
	case boil.AfterUpsertHook:
		followingTeacherAfterUpsertHooks = append(followingTeacherAfterUpsertHooks, followingTeacherHook)
	}
}

// One returns a single followingTeacher record from the query.
func (q followingTeacherQuery) One(ctx context.Context, exec boil.ContextExecutor) (*FollowingTeacher, error) {
	o := &FollowingTeacher{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model2: failed to execute a one query for following_teacher")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all FollowingTeacher records from the query.
func (q followingTeacherQuery) All(ctx context.Context, exec boil.ContextExecutor) (FollowingTeacherSlice, error) {
	var o []*FollowingTeacher

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model2: failed to assign all query results to FollowingTeacher slice")
	}

	if len(followingTeacherAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all FollowingTeacher records in the query.
func (q followingTeacherQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to count following_teacher rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q followingTeacherQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model2: failed to check if following_teacher exists")
	}

	return count > 0, nil
}

// FollowingTeachers retrieves all the records using an executor.
func FollowingTeachers(mods ...qm.QueryMod) followingTeacherQuery {
	mods = append(mods, qm.From("`following_teacher`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`following_teacher`.*"})
	}

	return followingTeacherQuery{q}
}

// FindFollowingTeacher retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindFollowingTeacher(ctx context.Context, exec boil.ContextExecutor, userID uint, teacherID uint, selectCols ...string) (*FollowingTeacher, error) {
	followingTeacherObj := &FollowingTeacher{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `following_teacher` where `user_id`=? AND `teacher_id`=?", sel,
	)

	q := queries.Raw(query, userID, teacherID)

	err := q.Bind(ctx, exec, followingTeacherObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model2: unable to select from following_teacher")
	}

	if err = followingTeacherObj.doAfterSelectHooks(ctx, exec); err != nil {
		return followingTeacherObj, err
	}

	return followingTeacherObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *FollowingTeacher) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model2: no following_teacher provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(followingTeacherColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	followingTeacherInsertCacheMut.RLock()
	cache, cached := followingTeacherInsertCache[key]
	followingTeacherInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			followingTeacherAllColumns,
			followingTeacherColumnsWithDefault,
			followingTeacherColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(followingTeacherType, followingTeacherMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(followingTeacherType, followingTeacherMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `following_teacher` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `following_teacher` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `following_teacher` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, followingTeacherPrimaryKeyColumns))
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
		return errors.Wrap(err, "model2: unable to insert into following_teacher")
	}

	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.UserID,
		o.TeacherID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "model2: unable to populate default values for following_teacher")
	}

CacheNoHooks:
	if !cached {
		followingTeacherInsertCacheMut.Lock()
		followingTeacherInsertCache[key] = cache
		followingTeacherInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the FollowingTeacher.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *FollowingTeacher) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	followingTeacherUpdateCacheMut.RLock()
	cache, cached := followingTeacherUpdateCache[key]
	followingTeacherUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			followingTeacherAllColumns,
			followingTeacherPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("model2: unable to update following_teacher, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `following_teacher` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, followingTeacherPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(followingTeacherType, followingTeacherMapping, append(wl, followingTeacherPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model2: unable to update following_teacher row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by update for following_teacher")
	}

	if !cached {
		followingTeacherUpdateCacheMut.Lock()
		followingTeacherUpdateCache[key] = cache
		followingTeacherUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q followingTeacherQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to update all for following_teacher")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to retrieve rows affected for following_teacher")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o FollowingTeacherSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), followingTeacherPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `following_teacher` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, followingTeacherPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to update all in followingTeacher slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to retrieve rows affected all in update all followingTeacher")
	}
	return rowsAff, nil
}

var mySQLFollowingTeacherUniqueColumns = []string{}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *FollowingTeacher) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model2: no following_teacher provided for upsert")
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

	nzDefaults := queries.NonZeroDefaultSet(followingTeacherColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLFollowingTeacherUniqueColumns, o)

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

	followingTeacherUpsertCacheMut.RLock()
	cache, cached := followingTeacherUpsertCache[key]
	followingTeacherUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			followingTeacherAllColumns,
			followingTeacherColumnsWithDefault,
			followingTeacherColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			followingTeacherAllColumns,
			followingTeacherPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("model2: unable to upsert following_teacher, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`following_teacher`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `following_teacher` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(followingTeacherType, followingTeacherMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(followingTeacherType, followingTeacherMapping, ret)
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
		return errors.Wrap(err, "model2: unable to upsert for following_teacher")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(followingTeacherType, followingTeacherMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "model2: unable to retrieve unique values for following_teacher")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "model2: unable to populate default values for following_teacher")
	}

CacheNoHooks:
	if !cached {
		followingTeacherUpsertCacheMut.Lock()
		followingTeacherUpsertCache[key] = cache
		followingTeacherUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single FollowingTeacher record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *FollowingTeacher) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("model2: no FollowingTeacher provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), followingTeacherPrimaryKeyMapping)
	sql := "DELETE FROM `following_teacher` WHERE `user_id`=? AND `teacher_id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete from following_teacher")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by delete for following_teacher")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q followingTeacherQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model2: no followingTeacherQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete all from following_teacher")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by deleteall for following_teacher")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o FollowingTeacherSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(followingTeacherBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), followingTeacherPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `following_teacher` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, followingTeacherPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete all from followingTeacher slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by deleteall for following_teacher")
	}

	if len(followingTeacherAfterDeleteHooks) != 0 {
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
func (o *FollowingTeacher) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindFollowingTeacher(ctx, exec, o.UserID, o.TeacherID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *FollowingTeacherSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := FollowingTeacherSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), followingTeacherPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `following_teacher`.* FROM `following_teacher` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, followingTeacherPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model2: unable to reload all in FollowingTeacherSlice")
	}

	*o = slice

	return nil
}

// FollowingTeacherExists checks if the FollowingTeacher row exists.
func FollowingTeacherExists(ctx context.Context, exec boil.ContextExecutor, userID uint, teacherID uint) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `following_teacher` where `user_id`=? AND `teacher_id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, userID, teacherID)
	}
	row := exec.QueryRowContext(ctx, sql, userID, teacherID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model2: unable to check if following_teacher exists")
	}

	return exists, nil
}

// Exists checks if the FollowingTeacher row exists.
func (o *FollowingTeacher) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return FollowingTeacherExists(ctx, exec, o.UserID, o.TeacherID)
}
