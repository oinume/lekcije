// Code generated by SQLBoiler 4.9.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// MPlan is an object representing the database table.
type MPlan struct {
	ID                   uint8       `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name                 string      `boil:"name" json:"name" toml:"name" yaml:"name"`
	InternalName         string      `boil:"internal_name" json:"internal_name" toml:"internal_name" yaml:"internal_name"`
	StripeTestProductID  null.String `boil:"stripe_test_product_id" json:"stripe_test_product_id,omitempty" toml:"stripe_test_product_id" yaml:"stripe_test_product_id,omitempty"`
	Price                int32       `boil:"price" json:"price" toml:"price" yaml:"price"`
	NotificationInterval uint8       `boil:"notification_interval" json:"notification_interval" toml:"notification_interval" yaml:"notification_interval"`
	ShowAd               uint8       `boil:"show_ad" json:"show_ad" toml:"show_ad" yaml:"show_ad"`
	MaxFollowingTeacher  uint8       `boil:"max_following_teacher" json:"max_following_teacher" toml:"max_following_teacher" yaml:"max_following_teacher"`
	CreatedAt            time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt            time.Time   `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *mPlanR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L mPlanL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var MPlanColumns = struct {
	ID                   string
	Name                 string
	InternalName         string
	StripeTestProductID  string
	Price                string
	NotificationInterval string
	ShowAd               string
	MaxFollowingTeacher  string
	CreatedAt            string
	UpdatedAt            string
}{
	ID:                   "id",
	Name:                 "name",
	InternalName:         "internal_name",
	StripeTestProductID:  "stripe_test_product_id",
	Price:                "price",
	NotificationInterval: "notification_interval",
	ShowAd:               "show_ad",
	MaxFollowingTeacher:  "max_following_teacher",
	CreatedAt:            "created_at",
	UpdatedAt:            "updated_at",
}

var MPlanTableColumns = struct {
	ID                   string
	Name                 string
	InternalName         string
	StripeTestProductID  string
	Price                string
	NotificationInterval string
	ShowAd               string
	MaxFollowingTeacher  string
	CreatedAt            string
	UpdatedAt            string
}{
	ID:                   "m_plan.id",
	Name:                 "m_plan.name",
	InternalName:         "m_plan.internal_name",
	StripeTestProductID:  "m_plan.stripe_test_product_id",
	Price:                "m_plan.price",
	NotificationInterval: "m_plan.notification_interval",
	ShowAd:               "m_plan.show_ad",
	MaxFollowingTeacher:  "m_plan.max_following_teacher",
	CreatedAt:            "m_plan.created_at",
	UpdatedAt:            "m_plan.updated_at",
}

// Generated where

type whereHelperuint8 struct{ field string }

func (w whereHelperuint8) EQ(x uint8) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperuint8) NEQ(x uint8) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperuint8) LT(x uint8) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperuint8) LTE(x uint8) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperuint8) GT(x uint8) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperuint8) GTE(x uint8) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperuint8) IN(slice []uint8) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperuint8) NIN(slice []uint8) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelpernull_String struct{ field string }

func (w whereHelpernull_String) EQ(x null.String) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_String) NEQ(x null.String) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_String) LT(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_String) LTE(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_String) GT(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_String) GTE(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpernull_String) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_String) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

type whereHelperint32 struct{ field string }

func (w whereHelperint32) EQ(x int32) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint32) NEQ(x int32) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint32) LT(x int32) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint32) LTE(x int32) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint32) GT(x int32) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint32) GTE(x int32) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint32) IN(slice []int32) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint32) NIN(slice []int32) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var MPlanWhere = struct {
	ID                   whereHelperuint8
	Name                 whereHelperstring
	InternalName         whereHelperstring
	StripeTestProductID  whereHelpernull_String
	Price                whereHelperint32
	NotificationInterval whereHelperuint8
	ShowAd               whereHelperuint8
	MaxFollowingTeacher  whereHelperuint8
	CreatedAt            whereHelpertime_Time
	UpdatedAt            whereHelpertime_Time
}{
	ID:                   whereHelperuint8{field: "`m_plan`.`id`"},
	Name:                 whereHelperstring{field: "`m_plan`.`name`"},
	InternalName:         whereHelperstring{field: "`m_plan`.`internal_name`"},
	StripeTestProductID:  whereHelpernull_String{field: "`m_plan`.`stripe_test_product_id`"},
	Price:                whereHelperint32{field: "`m_plan`.`price`"},
	NotificationInterval: whereHelperuint8{field: "`m_plan`.`notification_interval`"},
	ShowAd:               whereHelperuint8{field: "`m_plan`.`show_ad`"},
	MaxFollowingTeacher:  whereHelperuint8{field: "`m_plan`.`max_following_teacher`"},
	CreatedAt:            whereHelpertime_Time{field: "`m_plan`.`created_at`"},
	UpdatedAt:            whereHelpertime_Time{field: "`m_plan`.`updated_at`"},
}

