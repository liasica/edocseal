// Code generated by ent, DO NOT EDIT.

package certification

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"auroraride.com/edocseal/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Certification {
	return predicate.Certification(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Certification {
	return predicate.Certification(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Certification {
	return predicate.Certification(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Certification {
	return predicate.Certification(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Certification {
	return predicate.Certification(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Certification {
	return predicate.Certification(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Certification {
	return predicate.Certification(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Certification {
	return predicate.Certification(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Certification {
	return predicate.Certification(sql.FieldLTE(FieldID, id))
}

// IDCardNumber applies equality check predicate on the "id_card_number" field. It's identical to IDCardNumberEQ.
func IDCardNumber(v string) predicate.Certification {
	return predicate.Certification(sql.FieldEQ(FieldIDCardNumber, v))
}

// PrivatePath applies equality check predicate on the "private_path" field. It's identical to PrivatePathEQ.
func PrivatePath(v string) predicate.Certification {
	return predicate.Certification(sql.FieldEQ(FieldPrivatePath, v))
}

// CertPath applies equality check predicate on the "cert_path" field. It's identical to CertPathEQ.
func CertPath(v string) predicate.Certification {
	return predicate.Certification(sql.FieldEQ(FieldCertPath, v))
}

// ExpiresAt applies equality check predicate on the "expires_at" field. It's identical to ExpiresAtEQ.
func ExpiresAt(v time.Time) predicate.Certification {
	return predicate.Certification(sql.FieldEQ(FieldExpiresAt, v))
}

// IDCardNumberEQ applies the EQ predicate on the "id_card_number" field.
func IDCardNumberEQ(v string) predicate.Certification {
	return predicate.Certification(sql.FieldEQ(FieldIDCardNumber, v))
}

// IDCardNumberNEQ applies the NEQ predicate on the "id_card_number" field.
func IDCardNumberNEQ(v string) predicate.Certification {
	return predicate.Certification(sql.FieldNEQ(FieldIDCardNumber, v))
}

// IDCardNumberIn applies the In predicate on the "id_card_number" field.
func IDCardNumberIn(vs ...string) predicate.Certification {
	return predicate.Certification(sql.FieldIn(FieldIDCardNumber, vs...))
}

// IDCardNumberNotIn applies the NotIn predicate on the "id_card_number" field.
func IDCardNumberNotIn(vs ...string) predicate.Certification {
	return predicate.Certification(sql.FieldNotIn(FieldIDCardNumber, vs...))
}

// IDCardNumberGT applies the GT predicate on the "id_card_number" field.
func IDCardNumberGT(v string) predicate.Certification {
	return predicate.Certification(sql.FieldGT(FieldIDCardNumber, v))
}

// IDCardNumberGTE applies the GTE predicate on the "id_card_number" field.
func IDCardNumberGTE(v string) predicate.Certification {
	return predicate.Certification(sql.FieldGTE(FieldIDCardNumber, v))
}

// IDCardNumberLT applies the LT predicate on the "id_card_number" field.
func IDCardNumberLT(v string) predicate.Certification {
	return predicate.Certification(sql.FieldLT(FieldIDCardNumber, v))
}

// IDCardNumberLTE applies the LTE predicate on the "id_card_number" field.
func IDCardNumberLTE(v string) predicate.Certification {
	return predicate.Certification(sql.FieldLTE(FieldIDCardNumber, v))
}

// IDCardNumberContains applies the Contains predicate on the "id_card_number" field.
func IDCardNumberContains(v string) predicate.Certification {
	return predicate.Certification(sql.FieldContains(FieldIDCardNumber, v))
}

// IDCardNumberHasPrefix applies the HasPrefix predicate on the "id_card_number" field.
func IDCardNumberHasPrefix(v string) predicate.Certification {
	return predicate.Certification(sql.FieldHasPrefix(FieldIDCardNumber, v))
}

// IDCardNumberHasSuffix applies the HasSuffix predicate on the "id_card_number" field.
func IDCardNumberHasSuffix(v string) predicate.Certification {
	return predicate.Certification(sql.FieldHasSuffix(FieldIDCardNumber, v))
}

// IDCardNumberEqualFold applies the EqualFold predicate on the "id_card_number" field.
func IDCardNumberEqualFold(v string) predicate.Certification {
	return predicate.Certification(sql.FieldEqualFold(FieldIDCardNumber, v))
}

// IDCardNumberContainsFold applies the ContainsFold predicate on the "id_card_number" field.
func IDCardNumberContainsFold(v string) predicate.Certification {
	return predicate.Certification(sql.FieldContainsFold(FieldIDCardNumber, v))
}

// PrivatePathEQ applies the EQ predicate on the "private_path" field.
func PrivatePathEQ(v string) predicate.Certification {
	return predicate.Certification(sql.FieldEQ(FieldPrivatePath, v))
}

// PrivatePathNEQ applies the NEQ predicate on the "private_path" field.
func PrivatePathNEQ(v string) predicate.Certification {
	return predicate.Certification(sql.FieldNEQ(FieldPrivatePath, v))
}

// PrivatePathIn applies the In predicate on the "private_path" field.
func PrivatePathIn(vs ...string) predicate.Certification {
	return predicate.Certification(sql.FieldIn(FieldPrivatePath, vs...))
}

// PrivatePathNotIn applies the NotIn predicate on the "private_path" field.
func PrivatePathNotIn(vs ...string) predicate.Certification {
	return predicate.Certification(sql.FieldNotIn(FieldPrivatePath, vs...))
}

// PrivatePathGT applies the GT predicate on the "private_path" field.
func PrivatePathGT(v string) predicate.Certification {
	return predicate.Certification(sql.FieldGT(FieldPrivatePath, v))
}

// PrivatePathGTE applies the GTE predicate on the "private_path" field.
func PrivatePathGTE(v string) predicate.Certification {
	return predicate.Certification(sql.FieldGTE(FieldPrivatePath, v))
}

// PrivatePathLT applies the LT predicate on the "private_path" field.
func PrivatePathLT(v string) predicate.Certification {
	return predicate.Certification(sql.FieldLT(FieldPrivatePath, v))
}

// PrivatePathLTE applies the LTE predicate on the "private_path" field.
func PrivatePathLTE(v string) predicate.Certification {
	return predicate.Certification(sql.FieldLTE(FieldPrivatePath, v))
}

// PrivatePathContains applies the Contains predicate on the "private_path" field.
func PrivatePathContains(v string) predicate.Certification {
	return predicate.Certification(sql.FieldContains(FieldPrivatePath, v))
}

// PrivatePathHasPrefix applies the HasPrefix predicate on the "private_path" field.
func PrivatePathHasPrefix(v string) predicate.Certification {
	return predicate.Certification(sql.FieldHasPrefix(FieldPrivatePath, v))
}

// PrivatePathHasSuffix applies the HasSuffix predicate on the "private_path" field.
func PrivatePathHasSuffix(v string) predicate.Certification {
	return predicate.Certification(sql.FieldHasSuffix(FieldPrivatePath, v))
}

// PrivatePathEqualFold applies the EqualFold predicate on the "private_path" field.
func PrivatePathEqualFold(v string) predicate.Certification {
	return predicate.Certification(sql.FieldEqualFold(FieldPrivatePath, v))
}

// PrivatePathContainsFold applies the ContainsFold predicate on the "private_path" field.
func PrivatePathContainsFold(v string) predicate.Certification {
	return predicate.Certification(sql.FieldContainsFold(FieldPrivatePath, v))
}

// CertPathEQ applies the EQ predicate on the "cert_path" field.
func CertPathEQ(v string) predicate.Certification {
	return predicate.Certification(sql.FieldEQ(FieldCertPath, v))
}

// CertPathNEQ applies the NEQ predicate on the "cert_path" field.
func CertPathNEQ(v string) predicate.Certification {
	return predicate.Certification(sql.FieldNEQ(FieldCertPath, v))
}

// CertPathIn applies the In predicate on the "cert_path" field.
func CertPathIn(vs ...string) predicate.Certification {
	return predicate.Certification(sql.FieldIn(FieldCertPath, vs...))
}

// CertPathNotIn applies the NotIn predicate on the "cert_path" field.
func CertPathNotIn(vs ...string) predicate.Certification {
	return predicate.Certification(sql.FieldNotIn(FieldCertPath, vs...))
}

// CertPathGT applies the GT predicate on the "cert_path" field.
func CertPathGT(v string) predicate.Certification {
	return predicate.Certification(sql.FieldGT(FieldCertPath, v))
}

// CertPathGTE applies the GTE predicate on the "cert_path" field.
func CertPathGTE(v string) predicate.Certification {
	return predicate.Certification(sql.FieldGTE(FieldCertPath, v))
}

// CertPathLT applies the LT predicate on the "cert_path" field.
func CertPathLT(v string) predicate.Certification {
	return predicate.Certification(sql.FieldLT(FieldCertPath, v))
}

// CertPathLTE applies the LTE predicate on the "cert_path" field.
func CertPathLTE(v string) predicate.Certification {
	return predicate.Certification(sql.FieldLTE(FieldCertPath, v))
}

// CertPathContains applies the Contains predicate on the "cert_path" field.
func CertPathContains(v string) predicate.Certification {
	return predicate.Certification(sql.FieldContains(FieldCertPath, v))
}

// CertPathHasPrefix applies the HasPrefix predicate on the "cert_path" field.
func CertPathHasPrefix(v string) predicate.Certification {
	return predicate.Certification(sql.FieldHasPrefix(FieldCertPath, v))
}

// CertPathHasSuffix applies the HasSuffix predicate on the "cert_path" field.
func CertPathHasSuffix(v string) predicate.Certification {
	return predicate.Certification(sql.FieldHasSuffix(FieldCertPath, v))
}

// CertPathEqualFold applies the EqualFold predicate on the "cert_path" field.
func CertPathEqualFold(v string) predicate.Certification {
	return predicate.Certification(sql.FieldEqualFold(FieldCertPath, v))
}

// CertPathContainsFold applies the ContainsFold predicate on the "cert_path" field.
func CertPathContainsFold(v string) predicate.Certification {
	return predicate.Certification(sql.FieldContainsFold(FieldCertPath, v))
}

// ExpiresAtEQ applies the EQ predicate on the "expires_at" field.
func ExpiresAtEQ(v time.Time) predicate.Certification {
	return predicate.Certification(sql.FieldEQ(FieldExpiresAt, v))
}

// ExpiresAtNEQ applies the NEQ predicate on the "expires_at" field.
func ExpiresAtNEQ(v time.Time) predicate.Certification {
	return predicate.Certification(sql.FieldNEQ(FieldExpiresAt, v))
}

// ExpiresAtIn applies the In predicate on the "expires_at" field.
func ExpiresAtIn(vs ...time.Time) predicate.Certification {
	return predicate.Certification(sql.FieldIn(FieldExpiresAt, vs...))
}

// ExpiresAtNotIn applies the NotIn predicate on the "expires_at" field.
func ExpiresAtNotIn(vs ...time.Time) predicate.Certification {
	return predicate.Certification(sql.FieldNotIn(FieldExpiresAt, vs...))
}

// ExpiresAtGT applies the GT predicate on the "expires_at" field.
func ExpiresAtGT(v time.Time) predicate.Certification {
	return predicate.Certification(sql.FieldGT(FieldExpiresAt, v))
}

// ExpiresAtGTE applies the GTE predicate on the "expires_at" field.
func ExpiresAtGTE(v time.Time) predicate.Certification {
	return predicate.Certification(sql.FieldGTE(FieldExpiresAt, v))
}

// ExpiresAtLT applies the LT predicate on the "expires_at" field.
func ExpiresAtLT(v time.Time) predicate.Certification {
	return predicate.Certification(sql.FieldLT(FieldExpiresAt, v))
}

// ExpiresAtLTE applies the LTE predicate on the "expires_at" field.
func ExpiresAtLTE(v time.Time) predicate.Certification {
	return predicate.Certification(sql.FieldLTE(FieldExpiresAt, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Certification) predicate.Certification {
	return predicate.Certification(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Certification) predicate.Certification {
	return predicate.Certification(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Certification) predicate.Certification {
	return predicate.Certification(sql.NotPredicates(p))
}
