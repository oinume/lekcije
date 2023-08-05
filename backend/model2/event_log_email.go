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

// EventLogEmail is an object representing the database table.
type EventLogEmail struct {
	ID         uint      `boil:"id" json:"id" toml:"id" yaml:"id"`
	Datetime   time.Time `boil:"datetime" json:"datetime" toml:"datetime" yaml:"datetime"`
	Event      string    `boil:"event" json:"event" toml:"event" yaml:"event"`
	EmailType  string    `boil:"email_type" json:"email_type" toml:"email_type" yaml:"email_type"`
	UserID     uint      `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	UserAgent  string    `boil:"user_agent" json:"user_agent" toml:"user_agent" yaml:"user_agent"`
	TeacherIds string    `boil:"teacher_ids" json:"teacher_ids" toml:"teacher_ids" yaml:"teacher_ids"`
	URL        string    `boil:"url" json:"url" toml:"url" yaml:"url"`

	R *eventLogEmailR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L eventLogEmailL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var EventLogEmailColumns = struct {
	ID         string
	Datetime   string
	Event      string
	EmailType  string
	UserID     string
	UserAgent  string
	TeacherIds string
	URL        string
}{
	ID:         "id",
	Datetime:   "datetime",
	Event:      "event",
	EmailType:  "email_type",
	UserID:     "user_id",
	UserAgent:  "user_agent",
	TeacherIds: "teacher_ids",
	URL:        "url",
}

var EventLogEmailTableColumns = struct {
	ID         string
	Datetime   string
	Event      string
	EmailType  string
	UserID     string
	UserAgent  string
	TeacherIds string
	URL        string
}{
	ID:         "event_log_email.id",
	Datetime:   "event_log_email.datetime",
	Event:      "event_log_email.event",
	EmailType:  "event_log_email.email_type",
	UserID:     "event_log_email.user_id",
	UserAgent:  "event_log_email.user_agent",
	TeacherIds: "event_log_email.teacher_ids",
	URL:        "event_log_email.url",
}

// Generated where

type whereHelperuint struct{ field string }

func (w whereHelperuint) EQ(x uint) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperuint) NEQ(x uint) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperuint) LT(x uint) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperuint) LTE(x uint) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperuint) GT(x uint) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperuint) GTE(x uint) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperuint) IN(slice []uint) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperuint) NIN(slice []uint) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelpertime_Time struct{ field string }

func (w whereHelpertime_Time) EQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertime_Time) NEQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertime_Time) LT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertime_Time) LTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertime_Time) GT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertime_Time) GTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperstring) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperstring) NIN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var EventLogEmailWhere = struct {
	ID         whereHelperuint
	Datetime   whereHelpertime_Time
	Event      whereHelperstring
	EmailType  whereHelperstring
	UserID     whereHelperuint
	UserAgent  whereHelperstring
	TeacherIds whereHelperstring
	URL        whereHelperstring
}{
	ID:         whereHelperuint{field: "`event_log_email`.`id`"},
	Datetime:   whereHelpertime_Time{field: "`event_log_email`.`datetime`"},
	Event:      whereHelperstring{field: "`event_log_email`.`event`"},
	EmailType:  whereHelperstring{field: "`event_log_email`.`email_type`"},
	UserID:     whereHelperuint{field: "`event_log_email`.`user_id`"},
	UserAgent:  whereHelperstring{field: "`event_log_email`.`user_agent`"},
	TeacherIds: whereHelperstring{field: "`event_log_email`.`teacher_ids`"},
	URL:        whereHelperstring{field: "`event_log_email`.`url`"},
}

// EventLogEmailRels is where relationship names are stored.
var EventLogEmailRels = struct {
}{}

// eventLogEmailR is where relationships are stored.
type eventLogEmailR struct {
}

// NewStruct creates a new relationship struct
func (*eventLogEmailR) NewStruct() *eventLogEmailR {
	return &eventLogEmailR{}
}

// eventLogEmailL is where Load methods for each relationship are stored.
type eventLogEmailL struct{}

var (
	eventLogEmailAllColumns            = []string{"id", "datetime", "event", "email_type", "user_id", "user_agent", "teacher_ids", "url"}
	eventLogEmailColumnsWithoutDefault = []string{"datetime", "event", "email_type", "user_id", "user_agent", "teacher_ids", "url"}
	eventLogEmailColumnsWithDefault    = []string{"id"}
	eventLogEmailPrimaryKeyColumns     = []string{"id"}
	eventLogEmailGeneratedColumns      = []string{}
)

type (
	// EventLogEmailSlice is an alias for a slice of pointers to EventLogEmail.
	// This should almost always be used instead of []EventLogEmail.
	EventLogEmailSlice []*EventLogEmail
	// EventLogEmailHook is the signature for custom EventLogEmail hook methods
	EventLogEmailHook func(context.Context, boil.ContextExecutor, *EventLogEmail) error

	eventLogEmailQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	eventLogEmailType                 = reflect.TypeOf(&EventLogEmail{})
	eventLogEmailMapping              = queries.MakeStructMapping(eventLogEmailType)
	eventLogEmailPrimaryKeyMapping, _ = queries.BindMapping(eventLogEmailType, eventLogEmailMapping, eventLogEmailPrimaryKeyColumns)
	eventLogEmailInsertCacheMut       sync.RWMutex
	eventLogEmailInsertCache          = make(map[string]insertCache)
	eventLogEmailUpdateCacheMut       sync.RWMutex
	eventLogEmailUpdateCache          = make(map[string]updateCache)
	eventLogEmailUpsertCacheMut       sync.RWMutex
	eventLogEmailUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var eventLogEmailAfterSelectHooks []EventLogEmailHook

var eventLogEmailBeforeInsertHooks []EventLogEmailHook
var eventLogEmailAfterInsertHooks []EventLogEmailHook

var eventLogEmailBeforeUpdateHooks []EventLogEmailHook
var eventLogEmailAfterUpdateHooks []EventLogEmailHook

var eventLogEmailBeforeDeleteHooks []EventLogEmailHook
var eventLogEmailAfterDeleteHooks []EventLogEmailHook

var eventLogEmailBeforeUpsertHooks []EventLogEmailHook
var eventLogEmailAfterUpsertHooks []EventLogEmailHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *EventLogEmail) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range eventLogEmailAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *EventLogEmail) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range eventLogEmailBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *EventLogEmail) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range eventLogEmailAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *EventLogEmail) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range eventLogEmailBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *EventLogEmail) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range eventLogEmailAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *EventLogEmail) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range eventLogEmailBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *EventLogEmail) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range eventLogEmailAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *EventLogEmail) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range eventLogEmailBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *EventLogEmail) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range eventLogEmailAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddEventLogEmailHook registers your hook function for all future operations.
func AddEventLogEmailHook(hookPoint boil.HookPoint, eventLogEmailHook EventLogEmailHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		eventLogEmailAfterSelectHooks = append(eventLogEmailAfterSelectHooks, eventLogEmailHook)
	case boil.BeforeInsertHook:
		eventLogEmailBeforeInsertHooks = append(eventLogEmailBeforeInsertHooks, eventLogEmailHook)
	case boil.AfterInsertHook:
		eventLogEmailAfterInsertHooks = append(eventLogEmailAfterInsertHooks, eventLogEmailHook)
	case boil.BeforeUpdateHook:
		eventLogEmailBeforeUpdateHooks = append(eventLogEmailBeforeUpdateHooks, eventLogEmailHook)
	case boil.AfterUpdateHook:
		eventLogEmailAfterUpdateHooks = append(eventLogEmailAfterUpdateHooks, eventLogEmailHook)
	case boil.BeforeDeleteHook:
		eventLogEmailBeforeDeleteHooks = append(eventLogEmailBeforeDeleteHooks, eventLogEmailHook)
	case boil.AfterDeleteHook:
		eventLogEmailAfterDeleteHooks = append(eventLogEmailAfterDeleteHooks, eventLogEmailHook)
	case boil.BeforeUpsertHook:
		eventLogEmailBeforeUpsertHooks = append(eventLogEmailBeforeUpsertHooks, eventLogEmailHook)
	case boil.AfterUpsertHook:
		eventLogEmailAfterUpsertHooks = append(eventLogEmailAfterUpsertHooks, eventLogEmailHook)
	}
}

// One returns a single eventLogEmail record from the query.
func (q eventLogEmailQuery) One(ctx context.Context, exec boil.ContextExecutor) (*EventLogEmail, error) {
	o := &EventLogEmail{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model2: failed to execute a one query for event_log_email")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all EventLogEmail records from the query.
func (q eventLogEmailQuery) All(ctx context.Context, exec boil.ContextExecutor) (EventLogEmailSlice, error) {
	var o []*EventLogEmail

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model2: failed to assign all query results to EventLogEmail slice")
	}

	if len(eventLogEmailAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all EventLogEmail records in the query.
func (q eventLogEmailQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to count event_log_email rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q eventLogEmailQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model2: failed to check if event_log_email exists")
	}

	return count > 0, nil
}

// EventLogEmails retrieves all the records using an executor.
func EventLogEmails(mods ...qm.QueryMod) eventLogEmailQuery {
	mods = append(mods, qm.From("`event_log_email`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`event_log_email`.*"})
	}

	return eventLogEmailQuery{q}
}

// FindEventLogEmail retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindEventLogEmail(ctx context.Context, exec boil.ContextExecutor, iD uint, selectCols ...string) (*EventLogEmail, error) {
	eventLogEmailObj := &EventLogEmail{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `event_log_email` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, eventLogEmailObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model2: unable to select from event_log_email")
	}

	if err = eventLogEmailObj.doAfterSelectHooks(ctx, exec); err != nil {
		return eventLogEmailObj, err
	}

	return eventLogEmailObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *EventLogEmail) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model2: no event_log_email provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(eventLogEmailColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	eventLogEmailInsertCacheMut.RLock()
	cache, cached := eventLogEmailInsertCache[key]
	eventLogEmailInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			eventLogEmailAllColumns,
			eventLogEmailColumnsWithDefault,
			eventLogEmailColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(eventLogEmailType, eventLogEmailMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(eventLogEmailType, eventLogEmailMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `event_log_email` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `event_log_email` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `event_log_email` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, eventLogEmailPrimaryKeyColumns))
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
		return errors.Wrap(err, "model2: unable to insert into event_log_email")
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

	o.ID = uint(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == eventLogEmailMapping["id"] {
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
		return errors.Wrap(err, "model2: unable to populate default values for event_log_email")
	}

CacheNoHooks:
	if !cached {
		eventLogEmailInsertCacheMut.Lock()
		eventLogEmailInsertCache[key] = cache
		eventLogEmailInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the EventLogEmail.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *EventLogEmail) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	eventLogEmailUpdateCacheMut.RLock()
	cache, cached := eventLogEmailUpdateCache[key]
	eventLogEmailUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			eventLogEmailAllColumns,
			eventLogEmailPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("model2: unable to update event_log_email, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `event_log_email` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, eventLogEmailPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(eventLogEmailType, eventLogEmailMapping, append(wl, eventLogEmailPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model2: unable to update event_log_email row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by update for event_log_email")
	}

	if !cached {
		eventLogEmailUpdateCacheMut.Lock()
		eventLogEmailUpdateCache[key] = cache
		eventLogEmailUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q eventLogEmailQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to update all for event_log_email")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to retrieve rows affected for event_log_email")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o EventLogEmailSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), eventLogEmailPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `event_log_email` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, eventLogEmailPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to update all in eventLogEmail slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to retrieve rows affected all in update all eventLogEmail")
	}
	return rowsAff, nil
}

var mySQLEventLogEmailUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *EventLogEmail) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model2: no event_log_email provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(eventLogEmailColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLEventLogEmailUniqueColumns, o)

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

	eventLogEmailUpsertCacheMut.RLock()
	cache, cached := eventLogEmailUpsertCache[key]
	eventLogEmailUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			eventLogEmailAllColumns,
			eventLogEmailColumnsWithDefault,
			eventLogEmailColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			eventLogEmailAllColumns,
			eventLogEmailPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("model2: unable to upsert event_log_email, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`event_log_email`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `event_log_email` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(eventLogEmailType, eventLogEmailMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(eventLogEmailType, eventLogEmailMapping, ret)
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
		return errors.Wrap(err, "model2: unable to upsert for event_log_email")
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

	o.ID = uint(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == eventLogEmailMapping["id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(eventLogEmailType, eventLogEmailMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "model2: unable to retrieve unique values for event_log_email")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "model2: unable to populate default values for event_log_email")
	}

CacheNoHooks:
	if !cached {
		eventLogEmailUpsertCacheMut.Lock()
		eventLogEmailUpsertCache[key] = cache
		eventLogEmailUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single EventLogEmail record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *EventLogEmail) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("model2: no EventLogEmail provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), eventLogEmailPrimaryKeyMapping)
	sql := "DELETE FROM `event_log_email` WHERE `id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete from event_log_email")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by delete for event_log_email")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q eventLogEmailQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model2: no eventLogEmailQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete all from event_log_email")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by deleteall for event_log_email")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o EventLogEmailSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(eventLogEmailBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), eventLogEmailPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `event_log_email` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, eventLogEmailPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete all from eventLogEmail slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by deleteall for event_log_email")
	}

	if len(eventLogEmailAfterDeleteHooks) != 0 {
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
func (o *EventLogEmail) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindEventLogEmail(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *EventLogEmailSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := EventLogEmailSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), eventLogEmailPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `event_log_email`.* FROM `event_log_email` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, eventLogEmailPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model2: unable to reload all in EventLogEmailSlice")
	}

	*o = slice

	return nil
}

// EventLogEmailExists checks if the EventLogEmail row exists.
func EventLogEmailExists(ctx context.Context, exec boil.ContextExecutor, iD uint) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `event_log_email` where `id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model2: unable to check if event_log_email exists")
	}

	return exists, nil
}

// Exists checks if the EventLogEmail row exists.
func (o *EventLogEmail) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return EventLogEmailExists(ctx, exec, o.ID)
}
