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

// StatNewLessonNotifier is an object representing the database table.
type StatNewLessonNotifier struct {
	Date    time.Time `boil:"date" json:"date" toml:"date" yaml:"date"`
	Event   string    `boil:"event" json:"event" toml:"event" yaml:"event"`
	Count   uint      `boil:"count" json:"count" toml:"count" yaml:"count"`
	UuCount uint      `boil:"uu_count" json:"uu_count" toml:"uu_count" yaml:"uu_count"`

	R *statNewLessonNotifierR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L statNewLessonNotifierL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var StatNewLessonNotifierColumns = struct {
	Date    string
	Event   string
	Count   string
	UuCount string
}{
	Date:    "date",
	Event:   "event",
	Count:   "count",
	UuCount: "uu_count",
}

var StatNewLessonNotifierTableColumns = struct {
	Date    string
	Event   string
	Count   string
	UuCount string
}{
	Date:    "stat_new_lesson_notifier.date",
	Event:   "stat_new_lesson_notifier.event",
	Count:   "stat_new_lesson_notifier.count",
	UuCount: "stat_new_lesson_notifier.uu_count",
}

// Generated where

var StatNewLessonNotifierWhere = struct {
	Date    whereHelpertime_Time
	Event   whereHelperstring
	Count   whereHelperuint
	UuCount whereHelperuint
}{
	Date:    whereHelpertime_Time{field: "`stat_new_lesson_notifier`.`date`"},
	Event:   whereHelperstring{field: "`stat_new_lesson_notifier`.`event`"},
	Count:   whereHelperuint{field: "`stat_new_lesson_notifier`.`count`"},
	UuCount: whereHelperuint{field: "`stat_new_lesson_notifier`.`uu_count`"},
}

// StatNewLessonNotifierRels is where relationship names are stored.
var StatNewLessonNotifierRels = struct {
}{}

// statNewLessonNotifierR is where relationships are stored.
type statNewLessonNotifierR struct {
}

// NewStruct creates a new relationship struct
func (*statNewLessonNotifierR) NewStruct() *statNewLessonNotifierR {
	return &statNewLessonNotifierR{}
}

// statNewLessonNotifierL is where Load methods for each relationship are stored.
type statNewLessonNotifierL struct{}

var (
	statNewLessonNotifierAllColumns            = []string{"date", "event", "count", "uu_count"}
	statNewLessonNotifierColumnsWithoutDefault = []string{"date", "event", "count", "uu_count"}
	statNewLessonNotifierColumnsWithDefault    = []string{}
	statNewLessonNotifierPrimaryKeyColumns     = []string{"date", "event"}
	statNewLessonNotifierGeneratedColumns      = []string{}
)

type (
	// StatNewLessonNotifierSlice is an alias for a slice of pointers to StatNewLessonNotifier.
	// This should almost always be used instead of []StatNewLessonNotifier.
	StatNewLessonNotifierSlice []*StatNewLessonNotifier
	// StatNewLessonNotifierHook is the signature for custom StatNewLessonNotifier hook methods
	StatNewLessonNotifierHook func(context.Context, boil.ContextExecutor, *StatNewLessonNotifier) error

	statNewLessonNotifierQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	statNewLessonNotifierType                 = reflect.TypeOf(&StatNewLessonNotifier{})
	statNewLessonNotifierMapping              = queries.MakeStructMapping(statNewLessonNotifierType)
	statNewLessonNotifierPrimaryKeyMapping, _ = queries.BindMapping(statNewLessonNotifierType, statNewLessonNotifierMapping, statNewLessonNotifierPrimaryKeyColumns)
	statNewLessonNotifierInsertCacheMut       sync.RWMutex
	statNewLessonNotifierInsertCache          = make(map[string]insertCache)
	statNewLessonNotifierUpdateCacheMut       sync.RWMutex
	statNewLessonNotifierUpdateCache          = make(map[string]updateCache)
	statNewLessonNotifierUpsertCacheMut       sync.RWMutex
	statNewLessonNotifierUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var statNewLessonNotifierAfterSelectHooks []StatNewLessonNotifierHook

var statNewLessonNotifierBeforeInsertHooks []StatNewLessonNotifierHook
var statNewLessonNotifierAfterInsertHooks []StatNewLessonNotifierHook

var statNewLessonNotifierBeforeUpdateHooks []StatNewLessonNotifierHook
var statNewLessonNotifierAfterUpdateHooks []StatNewLessonNotifierHook

var statNewLessonNotifierBeforeDeleteHooks []StatNewLessonNotifierHook
var statNewLessonNotifierAfterDeleteHooks []StatNewLessonNotifierHook

var statNewLessonNotifierBeforeUpsertHooks []StatNewLessonNotifierHook
var statNewLessonNotifierAfterUpsertHooks []StatNewLessonNotifierHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *StatNewLessonNotifier) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statNewLessonNotifierAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *StatNewLessonNotifier) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statNewLessonNotifierBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *StatNewLessonNotifier) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statNewLessonNotifierAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *StatNewLessonNotifier) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statNewLessonNotifierBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *StatNewLessonNotifier) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statNewLessonNotifierAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *StatNewLessonNotifier) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statNewLessonNotifierBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *StatNewLessonNotifier) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statNewLessonNotifierAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *StatNewLessonNotifier) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statNewLessonNotifierBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *StatNewLessonNotifier) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statNewLessonNotifierAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddStatNewLessonNotifierHook registers your hook function for all future operations.
func AddStatNewLessonNotifierHook(hookPoint boil.HookPoint, statNewLessonNotifierHook StatNewLessonNotifierHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		statNewLessonNotifierAfterSelectHooks = append(statNewLessonNotifierAfterSelectHooks, statNewLessonNotifierHook)
	case boil.BeforeInsertHook:
		statNewLessonNotifierBeforeInsertHooks = append(statNewLessonNotifierBeforeInsertHooks, statNewLessonNotifierHook)
	case boil.AfterInsertHook:
		statNewLessonNotifierAfterInsertHooks = append(statNewLessonNotifierAfterInsertHooks, statNewLessonNotifierHook)
	case boil.BeforeUpdateHook:
		statNewLessonNotifierBeforeUpdateHooks = append(statNewLessonNotifierBeforeUpdateHooks, statNewLessonNotifierHook)
	case boil.AfterUpdateHook:
		statNewLessonNotifierAfterUpdateHooks = append(statNewLessonNotifierAfterUpdateHooks, statNewLessonNotifierHook)
	case boil.BeforeDeleteHook:
		statNewLessonNotifierBeforeDeleteHooks = append(statNewLessonNotifierBeforeDeleteHooks, statNewLessonNotifierHook)
	case boil.AfterDeleteHook:
		statNewLessonNotifierAfterDeleteHooks = append(statNewLessonNotifierAfterDeleteHooks, statNewLessonNotifierHook)
	case boil.BeforeUpsertHook:
		statNewLessonNotifierBeforeUpsertHooks = append(statNewLessonNotifierBeforeUpsertHooks, statNewLessonNotifierHook)
	case boil.AfterUpsertHook:
		statNewLessonNotifierAfterUpsertHooks = append(statNewLessonNotifierAfterUpsertHooks, statNewLessonNotifierHook)
	}
}

