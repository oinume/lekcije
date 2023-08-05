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

// LessonStatusLog is an object representing the database table.
type LessonStatusLog struct {
	ID        uint64    `boil:"id" json:"id" toml:"id" yaml:"id"`
	LessonID  uint64    `boil:"lesson_id" json:"lesson_id" toml:"lesson_id" yaml:"lesson_id"`
	Status    string    `boil:"status" json:"status" toml:"status" yaml:"status"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *lessonStatusLogR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L lessonStatusLogL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var LessonStatusLogColumns = struct {
	ID        string
	LessonID  string
	Status    string
	CreatedAt string
}{
	ID:        "id",
	LessonID:  "lesson_id",
	Status:    "status",
	CreatedAt: "created_at",
}

var LessonStatusLogTableColumns = struct {
	ID        string
	LessonID  string
	Status    string
	CreatedAt string
}{
	ID:        "lesson_status_log.id",
	LessonID:  "lesson_status_log.lesson_id",
	Status:    "lesson_status_log.status",
	CreatedAt: "lesson_status_log.created_at",
}

// Generated where

var LessonStatusLogWhere = struct {
	ID        whereHelperuint64
	LessonID  whereHelperuint64
	Status    whereHelperstring
	CreatedAt whereHelpertime_Time
}{
	ID:        whereHelperuint64{field: "`lesson_status_log`.`id`"},
	LessonID:  whereHelperuint64{field: "`lesson_status_log`.`lesson_id`"},
	Status:    whereHelperstring{field: "`lesson_status_log`.`status`"},
	CreatedAt: whereHelpertime_Time{field: "`lesson_status_log`.`created_at`"},
}

// LessonStatusLogRels is where relationship names are stored.
var LessonStatusLogRels = struct {
}{}

// lessonStatusLogR is where relationships are stored.
type lessonStatusLogR struct {
}

// NewStruct creates a new relationship struct
func (*lessonStatusLogR) NewStruct() *lessonStatusLogR {
	return &lessonStatusLogR{}
}

// lessonStatusLogL is where Load methods for each relationship are stored.
type lessonStatusLogL struct{}

var (
	lessonStatusLogAllColumns            = []string{"id", "lesson_id", "status", "created_at"}
	lessonStatusLogColumnsWithoutDefault = []string{"lesson_id", "status", "created_at"}
	lessonStatusLogColumnsWithDefault    = []string{"id"}
	lessonStatusLogPrimaryKeyColumns     = []string{"id"}
	lessonStatusLogGeneratedColumns      = []string{}
)

type (
	// LessonStatusLogSlice is an alias for a slice of pointers to LessonStatusLog.
	// This should almost always be used instead of []LessonStatusLog.
	LessonStatusLogSlice []*LessonStatusLog
	// LessonStatusLogHook is the signature for custom LessonStatusLog hook methods
	LessonStatusLogHook func(context.Context, boil.ContextExecutor, *LessonStatusLog) error

	lessonStatusLogQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	lessonStatusLogType                 = reflect.TypeOf(&LessonStatusLog{})
	lessonStatusLogMapping              = queries.MakeStructMapping(lessonStatusLogType)
	lessonStatusLogPrimaryKeyMapping, _ = queries.BindMapping(lessonStatusLogType, lessonStatusLogMapping, lessonStatusLogPrimaryKeyColumns)
	lessonStatusLogInsertCacheMut       sync.RWMutex
	lessonStatusLogInsertCache          = make(map[string]insertCache)
	lessonStatusLogUpdateCacheMut       sync.RWMutex
	lessonStatusLogUpdateCache          = make(map[string]updateCache)
	lessonStatusLogUpsertCacheMut       sync.RWMutex
	lessonStatusLogUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var lessonStatusLogAfterSelectHooks []LessonStatusLogHook

var lessonStatusLogBeforeInsertHooks []LessonStatusLogHook
var lessonStatusLogAfterInsertHooks []LessonStatusLogHook

var lessonStatusLogBeforeUpdateHooks []LessonStatusLogHook
var lessonStatusLogAfterUpdateHooks []LessonStatusLogHook

var lessonStatusLogBeforeDeleteHooks []LessonStatusLogHook
var lessonStatusLogAfterDeleteHooks []LessonStatusLogHook

var lessonStatusLogBeforeUpsertHooks []LessonStatusLogHook
var lessonStatusLogAfterUpsertHooks []LessonStatusLogHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *LessonStatusLog) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range lessonStatusLogAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *LessonStatusLog) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range lessonStatusLogBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *LessonStatusLog) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range lessonStatusLogAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *LessonStatusLog) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range lessonStatusLogBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *LessonStatusLog) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range lessonStatusLogAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *LessonStatusLog) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range lessonStatusLogBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *LessonStatusLog) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range lessonStatusLogAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *LessonStatusLog) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range lessonStatusLogBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *LessonStatusLog) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range lessonStatusLogAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddLessonStatusLogHook registers your hook function for all future operations.
func AddLessonStatusLogHook(hookPoint boil.HookPoint, lessonStatusLogHook LessonStatusLogHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		lessonStatusLogAfterSelectHooks = append(lessonStatusLogAfterSelectHooks, lessonStatusLogHook)
	case boil.BeforeInsertHook:
		lessonStatusLogBeforeInsertHooks = append(lessonStatusLogBeforeInsertHooks, lessonStatusLogHook)
	case boil.AfterInsertHook:
		lessonStatusLogAfterInsertHooks = append(lessonStatusLogAfterInsertHooks, lessonStatusLogHook)
	case boil.BeforeUpdateHook:
		lessonStatusLogBeforeUpdateHooks = append(lessonStatusLogBeforeUpdateHooks, lessonStatusLogHook)
	case boil.AfterUpdateHook:
		lessonStatusLogAfterUpdateHooks = append(lessonStatusLogAfterUpdateHooks, lessonStatusLogHook)
	case boil.BeforeDeleteHook:
		lessonStatusLogBeforeDeleteHooks = append(lessonStatusLogBeforeDeleteHooks, lessonStatusLogHook)
	case boil.AfterDeleteHook:
		lessonStatusLogAfterDeleteHooks = append(lessonStatusLogAfterDeleteHooks, lessonStatusLogHook)
	case boil.BeforeUpsertHook:
		lessonStatusLogBeforeUpsertHooks = append(lessonStatusLogBeforeUpsertHooks, lessonStatusLogHook)
	case boil.AfterUpsertHook:
		lessonStatusLogAfterUpsertHooks = append(lessonStatusLogAfterUpsertHooks, lessonStatusLogHook)
	}
}

// One returns a single lessonStatusLog record from the query.
func (q lessonStatusLogQuery) One(ctx context.Context, exec boil.ContextExecutor) (*LessonStatusLog, error) {
	o := &LessonStatusLog{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model2: failed to execute a one query for lesson_status_log")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all LessonStatusLog records from the query.
func (q lessonStatusLogQuery) All(ctx context.Context, exec boil.ContextExecutor) (LessonStatusLogSlice, error) {
	var o []*LessonStatusLog

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model2: failed to assign all query results to LessonStatusLog slice")
	}

	if len(lessonStatusLogAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all LessonStatusLog records in the query.
func (q lessonStatusLogQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to count lesson_status_log rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q lessonStatusLogQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model2: failed to check if lesson_status_log exists")
	}

	return count > 0, nil
}

// LessonStatusLogs retrieves all the records using an executor.
func LessonStatusLogs(mods ...qm.QueryMod) lessonStatusLogQuery {
	mods = append(mods, qm.From("`lesson_status_log`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`lesson_status_log`.*"})
	}

	return lessonStatusLogQuery{q}
}

// FindLessonStatusLog retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindLessonStatusLog(ctx context.Context, exec boil.ContextExecutor, iD uint64, selectCols ...string) (*LessonStatusLog, error) {
	lessonStatusLogObj := &LessonStatusLog{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `lesson_status_log` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, lessonStatusLogObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model2: unable to select from lesson_status_log")
	}

	if err = lessonStatusLogObj.doAfterSelectHooks(ctx, exec); err != nil {
		return lessonStatusLogObj, err
	}

	return lessonStatusLogObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *LessonStatusLog) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model2: no lesson_status_log provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(lessonStatusLogColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	lessonStatusLogInsertCacheMut.RLock()
	cache, cached := lessonStatusLogInsertCache[key]
	lessonStatusLogInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			lessonStatusLogAllColumns,
			lessonStatusLogColumnsWithDefault,
			lessonStatusLogColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(lessonStatusLogType, lessonStatusLogMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(lessonStatusLogType, lessonStatusLogMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `lesson_status_log` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `lesson_status_log` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `lesson_status_log` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, lessonStatusLogPrimaryKeyColumns))
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
	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "model2: unable to insert into lesson_status_log")
	}

	var lastID int64
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = uint64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == lessonStatusLogMapping["id"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.ID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "model2: unable to populate default values for lesson_status_log")
	}

CacheNoHooks:
	if !cached {
		lessonStatusLogInsertCacheMut.Lock()
		lessonStatusLogInsertCache[key] = cache
		lessonStatusLogInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the LessonStatusLog.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *LessonStatusLog) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	lessonStatusLogUpdateCacheMut.RLock()
	cache, cached := lessonStatusLogUpdateCache[key]
	lessonStatusLogUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			lessonStatusLogAllColumns,
			lessonStatusLogPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("model2: unable to update lesson_status_log, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `lesson_status_log` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, lessonStatusLogPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(lessonStatusLogType, lessonStatusLogMapping, append(wl, lessonStatusLogPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model2: unable to update lesson_status_log row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by update for lesson_status_log")
	}

	if !cached {
		lessonStatusLogUpdateCacheMut.Lock()
		lessonStatusLogUpdateCache[key] = cache
		lessonStatusLogUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q lessonStatusLogQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to update all for lesson_status_log")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to retrieve rows affected for lesson_status_log")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o LessonStatusLogSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), lessonStatusLogPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `lesson_status_log` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, lessonStatusLogPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to update all in lessonStatusLog slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to retrieve rows affected all in update all lessonStatusLog")
	}
	return rowsAff, nil
}

var mySQLLessonStatusLogUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *LessonStatusLog) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model2: no lesson_status_log provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(lessonStatusLogColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLLessonStatusLogUniqueColumns, o)

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

	lessonStatusLogUpsertCacheMut.RLock()
	cache, cached := lessonStatusLogUpsertCache[key]
	lessonStatusLogUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			lessonStatusLogAllColumns,
			lessonStatusLogColumnsWithDefault,
			lessonStatusLogColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			lessonStatusLogAllColumns,
			lessonStatusLogPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("model2: unable to upsert lesson_status_log, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`lesson_status_log`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `lesson_status_log` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(lessonStatusLogType, lessonStatusLogMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(lessonStatusLogType, lessonStatusLogMapping, ret)
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
	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "model2: unable to upsert for lesson_status_log")
	}

	var lastID int64
	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = uint64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == lessonStatusLogMapping["id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(lessonStatusLogType, lessonStatusLogMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "model2: unable to retrieve unique values for lesson_status_log")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "model2: unable to populate default values for lesson_status_log")
	}

CacheNoHooks:
	if !cached {
		lessonStatusLogUpsertCacheMut.Lock()
		lessonStatusLogUpsertCache[key] = cache
		lessonStatusLogUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single LessonStatusLog record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *LessonStatusLog) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("model2: no LessonStatusLog provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), lessonStatusLogPrimaryKeyMapping)
	sql := "DELETE FROM `lesson_status_log` WHERE `id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete from lesson_status_log")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by delete for lesson_status_log")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q lessonStatusLogQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model2: no lessonStatusLogQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete all from lesson_status_log")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by deleteall for lesson_status_log")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o LessonStatusLogSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(lessonStatusLogBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), lessonStatusLogPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `lesson_status_log` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, lessonStatusLogPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete all from lessonStatusLog slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by deleteall for lesson_status_log")
	}

	if len(lessonStatusLogAfterDeleteHooks) != 0 {
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
func (o *LessonStatusLog) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindLessonStatusLog(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *LessonStatusLogSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := LessonStatusLogSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), lessonStatusLogPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `lesson_status_log`.* FROM `lesson_status_log` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, lessonStatusLogPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model2: unable to reload all in LessonStatusLogSlice")
	}

	*o = slice

	return nil
}

// LessonStatusLogExists checks if the LessonStatusLog row exists.
func LessonStatusLogExists(ctx context.Context, exec boil.ContextExecutor, iD uint64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `lesson_status_log` where `id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model2: unable to check if lesson_status_log exists")
	}

	return exists, nil
}

// Exists checks if the LessonStatusLog row exists.
func (o *LessonStatusLog) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return LessonStatusLogExists(ctx, exec, o.ID)
}
