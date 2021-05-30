package db

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter, sort *racing.ListRacesRequestSortBy) ([]*racing.Race, error)

	// Get will return a race for given id.
	Get(id int64) (*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *racesRepo) List(filter *racing.ListRacesRequestFilter, sort *racing.ListRacesRequestSortBy) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, filter)
	query = r.applySortBy(query, sort)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

func (r *racesRepo) Get(id int64) (*racing.Race, error) {
	var (
		query string
		args  []interface{}
	)

	query = getRaceQueries()[singleRace]
	args = append(args, id)
	fmt.Println(query)
	fmt.Println(args...)

	row := r.db.QueryRow(query, args...)

	return r.scanRace(row)
}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	// if visible_only filter is passed through then return races which are flagged as visible true
	if filter.VisibleOnly {
		clauses = append(clauses, "visible IS TRUE")
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (r *racesRepo) applySortBy(query string, sort *racing.ListRacesRequestSortBy) string {
	var (
		orderByColumn string
		orderByType   string
	)

	if sort == nil {
		return query
	}

	// default sort by advertised_start_time
	orderByColumn = "advertised_start_time"
	// default sort by descending
	orderByType = "DESC"

	// Could be validated to make sure column exists and type is valid, skipping it for now
	if len(sort.Column) > 0 {
		orderByColumn = sort.Column
	}

	if len(sort.Type) > 0 {
		orderByType = sort.Type
	}

	query += fmt.Sprintf(" ORDER BY %s %s", orderByColumn, orderByType)

	return query
}

func (m *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart, &race.Status); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts

		races = append(races, &race)
	}

	return races, nil
}

func (m *racesRepo) scanRace(
	row *sql.Row,
) (*racing.Race, error) {
	var race racing.Race
	var advertisedStart time.Time

	if err := row.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart, &race.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	ts, err := ptypes.TimestampProto(advertisedStart)
	if err != nil {
		return nil, err
	}

	race.AdvertisedStartTime = ts

	return &race, nil
}