// One returns a single statNewLessonNotifier record from the query.
func (q statNewLessonNotifierQuery) One(ctx context.Context, exec boil.ContextExecutor) (*StatNewLessonNotifier, error) {
	o := &StatNewLessonNotifier{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model2: failed to execute a one query for stat_new_lesson_notifier")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all StatNewLessonNotifier records from the query.
func (q statNewLessonNotifierQuery) All(ctx context.Context, exec boil.ContextExecutor) (StatNewLessonNotifierSlice, error) {
	var o []*StatNewLessonNotifier

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model2: failed to assign all query results to StatNewLessonNotifier slice")
	}

	if len(statNewLessonNotifierAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all StatNewLessonNotifier records in the query.
func (q statNewLessonNotifierQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to count stat_new_lesson_notifier rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q statNewLessonNotifierQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model2: failed to check if stat_new_lesson_notifier exists")
	}

	return count > 0, nil
}

// StatNewLessonNotifiers retrieves all the records using an executor.
func StatNewLessonNotifiers(mods ...qm.QueryMod) statNewLessonNotifierQuery {
	mods = append(mods, qm.From("`stat_new_lesson_notifier`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`stat_new_lesson_notifier`.*"})
	}

	return statNewLessonNotifierQuery{q}
}

// FindStatNewLessonNotifier retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindStatNewLessonNotifier(ctx context.Context, exec boil.ContextExecutor, date time.Time, event string, selectCols ...string) (*StatNewLessonNotifier, error) {
	statNewLessonNotifierObj := &StatNewLessonNotifier{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `stat_new_lesson_notifier` where `date`=? AND `event`=?", sel,
	)

	q := queries.Raw(query, date, event)

	err := q.Bind(ctx, exec, statNewLessonNotifierObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model2: unable to select from stat_new_lesson_notifier")
	}

	if err = statNewLessonNotifierObj.doAfterSelectHooks(ctx, exec); err != nil {
		return statNewLessonNotifierObj, err
	}

	return statNewLessonNotifierObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *StatNewLessonNotifier) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model2: no stat_new_lesson_notifier provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(statNewLessonNotifierColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	statNewLessonNotifierInsertCacheMut.RLock()
	cache, cached := statNewLessonNotifierInsertCache[key]
	statNewLessonNotifierInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			statNewLessonNotifierAllColumns,
			statNewLessonNotifierColumnsWithDefault,
			statNewLessonNotifierColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(statNewLessonNotifierType, statNewLessonNotifierMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(statNewLessonNotifierType, statNewLessonNotifierMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `stat_new_lesson_notifier` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `stat_new_lesson_notifier` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `stat_new_lesson_notifier` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, statNewLessonNotifierPrimaryKeyColumns))
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
		return errors.Wrap(err, "model2: unable to insert into stat_new_lesson_notifier")
	}

	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.Date,
		o.Event,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "model2: unable to populate default values for stat_new_lesson_notifier")
	}

