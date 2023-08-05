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
	"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/volatiletech/strmangle"
)

// Teacher is an object representing the database table.
type Teacher struct {
	ID                uint              `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name              string            `boil:"name" json:"name" toml:"name" yaml:"name"`
	CountryID         int16             `boil:"country_id" json:"country_id" toml:"country_id" yaml:"country_id"`
	Gender            string            `boil:"gender" json:"gender" toml:"gender" yaml:"gender"`
	Birthday          time.Time         `boil:"birthday" json:"birthday" toml:"birthday" yaml:"birthday"`
	YearsOfExperience int8              `boil:"years_of_experience" json:"years_of_experience" toml:"years_of_experience" yaml:"years_of_experience"`
	FavoriteCount     uint              `boil:"favorite_count" json:"favorite_count" toml:"favorite_count" yaml:"favorite_count"`
	ReviewCount       uint              `boil:"review_count" json:"review_count" toml:"review_count" yaml:"review_count"`
	Rating            types.NullDecimal `boil:"rating" json:"rating,omitempty" toml:"rating" yaml:"rating,omitempty"`
	LastLessonAt      time.Time         `boil:"last_lesson_at" json:"last_lesson_at" toml:"last_lesson_at" yaml:"last_lesson_at"`
	FetchErrorCount   uint8             `boil:"fetch_error_count" json:"fetch_error_count" toml:"fetch_error_count" yaml:"fetch_error_count"`
	CreatedAt         time.Time         `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt         time.Time         `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *teacherR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L teacherL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var TeacherColumns = struct {
	ID                string
	Name              string
	CountryID         string
	Gender            string
	Birthday          string
	YearsOfExperience string
	FavoriteCount     string
	ReviewCount       string
	Rating            string
	LastLessonAt      string
	FetchErrorCount   string
	CreatedAt         string
	UpdatedAt         string
}{
	ID:                "id",
	Name:              "name",
	CountryID:         "country_id",
	Gender:            "gender",
	Birthday:          "birthday",
	YearsOfExperience: "years_of_experience",
	FavoriteCount:     "favorite_count",
	ReviewCount:       "review_count",
	Rating:            "rating",
	LastLessonAt:      "last_lesson_at",
	FetchErrorCount:   "fetch_error_count",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
}

var TeacherTableColumns = struct {
	ID                string
	Name              string
	CountryID         string
	Gender            string
	Birthday          string
	YearsOfExperience string
	FavoriteCount     string
	ReviewCount       string
	Rating            string
	LastLessonAt      string
	FetchErrorCount   string
	CreatedAt         string
	UpdatedAt         string
}{
	ID:                "teacher.id",
	Name:              "teacher.name",
	CountryID:         "teacher.country_id",
	Gender:            "teacher.gender",
	Birthday:          "teacher.birthday",
	YearsOfExperience: "teacher.years_of_experience",
	FavoriteCount:     "teacher.favorite_count",
	ReviewCount:       "teacher.review_count",
	Rating:            "teacher.rating",
	LastLessonAt:      "teacher.last_lesson_at",
	FetchErrorCount:   "teacher.fetch_error_count",
	CreatedAt:         "teacher.created_at",
	UpdatedAt:         "teacher.updated_at",
}

// Generated where

type whereHelperint16 struct{ field string }

func (w whereHelperint16) EQ(x int16) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint16) NEQ(x int16) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint16) LT(x int16) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint16) LTE(x int16) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint16) GT(x int16) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint16) GTE(x int16) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint16) IN(slice []int16) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint16) NIN(slice []int16) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelperint8 struct{ field string }

func (w whereHelperint8) EQ(x int8) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint8) NEQ(x int8) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint8) LT(x int8) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint8) LTE(x int8) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint8) GT(x int8) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint8) GTE(x int8) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint8) IN(slice []int8) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint8) NIN(slice []int8) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelpertypes_NullDecimal struct{ field string }

func (w whereHelpertypes_NullDecimal) EQ(x types.NullDecimal) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpertypes_NullDecimal) NEQ(x types.NullDecimal) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpertypes_NullDecimal) LT(x types.NullDecimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertypes_NullDecimal) LTE(x types.NullDecimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertypes_NullDecimal) GT(x types.NullDecimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertypes_NullDecimal) GTE(x types.NullDecimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpertypes_NullDecimal) IsNull() qm.QueryMod { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpertypes_NullDecimal) IsNotNull() qm.QueryMod {
	return qmhelper.WhereIsNotNull(w.field)
}

var TeacherWhere = struct {
	ID                whereHelperuint
	Name              whereHelperstring
	CountryID         whereHelperint16
	Gender            whereHelperstring
	Birthday          whereHelpertime_Time
	YearsOfExperience whereHelperint8
	FavoriteCount     whereHelperuint
	ReviewCount       whereHelperuint
	Rating            whereHelpertypes_NullDecimal
	LastLessonAt      whereHelpertime_Time
	FetchErrorCount   whereHelperuint8
	CreatedAt         whereHelpertime_Time
	UpdatedAt         whereHelpertime_Time
}{
	ID:                whereHelperuint{field: "`teacher`.`id`"},
	Name:              whereHelperstring{field: "`teacher`.`name`"},
	CountryID:         whereHelperint16{field: "`teacher`.`country_id`"},
	Gender:            whereHelperstring{field: "`teacher`.`gender`"},
	Birthday:          whereHelpertime_Time{field: "`teacher`.`birthday`"},
	YearsOfExperience: whereHelperint8{field: "`teacher`.`years_of_experience`"},
	FavoriteCount:     whereHelperuint{field: "`teacher`.`favorite_count`"},
	ReviewCount:       whereHelperuint{field: "`teacher`.`review_count`"},
	Rating:            whereHelpertypes_NullDecimal{field: "`teacher`.`rating`"},
	LastLessonAt:      whereHelpertime_Time{field: "`teacher`.`last_lesson_at`"},
	FetchErrorCount:   whereHelperuint8{field: "`teacher`.`fetch_error_count`"},
	CreatedAt:         whereHelpertime_Time{field: "`teacher`.`created_at`"},
	UpdatedAt:         whereHelpertime_Time{field: "`teacher`.`updated_at`"},
}

// TeacherRels is where relationship names are stored.
var TeacherRels = struct {
}{}

// teacherR is where relationships are stored.
type teacherR struct {
}

// NewStruct creates a new relationship struct
func (*teacherR) NewStruct() *teacherR {
	return &teacherR{}
}

// teacherL is where Load methods for each relationship are stored.
type teacherL struct{}

var (
	teacherAllColumns            = []string{"id", "name", "country_id", "gender", "birthday", "years_of_experience", "favorite_count", "review_count", "rating", "last_lesson_at", "fetch_error_count", "created_at", "updated_at"}
	teacherColumnsWithoutDefault = []string{"id", "name", "birthday", "last_lesson_at", "created_at", "updated_at"}
	teacherColumnsWithDefault    = []string{"country_id", "gender", "years_of_experience", "favorite_count", "review_count", "rating", "fetch_error_count"}
	teacherPrimaryKeyColumns     = []string{"id"}
	teacherGeneratedColumns      = []string{}
)

type (
	// TeacherSlice is an alias for a slice of pointers to Teacher.
	// This should almost always be used instead of []Teacher.
	TeacherSlice []*Teacher
	// TeacherHook is the signature for custom Teacher hook methods
	TeacherHook func(context.Context, boil.ContextExecutor, *Teacher) error

	teacherQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	teacherType                 = reflect.TypeOf(&Teacher{})
	teacherMapping              = queries.MakeStructMapping(teacherType)
	teacherPrimaryKeyMapping, _ = queries.BindMapping(teacherType, teacherMapping, teacherPrimaryKeyColumns)
	teacherInsertCacheMut       sync.RWMutex
	teacherInsertCache          = make(map[string]insertCache)
	teacherUpdateCacheMut       sync.RWMutex
	teacherUpdateCache          = make(map[string]updateCache)
	teacherUpsertCacheMut       sync.RWMutex
	teacherUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var teacherAfterSelectHooks []TeacherHook

var teacherBeforeInsertHooks []TeacherHook
var teacherAfterInsertHooks []TeacherHook

var teacherBeforeUpdateHooks []TeacherHook
var teacherAfterUpdateHooks []TeacherHook

var teacherBeforeDeleteHooks []TeacherHook
var teacherAfterDeleteHooks []TeacherHook

var teacherBeforeUpsertHooks []TeacherHook
var teacherAfterUpsertHooks []TeacherHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Teacher) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teacherAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Teacher) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teacherBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Teacher) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teacherAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Teacher) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teacherBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Teacher) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teacherAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Teacher) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teacherBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Teacher) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teacherAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Teacher) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teacherBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Teacher) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teacherAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddTeacherHook registers your hook function for all future operations.
func AddTeacherHook(hookPoint boil.HookPoint, teacherHook TeacherHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		teacherAfterSelectHooks = append(teacherAfterSelectHooks, teacherHook)
	case boil.BeforeInsertHook:
		teacherBeforeInsertHooks = append(teacherBeforeInsertHooks, teacherHook)
	case boil.AfterInsertHook:
		teacherAfterInsertHooks = append(teacherAfterInsertHooks, teacherHook)
	case boil.BeforeUpdateHook:
		teacherBeforeUpdateHooks = append(teacherBeforeUpdateHooks, teacherHook)
	case boil.AfterUpdateHook:
		teacherAfterUpdateHooks = append(teacherAfterUpdateHooks, teacherHook)
	case boil.BeforeDeleteHook:
		teacherBeforeDeleteHooks = append(teacherBeforeDeleteHooks, teacherHook)
	case boil.AfterDeleteHook:
		teacherAfterDeleteHooks = append(teacherAfterDeleteHooks, teacherHook)
	case boil.BeforeUpsertHook:
		teacherBeforeUpsertHooks = append(teacherBeforeUpsertHooks, teacherHook)
	case boil.AfterUpsertHook:
		teacherAfterUpsertHooks = append(teacherAfterUpsertHooks, teacherHook)
	}
}

// One returns a single teacher record from the query.
func (q teacherQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Teacher, error) {
	o := &Teacher{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model2: failed to execute a one query for teacher")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Teacher records from the query.
func (q teacherQuery) All(ctx context.Context, exec boil.ContextExecutor) (TeacherSlice, error) {
	var o []*Teacher

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model2: failed to assign all query results to Teacher slice")
	}

	if len(teacherAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Teacher records in the query.
func (q teacherQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to count teacher rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q teacherQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model2: failed to check if teacher exists")
	}

	return count > 0, nil
}

// Teachers retrieves all the records using an executor.
func Teachers(mods ...qm.QueryMod) teacherQuery {
	mods = append(mods, qm.From("`teacher`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`teacher`.*"})
	}

	return teacherQuery{q}
}

// FindTeacher retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindTeacher(ctx context.Context, exec boil.ContextExecutor, iD uint, selectCols ...string) (*Teacher, error) {
	teacherObj := &Teacher{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `teacher` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, teacherObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model2: unable to select from teacher")
	}

	if err = teacherObj.doAfterSelectHooks(ctx, exec); err != nil {
		return teacherObj, err
	}

	return teacherObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Teacher) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model2: no teacher provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(teacherColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	teacherInsertCacheMut.RLock()
	cache, cached := teacherInsertCache[key]
	teacherInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			teacherAllColumns,
			teacherColumnsWithDefault,
			teacherColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(teacherType, teacherMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(teacherType, teacherMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `teacher` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `teacher` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `teacher` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, teacherPrimaryKeyColumns))
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
		return errors.Wrap(err, "model2: unable to insert into teacher")
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
		return errors.Wrap(err, "model2: unable to populate default values for teacher")
	}

CacheNoHooks:
	if !cached {
		teacherInsertCacheMut.Lock()
		teacherInsertCache[key] = cache
		teacherInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Teacher.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Teacher) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	teacherUpdateCacheMut.RLock()
	cache, cached := teacherUpdateCache[key]
	teacherUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			teacherAllColumns,
			teacherPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("model2: unable to update teacher, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `teacher` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, teacherPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(teacherType, teacherMapping, append(wl, teacherPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model2: unable to update teacher row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by update for teacher")
	}

	if !cached {
		teacherUpdateCacheMut.Lock()
		teacherUpdateCache[key] = cache
		teacherUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q teacherQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to update all for teacher")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to retrieve rows affected for teacher")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o TeacherSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), teacherPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `teacher` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, teacherPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to update all in teacher slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to retrieve rows affected all in update all teacher")
	}
	return rowsAff, nil
}

var mySQLTeacherUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Teacher) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model2: no teacher provided for upsert")
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

	nzDefaults := queries.NonZeroDefaultSet(teacherColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLTeacherUniqueColumns, o)

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

	teacherUpsertCacheMut.RLock()
	cache, cached := teacherUpsertCache[key]
	teacherUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			teacherAllColumns,
			teacherColumnsWithDefault,
			teacherColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			teacherAllColumns,
			teacherPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("model2: unable to upsert teacher, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`teacher`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `teacher` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(teacherType, teacherMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(teacherType, teacherMapping, ret)
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
		return errors.Wrap(err, "model2: unable to upsert for teacher")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(teacherType, teacherMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "model2: unable to retrieve unique values for teacher")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "model2: unable to populate default values for teacher")
	}

CacheNoHooks:
	if !cached {
		teacherUpsertCacheMut.Lock()
		teacherUpsertCache[key] = cache
		teacherUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Teacher record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Teacher) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("model2: no Teacher provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), teacherPrimaryKeyMapping)
	sql := "DELETE FROM `teacher` WHERE `id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete from teacher")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by delete for teacher")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q teacherQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model2: no teacherQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete all from teacher")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by deleteall for teacher")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o TeacherSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(teacherBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), teacherPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `teacher` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, teacherPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model2: unable to delete all from teacher slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model2: failed to get rows affected by deleteall for teacher")
	}

	if len(teacherAfterDeleteHooks) != 0 {
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
func (o *Teacher) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindTeacher(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *TeacherSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := TeacherSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), teacherPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `teacher`.* FROM `teacher` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, teacherPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model2: unable to reload all in TeacherSlice")
	}

	*o = slice

	return nil
}

// TeacherExists checks if the Teacher row exists.
func TeacherExists(ctx context.Context, exec boil.ContextExecutor, iD uint) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `teacher` where `id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model2: unable to check if teacher exists")
	}

	return exists, nil
}

// Exists checks if the Teacher row exists.
func (o *Teacher) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return TeacherExists(ctx, exec, o.ID)
}
