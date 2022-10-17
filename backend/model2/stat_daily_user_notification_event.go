// Code generated by SQLBoiler 4.13.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

// StatDailyUserNotificationEvent is an object representing the database table.
type StatDailyUserNotificationEvent struct {
	Date   time.Time `boil:"date" json:"date" toml:"date" yaml:"date"`
	UserID uint      `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	Event  string    `boil:"event" json:"event" toml:"event" yaml:"event"`
	Count  uint      `boil:"count" json:"count" toml:"count" yaml:"count"`

	R *statDailyUserNotificationEventR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L statDailyUserNotificationEventL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var StatDailyUserNotificationEventColumns = struct {
	Date   string
	UserID string
	Event  string
	Count  string
}{
	Date:   "date",
	UserID: "user_id",
	Event:  "event",
	Count:  "count",
}

var StatDailyUserNotificationEventTableColumns = struct {
	Date   string
	UserID string
	Event  string
	Count  string
}{
	Date:   "stat_daily_user_notification_event.date",
	UserID: "stat_daily_user_notification_event.user_id",
	Event:  "stat_daily_user_notification_event.event",
	Count:  "stat_daily_user_notification_event.count",
}

// Generated where

var StatDailyUserNotificationEventWhere = struct {
	Date   whereHelpertime_Time
	UserID whereHelperuint
	Event  whereHelperstring
	Count  whereHelperuint
}{
	Date:   whereHelpertime_Time{field: "`stat_daily_user_notification_event`.`date`"},
	UserID: whereHelperuint{field: "`stat_daily_user_notification_event`.`user_id`"},
	Event:  whereHelperstring{field: "`stat_daily_user_notification_event`.`event`"},
	Count:  whereHelperuint{field: "`stat_daily_user_notification_event`.`count`"},
}

// StatDailyUserNotificationEventRels is where relationship names are stored.
var StatDailyUserNotificationEventRels = struct {
}{}

// statDailyUserNotificationEventR is where relationships are stored.
type statDailyUserNotificationEventR struct {
}

// NewStruct creates a new relationship struct
func (*statDailyUserNotificationEventR) NewStruct() *statDailyUserNotificationEventR {
	return &statDailyUserNotificationEventR{}
}

// statDailyUserNotificationEventL is where Load methods for each relationship are stored.
type statDailyUserNotificationEventL struct{}

var (
	statDailyUserNotificationEventAllColumns            = []string{"date", "user_id", "event", "count"}
	statDailyUserNotificationEventColumnsWithoutDefault = []string{"date", "user_id", "event", "count"}
	statDailyUserNotificationEventColumnsWithDefault    = []string{}
	statDailyUserNotificationEventPrimaryKeyColumns     = []string{"date", "user_id", "event"}
	statDailyUserNotificationEventGeneratedColumns      = []string{}
)

type (
	// StatDailyUserNotificationEventSlice is an alias for a slice of pointers to StatDailyUserNotificationEvent.
	// This should almost always be used instead of []StatDailyUserNotificationEvent.
	StatDailyUserNotificationEventSlice []*StatDailyUserNotificationEvent
	// StatDailyUserNotificationEventHook is the signature for custom StatDailyUserNotificationEvent hook methods
	StatDailyUserNotificationEventHook func(context.Context, boil.ContextExecutor, *StatDailyUserNotificationEvent) error

	statDailyUserNotificationEventQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	statDailyUserNotificationEventType                 = reflect.TypeOf(&StatDailyUserNotificationEvent{})
	statDailyUserNotificationEventMapping              = queries.MakeStructMapping(statDailyUserNotificationEventType)
	statDailyUserNotificationEventPrimaryKeyMapping, _ = queries.BindMapping(statDailyUserNotificationEventType, statDailyUserNotificationEventMapping, statDailyUserNotificationEventPrimaryKeyColumns)
	statDailyUserNotificationEventInsertCacheMut       sync.RWMutex
	statDailyUserNotificationEventInsertCache          = make(map[string]insertCache)
	statDailyUserNotificationEventUpdateCacheMut       sync.RWMutex
	statDailyUserNotificationEventUpdateCache          = make(map[string]updateCache)
	statDailyUserNotificationEventUpsertCacheMut       sync.RWMutex
	statDailyUserNotificationEventUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var statDailyUserNotificationEventAfterSelectHooks []StatDailyUserNotificationEventHook

var statDailyUserNotificationEventBeforeInsertHooks []StatDailyUserNotificationEventHook
var statDailyUserNotificationEventAfterInsertHooks []StatDailyUserNotificationEventHook

var statDailyUserNotificationEventBeforeUpdateHooks []StatDailyUserNotificationEventHook
var statDailyUserNotificationEventAfterUpdateHooks []StatDailyUserNotificationEventHook

var statDailyUserNotificationEventBeforeDeleteHooks []StatDailyUserNotificationEventHook
var statDailyUserNotificationEventAfterDeleteHooks []StatDailyUserNotificationEventHook

var statDailyUserNotificationEventBeforeUpsertHooks []StatDailyUserNotificationEventHook
var statDailyUserNotificationEventAfterUpsertHooks []StatDailyUserNotificationEventHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *StatDailyUserNotificationEvent) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statDailyUserNotificationEventAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *StatDailyUserNotificationEvent) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statDailyUserNotificationEventBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *StatDailyUserNotificationEvent) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statDailyUserNotificationEventAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *StatDailyUserNotificationEvent) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statDailyUserNotificationEventBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *StatDailyUserNotificationEvent) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statDailyUserNotificationEventAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *StatDailyUserNotificationEvent) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statDailyUserNotificationEventBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *StatDailyUserNotificationEvent) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statDailyUserNotificationEventAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *StatDailyUserNotificationEvent) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statDailyUserNotificationEventBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *StatDailyUserNotificationEvent) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range statDailyUserNotificationEventAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddStatDailyUserNotificationEventHook registers your hook function for all future operations.
func AddStatDailyUserNotificationEventHook(hookPoint boil.HookPoint, statDailyUserNotificationEventHook StatDailyUserNotificationEventHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		statDailyUserNotificationEventAfterSelectHooks = append(statDailyUserNotificationEventAfterSelectHooks, statDailyUserNotificationEventHook)
	case boil.BeforeInsertHook:
		statDailyUserNotificationEventBeforeInsertHooks = append(statDailyUserNotificationEventBeforeInsertHooks, statDailyUserNotificationEventHook)
	case boil.AfterInsertHook:
		statDailyUserNotificationEventAfterInsertHooks = append(statDailyUserNotificationEventAfterInsertHooks, statDailyUserNotificationEventHook)
	case boil.BeforeUpdateHook:
		statDailyUserNotificationEventBeforeUpdateHooks = append(statDailyUserNotificationEventBeforeUpdateHooks, statDailyUserNotificationEventHook)
	case boil.AfterUpdateHook:
		statDailyUserNotificationEventAfterUpdateHooks = append(statDailyUserNotificationEventAfterUpdateHooks, statDailyUserNotificationEventHook)
	case boil.BeforeDeleteHook:
		statDailyUserNotificationEventBeforeDeleteHooks = append(statDailyUserNotificationEventBeforeDeleteHooks, statDailyUserNotificationEventHook)
	case boil.AfterDeleteHook:
		statDailyUserNotificationEventAfterDeleteHooks = append(statDailyUserNotificationEventAfterDeleteHooks, statDailyUserNotificationEventHook)
	case boil.BeforeUpsertHook:
		statDailyUserNotificationEventBeforeUpsertHooks = append(statDailyUserNotificationEventBeforeUpsertHooks, statDailyUserNotificationEventHook)
	case boil.AfterUpsertHook:
		statDailyUserNotificationEventAfterUpsertHooks = append(statDailyUserNotificationEventAfterUpsertHooks, statDailyUserNotificationEventHook)
	}
}

// One returns a single statDailyUserNotificationEvent record from the query.
func (q statDailyUserNotificationEventQuery) One(ctx context.Context, exec boil.ContextExecutor) (*StatDailyUserNotificationEvent, error) {
	o := &StatDailyUserNotificationEvent{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model2: failed to execute a one query for stat_daily_user_notification_event")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all StatDailyUserNotificationEvent records from the query.
func (q statDailyUserNotificationEventQuery) All(ctx context.Context, exec boil.ContextExecutor) (StatDailyUserNotificationEventSlice, error) {
	var o []*StatDailyUserNotificationEvent

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model2: failed to assign all query results to StatDailyUserNotificationEvent slice")
	}

	if len(statDailyUserNotificationEventAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all StatDailyUserNotificationEvent records in the query.
func (q statDailyUserNotificationEventQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to count stat_daily_user_notification_event rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q statDailyUserNotificationEventQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model2: failed to check if stat_daily_user_notification_event exists")
	}

	return count > 0, nil
}

// StatDailyUserNotificationEvents retrieves all the records using an executor.
func StatDailyUserNotificationEvents(mods ...qm.QueryMod) statDailyUserNotificationEventQuery {
	mods = append(mods, qm.From("`stat_daily_user_notification_event`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`stat_daily_user_notification_event`.*"})
	}

	return statDailyUserNotificationEventQuery{q}
}

// FindStatDailyUserNotificationEvent retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindStatDailyUserNotificationEvent(ctx context.Context, exec boil.ContextExecutor, date time.Time, userID uint, event string, selectCols ...string) (*StatDailyUserNotificationEvent, error) {
	statDailyUserNotificationEventObj := &StatDailyUserNotificationEvent{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `stat_daily_user_notification_event` where `date`=? AND `user_id`=? AND `event`=?", sel,
	)

	q := queries.Raw(query, date, userID, event)

	err := q.Bind(ctx, exec, statDailyUserNotificationEventObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model2: unable to select from stat_daily_user_notification_event")
	}

	if err = statDailyUserNotificationEventObj.doAfterSelectHooks(ctx, exec); err != nil {
		return statDailyUserNotificationEventObj, err
	}

	return statDailyUserNotificationEventObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *StatDailyUserNotificationEvent) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model2: no stat_daily_user_notification_event provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(statDailyUserNotificationEventColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	statDailyUserNotificationEventInsertCacheMut.RLock()
	cache, cached := statDailyUserNotificationEventInsertCache[key]
	statDailyUserNotificationEventInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			statDailyUserNotificationEventAllColumns,
			statDailyUserNotificationEventColumnsWithDefault,
			statDailyUserNotificationEventColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(statDailyUserNotificationEventType, statDailyUserNotificationEventMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(statDailyUserNotificationEventType, statDailyUserNotificationEventMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `stat_daily_user_notification_event` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `stat_daily_user_notification_event` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `stat_daily_user_notification_event` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, statDailyUserNotificationEventPrimaryKeyColumns))
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
		return errors.Wrap(err, "model2: unable to insert into stat_daily_user_notification_event")
	}

	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.Date,
		o.UserID,
		o.Event,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "model2: unable to populate default values for stat_daily_user_notification_event")
	}

CacheNoHooks:
	if !cached {
		statDailyUserNotificationEventInsertCacheMut.Lock()
		statDailyUserNotificationEventInsertCache[key] = cache
		statDailyUserNotificationEventInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the StatDailyUserNotificationEvent.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *StatDailyUserNotificationEvent) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	statDailyUserNotificationEventUpdateCacheMut.RLock()
	cache, cached := statDailyUserNotificationEventUpdateCache[key]
	statDailyUserNotificationEventUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			statDailyUserNotificationEventAllColumns,
			statDailyUserNotificationEventPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("model2: unable to update stat_daily_user_notification_event, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `stat_daily_user_notification_event` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, statDailyUserNotificationEventPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(statDailyUserNotificationEventType, statDailyUserNotificationEventMapping, append(wl, statDailyUserNotificationEventPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model2: unable to update stat_daily_user_notification_event row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by update for stat_daily_user_notification_event")
	}

	if !cached {
		statDailyUserNotificationEventUpdateCacheMut.Lock()
		statDailyUserNotificationEventUpdateCache[key] = cache
		statDailyUserNotificationEventUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q statDailyUserNotificationEventQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to update all for stat_daily_user_notification_event")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to retrieve rows affected for stat_daily_user_notification_event")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o StatDailyUserNotificationEventSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), statDailyUserNotificationEventPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `stat_daily_user_notification_event` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, statDailyUserNotificationEventPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to update all in statDailyUserNotificationEvent slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to retrieve rows affected all in update all statDailyUserNotificationEvent")
	}
	return rowsAff, nil
}

var mySQLStatDailyUserNotificationEventUniqueColumns = []string{}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *StatDailyUserNotificationEvent) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model2: no stat_daily_user_notification_event provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(statDailyUserNotificationEventColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLStatDailyUserNotificationEventUniqueColumns, o)

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

	statDailyUserNotificationEventUpsertCacheMut.RLock()
	cache, cached := statDailyUserNotificationEventUpsertCache[key]
	statDailyUserNotificationEventUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			statDailyUserNotificationEventAllColumns,
			statDailyUserNotificationEventColumnsWithDefault,
			statDailyUserNotificationEventColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			statDailyUserNotificationEventAllColumns,
			statDailyUserNotificationEventPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("model2: unable to upsert stat_daily_user_notification_event, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`stat_daily_user_notification_event`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `stat_daily_user_notification_event` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(statDailyUserNotificationEventType, statDailyUserNotificationEventMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(statDailyUserNotificationEventType, statDailyUserNotificationEventMapping, ret)
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
		return errors.Wrap(err, "model2: unable to upsert for stat_daily_user_notification_event")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(statDailyUserNotificationEventType, statDailyUserNotificationEventMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "model2: unable to retrieve unique values for stat_daily_user_notification_event")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "model2: unable to populate default values for stat_daily_user_notification_event")
	}

CacheNoHooks:
	if !cached {
		statDailyUserNotificationEventUpsertCacheMut.Lock()
		statDailyUserNotificationEventUpsertCache[key] = cache
		statDailyUserNotificationEventUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single StatDailyUserNotificationEvent record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *StatDailyUserNotificationEvent) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("model2: no StatDailyUserNotificationEvent provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), statDailyUserNotificationEventPrimaryKeyMapping)
	sql := "DELETE FROM `stat_daily_user_notification_event` WHERE `date`=? AND `user_id`=? AND `event`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete from stat_daily_user_notification_event")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by delete for stat_daily_user_notification_event")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q statDailyUserNotificationEventQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model2: no statDailyUserNotificationEventQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete all from stat_daily_user_notification_event")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by deleteall for stat_daily_user_notification_event")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o StatDailyUserNotificationEventSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(statDailyUserNotificationEventBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), statDailyUserNotificationEventPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `stat_daily_user_notification_event` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, statDailyUserNotificationEventPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete all from statDailyUserNotificationEvent slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by deleteall for stat_daily_user_notification_event")
	}

	if len(statDailyUserNotificationEventAfterDeleteHooks) != 0 {
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
func (o *StatDailyUserNotificationEvent) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindStatDailyUserNotificationEvent(ctx, exec, o.Date, o.UserID, o.Event)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *StatDailyUserNotificationEventSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := StatDailyUserNotificationEventSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), statDailyUserNotificationEventPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `stat_daily_user_notification_event`.* FROM `stat_daily_user_notification_event` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, statDailyUserNotificationEventPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model2: unable to reload all in StatDailyUserNotificationEventSlice")
	}

	*o = slice

	return nil
}

// StatDailyUserNotificationEventExists checks if the StatDailyUserNotificationEvent row exists.
func StatDailyUserNotificationEventExists(ctx context.Context, exec boil.ContextExecutor, date time.Time, userID uint, event string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `stat_daily_user_notification_event` where `date`=? AND `user_id`=? AND `event`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, date, userID, event)
	}
	row := exec.QueryRowContext(ctx, sql, date, userID, event)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model2: unable to check if stat_daily_user_notification_event exists")
	}

	return exists, nil
}