CacheNoHooks:
	if !cached {
		statNewLessonNotifierInsertCacheMut.Lock()
		statNewLessonNotifierInsertCache[key] = cache
		statNewLessonNotifierInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the StatNewLessonNotifier.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *StatNewLessonNotifier) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	statNewLessonNotifierUpdateCacheMut.RLock()
	cache, cached := statNewLessonNotifierUpdateCache[key]
	statNewLessonNotifierUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			statNewLessonNotifierAllColumns,
			statNewLessonNotifierPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("model2: unable to update stat_new_lesson_notifier, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `stat_new_lesson_notifier` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, statNewLessonNotifierPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(statNewLessonNotifierType, statNewLessonNotifierMapping, append(wl, statNewLessonNotifierPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model2: unable to update stat_new_lesson_notifier row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by update for stat_new_lesson_notifier")
	}

	if !cached {
		statNewLessonNotifierUpdateCacheMut.Lock()
		statNewLessonNotifierUpdateCache[key] = cache
		statNewLessonNotifierUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q statNewLessonNotifierQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to update all for stat_new_lesson_notifier")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to retrieve rows affected for stat_new_lesson_notifier")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o StatNewLessonNotifierSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), statNewLessonNotifierPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `stat_new_lesson_notifier` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, statNewLessonNotifierPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to update all in statNewLessonNotifier slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to retrieve rows affected all in update all statNewLessonNotifier")
	}
	return rowsAff, nil
}

var mySQLStatNewLessonNotifierUniqueColumns = []string{}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *StatNewLessonNotifier) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model2: no stat_new_lesson_notifier provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(statNewLessonNotifierColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLStatNewLessonNotifierUniqueColumns, o)

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

	statNewLessonNotifierUpsertCacheMut.RLock()
	cache, cached := statNewLessonNotifierUpsertCache[key]
	statNewLessonNotifierUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			statNewLessonNotifierAllColumns,
			statNewLessonNotifierColumnsWithDefault,
			statNewLessonNotifierColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			statNewLessonNotifierAllColumns,
			statNewLessonNotifierPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("model2: unable to upsert stat_new_lesson_notifier, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`stat_new_lesson_notifier`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `stat_new_lesson_notifier` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(statNewLessonNotifierType, statNewLessonNotifierMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(statNewLessonNotifierType, statNewLessonNotifierMapping, ret)
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
		return errors.Wrap(err, "model2: unable to upsert for stat_new_lesson_notifier")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(statNewLessonNotifierType, statNewLessonNotifierMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "model2: unable to retrieve unique values for stat_new_lesson_notifier")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "model2: unable to populate default values for stat_new_lesson_notifier")
	}

CacheNoHooks:
	if !cached {
		statNewLessonNotifierUpsertCacheMut.Lock()
		statNewLessonNotifierUpsertCache[key] = cache
		statNewLessonNotifierUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single StatNewLessonNotifier record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *StatNewLessonNotifier) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("model2: no StatNewLessonNotifier provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), statNewLessonNotifierPrimaryKeyMapping)
	sql := "DELETE FROM `stat_new_lesson_notifier` WHERE `date`=? AND `event`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete from stat_new_lesson_notifier")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by delete for stat_new_lesson_notifier")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q statNewLessonNotifierQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model2: no statNewLessonNotifierQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete all from stat_new_lesson_notifier")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by deleteall for stat_new_lesson_notifier")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o StatNewLessonNotifierSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(statNewLessonNotifierBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), statNewLessonNotifierPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `stat_new_lesson_notifier` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, statNewLessonNotifierPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete all from statNewLessonNotifier slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by deleteall for stat_new_lesson_notifier")
	}

	if len(statNewLessonNotifierAfterDeleteHooks) != 0 {
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
func (o *StatNewLessonNotifier) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindStatNewLessonNotifier(ctx, exec, o.Date, o.Event)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *StatNewLessonNotifierSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := StatNewLessonNotifierSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), statNewLessonNotifierPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `stat_new_lesson_notifier`.* FROM `stat_new_lesson_notifier` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, statNewLessonNotifierPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model2: unable to reload all in StatNewLessonNotifierSlice")
	}

	*o = slice

	return nil
}

// StatNewLessonNotifierExists checks if the StatNewLessonNotifier row exists.
func StatNewLessonNotifierExists(ctx context.Context, exec boil.ContextExecutor, date time.Time, event string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `stat_new_lesson_notifier` where `date`=? AND `event`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, date, event)
	}
	row := exec.QueryRowContext(ctx, sql, date, event)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model2: unable to check if stat_new_lesson_notifier exists")
	}

	return exists, nil
}

// Exists checks if the StatNewLessonNotifier row exists.
func (o *StatNewLessonNotifier) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return StatNewLessonNotifierExists(ctx, exec, o.Date, o.Event)
}
