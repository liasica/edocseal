// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/liasica/edocseal/internal/ent/document"
	"github.com/liasica/edocseal/internal/ent/predicate"
	"github.com/liasica/edocseal/internal/model"
)

// DocumentUpdate is the builder for updating Document entities.
type DocumentUpdate struct {
	config
	hooks     []Hook
	mutation  *DocumentMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the DocumentUpdate builder.
func (du *DocumentUpdate) Where(ps ...predicate.Document) *DocumentUpdate {
	du.mutation.Where(ps...)
	return du
}

// SetHash sets the "hash" field.
func (du *DocumentUpdate) SetHash(s string) *DocumentUpdate {
	du.mutation.SetHash(s)
	return du
}

// SetNillableHash sets the "hash" field if the given value is not nil.
func (du *DocumentUpdate) SetNillableHash(s *string) *DocumentUpdate {
	if s != nil {
		du.SetHash(*s)
	}
	return du
}

// SetStatus sets the "status" field.
func (du *DocumentUpdate) SetStatus(d document.Status) *DocumentUpdate {
	du.mutation.SetStatus(d)
	return du
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (du *DocumentUpdate) SetNillableStatus(d *document.Status) *DocumentUpdate {
	if d != nil {
		du.SetStatus(*d)
	}
	return du
}

// SetTemplateID sets the "template_id" field.
func (du *DocumentUpdate) SetTemplateID(s string) *DocumentUpdate {
	du.mutation.SetTemplateID(s)
	return du
}

// SetNillableTemplateID sets the "template_id" field if the given value is not nil.
func (du *DocumentUpdate) SetNillableTemplateID(s *string) *DocumentUpdate {
	if s != nil {
		du.SetTemplateID(*s)
	}
	return du
}

// SetIDCardNumber sets the "id_card_number" field.
func (du *DocumentUpdate) SetIDCardNumber(s string) *DocumentUpdate {
	du.mutation.SetIDCardNumber(s)
	return du
}

// SetNillableIDCardNumber sets the "id_card_number" field if the given value is not nil.
func (du *DocumentUpdate) SetNillableIDCardNumber(s *string) *DocumentUpdate {
	if s != nil {
		du.SetIDCardNumber(*s)
	}
	return du
}

// SetExpiresAt sets the "expires_at" field.
func (du *DocumentUpdate) SetExpiresAt(t time.Time) *DocumentUpdate {
	du.mutation.SetExpiresAt(t)
	return du
}

// SetNillableExpiresAt sets the "expires_at" field if the given value is not nil.
func (du *DocumentUpdate) SetNillableExpiresAt(t *time.Time) *DocumentUpdate {
	if t != nil {
		du.SetExpiresAt(*t)
	}
	return du
}

// SetSignedURL sets the "signed_url" field.
func (du *DocumentUpdate) SetSignedURL(s string) *DocumentUpdate {
	du.mutation.SetSignedURL(s)
	return du
}

// SetNillableSignedURL sets the "signed_url" field if the given value is not nil.
func (du *DocumentUpdate) SetNillableSignedURL(s *string) *DocumentUpdate {
	if s != nil {
		du.SetSignedURL(*s)
	}
	return du
}

// ClearSignedURL clears the value of the "signed_url" field.
func (du *DocumentUpdate) ClearSignedURL() *DocumentUpdate {
	du.mutation.ClearSignedURL()
	return du
}

// SetUnsignedURL sets the "unsigned_url" field.
func (du *DocumentUpdate) SetUnsignedURL(s string) *DocumentUpdate {
	du.mutation.SetUnsignedURL(s)
	return du
}

// SetNillableUnsignedURL sets the "unsigned_url" field if the given value is not nil.
func (du *DocumentUpdate) SetNillableUnsignedURL(s *string) *DocumentUpdate {
	if s != nil {
		du.SetUnsignedURL(*s)
	}
	return du
}

// ClearUnsignedURL clears the value of the "unsigned_url" field.
func (du *DocumentUpdate) ClearUnsignedURL() *DocumentUpdate {
	du.mutation.ClearUnsignedURL()
	return du
}

// SetPaths sets the "paths" field.
func (du *DocumentUpdate) SetPaths(m *model.Paths) *DocumentUpdate {
	du.mutation.SetPaths(m)
	return du
}

// SetCreateAt sets the "create_at" field.
func (du *DocumentUpdate) SetCreateAt(t time.Time) *DocumentUpdate {
	du.mutation.SetCreateAt(t)
	return du
}

// SetNillableCreateAt sets the "create_at" field if the given value is not nil.
func (du *DocumentUpdate) SetNillableCreateAt(t *time.Time) *DocumentUpdate {
	if t != nil {
		du.SetCreateAt(*t)
	}
	return du
}

// Mutation returns the DocumentMutation object of the builder.
func (du *DocumentUpdate) Mutation() *DocumentMutation {
	return du.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (du *DocumentUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, du.sqlSave, du.mutation, du.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (du *DocumentUpdate) SaveX(ctx context.Context) int {
	affected, err := du.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (du *DocumentUpdate) Exec(ctx context.Context) error {
	_, err := du.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (du *DocumentUpdate) ExecX(ctx context.Context) {
	if err := du.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (du *DocumentUpdate) check() error {
	if v, ok := du.mutation.Status(); ok {
		if err := document.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Document.status": %w`, err)}
		}
	}
	if v, ok := du.mutation.IDCardNumber(); ok {
		if err := document.IDCardNumberValidator(v); err != nil {
			return &ValidationError{Name: "id_card_number", err: fmt.Errorf(`ent: validator failed for field "Document.id_card_number": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (du *DocumentUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *DocumentUpdate {
	du.modifiers = append(du.modifiers, modifiers...)
	return du
}

func (du *DocumentUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := du.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(document.Table, document.Columns, sqlgraph.NewFieldSpec(document.FieldID, field.TypeString))
	if ps := du.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := du.mutation.Hash(); ok {
		_spec.SetField(document.FieldHash, field.TypeString, value)
	}
	if value, ok := du.mutation.Status(); ok {
		_spec.SetField(document.FieldStatus, field.TypeEnum, value)
	}
	if value, ok := du.mutation.TemplateID(); ok {
		_spec.SetField(document.FieldTemplateID, field.TypeString, value)
	}
	if value, ok := du.mutation.IDCardNumber(); ok {
		_spec.SetField(document.FieldIDCardNumber, field.TypeString, value)
	}
	if value, ok := du.mutation.ExpiresAt(); ok {
		_spec.SetField(document.FieldExpiresAt, field.TypeTime, value)
	}
	if value, ok := du.mutation.SignedURL(); ok {
		_spec.SetField(document.FieldSignedURL, field.TypeString, value)
	}
	if du.mutation.SignedURLCleared() {
		_spec.ClearField(document.FieldSignedURL, field.TypeString)
	}
	if value, ok := du.mutation.UnsignedURL(); ok {
		_spec.SetField(document.FieldUnsignedURL, field.TypeString, value)
	}
	if du.mutation.UnsignedURLCleared() {
		_spec.ClearField(document.FieldUnsignedURL, field.TypeString)
	}
	if value, ok := du.mutation.Paths(); ok {
		_spec.SetField(document.FieldPaths, field.TypeJSON, value)
	}
	if value, ok := du.mutation.CreateAt(); ok {
		_spec.SetField(document.FieldCreateAt, field.TypeTime, value)
	}
	_spec.AddModifiers(du.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, du.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{document.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	du.mutation.done = true
	return n, nil
}

// DocumentUpdateOne is the builder for updating a single Document entity.
type DocumentUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *DocumentMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetHash sets the "hash" field.
func (duo *DocumentUpdateOne) SetHash(s string) *DocumentUpdateOne {
	duo.mutation.SetHash(s)
	return duo
}

// SetNillableHash sets the "hash" field if the given value is not nil.
func (duo *DocumentUpdateOne) SetNillableHash(s *string) *DocumentUpdateOne {
	if s != nil {
		duo.SetHash(*s)
	}
	return duo
}

// SetStatus sets the "status" field.
func (duo *DocumentUpdateOne) SetStatus(d document.Status) *DocumentUpdateOne {
	duo.mutation.SetStatus(d)
	return duo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (duo *DocumentUpdateOne) SetNillableStatus(d *document.Status) *DocumentUpdateOne {
	if d != nil {
		duo.SetStatus(*d)
	}
	return duo
}

// SetTemplateID sets the "template_id" field.
func (duo *DocumentUpdateOne) SetTemplateID(s string) *DocumentUpdateOne {
	duo.mutation.SetTemplateID(s)
	return duo
}

// SetNillableTemplateID sets the "template_id" field if the given value is not nil.
func (duo *DocumentUpdateOne) SetNillableTemplateID(s *string) *DocumentUpdateOne {
	if s != nil {
		duo.SetTemplateID(*s)
	}
	return duo
}

// SetIDCardNumber sets the "id_card_number" field.
func (duo *DocumentUpdateOne) SetIDCardNumber(s string) *DocumentUpdateOne {
	duo.mutation.SetIDCardNumber(s)
	return duo
}

// SetNillableIDCardNumber sets the "id_card_number" field if the given value is not nil.
func (duo *DocumentUpdateOne) SetNillableIDCardNumber(s *string) *DocumentUpdateOne {
	if s != nil {
		duo.SetIDCardNumber(*s)
	}
	return duo
}

// SetExpiresAt sets the "expires_at" field.
func (duo *DocumentUpdateOne) SetExpiresAt(t time.Time) *DocumentUpdateOne {
	duo.mutation.SetExpiresAt(t)
	return duo
}

// SetNillableExpiresAt sets the "expires_at" field if the given value is not nil.
func (duo *DocumentUpdateOne) SetNillableExpiresAt(t *time.Time) *DocumentUpdateOne {
	if t != nil {
		duo.SetExpiresAt(*t)
	}
	return duo
}

// SetSignedURL sets the "signed_url" field.
func (duo *DocumentUpdateOne) SetSignedURL(s string) *DocumentUpdateOne {
	duo.mutation.SetSignedURL(s)
	return duo
}

// SetNillableSignedURL sets the "signed_url" field if the given value is not nil.
func (duo *DocumentUpdateOne) SetNillableSignedURL(s *string) *DocumentUpdateOne {
	if s != nil {
		duo.SetSignedURL(*s)
	}
	return duo
}

// ClearSignedURL clears the value of the "signed_url" field.
func (duo *DocumentUpdateOne) ClearSignedURL() *DocumentUpdateOne {
	duo.mutation.ClearSignedURL()
	return duo
}

// SetUnsignedURL sets the "unsigned_url" field.
func (duo *DocumentUpdateOne) SetUnsignedURL(s string) *DocumentUpdateOne {
	duo.mutation.SetUnsignedURL(s)
	return duo
}

// SetNillableUnsignedURL sets the "unsigned_url" field if the given value is not nil.
func (duo *DocumentUpdateOne) SetNillableUnsignedURL(s *string) *DocumentUpdateOne {
	if s != nil {
		duo.SetUnsignedURL(*s)
	}
	return duo
}

// ClearUnsignedURL clears the value of the "unsigned_url" field.
func (duo *DocumentUpdateOne) ClearUnsignedURL() *DocumentUpdateOne {
	duo.mutation.ClearUnsignedURL()
	return duo
}

// SetPaths sets the "paths" field.
func (duo *DocumentUpdateOne) SetPaths(m *model.Paths) *DocumentUpdateOne {
	duo.mutation.SetPaths(m)
	return duo
}

// SetCreateAt sets the "create_at" field.
func (duo *DocumentUpdateOne) SetCreateAt(t time.Time) *DocumentUpdateOne {
	duo.mutation.SetCreateAt(t)
	return duo
}

// SetNillableCreateAt sets the "create_at" field if the given value is not nil.
func (duo *DocumentUpdateOne) SetNillableCreateAt(t *time.Time) *DocumentUpdateOne {
	if t != nil {
		duo.SetCreateAt(*t)
	}
	return duo
}

// Mutation returns the DocumentMutation object of the builder.
func (duo *DocumentUpdateOne) Mutation() *DocumentMutation {
	return duo.mutation
}

// Where appends a list predicates to the DocumentUpdate builder.
func (duo *DocumentUpdateOne) Where(ps ...predicate.Document) *DocumentUpdateOne {
	duo.mutation.Where(ps...)
	return duo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (duo *DocumentUpdateOne) Select(field string, fields ...string) *DocumentUpdateOne {
	duo.fields = append([]string{field}, fields...)
	return duo
}

// Save executes the query and returns the updated Document entity.
func (duo *DocumentUpdateOne) Save(ctx context.Context) (*Document, error) {
	return withHooks(ctx, duo.sqlSave, duo.mutation, duo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (duo *DocumentUpdateOne) SaveX(ctx context.Context) *Document {
	node, err := duo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (duo *DocumentUpdateOne) Exec(ctx context.Context) error {
	_, err := duo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (duo *DocumentUpdateOne) ExecX(ctx context.Context) {
	if err := duo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (duo *DocumentUpdateOne) check() error {
	if v, ok := duo.mutation.Status(); ok {
		if err := document.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Document.status": %w`, err)}
		}
	}
	if v, ok := duo.mutation.IDCardNumber(); ok {
		if err := document.IDCardNumberValidator(v); err != nil {
			return &ValidationError{Name: "id_card_number", err: fmt.Errorf(`ent: validator failed for field "Document.id_card_number": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (duo *DocumentUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *DocumentUpdateOne {
	duo.modifiers = append(duo.modifiers, modifiers...)
	return duo
}

func (duo *DocumentUpdateOne) sqlSave(ctx context.Context) (_node *Document, err error) {
	if err := duo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(document.Table, document.Columns, sqlgraph.NewFieldSpec(document.FieldID, field.TypeString))
	id, ok := duo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Document.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := duo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, document.FieldID)
		for _, f := range fields {
			if !document.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != document.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := duo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := duo.mutation.Hash(); ok {
		_spec.SetField(document.FieldHash, field.TypeString, value)
	}
	if value, ok := duo.mutation.Status(); ok {
		_spec.SetField(document.FieldStatus, field.TypeEnum, value)
	}
	if value, ok := duo.mutation.TemplateID(); ok {
		_spec.SetField(document.FieldTemplateID, field.TypeString, value)
	}
	if value, ok := duo.mutation.IDCardNumber(); ok {
		_spec.SetField(document.FieldIDCardNumber, field.TypeString, value)
	}
	if value, ok := duo.mutation.ExpiresAt(); ok {
		_spec.SetField(document.FieldExpiresAt, field.TypeTime, value)
	}
	if value, ok := duo.mutation.SignedURL(); ok {
		_spec.SetField(document.FieldSignedURL, field.TypeString, value)
	}
	if duo.mutation.SignedURLCleared() {
		_spec.ClearField(document.FieldSignedURL, field.TypeString)
	}
	if value, ok := duo.mutation.UnsignedURL(); ok {
		_spec.SetField(document.FieldUnsignedURL, field.TypeString, value)
	}
	if duo.mutation.UnsignedURLCleared() {
		_spec.ClearField(document.FieldUnsignedURL, field.TypeString)
	}
	if value, ok := duo.mutation.Paths(); ok {
		_spec.SetField(document.FieldPaths, field.TypeJSON, value)
	}
	if value, ok := duo.mutation.CreateAt(); ok {
		_spec.SetField(document.FieldCreateAt, field.TypeTime, value)
	}
	_spec.AddModifiers(duo.modifiers...)
	_node = &Document{config: duo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, duo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{document.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	duo.mutation.done = true
	return _node, nil
}