// MPlanRels is where relationship names are stored.
var MPlanRels = struct {
}{}

// mPlanR is where relationships are stored.
type mPlanR struct {
}

// NewStruct creates a new relationship struct
func (*mPlanR) NewStruct() *mPlanR {
	return &mPlanR{}
}

// mPlanL is where Load methods for each relationship are stored.
type mPlanL struct{}

var (
	mPlanAllColumns            = []string{"id", "name", "internal_name", "stripe_test_product_id", "price", "notification_interval", "show_ad", "max_following_teacher", "created_at", "updated_at"}
	mPlanColumnsWithoutDefault = []string{"id", "name", "internal_name", "stripe_test_product_id", "price", "notification_interval", "show_ad", "created_at", "updated_at"}
	mPlanColumnsWithDefault    = []string{"max_following_teacher"}
	mPlanPrimaryKeyColumns     = []string{"id"}
	mPlanGeneratedColumns      = []string{}
)

type (
	// MPlanSlice is an alias for a slice of pointers to MPlan.
	// This should almost always be used instead of []MPlan.
	MPlanSlice []*MPlan
	// MPlanHook is the signature for custom MPlan hook methods
	MPlanHook func(context.Context, boil.ContextExecutor, *MPlan) error

	mPlanQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	mPlanType                 = reflect.TypeOf(&MPlan{})
	mPlanMapping              = queries.MakeStructMapping(mPlanType)
	mPlanPrimaryKeyMapping, _ = queries.BindMapping(mPlanType, mPlanMapping, mPlanPrimaryKeyColumns)
	mPlanInsertCacheMut       sync.RWMutex
	mPlanInsertCache          = make(map[string]insertCache)
	mPlanUpdateCacheMut       sync.RWMutex
	mPlanUpdateCache          = make(map[string]updateCache)
	mPlanUpsertCacheMut       sync.RWMutex
	mPlanUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var mPlanAfterSelectHooks []MPlanHook

var mPlanBeforeInsertHooks []MPlanHook
var mPlanAfterInsertHooks []MPlanHook

var mPlanBeforeUpdateHooks []MPlanHook
var mPlanAfterUpdateHooks []MPlanHook

var mPlanBeforeDeleteHooks []MPlanHook
var mPlanAfterDeleteHooks []MPlanHook

var mPlanBeforeUpsertHooks []MPlanHook
var mPlanAfterUpsertHooks []MPlanHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *MPlan) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mPlanAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *MPlan) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mPlanBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *MPlan) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mPlanAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *MPlan) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mPlanBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *MPlan) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mPlanAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *MPlan) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mPlanBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *MPlan) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mPlanAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *MPlan) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mPlanBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *MPlan) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mPlanAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddMPlanHook registers your hook function for all future operations.
func AddMPlanHook(hookPoint boil.HookPoint, mPlanHook MPlanHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		mPlanAfterSelectHooks = append(mPlanAfterSelectHooks, mPlanHook)
	case boil.BeforeInsertHook:
		mPlanBeforeInsertHooks = append(mPlanBeforeInsertHooks, mPlanHook)
	case boil.AfterInsertHook:
		mPlanAfterInsertHooks = append(mPlanAfterInsertHooks, mPlanHook)
	case boil.BeforeUpdateHook:
		mPlanBeforeUpdateHooks = append(mPlanBeforeUpdateHooks, mPlanHook)
	case boil.AfterUpdateHook:
		mPlanAfterUpdateHooks = append(mPlanAfterUpdateHooks, mPlanHook)
	case boil.BeforeDeleteHook:
		mPlanBeforeDeleteHooks = append(mPlanBeforeDeleteHooks, mPlanHook)
	case boil.AfterDeleteHook:
		mPlanAfterDeleteHooks = append(mPlanAfterDeleteHooks, mPlanHook)
	case boil.BeforeUpsertHook:
		mPlanBeforeUpsertHooks = append(mPlanBeforeUpsertHooks, mPlanHook)
	case boil.AfterUpsertHook:
		mPlanAfterUpsertHooks = append(mPlanAfterUpsertHooks, mPlanHook)
	}
}

