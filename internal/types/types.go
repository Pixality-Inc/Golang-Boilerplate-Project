package types

import (
	"database/sql/driver"

	"github.com/pixality-inc/golang-core/json"
	"github.com/pixality-inc/golang-core/postgres"
	"github.com/pixality-inc/golang-core/util"
	uuid "github.com/satori/go.uuid"
)

// BookId

type BookId uuid.UUID // nolint:recvcheck

var emptyBookId = BookId(uuid.Nil)

func ParseBookId(s string) (BookId, error) {
	id, err := uuid.FromString(s)
	if err != nil {
		return emptyBookId, err
	}

	return BookId(id), nil
}

func (v BookId) UUID() uuid.UUID { return uuid.UUID(v) }

func (v BookId) String() string {
	return v.UUID().String()
}

func (v *BookId) Scan(value any) error {
	return postgres.ScanTypedIdUuid(value, func(value uuid.UUID) BookId { return BookId(value) }, &v)
}

func (v BookId) Value() (driver.Value, error) {
	return postgres.RenderSqlDriverValue(v)
}

func (v *BookId) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v *BookId) UnmarshalJSON(data []byte) error {
	value, err := util.UnmarshalJsonToId(data, emptyBookId, func(uuidValue uuid.UUID) BookId {
		return BookId(uuidValue)
	})
	if err != nil {
		return err
	}

	*v = value

	return nil
}
