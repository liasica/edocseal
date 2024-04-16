// Code generated by ent, DO NOT EDIT.

package document

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/liasica/edocseal/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id string) predicate.Document {
	return predicate.Document(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.Document {
	return predicate.Document(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.Document {
	return predicate.Document(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.Document {
	return predicate.Document(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.Document {
	return predicate.Document(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.Document {
	return predicate.Document(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.Document {
	return predicate.Document(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.Document {
	return predicate.Document(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.Document {
	return predicate.Document(sql.FieldLTE(FieldID, id))
}

// IDEqualFold applies the EqualFold predicate on the ID field.
func IDEqualFold(id string) predicate.Document {
	return predicate.Document(sql.FieldEqualFold(FieldID, id))
}

// IDContainsFold applies the ContainsFold predicate on the ID field.
func IDContainsFold(id string) predicate.Document {
	return predicate.Document(sql.FieldContainsFold(FieldID, id))
}

// Hash applies equality check predicate on the "hash" field. It's identical to HashEQ.
func Hash(v string) predicate.Document {
	return predicate.Document(sql.FieldEQ(FieldHash, v))
}

// TemplateID applies equality check predicate on the "template_id" field. It's identical to TemplateIDEQ.
func TemplateID(v string) predicate.Document {
	return predicate.Document(sql.FieldEQ(FieldTemplateID, v))
}

// IDCardNumber applies equality check predicate on the "id_card_number" field. It's identical to IDCardNumberEQ.
func IDCardNumber(v string) predicate.Document {
	return predicate.Document(sql.FieldEQ(FieldIDCardNumber, v))
}

// ExpiresAt applies equality check predicate on the "expires_at" field. It's identical to ExpiresAtEQ.
func ExpiresAt(v time.Time) predicate.Document {
	return predicate.Document(sql.FieldEQ(FieldExpiresAt, v))
}

// SignedURL applies equality check predicate on the "signed_url" field. It's identical to SignedURLEQ.
func SignedURL(v string) predicate.Document {
	return predicate.Document(sql.FieldEQ(FieldSignedURL, v))
}

// UnsignedURL applies equality check predicate on the "unsigned_url" field. It's identical to UnsignedURLEQ.
func UnsignedURL(v string) predicate.Document {
	return predicate.Document(sql.FieldEQ(FieldUnsignedURL, v))
}

// HashEQ applies the EQ predicate on the "hash" field.
func HashEQ(v string) predicate.Document {
	return predicate.Document(sql.FieldEQ(FieldHash, v))
}

// HashNEQ applies the NEQ predicate on the "hash" field.
func HashNEQ(v string) predicate.Document {
	return predicate.Document(sql.FieldNEQ(FieldHash, v))
}

// HashIn applies the In predicate on the "hash" field.
func HashIn(vs ...string) predicate.Document {
	return predicate.Document(sql.FieldIn(FieldHash, vs...))
}

// HashNotIn applies the NotIn predicate on the "hash" field.
func HashNotIn(vs ...string) predicate.Document {
	return predicate.Document(sql.FieldNotIn(FieldHash, vs...))
}

// HashGT applies the GT predicate on the "hash" field.
func HashGT(v string) predicate.Document {
	return predicate.Document(sql.FieldGT(FieldHash, v))
}

// HashGTE applies the GTE predicate on the "hash" field.
func HashGTE(v string) predicate.Document {
	return predicate.Document(sql.FieldGTE(FieldHash, v))
}

// HashLT applies the LT predicate on the "hash" field.
func HashLT(v string) predicate.Document {
	return predicate.Document(sql.FieldLT(FieldHash, v))
}

// HashLTE applies the LTE predicate on the "hash" field.
func HashLTE(v string) predicate.Document {
	return predicate.Document(sql.FieldLTE(FieldHash, v))
}

// HashContains applies the Contains predicate on the "hash" field.
func HashContains(v string) predicate.Document {
	return predicate.Document(sql.FieldContains(FieldHash, v))
}

// HashHasPrefix applies the HasPrefix predicate on the "hash" field.
func HashHasPrefix(v string) predicate.Document {
	return predicate.Document(sql.FieldHasPrefix(FieldHash, v))
}

// HashHasSuffix applies the HasSuffix predicate on the "hash" field.
func HashHasSuffix(v string) predicate.Document {
	return predicate.Document(sql.FieldHasSuffix(FieldHash, v))
}

// HashEqualFold applies the EqualFold predicate on the "hash" field.
func HashEqualFold(v string) predicate.Document {
	return predicate.Document(sql.FieldEqualFold(FieldHash, v))
}

// HashContainsFold applies the ContainsFold predicate on the "hash" field.
func HashContainsFold(v string) predicate.Document {
	return predicate.Document(sql.FieldContainsFold(FieldHash, v))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v Status) predicate.Document {
	return predicate.Document(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v Status) predicate.Document {
	return predicate.Document(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...Status) predicate.Document {
	return predicate.Document(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...Status) predicate.Document {
	return predicate.Document(sql.FieldNotIn(FieldStatus, vs...))
}

// TemplateIDEQ applies the EQ predicate on the "template_id" field.
func TemplateIDEQ(v string) predicate.Document {
	return predicate.Document(sql.FieldEQ(FieldTemplateID, v))
}

// TemplateIDNEQ applies the NEQ predicate on the "template_id" field.
func TemplateIDNEQ(v string) predicate.Document {
	return predicate.Document(sql.FieldNEQ(FieldTemplateID, v))
}

// TemplateIDIn applies the In predicate on the "template_id" field.
func TemplateIDIn(vs ...string) predicate.Document {
	return predicate.Document(sql.FieldIn(FieldTemplateID, vs...))
}

// TemplateIDNotIn applies the NotIn predicate on the "template_id" field.
func TemplateIDNotIn(vs ...string) predicate.Document {
	return predicate.Document(sql.FieldNotIn(FieldTemplateID, vs...))
}

// TemplateIDGT applies the GT predicate on the "template_id" field.
func TemplateIDGT(v string) predicate.Document {
	return predicate.Document(sql.FieldGT(FieldTemplateID, v))
}

// TemplateIDGTE applies the GTE predicate on the "template_id" field.
func TemplateIDGTE(v string) predicate.Document {
	return predicate.Document(sql.FieldGTE(FieldTemplateID, v))
}

// TemplateIDLT applies the LT predicate on the "template_id" field.
func TemplateIDLT(v string) predicate.Document {
	return predicate.Document(sql.FieldLT(FieldTemplateID, v))
}

// TemplateIDLTE applies the LTE predicate on the "template_id" field.
func TemplateIDLTE(v string) predicate.Document {
	return predicate.Document(sql.FieldLTE(FieldTemplateID, v))
}

// TemplateIDContains applies the Contains predicate on the "template_id" field.
func TemplateIDContains(v string) predicate.Document {
	return predicate.Document(sql.FieldContains(FieldTemplateID, v))
}

// TemplateIDHasPrefix applies the HasPrefix predicate on the "template_id" field.
func TemplateIDHasPrefix(v string) predicate.Document {
	return predicate.Document(sql.FieldHasPrefix(FieldTemplateID, v))
}

// TemplateIDHasSuffix applies the HasSuffix predicate on the "template_id" field.
func TemplateIDHasSuffix(v string) predicate.Document {
	return predicate.Document(sql.FieldHasSuffix(FieldTemplateID, v))
}

// TemplateIDEqualFold applies the EqualFold predicate on the "template_id" field.
func TemplateIDEqualFold(v string) predicate.Document {
	return predicate.Document(sql.FieldEqualFold(FieldTemplateID, v))
}

// TemplateIDContainsFold applies the ContainsFold predicate on the "template_id" field.
func TemplateIDContainsFold(v string) predicate.Document {
	return predicate.Document(sql.FieldContainsFold(FieldTemplateID, v))
}

// IDCardNumberEQ applies the EQ predicate on the "id_card_number" field.
func IDCardNumberEQ(v string) predicate.Document {
	return predicate.Document(sql.FieldEQ(FieldIDCardNumber, v))
}

// IDCardNumberNEQ applies the NEQ predicate on the "id_card_number" field.
func IDCardNumberNEQ(v string) predicate.Document {
	return predicate.Document(sql.FieldNEQ(FieldIDCardNumber, v))
}

// IDCardNumberIn applies the In predicate on the "id_card_number" field.
func IDCardNumberIn(vs ...string) predicate.Document {
	return predicate.Document(sql.FieldIn(FieldIDCardNumber, vs...))
}

// IDCardNumberNotIn applies the NotIn predicate on the "id_card_number" field.
func IDCardNumberNotIn(vs ...string) predicate.Document {
	return predicate.Document(sql.FieldNotIn(FieldIDCardNumber, vs...))
}

// IDCardNumberGT applies the GT predicate on the "id_card_number" field.
func IDCardNumberGT(v string) predicate.Document {
	return predicate.Document(sql.FieldGT(FieldIDCardNumber, v))
}

// IDCardNumberGTE applies the GTE predicate on the "id_card_number" field.
func IDCardNumberGTE(v string) predicate.Document {
	return predicate.Document(sql.FieldGTE(FieldIDCardNumber, v))
}

// IDCardNumberLT applies the LT predicate on the "id_card_number" field.
func IDCardNumberLT(v string) predicate.Document {
	return predicate.Document(sql.FieldLT(FieldIDCardNumber, v))
}

// IDCardNumberLTE applies the LTE predicate on the "id_card_number" field.
func IDCardNumberLTE(v string) predicate.Document {
	return predicate.Document(sql.FieldLTE(FieldIDCardNumber, v))
}

// IDCardNumberContains applies the Contains predicate on the "id_card_number" field.
func IDCardNumberContains(v string) predicate.Document {
	return predicate.Document(sql.FieldContains(FieldIDCardNumber, v))
}

// IDCardNumberHasPrefix applies the HasPrefix predicate on the "id_card_number" field.
func IDCardNumberHasPrefix(v string) predicate.Document {
	return predicate.Document(sql.FieldHasPrefix(FieldIDCardNumber, v))
}

// IDCardNumberHasSuffix applies the HasSuffix predicate on the "id_card_number" field.
func IDCardNumberHasSuffix(v string) predicate.Document {
	return predicate.Document(sql.FieldHasSuffix(FieldIDCardNumber, v))
}

// IDCardNumberEqualFold applies the EqualFold predicate on the "id_card_number" field.
func IDCardNumberEqualFold(v string) predicate.Document {
	return predicate.Document(sql.FieldEqualFold(FieldIDCardNumber, v))
}

// IDCardNumberContainsFold applies the ContainsFold predicate on the "id_card_number" field.
func IDCardNumberContainsFold(v string) predicate.Document {
	return predicate.Document(sql.FieldContainsFold(FieldIDCardNumber, v))
}

// ExpiresAtEQ applies the EQ predicate on the "expires_at" field.
func ExpiresAtEQ(v time.Time) predicate.Document {
	return predicate.Document(sql.FieldEQ(FieldExpiresAt, v))
}

// ExpiresAtNEQ applies the NEQ predicate on the "expires_at" field.
func ExpiresAtNEQ(v time.Time) predicate.Document {
	return predicate.Document(sql.FieldNEQ(FieldExpiresAt, v))
}

// ExpiresAtIn applies the In predicate on the "expires_at" field.
func ExpiresAtIn(vs ...time.Time) predicate.Document {
	return predicate.Document(sql.FieldIn(FieldExpiresAt, vs...))
}

// ExpiresAtNotIn applies the NotIn predicate on the "expires_at" field.
func ExpiresAtNotIn(vs ...time.Time) predicate.Document {
	return predicate.Document(sql.FieldNotIn(FieldExpiresAt, vs...))
}

// ExpiresAtGT applies the GT predicate on the "expires_at" field.
func ExpiresAtGT(v time.Time) predicate.Document {
	return predicate.Document(sql.FieldGT(FieldExpiresAt, v))
}

// ExpiresAtGTE applies the GTE predicate on the "expires_at" field.
func ExpiresAtGTE(v time.Time) predicate.Document {
	return predicate.Document(sql.FieldGTE(FieldExpiresAt, v))
}

// ExpiresAtLT applies the LT predicate on the "expires_at" field.
func ExpiresAtLT(v time.Time) predicate.Document {
	return predicate.Document(sql.FieldLT(FieldExpiresAt, v))
}

// ExpiresAtLTE applies the LTE predicate on the "expires_at" field.
func ExpiresAtLTE(v time.Time) predicate.Document {
	return predicate.Document(sql.FieldLTE(FieldExpiresAt, v))
}

// SignedURLEQ applies the EQ predicate on the "signed_url" field.
func SignedURLEQ(v string) predicate.Document {
	return predicate.Document(sql.FieldEQ(FieldSignedURL, v))
}

// SignedURLNEQ applies the NEQ predicate on the "signed_url" field.
func SignedURLNEQ(v string) predicate.Document {
	return predicate.Document(sql.FieldNEQ(FieldSignedURL, v))
}

// SignedURLIn applies the In predicate on the "signed_url" field.
func SignedURLIn(vs ...string) predicate.Document {
	return predicate.Document(sql.FieldIn(FieldSignedURL, vs...))
}

// SignedURLNotIn applies the NotIn predicate on the "signed_url" field.
func SignedURLNotIn(vs ...string) predicate.Document {
	return predicate.Document(sql.FieldNotIn(FieldSignedURL, vs...))
}

// SignedURLGT applies the GT predicate on the "signed_url" field.
func SignedURLGT(v string) predicate.Document {
	return predicate.Document(sql.FieldGT(FieldSignedURL, v))
}

// SignedURLGTE applies the GTE predicate on the "signed_url" field.
func SignedURLGTE(v string) predicate.Document {
	return predicate.Document(sql.FieldGTE(FieldSignedURL, v))
}

// SignedURLLT applies the LT predicate on the "signed_url" field.
func SignedURLLT(v string) predicate.Document {
	return predicate.Document(sql.FieldLT(FieldSignedURL, v))
}

// SignedURLLTE applies the LTE predicate on the "signed_url" field.
func SignedURLLTE(v string) predicate.Document {
	return predicate.Document(sql.FieldLTE(FieldSignedURL, v))
}

// SignedURLContains applies the Contains predicate on the "signed_url" field.
func SignedURLContains(v string) predicate.Document {
	return predicate.Document(sql.FieldContains(FieldSignedURL, v))
}

// SignedURLHasPrefix applies the HasPrefix predicate on the "signed_url" field.
func SignedURLHasPrefix(v string) predicate.Document {
	return predicate.Document(sql.FieldHasPrefix(FieldSignedURL, v))
}

// SignedURLHasSuffix applies the HasSuffix predicate on the "signed_url" field.
func SignedURLHasSuffix(v string) predicate.Document {
	return predicate.Document(sql.FieldHasSuffix(FieldSignedURL, v))
}

// SignedURLIsNil applies the IsNil predicate on the "signed_url" field.
func SignedURLIsNil() predicate.Document {
	return predicate.Document(sql.FieldIsNull(FieldSignedURL))
}

// SignedURLNotNil applies the NotNil predicate on the "signed_url" field.
func SignedURLNotNil() predicate.Document {
	return predicate.Document(sql.FieldNotNull(FieldSignedURL))
}

// SignedURLEqualFold applies the EqualFold predicate on the "signed_url" field.
func SignedURLEqualFold(v string) predicate.Document {
	return predicate.Document(sql.FieldEqualFold(FieldSignedURL, v))
}

// SignedURLContainsFold applies the ContainsFold predicate on the "signed_url" field.
func SignedURLContainsFold(v string) predicate.Document {
	return predicate.Document(sql.FieldContainsFold(FieldSignedURL, v))
}

// UnsignedURLEQ applies the EQ predicate on the "unsigned_url" field.
func UnsignedURLEQ(v string) predicate.Document {
	return predicate.Document(sql.FieldEQ(FieldUnsignedURL, v))
}

// UnsignedURLNEQ applies the NEQ predicate on the "unsigned_url" field.
func UnsignedURLNEQ(v string) predicate.Document {
	return predicate.Document(sql.FieldNEQ(FieldUnsignedURL, v))
}

// UnsignedURLIn applies the In predicate on the "unsigned_url" field.
func UnsignedURLIn(vs ...string) predicate.Document {
	return predicate.Document(sql.FieldIn(FieldUnsignedURL, vs...))
}

// UnsignedURLNotIn applies the NotIn predicate on the "unsigned_url" field.
func UnsignedURLNotIn(vs ...string) predicate.Document {
	return predicate.Document(sql.FieldNotIn(FieldUnsignedURL, vs...))
}

// UnsignedURLGT applies the GT predicate on the "unsigned_url" field.
func UnsignedURLGT(v string) predicate.Document {
	return predicate.Document(sql.FieldGT(FieldUnsignedURL, v))
}

// UnsignedURLGTE applies the GTE predicate on the "unsigned_url" field.
func UnsignedURLGTE(v string) predicate.Document {
	return predicate.Document(sql.FieldGTE(FieldUnsignedURL, v))
}

// UnsignedURLLT applies the LT predicate on the "unsigned_url" field.
func UnsignedURLLT(v string) predicate.Document {
	return predicate.Document(sql.FieldLT(FieldUnsignedURL, v))
}

// UnsignedURLLTE applies the LTE predicate on the "unsigned_url" field.
func UnsignedURLLTE(v string) predicate.Document {
	return predicate.Document(sql.FieldLTE(FieldUnsignedURL, v))
}

// UnsignedURLContains applies the Contains predicate on the "unsigned_url" field.
func UnsignedURLContains(v string) predicate.Document {
	return predicate.Document(sql.FieldContains(FieldUnsignedURL, v))
}

// UnsignedURLHasPrefix applies the HasPrefix predicate on the "unsigned_url" field.
func UnsignedURLHasPrefix(v string) predicate.Document {
	return predicate.Document(sql.FieldHasPrefix(FieldUnsignedURL, v))
}

// UnsignedURLHasSuffix applies the HasSuffix predicate on the "unsigned_url" field.
func UnsignedURLHasSuffix(v string) predicate.Document {
	return predicate.Document(sql.FieldHasSuffix(FieldUnsignedURL, v))
}

// UnsignedURLIsNil applies the IsNil predicate on the "unsigned_url" field.
func UnsignedURLIsNil() predicate.Document {
	return predicate.Document(sql.FieldIsNull(FieldUnsignedURL))
}

// UnsignedURLNotNil applies the NotNil predicate on the "unsigned_url" field.
func UnsignedURLNotNil() predicate.Document {
	return predicate.Document(sql.FieldNotNull(FieldUnsignedURL))
}

// UnsignedURLEqualFold applies the EqualFold predicate on the "unsigned_url" field.
func UnsignedURLEqualFold(v string) predicate.Document {
	return predicate.Document(sql.FieldEqualFold(FieldUnsignedURL, v))
}

// UnsignedURLContainsFold applies the ContainsFold predicate on the "unsigned_url" field.
func UnsignedURLContainsFold(v string) predicate.Document {
	return predicate.Document(sql.FieldContainsFold(FieldUnsignedURL, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Document) predicate.Document {
	return predicate.Document(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Document) predicate.Document {
	return predicate.Document(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Document) predicate.Document {
	return predicate.Document(sql.NotPredicates(p))
}
