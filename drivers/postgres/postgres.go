package postgres

import (
	// DRIVER: postgres
	"github.com/lib/pq"

	"github.com/xo/usql/drivers"
)

func init() {
	drivers.Register("postgres", drivers.Driver{
		N:   "pq",
		AD:  true,
		AMC: true,
		Syn: "postgres",
		V: func(db drivers.DB) (string, error) {
			var ver string
			err := db.QueryRow(`SHOW server_version`).Scan(&ver)
			if err != nil {
				return "", err
			}
			return "PostgreSQL " + ver, nil
		},
		ChPw: func(db drivers.DB, user, new, _ string) error {
			_, err := db.Exec(`ALTER USER ` + user + ` PASSWORD '` + new + `'`)
			return err
		},
		E: func(err error) (string, string) {
			if e, ok := err.(*pq.Error); ok {
				return string(e.Code), e.Message
			}
			return "", err.Error()
		},
		PwErr: func(err error) bool {
			if e, ok := err.(*pq.Error); ok {
				return e.Code.Name() == "invalid_password"
			}
			return false
		},
	}, "cockroachdb")
}
