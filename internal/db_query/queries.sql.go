// Code generated by pggen. DO NOT EDIT.

package db_query

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	GetPoktApplications(ctx context.Context, encryptionKey string) ([]GetPoktApplicationsRow, error)

	InsertPoktApplications(ctx context.Context, privateKey string, encryptionKey string) (pgconn.CommandTag, error)

	DeletePoktApplication(ctx context.Context, applicationID pgtype.UUID) (pgconn.CommandTag, error)

	GetChainConfigurations(ctx context.Context) ([]GetChainConfigurationsRow, error)
}

var _ Querier = &DBQuerier{}

type DBQuerier struct {
	conn  genericConn   // underlying Postgres transport to use
	types *typeResolver // resolve types by name
}

// genericConn is a connection like *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
type genericConn interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

// NewQuerier creates a DBQuerier that implements Querier.
func NewQuerier(conn genericConn) *DBQuerier {
	return &DBQuerier{conn: conn, types: newTypeResolver()}
}

// typeResolver looks up the pgtype.ValueTranscoder by Postgres type name.
type typeResolver struct {
	connInfo *pgtype.ConnInfo // types by Postgres type name
}

func newTypeResolver() *typeResolver {
	ci := pgtype.NewConnInfo()
	return &typeResolver{connInfo: ci}
}

// findValue find the OID, and pgtype.ValueTranscoder for a Postgres type name.
func (tr *typeResolver) findValue(name string) (uint32, pgtype.ValueTranscoder, bool) {
	typ, ok := tr.connInfo.DataTypeForName(name)
	if !ok {
		return 0, nil, false
	}
	v := pgtype.NewValue(typ.Value)
	return typ.OID, v.(pgtype.ValueTranscoder), true
}

// setValue sets the value of a ValueTranscoder to a value that should always
// work and panics if it fails.
func (tr *typeResolver) setValue(vt pgtype.ValueTranscoder, val interface{}) pgtype.ValueTranscoder {
	if err := vt.Set(val); err != nil {
		panic(fmt.Sprintf("set ValueTranscoder %T to %+v: %s", vt, val, err))
	}
	return vt
}

const getPoktApplicationsSQL = `SELECT id, pgp_sym_decrypt(encrypted_private_key, $1) AS decrypted_private_key
FROM pokt_applications;`

type GetPoktApplicationsRow struct {
	ID                  pgtype.UUID `json:"id"`
	DecryptedPrivateKey string      `json:"decrypted_private_key"`
}

// GetPoktApplications implements Querier.GetPoktApplications.
func (q *DBQuerier) GetPoktApplications(ctx context.Context, encryptionKey string) ([]GetPoktApplicationsRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GetPoktApplications")
	rows, err := q.conn.Query(ctx, getPoktApplicationsSQL, encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("query GetPoktApplications: %w", err)
	}
	defer rows.Close()
	items := []GetPoktApplicationsRow{}
	for rows.Next() {
		var item GetPoktApplicationsRow
		if err := rows.Scan(&item.ID, &item.DecryptedPrivateKey); err != nil {
			return nil, fmt.Errorf("scan GetPoktApplications row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close GetPoktApplications rows: %w", err)
	}
	return items, err
}

const insertPoktApplicationsSQL = `INSERT INTO pokt_applications (encrypted_private_key)
VALUES (pgp_sym_encrypt($1, $2));`

// InsertPoktApplications implements Querier.InsertPoktApplications.
func (q *DBQuerier) InsertPoktApplications(ctx context.Context, privateKey string, encryptionKey string) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertPoktApplications")
	cmdTag, err := q.conn.Exec(ctx, insertPoktApplicationsSQL, privateKey, encryptionKey)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertPoktApplications: %w", err)
	}
	return cmdTag, err
}

const deletePoktApplicationSQL = `DELETE FROM pokt_applications
WHERE id = $1;`

// DeletePoktApplication implements Querier.DeletePoktApplication.
func (q *DBQuerier) DeletePoktApplication(ctx context.Context, applicationID pgtype.UUID) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeletePoktApplication")
	cmdTag, err := q.conn.Exec(ctx, deletePoktApplicationSQL, applicationID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query DeletePoktApplication: %w", err)
	}
	return cmdTag, err
}

const getChainConfigurationsSQL = `SELECT * FROM chain_configurations;`

type GetChainConfigurationsRow struct {
	CreatedAt                        pgtype.Timestamp `json:"created_at"`
	UpdatedAt                        pgtype.Timestamp `json:"updated_at"`
	DeletedAt                        pgtype.Timestamp `json:"deleted_at"`
	ID                               pgtype.UUID      `json:"id"`
	ChainID                          pgtype.Varchar   `json:"chain_id"`
	PocketRequestTimeoutDuration     pgtype.Varchar   `json:"pocket_request_timeout_duration"`
	AltruistUrl                      pgtype.Varchar   `json:"altruist_url"`
	AltruistRequestTimeoutDuration   pgtype.Varchar   `json:"altruist_request_timeout_duration"`
	TopBucketP90latencyDuration      pgtype.Varchar   `json:"top_bucket_p90latency_duration"`
	HeightCheckBlockTolerance        *int32           `json:"height_check_block_tolerance"`
	DataIntegrityCheckLookbackHeight *int32           `json:"data_integrity_check_lookback_height"`
	FixedHeaders                     *pgtype.JSON     `json:"fixed_headers"`
}

// GetChainConfigurations implements Querier.GetChainConfigurations.
func (q *DBQuerier) GetChainConfigurations(ctx context.Context) ([]GetChainConfigurationsRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GetChainConfigurations")
	rows, err := q.conn.Query(ctx, getChainConfigurationsSQL)
	if err != nil {
		return nil, fmt.Errorf("query GetChainConfigurations: %w", err)
	}
	defer rows.Close()
	items := []GetChainConfigurationsRow{}
	for rows.Next() {
		var item GetChainConfigurationsRow
		if err := rows.Scan(&item.CreatedAt, &item.UpdatedAt, &item.DeletedAt, &item.ID, &item.ChainID, &item.PocketRequestTimeoutDuration, &item.AltruistUrl, &item.AltruistRequestTimeoutDuration, &item.TopBucketP90latencyDuration, &item.HeightCheckBlockTolerance, &item.DataIntegrityCheckLookbackHeight); err != nil {
			return nil, fmt.Errorf("scan GetChainConfigurations row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close GetChainConfigurations rows: %w", err)
	}
	return items, err
}

// textPreferrer wraps a pgtype.ValueTranscoder and sets the preferred encoding
// format to text instead binary (the default). pggen uses the text format
// when the OID is unknownOID because the binary format requires the OID.
// Typically occurs for unregistered types.
type textPreferrer struct {
	pgtype.ValueTranscoder
	typeName string
}

// PreferredParamFormat implements pgtype.ParamFormatPreferrer.
func (t textPreferrer) PreferredParamFormat() int16 { return pgtype.TextFormatCode }

func (t textPreferrer) NewTypeValue() pgtype.Value {
	return textPreferrer{ValueTranscoder: pgtype.NewValue(t.ValueTranscoder).(pgtype.ValueTranscoder), typeName: t.typeName}
}

func (t textPreferrer) TypeName() string {
	return t.typeName
}

// unknownOID means we don't know the OID for a type. This is okay for decoding
// because pgx call DecodeText or DecodeBinary without requiring the OID. For
// encoding parameters, pggen uses textPreferrer if the OID is unknown.
const unknownOID = 0