// One returns a single mPlan record from the query.
func (q mPlanQuery) One(ctx context.Context, exec boil.ContextExecutor) (*MPlan, error) {
	o := &MPlan{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model2: failed to execute a one query for m_plan")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all MPlan records from the query.
func (q mPlanQuery) All(ctx context.Context, exec boil.ContextExecutor) (MPlanSlice, error) {
	var o []*MPlan

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model2: failed to assign all query results to MPlan slice")
	}

	if len(mPlanAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all MPlan records in the query.
func (q mPlanQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to count m_plan rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q mPlanQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model2: failed to check if m_plan exists")
	}

	return count > 0, nil
}

// MPlans retrieves all the records using an executor.
func MPlans(mods ...qm.QueryMod) mPlanQuery {
	mods = append(mods, qm.From("`m_plan`"))
	return mPlanQuery{NewQuery(mods...)}
}

// FindMPlan retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindMPlan(ctx context.Context, exec boil.ContextExecutor, iD uint8, selectCols ...string) (*MPlan, error) {
	mPlanObj := &MPlan{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `m_plan` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, mPlanObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model2: unable to select from m_plan")
	}

	if err = mPlanObj.doAfterSelectHooks(ctx, exec); err != nil {
		return mPlanObj, err
	}

	return mPlanObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *MPlan) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model2: no m_plan provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(mPlanColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	mPlanInsertCacheMut.RLock()
	cache, cached := mPlanInsertCache[key]
	mPlanInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			mPlanAllColumns,
			mPlanColumnsWithDefault,
			mPlanColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(mPlanType, mPlanMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(mPlanType, mPlanMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `m_plan` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `m_plan` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `m_plan` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, mPlanPrimaryKeyColumns))
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
		return errors.Wrap(err, "model2: unable to insert into m_plan")
	}

	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
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
		return errors.Wrap(err, "model2: unable to populate default values for m_plan")
	}

CacheNoHooks:
	if !cached {
		mPlanInsertCacheMut.Lock()
		mPlanInsertCache[key] = cache
		mPlanInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the MPlan.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *MPlan) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	mPlanUpdateCacheMut.RLock()
	cache, cached := mPlanUpdateCache[key]
	mPlanUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			mPlanAllColumns,
			mPlanPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("model2: unable to update m_plan, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `m_plan` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, mPlanPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(mPlanType, mPlanMapping, append(wl, mPlanPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model2: unable to update m_plan row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by update for m_plan")
	}

	if !cached {
		mPlanUpdateCacheMut.Lock()
		mPlanUpdateCache[key] = cache
		mPlanUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q mPlanQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to update all for m_plan")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to retrieve rows affected for m_plan")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o MPlanSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mPlanPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `m_plan` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, mPlanPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to update all in mPlan slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to retrieve rows affected all in update all mPlan")
	}
	return rowsAff, nil
}

var mySQLMPlanUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *MPlan) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model2: no m_plan provided for upsert")
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

	nzDefaults := queries.NonZeroDefaultSet(mPlanColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLMPlanUniqueColumns, o)

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

	mPlanUpsertCacheMut.RLock()
	cache, cached := mPlanUpsertCache[key]
	mPlanUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			mPlanAllColumns,
			mPlanColumnsWithDefault,
			mPlanColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			mPlanAllColumns,
			mPlanPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("model2: unable to upsert m_plan, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`m_plan`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `m_plan` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(mPlanType, mPlanMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(mPlanType, mPlanMapping, ret)
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
		return errors.Wrap(err, "model2: unable to upsert for m_plan")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(mPlanType, mPlanMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "model2: unable to retrieve unique values for m_plan")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "model2: unable to populate default values for m_plan")
	}

CacheNoHooks:
	if !cached {
		mPlanUpsertCacheMut.Lock()
		mPlanUpsertCache[key] = cache
		mPlanUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single MPlan record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *MPlan) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("model2: no MPlan provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), mPlanPrimaryKeyMapping)
	sql := "DELETE FROM `m_plan` WHERE `id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete from m_plan")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by delete for m_plan")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q mPlanQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model2: no mPlanQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete all from m_plan")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by deleteall for m_plan")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o MPlanSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(mPlanBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mPlanPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `m_plan` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, mPlanPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete all from mPlan slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by deleteall for m_plan")
	}

	if len(mPlanAfterDeleteHooks) != 0 {
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
func (o *MPlan) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindMPlan(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *MPlanSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := MPlanSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mPlanPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `m_plan`.* FROM `m_plan` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, mPlanPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model2: unable to reload all in MPlanSlice")
	}

	*o = slice

	return nil
}

// MPlanExists checks if the MPlan row exists.
func MPlanExists(ctx context.Context, exec boil.ContextExecutor, iD uint8) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `m_plan` where `id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model2: unable to check if m_plan exists")
	}

	return exists, nil
}
